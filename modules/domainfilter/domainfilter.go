// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package domainfilter

import (
	"github.com/nlnwa/whatwg-url/url"
	"github.com/philippta/flyscrape"
)

func init() {
	flyscrape.RegisterModule(Module{})
}

type Module struct {
	URL            string   `json:"url"`
	URLs           []string `json:"urls"`
	AllowedDomains []string `json:"allowedDomains"`
	BlockedDomains []string `json:"blockedDomains"`

	active bool
}

func (Module) ModuleInfo() flyscrape.ModuleInfo {
	return flyscrape.ModuleInfo{
		ID:  "domainfilter",
		New: func() flyscrape.Module { return new(Module) },
	}
}

func (m *Module) Provision(v flyscrape.Context) {
	if m.URL != "" {
		if u, err := url.Parse(m.URL); err == nil {
			m.AllowedDomains = append(m.AllowedDomains, u.Host())
		}
	}
	for _, u := range m.URLs {
		if u, err := url.Parse(u); err == nil {
			m.AllowedDomains = append(m.AllowedDomains, u.Host())
		}
	}
}

func (m *Module) ValidateRequest(r *flyscrape.Request) bool {
	if m.disabled() {
		return true
	}

	u, err := url.Parse(r.URL)
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

func (m *Module) disabled() bool {
	return len(m.AllowedDomains) == 0 && len(m.BlockedDomains) == 0
}

var (
	_ flyscrape.RequestValidator = (*Module)(nil)
	_ flyscrape.Provisioner      = (*Module)(nil)
)
