// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package retry

import (
	"errors"
	"io"
	"net"
	"net/http"
	"slices"
	"strconv"
	"time"

	"github.com/philippta/flyscrape"
)

func init() {
	flyscrape.RegisterModule(Module{})
}

type Module struct {
	ticker    *time.Ticker
	semaphore chan struct{}

	RetryDelays []time.Duration
}

func (Module) ModuleInfo() flyscrape.ModuleInfo {
	return flyscrape.ModuleInfo{
		ID:  "retry",
		New: func() flyscrape.Module { return new(Module) },
	}
}

func (m *Module) Provision(flyscrape.Context) {
	if m.RetryDelays == nil {
		m.RetryDelays = defaultRetryDelays
	}
}

func (m *Module) AdaptTransport(t http.RoundTripper) http.RoundTripper {
	return flyscrape.RoundTripFunc(func(r *http.Request) (*http.Response, error) {
		resp, err := t.RoundTrip(r)
		if !shouldRetry(resp, err) {
			return resp, err
		}

		for _, delay := range m.RetryDelays {
			drainBody(resp, err)

			time.Sleep(retryAfter(resp, delay))

			resp, err = t.RoundTrip(r)
			if !shouldRetry(resp, err) {
				break
			}
		}

		return resp, err
	})
}

func shouldRetry(resp *http.Response, err error) bool {
	statusCodes := []int{
		http.StatusForbidden,
		http.StatusRequestTimeout,
		http.StatusTooEarly,
		http.StatusTooManyRequests,
		http.StatusInternalServerError,
		http.StatusBadGateway,
		http.StatusServiceUnavailable,
		http.StatusGatewayTimeout,
	}

	if resp != nil {
		if slices.Contains(statusCodes, resp.StatusCode) {
			return true
		}
	}
	if err == nil {
		return false
	}
	if _, ok := err.(net.Error); ok {
		return true
	}
	if errors.Is(err, io.ErrUnexpectedEOF) {
		return true
	}

	return false
}

func drainBody(resp *http.Response, err error) {
	if err == nil && resp != nil && resp.Body != nil {
		io.Copy(io.Discard, resp.Body)
		resp.Body.Close()
	}
}

func retryAfter(resp *http.Response, fallback time.Duration) time.Duration {
	if resp == nil {
		return fallback
	}

	timeexp := resp.Header.Get("Retry-After")
	if timeexp == "" {
		return fallback
	}

	if seconds, err := strconv.Atoi(timeexp); err == nil {
		return time.Duration(seconds) * time.Second
	}

	formats := []string{
		time.RFC1123, // HTTP Spec
		time.RFC1123Z,
		time.ANSIC,
		time.UnixDate,
		time.RubyDate,
		time.RFC822,
		time.RFC822Z,
		time.RFC850,
		time.RFC3339,
	}
	for _, format := range formats {
		if t, err := time.Parse(format, timeexp); err == nil {
			return t.Sub(time.Now())
		}
	}

	return fallback
}

var defaultRetryDelays = []time.Duration{
	1 * time.Second,
	2 * time.Second,
	5 * time.Second,
	10 * time.Second,
}

var (
	_ flyscrape.TransportAdapter = (*Module)(nil)
	_ flyscrape.Provisioner      = (*Module)(nil)
)
