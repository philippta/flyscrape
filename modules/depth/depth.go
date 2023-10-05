// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package depth

import (
	"github.com/philippta/flyscrape"
)

func init() {
	flyscrape.RegisterModule(Module{})
}

type Module struct {
	Depth int `json:"depth"`
}

func (Module) ModuleInfo() flyscrape.ModuleInfo {
	return flyscrape.ModuleInfo{
		ID:  "depth",
		New: func() flyscrape.Module { return new(Module) },
	}
}

func (m *Module) ValidateRequest(r *flyscrape.Request) bool {
	return r.Depth <= m.Depth
}

var _ flyscrape.RequestValidator = (*Module)(nil)
