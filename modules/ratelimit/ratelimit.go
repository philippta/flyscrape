// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package ratelimit

import (
	"math"
	"net/http"
	"time"

	"github.com/philippta/flyscrape"
)

func init() {
	flyscrape.RegisterModule(Module{})
}

type Module struct {
	Rate        int `json:"rate"`
	Concurrency int `json:"concurrency"`

	ticker      *time.Ticker
	ratelimit   chan struct{}
	concurrency chan struct{}
}

func (Module) ModuleInfo() flyscrape.ModuleInfo {
	return flyscrape.ModuleInfo{
		ID:  "ratelimit",
		New: func() flyscrape.Module { return new(Module) },
	}
}

func (m *Module) Provision(v flyscrape.Context) {
	if m.rateLimitEnabled() {
		rate := time.Duration(float64(time.Minute) / float64(m.Rate))
		m.ticker = time.NewTicker(rate)
		m.ratelimit = make(chan struct{}, int(math.Max(float64(m.Rate)/10, 1)))

		go func() {
			m.ratelimit <- struct{}{}
			for range m.ticker.C {
				m.ratelimit <- struct{}{}
			}
		}()
	}

	if m.concurrencyEnabled() {
		m.concurrency = make(chan struct{}, m.Concurrency)
		for i := 0; i < m.Concurrency; i++ {
			m.concurrency <- struct{}{}
		}
	}
}

func (m *Module) AdaptTransport(t http.RoundTripper) http.RoundTripper {
	return flyscrape.RoundTripFunc(func(r *http.Request) (*http.Response, error) {
		if m.rateLimitEnabled() {
			<-m.ratelimit
		}

		if m.concurrencyEnabled() {
			<-m.concurrency
			defer func() { m.concurrency <- struct{}{} }()
		}

		return t.RoundTrip(r)
	})
}

func (m *Module) Finalize() {
	if m.rateLimitEnabled() {
		m.ticker.Stop()
	}
}

func (m *Module) rateLimitEnabled() bool {
	return m.Rate != 0
}

func (m *Module) concurrencyEnabled() bool {
	return m.Concurrency > 0
}

var (
	_ flyscrape.TransportAdapter = (*Module)(nil)
	_ flyscrape.Provisioner      = (*Module)(nil)
	_ flyscrape.Finalizer        = (*Module)(nil)
)
