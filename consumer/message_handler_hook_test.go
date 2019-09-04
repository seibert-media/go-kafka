// Copyright (c) 2019 //SEIBERT/MEDIA GmbH All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package consumer_test

import (
	"context"

	"github.com/Shopify/sarama"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/seibert-media/go-kafka/consumer"
	"github.com/seibert-media/go-kafka/mocks"
)

var _ = Describe("MessageHandlerHook", func() {
	var messageHandler consumer.MessageHandler
	var subMessageHandler *mocks.MessageHandler
	var preMessageHandler *mocks.MessageHandler
	var postMessageHandler *mocks.MessageHandler
	var msg *sarama.ConsumerMessage
	var err error
	BeforeEach(func() {
		msg = &sarama.ConsumerMessage{}
		preMessageHandler = &mocks.MessageHandler{}
		subMessageHandler = &mocks.MessageHandler{}
		postMessageHandler = &mocks.MessageHandler{}
	})
	Context("with all", func() {
		BeforeEach(func() {
			messageHandler = consumer.NewMessageHandlerHook(preMessageHandler, subMessageHandler, postMessageHandler)
			err = messageHandler.ConsumeMessage(context.Background(), msg)
		})
		It("return no error", func() {
			Expect(err).To(BeNil())
		})
		It("calls preMessageHandler", func() {
			Expect(preMessageHandler.ConsumeMessageCallCount()).To(Equal(1))
		})
		It("calls subMessageHandler", func() {
			Expect(subMessageHandler.ConsumeMessageCallCount()).To(Equal(1))
		})
		It("calls postMessageHandler", func() {
			Expect(postMessageHandler.ConsumeMessageCallCount()).To(Equal(1))
		})
	})
	Context("with no pre and no post", func() {
		BeforeEach(func() {
			messageHandler = consumer.NewMessageHandlerHook(nil, subMessageHandler, nil)
			err = messageHandler.ConsumeMessage(context.Background(), msg)
		})
		It("return no error", func() {
			Expect(err).To(BeNil())
		})
		It("calls preMessageHandler", func() {
			Expect(preMessageHandler.ConsumeMessageCallCount()).To(Equal(0))
		})
		It("calls subMessageHandler", func() {
			Expect(subMessageHandler.ConsumeMessageCallCount()).To(Equal(1))
		})
		It("calls postMessageHandler not", func() {
			Expect(postMessageHandler.ConsumeMessageCallCount()).To(Equal(0))
		})
	})
	Context("with pre and no post", func() {
		BeforeEach(func() {
			messageHandler = consumer.NewMessageHandlerHook(preMessageHandler, subMessageHandler, nil)
			err = messageHandler.ConsumeMessage(context.Background(), msg)
		})
		It("return no error", func() {
			Expect(err).To(BeNil())
		})
		It("calls preMessageHandler", func() {
			Expect(preMessageHandler.ConsumeMessageCallCount()).To(Equal(1))
		})
		It("calls subMessageHandler", func() {
			Expect(subMessageHandler.ConsumeMessageCallCount()).To(Equal(1))
		})
		It("calls postMessageHandler not", func() {
			Expect(postMessageHandler.ConsumeMessageCallCount()).To(Equal(0))
		})
	})
	Context("with no pre and post", func() {
		BeforeEach(func() {
			messageHandler = consumer.NewMessageHandlerHook(nil, subMessageHandler, postMessageHandler)
			err = messageHandler.ConsumeMessage(context.Background(), msg)
		})
		It("return no error", func() {
			Expect(err).To(BeNil())
		})
		It("calls preMessageHandler not", func() {
			Expect(preMessageHandler.ConsumeMessageCallCount()).To(Equal(0))
		})
		It("calls subMessageHandler", func() {
			Expect(subMessageHandler.ConsumeMessageCallCount()).To(Equal(1))
		})
		It("calls postMessageHandler", func() {
			Expect(postMessageHandler.ConsumeMessageCallCount()).To(Equal(1))
		})
	})
})
