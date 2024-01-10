// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package proxy_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/philippta/flyscrape"
	"github.com/philippta/flyscrape/modules/proxy"
	"github.com/philippta/flyscrape/modules/starturl"
	"github.com/stretchr/testify/require"
)

func TestProxy(t *testing.T) {
	var called bool
	p := newProxy(func() { called = true })
	defer p.Close()

	mods := []flyscrape.Module{
		&starturl.Module{URL: "http://www.example.com"},
		&proxy.Module{
			Proxies: []string{p.URL},
		},
	}

	scraper := flyscrape.NewScraper()
	scraper.Modules = mods
	scraper.Run()

	require.True(t, called)
}

func TestProxyMultiple(t *testing.T) {
	calls := []int{0, 0, 0}
	p0 := newProxy(func() { calls[0]++ })
	p1 := newProxy(func() { calls[1]++ })
	p2 := newProxy(func() { calls[2]++ })
	defer p0.Close()
	defer p1.Close()
	defer p2.Close()

	mod := &proxy.Module{Proxies: []string{p0.URL, p1.URL}, Proxy: p2.URL}
	mod.Provision(nil)
	trans := mod.AdaptTransport(nil)

	req := httptest.NewRequest("GET", "http://www.example.com/", nil)

	for i := 0; i < 50; i++ {
		resp, err := trans.RoundTrip(req)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, resp.StatusCode)
	}

	require.Greater(t, calls[0], 1)
	require.Greater(t, calls[1], 1)
	require.Greater(t, calls[2], 1)
	require.Equal(t, 50, calls[0]+calls[1]+calls[2])
}

func newProxy(f func()) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		f()
		w.Write([]byte("response from proxy"))
	}))
}
