// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package flyscrape

import (
	"crypto/tls"
	"io"
	"net/http"
	"net/url"

	"github.com/cornelk/hashmap"
)

const userAgent = "flyscrape/0.1"

func ProxiedFetch(proxyURL string) FetchFunc {
	pu, err := url.Parse(proxyURL)
	if err != nil {
		panic("invalid proxy url")
	}

	client := http.Client{
		Transport: &http.Transport{
			Proxy:           http.ProxyURL(pu),
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	return func(url string) (string, error) {
		resp, err := client.Get(url)
		if err != nil {
			return "", err
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return "", err
		}

		html := string(body)
		return html, nil
	}
}

func CachedFetch(fetch FetchFunc) FetchFunc {
	cache := hashmap.New[string, string]()

	return func(url string) (string, error) {
		if html, ok := cache.Get(url); ok {
			return html, nil
		}

		html, err := fetch(url)
		if err != nil {
			return "", err
		}

		cache.Set(url, html)
		return html, nil
	}
}

func Fetch() FetchFunc {
	return func(url string) (string, error) {
		req, err := http.NewRequest(http.MethodGet, url, nil)
		if err != nil {
			return "", err
		}

		req.Header.Set("User-Agent", userAgent)

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return "", err
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return "", err
		}

		html := string(body)
		return html, nil
	}
}
