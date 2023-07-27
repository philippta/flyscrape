package js

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"time"

	"flyscrape/js/jsbundle"

	"github.com/evanw/esbuild/pkg/api"
	v8 "rogchap.com/v8go"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

type RunOptions struct {
	HTML string
}

type RunFunc func(RunOptions) any

type Options struct {
	URL string `json:"url"`
}

func Compile(file string) (*Options, RunFunc, error) {
	src, err := build(file)
	if err != nil {
		return nil, nil, err
	}
	os.WriteFile("out.js", []byte(src), 0o644)
	return vm(src)
}

func build(file string) (string, error) {
	dir, err := os.MkdirTemp("", "flyscrape")
	if err != nil {
		return "", err
	}
	defer os.RemoveAll(dir)

	tmpfile := filepath.Join(dir, "flyscrape.js")
	if err := os.WriteFile(tmpfile, jsbundle.Flyscrape, 0o644); err != nil {
		return "", err
	}

	resolve := api.Plugin{
		Name: "flyscrape",
		Setup: func(build api.PluginBuild) {
			build.OnResolve(api.OnResolveOptions{
				Filter: "^flyscrape$",
			}, func(ora api.OnResolveArgs) (api.OnResolveResult, error) {
				return api.OnResolveResult{Path: tmpfile}, nil
			})
		},
	}

	res := api.Build(api.BuildOptions{
		EntryPoints: []string{file},
		Bundle:      true,
		Platform:    api.PlatformNode,
		Plugins:     []api.Plugin{resolve},
	})

	var errs []error
	for _, msg := range res.Errors {
		errs = append(errs, fmt.Errorf("%s", msg.Text))
	}
	if len(res.Errors) > 0 {
		return "", errors.Join(errs...)
	}

	out := string(res.OutputFiles[0].Contents)
	return out, nil
}

func vm(src string) (*Options, RunFunc, error) {
	os.WriteFile("out.js", []byte(src), 0o644)

	ctx := v8.NewContext()
	ctx.RunScript("var module = {}", "main.js")
	if _, err := ctx.RunScript(src, "main.js"); err != nil {
		return nil, nil, fmt.Errorf("run bundled js: %w", err)
	}

	val, err := ctx.RunScript("module.exports.options", "main.js")
	if err != nil {
		return nil, nil, fmt.Errorf("export options: %w", err)
	}
	options, err := val.AsObject()
	if err != nil {
		return nil, nil, fmt.Errorf("cast options as object: %w", err)
	}

	var opts Options
	url, err := options.Get("url")
	if err != nil {
		return nil, nil, fmt.Errorf("getting url from options: %w", err)
	}
	opts.URL = url.String()

	run := func(ro RunOptions) any {
		suffix := randSeq(10)
		ctx.Global().Set("html_"+suffix, ro.HTML)
		data, err := ctx.RunScript(fmt.Sprintf(`JSON.stringify(module.exports.default({html: html_%s}))`, suffix), "main.js")
		if err != nil {
			return err.Error()
		}

		var obj any
		if err := json.Unmarshal([]byte(data.String()), &obj); err != nil {
			return err.Error()
		}

		return obj
	}
	return &opts, run, nil
}

func randSeq(n int) string {
	letters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
