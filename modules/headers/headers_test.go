// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package headers_test

import (
	"net/http"
	"strings"
	"testing"

	"github.com/philippta/flyscrape"
	"github.com/philippta/flyscrape/modules/headers"
	"github.com/philippta/flyscrape/modules/hook"
	"github.com/philippta/flyscrape/modules/starturl"
	"github.com/stretchr/testify/require"
)

func TestHeaders(t *testing.T) {
	gotHeaders := map[string]string{}
	sentHeaders := map[string]string{
		"Authorization": "Basic ZGVtbzpwQDU1dzByZA==",
		"User-Agent":    "Gecko/1.0",
	}

	mods := []flyscrape.Module{
		hook.Module{
			AdaptTransportFn: func(rt http.RoundTripper) http.RoundTripper {
				return flyscrape.RoundTripFunc(func(r *http.Request) (*http.Response, error) {
					for k := range r.Header {
						gotHeaders[k] = r.Header.Get(k)
					}
					return flyscrape.MockResponse(200, "")
				})
			},
		},
		&starturl.Module{URL: "http://www.example.com"},
		&headers.Module{
			Headers: sentHeaders,
		},
	}

	scraper := flyscrape.NewScraper()
	scraper.Modules = mods
	scraper.Run()

	require.Equal(t, sentHeaders, gotHeaders)
}

func TestHeadersRandomUserAgent(t *testing.T) {
	var userAgent string
	mods := []flyscrape.Module{
		hook.Module{
			AdaptTransportFn: func(rt http.RoundTripper) http.RoundTripper {
				return flyscrape.RoundTripFunc(func(r *http.Request) (*http.Response, error) {
					userAgent = r.Header.Get("User-Agent")
					return flyscrape.MockResponse(200, "")
				})
			},
		},
		&starturl.Module{URL: "http://www.example.com"},
		&headers.Module{},
	}

	scraper := flyscrape.NewScraper()
	scraper.Modules = mods
	scraper.Run()

	require.NotEmpty(t, userAgent)
	require.True(t, strings.HasPrefix(userAgent, "Mozilla/5.0 ("))
}
