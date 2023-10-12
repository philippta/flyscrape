// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package flyscrape

import (
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/PuerkitoBio/goquery"
	"github.com/dop251/goja"
	"github.com/dop251/goja_nodejs/console"
	"github.com/dop251/goja_nodejs/require"
	"github.com/evanw/esbuild/pkg/api"
)

//go:embed template.js
var ScriptTemplate []byte

type Config []byte

type ScrapeParams struct {
	HTML string
	URL  string
}

type ScrapeFunc func(ScrapeParams) (any, error)

type TransformError struct {
	Line   int
	Column int
	Text   string
}

func (err TransformError) Error() string {
	return fmt.Sprintf("%d:%d: %s", err.Line, err.Column, err.Text)
}

func Compile(src string) (Config, ScrapeFunc, error) {
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

func vm(src string) (Config, ScrapeFunc, error) {
	vm := goja.New()

	registry := &require.Registry{}
	registry.Enable(vm)

	console.Enable(vm)

	if _, err := vm.RunString(removeIIFE(src)); err != nil {
		return nil, nil, fmt.Errorf("running user script: %w", err)
	}

	cfg, err := vm.RunString("JSON.stringify(config)")
	if err != nil {
		return nil, nil, fmt.Errorf("reading config: %w", err)
	}

	var c atomic.Uint64
	var lock sync.Mutex

	scrape := func(p ScrapeParams) (any, error) {
		lock.Lock()
		defer lock.Unlock()

		doc, err := goquery.NewDocumentFromReader(strings.NewReader(p.HTML))
		if err != nil {
			log.Println(err)
			return nil, err
		}

		baseurl, err := url.Parse(p.URL)
		if err != nil {
			log.Println(err)
			return nil, err
		}

		suffix := strconv.FormatUint(c.Add(1), 10)
		vm.Set("url_"+suffix, p.URL)
		vm.Set("doc_"+suffix, wrap(vm, doc.Selection))
		vm.Set("absurl_"+suffix, func(ref string) string {
			abs, err := baseurl.Parse(ref)
			if err != nil {
				log.Println(err)
				return ref
			}
			return abs.String()
		})

		data, err := vm.RunString(fmt.Sprintf(`JSON.stringify(stdin_default({doc: doc_%s, url: url_%s, absoluteURL: absurl_%s}))`, suffix, suffix, suffix))
		if err != nil {
			log.Println(err)
			return nil, err
		}

		var obj any
		if err := json.Unmarshal([]byte(data.String()), &obj); err != nil {
			log.Println(err)
			return nil, err
		}

		return obj, nil
	}

	return Config(cfg.String()), scrape, nil
}

func wrap(vm *goja.Runtime, sel *goquery.Selection) map[string]any {
	o := map[string]any{}
	o["WARNING"] = "Forgot to call text(), html() or attr()?"
	o["text"] = sel.Text
	o["html"] = func() string { h, _ := goquery.OuterHtml(sel); return h }
	o["attr"] = func(name string) string { v, _ := sel.Attr(name); return v }
	o["hasAttr"] = func(name string) bool { _, ok := sel.Attr(name); return ok }
	o["hasClass"] = sel.HasClass
	o["length"] = sel.Length()
	o["first"] = func() map[string]any { return wrap(vm, sel.First()) }
	o["last"] = func() map[string]any { return wrap(vm, sel.Last()) }
	o["get"] = func(index int) map[string]any { return wrap(vm, sel.Eq(index)) }
	o["find"] = func(s string) map[string]any { return wrap(vm, sel.Find(s)) }
	o["next"] = func() map[string]any { return wrap(vm, sel.Next()) }
	o["prev"] = func() map[string]any { return wrap(vm, sel.Prev()) }
	o["siblings"] = func() map[string]any { return wrap(vm, sel.Siblings()) }
	o["children"] = func() map[string]any { return wrap(vm, sel.Children()) }
	o["parent"] = func() map[string]any { return wrap(vm, sel.Parent()) }
	o["map"] = func(callback func(map[string]any, int) any) []any {
		var vals []any
		sel.Map(func(i int, s *goquery.Selection) string {
			vals = append(vals, callback(wrap(vm, s), i))
			return ""
		})
		return vals
	}
	o["filter"] = func(callback func(map[string]any, int) bool) []any {
		var vals []any
		sel.Each(func(i int, s *goquery.Selection) {
			el := wrap(vm, s)
			ok := callback(el, i)
			if ok {
				vals = append(vals, el)
			}
		})
		return vals
	}
	return o
}

func removeIIFE(s string) string {
	s = strings.TrimPrefix(s, "(() => {\n")
	s = strings.TrimSuffix(s, "})();\n")
	return s
}
