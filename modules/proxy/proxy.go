// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package proxy

import (
	"crypto/tls"
	"math/rand"
	"net/http"
	"net/url"

	"github.com/philippta/flyscrape"
)

func init() {
	flyscrape.RegisterModule(Module{})
}

type Module struct {
	Proxies []string `json:"proxies"`
	Proxy   string   `json:"proxy"`

	transports []*http.Transport
}

func (Module) ModuleInfo() flyscrape.ModuleInfo {
	return flyscrape.ModuleInfo{
		ID:  "proxy",
		New: func() flyscrape.Module { return new(Module) },
	}
}

func (m *Module) Provision(ctx flyscrape.Context) {
	if m.disabled() {
		return
	}

	for _, purl := range append(m.Proxies, m.Proxy) {
		if parsed, err := url.Parse(purl); err == nil {
			m.transports = append(m.transports, &http.Transport{
				Proxy:           http.ProxyURL(parsed),
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			})
		}

	}
}

func (m *Module) AdaptTransport(t http.RoundTripper) http.RoundTripper {
	if m.disabled() {
		return t
	}

	return flyscrape.RoundTripFunc(func(r *http.Request) (*http.Response, error) {
		transport := m.transports[rand.Intn(len(m.transports))]
		return transport.RoundTrip(r)
	})
}

func (m *Module) disabled() bool {
	return len(m.Proxies) == 0
}
