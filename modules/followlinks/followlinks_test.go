// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package followlinks_test

import (
	"testing"

	"github.com/philippta/flyscrape"
	"github.com/philippta/flyscrape/modules/followlinks"
	"github.com/philippta/flyscrape/modules/starturl"
	"github.com/stretchr/testify/require"
)

func TestFollowLinks(t *testing.T) {
	scraper := flyscrape.NewScraper()
	scraper.LoadModule(&starturl.Module{URL: "http://www.example.com/foo/bar"})
	scraper.LoadModule(&followlinks.Module{})

	scraper.SetTransport(flyscrape.MockTransport(200, `
        <a href="/baz">Baz</a>
        <a href="baz">Baz</a>
        <a href="http://www.google.com">Google</a>`))

	var urls []string
	scraper.OnRequest(func(req *flyscrape.Request) {
		urls = append(urls, req.URL)
	})

	scraper.Run()

	require.Len(t, urls, 5)
	require.Contains(t, urls, "http://www.example.com/baz")
	require.Contains(t, urls, "http://www.example.com/foo/bar")
	require.Contains(t, urls, "http://www.example.com/foo/baz")
	require.Contains(t, urls, "http://www.google.com")
	require.Contains(t, urls, "http://www.google.com/baz")
}
