// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package cmd

import (
	"fmt"
	"runtime/debug"

	"github.com/philippta/flyscrape"
)

type VersionCommand struct{}

func (c *VersionCommand) Run(args []string) error {
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return fmt.Errorf("no build info found")
	}

	var os, arch, version string
	for _, setting := range info.Settings {
		switch setting.Key {
		case "GOARCH":
			arch = setting.Value
		case "GOOS":
			os = setting.Value
		case "vcs.revision":
			version = "v0.0.0-" + setting.Value
		}
	}

	if flyscrape.Version != "" {
		version = flyscrape.Version
	}

	fmt.Printf("flyscrape %s %s/%s\n", version, os, arch)
	return nil
}
