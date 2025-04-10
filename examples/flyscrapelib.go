package main

import "github.com/philippta/flyscrape/flyscrape"

var config = flyscrape.Config{
	URL:      "https://news.ycombinator.com/",
	Browser:  true,
	Headless: false,
}

func scrape(ctx flyscrape.Context) any {
	var (
		post         = ctx.Doc.Find(".athing.submission").First()
		title        = post.Find(".titleline > a").Text()
		commentsLink = post.Next().Find("a").Last().Attr("href")
		comments     = ctx.Scrape(commentsLink, scrapeComments)
	)

	return flyscrape.M{
		"title":    title,
		"comments": comments,
	}
}

func scrapeComments(ctx flyscrape.Context) any {
	var comments []flyscrape.M

	for _, comment := range ctx.Doc.Find(".comtr").All() {
		comments = append(comments, flyscrape.M{
			"author": comment.Find(".hnuser").Text(),
			"text":   comment.Find(".commtext").Text(),
		})

	}

	return comments
}

func main() {
	flyscrape.Run(config, scrape)
}
