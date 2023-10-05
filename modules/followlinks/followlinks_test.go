// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package followlinks_test

import (
	"net/http"
	"sync"
	"testing"

	"github.com/philippta/flyscrape"
	"github.com/philippta/flyscrape/modules/followlinks"
	"github.com/philippta/flyscrape/modules/hook"
	"github.com/philippta/flyscrape/modules/starturl"
	"github.com/stretchr/testify/require"
)

func TestFollowLinks(t *testing.T) {
	var urls []string
	var mu sync.Mutex

	scraper := flyscrape.NewScraper()
	scraper.LoadModule(&starturl.Module{URL: "http://www.example.com/foo/bar"})
	scraper.LoadModule(&followlinks.Module{})

	scraper.LoadModule(hook.Module{
		AdaptTransportFn: func(rt http.RoundTripper) http.RoundTripper {
			return flyscrape.MockTransport(200, `
				<a href="/baz">Baz</a>
				<a href="baz">Baz</a>
				<a href="http://www.google.com">Google</a>`)
		},
		ReceiveResponseFn: func(r *flyscrape.Response) {
			mu.Lock()
			urls = append(urls, r.Request.URL)
			mu.Unlock()
		},
	})

	scraper.Run()

	require.Len(t, urls, 5)
	require.Contains(t, urls, "http://www.example.com/baz")
	require.Contains(t, urls, "http://www.example.com/foo/bar")
	require.Contains(t, urls, "http://www.example.com/foo/baz")
	require.Contains(t, urls, "http://www.google.com")
	require.Contains(t, urls, "http://www.google.com/baz")
}
