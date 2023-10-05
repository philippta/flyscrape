// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package followlinks

import (
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/philippta/flyscrape"
)

func init() {
	flyscrape.RegisterModule(Module{})
}

type Module struct{}

func (Module) ModuleInfo() flyscrape.ModuleInfo {
	return flyscrape.ModuleInfo{
		ID:  "followlinks",
		New: func() flyscrape.Module { return new(Module) },
	}
}

func (m *Module) ReceiveResponse(resp *flyscrape.Response) {
	for _, link := range parseLinks(string(resp.Body), resp.Request.URL) {
		resp.Visit(link)
	}
}

func parseLinks(html string, origin string) []string {
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
	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		link, _ := s.Attr("href")

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

	return links
}

func isValidLink(link *url.URL) bool {
	if link.Scheme != "" && link.Scheme != "http" && link.Scheme != "https" {
		return false
	}

	return true
}

var _ flyscrape.ResponseReceiver = (*Module)(nil)
