// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package headers

import (
	_ "embed"
	"fmt"
	"math/rand"
	"strings"
)

//go:generate bash -c "flyscrape run ../../examples/useragents/chrome.js | jq -r '[.[].data] | flatten | .[]' | sort -nr | uniq > versions_chrome.txt"
//go:generate bash -c "flyscrape run ../../examples/useragents/firefox.js | jq -r '[.[].data] | flatten | .[]' | sort -nr | uniq > versions_firefox.txt"
//go:generate bash -c "flyscrape run ../../examples/useragents/edge.js | jq -r '[.[].data] | flatten | .[]' | sort -nr | uniq > versions_edge.txt"
//go:generate bash -c "flyscrape run ../../examples/useragents/opera.js | jq -r '[.[].data] | flatten | .[]' | sort -nr | uniq > versions_opera.txt"

//go:embed versions_chrome.txt
var versionsChromeRaw string
var versionsChrome = strings.Split("\n", strings.TrimSpace(versionsChromeRaw))

//go:embed versions_firefox.txt
var versionsFirefoxRaw string
var versionsFirefox = strings.Split("\n", strings.TrimSpace(versionsFirefoxRaw))

//go:embed versions_edge.txt
var versionsEdgeRaw string
var versionsEdge = strings.Split("\n", strings.TrimSpace(versionsEdgeRaw))

//go:embed versions_opera.txt
var versionsOperaRaw string
var versionsOpera = strings.Split("\n", strings.TrimSpace(versionsOperaRaw))

//go:embed versions_macos.txt
var versionsMacOSRaw string
var versionsMacOS = strings.Split("\n", strings.TrimSpace(versionsMacOSRaw))

//go:embed versions_windows.txt
var versionsWindowsRaw string
var versionsWindows = strings.Split("\n", strings.TrimSpace(versionsWindowsRaw))

//go:embed versions_linux.txt
var versionsLinuxRaw string
var versionsLinux = strings.Split("\n", strings.TrimSpace(versionsLinuxRaw))

func randomUAChrome() string {
	f := "Mozilla/5.0 (%s) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/%s Safari/537.36"
	return fmt.Sprintf(f, randomOS(), random(versionsChrome))
}

func randomUAFirefox() string {
	f := "Mozilla/5.0 (%s; rv:%s) Gecko/20100101 Firefox/%s"
	ver := random(versionsFirefox)
	return fmt.Sprintf(f, randomOS(), ver, ver)
}

func randomUAEdge() string {
	f := "Mozilla/5.0 (%s) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/123.0.0.0 Safari/537.36 Edg/%s"
	return fmt.Sprintf(f, randomOS(), random(versionsEdge))
}

func randomUAOpera() string {
	f := "Mozilla/5.0 (%s) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/123.0.0.0 Safari/537.36 OPR/"
	return fmt.Sprintf(f, randomOS(), random(versionsOpera))
}

func randomUserAgent() string {
	switch rand.Intn(4) {
	case 0:
		return randomUAChrome()
	case 1:
		return randomUAFirefox()
	case 2:
		return randomUAEdge()
	case 3:
		return randomUAOpera()
	}
	panic("rand.Intn is broken")
}

func randomOS() string {
	switch rand.Intn(3) {
	case 0:
		return random(versionsMacOS)
	case 1:
		return random(versionsWindows)
	case 2:
		return random(versionsLinux)
	}
	panic("rand.Intn is broken")
}

func random(ss []string) string {
	return ss[rand.Intn(len(ss))]
}
