// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/philippta/flyscrape"
)

type NewCommand struct{}

func (c *NewCommand) Run(args []string) error {
	fs := flag.NewFlagSet("flyscrape-new", flag.ContinueOnError)
	fs.Usage = c.Usage

	if err := fs.Parse(args); err != nil {
		return err
	} else if fs.NArg() == 0 || fs.Arg(0) == "" {
		c.Usage()
		return flag.ErrHelp
	} else if fs.NArg() > 1 {
		return fmt.Errorf("too many arguments")
	}

	script := fs.Arg(0)
	if _, err := os.Stat(script); err == nil {
		return fmt.Errorf("script already exists")
	}

	if err := os.WriteFile(script, flyscrape.ScriptTemplate, 0o644); err != nil {
		return fmt.Errorf("failed to create script %q: %w", script, err)
	}

	fmt.Printf("Scraping script %v created.\n", script)
	return nil
}

func (c *NewCommand) Usage() {
	fmt.Println(`
The new command creates a new scraping script.

Usage:

    flyscrape new SCRIPT

Examples:

    # Create a new scraping script.
    $ flyscrape new example.js
`[1:])
}
