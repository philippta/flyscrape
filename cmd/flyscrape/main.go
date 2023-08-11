package main

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"flyscrape"
)

func main() {
	if len(os.Args) != 2 {
		exit("Please provide a file to run.")
	}

	src, err := os.ReadFile(os.Args[1])
	if err != nil {
		exit(fmt.Sprintf("Error reading file: %v", err))
	}

	opts, scrape, err := flyscrape.Compile(string(src))
	if err != nil {
		exit(fmt.Sprintf("Error compiling JavaScript file: %v", err))
	}

	svc := flyscrape.Scraper{
		ScrapeOptions: opts,
		ScrapeFunc:    scrape,
		Concurrency:   5,
		FetchFunc: func(url string) (string, error) {
			resp, err := http.Get(url)
			if err != nil {
				return "", err
			}
			defer resp.Body.Close()

			data, err := io.ReadAll(resp.Body)
			if err != nil {
				return "", err
			}
			return string(data), nil
		},
	}

	count := 0
	start := time.Now()
	for result := range svc.Scrape() {
		data, _ := json.MarshalIndent(result, "", "  ")
		fmt.Println(string(data))
		count++
	}
	fmt.Printf("Scraped %d websites in %v\n", count, time.Since(start))
}

func exit(msg string) {
	fmt.Fprintln(os.Stderr, msg)
	os.Exit(1)
}
