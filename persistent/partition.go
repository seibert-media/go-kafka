// Copyright (c) 2018 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package persistent

import (
	"bytes"
	"encoding/binary"
)

type Partition int32

func (o Partition) Int32() int32 {
	return int32(o)
}

func (o Partition) Bytes() []byte {
	result := &bytes.Buffer{}
	binary.Write(result, binary.BigEndian, o.Int32())
	return result.Bytes()
}

func PartitionFromBytes(content []byte) Partition {
	var result int32
	binary.Read(bytes.NewBuffer(content), binary.BigEndian, &result)
	return Partition(result)
}
