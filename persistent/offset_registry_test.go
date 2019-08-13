// Copyright (c) 2018 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package persistent_test

import (
	"io/ioutil"
	"os"

	"github.com/boltdb/bolt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/seibert-media/go-kafka/persistent"
)

var _ = Describe("OffsetRegistry", func() {
	var filename string
	var db *bolt.DB
	bucketName := []byte("offset")
	BeforeEach(func() {
		file, err := ioutil.TempFile("", "")
		Expect(err).To(BeNil())
		filename = file.Name()
		db, err = bolt.Open(filename, 0600, nil)
		Expect(err).To(BeNil())

		_ = db.Update(func(tx *bolt.Tx) error {
			_, _ = tx.CreateBucket(bucketName)
			return nil
		})

	})
	AfterEach(func() {
		_ = os.Remove(filename)
	})
	It("return error if not exits", func() {
		_ = db.View(func(tx *bolt.Tx) error {
			offsetRegistry := persistent.NewOffsetRegistry(
				tx,
				bucketName,
			)
			_, err := offsetRegistry.Get(1)
			Expect(err).NotTo(BeNil())
			return nil
		})
	})
	It("saves offset", func() {
		offset := int64(42)
		_ = db.Update(func(tx *bolt.Tx) error {
			offsetRegistry := persistent.NewOffsetRegistry(
				tx,
				bucketName,
			)
			_ = offsetRegistry.Set(2, offset)
			return nil
		})
		_ = db.View(func(tx *bolt.Tx) error {
			offsetRegistry := persistent.NewOffsetRegistry(
				tx,
				bucketName,
			)
			offset, err := offsetRegistry.Get(2)
			Expect(err).To(BeNil())
			Expect(offset).To(Equal(offset))
			return nil
		})
	})
})
