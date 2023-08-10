package js_test

import (
	"io"
	"net/http"
	"testing"

	"flyscrape/flyscrape"
	"flyscrape/js"

	"github.com/stretchr/testify/require"
)

func TestV8(t *testing.T) {
	opts, run, err := js.Compile("../examples/esbuild.github.io.js")
	require.NoError(t, err)
	require.NotNil(t, opts)
	require.NotNil(t, run)

	html := fetch(opts.URL)
	json, err := run(flyscrape.ScrapeParams{
		HTML: html,
	})

	require.NoError(t, err)
	t.Log(json)
}

func fetch(url string) string {
	resp, err := http.Get(url)
	if err != nil {
		return ""
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return ""
	}
	return string(b)
}
