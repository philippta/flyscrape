// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package cache

import (
	"bufio"
	"bytes"
	"net/http"
	"net/http/httputil"

	"github.com/cornelk/hashmap"
	"github.com/philippta/flyscrape"
)

func init() {
	flyscrape.RegisterModule(Module{})
}

type Module struct {
	Cache string `json:"cache"`

	cache *hashmap.Map[string, []byte]
}

func (Module) ModuleInfo() flyscrape.ModuleInfo {
	return flyscrape.ModuleInfo{
		ID:  "cache",
		New: func() flyscrape.Module { return new(Module) },
	}
}

func (m *Module) Provision(flyscrape.Context) {
	if m.disabled() {
		return
	}
	if m.cache == nil {
		m.cache = hashmap.New[string, []byte]()
	}
}

func (m *Module) AdaptTransport(t http.RoundTripper) http.RoundTripper {
	if m.disabled() {
		return t
	}

	return flyscrape.RoundTripFunc(func(r *http.Request) (*http.Response, error) {
		key := cacheKey(r)

		if b, ok := m.cache.Get(key); ok {
			if resp, err := http.ReadResponse(bufio.NewReader(bytes.NewReader(b)), r); err == nil {
				return resp, nil
			}
		}

		resp, err := t.RoundTrip(r)
		if err != nil {
			return resp, err
		}

		encoded, err := httputil.DumpResponse(resp, true)
		if err != nil {
			return resp, err
		}

		m.cache.Set(key, encoded)
		return resp, nil
	})
}

func (m *Module) disabled() bool {
	return m.Cache == ""
}

func cacheKey(r *http.Request) string {
	return r.Method + " " + r.URL.String()
}
