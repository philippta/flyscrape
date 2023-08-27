// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package flyscrape_test

import (
	"testing"

	"flyscrape"

	"github.com/stretchr/testify/require"
)

var html = `
<html>
    <body>
        <main>
            <h1>Plugins</h1>
            <p>The plugin API allows you to inject code into various parts of the build process.</p>
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
        body: $("p").text()
    }
}
`

func TestV8(t *testing.T) {
	opts, run, err := flyscrape.Compile(script)
	require.NoError(t, err)
	require.NotNil(t, opts)
	require.NotNil(t, run)

	extract, err := run(flyscrape.ScrapeParams{
		HTML: html,
	})

	require.NoError(t, err)
	require.Equal(t, "Plugins", extract.(map[string]any)["headline"])
	require.Equal(t, "The plugin API allows you to inject code into various parts of the build process.", extract.(map[string]any)["body"])
}
