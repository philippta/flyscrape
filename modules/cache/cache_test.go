// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package cache_test

import (
	"net/http"
	"testing"

	"github.com/philippta/flyscrape"
	"github.com/philippta/flyscrape/modules/cache"
	"github.com/philippta/flyscrape/modules/hook"
	"github.com/philippta/flyscrape/modules/starturl"
	"github.com/stretchr/testify/require"
)

func TestCache(t *testing.T) {
	cachemod := &cache.Module{Cache: "memory"}
	calls := 0

	for i := 0; i < 2; i++ {
		scraper := flyscrape.NewScraper()
		scraper.LoadModule(&starturl.Module{URL: "http://www.example.com"})
		scraper.LoadModule(hook.Module{
			AdaptTransportFn: func(rt http.RoundTripper) http.RoundTripper {
				return flyscrape.RoundTripFunc(func(r *http.Request) (*http.Response, error) {
					calls++
					return flyscrape.MockResponse(200, "foo")
				})
			},
		})
		scraper.LoadModule(cachemod)
		scraper.Run()
	}

	require.Equal(t, 1, calls)
}
