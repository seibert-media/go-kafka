// Copyright (c) 2018 //SEIBERT/MEDIA GmbH All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package consumer

import (
	"context"

	"github.com/Shopify/sarama"
)

//go:generate counterfeiter -o ../mocks/messagehandler.go --fake-name MessageHandler . MessageHandler

// MessageHandler is responsible for handling arriving Kakfa messages.
type MessageHandler interface {

	// ConsumeMessage is called for each message.
	ConsumeMessage(ctx context.Context, msg *sarama.ConsumerMessage) error
}
