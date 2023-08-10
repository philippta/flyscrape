package flyscrape_test

import (
	"sort"
	"testing"

	"flyscrape/flyscrape"

	"github.com/stretchr/testify/require"
)

func TestScrape(t *testing.T) {
	svc := flyscrape.Service{
		Concurrency: 10,
		ScrapeOptions: flyscrape.ScrapeOptions{
			URL:   "http://example.com/foo/bar",
			Depth: 1,
		},
		ScrapeFunc: func(params flyscrape.ScrapeParams) (flyscrape.M, error) {
			return flyscrape.M{
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

	require.Len(t, urls, 5)
	require.Equal(t, "http://example.com/baz", urls[0])
	require.Equal(t, "http://example.com/foo/bar", urls[1])
	require.Equal(t, "http://example.com/foo/baz", urls[2])
	require.Equal(t, "http://www.google.com/", urls[3])
	require.Equal(t, "http://www.google.com/baz", urls[4])
}

func TestFindLinks(t *testing.T) {
	origin := "http://example.com/foo/bar"
	html := `
        <html>
            <body>
                <a href="/baz">Baz</a>
                <a href="baz">Baz</a>
                <a href="http://www.google.com">Google</a>
                <a href="javascript:void(0)">Google</a>
                <a href="/foo#hello">Anchor</a>
            </body>
        </html>`

	links := flyscrape.Links(html, origin)
	require.Len(t, links, 4)
	require.Equal(t, "http://example.com/baz", links[0])
	require.Equal(t, "http://example.com/foo/baz", links[1])
	require.Equal(t, "http://www.google.com/", links[2])
	require.Equal(t, "http://example.com/foo", links[3])
}
