// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package headers_test

import (
	"fmt"
	"net/http"
	"reflect"
	"testing"

	"github.com/philippta/flyscrape"
	"github.com/philippta/flyscrape/modules/headers"
	"github.com/philippta/flyscrape/modules/hook"
	"github.com/philippta/flyscrape/modules/starturl"
	"github.com/stretchr/testify/require"
)

func TestHeaders(t *testing.T) {
	testCases := []struct {
		name        string
		headersFn   func() headers.Module
		wantHeaders map[string][]string
	}{
		{
			name: "empty custom headers",
			headersFn: func() headers.Module {
				return headers.Module{
					Headers: map[string]string{},
				}
			},
			wantHeaders: map[string][]string{"User-Agent": {"flyscrape/0.1"}},
		},
		{
			name: "non-empty custom headers",
			headersFn: func() headers.Module {
				return headers.Module{
					Headers: map[string]string{
						"Basic":      "ZGVtbzpwQDU1dzByZA==",
						"User-Agent": "Gecko/1.0",
					},
				}
			},
			wantHeaders: map[string][]string{
				"Basic":      {"ZGVtbzpwQDU1dzByZA=="},
				"User-Agent": {"Gecko/1.0"},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var headers map[string][]string

			mods := []flyscrape.Module{
				&starturl.Module{URL: "http://www.example.com"},
				hook.Module{
					AdaptTransportFn: func(rt http.RoundTripper) http.RoundTripper {
						return flyscrape.RoundTripFunc(func(r *http.Request) (*http.Response, error) {
							headers = r.Header
							return rt.RoundTrip(r)
						})
					},
				},
				tc.headersFn(),
			}

			scraper := flyscrape.NewScraper()
			scraper.Modules = mods
			scraper.Run()

			require.Truef(
				t,
				reflect.DeepEqual(tc.wantHeaders, headers),
				fmt.Sprintf("%v does not equal to %v", tc.wantHeaders, headers),
			)
		})
	}
}
