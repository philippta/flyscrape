// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/philippta/flyscrape"
	"github.com/philippta/flyscrape/modules/hook"
)

type RunCommand struct{}

func (c *RunCommand) Run(args []string) error {
	fs := flag.NewFlagSet("flyscrape-run", flag.ContinueOnError)
	fs.Usage = c.Usage

	if err := fs.Parse(args); err != nil {
		return err
	} else if fs.NArg() == 0 || fs.Arg(0) == "" {
		return fmt.Errorf("script path required")
	} else if fs.NArg() > 1 {
		return fmt.Errorf("too many arguments")
	}

	script := fs.Arg(0)
	src, err := os.ReadFile(script)
	if err != nil {
		return fmt.Errorf("failed to read script %q: %w", script, err)
	}

	cfg, scrape, err := flyscrape.Compile(string(src))
	if err != nil {
		return fmt.Errorf("failed to compile script: %w", err)
	}

	scraper := flyscrape.NewScraper()
	scraper.ScrapeFunc = scrape
	scraper.Script = script

	flyscrape.LoadModules(scraper, cfg)

	count := 0
	start := time.Now()

	scraper.LoadModule(hook.Module{
		ReceiveResponseFn: func(r *flyscrape.Response) {
			count++
		},
	})

	scraper.Run()

	log.Printf("Scraped %d websites in %v\n", count, time.Since(start))
	return nil
}

func (c *RunCommand) Usage() {
	fmt.Println(`
The run command runs the scraping script.

Usage:

    flyscrape run SCRIPT


Examples:

    # Run the script.
    $ flyscrape run example.js
`[1:])
}
