package flyscrape

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/philippta/flyscrape"

	_ "github.com/philippta/flyscrape/modules/browser"
	_ "github.com/philippta/flyscrape/modules/cache"
	_ "github.com/philippta/flyscrape/modules/cookies"
	_ "github.com/philippta/flyscrape/modules/depth"
	_ "github.com/philippta/flyscrape/modules/domainfilter"
	_ "github.com/philippta/flyscrape/modules/followlinks"
	_ "github.com/philippta/flyscrape/modules/headers"
	_ "github.com/philippta/flyscrape/modules/output/json"
	_ "github.com/philippta/flyscrape/modules/output/ndjson"
	_ "github.com/philippta/flyscrape/modules/proxy"
	_ "github.com/philippta/flyscrape/modules/ratelimit"
	_ "github.com/philippta/flyscrape/modules/retry"
	_ "github.com/philippta/flyscrape/modules/starturl"
	_ "github.com/philippta/flyscrape/modules/urlfilter"
)

type M map[string]any

type ConfigOutput struct {
	Format string `json:"format"`
	File   string `json:"file"`
}

type Config struct {
	URL            string            `json:"url"`
	URLs           []string          `json:"urls"`
	Browser        bool              `json:"browser"`
	Headless       bool              `json:"headless"`
	Depth          int               `json:"depth"`
	Follow         []string          `json:"follow"`
	AllowedDomains []string          `json:"allowedDomains"`
	BlockedDomains []string          `json:"blockedDomains"`
	AllowedURLs    []string          `json:"allowedURLs"`
	BlockedURLs    []string          `json:"blockedURLs"`
	Rate           int               `json:"rate"`
	Concurrency    int               `json:"concurrency"`
	Proxy          string            `json:"proxy"`
	Proxies        []string          `json:"proxies"`
	Cache          string            `json:"cache"`
	Cookies        string            `json:"cookies"`
	Headers        map[string]string `json:"headers"`
	Output         ConfigOutput      `json:"output"`
}

type Context struct {
	URL         string
	Doc         *Document
	AbsoluteURL func(url string) string
	Scrape      func(url string, fn ScrapeFunc) any
}

type ScrapeFunc func(ctx Context) any

func Run(cfg Config, fn ScrapeFunc) {
	jsoncfg, _ := json.Marshal(cfg)

	scraper := flyscrape.NewScraper()
	scraper.ScrapeFunc = scrapeFunc(fn)
	scraper.Script = "main.go"
	scraper.Modules = flyscrape.LoadModules(jsoncfg)
	scraper.Run()
}

func scrapeFunc(fn ScrapeFunc) flyscrape.ScrapeFunc {
	return func(p flyscrape.ScrapeParams) (any, error) {
		doc, err := goquery.NewDocumentFromReader(strings.NewReader(p.HTML))
		if err != nil {
			return nil, fmt.Errorf("new document from reader: %w", err)
		}

		baseurl, err := url.Parse(p.URL)
		if err != nil {
			return nil, fmt.Errorf("parse url: %w", err)
		}

		absoluteURL := func(ref string) string {
			abs, err := baseurl.Parse(ref)
			if err != nil {
				return ref
			}
			return abs.String()
		}

		ctx := Context{
			URL:         p.URL,
			Doc:         NewDocument(doc.Selection),
			AbsoluteURL: absoluteURL,
			Scrape: func(url string, sfn ScrapeFunc) any {
				url = absoluteURL(url)

				html, err := p.Process(url)
				if err != nil {
					return M{"error": err.Error()}
				}

				newp := flyscrape.ScrapeParams{
					HTML:    string(html),
					URL:     url,
					Process: p.Process,
				}

				data, err := scrapeFunc(sfn)(newp)
				if err != nil {
					return M{"error": err.Error()}
				}

				return data
			},
		}

		return fn(ctx), nil
	}
}

type Document struct {
	Selection *goquery.Selection
}

func NewDocument(sel *goquery.Selection) *Document {
	return &Document{Selection: sel}
}

func (d *Document) Text() string {
	return d.Selection.Text()
}

func (d *Document) Name() string {
	if d.Selection.Length() > 0 {
		return d.Selection.Get(0).Data
	}
	return ""
}

func (d *Document) Html() string {
	h, _ := goquery.OuterHtml(d.Selection)
	return h
}

func (d *Document) Attr(name string) string {
	v, _ := d.Selection.Attr(name)
	return v
}

func (d *Document) HasAttr(name string) bool {
	_, ok := d.Selection.Attr(name)
	return ok
}

func (d *Document) HasClass(className string) bool {
	return d.Selection.HasClass(className)
}

func (d *Document) Length() int {
	return d.Selection.Length()
}

func (d *Document) First() *Document {
	return NewDocument(d.Selection.First())
}

func (d *Document) Last() *Document {
	return NewDocument(d.Selection.Last())
}

func (d *Document) Get(index int) *Document {
	return NewDocument(d.Selection.Eq(index))
}

func (d *Document) Find(s string) *Document {
	return NewDocument(d.Selection.Find(s))
}

func (d *Document) Next() *Document {
	return NewDocument(d.Selection.Next())
}

func (d *Document) NextAll() *Document {
	return NewDocument(d.Selection.NextAll())
}

func (d *Document) NextUntil(s string) *Document {
	return NewDocument(d.Selection.NextUntil(s))
}

func (d *Document) Prev() *Document {
	return NewDocument(d.Selection.Prev())
}

func (d *Document) PrevAll() *Document {
	return NewDocument(d.Selection.PrevAll())
}

func (d *Document) PrevUntil(s string) *Document {
	return NewDocument(d.Selection.PrevUntil(s))
}

func (d *Document) Siblings() *Document {
	return NewDocument(d.Selection.Siblings())
}

func (d *Document) Children() *Document {
	return NewDocument(d.Selection.Children())
}

func (d *Document) Parent() *Document {
	return NewDocument(d.Selection.Parent())
}

func (d *Document) Map(callback func(*Document, int) any) []any {
	var vals []any
	d.Selection.Map(func(i int, s *goquery.Selection) string {
		vals = append(vals, callback(NewDocument(s), i))
		return ""
	})
	return vals
}

func (d *Document) Filter(callback func(*Document, int) bool) []*Document {
	var vals []*Document
	d.Selection.Each(func(i int, s *goquery.Selection) {
		el := NewDocument(s)
		ok := callback(el, i)
		if ok {
			vals = append(vals, el)
		}
	})
	return vals
}

func (d *Document) All() []*Document {
	ss := make([]*Document, 0, d.Length())
	d.Selection.Each(func(i int, s *goquery.Selection) {
		ss = append(ss, NewDocument(s))
	})
	return ss
}
