// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package urlfilter_test

import (
	"net/http"
	"sync"
	"testing"

	"github.com/philippta/flyscrape"
	"github.com/philippta/flyscrape/modules/followlinks"
	"github.com/philippta/flyscrape/modules/hook"
	"github.com/philippta/flyscrape/modules/starturl"
	"github.com/philippta/flyscrape/modules/urlfilter"
	"github.com/stretchr/testify/require"
)

func TestURLFilterAllowed(t *testing.T) {
	var urls []string
	var mu sync.Mutex

	scraper := flyscrape.NewScraper()
	scraper.LoadModule(&starturl.Module{URL: "http://www.example.com/"})
	scraper.LoadModule(&followlinks.Module{})
	scraper.LoadModule(&urlfilter.Module{
		URL:         "http://www.example.com/",
		AllowedURLs: []string{`/foo\?id=\d+`, `/bar$`},
	})
	scraper.LoadModule(hook.Module{
		AdaptTransportFn: func(rt http.RoundTripper) http.RoundTripper {
			return flyscrape.MockTransport(200, `
				<a href="foo?id=123">123</a>
				<a href="foo?id=ABC">ABC</a>
				<a href="/bar">bar</a>
				<a href="/barz">barz</a>`)
		},
		ReceiveResponseFn: func(r *flyscrape.Response) {
			mu.Lock()
			urls = append(urls, r.Request.URL)
			mu.Unlock()
		},
	})

	scraper.Run()

	require.Len(t, urls, 3)
	require.Contains(t, urls, "http://www.example.com/")
	require.Contains(t, urls, "http://www.example.com/foo?id=123")
	require.Contains(t, urls, "http://www.example.com/bar")
}

func TestURLFilterBlocked(t *testing.T) {
	var urls []string
	var mu sync.Mutex

	scraper := flyscrape.NewScraper()
	scraper.LoadModule(&starturl.Module{URL: "http://www.example.com/"})
	scraper.LoadModule(&followlinks.Module{})
	scraper.LoadModule(&urlfilter.Module{
		URL:         "http://www.example.com/",
		BlockedURLs: []string{`/foo\?id=\d+`, `/bar$`},
	})
	scraper.LoadModule(hook.Module{
		AdaptTransportFn: func(rt http.RoundTripper) http.RoundTripper {
			return flyscrape.MockTransport(200, `
				<a href="foo?id=123">123</a>
				<a href="foo?id=ABC">ABC</a>
				<a href="/bar">bar</a>
				<a href="/barz">barz</a>`)
		},
		ReceiveResponseFn: func(r *flyscrape.Response) {
			mu.Lock()
			urls = append(urls, r.Request.URL)
			mu.Unlock()
		},
	})

	scraper.Run()

	require.Len(t, urls, 3)
	require.Contains(t, urls, "http://www.example.com/")
	require.Contains(t, urls, "http://www.example.com/foo?id=ABC")
	require.Contains(t, urls, "http://www.example.com/barz")
}
