package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/bjluckow/gitignorer/internal/fetch"
	"github.com/bjluckow/gitignorer/internal/match"
)

func main() {
	writeFile := flag.Bool("o", false, "write output to ./.gitignore")
	appendFile := flag.Bool("a", false, "append output to ./gitignore")
	flag.Usage = usage
	flag.Parse()

	args := flag.Args()

	if len(args) < 1 {
		usage()
		os.Exit(1)
	}

	if args[0] == "list" {
		list()
		return
	}

	var out io.Writer = os.Stdout
	if *writeFile {
		f, err := os.Create(".gitignore")
		if err != nil {
			fatal(err)
		}
		defer f.Close()
		out = f
	} else if *appendFile {
		f, err := os.OpenFile(".gitignore", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			fatal(err)
		}
		defer f.Close()
		out = f
	}
	generate(os.Args[1:], out)
}

func usage() {
	fmt.Println(`usage:
	gitignorer list
	gitignorer <template1 template2...> (e.g. "go python node")
		-o  write output to ./.gitignore
		-a  append output to ./.gitignore
	`)
}

func fatal(err error) {
	fmt.Fprintln(os.Stderr, "error: ", err)
	os.Exit(1)
}

func generate(args []string, out io.Writer) {
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
		fmt.Fprintf(out, "# === %s ===\n", t.Name)
		fmt.Fprintln(out, body)
		fmt.Fprintln(out)
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
