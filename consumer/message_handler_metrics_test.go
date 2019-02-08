// Copyright (c) 2019 //SEIBERT/MEDIA GmbH All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package consumer_test

import (
	"context"

	"github.com/Shopify/sarama"
	"github.com/seibert-media/go-kafka/mocks"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/seibert-media/go-kafka/consumer"
)

var _ = Describe("MessageHandler Metrics", func() {
	It("calls sub message handler", func() {
		subMessageHandler := &mocks.MessageHandler{}
		messageHandler := consumer.NewMetricsMessageHandler("a", "b", subMessageHandler)
		err := messageHandler.ConsumeMessage(context.Background(), &sarama.ConsumerMessage{
			Key:   []byte("hello"),
			Value: []byte("world"),
		})
		Expect(err).To(BeNil())
		Expect(subMessageHandler.ConsumeMessageCallCount()).To(Equal(1))
		ctx, message := subMessageHandler.ConsumeMessageArgsForCall(0)
		Expect(ctx).NotTo(BeNil())
		Expect(string(message.Key)).To(Equal("hello"))
		Expect(string(message.Value)).To(Equal("world"))
	})
})
