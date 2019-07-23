// Copyright (c) 2019 //SEIBERT/MEDIA GmbH All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package schema

import (
	"bytes"
	"io"

	"github.com/Shopify/sarama"
	"github.com/golang/glog"
	"github.com/pkg/errors"
)

//go:generate counterfeiter -o ../mocks/schema_avro.go --fake-name SchemaAvro . Avro

// Avro contains data and schema.
type Avro interface {
	Serialize(w io.Writer) error
	Schema() string
}

//go:generate counterfeiter -o ../mocks/schema_avro_encoder_builder.go --fake-name SchemaAvroEncoderBuilder . AvroEncoderBuilder

// AvroEncoderBuilder create a sarama.Encoder for the given Avro object and subject.
type AvroEncoderBuilder interface {
	BuildEncoder(subject string, avro Avro) (sarama.Encoder, error)
}

type avroEncoderBuilder struct {
	schemaRegistry Registry
}

// NewAvroSchemaEncoderBuilder returns a AvroEncoderBuilder.
func NewAvroEncoderBuilder(
	schemaRegistry Registry,
) AvroEncoderBuilder {
	return &avroEncoderBuilder{
		schemaRegistry: schemaRegistry,
	}
}

// BuildEncoder return a sarama Encoder.
// subject = TOPICNAME-value in for example.
func (a *avroEncoderBuilder) BuildEncoder(subject string, avro Avro) (sarama.Encoder, error) {
	schemaId, err := a.schemaRegistry.SchemaId(subject, avro.Schema())
	if err != nil {
		return nil, errors.Wrap(err, "get schema id failed")
	}
	avroBytes := &bytes.Buffer{}
	if err := avro.Serialize(avroBytes); err != nil {
		return nil, errors.Wrap(err, "encode failed")
	}
	if glog.V(4) {
		glog.Infof("send %s", avroBytes.String())
	}
	return &AvroEncoder{
		SchemaId: schemaId,
		Content:  avroBytes.Bytes(),
	}, nil
}
