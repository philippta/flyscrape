package flyscrape_test

import (
	"testing"

	"github.com/philippta/flyscrape/flyscrape"
)

func TestRun(t *testing.T) {
	cfg := flyscrape.Config{
		URL:      "https://news.ycombinator.com/",
		Browser:  true,
		Headless: false,
		Output: flyscrape.ConfigOutput{
			Format: "ndjson",
		},
	}

	fn := func(ctx flyscrape.Context) any {
		post := ctx.Doc.Find(".athing.submission").First()
		title := post.Find(".titleline > a").Text()
		commentsLink := post.Next().Find("a").Last().Attr("href")

		comments := ctx.Scrape(commentsLink, func(ctx flyscrape.Context) any {
			var comments []flyscrape.M

			for _, comment := range ctx.Doc.Find(".comtr").All() {
				comments = append(comments, flyscrape.M{
					"author": comment.Find(".hnuser").Text(),
					"text":   comment.Find(".commtext").Text(),
				})

			}

			return comments
		})

		return flyscrape.M{
			"title":    title,
			"comments": comments,
		}
	}

	flyscrape.Run(cfg, fn)
}
