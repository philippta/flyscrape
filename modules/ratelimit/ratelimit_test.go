// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package ratelimit_test

import (
	"net/http"
	"sync"
	"testing"
	"time"

	"github.com/philippta/flyscrape"
	"github.com/philippta/flyscrape/modules/followlinks"
	"github.com/philippta/flyscrape/modules/hook"
	"github.com/philippta/flyscrape/modules/ratelimit"
	"github.com/philippta/flyscrape/modules/starturl"
	"github.com/stretchr/testify/require"
)

func TestRatelimit(t *testing.T) {
	var times []time.Time
	var mu sync.Mutex

	mods := []flyscrape.Module{
		&starturl.Module{URL: "http://www.example.com"},
		&followlinks.Module{},
		hook.Module{
			AdaptTransportFn: func(rt http.RoundTripper) http.RoundTripper {
				return flyscrape.MockTransport(200, `<a href="foo">foo</a>`)
			},
			ReceiveResponseFn: func(r *flyscrape.Response) {
				mu.Lock()
				times = append(times, time.Now())
				mu.Unlock()
			},
		},
		&ratelimit.Module{
			Rate: 240,
		},
	}

	start := time.Now()
	scraper := flyscrape.NewScraper()
	scraper.Modules = mods
	scraper.Run()

	first := times[0].Add(-250 * time.Millisecond)
	second := times[1].Add(-500 * time.Millisecond)

	require.Less(t, first.Sub(start), 250*time.Millisecond)
	require.Less(t, second.Sub(start), 250*time.Millisecond)

	require.Less(t, start.Sub(first), 250*time.Millisecond)
	require.Less(t, start.Sub(second), 250*time.Millisecond)
}

func TestRatelimitConcurrency(t *testing.T) {
	var times []time.Time
	var mu sync.Mutex

	mods := []flyscrape.Module{
		&starturl.Module{URL: "http://www.example.com"},
		&followlinks.Module{},
		hook.Module{
			AdaptTransportFn: func(rt http.RoundTripper) http.RoundTripper {
				return flyscrape.RoundTripFunc(func(r *http.Request) (*http.Response, error) {
					mu.Lock()
					times = append(times, time.Now())
					mu.Unlock()

					time.Sleep(10 * time.Millisecond)
					return flyscrape.MockResponse(200, `
						<a href="foo"></a>
						<a href="bar"></a>
						<a href="baz"></a>
						<a href="qux"></a>
					`)
				})
			},
		},
		&ratelimit.Module{
			Concurrency: 2,
		},
	}

	scraper := flyscrape.NewScraper()
	scraper.Modules = mods
	scraper.Run()

	require.Len(t, times, 5)
	require.Less(t, times[2].Sub(times[1]), time.Millisecond)
	require.Less(t, times[4].Sub(times[3]), time.Millisecond)
}
