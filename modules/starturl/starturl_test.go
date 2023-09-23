// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package starturl_test

import (
	"testing"

	"github.com/philippta/flyscrape"
	"github.com/philippta/flyscrape/modules/starturl"
	"github.com/stretchr/testify/require"
)

func TestFollowLinks(t *testing.T) {
	scraper := flyscrape.NewScraper()
	scraper.LoadModule(&starturl.Module{URL: "http://www.example.com/foo/bar"})
	scraper.SetTransport(flyscrape.MockTransport(200, ""))

	var url string
	var depth int
	scraper.OnRequest(func(req *flyscrape.Request) {
		url = req.URL
		depth = req.Depth
	})

	scraper.Run()

	require.Equal(t, "http://www.example.com/foo/bar", url)
	require.Equal(t, 0, depth)
}
