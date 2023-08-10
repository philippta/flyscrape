package js

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	"flyscrape/flyscrape"
	"flyscrape/js/jsbundle"

	"github.com/evanw/esbuild/pkg/api"
	v8 "rogchap.com/v8go"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func Compile(file string) (*flyscrape.ScrapeOptions, flyscrape.ScrapeFunc, error) {
	src, err := build(file)
	if err != nil {
		return nil, nil, err
	}
	return vm(src)
}

func build(file string) (string, error) {
	data, err := os.ReadFile(file)
	if err != nil {
		return "", fmt.Errorf("read %q: %w", file, err)
	}

	res := api.Transform(string(data), api.TransformOptions{
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

func vm(src string) (*flyscrape.ScrapeOptions, flyscrape.ScrapeFunc, error) {
	ctx := v8.NewContext()
	ctx.RunScript("var module = {}", "main.js")

	if _, err := ctx.RunScript(removeIIFE(jsbundle.Flyscrape), "main.js"); err != nil {
		return nil, nil, fmt.Errorf("run bundled js: %w", err)
	}
	if _, err := ctx.RunScript(`const require = () => require_flyscrape();`, "main.js"); err != nil {
		return nil, nil, fmt.Errorf("define require: %w", err)
	}
	if _, err := ctx.RunScript(removeIIFE(src), "main.js"); err != nil {
		return nil, nil, fmt.Errorf("userscript: %w", err)
	}

	var opts flyscrape.ScrapeOptions

	url, err := ctx.RunScript("options.url", "main.js")
	if err != nil {
		return nil, nil, fmt.Errorf("eval: options.url: %w", err)
	}
	opts.URL = url.String()

	depth, err := ctx.RunScript("options.depth", "main.js")
	if err != nil {
		return nil, nil, fmt.Errorf("eval: options.depth: %w", err)
	}
	opts.Depth = int(depth.Integer())

	scrape := func(params flyscrape.ScrapeParams) (flyscrape.M, error) {
		suffix := randSeq(10)
		ctx.Global().Set("html_"+suffix, params.HTML)
		data, err := ctx.RunScript(fmt.Sprintf(`JSON.stringify(stdin_default({html: html_%s}))`, suffix), "main.js")
		if err != nil {
			return nil, err
		}

		var obj flyscrape.M
		if err := json.Unmarshal([]byte(data.String()), &obj); err != nil {
			return nil, err
		}

		return obj, nil
	}

	return &opts, scrape, nil
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
