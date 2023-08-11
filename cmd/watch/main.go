package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"flyscrape"

	"github.com/cornelk/hashmap"
	"github.com/inancgumus/screen"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Please provide a file to run.")
		os.Exit(1)
	}

	cache := hashmap.New[string, string]()

	err := flyscrape.Watch(os.Args[1], func(s string) error {
		opts, scrape, err := flyscrape.Compile(s)
		if err == nil {
			run(cache, opts, scrape)
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
}

func run(cache *hashmap.Map[string, string], opts flyscrape.ScrapeOptions, fn flyscrape.ScrapeFunc) {
	opts.Depth = 0

	svc := flyscrape.Scraper{
		Concurrency:   20,
		ScrapeOptions: opts,
		ScrapeFunc:    fn,
		FetchFunc: func(url string) (string, error) {
			if html, ok := cache.Get(url); ok {
				return html, nil
			}
			html, err := fetch(url)
			if err != nil {
				return "", err
			}
			cache.Set(url, html)
			return html, nil
		},
	}

	result := <-svc.Scrape()
	if result.Error != nil {
		fmt.Println(result.Error)
	}

	screen.Clear()
	screen.MoveTopLeft()

	enc := json.NewEncoder(os.Stdout)
	enc.SetEscapeHTML(false)
	enc.SetIndent("", "   ")
	enc.Encode(result)
}

func fetch(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}
