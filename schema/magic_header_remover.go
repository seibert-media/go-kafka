// Copyright (c) 2018 //SEIBERT/MEDIA GmbH All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package schema

import (
	"fmt"
	"io"
)

const headerLength = 5

// RemoveMagicHeader reads the five bytes from the reader to remove the magic header.
func RemoveMagicHeader(reader io.Reader) error {
	length, err := reader.Read(make([]byte, headerLength))
	if err != nil {
		return err
	}
	if length != headerLength {
		return fmt.Errorf("read %d bytes from reader failed", headerLength)
	}
	return nil
}
