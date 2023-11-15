// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package flyscrape

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime"
	"net/http"
	gourl "net/url"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/sync/errgroup"
)

func NewJSLibrary(client *http.Client) (imports Imports, wait func()) {
	downloads := &errgroup.Group{}

	// Allow 5 parallel downloads. Why 5?
	// Docker downloads 3 layers in parallel.
	// My Chrome downloads up to 6 files in parallel.
	// 5 feels like a reasonable number.
	downloads.SetLimit(5)

	im := Imports{
		"flyscrape": map[string]any{
			"parse": jsParse(),
		},
		"flyscrape/http": map[string]any{
			"get":      jsHTTPGet(client),
			"postForm": jsHTTPPostForm(client),
			"postJSON": jsHTTPPostJSON(client),
			"download": jsHTTPDownload(client, downloads),
		},
	}

	return im, func() { downloads.Wait() }
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

func jsHTTPDownload(client *http.Client, g *errgroup.Group) func(url string, dst string) {
	fileExists := func(name string) bool {
		_, err := os.Stat(name)
		return err == nil
	}

	isDir := func(path string) bool {
		if strings.HasSuffix(path, "/") {
			return true
		}
		if filepath.Ext(path) == "" {
			return true
		}
		s, err := os.Stat(path)
		return err == nil && s.IsDir()
	}

	suggestedFilename := func(url, contentDisp string) string {
		filename := filepath.Base(url)

		if contentDisp == "" {
			return filename
		}

		_, params, err := mime.ParseMediaType(contentDisp)
		if err != nil {
			return filename
		}

		name, ok := params["filename"]
		if !ok || name == "" {
			return filename
		}

		return filepath.Base(name)
	}

	return func(url string, dst string) {
		g.Go(func() error {
			req, err := http.NewRequest("GET", url, nil)
			if err != nil {
				log.Printf("error downloading file %q: %v", url, err)
				return nil
			}
			req.Header.Add(HeaderBypassCache, "true")

			resp, err := client.Do(req)
			if err != nil {
				log.Printf("error downloading file %q: %v", url, err)
				return nil
			}
			defer resp.Body.Close()

			if resp.StatusCode < 200 || resp.StatusCode >= 300 {
				log.Printf("error downloading file %q: unexpected status code %d", url, resp.StatusCode)
				return nil
			}

			dst, err = filepath.Abs(dst)
			if err != nil {
				log.Printf("error downloading file %q: abs path failed: %v", url, err)
				return nil
			}

			if isDir(dst) {
				name := suggestedFilename(url, resp.Header.Get("Content-Disposition"))
				dst = filepath.Join(dst, name)
			}

			if fileExists(dst) {
				return nil
			}

			os.MkdirAll(filepath.Dir(dst), 0o755)
			f, err := os.Create(dst)
			if err != nil {
				log.Printf("error downloading file %q: file save failed: %v", url, err)
				return nil
			}
			defer f.Close()

			io.Copy(f, resp.Body)
			return nil
		})
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
