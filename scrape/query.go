package scrape

import (
	"strings"

	"github.com/PuerkitoBio/goquery"
)

var emptyDoc, _ = goquery.NewDocumentFromReader(strings.NewReader(""))

func Doc(html string) *goquery.Selection {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return emptyDoc.Selection
	}
	return doc.Selection
}

func Query(s *goquery.Selection, selector string) string {
	val := s.Find(selector).First().Text()
	return strings.TrimSpace(val)
}

func QueryAttr(s *goquery.Selection, selector, attr string) string {
	val := s.Find(selector).First().AttrOr(attr, "")
	return strings.TrimSpace(val)
}

func QueryHTML(s *goquery.Selection, selector string) string {
	val, err := goquery.OuterHtml(s.Find(selector))
	if err != nil {
		return ""
	}
	return strings.TrimSpace(val)
}

func QueryFunc(s *goquery.Selection, selector string, f func(*goquery.Selection)) {
	s.Find(selector).Each(func(i int, s *goquery.Selection) {
		f(s)
	})
}
