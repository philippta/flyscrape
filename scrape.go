// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package flyscrape

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/cookiejar"
	"strings"
	"sync"

	"github.com/cornelk/hashmap"
)

type Context interface {
	ScriptName() string
	Visit(url string)
	MarkVisited(url string)
	MarkUnvisited(url string)
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
	return &Scraper{}
}

type Scraper struct {
	ScrapeFunc ScrapeFunc
	Script     string
	Modules    []Module
	Client     *http.Client

	wg      sync.WaitGroup
	jobs    chan target
	visited *hashmap.Map[string, struct{}]
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

func (s *Scraper) ScriptName() string {
	return s.Script
}

func (s *Scraper) Run() {
	s.jobs = make(chan target, 1<<20)
	s.visited = hashmap.New[string, struct{}]()

	s.initClient()

	for _, mod := range s.Modules {
		if v, ok := mod.(Provisioner); ok {
			v.Provision(s)
		}
	}

	for _, mod := range s.Modules {
		if v, ok := mod.(TransportAdapter); ok {
			s.Client.Transport = v.AdaptTransport(s.Client.Transport)
		}
	}

	s.scrape()
	s.wg.Wait()
	close(s.jobs)

	for _, mod := range s.Modules {
		if v, ok := mod.(Finalizer); ok {
			v.Finalize()
		}
	}
}

func (s *Scraper) initClient() {
	if s.Client == nil {
		s.Client = &http.Client{}
	}
	if s.Client.Jar == nil {
		s.Client.Jar, _ = cookiejar.New(nil)
	}
	if s.Client.Transport == nil {
		s.Client.Transport = http.DefaultTransport
	}
}

func (s *Scraper) scrape() {
	for i := 0; i < 500; i++ {
		go func() {
			for job := range s.jobs {
				s.process(job.url, job.depth)
				s.wg.Done()
			}
		}()
	}
}

func (s *Scraper) process(url string, depth int) {
	request := &Request{
		Method:  http.MethodGet,
		URL:     url,
		Headers: http.Header{},
		Cookies: s.Client.Jar,
		Depth:   depth,
	}

	response := &Response{
		Request: request,
		Visit: func(url string) {
			s.enqueueJob(url, depth+1)
		},
	}

	for _, mod := range s.Modules {
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

	for _, mod := range s.Modules {
		if v, ok := mod.(RequestValidator); ok {
			if !v.ValidateRequest(request) {
				return
			}
		}
	}

	defer func() {
		for _, mod := range s.Modules {
			if v, ok := mod.(ResponseReceiver); ok {
				v.ReceiveResponse(response)
			}
		}
	}()

	resp, err := s.Client.Do(req)
	if err != nil {
		response.Error = err
		return
	}
	defer resp.Body.Close()

	response.StatusCode = resp.StatusCode
	response.Headers = resp.Header

	if response.StatusCode < 200 || response.StatusCode >= 300 {
		response.Error = fmt.Errorf("%d %s", response.StatusCode, http.StatusText(response.StatusCode))
	}

	response.Body, err = io.ReadAll(resp.Body)
	if err != nil {
		response.Error = err
		return
	}

	if s.ScrapeFunc != nil {
		func() {
			defer func() {
				if r := recover(); r != nil {
					log.Println(r)
				}
			}()

			p := ScrapeParams{
				HTML:    string(response.Body),
				URL:     request.URL,
				Process: s.processImmediate,
				Follow: func(url string) {
					s.enqueueJob(url, depth+1)
				},
			}

			response.Data, err = s.ScrapeFunc(p)
			if err != nil {
				response.Error = err
				return
			}
		}()
	}
}

func (s *Scraper) processImmediate(url string) ([]byte, error) {
	request := &Request{
		Method:  http.MethodGet,
		URL:     url,
		Headers: http.Header{},
		Cookies: s.Client.Jar,
	}

	for _, mod := range s.Modules {
		if v, ok := mod.(RequestBuilder); ok {
			v.BuildRequest(request)
		}
	}

	req, err := http.NewRequest(request.Method, request.URL, nil)
	if err != nil {
		return nil, err
	}
	req.Header = request.Headers

	for _, mod := range s.Modules {
		if v, ok := mod.(RequestValidator); ok {
			if !v.ValidateRequest(request) {
				return nil, nil
			}
		}
	}

	resp, err := s.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("%d %s", resp.StatusCode, http.StatusText(resp.StatusCode))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func (s *Scraper) enqueueJob(url string, depth int) {
	url = strings.TrimSpace(url)
	if url == "" {
		return
	}

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
