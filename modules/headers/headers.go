// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package headers

import (
	"net/http"

	"github.com/philippta/flyscrape"
)

func init() {
	flyscrape.RegisterModule(Module{})
}

type Module struct {
	Headers map[string]string `json:"headers"`
}

func (Module) ModuleInfo() flyscrape.ModuleInfo {
	return flyscrape.ModuleInfo{
		ID:  "headers",
		New: func() flyscrape.Module { return new(Module) },
	}
}

func (m Module) AdaptTransport(t http.RoundTripper) http.RoundTripper {
	if len(m.Headers) == 0 {
		return t
	}

	return flyscrape.RoundTripFunc(func(r *http.Request) (*http.Response, error) {
		h := http.Header{}
		for k, v := range m.Headers {
			h.Set(k, v)
		}

		r.Header = h

		return t.RoundTrip(r)
	})
}

var _ flyscrape.TransportAdapter = Module{}
