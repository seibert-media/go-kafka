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

type OffsetConsumer struct {
	MessageHandler MessageHandler
	KafkaBrokers   string
	KafkaTopic     string
	KafkaGroup     string
}

func (o *OffsetConsumer) Consume(ctx context.Context) error {
	glog.V(0).Infof("import to %s started", o.KafkaTopic)

	config := sarama.NewConfig()
	config.Version = sarama.V2_0_0_0
	config.Consumer.Offsets.Initial = sarama.OffsetOldest
	config.Consumer.Return.Errors = true

	client, err := sarama.NewClient(strings.Split(o.KafkaBrokers, ","), config)
	if err != nil {
		return errors.Wrapf(err, "create kafka client with brokers %s failed", o.KafkaBrokers)
	}
	defer client.Close()

	consumer, err := sarama.NewConsumerFromClient(client)
	if err != nil {
		return errors.Wrapf(err, "create consumer with brokers %s failed", o.KafkaBrokers)
	}
	defer consumer.Close()

	offsetManager, err := sarama.NewOffsetManagerFromClient(o.KafkaGroup, client)
	if err != nil {
		return errors.Wrapf(err, "create offsetManager for group %s failed", o.KafkaGroup)
	}
	defer offsetManager.Close()

	partitions, err := consumer.Partitions(o.KafkaTopic)
	if err != nil {
		return errors.Wrapf(err, "get partitions for topic %s failed", o.KafkaTopic)
	}
	glog.V(2).Infof("found kafka partitions: %v", partitions)

	ctx, cancel := context.WithCancel(ctx)

	var wg sync.WaitGroup
	for _, partition := range partitions {
		wg.Add(1)
		go func(partition int32) {
			defer wg.Done()

			glog.V(1).Infof("consume topic %s partition %d started", o.KafkaTopic, partition)
			defer glog.V(1).Infof("consume topic %s partition %d finished", o.KafkaTopic, partition)

			partitionOffsetManager, err := offsetManager.ManagePartition(o.KafkaTopic, partition)
			if err != nil {
				glog.Warningf("create partitionOffsetManager for topic %s failed: %v", o.KafkaTopic, err)
				cancel()
				return
			}
			defer partitionOffsetManager.Close()

			nextOffset, metadata := partitionOffsetManager.NextOffset()
			glog.V(2).Infof("offset: %d %s", nextOffset, metadata)

			partitionConsumer, err := consumer.ConsumePartition(o.KafkaTopic, partition, nextOffset)
			if err != nil {
				glog.Warningf("create partitionConsumer for topic %s failed: %v", o.KafkaTopic, err)
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
					if err := o.MessageHandler.ConsumeMessage(ctx, msg); err != nil {
						glog.V(1).Infof("consume message %d failed: %v", msg.Offset, err)
						continue
					}
					partitionOffsetManager.MarkOffset(msg.Offset+1, "")
					glog.V(2).Infof("message %d consumed successful", msg.Offset)
				}
			}
		}(partition)
	}
	wg.Wait()
	glog.V(0).Infof("import to %s finish", o.KafkaTopic)
	return nil
}
