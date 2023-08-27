// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package js

import _ "embed"

//go:generate esbuild --bundle --platform=node --format=iife --outfile=flyscrape_bundle.js flyscrape.js

//go:embed flyscrape_bundle.js
var Flyscrape string

//go:embed template.js
var Template []byte
