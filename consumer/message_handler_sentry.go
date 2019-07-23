// Copyright (c) 2019 //SEIBERT/MEDIA GmbH All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package consumer

import (
	"context"
	"fmt"

	"github.com/Shopify/sarama"
	"github.com/getsentry/raven-go"
	"github.com/golang/glog"
)

func SendErrorsToSentry(messageHandler MessageHandler) MessageHandler {
	return MessageHandlerFunc(func(ctx context.Context, msg *sarama.ConsumerMessage) error {
		if err := messageHandler.ConsumeMessage(ctx, msg); err != nil {
			glog.Warningf("consume message %d in topic %s failed: %v", msg.Offset, msg.Topic, err)
			raven.CaptureErrorAndWait(
				err,
				DataFromError(
					AddDataToError(
						err,
						map[string]string{
							"topic":     msg.Topic,
							"partition": fmt.Sprintf("%d", msg.Partition),
							"offset":    fmt.Sprintf("%d", msg.Offset),
						},
					),
				),
			)
		}
		return nil
	})
}
