// Copyright (c) 2019 //SEIBERT/MEDIA GmbH All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package consumer

import (
	"context"

	"github.com/Shopify/sarama"
)

// MessageHandlerFunc allow use a function as MessageHandler.
type MessageHandlerFunc func(ctx context.Context, msg *sarama.ConsumerMessage) error

// ConsumeMessage forward to the function.
func (m MessageHandlerFunc) ConsumeMessage(ctx context.Context, msg *sarama.ConsumerMessage) error {
	return m(ctx, msg)
}
