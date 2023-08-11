package flyscrape_test

import (
	"os"
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

func TestV8(t *testing.T) {
	data, err := os.ReadFile("examples/esbuild.github.io.js")
	require.NoError(t, err)

	opts, run, err := flyscrape.Compile(string(data))
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
