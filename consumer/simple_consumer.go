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

// NewSimpleConsumer consumes messages in all partitions.
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
			newPartitionConsumer := NewPartitionConsumer(s.messageHandler, s.client, s.topic, partition, sarama.OffsetOldest)
			if err := newPartitionConsumer.Consume(ctx); err != nil {
				glog.Warningf("consumer partition failed: %v", err)
				cancel()
			}
		}(partition)
	}
	wg.Wait()
	glog.V(3).Infof("import to %s finish", s.topic)
	return nil
}
