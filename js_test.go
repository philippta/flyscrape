// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package flyscrape_test

import (
	"encoding/json"
	"testing"

	"github.com/philippta/flyscrape"
	"github.com/stretchr/testify/require"
)

var html = `
<html>
    <body>
        <main>
            <h1>headline</h1>
            <p>paragraph</p>
        </main>
    </body>
</html>`

var script = `
import { parse } from "flyscrape";

export const options = {
    url: "https://localhost/",
}

export default function({ html, url }) {
    const $ = parse(html);

    return {
        headline: $("h1").text(),
        body: $("p").text(),
        url: url,
    }
}
`

func TestJSScrape(t *testing.T) {
	opts, run, err := flyscrape.Compile(script)
	require.NoError(t, err)
	require.NotNil(t, opts)
	require.NotNil(t, run)

	result, err := run(flyscrape.ScrapeParams{
		HTML: html,
		URL:  "http://localhost/",
	})

	require.NoError(t, err)

	m, ok := result.(map[string]any)
	require.True(t, ok)
	require.Equal(t, "headline", m["headline"])
	require.Equal(t, "paragraph", m["body"])
	require.Equal(t, "http://localhost/", m["url"])
}

func TestJSCompileError(t *testing.T) {
	opts, run, err := flyscrape.Compile("import foo;")
	require.Error(t, err)
	require.Empty(t, opts)
	require.Nil(t, run)

	var terr flyscrape.TransformError
	require.ErrorAs(t, err, &terr)

	require.Equal(t, terr, flyscrape.TransformError{
		Line:   1,
		Column: 10,
		Text:   `Expected "from" but found ";"`,
	})
}

func TestJSOptions(t *testing.T) {
	js := `
    export const options = {
        url: 'http://localhost/',
        depth: 5,
        allowedDomains: ['example.com'],
    }
    export default function() {}
    `
	rawOpts, _, err := flyscrape.Compile(js)
	require.NoError(t, err)

	type options struct {
		URL            string   `json:"url"`
		Depth          int      `json:"depth"`
		AllowedDomains []string `json:"allowedDomains"`
	}

	var opts options
	err = json.Unmarshal(rawOpts, &opts)
	require.NoError(t, err)

	require.Equal(t, options{
		URL:            "http://localhost/",
		Depth:          5,
		AllowedDomains: []string{"example.com"},
	}, opts)
}
