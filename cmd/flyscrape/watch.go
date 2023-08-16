package main

import (
	"flag"
	"fmt"
	"log"

	"flyscrape"

	"github.com/inancgumus/screen"
)

type WatchCommand struct{}

func (c *WatchCommand) Run(args []string) error {
	fs := flag.NewFlagSet("flyscrape-watch", flag.ContinueOnError)
	fs.Usage = c.Usage

	if err := fs.Parse(args); err != nil {
		return err
	} else if fs.NArg() == 0 || fs.Arg(0) == "" {
		return fmt.Errorf("script path required")
	} else if fs.NArg() > 1 {
		return fmt.Errorf("too many arguments")
	}

	fetch := flyscrape.CachedFetch()
	script := fs.Arg(0)

	flyscrape.Watch(script, func(s string) error {
		opts, scrape, err := flyscrape.Compile(s)
		if err != nil {
			log.Println(err)
			// ignore compilation errors
			return nil
		}

		opts.Depth = 0
		scr := flyscrape.Scraper{
			ScrapeOptions: opts,
			ScrapeFunc:    scrape,
			FetchFunc:     fetch,
		}

		result := <-scr.Scrape()
		if result.Error != nil {
			log.Println(result.Error)
			return nil
		}

		screen.Clear()
		screen.MoveTopLeft()
		flyscrape.PrettyPrint(result)
		return nil
	})

	return nil
}

func (c *WatchCommand) Usage() {
	fmt.Println(`
The watch command watches the scraping script and re-runs it on any change.
Recursive scraping is disabled in this mode, only the initial URL will be scraped.

Usage:

    flyscrape watch SCRIPT


Examples:

    # Run and watch script.
    $ flyscrape watch example.js
`[1:])
}
