// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package flyscrape_test

import (
	"encoding/json"
	"testing"

	"github.com/dop251/goja"
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
	exports, err := flyscrape.Compile(script, nil)
	require.NoError(t, err)
	require.NotNil(t, exports)
	require.NotEmpty(t, exports.Config)

	result, err := exports.Scrape(flyscrape.ScrapeParams{
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

func TestJSScrapeObject(t *testing.T) {
	js := `
    export default function() {
        return {foo: "bar"}
    }
    `
	exports, err := flyscrape.Compile(js, nil)
	require.NoError(t, err)

	result, err := exports.Scrape(flyscrape.ScrapeParams{
		HTML: html,
		URL:  "http://localhost/",
	})
	require.NoError(t, err)

	m, ok := result.(map[string]any)
	require.True(t, ok)
	require.Equal(t, "bar", m["foo"])
}

func TestJSScrapeNull(t *testing.T) {
	js := `
    export default function() {
        return null
    }
    `
	exports, err := flyscrape.Compile(js, nil)
	require.NoError(t, err)

	result, err := exports.Scrape(flyscrape.ScrapeParams{
		HTML: html,
		URL:  "http://localhost/",
	})
	require.NoError(t, err)
	require.Nil(t, result)
}

func TestJSScrapeString(t *testing.T) {
	js := `
    export default function() {
        return "foo"
    }
    `
	exports, err := flyscrape.Compile(js, nil)
	require.NoError(t, err)

	result, err := exports.Scrape(flyscrape.ScrapeParams{
		HTML: html,
		URL:  "http://localhost/",
	})
	require.NoError(t, err)

	m, ok := result.(string)
	require.True(t, ok)
	require.Equal(t, "foo", m)
}

func TestJSScrapeArray(t *testing.T) {
	js := `
    export default function() {
        return [1,2,3]
    }
    `
	exports, err := flyscrape.Compile(js, nil)
	require.NoError(t, err)

	result, err := exports.Scrape(flyscrape.ScrapeParams{
		HTML: html,
		URL:  "http://localhost/",
	})
	require.NoError(t, err)

	m, ok := result.([]any)
	require.True(t, ok)
	require.Equal(t, int64(1), m[0])
	require.Equal(t, int64(2), m[1])
	require.Equal(t, int64(3), m[2])
}

func TestJSCompileError(t *testing.T) {
	exports, err := flyscrape.Compile("import foo;", nil)
	require.Error(t, err)
	require.Nil(t, exports)

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
	exports, err := flyscrape.Compile(js, nil)
	require.NoError(t, err)
	require.NotNil(t, exports)
	require.NotEmpty(t, exports.Config())

	type config struct {
		URL            string   `json:"url"`
		Depth          int      `json:"depth"`
		AllowedDomains []string `json:"allowedDomains"`
	}

	var cfg config
	err = json.Unmarshal(exports.Config(), &cfg)
	require.NoError(t, err)

	require.Equal(t, config{
		URL:            "http://localhost/",
		Depth:          5,
		AllowedDomains: []string{"example.com"},
	}, cfg)
}

func TestJSImports(t *testing.T) {
	js := `
    import A from "pkg-a"
    import { bar } from "pkg-a/pkg-b"

    export const config = {}
    export default function() {}

	export const a = A.foo
	export const b = bar()
    `
	imports := flyscrape.Imports{
		"pkg-a": map[string]any{
			"foo": 10,
		},
		"pkg-a/pkg-b": map[string]any{
			"bar": func() string {
				return "baz"
			},
		},
	}

	exports, err := flyscrape.Compile(js, imports)
	require.NoError(t, err)
	require.NotNil(t, exports)

	require.Equal(t, int64(10), exports["a"].(int64))
	require.Equal(t, "baz", exports["b"].(string))
}

func TestJSArbitraryFunction(t *testing.T) {
	js := `
    export const config = {}
    export default function() {}
    export function foo() {
        return "bar";
    }
    `
	exports, err := flyscrape.Compile(js, nil)
	require.NoError(t, err)
	require.NotNil(t, exports)

	foo := func() string {
		fn := exports["foo"].(func(goja.FunctionCall) goja.Value)
		return fn(goja.FunctionCall{}).String()
	}

	require.Equal(t, "bar", foo())
}

func TestJSArbitraryConstString(t *testing.T) {
	js := `
    export const config = {}
    export default function() {}
    export const foo = "bar"
    `
	exports, err := flyscrape.Compile(js, nil)
	require.NoError(t, err)
	require.NotNil(t, exports)

	require.Equal(t, "bar", exports["foo"].(string))
}
