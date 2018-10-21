// Copyright (c) 2018 //SEIBERT/MEDIA GmbH All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package schema_test

import (
	"bytes"
	"io/ioutil"
	"net/http"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/seibert-media/go-kafka/mocks"
	"github.com/seibert-media/go-kafka/schema"
)

var _ = Describe("SchemaRegistry", func() {
	It("returns schema id", func() {
		httpClient := &mocks.HttpClient{}
		httpClient.DoReturns(&http.Response{
			StatusCode: 200,
			Body:       ioutil.NopCloser(bytes.NewBufferString(`{"id":42}`)),
		}, nil)
		var expectedSchemaId uint32 = 42
		subject := "mytopic-value"
		schemaRegistry := schema.Registry{
			SchemaRegistryUrl: "http://schema-registry.example.com",
			HttpClient:        httpClient,
		}
		schemaId, err := schemaRegistry.SchemaId(subject, "schema-1")
		Expect(err).NotTo(HaveOccurred())
		Expect(schemaId).To(Equal(expectedSchemaId))
	})
})
