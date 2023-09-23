// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package flyscrape

import (
	"io"
	"log"
	"net/http"
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

type Visitor interface {
	Visit(url string)
	MarkVisited(url string)
}

type (
	Request struct {
		URL   string
		Depth int
	}

	Response struct {
		ScrapeResult
		HTML  string
		Visit func(url string)
	}

	target struct {
		url   string
		depth int
	}
)

type Scraper struct {
	ScrapeFunc ScrapeFunc

	opts    Options
	wg      sync.WaitGroup
	jobs    chan target
	visited *hashmap.Map[string, struct{}]
	modules *hashmap.Map[string, Module]

	canRequestHandlers []func(url string, depth int) bool
	onRequestHandlers  []func(*Request)
	onResponseHandlers []func(*Response)
	onCompleteHandlers []func()
	transport          func(*http.Request) (*http.Response, error)
}

func NewScraper() *Scraper {
	s := &Scraper{
		jobs:    make(chan target, 1024),
		visited: hashmap.New[string, struct{}](),
		modules: hashmap.New[string, Module](),
		transport: func(r *http.Request) (*http.Response, error) {
			r.Header.Set("User-Agent", "flyscrape/0.1")
			return http.DefaultClient.Do(r)
		},
	}
	return s
}

func (s *Scraper) LoadModule(mod Module) {
	if v, ok := mod.(Transport); ok {
		s.SetTransport(v.Transport)
	}

	if v, ok := mod.(CanRequest); ok {
		s.CanRequest(v.CanRequest)
	}

	if v, ok := mod.(OnRequest); ok {
		s.OnRequest(v.OnRequest)
	}

	if v, ok := mod.(OnResponse); ok {
		s.OnResponse(v.OnResponse)
	}

	if v, ok := mod.(OnLoad); ok {
		v.OnLoad(s)
	}

	if v, ok := mod.(OnComplete); ok {
		s.OnComplete(v.OnComplete)
	}
}

func (s *Scraper) Visit(url string) {
	s.enqueueJob(url, 0)
}

func (s *Scraper) MarkVisited(url string) {
	s.visited.Insert(url, struct{}{})
}

func (s *Scraper) SetTransport(f func(r *http.Request) (*http.Response, error)) {
	s.transport = f
}

func (s *Scraper) CanRequest(f func(url string, depth int) bool) {
	s.canRequestHandlers = append(s.canRequestHandlers, f)
}

func (s *Scraper) OnRequest(f func(req *Request)) {
	s.onRequestHandlers = append(s.onRequestHandlers, f)
}

func (s *Scraper) OnResponse(f func(resp *Response)) {
	s.onResponseHandlers = append(s.onResponseHandlers, f)
}

func (s *Scraper) OnComplete(f func()) {
	s.onCompleteHandlers = append(s.onCompleteHandlers, f)
}

func (s *Scraper) Run() {
	go s.worker()
	s.wg.Wait()
	close(s.jobs)

	for _, handler := range s.onCompleteHandlers {
		handler()
	}
}

func (s *Scraper) worker() {
	for job := range s.jobs {
		go func(job target) {
			defer s.wg.Done()

			for _, handler := range s.canRequestHandlers {
				if !handler(job.url, job.depth) {
					return
				}
			}

			res, html := s.process(job)
			for _, handler := range s.onResponseHandlers {
				handler(&Response{
					ScrapeResult: res,
					HTML:         html,
					Visit: func(url string) {
						s.enqueueJob(url, job.depth+1)
					},
				})
			}
		}(job)
	}
}

func (s *Scraper) process(job target) (res ScrapeResult, html string) {
	res.URL = job.url
	res.Timestamp = time.Now()

	req, err := http.NewRequest(http.MethodGet, job.url, nil)
	if err != nil {
		res.Error = err
		return
	}

	for _, handler := range s.onRequestHandlers {
		handler(&Request{URL: job.url, Depth: job.depth})
	}

	resp, err := s.transport(req)
	if err != nil {
		res.Error = err
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		res.Error = err
		return
	}

	html = string(body)

	if s.ScrapeFunc != nil {
		res.Data, err = s.ScrapeFunc(ScrapeParams{HTML: html, URL: job.url})
		if err != nil {
			res.Error = err
			return
		}
	}

	return
}

func (s *Scraper) enqueueJob(url string, depth int) {
	if _, ok := s.visited.Get(url); ok {
		return
	}

	s.wg.Add(1)
	select {
	case s.jobs <- target{url: url, depth: depth}:
		s.MarkVisited(url)
	default:
		log.Println("queue is full, can't add url:", url)
		s.wg.Done()
	}
}

func ParseLinks(html string, origin string) []string {
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
