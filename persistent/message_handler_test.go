// Copyright (c) 2018 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package persistent_test

import (
	"io/ioutil"
	"os"

	"github.com/Shopify/sarama"
	"github.com/boltdb/bolt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pkg/errors"
	"github.com/seibert-media/go-kafka/mocks"
	"github.com/seibert-media/go-kafka/persistent"
)

var _ = Describe("MessageHandler", func() {
	offsetBucketName := []byte("offset")
	var filename string
	var db *bolt.DB

	BeforeEach(func() {
		file, err := ioutil.TempFile("", "")
		Expect(err).To(BeNil())
		filename = file.Name()
		db, err = bolt.Open(filename, 0600, nil)
		Expect(err).To(BeNil())

		_ = db.Update(func(tx *bolt.Tx) error {
			_, _ = tx.CreateBucket(offsetBucketName)
			_, _ = tx.CreateBucket([]byte("version"))
			return nil
		})
	})

	AfterEach(func() {
		_ = os.Remove(filename)
	})

	It("return -2 if next offset not found", func() {
		persistentMessageHandler := persistent.MessageHandler{
			DB:               db,
			OffsetBucketName: offsetBucketName,
		}
		offset, err := persistentMessageHandler.NextOffset(0)
		Expect(err).To(BeNil())
		Expect(offset).To(Equal(int64(-2)))
	})

	It("return last offset+1 after the first message", func() {
		persistentMessageHandler := persistent.MessageHandler{
			DB:               db,
			OffsetBucketName: offsetBucketName,
			MessageHandler:   &mocks.PersistentMessageHandler{},
		}
		err := persistentMessageHandler.HandleMessage(0, &sarama.ConsumerMessage{
			Offset: 1000,
			Value:  []byte{},
		})
		Expect(err).To(BeNil())

		nextOffset, err := persistentMessageHandler.NextOffset(0)
		Expect(err).To(BeNil())
		Expect(nextOffset).To(Equal(int64(1001)))
	})

	It("does not change offset if error", func() {
		handler := &mocks.PersistentMessageHandler{}
		handler.HandleMessageReturns(errors.New("banana"))
		persistentMessageHandler := persistent.MessageHandler{
			DB:               db,
			OffsetBucketName: offsetBucketName,
			MessageHandler:   handler,
		}
		err := persistentMessageHandler.HandleMessage(0, &sarama.ConsumerMessage{
			Offset: 1000,
			Value:  []byte{},
		})
		Expect(err).NotTo(BeNil())

		nextOffset, err := persistentMessageHandler.NextOffset(0)
		Expect(err).To(BeNil())
		Expect(nextOffset).To(Equal(int64(-2)))
	})

	It("call subhandler with message", func() {
		messageHandler := &mocks.PersistentMessageHandler{}
		persistentMessageHandler := persistent.MessageHandler{
			DB:               db,
			OffsetBucketName: offsetBucketName,
			MessageHandler:   messageHandler,
		}
		value := []byte("hello world")
		var offset int64 = 1000
		err := persistentMessageHandler.HandleMessage(0, &sarama.ConsumerMessage{
			Offset: offset,
			Value:  value,
		})
		Expect(err).To(BeNil())
		Expect(messageHandler.HandleMessageCallCount()).To(Equal(1))
		tx, message := messageHandler.HandleMessageArgsForCall(0)
		Expect(tx).NotTo(BeNil())
		Expect(message).NotTo(BeNil())
		Expect(message.Value).To(Equal(value))
		Expect(message.Offset).To(Equal(offset))
	})
})
