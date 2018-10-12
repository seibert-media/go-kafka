// Copyright (c) 2018 //SEIBERT/MEDIA GmbH
//
// Permission is hereby granted, free of charge, to any person obtaining a copy of
// this software and associated documentation files (the "Software"), to deal in
// the Software without restriction, including without limitation the rights to
// use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of
// the Software, and to permit persons to whom the Software is furnished to do so,
// subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
// FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
// COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
// IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
// CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

package schema

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"github.com/golang/glog"
	"github.com/pkg/errors"
)

//go:generate counterfeiter -o ../mocks/http_client.go --fake-name HttpClient . httpClient
type httpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type Registry struct {
	SchemaRegistryUrl string
	HttpClient        httpClient

	mux   sync.Mutex
	cache map[string]uint32
}

// SchemaId return the id for the given schema json
func (s *Registry) SchemaId(subject string, schema string) (uint32, error) {
	s.mux.Lock()
	defer s.mux.Unlock()
	if s.cache == nil {
		s.cache = make(map[string]uint32)
	}
	id, ok := s.cache[schema]
	if ok {
		glog.V(2).Infof("cache hit return %d", id)
		return id, nil
	}
	input := struct {
		Schema string `json:"schema"`
	}{
		Schema: schema,
	}
	body := &bytes.Buffer{}
	if err := json.NewEncoder(body).Encode(input); err != nil {
		return 0, errors.Wrap(err, "encode json failed")
	}
	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/subjects/%s/versions", s.SchemaRegistryUrl, subject), body)
	if err != nil {
		return 0, errors.Wrap(err, "create request failed")
	}
	req.Header.Add("Content-Type", "application/vnd.schemaregistry.v1+json")
	resp, err := s.HttpClient.Do(req)
	if err != nil {
		return 0, errors.Wrap(err, "http request failed")
	}
	if resp.StatusCode/100 != 2 {
		return 0, errors.New("status code != 2xx")
	}
	defer resp.Body.Close()
	var output struct {
		Id uint32 `json:"id"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&output); err != nil {
		return 0, errors.Wrap(err, "decode response into json failed")
	}
	if output.Id == 0 {
		return 0, errors.New("get id from schema registry failed")
	}
	s.cache[schema] = output.Id
	glog.V(2).Infof("got %d from schema registry", output.Id)
	return output.Id, nil
}
