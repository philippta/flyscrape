// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package cache

import (
	"bufio"
	"bytes"
	"net/http"
	"net/http/httputil"
	"path/filepath"
	"strings"

	"github.com/philippta/flyscrape"
)

func init() {
	flyscrape.RegisterModule(Module{})
}

type Module struct {
	Cache string `json:"cache"`

	store Store
}

func (Module) ModuleInfo() flyscrape.ModuleInfo {
	return flyscrape.ModuleInfo{
		ID:  "cache",
		New: func() flyscrape.Module { return new(Module) },
	}
}

func (m *Module) Provision(ctx flyscrape.Context) {
	switch {
	case m.Cache == "file":
		file := replaceExt(ctx.ScriptName(), ".cache")
		m.store = NewBoltStore(file)

	case strings.HasPrefix(m.Cache, "file:"):
		m.store = NewBoltStore(strings.TrimPrefix(m.Cache, "file:"))
	}
}

func (m *Module) AdaptTransport(t http.RoundTripper) http.RoundTripper {
	if m.store == nil {
		return t
	}
	return flyscrape.RoundTripFunc(func(r *http.Request) (*http.Response, error) {
		if nocache(r) {
			return t.RoundTrip(r)
		}

		key := r.Method + " " + r.URL.String()
		if b, ok := m.store.Get(key); ok {
			if resp, err := http.ReadResponse(bufio.NewReader(bytes.NewReader(b)), r); err == nil {
				return resp, nil
			}
		}

		resp, err := t.RoundTrip(r)
		if err != nil {
			return resp, err
		}

		// Avoid caching when running into rate limits or
		// when the page errored.
		if resp.StatusCode < 200 || resp.StatusCode > 299 {
			return resp, err
		}

		encoded, err := httputil.DumpResponse(resp, true)
		if err != nil {
			return resp, err
		}

		m.store.Set(key, encoded)
		return resp, nil
	})
}

func (m *Module) Finalize() {
	if v, ok := m.store.(interface{ Close() }); ok {
		v.Close()
	}
}

func nocache(r *http.Request) bool {
	if r.Header.Get(flyscrape.HeaderBypassCache) != "" {
		r.Header.Del(flyscrape.HeaderBypassCache)
		return true
	}
	return false
}

func replaceExt(filePath string, newExt string) string {
	ext := filepath.Ext(filePath)
	if ext != "" {
		fileNameWithoutExt := filePath[:len(filePath)-len(ext)]
		newFilePath := fileNameWithoutExt + newExt
		return newFilePath
	}
	return filePath + newExt
}

type Store interface {
	Get(key string) ([]byte, bool)
	Set(key string, value []byte)
}
