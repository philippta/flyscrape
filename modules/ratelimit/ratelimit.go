// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package ratelimit

import (
	"time"

	"github.com/philippta/flyscrape"
)

func init() {
	flyscrape.RegisterModule(new(Module))
}

type Module struct {
	Rate float64 `json:"rate"`

	ticker    *time.Ticker
	semaphore chan struct{}
}

func (m *Module) ID() string {
	return "ratelimit"
}

func (m *Module) OnLoad(v flyscrape.Visitor) {
	rate := time.Duration(float64(time.Second) / m.Rate)

	m.ticker = time.NewTicker(rate)
	m.semaphore = make(chan struct{}, 1)

	go func() {
		for range m.ticker.C {
			m.semaphore <- struct{}{}
		}
	}()
}

func (m *Module) OnRequest(_ *flyscrape.Request) {
	<-m.semaphore
}

func (m *Module) OnComplete() {
	m.ticker.Stop()
}

var (
	_ flyscrape.Module     = (*Module)(nil)
	_ flyscrape.OnRequest  = (*Module)(nil)
	_ flyscrape.OnLoad     = (*Module)(nil)
	_ flyscrape.OnComplete = (*Module)(nil)
)
