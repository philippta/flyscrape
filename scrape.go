// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package flyscrape

import (
	"log"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/cornelk/hashmap"
	"github.com/nlnwa/whatwg-url/url"
)

type ScrapeParams struct {
	HTML string
	URL  string
}

type ScrapeOptions struct {
	URL            string   `json:"url"`
	AllowedDomains []string `json:"allowedDomains"`
	BlockedDomains []string `json:"blockedDomains"`
	Depth          int      `json:"depth"`
	Rate           float64  `json:"rate"`
}

type ScrapeResult struct {
	URL       string    `json:"url"`
	Data      any       `json:"data,omitempty"`
	Links     []string  `json:"-"`
	Error     error     `json:"error,omitempty"`
	Timestamp time.Time `json:"timestamp"`
}

func (s *ScrapeResult) omit() bool {
	return s.Error == nil && s.Data == nil
}

type ScrapeFunc func(ScrapeParams) (any, error)

type FetchFunc func(url string) (string, error)

type target struct {
	url   string
	depth int
}

type Scraper struct {
	ScrapeOptions ScrapeOptions
	ScrapeFunc    ScrapeFunc
	FetchFunc     FetchFunc

	visited *hashmap.Map[string, struct{}]
	wg      *sync.WaitGroup
	jobs    chan target
	results chan ScrapeResult
}

func (s *Scraper) init() {
	s.visited = hashmap.New[string, struct{}]()
	s.wg = &sync.WaitGroup{}
	s.jobs = make(chan target, 1024)
	s.results = make(chan ScrapeResult)

	if s.FetchFunc == nil {
		s.FetchFunc = Fetch()
	}

	if s.ScrapeOptions.Rate == 0 {
		s.ScrapeOptions.Rate = 100
	}

	if len(s.ScrapeOptions.AllowedDomains) == 0 {
		u, err := url.Parse(s.ScrapeOptions.URL)
		if err == nil {
			s.ScrapeOptions.AllowedDomains = []string{u.Host()}
		}
	}
}

func (s *Scraper) Scrape() <-chan ScrapeResult {
	s.init()
	s.enqueueJob(s.ScrapeOptions.URL, s.ScrapeOptions.Depth)

	go s.worker()
	go s.waitClose()

	return s.results
}

func (s *Scraper) worker() {
	var (
		rate      = time.Duration(float64(time.Second) / s.ScrapeOptions.Rate)
		leakyjobs = leakychan(s.jobs, rate)
	)

	for job := range leakyjobs {
		go func(job target) {
			defer s.wg.Done()

			res := s.process(job)
			if !res.omit() {
				s.results <- res
			}

			if job.depth <= 0 {
				return
			}

			for _, l := range res.Links {
				if _, ok := s.visited.Get(l); ok {
					continue
				}

				if !s.isURLAllowed(l) {
					continue
				}

				s.enqueueJob(l, job.depth-1)
			}
		}(job)
	}
}

func (s *Scraper) process(job target) (res ScrapeResult) {
	res.URL = job.url
	res.Timestamp = time.Now()

	html, err := s.FetchFunc(job.url)
	if err != nil {
		res.Error = err
		return
	}

	res.Links = links(html, job.url)
	res.Data, err = s.ScrapeFunc(ScrapeParams{HTML: html, URL: job.url})
	if err != nil {
		res.Error = err
		return
	}

	return
}

func (s *Scraper) enqueueJob(url string, depth int) {
	s.wg.Add(1)
	select {
	case s.jobs <- target{url: url, depth: depth}:
		s.visited.Set(url, struct{}{})
	default:
		log.Println("queue is full, can't add url:", url)
		s.wg.Done()
	}
}

func (s *Scraper) isURLAllowed(rawurl string) bool {
	u, err := url.Parse(rawurl)
	if err != nil {
		return false
	}

	host := u.Host()
	ok := false

	for _, domain := range s.ScrapeOptions.AllowedDomains {
		if domain == "*" || host == domain {
			ok = true
			break
		}
	}

	for _, domain := range s.ScrapeOptions.BlockedDomains {
		if host == domain {
			ok = false
			break
		}
	}

	return ok
}

func (s *Scraper) waitClose() {
	s.wg.Wait()
	close(s.jobs)
	close(s.results)
}

func links(html string, origin string) []string {
	var links []string
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return nil
	}

	urlParser := url.NewParser(url.WithPercentEncodeSinglePercentSign())

	uniqueLinks := make(map[string]bool)
	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		link, _ := s.Attr("href")

		parsedLink, err := urlParser.ParseRef(origin, link)
		if err != nil || !isValidLink(parsedLink) {
			return
		}

		absLink := parsedLink.Href(true)

		if !uniqueLinks[absLink] {
			links = append(links, absLink)
			uniqueLinks[absLink] = true
		}
	})

	return links
}

func isValidLink(link *url.Url) bool {
	if link.Scheme() != "" && link.Scheme() != "http" && link.Scheme() != "https" {
		return false
	}

	if strings.HasPrefix(link.String(), "javascript:") {
		return false
	}

	return true
}

func leakychan[T any](in chan T, rate time.Duration) chan T {
	ticker := time.NewTicker(rate)
	sem := make(chan struct{}, 1)
	c := make(chan T)

	go func() {
		for range ticker.C {
			sem <- struct{}{}
		}
	}()

	go func() {
		for v := range in {
			<-sem
			c <- v
		}
		ticker.Stop()
		close(c)
	}()

	return c
}
