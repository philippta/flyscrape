// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package flyscrape_test

import (
	"encoding/json"
	"net/http"
	"os"
	"testing"

	"github.com/philippta/flyscrape"
	"github.com/stretchr/testify/require"
)

func TestJSLibParse(t *testing.T) {
	script := `
    import { parse } from "flyscrape"

    const doc = parse('<div class=foo>Hello world</div>')
    export const text = doc.find(".foo").text()

    export default function () {}
    `

	client := &http.Client{
		Transport: flyscrape.MockTransport(200, html),
	}

	options := flyscrape.BuildOptions("")

	imports, _ := flyscrape.NewJSLibrary(client)
	exports, err := flyscrape.Compile(script, imports, options)
	require.NoError(t, err)

	h, ok := exports["text"].(string)
	require.True(t, ok)
	require.Equal(t, "Hello world", h)
}

func TestJSLibHTTPGet(t *testing.T) {
	script := `
    import http from "flyscrape/http"

    const res = http.get("https://example.com")

    export const body = res.body;
    export const status = res.status;
    export const error = res.error;
    export const headers = res.headers;

    export default function () {}
    `

	client := &http.Client{
		Transport: flyscrape.MockTransport(200, html),
	}

	options := flyscrape.BuildOptions("")

	imports, _ := flyscrape.NewJSLibrary(client)
	exports, err := flyscrape.Compile(script, imports, options)
	require.NoError(t, err)

	body, ok := exports["body"].(string)
	require.True(t, ok)
	require.Equal(t, html, body)

	status, ok := exports["status"].(int64)
	require.True(t, ok)
	require.Equal(t, int64(200), status)

	error, ok := exports["error"].(string)
	require.True(t, ok)
	require.Equal(t, "", error)

	headers, ok := exports["headers"].(map[string]any)
	require.True(t, ok)
	require.NotEmpty(t, headers)
}

func TestJSLibHTTPPostForm(t *testing.T) {
	script := `
    import http from "flyscrape/http"

    const res = http.postForm("https://example.com", {
        username: "foo",
        password: "bar",
        arr: [1,2,3],
    })

    export const body = res.body;
    export const status = res.status;
    export const error = res.error;
    export const headers = res.headers;

    export default function () {}
    `

	client := &http.Client{
		Transport: flyscrape.RoundTripFunc(func(r *http.Request) (*http.Response, error) {
			require.Equal(t, "POST", r.Method)
			require.Equal(t, "application/x-www-form-urlencoded", r.Header.Get("Content-Type"))
			require.Equal(t, "foo", r.FormValue("username"))
			require.Equal(t, "bar", r.FormValue("password"))
			require.Len(t, r.Form["arr"], 3)

			return flyscrape.MockResponse(400, "Bad Request")
		}),
	}

	options := flyscrape.BuildOptions("")

	imports, _ := flyscrape.NewJSLibrary(client)
	exports, err := flyscrape.Compile(script, imports, options)
	require.NoError(t, err)

	body, ok := exports["body"].(string)
	require.True(t, ok)
	require.Equal(t, "Bad Request", body)

	status, ok := exports["status"].(int64)
	require.True(t, ok)
	require.Equal(t, int64(400), status)

	error, ok := exports["error"].(string)
	require.True(t, ok)
	require.Equal(t, "", error)

	headers, ok := exports["headers"].(map[string]any)
	require.True(t, ok)
	require.NotEmpty(t, headers)
}

func TestJSLibHTTPPostJSON(t *testing.T) {
	script := `
    import http from "flyscrape/http"

    const res = http.postJSON("https://example.com", {
        username: "foo",
        password: "bar",
    })

    export const body = res.body;
    export const status = res.status;
    export const error = res.error;
    export const headers = res.headers;

    export default function () {}
    `

	client := &http.Client{
		Transport: flyscrape.RoundTripFunc(func(r *http.Request) (*http.Response, error) {
			require.Equal(t, "POST", r.Method)
			require.Equal(t, "application/json", r.Header.Get("Content-Type"))

			m := map[string]any{}
			json.NewDecoder(r.Body).Decode(&m)
			require.Equal(t, "foo", m["username"])
			require.Equal(t, "bar", m["password"])

			return flyscrape.MockResponse(400, "Bad Request")
		}),
	}

	options := flyscrape.BuildOptions("")

	imports, _ := flyscrape.NewJSLibrary(client)
	exports, err := flyscrape.Compile(script, imports, options)
	require.NoError(t, err)

	body, ok := exports["body"].(string)
	require.True(t, ok)
	require.Equal(t, "Bad Request", body)

	status, ok := exports["status"].(int64)
	require.True(t, ok)
	require.Equal(t, int64(400), status)

	error, ok := exports["error"].(string)
	require.True(t, ok)
	require.Equal(t, "", error)

	headers, ok := exports["headers"].(map[string]any)
	require.True(t, ok)
	require.NotEmpty(t, headers)
}

func TestJSLibHTTPDownload(t *testing.T) {
	cwd, err := os.Getwd()
	require.NoError(t, err)

	tmpdir, err := os.MkdirTemp("", "http-download")
	require.NoError(t, err)

	defer os.RemoveAll(tmpdir)
	defer os.Chdir(cwd)
	os.Chdir(tmpdir)

	script := `
    import http from "flyscrape/http";

    http.download("https://example.com/foo.txt", "foo.txt");
    http.download("https://example.com/foo.txt", "dir/my-foo.txt");
    http.download("https://example.com/bar.txt", "dir/");
    http.download("https://example.com/baz.txt", "dir");
    http.download("https://example.com/content-disposition", ".");
    http.download("https://example.com/hack.txt", ".");
    http.download("https://example.com/no-dest.txt");
    http.download("https://example.com/404.txt");
    `

	nreqs := 0
	client := &http.Client{
		Transport: flyscrape.RoundTripFunc(func(r *http.Request) (*http.Response, error) {
			nreqs++

			if r.URL.Path == "/content-disposition" {
				resp, err := flyscrape.MockResponse(200, "hello world")
				resp.Header.Set("Content-Disposition", `attachment; filename="qux.txt"`)
				return resp, err
			}
			if r.URL.Path == "/hack.txt" {
				resp, err := flyscrape.MockResponse(200, "hello world")
				resp.Header.Set("Content-Disposition", `attachment; filename="../../hack.txt"`)
				return resp, err
			}
			if r.URL.Path == "/404.txt" {
				resp, err := flyscrape.MockResponse(404, "hello world")
				return resp, err
			}

			return flyscrape.MockResponse(200, "hello world")
		}),
	}

	options := flyscrape.BuildOptions("")

	imports, wait := flyscrape.NewJSLibrary(client)
	_, err = flyscrape.Compile(script, imports, options)
	require.NoError(t, err)

	wait()

	require.Equal(t, nreqs, 8)
	require.FileExists(t, "foo.txt")
	require.FileExists(t, "dir/my-foo.txt")
	require.FileExists(t, "dir/bar.txt")
	require.FileExists(t, "dir/baz.txt")
	require.FileExists(t, "qux.txt")
	require.FileExists(t, "hack.txt")
	require.FileExists(t, "no-dest.txt")
	require.NoFileExists(t, "404.txt")
}
