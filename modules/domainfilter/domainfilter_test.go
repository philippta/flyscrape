// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package domainfilter_test

import (
	"testing"

	"github.com/philippta/flyscrape"
	"github.com/philippta/flyscrape/modules/domainfilter"
	"github.com/philippta/flyscrape/modules/followlinks"
	"github.com/philippta/flyscrape/modules/starturl"
	"github.com/stretchr/testify/require"
)

func TestDomainfilterAllowed(t *testing.T) {
	scraper := flyscrape.NewScraper()
	scraper.LoadModule(&starturl.Module{URL: "http://www.example.com"})
	scraper.LoadModule(&followlinks.Module{})
	scraper.LoadModule(&domainfilter.Module{
		URL:            "http://www.example.com",
		AllowedDomains: []string{"www.google.com"},
	})

	scraper.SetTransport(flyscrape.MockTransport(200, `
        <a href="http://www.google.com">Google</a>
        <a href="http://www.duckduckgo.com">DuckDuckGo</a>`))

	var urls []string
	scraper.OnRequest(func(req *flyscrape.Request) {
		urls = append(urls, req.URL)
	})

	scraper.Run()

	require.Len(t, urls, 2)
	require.Contains(t, urls, "http://www.example.com")
	require.Contains(t, urls, "http://www.google.com")
}

func TestDomainfilterAllowedAll(t *testing.T) {
	scraper := flyscrape.NewScraper()
	scraper.LoadModule(&starturl.Module{URL: "http://www.example.com"})
	scraper.LoadModule(&followlinks.Module{})
	scraper.LoadModule(&domainfilter.Module{
		URL:            "http://www.example.com",
		AllowedDomains: []string{"*"},
	})

	scraper.SetTransport(flyscrape.MockTransport(200, `
        <a href="http://www.google.com">Google</a>
        <a href="http://www.duckduckgo.com">DuckDuckGo</a>`))

	var urls []string
	scraper.OnRequest(func(req *flyscrape.Request) {
		urls = append(urls, req.URL)
	})

	scraper.Run()

	require.Len(t, urls, 3)
	require.Contains(t, urls, "http://www.example.com")
	require.Contains(t, urls, "http://www.duckduckgo.com")
	require.Contains(t, urls, "http://www.google.com")
}

func TestDomainfilterBlocked(t *testing.T) {
	scraper := flyscrape.NewScraper()
	scraper.LoadModule(&starturl.Module{URL: "http://www.example.com"})
	scraper.LoadModule(&followlinks.Module{})
	scraper.LoadModule(&domainfilter.Module{
		URL:            "http://www.example.com",
		AllowedDomains: []string{"*"},
		BlockedDomains: []string{"www.google.com"},
	})

	scraper.SetTransport(flyscrape.MockTransport(200, `
        <a href="http://www.google.com">Google</a>
        <a href="http://www.duckduckgo.com">DuckDuckGo</a>`))

	var urls []string
	scraper.OnRequest(func(req *flyscrape.Request) {
		urls = append(urls, req.URL)
	})

	scraper.Run()

	require.Len(t, urls, 2)
	require.Contains(t, urls, "http://www.example.com")
	require.Contains(t, urls, "http://www.duckduckgo.com")
}
