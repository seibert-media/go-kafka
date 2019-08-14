// Copyright (c) 2019 //SEIBERT/MEDIA GmbH All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package consumer

import (
	"context"

	"github.com/Shopify/sarama"
	"github.com/golang/glog"
)

func NewMessageHandlerHook(
	preMessageHandler MessageHandler,
	messageHandler MessageHandler,
	postMessageHandler MessageHandler,
) MessageHandler {
	return &messageHandlerHook{
		preMessageHandler:  preMessageHandler,
		messageHandler:     messageHandler,
		postMessageHandler: postMessageHandler,
	}
}

type messageHandlerHook struct {
	preMessageHandler  MessageHandler
	messageHandler     MessageHandler
	postMessageHandler MessageHandler
}

func (m *messageHandlerHook) ConsumeMessage(ctx context.Context, msg *sarama.ConsumerMessage) error {
	if m.preMessageHandler != nil {
		if err := m.preMessageHandler.ConsumeMessage(ctx, msg); err != nil {
			glog.Warningf("post message handler failed: %v", err)
		}
	}
	err := m.messageHandler.ConsumeMessage(ctx, msg)
	if m.preMessageHandler != nil {
		if err := m.postMessageHandler.ConsumeMessage(ctx, msg); err != nil {
			glog.Warningf("post message handler failed: %v", err)
		}
	}
	return err
}
