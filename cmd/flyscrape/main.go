package main

import (
	_ "embed"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

func main() {
	log.SetFlags(0)

	m := &Main{}
	if err := m.Run(os.Args[1:]); err == flag.ErrHelp {
		os.Exit(1)
	} else if err != nil {
		log.Println(err)
		os.Exit(1)
	}
}

type Main struct{}

func (m *Main) Run(args []string) error {
	var cmd string
	if len(args) > 0 {
		cmd, args = args[0], args[1:]
	}

	switch cmd {
	case "new":
		return (&NewCommand{}).Run(args)
	case "run":
		return (&RunCommand{}).Run(args)
	case "watch":
		return (&WatchCommand{}).Run(args)
	default:
		if cmd == "" || cmd == "help" || strings.HasPrefix(cmd, "-") {
			m.Usage()
			return flag.ErrHelp
		}
		return fmt.Errorf("flyscrape %s: unknown command", cmd)
	}
}

func (m *Main) Usage() {
	fmt.Println(`
flyscrape is an elegant scraping tool for efficiently extracting data from websites.

Usage:

    flyscrape <command> [arguments]

Commands:
    
    new    creates a sample scraping script
    run    runs a scraping script
    watch  watches and re-runs a scraping script
`[1:])
}
