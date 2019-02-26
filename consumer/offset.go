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
	glog.V(0).Infof("import to %s started", o.topic)

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

			nextOffset, metadata := partitionOffsetManager.NextOffset()
			glog.V(2).Infof("offset: %d %s", nextOffset, metadata)

			partitionConsumer, err := consumer.ConsumePartition(o.topic, partition, nextOffset)
			if err != nil {
				glog.Warningf("create partitionConsumer for topic %s failed: %v", o.topic, err)
				cancel()
				return
			}
			defer partitionConsumer.Close()
			for {
				select {
				case <-ctx.Done():
					return
				case err := <-partitionConsumer.Errors():
					glog.Warningf("get error %v", err)
					cancel()
					return
				case msg := <-partitionConsumer.Messages():
					if glog.V(4) {
						glog.Infof("handle message: %s", string(msg.Value))
					}
					if err := o.messageHandler.ConsumeMessage(ctx, msg); err != nil {
						glog.V(1).Infof("consume message %d failed: %v", msg.Offset, err)
						continue
					}
					partitionOffsetManager.MarkOffset(msg.Offset+1, "")
					glog.V(3).Infof("message %d consumed from partition %d in topic %s successful", msg.Offset, partition, o.topic)
				}
			}
		}(partition)
	}
	wg.Wait()
	glog.V(0).Infof("import to %s finish", o.topic)
	return nil
}
