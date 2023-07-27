package flyscrape

import (
	"fmt"
	"io"
	"net/http"

	"flyscrape/scrape"
)

type Service struct{}

func (s *Service) ScrapeURL(url string, params map[string]any) (any, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("fetch %q: %w", url, err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read body %q: %w", url, err)
	}

	return scrape.Parse(string(body), params), nil
}
