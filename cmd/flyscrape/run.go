// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package main

import (
	"flag"
	"fmt"

	"github.com/philippta/flyscrape"
)

type RunCommand struct{}

func (c *RunCommand) Run(args []string) error {
	fs := flag.NewFlagSet("flyscrape-run", flag.ContinueOnError)
	fs.Usage = c.Usage

	if err := fs.Parse(args); err != nil {
		return err
	} else if fs.NArg() == 0 || fs.Arg(0) == "" {
		c.Usage()
		return flag.ErrHelp
	} else if fs.NArg() > 1 {
		return fmt.Errorf("too many arguments")
	}

	return flyscrape.Run(fs.Arg(0))
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
