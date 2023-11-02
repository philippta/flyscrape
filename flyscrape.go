// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package flyscrape

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/inancgumus/screen"
)

func Run(file string) error {
	src, err := os.ReadFile(file)
	if err != nil {
		return fmt.Errorf("failed to read script %q: %w", file, err)
	}

	client := &http.Client{}

	exports, err := Compile(string(src), NewJSLibrary(client))
	if err != nil {
		return fmt.Errorf("failed to compile script: %w", err)
	}

	scraper := NewScraper()
	scraper.ScrapeFunc = exports.Scrape
	scraper.SetupFunc = exports.Setup
	scraper.Script = file
	scraper.Client = client
	scraper.Modules = LoadModules(exports.Config())

	scraper.Run()
	return nil
}

func Dev(file string) error {
	cachefile, err := newCacheFile()
	if err != nil {
		return fmt.Errorf("failed to create cache file: %w", err)
	}

	trapsignal(func() {
		os.RemoveAll(cachefile)
	})

	fn := func(s string) error {
		client := &http.Client{}

		exports, err := Compile(s, NewJSLibrary(client))
		if err != nil {
			printCompileErr(file, err)
			return nil
		}

		cfg := exports.Config()
		cfg = updateCfg(cfg, "depth", 0)
		cfg = updateCfg(cfg, "cache", "file:"+cachefile)

		scraper := NewScraper()
		scraper.ScrapeFunc = exports.Scrape
		scraper.SetupFunc = exports.Setup
		scraper.Script = file
		scraper.Client = client
		scraper.Modules = LoadModules(cfg)

		screen.Clear()
		screen.MoveTopLeft()
		scraper.Run()

		return nil
	}

	if err := Watch(file, fn); err != nil && err != StopWatch {
		return fmt.Errorf("failed to watch script %q: %w", file, err)
	}
	return nil
}

func printCompileErr(script string, err error) {
	screen.Clear()
	screen.MoveTopLeft()

	if errs, ok := err.(interface{ Unwrap() []error }); ok {
		for _, err := range errs.Unwrap() {
			log.Printf("%s:%v\n", script, err)
		}
	} else {
		log.Println(err)
	}
}

func updateCfg(cfg Config, key string, value any) Config {
	var m map[string]any
	if err := json.Unmarshal(cfg, &m); err != nil {
		return cfg
	}

	m[key] = value

	b, err := json.Marshal(m)
	if err != nil {
		return cfg
	}

	return b
}

func newCacheFile() (string, error) {
	cachedir, err := os.MkdirTemp("", "flyscrape-cache")
	if err != nil {
		return "", err
	}
	return filepath.Join(cachedir, "dev.cache"), nil
}

func trapsignal(f func()) {
	sig := make(chan os.Signal, 2)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-sig
		f()
		os.Exit(0)
	}()
}
