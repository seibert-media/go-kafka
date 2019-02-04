// Copyright (c) 2019 //SEIBERT/MEDIA GmbH All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package consumer_test

import (
	"context"

	"github.com/Shopify/sarama"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pkg/errors"
	"github.com/seibert-media/go-kafka/consumer"
)

var _ = Describe("MessageHandler", func() {
	It("return nil error", func() {
		var messageHandler consumer.MessageHandler = consumer.MessageHandlerFunc(func(ctx context.Context, msg *sarama.ConsumerMessage) error {
			return nil
		})
		err := messageHandler.ConsumeMessage(context.Background(), &sarama.ConsumerMessage{})
		Expect(err).To(BeNil())
	})
	It("return error", func() {
		var messageHandler consumer.MessageHandler = consumer.MessageHandlerFunc(func(ctx context.Context, msg *sarama.ConsumerMessage) error {
			return errors.New("banana")
		})
		err := messageHandler.ConsumeMessage(context.Background(), &sarama.ConsumerMessage{})
		Expect(err).NotTo(BeNil())
	})
	It("forward args", func() {
		var messageHandler consumer.MessageHandler = consumer.MessageHandlerFunc(func(ctx context.Context, msg *sarama.ConsumerMessage) error {
			Expect(ctx).NotTo(BeNil())
			Expect(msg).NotTo(BeNil())
			return nil
		})
		err := messageHandler.ConsumeMessage(context.Background(), &sarama.ConsumerMessage{})
		Expect(err).To(BeNil())
	})
})
