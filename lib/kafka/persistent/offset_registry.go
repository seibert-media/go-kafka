// Copyright (c) 2018 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package persistent

import (
	"github.com/boltdb/bolt"
	"github.com/pkg/errors"
)

// OffsetRegistry save and load the current offset from Bolt.
type OffsetRegistry interface {
	Get(partition int32) (int64, error)
	Set(partition int32, offset int64) error
}

func NewOffsetRegistry(
	tx *bolt.Tx,
	bucketName []byte,
) OffsetRegistry {
	return &offsetRegistry{
		tx:         tx,
		bucketName: bucketName,
	}
}

type offsetRegistry struct {
	tx         *bolt.Tx
	bucketName []byte
}

// Get offset for the given partition.
func (o *offsetRegistry) Get(partition int32) (int64, error) {
	bucket := o.tx.Bucket(o.bucketName)
	if bucket == nil {
		return 0, errors.New("bucket does not exists")
	}
	bytes := bucket.Get(Partition(partition).Bytes())
	if bytes == nil {
		return 0, errors.New("get offest failed")
	}
	return OffsetFromBytes(bytes).Int64(), nil
}

// Set offset for the given partition.
func (o *offsetRegistry) Set(partition int32, offset int64) error {
	bucket := o.tx.Bucket(o.bucketName)
	if bucket == nil {
		return errors.New("bucket does not exists")
	}
	return bucket.Put(Partition(partition).Bytes(), Offset(offset).Bytes())
}
