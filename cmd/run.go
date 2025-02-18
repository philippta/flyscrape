// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package cmd

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
	}

	cfg, err := parseConfigArgs(fs.Args()[1:])
	if err != nil {
		return fmt.Errorf("error parsing config flags: %w", err)
	}

	return flyscrape.Run(fs.Arg(0), cfg)
}

func (c *RunCommand) Usage() {
	fmt.Println(`
The run command runs the scraping script.

Usage:

    flyscrape run SCRIPT [config flags]

Examples:

    # Run the script.
    $ flyscrape run example.js

    # Set the URL as argument.
    $ flyscrape run example.js --url "http://other.com"

    # Enable proxy support.
    $ flyscrape run example.js --proxies "http://someproxy:8043"

    # Follow paginated links.
    $ flyscrape run example.js --depth 5 --follow ".next-button > a"

    # Set the output format to ndjson.
    $ flyscrape run example.js --output.format ndjson

    # Write the output to a file.
    $ flyscrape run example.js --output.file results.json
`[1:])
}
