// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package flyscrape

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/evanw/esbuild/pkg/api"
	"github.com/inancgumus/screen"
	"github.com/tidwall/sjson"
)

func Run(file string, overrides map[string]any) error {
	src, err := os.ReadFile(file)
	if err != nil {
		return fmt.Errorf("failed to read script %q: %w", file, err)
	}

	client := &http.Client{}

	imports, wait := NewJSLibrary(client)
	defer wait()

	options := BuildOptions(file)

	exports, err := Compile(string(src), imports, options)
	if err != nil {
		return fmt.Errorf("failed to compile script: %w", err)
	}

	cfg := exports.Config()
	cfg = updateCfgMultiple(cfg, overrides)

	scraper := NewScraper()
	scraper.ScrapeFunc = exports.Scrape
	scraper.SetupFunc = exports.Setup
	scraper.Script = file
	scraper.Client = client
	scraper.Modules = LoadModules(cfg)

	scraper.Run()
	return nil
}

func Dev(file string, overrides map[string]any) error {
	cachefile, err := newCacheFile()
	if err != nil {
		return fmt.Errorf("failed to create cache file: %w", err)
	}

	trapsignal(func() {
		os.RemoveAll(cachefile)
	})

	fn := func(s string, file string) error {
		client := &http.Client{}

		imports, wait := NewJSLibrary(client)
		defer wait()

		options := BuildOptions(file)

		exports, err := Compile(s, imports, options)
		if err != nil {
			printCompileErr(file, err)
			return nil
		}

		cfg := exports.Config()
		cfg = updateCfgMultiple(cfg, overrides)
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
	newcfg, err := sjson.Set(string(cfg), key, value)
	if err != nil {
		return cfg
	}
	return Config(newcfg)
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

func updateCfgMultiple(cfg Config, updates map[string]any) Config {
	c := string(cfg)

	for k, v := range updates {
		nc, err := sjson.Set(c, k, v)
		if err != nil {
			continue
		}
		c = nc
	}

	return []byte(c)
}

func BuildOptions(fileName string) api.TransformOptions {
	options := api.TransformOptions{
		Platform: api.PlatformNode,
		Format:   api.FormatCommonJS,
		Loader:   api.LoaderJS,
	}

	if len(fileName) < 3 {
		return options
	}

	if fileName[len(fileName)-3:] == ".ts" {
		options.Loader = api.LoaderTS
	}
	return options
}
