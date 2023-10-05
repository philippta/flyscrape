// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package flyscrape

import (
	"io"
	"log"
	"net/http"
	"net/http/cookiejar"
	"slices"
	"sync"

	"github.com/cornelk/hashmap"
)

type FetchFunc func(url string) (string, error)

type Context interface {
	Visit(url string)
	MarkVisited(url string)
	MarkUnvisited(url string)
	DisableModule(id string)
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

func NewScraper() *Scraper {
	return &Scraper{
		jobs:    make(chan target, 1024),
		visited: hashmap.New[string, struct{}](),
	}
}

type Scraper struct {
	ScrapeFunc ScrapeFunc

	wg      sync.WaitGroup
	jobs    chan target
	visited *hashmap.Map[string, struct{}]

	modules   []Module
	moduleIDs []string
	client    *http.Client
}

func (s *Scraper) LoadModule(mod Module) {
	id := mod.ModuleInfo().ID

	s.modules = append(s.modules, mod)
	s.moduleIDs = append(s.moduleIDs, id)
}

func (s *Scraper) DisableModule(id string) {
	idx := slices.Index(s.moduleIDs, id)
	if idx == -1 {
		return
	}
	s.modules = slices.Delete(s.modules, idx, idx+1)
	s.moduleIDs = slices.Delete(s.moduleIDs, idx, idx+1)
}

func (s *Scraper) Visit(url string) {
	s.enqueueJob(url, 0)
}

func (s *Scraper) MarkVisited(url string) {
	s.visited.Insert(url, struct{}{})
}

func (s *Scraper) MarkUnvisited(url string) {
	s.visited.Del(url)
}

func (s *Scraper) Run() {
	for _, mod := range s.modules {
		if v, ok := mod.(Provisioner); ok {
			v.Provision(s)
		}
	}

	s.initClient()
	go s.scrape()
	s.wg.Wait()
	close(s.jobs)

	for _, mod := range s.modules {
		if v, ok := mod.(Finalizer); ok {
			v.Finalize()
		}
	}
}

func (s *Scraper) initClient() {
	jar, _ := cookiejar.New(nil)
	s.client = &http.Client{Jar: jar}

	for _, mod := range s.modules {
		if v, ok := mod.(TransportAdapter); ok {
			s.client.Transport = v.AdaptTransport(s.client.Transport)
		}
	}
}

func (s *Scraper) scrape() {
	for job := range s.jobs {
		job := job
		go func() {
			s.process(job.url, job.depth)
			s.wg.Done()
		}()
	}
}

func (s *Scraper) process(url string, depth int) {
	request := &Request{
		Method:  http.MethodGet,
		URL:     url,
		Headers: defaultHeaders(),
		Cookies: s.client.Jar,
		Depth:   depth,
	}

	response := &Response{
		Request: request,
		Visit: func(url string) {
			s.enqueueJob(url, depth+1)
		},
	}

	for _, mod := range s.modules {
		if v, ok := mod.(RequestBuilder); ok {
			v.BuildRequest(request)
		}
	}

	req, err := http.NewRequest(request.Method, request.URL, nil)
	if err != nil {
		response.Error = err
		return
	}
	req.Header = request.Headers

	for _, mod := range s.modules {
		if v, ok := mod.(RequestValidator); ok {
			if !v.ValidateRequest(request) {
				return
			}
		}
	}

	defer func() {
		for _, mod := range s.modules {
			if v, ok := mod.(ResponseReceiver); ok {
				v.ReceiveResponse(response)
			}
		}
	}()

	resp, err := s.client.Do(req)
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

func defaultHeaders() http.Header {
	h := http.Header{}
	h.Set("User-Agent", "flyscrape/0.1")

	return h
}
