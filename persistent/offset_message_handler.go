// Copyright (c) 2018 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package persistent

import (
	"context"
	"github.com/Shopify/sarama"
	"github.com/boltdb/bolt"
	"github.com/pkg/errors"
	"github.com/seibert-media/go-kafka/consumer"
)

func NewOffsetMessageHandler(
	db *bolt.DB,
	offsetBucketName []byte,
	messageHandler MessageHandler,
) consumer.MessageHandler {
	return &persitentMessageHandler{
		db:               db,
		offsetBucketName: offsetBucketName,
		messageHandler:   messageHandler,
	}
}

// MessageHandler opens for each Kakfa message a Bolt transaction witch is used for processing the message and storing the offset. So both will be saved or none.
type persitentMessageHandler struct {
	db               *bolt.DB
	offsetBucketName []byte
	messageHandler   MessageHandler
}

// NextOffset returns the offset to start consume from.
func (f *persitentMessageHandler) NextOffset(partition int32) (int64, error) {
	offset := sarama.OffsetOldest
	_ = f.db.View(func(tx *bolt.Tx) error {
		offsetRegistry := NewOffsetRegistry(
			tx,
			f.offsetBucketName,
		)
		value, err := offsetRegistry.Get(partition)
		if err != nil {
			return err
		}
		offset = value
		return nil
	})
	return offset, nil
}

// HandleMessage updates offset and call messagehandler for each message.
func (f *persitentMessageHandler) ConsumeMessage(ctx context.Context, msg *sarama.ConsumerMessage) error {
	return f.db.Update(func(tx *bolt.Tx) error {
		offsetRegistry := NewOffsetRegistry(
			tx,
			f.offsetBucketName,
		)
		if err := offsetRegistry.Set(msg.Partition, msg.Offset+1); err != nil {
			return errors.Wrap(err, "set offest failed")
		}
		return f.messageHandler.ConsumeMessage(ctx, tx, msg)
	})
}
