// Copyright (c) 2018 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package persistent

import (
	"github.com/Shopify/sarama"
	"github.com/boltdb/bolt"
	"github.com/pkg/errors"
)

//go:generate counterfeiter -o ../mocks/persistent_message_handler.go --fake-name PersistentMessageHandler . PersistentMessageHandler
type PersistentMessageHandler interface {
	HandleMessage(tx *bolt.Tx, msg *sarama.ConsumerMessage) error
}

type MessageHandler struct {
	DB               *bolt.DB
	OffsetBucketName []byte
	MessageHandler   PersistentMessageHandler
}

func (f *MessageHandler) NextOffset(partition int32) (int64, error) {
	offset := sarama.OffsetOldest
	f.DB.View(func(tx *bolt.Tx) error {
		offsetRegistry := OffsetRegistry{
			BucketName: f.OffsetBucketName,
			Tx:         tx,
		}
		value, err := offsetRegistry.Get(partition)
		if err != nil {
			return err
		}
		offset = value
		return nil
	})
	return offset, nil
}

func (f *MessageHandler) HandleMessage(partition int32, msg *sarama.ConsumerMessage) error {
	return f.DB.Update(func(tx *bolt.Tx) error {
		offsetRegistry := OffsetRegistry{
			BucketName: f.OffsetBucketName,
			Tx:         tx,
		}
		if err := offsetRegistry.Set(partition, msg.Offset+1); err != nil {
			return errors.Wrap(err, "set offest failed")
		}
		return f.MessageHandler.HandleMessage(tx, msg)
	})
}
