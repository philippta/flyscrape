// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package flyscrape

import (
	"fmt"
	"io"
	"net/http"
	gourl "net/url"
	"strings"
)

func NewJSLibrary(client *http.Client) Imports {
	return Imports{
		"flyscrape": map[string]any{
			"fetchDocument": jsFetchDocument(client),
			"submitForm":    jsSubmitForm(client),
		},
	}
}

func jsFetchDocument(client *http.Client) func(url string) map[string]any {
	return func(url string) map[string]any {
		resp, err := client.Get(url)
		if err != nil {
			return nil
		}
		defer resp.Body.Close()

		var b strings.Builder
		if _, err := io.Copy(&b, resp.Body); err != nil {
			return nil
		}

		doc, err := DocumentFromString(b.String())
		if err != nil {
			return nil
		}

		return doc
	}
}

func jsSubmitForm(client *http.Client) func(url string, data map[string]any) map[string]any {
	return func(url string, data map[string]any) map[string]any {
		form := gourl.Values{}
		for k, v := range data {
			form.Set(k, fmt.Sprintf("%v", v))
		}

		resp, err := client.PostForm(url, form)
		if err != nil {
			return nil
		}
		defer resp.Body.Close()

		var b strings.Builder
		if _, err := io.Copy(&b, resp.Body); err != nil {
			return nil
		}

		doc, err := DocumentFromString(b.String())
		if err != nil {
			return nil
		}

		return doc
	}
}
