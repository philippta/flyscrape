// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package hook

import (
	"net/http"

	"github.com/philippta/flyscrape"
)

type Module struct {
	AdaptTransportFn  func(http.RoundTripper) http.RoundTripper
	ValidateRequestFn func(*flyscrape.Request) bool
	BuildRequestFn    func(*flyscrape.Request)
	ReceiveResponseFn func(*flyscrape.Response)
	ProvisionFn       func(flyscrape.Context)
	FinalizeFn        func()
}

func (Module) ModuleInfo() flyscrape.ModuleInfo {
	return flyscrape.ModuleInfo{
		ID:  "hook",
		New: func() flyscrape.Module { return new(Module) },
	}
}

func (m Module) AdaptTransport(t http.RoundTripper) http.RoundTripper {
	if m.AdaptTransportFn == nil {
		return t
	}
	return m.AdaptTransportFn(t)
}

func (m Module) ValidateRequest(r *flyscrape.Request) bool {
	if m.ValidateRequestFn == nil {
		return true
	}
	return m.ValidateRequestFn(r)
}

func (m Module) BuildRequest(r *flyscrape.Request) {
	if m.BuildRequestFn == nil {
		return
	}
	m.BuildRequestFn(r)
}

func (m Module) ReceiveResponse(r *flyscrape.Response) {
	if m.ReceiveResponseFn == nil {
		return
	}
	m.ReceiveResponseFn(r)
}

func (m Module) Provision(ctx flyscrape.Context) {
	if m.ProvisionFn == nil {
		return
	}
	m.ProvisionFn(ctx)
}

func (m Module) Finalize() {
	if m.FinalizeFn == nil {
		return
	}
	m.FinalizeFn()
}

var (
	_ flyscrape.TransportAdapter = Module{}
	_ flyscrape.RequestValidator = Module{}
	_ flyscrape.RequestBuilder   = Module{}
	_ flyscrape.ResponseReceiver = Module{}
	_ flyscrape.Provisioner      = Module{}
	_ flyscrape.Finalizer        = Module{}
)
