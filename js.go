// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package flyscrape

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"strings"

	"github.com/evanw/esbuild/pkg/api"
	"github.com/philippta/flyscrape/js"
	v8 "rogchap.com/v8go"
)

type Options []byte

type TransformError struct {
	Line   int
	Column int
	Text   string
}

func (err TransformError) Error() string {
	return fmt.Sprintf("%d:%d: %s", err.Line, err.Column, err.Text)
}

func Compile(src string) (Options, ScrapeFunc, error) {
	src, err := build(src)
	if err != nil {
		return nil, nil, err
	}
	return vm(src)
}

func build(src string) (string, error) {
	res := api.Transform(src, api.TransformOptions{
		Platform: api.PlatformBrowser,
		Format:   api.FormatIIFE,
	})

	var errs []error
	for _, msg := range res.Errors {
		err := TransformError{Text: msg.Text}
		if msg.Location != nil {
			err.Line = msg.Location.Line
			err.Column = msg.Location.Column
		}
		errs = append(errs, err)
	}
	if len(res.Errors) > 0 {
		return "", errors.Join(errs...)
	}

	return string(res.Code), nil
}

func vm(src string) (Options, ScrapeFunc, error) {
	ctx := v8.NewContext()
	ctx.RunScript("var module = {}", "main.js")

	if _, err := ctx.RunScript(removeIIFE(js.Flyscrape), "main.js"); err != nil {
		return nil, nil, fmt.Errorf("running flyscrape bundle: %w", err)
	}
	if _, err := ctx.RunScript(`const require = () => require_flyscrape();`, "main.js"); err != nil {
		return nil, nil, fmt.Errorf("creating require function: %w", err)
	}
	if _, err := ctx.RunScript(removeIIFE(src), "main.js"); err != nil {
		return nil, nil, fmt.Errorf("running user script: %w", err)
	}

	cfg, err := ctx.RunScript("JSON.stringify(options)", "main.js")
	if err != nil {
		return nil, nil, fmt.Errorf("reading options: %w", err)
	}
	if !cfg.IsString() {
		return nil, nil, fmt.Errorf("options is not a string")
	}

	scrape := func(params ScrapeParams) (any, error) {
		suffix := randSeq(16)
		ctx.Global().Set("html_"+suffix, params.HTML)
		ctx.Global().Set("url_"+suffix, params.URL)
		data, err := ctx.RunScript(fmt.Sprintf(`JSON.stringify(stdin_default({html: html_%s, url: url_%s}))`, suffix, suffix), "main.js")
		if err != nil {
			return nil, err
		}

		var obj any
		if err := json.Unmarshal([]byte(data.String()), &obj); err != nil {
			return nil, err
		}

		return obj, nil
	}

	return Options(cfg.String()), scrape, nil
}

func randSeq(n int) string {
	letters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func removeIIFE(s string) string {
	s = strings.TrimPrefix(s, "(() => {\n")
	s = strings.TrimSuffix(s, "})();\n")
	return s
}
