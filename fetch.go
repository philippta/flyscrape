// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package flyscrape

import (
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/cornelk/hashmap"
)

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

	return func(req *http.Request) (string, error) {
		resp, err := client.Do(req)
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

	return func(req *http.Request) (string, error) {
		reqKey := fmt.Sprintf("%s %s", req.Method, req.URL.String())
		if html, ok := cache.Get(reqKey); ok {
			return html, nil
		}

		html, err := fetch(req)
		if err != nil {
			return "", err
		}

		cache.Set(reqKey, html)
		return html, nil
	}
}

func Fetch() FetchFunc {
	return func(req *http.Request) (string, error) {
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
