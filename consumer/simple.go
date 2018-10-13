// Copyright (c) 2018 //SEIBERT/MEDIA GmbH All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package consumer

import (
	"context"
	"strings"
	"sync"

	"github.com/Shopify/sarama"
	"github.com/golang/glog"
	"github.com/pkg/errors"
)

type SimpleConsumer struct {
	MessageHandler MessageHandler
	KafkaBrokers   string
	KafkaTopic     string
}

func (s *SimpleConsumer) Consume(ctx context.Context) error {
	glog.V(0).Infof("import to %s started", s.KafkaTopic)

	config := sarama.NewConfig()
	config.Version = sarama.V2_0_0_0
	config.Consumer.Offsets.Initial = sarama.OffsetOldest
	config.Consumer.Return.Errors = true

	client, err := sarama.NewClient(strings.Split(s.KafkaBrokers, ","), config)
	if err != nil {
		return errors.Wrapf(err, "create kafka client with brokers %s failed", s.KafkaBrokers)
	}
	defer client.Close()

	consumer, err := sarama.NewConsumerFromClient(client)
	if err != nil {
		return errors.Wrapf(err, "create consumer with brokers %s failed", s.KafkaBrokers)
	}
	defer consumer.Close()

	partitions, err := consumer.Partitions(s.KafkaTopic)
	if err != nil {
		return errors.Wrapf(err, "get partitions for topic %s failed", s.KafkaTopic)
	}
	glog.V(2).Infof("found kafka partitions: %v", partitions)

	ctx, cancel := context.WithCancel(ctx)

	var wg sync.WaitGroup
	for _, partition := range partitions {
		wg.Add(1)
		go func(partition int32) {
			defer wg.Done()

			glog.V(1).Infof("consume topic %s partition %d started", s.KafkaTopic, partition)
			defer glog.V(1).Infof("consume topic %s partition %d finished", s.KafkaTopic, partition)

			partitionConsumer, err := consumer.ConsumePartition(s.KafkaTopic, partition, sarama.OffsetOldest)
			if err != nil {
				glog.Warningf("create partitionConsumer for topic %s failed: %v", s.KafkaTopic, err)
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
					if err := s.MessageHandler.ConsumeMessage(ctx, msg); err != nil {
						glog.V(1).Infof("consume message %d failed: %v", msg.Offset, err)
						continue
					}
					glog.V(2).Infof("message %d consumed successful", msg.Offset)
				}
			}
		}(partition)
	}
	wg.Wait()
	glog.V(0).Infof("import to %s finish", s.KafkaTopic)
	return nil
}
