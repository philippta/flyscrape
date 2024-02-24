// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package browser_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/philippta/flyscrape"
	"github.com/philippta/flyscrape/modules/browser"
	"github.com/philippta/flyscrape/modules/headers"
	"github.com/philippta/flyscrape/modules/hook"
	"github.com/philippta/flyscrape/modules/starturl"
	"github.com/stretchr/testify/require"
)

func TestBrowser(t *testing.T) {
	var called bool

	srv := newServer(func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.Write([]byte(`<h1>Hello Browser</h1><a href="foo">Foo</a>`))
	})
	defer srv.Close()

	var body string

	mods := []flyscrape.Module{
		&starturl.Module{URL: srv.URL},
		&browser.Module{Browser: true},
		&hook.Module{
			ReceiveResponseFn: func(r *flyscrape.Response) {
				body = string(r.Body)
			},
		},
	}

	scraper := flyscrape.NewScraper()
	scraper.Modules = mods
	scraper.Run()

	require.True(t, called)
	require.Contains(t, body, "Hello Browser")
}
func TestBrowserStatusCode(t *testing.T) {
	srv := newServer(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(404)
	})
	defer srv.Close()

	var statusCode int

	mods := []flyscrape.Module{
		&starturl.Module{URL: srv.URL},
		&browser.Module{Browser: true},
		&hook.Module{
			ReceiveResponseFn: func(r *flyscrape.Response) {
				statusCode = r.StatusCode
			},
		},
	}

	scraper := flyscrape.NewScraper()
	scraper.Modules = mods
	scraper.Run()

	require.Equal(t, 404, statusCode)
}

func TestBrowserRequestHeader(t *testing.T) {
	srv := newServer(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(r.Header.Get("User-Agent")))
	})
	defer srv.Close()

	var body string

	mods := []flyscrape.Module{
		&starturl.Module{URL: srv.URL},
		&browser.Module{Browser: true},
		&headers.Module{
			Headers: map[string]string{
				"User-Agent": "custom-headers",
			},
		},
		&hook.Module{
			ReceiveResponseFn: func(r *flyscrape.Response) {
				body = string(r.Body)
			},
		},
	}

	scraper := flyscrape.NewScraper()
	scraper.Modules = mods
	scraper.Run()

	require.Contains(t, body, "custom-headers")
}

func TestBrowserResponseHeader(t *testing.T) {
	srv := newServer(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Foo", "bar")
	})
	defer srv.Close()

	var header string

	mods := []flyscrape.Module{
		&starturl.Module{URL: srv.URL},
		&browser.Module{Browser: true},
		&hook.Module{
			ReceiveResponseFn: func(r *flyscrape.Response) {
				header = r.Headers.Get("Foo")
			},
		},
	}

	scraper := flyscrape.NewScraper()
	scraper.Modules = mods
	scraper.Run()

	require.Equal(t, header, "bar")
}

func TestBrowserUnsetFlyscrapeUserAgent(t *testing.T) {
	srv := newServer(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(r.Header.Get("User-Agent")))
	})
	defer srv.Close()

	var body string

	mods := []flyscrape.Module{
		&starturl.Module{URL: srv.URL},
		&browser.Module{Browser: true},
		&hook.Module{
			ReceiveResponseFn: func(r *flyscrape.Response) {
				body = string(r.Body)
			},
		},
	}

	scraper := flyscrape.NewScraper()
	scraper.Modules = mods
	scraper.Run()

	fmt.Println(body)
	require.Contains(t, body, "Mozilla/5.0")
	require.NotContains(t, body, "flyscrape")
}

func newServer(f func(http.ResponseWriter, *http.Request)) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		f(w, r)
	}))
}
