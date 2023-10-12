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
export const config = {
    url: "https://localhost/",
}

export default function({ doc, url }) {
    return {
        headline: doc.find("h1").text(),
        body: doc.find("p").text(),
        url: url,
    }
}
`

func TestJSScrape(t *testing.T) {
	cfg, run, err := flyscrape.Compile(script)
	require.NoError(t, err)
	require.NotNil(t, cfg)
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
	cfg, run, err := flyscrape.Compile("import foo;")
	require.Error(t, err)
	require.Empty(t, cfg)
	require.Nil(t, run)

	var terr flyscrape.TransformError
	require.ErrorAs(t, err, &terr)

	require.Equal(t, terr, flyscrape.TransformError{
		Line:   1,
		Column: 10,
		Text:   `Expected "from" but found ";"`,
	})
}

func TestJSConfig(t *testing.T) {
	js := `
    export const config = {
        url: 'http://localhost/',
        depth: 5,
        allowedDomains: ['example.com'],
    }
    export default function() {}
    `
	rawCfg, _, err := flyscrape.Compile(js)
	require.NoError(t, err)

	type config struct {
		URL            string   `json:"url"`
		Depth          int      `json:"depth"`
		AllowedDomains []string `json:"allowedDomains"`
	}

	var cfg config
	err = json.Unmarshal(rawCfg, &cfg)
	require.NoError(t, err)

	require.Equal(t, config{
		URL:            "http://localhost/",
		Depth:          5,
		AllowedDomains: []string{"example.com"},
	}, cfg)
}
