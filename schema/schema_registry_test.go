package schema_test

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/seibert-media/go-kafka/mocks"
	"github.com/seibert-media/go-kafka/schema"
)

func TestRegistryReturnSchemaId(t *testing.T) {
	httpClient := &mocks.HttpClient{}
	httpClient.DoReturns(&http.Response{
		StatusCode: 200,
		Body:       ioutil.NopCloser(bytes.NewBufferString(`{"id":42}`)),
	}, nil)
	object := &mocks.HasSchema{}
	object.SchemaReturns("schema-1")
	var expectedSchemaId uint32 = 42
	subject := "mytopic-value"
	schemaRegistry := schema.Registry{
		SchemaRegistryUrl: "http://schema-registry.example.com",
		HttpClient:        httpClient,
	}
	schemaId, err := schemaRegistry.SchemaId(subject, object)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if expectedSchemaId != schemaId {
		t.Fatalf("expected schemaId %d but got %d", expectedSchemaId, schemaId)
	}
}
