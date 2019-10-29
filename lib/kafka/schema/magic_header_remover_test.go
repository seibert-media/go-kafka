// Copyright (c) 2018 //SEIBERT/MEDIA GmbH All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package schema_test

import (
	"bytes"
	"io"
	"io/ioutil"
	"testing"

	"github.com/seibert-media/go-kafka/schema"
)

func TestRemoveMagicHeader(t *testing.T) {
	expectedContent := "hello world"
	var reader io.Reader = bytes.NewBuffer(append([]byte{0, 0, 0, 0, 0}, []byte(expectedContent)...))
	err := schema.RemoveMagicHeader(reader)
	if err != nil {
		t.Fatal(err)
	}
	bytes, err := ioutil.ReadAll(reader)
	if err != nil {
		t.Fatal(err)
	}
	content := string(bytes)
	if expectedContent != content {
		t.Fatalf("expected %s but got %s", expectedContent, content)
	}
}

func TestRemoveMagicHeaderReturnError(t *testing.T) {
	var reader io.Reader = bytes.NewBuffer([]byte{0})
	err := schema.RemoveMagicHeader(reader)
	if err == nil {
		t.Fatal("error expected:")
	}
}
