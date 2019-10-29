// Copyright (c) 2019 //SEIBERT/MEDIA GmbH All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package consumer_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pkg/errors"
	"github.com/seibert-media/go-kafka/consumer"
)

var _ = Describe("DataError", func() {
	It("returns nothing for standard error", func() {
		err := errors.New("banana")
		data := consumer.DataFromError(err)
		Expect(data).To(HaveLen(0))
		Expect(err.Error()).To(Equal("banana"))
	})
	It("returns data for DataError", func() {
		err := consumer.AddDataToError(errors.New("banana"), map[string]string{"hello": "world"})
		data := consumer.DataFromError(err)
		Expect(data).To(HaveLen(1))
		Expect(data).To(HaveKeyWithValue("hello", "world"))
		Expect(err.Error()).To(Equal("banana"))
	})
	It("returns data if DataError is wrapped", func() {
		err := errors.Wrap(consumer.AddDataToError(errors.New("banana"), map[string]string{"hello": "world"}), "foo bar")
		data := consumer.DataFromError(err)
		Expect(data).To(HaveLen(1))
		Expect(data).To(HaveKeyWithValue("hello", "world"))
		Expect(err.Error()).To(Equal("foo bar: banana"))
	})
	It("combines data", func() {
		err := consumer.AddDataToError(errors.New("banana"), map[string]string{"hello": "world"})
		err = consumer.AddDataToError(err, map[string]string{"hallo": "welt"})
		data := consumer.DataFromError(err)
		Expect(data).To(HaveLen(2))
		Expect(data).To(HaveKeyWithValue("hello", "world"))
		Expect(data).To(HaveKeyWithValue("hallo", "welt"))
		Expect(err.Error()).To(Equal("banana"))
	})
})
