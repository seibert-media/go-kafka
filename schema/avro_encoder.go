// Copyright (c) 2018 //SEIBERT/MEDIA GmbH All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package schema

import (
	"encoding/binary"
)

type AvroEncoder struct {
	SchemaId uint32
	Content  []byte
}

func (a *AvroEncoder) Encode() ([]byte, error) {
	bs := make([]byte, 4)
	binary.BigEndian.PutUint32(bs, a.SchemaId)
	header := append([]byte{0}, bs...)
	return append(header, a.Content...), nil
}

func (a *AvroEncoder) Length() int {
	return 5 + len(a.Content)
}
