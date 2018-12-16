// Copyright (c) 2018 //SEIBERT/MEDIA GmbH All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

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

//go:generate counterfeiter -o ../mocks/http_client.go --fake-name HttpClient . HttpClient

// HttpClient for schema registry.
type HttpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// Registry gets the schemaId from the schema registry.
type Registry struct {
	SchemaRegistryUrl string
	HttpClient        HttpClient

	mux   sync.Mutex
	cache map[string]uint32
}

// SchemaId return the id for the given schema json.
func (s *Registry) SchemaId(subject string, schema string) (uint32, error) {
	s.mux.Lock()
	defer s.mux.Unlock()
	if s.cache == nil {
		s.cache = make(map[string]uint32)
	}
	id, ok := s.cache[schema]
	if ok {
		glog.V(4).Infof("cache hit return %d", id)
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
		if glog.V(2) {
			glog.Infof("%s to %s failed with status %d. content: %s", req.Method, req.URL.String(), resp.StatusCode, body.String())
		}
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
	glog.V(3).Infof("got %d from schema registry", output.Id)
	return output.Id, nil
}
