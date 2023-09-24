// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package jsonprinter

import (
	"fmt"
	"time"

	"github.com/philippta/flyscrape"
)

func init() {
	flyscrape.RegisterModule(new(Module))
}

type Module struct {
	first bool
}

func (m *Module) OnResponse(resp *flyscrape.Response) {
	if resp.Error == nil && resp.Data == nil {
		return
	}

	if m.first {
		fmt.Println("[")
	} else {
		fmt.Println(",")
	}

	o := output{
		URL:       resp.Request.URL,
		Data:      resp.Data,
		Error:     resp.Error,
		Timestamp: time.Now(),
	}

	fmt.Print(flyscrape.PrettyPrint(o, "  "))
}

func (m *Module) OnComplete() {
	fmt.Println("\n]")
}

type output struct {
	URL       string    `json:"url,omitempty"`
	Data      any       `json:"data,omitempty"`
	Error     error     `json:"error,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
}

var (
	_ flyscrape.OnResponse = (*Module)(nil)
	_ flyscrape.OnComplete = (*Module)(nil)
)
