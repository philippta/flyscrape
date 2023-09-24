// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package domainfilter

import (
	"github.com/nlnwa/whatwg-url/url"
	"github.com/philippta/flyscrape"
)

func init() {
	flyscrape.RegisterModule(new(Module))
}

type Module struct {
	URL            string   `json:"url"`
	AllowedDomains []string `json:"allowedDomains"`
	BlockedDomains []string `json:"blockedDomains"`
}

func (m *Module) OnLoad(v flyscrape.Visitor) {
	if u, err := url.Parse(m.URL); err == nil {
		m.AllowedDomains = append(m.AllowedDomains, u.Host())
	}
}

func (m *Module) CanRequest(rawurl string, depth int) bool {
	u, err := url.Parse(rawurl)
	if err != nil {
		return false
	}

	host := u.Host()
	ok := false

	for _, domain := range m.AllowedDomains {
		if domain == "*" || host == domain {
			ok = true
			break
		}
	}

	for _, domain := range m.BlockedDomains {
		if host == domain {
			ok = false
			break
		}
	}

	return ok
}

var (
	_ flyscrape.CanRequest = (*Module)(nil)
	_ flyscrape.OnLoad     = (*Module)(nil)
)
