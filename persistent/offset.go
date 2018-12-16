// Copyright (c) 2018 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package persistent

import (
	"bytes"
	"encoding/binary"
)

// Offset in the Kafka topic.
type Offset int64

// Int64 value for the offset.
func (o Offset) Int64() int64 {
	return int64(o)
}

// Bytes representation for the offset.
func (o Offset) Bytes() []byte {
	result := &bytes.Buffer{}
	binary.Write(result, binary.BigEndian, o.Int64())
	return result.Bytes()
}

// OffsetFromBytes returns the offset for the given bytes.
func OffsetFromBytes(content []byte) Offset {
	var result int64
	_ = binary.Read(bytes.NewBuffer(content), binary.BigEndian, &result)
	return Offset(result)
}
