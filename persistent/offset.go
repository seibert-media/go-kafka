// Copyright (c) 2018 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package persistent

import (
	"bytes"
	"encoding/binary"
)

type Offset int64

func (o Offset) Int64() int64 {
	return int64(o)
}

func (o Offset) Bytes() []byte {
	result := &bytes.Buffer{}
	binary.Write(result, binary.BigEndian, o.Int64())
	return result.Bytes()
}

func OffsetFromBytes(content []byte) Offset {
	var result int64
	binary.Read(bytes.NewBuffer(content), binary.BigEndian, &result)
	return Offset(result)
}
