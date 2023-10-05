// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package jsonprint

import (
	"fmt"
	"time"

	"github.com/philippta/flyscrape"
)

func init() {
	flyscrape.RegisterModule(Module{})
}

type Module struct {
	once bool
}

func (Module) ModuleInfo() flyscrape.ModuleInfo {
	return flyscrape.ModuleInfo{
		ID:  "jsonprint",
		New: func() flyscrape.Module { return new(Module) },
	}
}

func (m *Module) ReceiveResponse(resp *flyscrape.Response) {
	if resp.Error == nil && resp.Data == nil {
		return
	}

	if !m.once {
		fmt.Println("[")
		m.once = true
	} else {
		fmt.Println(",")
	}

	o := output{
		URL:       resp.Request.URL,
		Data:      resp.Data,
		Error:     resp.Error,
		Timestamp: time.Now(),
	}

	fmt.Print(flyscrape.Prettify(o, "  "))
}

func (m *Module) Finalize() {
	fmt.Println("\n]")
}

type output struct {
	URL       string    `json:"url,omitempty"`
	Data      any       `json:"data,omitempty"`
	Error     error     `json:"error,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
}

var (
	_ flyscrape.ResponseReceiver = (*Module)(nil)
	_ flyscrape.Finalizer        = (*Module)(nil)
)
