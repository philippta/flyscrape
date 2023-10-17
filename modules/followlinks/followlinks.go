// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package followlinks

import (
	"net/url"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/philippta/flyscrape"
)

func init() {
	flyscrape.RegisterModule(Module{})
}

type Module struct {
	Follow []string `json:"follow"`
}

func (Module) ModuleInfo() flyscrape.ModuleInfo {
	return flyscrape.ModuleInfo{
		ID:  "followlinks",
		New: func() flyscrape.Module { return new(Module) },
	}
}

func (m *Module) Provision(ctx flyscrape.Context) {
	if len(m.Follow) == 0 {
		m.Follow = []string{"a[href]"}
	}
}

func (m *Module) ReceiveResponse(resp *flyscrape.Response) {
	for _, link := range m.parseLinks(string(resp.Body), resp.Request.URL) {
		resp.Visit(link)
	}
}

func (m *Module) parseLinks(html string, origin string) []string {
	var links []string
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return nil
	}

	originurl, err := url.Parse(origin)
	if err != nil {
		return nil
	}

	uniqueLinks := make(map[string]bool)

	for _, selector := range m.Follow {
		attr := parseSelectorAttr(selector)
		doc.Find(selector).Each(func(i int, s *goquery.Selection) {
			link, _ := s.Attr(attr)

			parsedLink, err := originurl.Parse(link)

			if err != nil || !isValidLink(parsedLink) {
				return
			}

			absLink := parsedLink.String()

			if !uniqueLinks[absLink] {
				links = append(links, absLink)
				uniqueLinks[absLink] = true
			}
		})
	}

	return links
}

func isValidLink(link *url.URL) bool {
	if link.Scheme != "" && link.Scheme != "http" && link.Scheme != "https" {
		return false
	}

	return true
}

func parseSelectorAttr(sel string) string {
	matches := selectorExpr.FindAllString(sel, -1)
	if len(matches) == 0 {
		return "href"
	}

	attr := attrExpr.FindString(matches[len(matches)-1])
	if attr == "" {
		return "href"
	}

	return attr
}

var (
	_ flyscrape.Provisioner      = (*Module)(nil)
	_ flyscrape.ResponseReceiver = (*Module)(nil)
)

var (
	selectorExpr = regexp.MustCompile(`\[(.*?)\]`)
	attrExpr     = regexp.MustCompile(`[\w-]+`)
)
