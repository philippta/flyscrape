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
}

func (Module) ModuleInfo() flyscrape.ModuleInfo {
	return flyscrape.ModuleInfo{
		ID:  "browser",
		New: func() flyscrape.Module { return new(Module) },
	}
}

func (m Module) AdaptTransport(t http.RoundTripper) http.RoundTripper {
	if !m.Browser {
		return t
	}

	headless := true
	if m.Headless != nil {
		headless = *m.Headless
	}

	ct, err := chromeTransport(headless)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	return ct
}

func chromeTransport(headless bool) (flyscrape.RoundTripFunc, error) {
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
			if h == "User-Agent" && strings.HasPrefix(r.UserAgent(), "flyscrape") {
				continue
			}
			page.MustSetExtraHeaders(h, r.Header.Get(h))
		}

		if err := page.Navigate(r.URL.String()); err != nil {
			return nil, err
		}

		if err := page.WaitStable(time.Second); err != nil {
			return nil, err
		}

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
	}, nil
}

var _ flyscrape.TransportAdapter = Module{}
