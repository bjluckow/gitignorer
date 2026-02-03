package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/bjluckow/gitignorer/internal/fetch"
	"github.com/bjluckow/gitignorer/internal/match"
)

func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(1)
	}

	if os.Args[1] == "list" {
		list()
		return
	} else {
		generate(os.Args[1:])
	}
}

func usage() {
	fmt.Println(`usage:
	gitignorer <template1 template2...> (e.g. "go python node")
	gitignorer list
	`)
}

func fatal(err error) {
	fmt.Fprintln(os.Stderr, "error: ", err)
	os.Exit(1)
}

func generate(args []string) {
	templates, err := fetch.LoadIndex()
	if err != nil {
		fatal(err)
	}

	matches := match.Templates(args, templates)
	if len(matches) == 0 {
		fatal(errors.New("no templates matched"))
	}

	for _, t := range matches {
		body, err := fetch.FetchTemplate(t)
		if err != nil {
			fatal(err)
		}
		fmt.Printf("# === %s ===\n", t.Name)
		fmt.Println(body)
		fmt.Println()
	}
}

func list() {
	templates, err := fetch.LoadIndex()
	if err != nil {
		fatal(err)
	}

	for _, t := range templates {
		fmt.Println(t.Name)
	}
}
