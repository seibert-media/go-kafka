// Copyright (c) 2019 //SEIBERT/MEDIA GmbH All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package consumer

import (
	"context"

	"github.com/Shopify/sarama"
	"github.com/prometheus/client_golang/prometheus"
)

// NewMetricsMessageHandler is a MessageHandler adapter that create Prometheus metrics for started, completed and failed.
func NewMetricsMessageHandler(
	namespace string,
	subsystem string,
	messageHandler MessageHandler,
) MessageHandler {
	started := prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: namespace,
		Subsystem: subsystem,
		Name:      "started",
		Help:      "started",
	})
	completed := prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: namespace,
		Subsystem: subsystem,
		Name:      "completed",
		Help:      "completed",
	})
	failed := prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: namespace,
		Subsystem: subsystem,
		Name:      "failed",
		Help:      "failed",
	})
	prometheus.MustRegister(started, completed, failed)
	return &metricsMessageHandler{
		started:        started,
		completed:      completed,
		failed:         failed,
		messageHandler: messageHandler,
	}
}

type metricsMessageHandler struct {
	messageHandler MessageHandler
	started        prometheus.Gauge
	completed      prometheus.Gauge
	failed         prometheus.Gauge
}

func (m *metricsMessageHandler) ConsumeMessage(ctx context.Context, msg *sarama.ConsumerMessage) error {
	m.started.Inc()
	err := m.messageHandler.ConsumeMessage(ctx, msg)
	if err != nil {
		m.failed.Inc()
		return err
	}
	m.completed.Inc()
	return nil
}
