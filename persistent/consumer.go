// Copyright (c) 2018 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package persistent

import (
	"context"
	"strings"
	"sync"

	"github.com/Shopify/sarama"
	"github.com/golang/glog"
	"github.com/pkg/errors"
)

// Consumer for with customer offset manager.
type Consumer struct {
	KafkaBrokers  string
	KafkaTopic    string
	OffsetManager interface {
		NextOffset(partition int32) (int64, error)
		HandleMessage(partition int32, msg *sarama.ConsumerMessage) error
	}
}

// Consume all messages until context is canceled.
func (o *Consumer) Consume(ctx context.Context) error {
	glog.V(3).Infof("import to %s started", o.KafkaTopic)

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

	partitions, err := consumer.Partitions(o.KafkaTopic)
	if err != nil {
		return errors.Wrapf(err, "get partitions for topic %s failed", o.KafkaTopic)
	}
	glog.V(3).Infof("found kafka partitions: %v", partitions)

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	var wg sync.WaitGroup
	for _, partition := range partitions {
		wg.Add(1)
		go func(partition int32) {
			defer wg.Done()

			glog.V(3).Infof("consume topic %s partition %d started", o.KafkaTopic, partition)
			defer glog.V(3).Infof("consume topic %s partition %d finished", o.KafkaTopic, partition)

			offset, err := o.OffsetManager.NextOffset(partition)
			if err != nil {
				glog.Warningf("get offset failed: %v", err)
				cancel()
				return
			}
			glog.V(3).Infof("got offset %d for partition %d", offset, partition)

			partitionConsumer, err := consumer.ConsumePartition(o.KafkaTopic, partition, offset)
			if err != nil {
				cancel()
				glog.Warningf("create partitionConsumer for topic %s failed: %v", o.KafkaTopic, err)
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
					if err := o.OffsetManager.HandleMessage(partition, msg); err != nil {
						glog.Warningf("handle message failed %v", err)
					}
					glog.V(3).Infof("message %d consumed successful", msg.Offset)
				}
			}
		}(partition)
	}
	wg.Wait()
	glog.V(3).Infof("import to %s finish", o.KafkaTopic)
	return nil
}
