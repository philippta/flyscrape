package flyscrape

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"flyscrape/js"

	"github.com/evanw/esbuild/pkg/api"
	v8 "rogchap.com/v8go"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func Compile(src string) (ScrapeOptions, ScrapeFunc, error) {
	src, err := build(src)
	if err != nil {
		return ScrapeOptions{}, nil, err
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
		errs = append(errs, fmt.Errorf("%s", msg.Text))
	}
	if len(res.Errors) > 0 {
		return "", errors.Join(errs...)
	}

	return string(res.Code), nil
}

func vm(src string) (ScrapeOptions, ScrapeFunc, error) {
	ctx := v8.NewContext()
	ctx.RunScript("var module = {}", "main.js")

	if _, err := ctx.RunScript(removeIIFE(js.Flyscrape), "main.js"); err != nil {
		return ScrapeOptions{}, nil, fmt.Errorf("running flyscrape bundle: %w", err)
	}
	if _, err := ctx.RunScript(`const require = () => require_flyscrape();`, "main.js"); err != nil {
		return ScrapeOptions{}, nil, fmt.Errorf("creating require function: %w", err)
	}
	if _, err := ctx.RunScript(removeIIFE(src), "main.js"); err != nil {
		return ScrapeOptions{}, nil, fmt.Errorf("running user script: %w", err)
	}

	var opts ScrapeOptions
	optsJSON, err := ctx.RunScript("JSON.stringify(options)", "main.js")
	if err != nil {
		return ScrapeOptions{}, nil, fmt.Errorf("reading options: %w", err)
	}
	if err := json.Unmarshal([]byte(optsJSON.String()), &opts); err != nil {
		return ScrapeOptions{}, nil, fmt.Errorf("decoding options json: %w", err)
	}

	scrape := func(params ScrapeParams) (any, error) {
		suffix := randSeq(10)
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

	return opts, scrape, nil
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
