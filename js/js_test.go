package js_test

import (
	"io"
	"net/http"
	"testing"

	"flyscrape/js"

	"github.com/stretchr/testify/require"
)

func TestV8(t *testing.T) {
	opts, run, err := js.Compile("../examples/esbuild.github.io.js")
	require.NoError(t, err)

	html := fetch(opts.URL)
	json := run(js.RunOptions{
		HTML: html,
	})

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
