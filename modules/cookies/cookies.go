// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package cookies

import (
	"net/http"
	"slices"

	"github.com/browserutils/kooky"
	_ "github.com/browserutils/kooky/browser/chrome"
	_ "github.com/browserutils/kooky/browser/edge"
	_ "github.com/browserutils/kooky/browser/firefox"
	"github.com/philippta/flyscrape"
)

var supportedBrowsers = []string{
	"chrome",
	"edge",
	"firefox",
}

func init() {
	flyscrape.RegisterModule(Module{})
}

type Module struct {
	Cookies string `json:"cookies"`
}

func (Module) ModuleInfo() flyscrape.ModuleInfo {
	return flyscrape.ModuleInfo{
		ID:  "cookies",
		New: func() flyscrape.Module { return new(Module) },
	}
}

func (m Module) AdaptTransport(t http.RoundTripper) http.RoundTripper {
	if !slices.Contains(supportedBrowsers, m.Cookies) {
		return t
	}

	var stores []kooky.CookieStore
	for _, store := range kooky.FindAllCookieStores() {
		if store.Browser() == m.Cookies && store.IsDefaultProfile() {
			stores = append(stores, store)
		}
	}

	return flyscrape.RoundTripFunc(func(r *http.Request) (*http.Response, error) {
		for _, store := range stores {
			for _, cookie := range store.Cookies(r.URL) {
				r.AddCookie(cookie)
			}
		}
		return t.RoundTrip(r)
	})
}

var _ flyscrape.TransportAdapter = Module{}
