package flyscrape

import (
	"log"
	"strings"
	"sync"

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
	AllowedDomains []string `json:"allowed_domains"`
	Depth          int      `json:"depth"`
}

type ScrapeResult struct {
	URL  string
	Data M
	Err  error
}

type (
	M          map[string]any
	ScrapeFunc func(ScrapeParams) (M, error)
	FetchFunc  func(url string) (string, error)
)

type Service struct {
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
	data  M
	links []string
	err   error
}

func (s *Service) Scrape() <-chan ScrapeResult {
	if s.Concurrency == 0 {
		s.Concurrency = 1
	}

	jobs := make(chan target, 10)
	results := make(chan result)
	scraperesults := make(chan ScrapeResult)
	s.visited = hashmap.New[string, struct{}]()
	s.wg = &sync.WaitGroup{}

	for i := 0; i < s.Concurrency; i++ {
		go s.worker(i, jobs, results)
	}

	s.wg.Add(1)
	jobs <- target{url: s.ScrapeOptions.URL, depth: s.ScrapeOptions.Depth}

	go func() {
		s.wg.Wait()
		close(jobs)
		close(results)
	}()

	go func() {
		for res := range results {
			scraperesults <- ScrapeResult{
				URL:  res.url,
				Data: res.data,
				Err:  res.err,
			}
		}
		close(scraperesults)
	}()

	return scraperesults
}

func (s *Service) worker(id int, jobs chan target, results chan<- result) {
	for j := range jobs {
		res := s.process(j)

		for _, l := range res.links {
			if _, ok := s.visited.Get(l); ok {
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

		results <- res
		s.wg.Done()
	}
}

func (s *Service) process(job target) result {
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
