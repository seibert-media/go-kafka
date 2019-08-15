// Copyright (c) 2019 //SEIBERT/MEDIA GmbH All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

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

// Storage saves all write to a Kafka topic and caches reads in a local Bolt database.
type Storage interface {
	// Get value for the given key.
	Get(ctx context.Context, key []byte) ([]byte, error)

	// Set a value for the given key.
	Set(ctx context.Context, key []byte, value []byte) error

	// Read fill Bolt db until context is canceled.
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
	err = s.db.View(func(tx *bolt.Tx) error {
		if bucket := tx.Bucket([]byte(storageBucket)); bucket != nil {
			result = bucket.Get(key)
		}
		return nil
	})
	glog.V(2).Infof("get %s=%s to kafka", string(key), string(result))
	return
}

// Set send each value to a topic.
func (s *storage) Set(ctx context.Context, key []byte, value []byte) error {
	glog.V(2).Infof("send %s=%s to kafka", string(key), string(value))
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
	err := s.db.Update(func(tx *bolt.Tx) error {
		if _, err := tx.CreateBucketIfNotExists([]byte(storageBucket)); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return errors.Wrap(err, "create buckets failed")
	}
	return NewOffsetConsumer(
		MessageHandlerFunc(func(ctx context.Context, tx *bolt.Tx, msg *sarama.ConsumerMessage) error {
			bucket := tx.Bucket([]byte(storageBucket))
			if bucket == nil {
				return errors.Errorf("get bucket %s failed", storageBucket)
			}
			glog.V(2).Infof("save %s=%s from kafka to bolt", string(msg.Key), string(msg.Value))
			return bucket.Put(msg.Key, msg.Value)
		}),
		s.client,
		s.topic,
		s.db,
		[]byte(offsetBucket),
	).Consume(ctx)
}
