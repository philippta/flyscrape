package flyscrape

import (
	"encoding/json"
	"net/http"
)

type Module any

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

func RegisterModule(mod Module) {
	globalModules = append(globalModules, mod)
}

func LoadModules(s *Scraper, cfg Config) {
	for _, mod := range globalModules {
		json.Unmarshal(cfg, mod)
		s.LoadModule(mod)
	}
}

var globalModules = []Module{}
