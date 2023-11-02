// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package flyscrape

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	gourl "net/url"
	"strings"
)

func NewJSLibrary(client *http.Client) Imports {
	return Imports{
		"flyscrape": map[string]any{
			"parse": jsParse(),
		},
		"flyscrape/http": map[string]any{
			"get":      jsHTTPGet(client),
			"postForm": jsHTTPPostForm(client),
			"postJSON": jsHTTPPostJSON(client),
		},
	}
}

func jsParse() func(html string) map[string]any {
	return func(html string) map[string]any {
		doc, err := DocumentFromString(html)
		if err != nil {
			return nil
		}
		return doc
	}
}

func jsHTTPGet(client *http.Client) func(url string) map[string]any {
	return func(url string) map[string]any {
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return map[string]any{"error": err.Error()}
		}
		return jsFetch(client, req)
	}
}

func jsHTTPPostForm(client *http.Client) func(url string, form map[string]any) map[string]any {
	return func(url string, form map[string]any) map[string]any {
		vals := gourl.Values{}
		for k, v := range form {
			switch v := v.(type) {
			case []any:
				for _, v := range v {
					vals.Add(k, fmt.Sprintf("%v", v))
				}
			default:
				vals.Add(k, fmt.Sprintf("%v", v))
			}
		}

		req, err := http.NewRequest("POST", url, strings.NewReader(vals.Encode()))
		if err != nil {
			return map[string]any{"error": err.Error()}
		}
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		return jsFetch(client, req)
	}
}

func jsHTTPPostJSON(client *http.Client) func(url string, data any) map[string]any {
	return func(url string, data any) map[string]any {
		b, _ := json.Marshal(data)

		req, err := http.NewRequest("POST", url, bytes.NewReader(b))
		if err != nil {
			return map[string]any{"error": err.Error()}
		}
		req.Header.Set("Content-Type", "application/json")

		return jsFetch(client, req)
	}
}

func jsFetch(client *http.Client, req *http.Request) (obj map[string]any) {
	obj = map[string]any{
		"body":    "",
		"status":  0,
		"headers": map[string]any{},
		"error":   "",
	}

	resp, err := client.Do(req)
	if err != nil {
		obj["error"] = err.Error()
		return
	}
	defer resp.Body.Close()

	obj["status"] = resp.StatusCode

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		obj["error"] = err.Error()
		return
	}

	obj["body"] = string(b)

	headers := map[string]any{}
	for name := range resp.Header {
		headers[name] = resp.Header.Get(name)
	}
	obj["headers"] = headers

	return
}
