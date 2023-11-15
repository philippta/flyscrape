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

func TestStartURL_MultipleStartingURLs(t *testing.T) {
	testCases := []struct {
		name          string
		startURLModFn func() *starturl.Module
		urls          []string
	}{
		{
			name: ".URL and .URLs",
			startURLModFn: func() *starturl.Module {
				return &starturl.Module{
					URL: "http://www.example.com/foo",
					URLs: []string{
						"http://www.example.com/bar",
						"http://www.example.com/baz",
					},
				}
			},
			urls: []string{
				"http://www.example.com/foo",
				"http://www.example.com/bar",
				"http://www.example.com/baz",
			},
		},
		{
			name: "only .URL",
			startURLModFn: func() *starturl.Module {
				return &starturl.Module{
					URL: "http://www.example.com/foo",
				}
			},
			urls: []string{
				"http://www.example.com/foo",
			},
		},
		{
			name: "only .URLs",
			startURLModFn: func() *starturl.Module {
				return &starturl.Module{
					URLs: []string{
						"http://www.example.com/bar",
						"http://www.example.com/baz",
					},
				}
			},
			urls: []string{
				"http://www.example.com/bar",
				"http://www.example.com/baz",
			},
		},
		{
			name: "empty",
			startURLModFn: func() *starturl.Module {
				return &starturl.Module{}
			},
			urls: []string{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			urls := []string{}

			mods := []flyscrape.Module{
				tc.startURLModFn(),
				hook.Module{
					AdaptTransportFn: func(rt http.RoundTripper) http.RoundTripper {
						return flyscrape.MockTransport(http.StatusOK, "")
					},
					BuildRequestFn: func(r *flyscrape.Request) {
						urls = append(urls, r.URL)
					},
				},
			}

			scraper := flyscrape.NewScraper()
			scraper.Modules = mods
			scraper.Run()

			require.ElementsMatch(t, tc.urls, urls)
		})
	}
}
