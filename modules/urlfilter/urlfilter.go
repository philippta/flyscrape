// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package urlfilter

import (
	"regexp"

	"github.com/philippta/flyscrape"
)

func init() {
	flyscrape.RegisterModule(new(Module))
}

type Module struct {
	URL         string   `json:"url"`
	AllowedURLs []string `json:"allowedURLs"`
	BlockedURLs []string `json:"blockedURLs"`

	allowedURLsRE []*regexp.Regexp
	blockedURLsRE []*regexp.Regexp
}

func (m *Module) OnLoad(v flyscrape.Visitor) {
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

func (m *Module) CanRequest(rawurl string, depth int) bool {
	// allow root url
	if rawurl == m.URL {
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
		if re.MatchString(rawurl) {
			ok = true
			break
		}
	}

	for _, re := range m.blockedURLsRE {
		if re.MatchString(rawurl) {
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
