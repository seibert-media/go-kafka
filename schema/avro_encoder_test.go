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
	"reflect"
	"testing"

	"github.com/seibert-media/go-kafka/schema"
)

func TestAvroEncoder(t *testing.T) {
	testcases := []struct {
		name            string
		schemaId        uint32
		content         []byte
		expectedLength  int
		expectedError   error
		expectedContent []byte
	}{
		{
			name:            "simple",
			schemaId:        123,
			content:         []byte("hello"),
			expectedLength:  5 + 5,
			expectedError:   nil,
			expectedContent: append([]byte{0, 0, 0, 0, 123}, []byte("hello")...),
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			encoder := schema.AvroEncoder{
				SchemaId: 123,
				Content:  []byte("hello"),
			}
			length := encoder.Length()
			if tc.expectedLength != length {
				t.Errorf("expect length %d but got %d", tc.expectedLength, length)
			}
			content, err := encoder.Encode()
			if tc.expectedError != err {
				t.Errorf("expect length %d but got %d", tc.expectedLength, length)
			}
			if !reflect.DeepEqual(tc.expectedContent, content) {
				t.Errorf("expected content %v  but got %v", tc.expectedContent, content)
			}
		})
	}
}
