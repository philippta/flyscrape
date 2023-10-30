// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package starturl_test

import (
	"net/http"
	"testing"

	"github.com/philippta/flyscrape"
	"github.com/philippta/flyscrape/modules/hook"
	"github.com/philippta/flyscrape/modules/starturl"
	"github.com/stretchr/testify/require"
)

func TestStartURL(t *testing.T) {
	var url string
	var depth int

	mods := []flyscrape.Module{
		&starturl.Module{URL: "http://www.example.com/foo/bar"},
		hook.Module{
			AdaptTransportFn: func(rt http.RoundTripper) http.RoundTripper {
				return flyscrape.MockTransport(200, "")
			},
			BuildRequestFn: func(r *flyscrape.Request) {
				url = r.URL
				depth = r.Depth
			},
		},
	}

	scraper := flyscrape.NewScraper()
	scraper.Modules = mods
	scraper.Run()

	require.Equal(t, "http://www.example.com/foo/bar", url)
	require.Equal(t, 0, depth)
}
