// Copyright (c) 2018 //SEIBERT/MEDIA GmbH All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package schema

import "net/http"

//go:generate counterfeiter -o ../mocks/http_client.go --fake-name HttpClient . httpClient
type httpClient interface {
	Do(req *http.Request) (*http.Response, error)
}
