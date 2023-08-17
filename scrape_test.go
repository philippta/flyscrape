package flyscrape_test

import (
	"sort"
	"testing"

	"flyscrape"

	"github.com/stretchr/testify/require"
)

func TestScrape(t *testing.T) {
	svc := flyscrape.Scraper{
		ScrapeOptions: flyscrape.ScrapeOptions{
			URL:            "http://example.com/foo/bar",
			Depth:          1,
			AllowedDomains: []string{"example.com", "www.google.com"},
		},
		ScrapeFunc: func(params flyscrape.ScrapeParams) (any, error) {
			return map[string]any{
				"url": params.URL,
			}, nil
		},
		FetchFunc: func(url string) (string, error) {
			return `<html>
                <body>
                    <a href="/baz">Baz</a>
                    <a href="baz">Baz</a>
                    <a href="http://www.google.com">Google</a>
                </body>
            </html>`, nil
		},
	}

	var urls []string
	for res := range svc.Scrape() {
		urls = append(urls, res.URL)
	}
	sort.Strings(urls)

	require.Len(t, urls, 4)
	require.Equal(t, "http://example.com/baz", urls[0])
	require.Equal(t, "http://example.com/foo/bar", urls[1])
	require.Equal(t, "http://example.com/foo/baz", urls[2])
	require.Equal(t, "http://www.google.com/", urls[3])
}
