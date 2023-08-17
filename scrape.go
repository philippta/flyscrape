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
	URL   string   `json:"url"`
	Data  any      `json:"data,omitempty"`
	Links []string `json:"-"`
	Error error    `json:"error,omitempty"`
}

type (
	ScrapeFunc func(ScrapeParams) (any, error)
	FetchFunc  func(url string) (string, error)
)

type Scraper struct {
	ScrapeOptions ScrapeOptions
	ScrapeFunc    ScrapeFunc
	FetchFunc     FetchFunc
	Concurrency   int

	visited *hashmap.Map[string, struct{}]
	wg      *sync.WaitGroup
}

type target struct {
	url   string
	depth int
}

type result struct {
	url   string
	data  any
	links []string
	err   error
}

func (s *Scraper) Scrape() <-chan ScrapeResult {
	if s.Concurrency == 0 {
		s.Concurrency = 1
	}
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

	jobs := make(chan target, 1024)
	results := make(chan result)
	scraperesults := make(chan ScrapeResult)
	s.visited = hashmap.New[string, struct{}]()
	s.wg = &sync.WaitGroup{}

	go s.worker(jobs, results)

	s.wg.Add(1)
	s.visited.Set(s.ScrapeOptions.URL, struct{}{})
	jobs <- target{url: s.ScrapeOptions.URL, depth: s.ScrapeOptions.Depth}

	go func() {
		s.wg.Wait()
		close(jobs)
		close(results)
	}()

	go func() {
		for res := range results {
			scraperesults <- ScrapeResult{
				URL:   res.url,
				Data:  res.data,
				Links: res.links,
				Error: res.err,
			}
		}
		close(scraperesults)
	}()

	return scraperesults
}

func (s *Scraper) worker(jobs chan target, results chan<- result) {
	rate := time.Duration(float64(time.Second) / s.ScrapeOptions.Rate)
	for j := range leakychan(jobs, rate) {
		j := j
		go func() {
			res := s.process(j)

			if j.depth > 0 {
				for _, l := range res.links {
					if _, ok := s.visited.Get(l); ok {
						continue
					}

					if !s.isURLAllowed(l) {
						continue
					}

					s.wg.Add(1)
					select {
					case jobs <- target{url: l, depth: j.depth - 1}:
						s.visited.Set(l, struct{}{})
					default:
						log.Println("queue is full, can't add url:", l)
						s.wg.Done()
					}
				}
			}

			if res.err != nil || res.data != nil {
				results <- res
			}
			s.wg.Done()
		}()
	}
}

func (s *Scraper) process(job target) result {
	html, err := s.FetchFunc(job.url)
	if err != nil {
		return result{url: job.url, err: err}
	}

	links := Links(html, job.url)
	data, err := s.ScrapeFunc(ScrapeParams{HTML: html, URL: job.url})
	if err != nil {
		return result{url: job.url, links: links, err: err}
	}

	return result{url: job.url, data: data, links: links}
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

func Links(html string, origin string) []string {
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
