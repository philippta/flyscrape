// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package main

import (
	"flag"
	"fmt"

	"github.com/philippta/flyscrape"
)

type DevCommand struct{}

func (c *DevCommand) Run(args []string) error {
	fs := flag.NewFlagSet("flyscrape-dev", flag.ContinueOnError)
	fs.Usage = c.Usage

	if err := fs.Parse(args); err != nil {
		return err
	} else if fs.NArg() == 0 || fs.Arg(0) == "" {
		c.Usage()
		return flag.ErrHelp
	} else if fs.NArg() > 1 {
		return fmt.Errorf("too many arguments")
	}

	return flyscrape.Dev(fs.Arg(0))
}

func (c *DevCommand) Usage() {
	fmt.Println(`
The dev command watches the scraping script and re-runs it on any change.
Recursive scraping is disabled in this mode, only the initial URL will be scraped.

Usage:

    flyscrape dev SCRIPT

Examples:

    # Run and watch script.
    $ flyscrape dev example.js
`[1:])
}
