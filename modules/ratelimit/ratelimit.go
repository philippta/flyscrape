// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package ratelimit

import (
	"time"

	"github.com/philippta/flyscrape"
)

func init() {
	flyscrape.RegisterModule(Module{})
}

type Module struct {
	Rate float64 `json:"rate"`

	ticker    *time.Ticker
	semaphore chan struct{}
}

func (Module) ModuleInfo() flyscrape.ModuleInfo {
	return flyscrape.ModuleInfo{
		ID:  "ratelimit",
		New: func() flyscrape.Module { return new(Module) },
	}
}

func (m *Module) Provision(v flyscrape.Context) {
	if m.disabled() {
		return
	}

	rate := time.Duration(float64(time.Second) / m.Rate)

	m.ticker = time.NewTicker(rate)
	m.semaphore = make(chan struct{}, 1)

	go func() {
		for range m.ticker.C {
			m.semaphore <- struct{}{}
		}
	}()
}

func (m *Module) BuildRequest(_ *flyscrape.Request) {
	if m.disabled() {
		return
	}
	<-m.semaphore
}

func (m *Module) Finalize() {
	if m.disabled() {
		return
	}
	m.ticker.Stop()
}

func (m *Module) disabled() bool {
	return m.Rate == 0
}

var (
	_ flyscrape.RequestBuilder = (*Module)(nil)
	_ flyscrape.Provisioner    = (*Module)(nil)
	_ flyscrape.Finalizer      = (*Module)(nil)
)
