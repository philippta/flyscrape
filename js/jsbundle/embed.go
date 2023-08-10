package jsbundle

import _ "embed"

//go:generate esbuild --bundle --platform=node --format=iife --outfile=flyscrape.js ../../nodestuff/src/flyscrape.js

//go:embed flyscrape.js
var Flyscrape string
