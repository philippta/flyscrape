// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package jsonprinter

import (
	"fmt"

	"github.com/philippta/flyscrape"
)

func init() {
	flyscrape.RegisterModule(new(Module))
}

type Module struct {
	first bool
}

func (m *Module) ID() string {
	return "jsonprinter"
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

	fmt.Print(flyscrape.PrettyPrint(resp.ScrapeResult, "  "))
}

func (m *Module) OnComplete() {
	fmt.Println("\n]")
}

var (
	_ flyscrape.Module     = (*Module)(nil)
	_ flyscrape.OnResponse = (*Module)(nil)
	_ flyscrape.OnComplete = (*Module)(nil)
)
