package main

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"flyscrape/js"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "Please provide a file to run.")
		os.Exit(1)
	}

	opts, run, err := js.Compile(os.Args[1])
	if err != nil {
		panic(err)
	}

	resp, err := http.Get(opts.URL)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	out := run(js.RunOptions{HTML: string(body)})

	j, err := json.MarshalIndent(out, "", "  ")
	if err != nil {
		panic(err)
	}

	fmt.Println(string(j))
}
