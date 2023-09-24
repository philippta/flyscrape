// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package flyscrape

import (
	"io"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
	"sync"

	gourl "net/url"

	"github.com/PuerkitoBio/goquery"
	"github.com/cornelk/hashmap"
)

type ScrapeParams struct {
	HTML string
	URL  string
}

type ScrapeFunc func(ScrapeParams) (any, error)

type FetchFunc func(url string) (string, error)

type Visitor interface {
	Visit(url string)
	MarkVisited(url string)
}

type Request struct {
	Method  string
	URL     string
	Headers http.Header
	Cookies http.CookieJar
	Depth   int
}

type Response struct {
	StatusCode int
	Headers    http.Header
	Body       []byte
	Data       any
	Error      error
	Request    *Request

	Visit func(url string)
}

type target struct {
	url   string
	depth int
}

type Scraper struct {
	ScrapeFunc ScrapeFunc

	cfg       Config
	wg        sync.WaitGroup
	jobs      chan target
	visited   *hashmap.Map[string, struct{}]
	modules   *hashmap.Map[string, Module]
	cookieJar *cookiejar.Jar

	canRequestHandlers []func(url string, depth int) bool
	onRequestHandlers  []func(*Request)
	onResponseHandlers []func(*Response)
	onCompleteHandlers []func()
	transport          func(*http.Request) (*http.Response, error)
}

func NewScraper() *Scraper {
	jar, _ := cookiejar.New(nil)
	s := &Scraper{
		jobs:    make(chan target, 1024),
		visited: hashmap.New[string, struct{}](),
		modules: hashmap.New[string, Module](),
		transport: func(r *http.Request) (*http.Response, error) {
			r.Header.Set("User-Agent", "flyscrape/0.1")
			return http.DefaultClient.Do(r)
		},
		cookieJar: jar,
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

			s.process(job.url, job.depth)
		}(job)
	}
}

func (s *Scraper) process(url string, depth int) {
	request := &Request{
		Method:  http.MethodGet,
		URL:     url,
		Headers: http.Header{},
		Cookies: s.cookieJar,
	}

	response := &Response{
		Request: request,
		Visit: func(url string) {
			s.enqueueJob(url, depth+1)
		},
	}

	defer func() {
		for _, handler := range s.onResponseHandlers {
			handler(response)
		}
	}()

	req, err := http.NewRequest(request.Method, request.URL, nil)
	if err != nil {
		response.Error = err
		return
	}
	req.Header = request.Headers

	for _, handler := range s.onRequestHandlers {
		handler(request)
	}

	resp, err := s.transport(req)
	if err != nil {
		response.Error = err
		return
	}
	defer resp.Body.Close()

	response.StatusCode = resp.StatusCode
	response.Headers = resp.Header

	response.Body, err = io.ReadAll(resp.Body)
	if err != nil {
		response.Error = err
		return
	}

	if s.ScrapeFunc != nil {
		response.Data, err = s.ScrapeFunc(ScrapeParams{HTML: string(response.Body), URL: request.URL})
		if err != nil {
			response.Error = err
			return
		}
	}
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

	originurl, err := url.Parse(origin)
	if err != nil {
		return nil
	}

	uniqueLinks := make(map[string]bool)
	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		link, _ := s.Attr("href")

		parsedLink, err := originurl.Parse(link)

		if err != nil || !isValidLink(parsedLink) {
			return
		}

		absLink := parsedLink.String()

		if !uniqueLinks[absLink] {
			links = append(links, absLink)
			uniqueLinks[absLink] = true
		}
	})

	return links
}

func isValidLink(link *gourl.URL) bool {
	if link.Scheme != "" && link.Scheme != "http" && link.Scheme != "https" {
		return false
	}

	return true
}
