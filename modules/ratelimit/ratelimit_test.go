// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package ratelimit_test

import (
	"testing"
	"time"

	"github.com/philippta/flyscrape"
	"github.com/philippta/flyscrape/modules/followlinks"
	"github.com/philippta/flyscrape/modules/ratelimit"
	"github.com/philippta/flyscrape/modules/starturl"
	"github.com/stretchr/testify/require"
)

func TestRatelimit(t *testing.T) {
	scraper := flyscrape.NewScraper()
	scraper.LoadModule(&starturl.Module{URL: "http://www.example.com/"})
	scraper.LoadModule(&followlinks.Module{})
	scraper.LoadModule(&ratelimit.Module{
		Rate: 100,
	})

	scraper.SetTransport(flyscrape.MockTransport(200, `<a href="foo">foo</a>`))

	var times []time.Time
	scraper.OnRequest(func(req *flyscrape.Request) {
		times = append(times, time.Now())
	})

	start := time.Now()

	scraper.Run()

	first := times[0].Add(-10 * time.Millisecond)
	second := times[1].Add(-20 * time.Millisecond)

	require.Less(t, first.Sub(start), 2*time.Millisecond)
	require.Less(t, second.Sub(start), 2*time.Millisecond)

	require.Less(t, start.Sub(first), 2*time.Millisecond)
	require.Less(t, start.Sub(second), 2*time.Millisecond)
}
