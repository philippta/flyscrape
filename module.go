package flyscrape

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
)

type Module interface {
	ID() string
}

type Transport interface {
	Transport(*http.Request) (*http.Response, error)
}

type CanRequest interface {
	CanRequest(url string, depth int) bool
}

type OnRequest interface {
	OnRequest(*Request)
}
type OnResponse interface {
	OnResponse(*Response)
}

type OnLoad interface {
	OnLoad(Visitor)
}

type OnComplete interface {
	OnComplete()
}

func RegisterModule(m Module) {
	id := m.ID()
	if id == "" {
		panic("module id is missing")
	}

	globalModulesMu.Lock()
	defer globalModulesMu.Unlock()

	if _, ok := globalModules[id]; ok {
		panic(fmt.Sprintf("module %s already registered", id))
	}
	globalModules[id] = m
}

func LoadModules(s *Scraper, opts Options) {
	globalModulesMu.RLock()
	defer globalModulesMu.RUnlock()

	for _, mod := range globalModules {
		json.Unmarshal(opts, mod)
		s.LoadModule(mod)
	}
}

var (
	globalModules   = map[string]Module{}
	globalModulesMu sync.RWMutex
)
