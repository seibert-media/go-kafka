// Copyright (c) 2018 //SEIBERT/MEDIA GmbH All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package persistent

import (
	"context"
	"github.com/boltdb/bolt"
	"github.com/seibert-media/go-kafka/consumer"
	"sync"

	"github.com/Shopify/sarama"
	"github.com/golang/glog"
	"github.com/pkg/errors"
)

// NewOffsetConsumer return an consumer with offset tracking.
func NewOffsetConsumer(
	messageHandler MessageHandler,
	client sarama.Client,
	topic string,
	db *bolt.DB,
	offsetBucketName []byte,
) consumer.Consumer {
	return &offsetConsumer{
		db:               db,
		offsetBucketName: offsetBucketName,
		messageHandler:   messageHandler,
		client:           client,
		topic:            topic,
	}
}

// OffsetConsumer consumes the configured Kafka topic and calls the messagehandler for each message.
type offsetConsumer struct {
	db               *bolt.DB
	offsetBucketName []byte
	messageHandler   MessageHandler
	client           sarama.Client
	topic            string
	group            string
}

// Consume all messages until context is canceled.
func (o *offsetConsumer) Consume(ctx context.Context) error {
	glog.V(0).Infof("consume topic %s started", o.topic)

	con, err := sarama.NewConsumerFromClient(o.client)
	if err != nil {
		return errors.Wrapf(err, "create consumer failed")
	}
	defer con.Close()

	offsetManager, err := sarama.NewOffsetManagerFromClient(o.group, o.client)
	if err != nil {
		return errors.Wrapf(err, "create offsetManager for group %s failed", o.group)
	}
	defer offsetManager.Close()

	partitions, err := con.Partitions(o.topic)
	if err != nil {
		return errors.Wrapf(err, "get partitions for topic %s failed", o.topic)
	}
	glog.V(2).Infof("found kafka partitions: %v", partitions)

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	var wg sync.WaitGroup
	for _, partition := range partitions {
		wg.Add(1)
		go func(partition int32) {
			defer wg.Done()

			glog.V(1).Infof("consume topic %s partition %d started", o.topic, partition)
			defer glog.V(1).Infof("consume topic %s partition %d finished", o.topic, partition)

			var nextOffset int64
			err := o.db.View(func(tx *bolt.Tx) error {
				nextOffset, err = NewOffsetRegistry(tx, o.offsetBucketName).Get(partition)
				return err
			})
			if err != nil {
				nextOffset = sarama.OffsetOldest
			}

			glog.V(2).Infof("topic %s with group %s start from offset %d", o.topic, o.group, nextOffset)

			partitionConsumer := consumer.NewPartitionConsumer(
				NewOffsetMessageHandler(
					o.db,
					o.offsetBucketName,
					o.messageHandler,
				),
				o.client,
				o.topic,
				partition,
				nextOffset,
			)
			if err := partitionConsumer.Consume(ctx); err != nil {
				glog.Warningf("consumer partition %d for topic %s failed: %v", partition, o.topic, err)
				cancel()
				return
			}
		}(partition)
	}
	wg.Wait()
	glog.V(0).Infof("consume topic %s finish", o.topic)
	return nil
}
