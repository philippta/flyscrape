// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package mongodb

import (
	"context"
	"errors"
	"log"
	"os"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"

	"github.com/philippta/flyscrape"
)

var (
	DefaultMaxPoolSize   = 100
	DefaultBatchSize     = 100
	DefaultFlushInterval = 10 * time.Second
	DefaultTimeout       = 30 * time.Second
	DefaultMaxRetries    = 3
)

func init() {
	flyscrape.RegisterModule(&Module{})
}

type Module struct {
	SourceId string `json:"sourceId"`
	Output   struct {
		MongoDB struct {
			URI         string `json:"uri"`
			Database    string `json:"database"`
			Collection  string `json:"collection"`
			MaxPoolSize int    `json:"maxPoolSize,omitempty"`
		} `json:"mongodb"`
	} `json:"output"`
	Concurrency int `json:"concurrency"`

	client      *mongo.Client
	collection  *mongo.Collection
	maxPoolSize int

	buf []interface{}
	mu  *sync.Mutex

	ticker      *time.Ticker
	done        chan struct{}
	concurrency chan struct{}
}

func (m *Module) ModuleInfo() flyscrape.ModuleInfo {
	return flyscrape.ModuleInfo{
		ID:  "output.mongodb",
		New: func() flyscrape.Module { return new(Module) },
	}
}

func (m *Module) Provision(ctx flyscrape.Context) {
	if m.disabled() {
		return
	}

	m.mu = &sync.Mutex{}

	m.maxPoolSize = DefaultMaxPoolSize
	if m.Output.MongoDB.MaxPoolSize != 0 {
		m.maxPoolSize = m.Output.MongoDB.MaxPoolSize
	}

	if m.concurrencyEnabled() {
		m.concurrency = make(chan struct{}, m.Concurrency)
		for i := 0; i < m.Concurrency; i++ {
			m.concurrency <- struct{}{}
		}
	}

	ctxTimeout, cancel := context.WithTimeout(context.Background(), DefaultTimeout)
	defer cancel()
	client, err := mongo.Connect(ctxTimeout, options.Client().ApplyURI(m.Output.MongoDB.URI).SetMaxPoolSize(uint64(m.maxPoolSize)))

	if err != nil {
		log.Printf("failed to connect to MongoDB: %v", err)
		os.Exit(1)
	}

	if err := client.Ping(ctxTimeout, readpref.Primary()); err != nil {
		log.Printf("failed to ping MongoDB: %v", err)
		os.Exit(1)
	}

	m.client = client
	m.collection = client.Database(m.Output.MongoDB.Database).Collection(m.Output.MongoDB.Collection)
	m.buf = make([]interface{}, 0, DefaultBatchSize)
	m.done = make(chan struct{})
	m.ticker = time.NewTicker(DefaultFlushInterval)

	go func() {
		for {
			select {
			case <-m.ticker.C:
				m.flushBuffer()
			case <-m.done:
				return
			}
		}
	}()
}

func (m *Module) ReceiveResponse(resp *flyscrape.Response) {
	if m.disabled() {
		return
	}

	if resp.Data == nil && resp.Error == nil {
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

	m.buf = append(m.buf, o)
	if len(m.buf) >= DefaultBatchSize {
		go m.flushBuffer()
	}
}

func (m *Module) Finalize() {
	if m.disabled() {
		return
	}

	m.ticker.Stop()

	ctxTimeout, cancel := context.WithTimeout(context.Background(), DefaultTimeout)
	defer cancel()

	done := make(chan struct{})
	go func() {
		m.flushBuffer()
		close(done)
	}()
	select {
	case <-done:
	case <-ctxTimeout.Done():
	}

	if err := m.client.Disconnect(ctxTimeout); err != nil {
		log.Printf("failed to disconnect from MongoDB: %v", err)
	}

	close(m.done)
}

func (m *Module) disabled() bool {
	return m.Output.MongoDB.URI == "" || m.Output.MongoDB.Database == "" || m.Output.MongoDB.Collection == ""
}

type output struct {
	URL       string    `json:"url,omitempty"`
	Data      any       `json:"data,omitempty"`
	Error     string    `json:"error,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
}

func (m *Module) flushBuffer() {
	if m.concurrencyEnabled() {
		<-m.concurrency
		defer func() { m.concurrency <- struct{}{} }()
	}
	m.mu.Lock()
	defer m.mu.Unlock()

	if len(m.buf) == 0 {
		return
	}

	var err error
	var res *mongo.InsertManyResult

	for i := 0; i < DefaultMaxRetries; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), DefaultTimeout)
		defer cancel()

		log.Printf("attempt %d to insert %d documents to MongoDB", i+1, len(m.buf))
		res, err = m.collection.InsertMany(ctx, m.buf)
		if err == nil {
			log.Printf("successfully inserted %d documents to MongoDB", len(res.InsertedIDs))
			m.buf = m.buf[:0]
			return
		}

		if errors.Is(err, context.DeadlineExceeded) {
			log.Printf("operation timed out, retrying...")
			continue
		}

		log.Printf("failed to insert documents to MongoDB: %v", err)
		break
	}

	if err != nil {
		log.Printf("failed to insert %d documents after %d retries: %v", len(m.buf), DefaultMaxRetries, err)
	}
	m.buf = m.buf[:0]
}

func (m *Module) concurrencyEnabled() bool {
	return m.Concurrency > 0
}

var (
	_ flyscrape.Provisioner      = (*Module)(nil)
	_ flyscrape.ResponseReceiver = (*Module)(nil)
	_ flyscrape.Finalizer        = (*Module)(nil)
)
