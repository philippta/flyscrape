// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package flyscrape

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

const HeaderBypassCache = "X-Flyscrape-Bypass-Cache"

func Prettify(v any, prefix string) string {
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetEscapeHTML(false)
	enc.SetIndent(prefix, "  ")
	enc.Encode(v)
	return prefix + strings.TrimSuffix(buf.String(), "\n")
}

type RoundTripFunc func(*http.Request) (*http.Response, error)

func (f RoundTripFunc) RoundTrip(r *http.Request) (*http.Response, error) {
	return f(r)
}

func MockTransport(statusCode int, html string) RoundTripFunc {
	return func(*http.Request) (*http.Response, error) {
		return MockResponse(statusCode, html)
	}
}

func MockResponse(statusCode int, html string) (*http.Response, error) {
	return &http.Response{
		StatusCode: statusCode,
		Status:     fmt.Sprintf("%d %s", statusCode, http.StatusText(statusCode)),
		Body:       io.NopCloser(strings.NewReader(html)),
		Header:     http.Header{"Content-Type": []string{"text/html"}},
	}, nil
}
