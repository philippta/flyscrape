// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package starturl

import (
	"github.com/philippta/flyscrape"
)

func init() {
	flyscrape.RegisterModule(new(Module))
}

type Module struct {
	URL string `json:"url"`
}

func (m *Module) OnLoad(v flyscrape.Visitor) {
	v.Visit(m.URL)
}

var _ flyscrape.OnLoad = (*Module)(nil)
