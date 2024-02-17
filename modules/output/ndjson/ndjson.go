// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package ndjson

import (
	"encoding/json"
	"io"
	"log"
	"os"
	"sync"
	"time"

	"github.com/philippta/flyscrape"
)

func init() {
	flyscrape.RegisterModule(Module{})
}

type Module struct {
	Output struct {
		Format string `json:"format"`
		File   string `json:"file"`
	} `json:"output"`

	w  io.WriteCloser
	mu *sync.Mutex
}

func (Module) ModuleInfo() flyscrape.ModuleInfo {
	return flyscrape.ModuleInfo{
		ID:  "output.ndjson",
		New: func() flyscrape.Module { return new(Module) },
	}
}

func (m *Module) Provision(ctx flyscrape.Context) {
	if m.disabled() {
		return
	}

	m.mu = &sync.Mutex{}

	if m.Output.File == "" {
		m.w = nopCloser{os.Stdout}
		return
	}

	f, err := os.Create(m.Output.File)
	if err != nil {
		log.Printf("failed to create file %q: %v", m.Output.File, err)
		os.Exit(1)
	}
	m.w = f
}

func (m *Module) ReceiveResponse(resp *flyscrape.Response) {
	if m.disabled() {
		return
	}

	if resp.Error == nil && resp.Data == nil {
		return
	}

	o := output{
		URL:       resp.Request.URL,
		Data:      resp.Data,
		Timestamp: time.Now(),
	}
	if resp.Error != nil {
		o.Error = resp.Error.Error()
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	enc := json.NewEncoder(m.w)
	enc.SetEscapeHTML(false)
	enc.Encode(o)
}

func (m *Module) Finalize() {
	if m.disabled() {
		return
	}
	m.w.Close()
}

func (m *Module) disabled() bool {
	return m.Output.Format != "ndjson"
}

type output struct {
	URL       string    `json:"url,omitempty"`
	Data      any       `json:"data,omitempty"`
	Error     string    `json:"error,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
}

type nopCloser struct {
	io.Writer
}

func (c nopCloser) Write(p []byte) (n int, err error) {
	return c.Writer.Write(p)
}

func (c nopCloser) Close() error {
	return nil
}

var (
	_ flyscrape.Provisioner      = (*Module)(nil)
	_ flyscrape.ResponseReceiver = (*Module)(nil)
	_ flyscrape.Finalizer        = (*Module)(nil)
)
