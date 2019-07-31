// Copyright (c) 2019 //SEIBERT/MEDIA GmbH All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package consumer_test

import (
	"context"

	"github.com/Shopify/sarama"
	saramamocks "github.com/Shopify/sarama/mocks"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pkg/errors"
	"github.com/seibert-media/go-kafka/consumer"
	"github.com/seibert-media/go-kafka/mocks"
)

var _ = Describe("MessageHandler Metrics", func() {
	var messageHandler consumer.MessageHandler
	var producer *saramamocks.SyncProducer
	var subMessageHandler *mocks.MessageHandler
	BeforeEach(func() {
		producer = saramamocks.NewSyncProducer(GinkgoT(), nil)
		subMessageHandler = &mocks.MessageHandler{}
		messageHandler = consumer.NewRetryMessageHandler(subMessageHandler, producer, 1)
	})
	Context("handle message successful", func() {
		BeforeEach(func() {
			subMessageHandler.ConsumeMessageReturns(nil)
		})
		It("does not produce message on success", func() {
			counter := 0
			producer.ExpectSendMessageWithCheckerFunctionAndSucceed(func(val []byte) error {
				counter++
				return nil
			})
			err := messageHandler.ConsumeMessage(context.Background(), &sarama.ConsumerMessage{
				Key:   []byte("hello"),
				Value: []byte("world"),
			})
			Expect(err).To(BeNil())
			Expect(counter).To(Equal(0))
		})
	})
	Context("handle message fails", func() {
		BeforeEach(func() {
			subMessageHandler.ConsumeMessageReturns(errors.New("banana"))
		})
		It("does not produce message if header is > retry limit", func() {
			counter := 0
			producer.ExpectSendMessageWithCheckerFunctionAndSucceed(func(val []byte) error {
				counter++
				return nil
			})
			err := messageHandler.ConsumeMessage(context.Background(), &sarama.ConsumerMessage{
				Key:   []byte("hello"),
				Value: []byte("world"),
				Headers: []*sarama.RecordHeader{
					{
						Key:   []byte("retry-limit"),
						Value: []byte("1337"),
					},
				},
			})
			Expect(err).NotTo(BeNil())
			Expect(counter).To(Equal(0))
		})
		It("does not produce message if header is == retry limit", func() {
			counter := 0
			producer.ExpectSendMessageWithCheckerFunctionAndSucceed(func(val []byte) error {
				counter++
				return nil
			})
			err := messageHandler.ConsumeMessage(context.Background(), &sarama.ConsumerMessage{
				Key:   []byte("hello"),
				Value: []byte("world"),
				Headers: []*sarama.RecordHeader{
					{
						Key:   []byte("retry-limit"),
						Value: []byte("1"),
					},
				},
			})
			Expect(err).NotTo(BeNil())
			Expect(counter).To(Equal(0))
		})
		It("sends message on error if header is not existing", func() {
			counter := 0
			producer.ExpectSendMessageWithCheckerFunctionAndSucceed(func(val []byte) error {
				counter++
				return nil
			})
			err := messageHandler.ConsumeMessage(context.Background(), &sarama.ConsumerMessage{
				Key:   []byte("hello"),
				Value: []byte("world"),
			})
			Expect(err).To(BeNil())
			Expect(counter).To(Equal(1))
		})
		It("sends message on error if header is 0", func() {
			counter := 0
			producer.ExpectSendMessageWithCheckerFunctionAndSucceed(func(val []byte) error {
				counter++
				return nil
			})
			err := messageHandler.ConsumeMessage(context.Background(), &sarama.ConsumerMessage{
				Key:   []byte("hello"),
				Value: []byte("world"),
				Headers: []*sarama.RecordHeader{
					{
						Key:   []byte("retry-limit"),
						Value: []byte("0"),
					},
				},
			})
			Expect(err).To(BeNil())
			Expect(counter).To(Equal(1))
		})
	})
})
