// Copyright (c) 2018 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package persistent

import (
	"context"
	"github.com/Shopify/sarama"
	"github.com/boltdb/bolt"
)

//go:generate counterfeiter -o ../mocks/persistent-message-handler.go --fake-name PersistentMessageHandler . PersistentMessageHandler

// MessageHandler handles Kafka messages with an open database connection.
type MessageHandler interface {
	// HandleMessage with a open Bolt transaction.
	ConsumeMessage(ctx context.Context, tx *bolt.Tx, msg *sarama.ConsumerMessage) error
}
