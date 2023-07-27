package jsbundle

import _ "embed"

//go:generate esbuild ../../nodestuff/src/flyscrape.js --bundle --platform=node --outfile=flyscrape.js

//go:embed flyscrape.js
var Flyscrape []byte
