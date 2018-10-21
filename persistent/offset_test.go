// Copyright (c) 2018 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package persistent_test

import (
	"strconv"
	"testing"

	"github.com/seibert-media/go-kafka/persistent"
)

func TestOffset(t *testing.T) {
	tests := []struct {
		value int64
	}{
		{
			value: 0,
		},
		{
			value: 123,
		},
		{
			value: 1337,
		},
	}
	for _, test := range tests {
		t.Run(strconv.Itoa(int(test.value)), func(t *testing.T) {
			bytes := persistent.Offset(test.value).Bytes()
			offset := persistent.OffsetFromBytes(bytes)
			if test.value != offset.Int64() {
				t.Fatalf("expected %d got %d", test.value, offset.Int64())
			}
		})
	}
}
