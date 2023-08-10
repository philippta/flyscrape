package main

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"flyscrape/flyscrape"
	"flyscrape/js"
)

func main() {
	if len(os.Args) != 2 {
		exit("Please provide a file to run.")
	}

	opts, scrape, err := js.Compile(os.Args[1])
	if err != nil {
		exit(fmt.Sprintf("Error compiling JavaScript file: %v", err))
	}

	svc := flyscrape.Service{
		ScrapeOptions: *opts,
		ScrapeFunc:    scrape,
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
	results := svc.Scrape()
	if err != nil {
	}
	fmt.Printf("%T\n", results[0])

	data, _ := json.MarshalIndent(results, "", "  ")
	fmt.Println(string(data))
	return
}

func exit(msg string) {
	fmt.Fprintln(os.Stderr, msg)
	os.Exit(1)
}
