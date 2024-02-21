// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package main

import (
	_ "embed"
	"flag"
	"log"
	"os"

	"github.com/philippta/flyscrape/cmd"
	_ "github.com/philippta/flyscrape/modules/cache"
	_ "github.com/philippta/flyscrape/modules/cookies"
	_ "github.com/philippta/flyscrape/modules/depth"
	_ "github.com/philippta/flyscrape/modules/domainfilter"
	_ "github.com/philippta/flyscrape/modules/followlinks"
	_ "github.com/philippta/flyscrape/modules/headers"
	_ "github.com/philippta/flyscrape/modules/output/json"
	_ "github.com/philippta/flyscrape/modules/output/ndjson"
	_ "github.com/philippta/flyscrape/modules/proxy"
	_ "github.com/philippta/flyscrape/modules/ratelimit"
	_ "github.com/philippta/flyscrape/modules/retry"
	_ "github.com/philippta/flyscrape/modules/starturl"
	_ "github.com/philippta/flyscrape/modules/urlfilter"
)

func main() {
	log.SetFlags(0)

	if err := (&cmd.Main{}).Run(os.Args[1:]); err != nil {
		if err != flag.ErrHelp {
			log.Println(err)
		}
		os.Exit(1)
	}
}
