// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package urlfilter

import (
	"regexp"

	"github.com/philippta/flyscrape"
)

func init() {
	flyscrape.RegisterModule(Module{})
}

type Module struct {
	URL         string   `json:"url"`
	AllowedURLs []string `json:"allowedURLs"`
	BlockedURLs []string `json:"blockedURLs"`

	allowedURLsRE []*regexp.Regexp
	blockedURLsRE []*regexp.Regexp
}

func (Module) ModuleInfo() flyscrape.ModuleInfo {
	return flyscrape.ModuleInfo{
		ID:  "urlfilter",
		New: func() flyscrape.Module { return new(Module) },
	}
}

func (m *Module) Provision(v flyscrape.Context) {
	if m.disabled() {
		return
	}

	for _, pat := range m.AllowedURLs {
		re, err := regexp.Compile(pat)
		if err != nil {
			continue
		}
		m.allowedURLsRE = append(m.allowedURLsRE, re)
	}

	for _, pat := range m.BlockedURLs {
		re, err := regexp.Compile(pat)
		if err != nil {
			continue
		}
		m.blockedURLsRE = append(m.blockedURLsRE, re)
	}
}

func (m *Module) ValidateRequest(r *flyscrape.Request) bool {
	if m.disabled() {
		return true
	}

	// allow root url
	if r.URL == m.URL {
		return true
	}

	// allow if no filter is set
	if len(m.allowedURLsRE) == 0 && len(m.blockedURLsRE) == 0 {
		return true
	}

	ok := false
	if len(m.allowedURLsRE) == 0 {
		ok = true
	}

	for _, re := range m.allowedURLsRE {
		if re.MatchString(r.URL) {
			ok = true
			break
		}
	}

	for _, re := range m.blockedURLsRE {
		if re.MatchString(r.URL) {
			ok = false
			break
		}
	}

	return ok
}

func (m *Module) disabled() bool {
	return len(m.AllowedURLs) == 0 && len(m.BlockedURLs) == 0
}

var (
	_ flyscrape.RequestValidator = (*Module)(nil)
	_ flyscrape.Provisioner      = (*Module)(nil)
)
