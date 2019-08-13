package persistent

import (
	"context"
	"github.com/Shopify/sarama"
	"github.com/boltdb/bolt"
	"github.com/golang/glog"
	"github.com/pkg/errors"
)

const storageBucket = "storage"
const offsetBucket = "offset"

type Storage interface {
	Get(ctx context.Context, key []byte) ([]byte, error)
	Set(ctx context.Context, key []byte, value []byte) error
	Read(ctx context.Context) error
}

func NewStorage(
	producer sarama.SyncProducer,
	client sarama.Client,
	db *bolt.DB,
	topic string,
) Storage {
	return &storage{
		producer: producer,
		client:   client,
		db:       db,
		topic:    topic,
	}
}

type storage struct {
	producer sarama.SyncProducer
	client   sarama.Client
	db       *bolt.DB
	topic    string
}

// Get last seen value from Bolt.
func (s *storage) Get(ctx context.Context, key []byte) (result []byte, err error) {
	err = s.db.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte(storageBucket))
		if err != nil {
			return errors.Wrapf(err, "create bucket %s failed", storageBucket)
		}
		result = bucket.Get(key)
		return nil
	})
	return
}

// Set send each value to a topic.
func (s *storage) Set(ctx context.Context, key []byte, value []byte) error {
	partition, offset, err := s.producer.SendMessage(&sarama.ProducerMessage{
		Topic: s.topic,
		Key:   sarama.ByteEncoder(key),
		Value: sarama.ByteEncoder(value),
	})
	if err != nil {
		return errors.Wrap(err, "send message failed")
	}
	glog.V(3).Infof("send message successful to %s with partition %d offset %d", s.topic, partition, offset)
	return nil
}

// Read topic and save keys and values to Bolt.
func (s *storage) Read(ctx context.Context) error {
	return NewOffsetConsumer(
		MessageHandlerFunc(func(ctx context.Context, tx *bolt.Tx, msg *sarama.ConsumerMessage) error {
			bucket, err := tx.CreateBucketIfNotExists([]byte(storageBucket))
			if err != nil {
				return errors.Wrapf(err, "create bucket %s failed", storageBucket)
			}
			return bucket.Put(msg.Key, msg.Value)
		}),
		s.client,
		s.topic,
		s.db,
		[]byte(offsetBucket),
	).Consume(ctx)
}
