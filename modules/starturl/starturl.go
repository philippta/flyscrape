// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package starturl

import (
	"github.com/philippta/flyscrape"
)

func init() {
	flyscrape.RegisterModule(Module{})
}

type Module struct {
	URL string `json:"url"`
}

func (Module) ModuleInfo() flyscrape.ModuleInfo {
	return flyscrape.ModuleInfo{
		ID:  "starturl",
		New: func() flyscrape.Module { return new(Module) },
	}
}

func (m *Module) Provision(ctx flyscrape.Context) {
	if m.disabled() {
		return
	}
	ctx.Visit(m.URL)
}

func (m *Module) disabled() bool {
	return m.URL == ""
}

var _ flyscrape.Provisioner = (*Module)(nil)
