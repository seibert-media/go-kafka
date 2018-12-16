// Copyright (c) 2018 //SEIBERT/MEDIA GmbH All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package consumer

import "context"

//go:generate counterfeiter -o ../mocks/consumer.go --fake-name Consumer . Consumer

// Consumer is the common interface for all consumers.
type Consumer interface {

	// Consume messages until context is canceled.
	Consume(ctx context.Context) error
}
