package js

import _ "embed"

//go:generate esbuild --bundle --platform=node --format=iife --outfile=flyscrape_bundle.js flyscrape.js

//go:embed flyscrape_bundle.js
var Flyscrape string

//go:embed template.js
var Template []byte
