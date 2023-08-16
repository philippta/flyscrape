package flyscrape

import (
	"io"
	"net/http"

	"github.com/cornelk/hashmap"
)

func CachedFetch() FetchFunc {
	cache := hashmap.New[string, string]()

	return func(url string) (string, error) {
		if html, ok := cache.Get(url); ok {
			return html, nil
		}

		resp, err := http.Get(url)
		if err != nil {
			return "", err
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return "", err
		}

		html := string(body)
		cache.Set(url, html)
		return html, nil
	}
}

func Fetch() FetchFunc {
	return func(url string) (string, error) {
		resp, err := http.Get(url)
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
