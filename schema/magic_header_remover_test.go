// Copyright (c) 2018 //SEIBERT/MEDIA GmbH
//
// Permission is hereby granted, free of charge, to any person obtaining a copy of
// this software and associated documentation files (the "Software"), to deal in
// the Software without restriction, including without limitation the rights to
// use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of
// the Software, and to permit persons to whom the Software is furnished to do so,
// subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
// FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
// COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
// IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
// CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

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
