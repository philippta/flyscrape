// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package flyscrape_test

import (
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
	"testing"

	"github.com/philippta/flyscrape"
	"github.com/stretchr/testify/require"
)

func TestJSLibFetchDocument(t *testing.T) {
	script := `
    import { fetchDocument } from "flyscrape"

    const doc = fetchDocument("https://example.com")
    export const headline = doc.find("h1").text()
    `

	client := &http.Client{
		Transport: flyscrape.MockTransport(200, html),
	}

	exports, err := flyscrape.Compile(script, flyscrape.NewJSLibrary(client))
	require.NoError(t, err)

	h, ok := exports["headline"].(string)
	require.True(t, ok)
	require.Equal(t, "headline", h)
}

func TestJSLibSubmitForm(t *testing.T) {
	script := `
    import { submitForm } from "flyscrape"

    const doc = submitForm("https://example.com", {
        "username": "foo",
        "password": "bar",
    })

    export const text = doc.find("div").text()
    `

	var username, password string

	jar, _ := cookiejar.New(nil)
	client := &http.Client{
		Jar: jar,
		Transport: flyscrape.RoundTripFunc(func(r *http.Request) (*http.Response, error) {
			username = r.FormValue("username")
			password = r.FormValue("password")

			resp := &http.Response{
				StatusCode: 200,
				Status:     fmt.Sprintf("%d %s", 200, http.StatusText(200)),
				Body:       io.NopCloser(strings.NewReader(`<div>Login successful</div>`)),
				Header:     http.Header{},
			}

			cookie := http.Cookie{
				Name:   "example",
				Value:  "Hello world!",
				Path:   "/",
				MaxAge: 3600,
			}

			resp.Header.Add("Set-Cookie", cookie.String())
			return resp, nil
		}),
	}

	exports, err := flyscrape.Compile(script, flyscrape.NewJSLibrary(client))
	require.NoError(t, err)

	text, ok := exports["text"].(string)
	require.True(t, ok)
	require.Equal(t, "Login successful", text)
	require.Equal(t, "foo", username)
	require.Equal(t, "bar", password)

	u, _ := url.Parse("https://example.com")
	require.NotEmpty(t, jar.Cookies(u))
}
