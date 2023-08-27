// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package flyscrape_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"flyscrape"

	"github.com/stretchr/testify/require"
)

func TestScrapeFollowLinks(t *testing.T) {
	scr := flyscrape.Scraper{
		ScrapeOptions: flyscrape.ScrapeOptions{
			URL:          "http://www.example.com/foo/bar",
			Depth:        1,
			AllowDomains: []string{"www.google.com"},
		},
		ScrapeFunc: func(params flyscrape.ScrapeParams) (any, error) {
			return "foobar", nil
		},
		FetchFunc: func(url string) (string, error) {
			return `<a href="/baz">Baz</a>
                    <a href="baz">Baz</a>
                    <a href="http://www.google.com">Google</a>`, nil
		},
	}

	urls := make(map[string]struct{})
	for res := range scr.Scrape() {
		urls[res.URL] = struct{}{}
	}

	require.Len(t, urls, 4)
	require.Contains(t, urls, "http://www.example.com/baz")
	require.Contains(t, urls, "http://www.example.com/foo/bar")
	require.Contains(t, urls, "http://www.example.com/foo/baz")
	require.Contains(t, urls, "http://www.google.com/")
}

func TestScrapeDepth(t *testing.T) {
	scr := flyscrape.Scraper{
		ScrapeOptions: flyscrape.ScrapeOptions{
			URL:          "http://www.example.com/",
			Depth:        2,
			AllowDomains: []string{"*"},
		},
		ScrapeFunc: func(params flyscrape.ScrapeParams) (any, error) {
			return "foobar", nil
		},
		FetchFunc: func(url string) (string, error) {
			switch url {
			case "http://www.example.com/":
				return `<a href="http://www.google.com">Google</a>`, nil
			case "http://www.google.com/":
				return `<a href="http://www.duckduckgo.com">DuckDuckGo</a>`, nil
			case "http://www.duckduckgo.com/":
				return `<a href="http://www.example.com">Example</a>`, nil
			}
			return "", nil
		},
	}

	urls := make(map[string]struct{})
	for res := range scr.Scrape() {
		urls[res.URL] = struct{}{}
	}

	require.Len(t, urls, 3)
	require.Contains(t, urls, "http://www.example.com/")
	require.Contains(t, urls, "http://www.google.com/")
	require.Contains(t, urls, "http://www.duckduckgo.com/")
}

func TestScrapeAllowDomains(t *testing.T) {
	scr := flyscrape.Scraper{
		ScrapeOptions: flyscrape.ScrapeOptions{
			URL:          "http://www.example.com/",
			Depth:        1,
			AllowDomains: []string{"www.google.com"},
		},
		ScrapeFunc: func(params flyscrape.ScrapeParams) (any, error) {
			return "foobar", nil
		},
		FetchFunc: func(url string) (string, error) {
			return `<a href="http://www.google.com">Google</a>
                    <a href="http://www.duckduckgo.com">DuckDuckGo</a>`, nil
		},
	}

	urls := make(map[string]struct{})
	for res := range scr.Scrape() {
		urls[res.URL] = struct{}{}
	}

	require.Len(t, urls, 2)
	require.Contains(t, urls, "http://www.example.com/")
	require.Contains(t, urls, "http://www.google.com/")
}

func TestScrapeAllowDomainsAll(t *testing.T) {
	scr := flyscrape.Scraper{
		ScrapeOptions: flyscrape.ScrapeOptions{
			URL:          "http://www.example.com/",
			Depth:        1,
			AllowDomains: []string{"*"},
		},
		ScrapeFunc: func(params flyscrape.ScrapeParams) (any, error) {
			return "foobar", nil
		},
		FetchFunc: func(url string) (string, error) {
			return `<a href="http://www.google.com">Google</a>
                    <a href="http://www.duckduckgo.com">DuckDuckGo</a>`, nil
		},
	}

	urls := make(map[string]struct{})
	for res := range scr.Scrape() {
		urls[res.URL] = struct{}{}
	}

	require.Len(t, urls, 3)
	require.Contains(t, urls, "http://www.example.com/")
	require.Contains(t, urls, "http://www.duckduckgo.com/")
	require.Contains(t, urls, "http://www.google.com/")
}

func TestScrapeDenyDomains(t *testing.T) {
	scr := flyscrape.Scraper{
		ScrapeOptions: flyscrape.ScrapeOptions{
			URL:          "http://www.example.com/",
			Depth:        1,
			AllowDomains: []string{"*"},
			DenyDomains:  []string{"www.google.com"},
		},
		ScrapeFunc: func(params flyscrape.ScrapeParams) (any, error) {
			return "foobar", nil
		},
		FetchFunc: func(url string) (string, error) {
			return `<a href="http://www.google.com">Google</a>
                    <a href="http://www.duckduckgo.com">DuckDuckGo</a>`, nil
		},
	}

	urls := make(map[string]struct{})
	for res := range scr.Scrape() {
		urls[res.URL] = struct{}{}
	}

	require.Len(t, urls, 2)
	require.Contains(t, urls, "http://www.example.com/")
	require.Contains(t, urls, "http://www.duckduckgo.com/")
}

func TestScrapeAllowURLs(t *testing.T) {
	scr := flyscrape.Scraper{
		ScrapeOptions: flyscrape.ScrapeOptions{
			URL:       "http://www.example.com/",
			Depth:     1,
			AllowURLs: []string{`/foo\?id=\d+`, `/bar$`},
		},
		ScrapeFunc: func(params flyscrape.ScrapeParams) (any, error) {
			return "foobar", nil
		},
		FetchFunc: func(url string) (string, error) {
			return `<a href="foo?id=123">123</a>
			        <a href="foo?id=ABC">ABC</a>
			        <a href="/bar">bar</a>
                    <a href="/barz">barz</a>`, nil
		},
	}

	urls := make(map[string]struct{})
	for res := range scr.Scrape() {
		urls[res.URL] = struct{}{}
	}

	require.Len(t, urls, 3)
	require.Contains(t, urls, "http://www.example.com/")
	require.Contains(t, urls, "http://www.example.com/foo?id=123")
	require.Contains(t, urls, "http://www.example.com/bar")
}

func TestScrapeRate(t *testing.T) {
	scr := flyscrape.Scraper{
		ScrapeOptions: flyscrape.ScrapeOptions{
			URL:   "http://www.example.com/",
			Depth: 1,
			Rate:  100, // every 10ms
		},
		ScrapeFunc: func(params flyscrape.ScrapeParams) (any, error) {
			return "foobar", nil
		},
		FetchFunc: func(url string) (string, error) {
			return `<a href="foo">foo</a>`, nil
		},
	}

	res := scr.Scrape()

	start := time.Now()
	<-res
	first := time.Now().Add(-10 * time.Millisecond)
	<-res
	second := time.Now().Add(-20 * time.Millisecond)

	require.Less(t, first.Sub(start), 2*time.Millisecond)
	require.Less(t, second.Sub(start), 2*time.Millisecond)
}

func TestScrapeProxy(t *testing.T) {
	proxyCalled := false
	proxy := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		proxyCalled = true
		w.Write([]byte(`<a href="http://www.google.com">Google</a>`))
	}))

	scr := flyscrape.Scraper{
		ScrapeOptions: flyscrape.ScrapeOptions{
			URL:   "http://www.example.com/",
			Proxy: proxy.URL,
		},
		ScrapeFunc: func(params flyscrape.ScrapeParams) (any, error) {
			return "foobar", nil
		},
	}

	res := <-scr.Scrape()

	require.True(t, proxyCalled)
	require.Equal(t, "http://www.example.com/", res.URL)
}
