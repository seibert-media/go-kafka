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

// NewSimpleConsumer return an basic Consumer.
func NewSimpleConsumer(
	messageHandler MessageHandler,
	client sarama.Client,
	topic string,
) Consumer {
	return &simpleConsumer{
		messageHandler: messageHandler,
		client:         client,
		topic:          topic,
	}
}

// SimpleConsumer consume all new messages in the configured topic and calls messagehandler for each.
type simpleConsumer struct {
	messageHandler MessageHandler
	client         sarama.Client
	topic          string
}

// Consume all messages until context is canceled.
func (s *simpleConsumer) Consume(ctx context.Context) error {
	glog.V(3).Infof("import to %s started", s.topic)

	consumer, err := sarama.NewConsumerFromClient(s.client)
	if err != nil {
		return errors.Wrapf(err, "create consumer failed")
	}
	defer consumer.Close()

	partitions, err := consumer.Partitions(s.topic)
	if err != nil {
		return errors.Wrapf(err, "get partitions for topic %s failed", s.topic)
	}
	glog.V(3).Infof("found kafka partitions: %v", partitions)

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	var wg sync.WaitGroup
	for _, partition := range partitions {
		wg.Add(1)
		go func(partition int32) {
			defer wg.Done()

			glog.V(3).Infof("consume topic %s partition %d started", s.topic, partition)
			defer glog.V(3).Infof("consume topic %s partition %d finished", s.topic, partition)

			partitionConsumer, err := consumer.ConsumePartition(s.topic, partition, sarama.OffsetOldest)
			if err != nil {
				glog.Warningf("create partitionConsumer for topic %s failed: %v", s.topic, err)
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
					if err := s.messageHandler.ConsumeMessage(ctx, msg); err != nil {
						glog.V(1).Infof("consume message %d failed: %v", msg.Offset, err)
						continue
					}
					glog.V(3).Infof("message %d consumed from partition %d in topic %s successful", msg.Offset, partition, s.topic)
				}
			}
		}(partition)
	}
	wg.Wait()
	glog.V(3).Infof("import to %s finish", s.topic)
	return nil
}
