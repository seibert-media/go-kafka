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

type httpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type Registry struct {
	SchemaRegistryUrl string
	HttpClient        httpClient

	mux   sync.Mutex
	cache map[string]uint32
}

func (s *Registry) SchemaId(subject string, object interface {
	Schema() string
}) (uint32, error) {
	s.mux.Lock()
	defer s.mux.Unlock()
	if s.cache == nil {
		s.cache = make(map[string]uint32)
	}
	id, ok := s.cache[object.Schema()]
	if ok {
		glog.V(2).Infof("cache hit return %d", id)
		return id, nil
	}
	input := struct {
		Schema string `json:"schema"`
	}{
		Schema: object.Schema(),
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
		return 0, errors.Wrap(err, "status code != 2xx")
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
	s.cache[object.Schema()] = output.Id
	glog.V(2).Infof("got %d from schema registry", output.Id)
	return output.Id, nil
}
