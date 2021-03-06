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

type RavenClient interface {
	CaptureErrorAndWait(err error, tags map[string]string, interfaces ...raven.Interface) string
}

func SendErrorsToSentry(messageHandler MessageHandler, ravenClient RavenClient) MessageHandler {
	return MessageHandlerFunc(func(ctx context.Context, msg *sarama.ConsumerMessage) error {
		if err := messageHandler.ConsumeMessage(ctx, msg); err != nil {
			data := DataFromError(
				AddDataToError(
					err,
					map[string]string{
						"topic":     msg.Topic,
						"partition": fmt.Sprintf("%d", msg.Partition),
						"offset":    fmt.Sprintf("%d", msg.Offset),
					},
				),
			)

			if glog.V(2) {
				glog.Warningf("consume message %d in topic %s failed: %v %+v", msg.Offset, msg.Topic, err, data)
			} else {
				glog.Warningf("consume message %d in topic %s failed: %v", msg.Offset, msg.Topic, err)
			}

			ravenClient.CaptureErrorAndWait(
				err,
				data,
			)
		}
		return nil
	})
}
