// Copyright (c) 2018 //SEIBERT/MEDIA GmbH All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package consumer

import (
	"context"
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
	group string,
) Consumer {
	return &offsetConsumer{
		messageHandler: messageHandler,
		client:         client,
		topic:          topic,
		group:          group,
	}
}

// OffsetConsumer consumes the configured Kafka topic and calls the messagehandler for each message.
type offsetConsumer struct {
	messageHandler MessageHandler
	client         sarama.Client
	topic          string
	group          string
}

// Consume all messages until context is canceled.
func (o *offsetConsumer) Consume(ctx context.Context) error {
	glog.V(0).Infof("consume topic %s started", o.topic)

	consumer, err := sarama.NewConsumerFromClient(o.client)
	if err != nil {
		return errors.Wrapf(err, "create consumer failed")
	}
	defer consumer.Close()

	offsetManager, err := sarama.NewOffsetManagerFromClient(o.group, o.client)
	if err != nil {
		return errors.Wrapf(err, "create offsetManager for group %s failed", o.group)
	}
	defer offsetManager.Close()

	partitions, err := consumer.Partitions(o.topic)
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

			partitionOffsetManager, err := offsetManager.ManagePartition(o.topic, partition)
			if err != nil {
				glog.Warningf("create partitionOffsetManager for topic %s failed: %v", o.topic, err)
				cancel()
				return
			}
			defer partitionOffsetManager.Close()

			nextOffset, _ := partitionOffsetManager.NextOffset()
			glog.V(2).Infof("topic %s with group %s start from offset %d", o.topic, o.group, nextOffset)

			messageHandler := NewMessageHandlerHook(
				nil,
				o.messageHandler,
				MessageHandlerFunc(func(ctx context.Context, msg *sarama.ConsumerMessage) error {
					partitionOffsetManager.MarkOffset(msg.Offset+1, "")
					return nil
				}),
			)
			partitionConsumer := NewPartitionConsumer(messageHandler, o.client, o.topic, partition, nextOffset)
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
