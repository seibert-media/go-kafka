// Copyright (c) 2018 //SEIBERT/MEDIA GmbH All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package schema_test

import (
	"fmt"
	"net/http"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"
	"github.com/seibert-media/go-kafka/schema"
)

var _ = Describe("SchemaRegistry", func() {
	var server *ghttp.Server
	BeforeEach(func() {
		server = ghttp.NewServer()
		server.RouteToHandler(http.MethodPost, "/subjects/mytopic-value/versions",
			func(resp http.ResponseWriter, req *http.Request) {
				resp.WriteHeader(200)
				fmt.Fprintf(resp, `{"id":42}`)
			},
		)
	})
	AfterEach(func() {
		server.Close()
	})
	It("returns schema id", func() {
		var expectedSchemaId uint32 = 42
		subject := "mytopic-value"
		schemaRegistry := schema.NewRegistry(http.DefaultClient, server.URL())
		schemaId, err := schemaRegistry.SchemaId(subject, "schema-1")
		Expect(err).NotTo(HaveOccurred())
		Expect(schemaId).To(Equal(expectedSchemaId))
	})
})
