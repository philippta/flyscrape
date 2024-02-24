// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package flyscrape

import (
	"encoding/json"
	"net/http"
	"sync"
)

type Module interface {
	ModuleInfo() ModuleInfo
}

type ModuleInfo struct {
	ID  string
	New func() Module
}

type TransportAdapter interface {
	AdaptTransport(http.RoundTripper) http.RoundTripper
}

type RequestValidator interface {
	ValidateRequest(*Request) bool
}

type RequestBuilder interface {
	BuildRequest(*Request)
}

type ResponseReceiver interface {
	ReceiveResponse(*Response)
}

type Provisioner interface {
	Provision(Context)
}

type Finalizer interface {
	Finalize()
}

func RegisterModule(mod Module) {
	modulesMu.Lock()
	defer modulesMu.Unlock()

	id := mod.ModuleInfo().ID
	if _, ok := modules[id]; ok {
		panic("module with id: " + id + " already registered")
	}
	modules[mod.ModuleInfo().ID] = mod
}

func LoadModules(cfg Config) []Module {
	modulesMu.RLock()
	defer modulesMu.RUnlock()

	loaded := map[string]struct{}{}
	mods := []Module{}

	// load standard modules in order
	for _, id := range moduleOrder {
		mod := modules[id].ModuleInfo().New()
		if err := json.Unmarshal(cfg, mod); err != nil {
			panic("failed to decode config: " + err.Error())
		}
		mods = append(mods, mod)
		loaded[id] = struct{}{}
	}

	// load custom modules
	for id := range modules {
		if _, ok := loaded[id]; ok {
			continue
		}
		mod := modules[id].ModuleInfo().New()
		if err := json.Unmarshal(cfg, mod); err != nil {
			panic("failed to decode config: " + err.Error())
		}
		mods = append(mods, mod)
		loaded[id] = struct{}{}
	}

	return mods
}

var (
	modules   = map[string]Module{}
	modulesMu sync.RWMutex

	moduleOrder = []string{
		// Transport adapters must be loaded in a specific order.
		// All other modules can be loaded in any order.
		"proxy",
		"browser",
		"retry",
		"ratelimit",
		"cache",
		"cookies",
		"headers",
	}
)
