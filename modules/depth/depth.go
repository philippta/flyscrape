// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package depth

import (
	"github.com/philippta/flyscrape"
)

func init() {
	flyscrape.RegisterModule(new(Module))
}

type Module struct {
	Depth int `json:"depth"`
}

func (m *Module) ID() string {
	return "depth"
}

func (m *Module) CanRequest(url string, depth int) bool {
	return depth <= m.Depth
}

var (
	_ flyscrape.Module     = (*Module)(nil)
	_ flyscrape.CanRequest = (*Module)(nil)
)
