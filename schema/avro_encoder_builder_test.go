// Copyright (c) 2019 //SEIBERT/MEDIA GmbH All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package schema_test

import (
	"io"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pkg/errors"
	"github.com/seibert-media/go-kafka/mocks"
	"github.com/seibert-media/go-kafka/schema"
)

var avroSerialized = []byte("hello world")
var magicHeaderLength = 5

var _ = Describe("AvroEncoderBuilder", func() {
	var avro *mocks.SchemaAvro
	var schemaRegistry *mocks.SchemaRegistry
	var avroEncoderBuilder schema.AvroEncoderBuilder
	BeforeEach(func() {
		avro = &mocks.SchemaAvro{}
		avro.SchemaReturns("{}")
		avro.SerializeStub = func(writer io.Writer) error {
			writer.Write(avroSerialized)
			return nil
		}
		schemaRegistry = &mocks.SchemaRegistry{}
		avroEncoderBuilder = schema.NewAvroEncoderBuilder(schemaRegistry)
	})
	Context("schema registry returns id", func() {
		BeforeEach(func() {
			schemaRegistry.SchemaIdReturns(123, nil)
		})
		It("returns no error", func() {
			_, err := avroEncoderBuilder.BuildEncoder("mysubject", avro)
			Expect(err).To(BeNil())
		})
		It("length contains magic header + content", func() {
			encoder, err := avroEncoderBuilder.BuildEncoder("mysubject", avro)
			Expect(err).To(BeNil())
			Expect(encoder.Length()).To(Equal(magicHeaderLength + len(avroSerialized)))
		})
		It("length contains magic header + content", func() {
			encoder, err := avroEncoderBuilder.BuildEncoder("mysubject", avro)
			Expect(err).To(BeNil())
			bytes, err := encoder.Encode()
			Expect(err).To(BeNil())
			Expect(len(bytes)).To(Equal(magicHeaderLength + len(avroSerialized)))
		})
		It("encoder contains schemaId", func() {
			encoder, err := avroEncoderBuilder.BuildEncoder("mysubject", avro)
			Expect(err).To(BeNil())
			avroEncoder, ok := encoder.(*schema.AvroEncoder)
			Expect(ok).To(BeTrue())
			Expect(avroEncoder.SchemaId).To(Equal(uint32(123)))
		})
		It("encoder contains content", func() {
			encoder, err := avroEncoderBuilder.BuildEncoder("mysubject", avro)
			Expect(err).To(BeNil())
			avroEncoder, ok := encoder.(*schema.AvroEncoder)
			Expect(ok).To(BeTrue())
			Expect(avroEncoder.Content).To(Equal(avroSerialized))
		})
		Context("avro encode failed", func() {
			BeforeEach(func() {
				avro.SerializeReturns(errors.New("banana"))
			})
			It("returns error", func() {
				_, err := avroEncoderBuilder.BuildEncoder("mysubject", avro)
				Expect(err).NotTo(BeNil())
			})
		})
	})
	Context("schema registry returns error", func() {
		BeforeEach(func() {
			schemaRegistry.SchemaIdReturns(0, errors.New("banana"))
		})
		It("returns error", func() {
			_, err := avroEncoderBuilder.BuildEncoder("mysubject", avro)
			Expect(err).NotTo(BeNil())
		})
	})
})
