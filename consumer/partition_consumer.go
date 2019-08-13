// Copyright (c) 2018 //SEIBERT/MEDIA GmbH All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package consumer

import (
	"context"
	"github.com/Shopify/sarama"
	"github.com/golang/glog"
	"github.com/pkg/errors"
)

// NewPartitionConsumer returns a consumer for the given topic, partition and offset.
func NewPartitionConsumer(
	messageHandler MessageHandler,
	client sarama.Client,
	topic string,
	partition int32,
	offset int64,
) Consumer {
	return &partitionConsumer{
		messageHandler: messageHandler,
		client:         client,
		topic:          topic,
		partition:      partition,
		offset:         offset,
	}
}

// PartitionConsumer consume all new messages in the configured topic and calls messagehandler for each.
type partitionConsumer struct {
	messageHandler MessageHandler
	client         sarama.Client
	topic          string
	partition      int32
	offset         int64
}

// Consume all messages until context is canceled.
func (p *partitionConsumer) Consume(ctx context.Context) error {
	glog.V(3).Infof("consume topic %s partition %d started", p.topic, p.partition)
	defer glog.V(3).Infof("consume topic %s partition %d finished", p.topic, p.partition)

	consumer, err := sarama.NewConsumerFromClient(p.client)
	if err != nil {
		return errors.Wrapf(err, "create consumer failed")
	}
	defer consumer.Close()

	partitionConsumer, err := consumer.ConsumePartition(p.topic, p.partition, p.offset)
	if err != nil {
		return errors.Wrapf(err, "create partitionConsumer for topic %s failed", p.topic)
	}
	defer partitionConsumer.Close()
	for {
		select {
		case <-ctx.Done():
			return nil
		case err := <-partitionConsumer.Errors():
			return err
		case msg, ok := <-partitionConsumer.Messages():
			if !ok {
				glog.V(4).Infof("message channel closed")
				return nil
			}
			if glog.V(4) {
				glog.Infof("handle message: %s", string(msg.Value))
			}
			if err := p.messageHandler.ConsumeMessage(ctx, msg); err != nil {
				glog.V(1).Infof("consume message %d failed: %v", msg.Offset, err)
				continue
			}
			glog.V(3).Infof("message %d consumed from partition %d in topic %s successful", msg.Offset, p.partition, p.topic)
		}
	}
}
