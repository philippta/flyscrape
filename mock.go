package flyscrape

import (
	"fmt"
	"io"
	"net/http"
	"strings"
)

func MockTransport(statusCode int, html string) func(*http.Request) (*http.Response, error) {
	return func(*http.Request) (*http.Response, error) {
		return MockResponse(statusCode, html)
	}
}

func MockResponse(statusCode int, html string) (*http.Response, error) {
	return &http.Response{
		StatusCode: statusCode,
		Status:     fmt.Sprintf("%d %s", statusCode, http.StatusText(statusCode)),
		Body:       io.NopCloser(strings.NewReader(html)),
	}, nil
}
