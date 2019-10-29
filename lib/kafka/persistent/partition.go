// Copyright (c) 2018 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package persistent

import (
	"bytes"
	"encoding/binary"
)

// Partition in Kafka.
type Partition int32

// Int32 value of the partition.
func (o Partition) Int32() int32 {
	return int32(o)
}

// Bytes representation for the partion.
func (o Partition) Bytes() []byte {
	result := &bytes.Buffer{}
	binary.Write(result, binary.BigEndian, o.Int32())
	return result.Bytes()
}

// PartitionFromBytes returns the partition for the given bytes.
func PartitionFromBytes(content []byte) Partition {
	var result int32
	_ = binary.Read(bytes.NewBuffer(content), binary.BigEndian, &result)
	return Partition(result)
}
