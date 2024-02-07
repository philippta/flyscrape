// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package retry_test

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"testing"
	"time"

	"github.com/philippta/flyscrape"
	"github.com/philippta/flyscrape/modules/followlinks"
	"github.com/philippta/flyscrape/modules/hook"
	"github.com/philippta/flyscrape/modules/retry"
	"github.com/philippta/flyscrape/modules/starturl"
	"github.com/stretchr/testify/require"
)

func TestRetry(t *testing.T) {
	t.Parallel()
	var count int

	mods := []flyscrape.Module{
		&starturl.Module{URL: "http://www.example.com"},
		&followlinks.Module{},
		hook.Module{
			AdaptTransportFn: func(rt http.RoundTripper) http.RoundTripper {
				return flyscrape.RoundTripFunc(func(r *http.Request) (*http.Response, error) {
					count++
					return flyscrape.MockResponse(http.StatusServiceUnavailable, "service unavailable")
				})
			},
		},
		&retry.Module{
			RetryDelays: []time.Duration{
				100 * time.Millisecond,
				200 * time.Millisecond,
			},
		},
	}

	scraper := flyscrape.NewScraper()
	scraper.Modules = mods
	scraper.Run()

	require.Equal(t, 3, count)
}

func TestRetryStatusCodes(t *testing.T) {
	t.Parallel()

	tests := []struct {
		statusCode int
		retry      bool
	}{
		{statusCode: http.StatusBadGateway, retry: true},
		{statusCode: http.StatusTooManyRequests, retry: true},
		{statusCode: http.StatusBadRequest, retry: false},
		{statusCode: http.StatusOK, retry: false},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("%s_%t", http.StatusText(test.statusCode), test.retry), func(t *testing.T) {
			t.Parallel()
			var count int
			mods := []flyscrape.Module{
				&starturl.Module{URL: "http://www.example.com"},
				&followlinks.Module{},
				hook.Module{
					AdaptTransportFn: func(rt http.RoundTripper) http.RoundTripper {
						return flyscrape.RoundTripFunc(func(r *http.Request) (*http.Response, error) {
							count++
							return flyscrape.MockResponse(test.statusCode, http.StatusText(test.statusCode))
						})
					},
				},
				&retry.Module{
					RetryDelays: []time.Duration{
						100 * time.Millisecond,
						200 * time.Millisecond,
					},
				},
			}

			scraper := flyscrape.NewScraper()
			scraper.Modules = mods
			scraper.Run()

			if test.retry {
				require.NotEqual(t, 1, count)
			} else {
				require.Equal(t, 1, count)
			}
		})
	}
}

func TestRetryErrors(t *testing.T) {
	t.Parallel()

	tests := []struct {
		error error
	}{
		{error: &net.OpError{}},
		{error: io.ErrUnexpectedEOF},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("%T", test.error), func(t *testing.T) {
			t.Parallel()
			var count int
			mods := []flyscrape.Module{
				&starturl.Module{URL: "http://www.example.com"},
				&followlinks.Module{},
				hook.Module{
					AdaptTransportFn: func(rt http.RoundTripper) http.RoundTripper {
						return flyscrape.RoundTripFunc(func(r *http.Request) (*http.Response, error) {
							return nil, test.error
						})
					},
				},
				&retry.Module{
					RetryDelays: []time.Duration{
						100 * time.Millisecond,
						200 * time.Millisecond,
					},
				},
			}

			scraper := flyscrape.NewScraper()
			scraper.Modules = mods
			scraper.Run()

			require.NotEqual(t, 1, count)
		})
	}
}
