// Copyright (c) 2018 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package persistent

import (
	"github.com/boltdb/bolt"
	"github.com/pkg/errors"
)

// OffsetRegistry save and load the current offset from Bolt.
type OffsetRegistry struct {
	Tx         *bolt.Tx
	BucketName []byte
}

// Get offset for the given partition.
func (o *OffsetRegistry) Get(partition int32) (int64, error) {
	bucket := o.Tx.Bucket(o.BucketName)
	bytes := bucket.Get(Partition(partition).Bytes())
	if bytes == nil {
		return 0, errors.New("get offest failed")
	}
	return OffsetFromBytes(bytes).Int64(), nil
}

// Set offset for the given partition.
func (o *OffsetRegistry) Set(partition int32, offset int64) error {
	offsetBucket := o.Tx.Bucket(o.BucketName)
	return offsetBucket.Put(Partition(partition).Bytes(), Offset(offset).Bytes())
}
