// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package depth_test

import (
	"net/http"
	"sync"
	"testing"

	"github.com/philippta/flyscrape"
	"github.com/philippta/flyscrape/modules/depth"
	"github.com/philippta/flyscrape/modules/followlinks"
	"github.com/philippta/flyscrape/modules/hook"
	"github.com/philippta/flyscrape/modules/starturl"
	"github.com/stretchr/testify/require"
)

func TestDepth(t *testing.T) {
	var urls []string
	var mu sync.Mutex

	scraper := flyscrape.NewScraper()
	scraper.LoadModule(&starturl.Module{URL: "http://www.example.com"})
	scraper.LoadModule(&followlinks.Module{})
	scraper.LoadModule(&depth.Module{Depth: 2})
	scraper.LoadModule(hook.Module{
		AdaptTransportFn: func(rt http.RoundTripper) http.RoundTripper {
			return flyscrape.RoundTripFunc(func(r *http.Request) (*http.Response, error) {
				switch r.URL.String() {
				case "http://www.example.com":
					return flyscrape.MockResponse(200, `<a href="http://www.google.com">Google</a>`)
				case "http://www.google.com":
					return flyscrape.MockResponse(200, `<a href="http://www.duckduckgo.com">DuckDuckGo</a>`)
				case "http://www.duckduckgo.com":
					return flyscrape.MockResponse(200, `<a href="http://www.example.com">Example</a>`)
				}
				return flyscrape.MockResponse(200, "")
			})
		},
		ReceiveResponseFn: func(r *flyscrape.Response) {
			mu.Lock()
			urls = append(urls, r.Request.URL)
			mu.Unlock()
		},
	})

	scraper.Run()

	require.Len(t, urls, 3)
	require.Contains(t, urls, "http://www.example.com")
	require.Contains(t, urls, "http://www.google.com")
	require.Contains(t, urls, "http://www.duckduckgo.com")
}
