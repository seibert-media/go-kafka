package persistent_test

import (
	"context"
	"io/ioutil"
	"os"

	"github.com/Shopify/sarama"
	saramamocks "github.com/Shopify/sarama/mocks"
	"github.com/boltdb/bolt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/seibert-media/go-kafka/persistent"
)

var _ = Describe("OffsetRegistry", func() {
	var producer *saramamocks.SyncProducer
	var client sarama.Client
	var db *bolt.DB
	var topic string
	var filename string
	var storage persistent.Storage
	BeforeEach(func() {
		producer = &saramamocks.SyncProducer{}

		file, err := ioutil.TempFile("", "")
		Expect(err).To(BeNil())
		filename = file.Name()
		db, err = bolt.Open(filename, 0600, nil)
		Expect(err).To(BeNil())

		// todo: howto mock?
		client = nil

		topic = "my-topic"

		storage = persistent.NewStorage(
			producer,
			client,
			db,
			topic,
		)
	})
	AfterEach(func() {
		_ = os.Remove(filename)
	})
	It("return no error if value not found", func() {
		value, err := storage.Get(context.Background(), []byte("key"))
		Expect(err).To(BeNil())
		Expect(len(value)).To(Equal(0))
	})
	It("sends message on set value", func() {
		counter := 0
		producer.ExpectSendMessageWithCheckerFunctionAndSucceed(func(val []byte) error {
			counter++
			return nil
		})
		err := storage.Set(context.Background(), []byte("key"), []byte("value"))
		Expect(err).To(BeNil())
		Expect(counter).To(Equal(1))
	})
})
