// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package followlinks_test

import (
	"net/http"
	"sync"
	"testing"

	"github.com/philippta/flyscrape"
	"github.com/philippta/flyscrape/modules/followlinks"
	"github.com/philippta/flyscrape/modules/hook"
	"github.com/philippta/flyscrape/modules/starturl"
	"github.com/stretchr/testify/require"
)

func TestFollowLinks(t *testing.T) {
	var urls []string
	var mu sync.Mutex

	mods := []flyscrape.Module{
		&starturl.Module{URL: "http://www.example.com/foo/bar"},
		&followlinks.Module{},
		hook.Module{
			AdaptTransportFn: func(rt http.RoundTripper) http.RoundTripper {
				return flyscrape.MockTransport(200, `
				<a href="/baz">Baz</a>
				<a href="baz">Baz</a>
				<a href="http://www.google.com">Google</a>`)
			},
			ReceiveResponseFn: func(r *flyscrape.Response) {
				mu.Lock()
				urls = append(urls, r.Request.URL)
				mu.Unlock()
			},
		},
	}

	scraper := flyscrape.NewScraper()
	scraper.Modules = mods
	scraper.Run()

	require.Len(t, urls, 5)
	require.Contains(t, urls, "http://www.example.com/baz")
	require.Contains(t, urls, "http://www.example.com/foo/bar")
	require.Contains(t, urls, "http://www.example.com/foo/baz")
	require.Contains(t, urls, "http://www.google.com")
	require.Contains(t, urls, "http://www.google.com/baz")
}

func TestFollowSelector(t *testing.T) {
	var urls []string
	var mu sync.Mutex

	mods := []flyscrape.Module{
		&starturl.Module{URL: "http://www.example.com/foo/bar"},
		&followlinks.Module{
			Follow: &[]string{".next a[href]"},
		},
		hook.Module{
			AdaptTransportFn: func(rt http.RoundTripper) http.RoundTripper {
				return flyscrape.MockTransport(200, `
				<a href="/baz">Baz</a>
				<a href="baz">Baz</a>
                <div class="next">
				    <a href="http://www.google.com">Google</a>
                </div>`)
			},
			ReceiveResponseFn: func(r *flyscrape.Response) {
				mu.Lock()
				urls = append(urls, r.Request.URL)
				mu.Unlock()
			},
		},
	}

	scraper := flyscrape.NewScraper()
	scraper.Modules = mods
	scraper.Run()

	require.Len(t, urls, 2)
	require.Contains(t, urls, "http://www.example.com/foo/bar")
	require.Contains(t, urls, "http://www.google.com")
}

func TestFollowDataAttr(t *testing.T) {
	var urls []string
	var mu sync.Mutex

	mods := []flyscrape.Module{
		&starturl.Module{URL: "http://www.example.com/foo/bar"},
		&followlinks.Module{
			Follow: &[]string{"[data-url]"},
		},
		hook.Module{
			AdaptTransportFn: func(rt http.RoundTripper) http.RoundTripper {
				return flyscrape.MockTransport(200, `
				<a href="/baz">Baz</a>
				<a href="baz">Baz</a>
				<div data-url="http://www.google.com">Google</div>`)
			},
			ReceiveResponseFn: func(r *flyscrape.Response) {
				mu.Lock()
				urls = append(urls, r.Request.URL)
				mu.Unlock()
			},
		},
	}

	scraper := flyscrape.NewScraper()
	scraper.Modules = mods
	scraper.Run()

	require.Len(t, urls, 2)
	require.Contains(t, urls, "http://www.example.com/foo/bar")
	require.Contains(t, urls, "http://www.google.com")
}

func TestFollowMultiple(t *testing.T) {
	var urls []string
	var mu sync.Mutex

	mods := []flyscrape.Module{
		&starturl.Module{URL: "http://www.example.com/foo/bar"},
		&followlinks.Module{
			Follow: &[]string{"a.prev", "a.next"},
		},
		hook.Module{
			AdaptTransportFn: func(rt http.RoundTripper) http.RoundTripper {
				return flyscrape.MockTransport(200, `
				<a href="/baz">Baz</a>
				<a class="prev" href="a">a</a>
				<a class="next" href="b">b</a>`)
			},
			ReceiveResponseFn: func(r *flyscrape.Response) {
				mu.Lock()
				urls = append(urls, r.Request.URL)
				mu.Unlock()
			},
		},
	}

	scraper := flyscrape.NewScraper()
	scraper.Modules = mods
	scraper.Run()

	require.Len(t, urls, 3)
	require.Contains(t, urls, "http://www.example.com/foo/bar")
	require.Contains(t, urls, "http://www.example.com/foo/a")
	require.Contains(t, urls, "http://www.example.com/foo/b")
}

func TestFollowNoFollow(t *testing.T) {
	var urls []string
	var mu sync.Mutex

	mods := []flyscrape.Module{
		&starturl.Module{URL: "http://www.example.com/foo/bar"},
		&followlinks.Module{
			Follow: &[]string{},
		},
		hook.Module{
			AdaptTransportFn: func(rt http.RoundTripper) http.RoundTripper {
				return flyscrape.MockTransport(200, `
				<a href="/baz">Baz</a>
				<a href="baz">Baz</a>
                <div class="next">
				    <a href="http://www.google.com">Google</a>
                </div>`)
			},
			ReceiveResponseFn: func(r *flyscrape.Response) {
				mu.Lock()
				urls = append(urls, r.Request.URL)
				mu.Unlock()
			},
		},
	}

	scraper := flyscrape.NewScraper()
	scraper.Modules = mods
	scraper.Run()

	require.Len(t, urls, 1)
	require.Contains(t, urls, "http://www.example.com/foo/bar")
}
