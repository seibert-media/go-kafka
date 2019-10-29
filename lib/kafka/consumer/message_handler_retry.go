// Copyright (c) 2019 //SEIBERT/MEDIA GmbH All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package consumer

import (
	"context"
	"strconv"

	"github.com/Shopify/sarama"
	"github.com/golang/glog"
	"github.com/pkg/errors"
)

const retryCounterHeaderField = "retry-limit"

func NewRetryMessageHandler(
	messageHandler MessageHandler,
	producer sarama.SyncProducer,
	retries int,
) MessageHandler {
	return &retryMessageHandler{
		messageHandler: messageHandler,
		retryLimit:     retries,
		producer:       producer,
	}
}

type retryMessageHandler struct {
	messageHandler MessageHandler
	retryLimit     int
	producer       sarama.SyncProducer
}

func (r *retryMessageHandler) ConsumeMessage(ctx context.Context, msg *sarama.ConsumerMessage) error {
	if err := r.messageHandler.ConsumeMessage(ctx, msg); err != nil {
		if r.retryLimitReached(msg) {
			glog.V(3).Infof("retry limit reached => return error")
			return errors.Wrap(err, "retry limit reached")
		}
		r.increaseRetryCounter(msg)
		if err := r.reinjectMsg(msg); err != nil {
			return errors.Wrap(err, "reinject message failed")
		}
	}
	return nil
}

func (r *retryMessageHandler) reinjectMsg(msg *sarama.ConsumerMessage) error {
	partition, offset, err := r.producer.SendMessage(r.convertMsg(msg))
	if err != nil {
		return errors.Wrap(err, "reinject msg failed")
	}
	glog.V(3).Infof("send message successful to %s with partition %d offset %d", msg.Topic, partition, offset)
	return nil
}

func (r *retryMessageHandler) convertMsg(msg *sarama.ConsumerMessage) *sarama.ProducerMessage {
	var headers []sarama.RecordHeader
	for _, header := range msg.Headers {
		headers = append(headers, *header)
	}
	return &sarama.ProducerMessage{
		Topic:   msg.Topic,
		Key:     sarama.ByteEncoder(msg.Key),
		Value:   sarama.ByteEncoder(msg.Value),
		Headers: headers,
	}
}

func (r *retryMessageHandler) increaseRetryCounter(msg *sarama.ConsumerMessage) {
	retryCounter := r.retryCounter(msg)
	retryCounter++
	r.setRetryCounter(msg, retryCounter)
}

func (r *retryMessageHandler) retryLimitReached(msg *sarama.ConsumerMessage) bool {
	return r.retryCounter(msg) >= r.retryLimit
}

func (r *retryMessageHandler) retryCounter(msg *sarama.ConsumerMessage) int {
	for _, header := range r.header(msg, retryCounterHeaderField) {
		i, err := strconv.Atoi(string(header.Value))
		if err != nil {
			return 0
		}
		return i
	}
	return 0
}

func (r *retryMessageHandler) header(msg *sarama.ConsumerMessage, field string) []*sarama.RecordHeader {
	var result []*sarama.RecordHeader
	for _, header := range msg.Headers {
		if string(header.Key) == retryCounterHeaderField {
			result = append(result, header)
		}
	}
	return result
}

func (r *retryMessageHandler) setRetryCounter(msg *sarama.ConsumerMessage, retryCounter int) {
	headers := r.header(msg, retryCounterHeaderField)
	if len(headers) == 0 {
		msg.Headers = append(msg.Headers, &sarama.RecordHeader{
			Key:   []byte(retryCounterHeaderField),
			Value: []byte(strconv.Itoa(retryCounter)),
		})
		return
	}
	for _, header := range headers {
		header.Value = []byte(strconv.Itoa(retryCounter))
	}
}
