// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package browser

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/go-rod/rod/lib/proto"
	"github.com/philippta/flyscrape"
)

func init() {
	flyscrape.RegisterModule(Module{})
}

type Module struct {
	Browser  bool  `json:"browser"`
	Headless *bool `json:"headless"`

	browser *rod.Browser
}

func (Module) ModuleInfo() flyscrape.ModuleInfo {
	return flyscrape.ModuleInfo{
		ID:  "browser",
		New: func() flyscrape.Module { return new(Module) },
	}
}

func (m *Module) AdaptTransport(t http.RoundTripper) http.RoundTripper {
	if !m.Browser {
		return t
	}

	headless := true
	if m.Headless != nil {
		headless = *m.Headless
	}

	browser, err := newBrowser(headless)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	m.browser = browser

	return chromeTransport(browser)
}

func (m *Module) Finalize() {
	if m.browser != nil {
		m.browser.Close()
	}
}

func newBrowser(headless bool) (*rod.Browser, error) {
	serviceURL, err := launcher.New().
		Headless(headless).
		Launch()
	if err != nil {
		return nil, fmt.Errorf("failed to launch browser: %w", err)
	}

	browser := rod.New().ControlURL(serviceURL).NoDefaultDevice()
	if err := browser.Connect(); err != nil {
		return nil, fmt.Errorf("failed to connect to browser: %w", err)
	}

	return browser, nil
}

func chromeTransport(browser *rod.Browser) flyscrape.RoundTripFunc {
	return func(r *http.Request) (*http.Response, error) {
		select {
		case <-r.Context().Done():
			return nil, r.Context().Err()
		default:
		}

		page := browser.MustPage()
		defer page.Close()

		var once sync.Once
		var networkResponse *proto.NetworkResponse
		go page.EachEvent(func(e *proto.NetworkResponseReceived) {
			if e.Type != proto.NetworkResourceTypeDocument {
				return
			}
			once.Do(func() {
				networkResponse = e.Response
			})
		})()

		page = page.Context(r.Context())

		for h := range r.Header {
			if h == "Cookie" {
				continue
			}
			if h == "User-Agent" && strings.HasPrefix(r.UserAgent(), "flyscrape") {
				continue
			}
			page.MustSetExtraHeaders(h, r.Header.Get(h))
		}

		page.SetCookies(parseCookies(r))

		if err := page.Navigate(r.URL.String()); err != nil {
			return nil, err
		}

		timeout := page.Timeout(10 * time.Second)
		timeout.WaitLoad()
		timeout.WaitDOMStable(300*time.Millisecond, 0)
		timeout.WaitRequestIdle(time.Second, nil, nil, nil)

		html, err := page.HTML()
		if err != nil {
			return nil, err
		}

		resp := &http.Response{
			StatusCode: 200,
			Status:     "200 OK",
			Body:       io.NopCloser(strings.NewReader(html)),
			Header:     http.Header{"Content-Type": []string{"text/html"}},
		}

		if networkResponse != nil {
			resp.StatusCode = networkResponse.Status
			resp.Status = networkResponse.StatusText
			resp.Header = http.Header{}

			for k, v := range networkResponse.Headers {
				resp.Header.Set(k, v.String())
			}
		}

		return resp, err
	}
}

func parseCookies(r *http.Request) []*proto.NetworkCookieParam {
	rawCookie := r.Header.Get("Cookie")
	if rawCookie == "" {
		return nil
	}

	header := http.Header{}
	header.Add("Cookie", rawCookie)
	request := http.Request{Header: header}

	domainSegs := strings.Split(r.URL.Hostname(), ".")
	if len(domainSegs) < 2 {
		return nil
	}

	domain := "." + strings.Join(domainSegs[len(domainSegs)-2:], ".")

	var cookies []*proto.NetworkCookieParam
	for _, cookie := range request.Cookies() {
		cookies = append(cookies, &proto.NetworkCookieParam{
			Name:     cookie.Name,
			Value:    cookie.Value,
			Domain:   domain,
			Path:     "/",
			Secure:   false,
			HTTPOnly: false,
			SameSite: "Lax",
			Expires:  -1,
			URL:      r.URL.String(),
		})
	}

	return cookies
}

var (
	_ flyscrape.TransportAdapter = &Module{}
	_ flyscrape.Finalizer        = &Module{}
)
