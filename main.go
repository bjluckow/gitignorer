package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
)

func main() {
	flag.Usage = globalUsage
	flag.Parse()

	if flag.NArg() < 1 {
		globalUsage()
		os.Exit(1)
	}

	switch flag.Arg(0) {
	case "fetch":
		runFetch(flag.Args()[1:])
	case "list":
		runList(flag.Args()[1:])
	case "clean":
		runClean(flag.Args()[1:])
	case "cache":
    	runCache(flag.Args()[1:])
	default:
		fmt.Fprintf(os.Stderr, "unknown subcommand: %s\n", flag.Arg(0))
		globalUsage()
		os.Exit(1)
	}
}

func runFetch(args []string) {
	fs := flag.NewFlagSet("fetch", flag.ExitOnError)
	write   := fs.Bool("w", false, "write output to ./.gitignore")
	append_ := fs.Bool("a", false, "append output to ./.gitignore")
	refresh := fs.Bool("r", false, "refresh cached templates")
	doClean := fs.Bool("c", false, "clean output before writing")
	outPath := fs.String("o", "", "write output to a custom path")
	fs.Usage = func() {
		fmt.Println("usage: gitignorer fetch [-w] [-a] [-r] [-c] [-o path] <template1 template2 ...>")
		fs.PrintDefaults()
	}
	fs.Parse(args)

	templates := fs.Args()
	if len(templates) == 0 {
		fs.Usage()
		os.Exit(1)
	}

	cache, err := loadOrFetch(*refresh)
	if err != nil {
		fatal(err)
	}

	matches := matchTemplates(templates, cache)
	if len(matches) == 0 {
		fatal(errors.New("no templates matched"))
	}

	// resolve output writer
	var out io.Writer = os.Stdout
	if *outPath != "" {
		f, err := os.Create(*outPath)
		if err != nil {
			fatal(err)
		}
		defer f.Close()
		out = f
	} else if *write {
		f, err := os.Create(".gitignore")
		if err != nil {
			fatal(err)
		}
		defer f.Close()
		out = f
	} else if *append_ {
		f, err := os.OpenFile(".gitignore", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			fatal(err)
		}
		defer f.Close()
		out = f
	}

	var sb strings.Builder
	for _, t := range matches {
		if t.Content == "" || *refresh {
			t, err = fetchTemplate(t)
			if err != nil {
				fatal(err)
			}
			cache[t.Name] = t
		}
		fmt.Fprintf(&sb, "# === %s ===\n", t.Name)
		fmt.Fprintln(&sb, t.Content)
		fmt.Fprintln(&sb)
	}

	if err := saveCache(cache); err != nil {
		fmt.Fprintln(os.Stderr, "warning: could not save cache:", err)
	}

	output := sb.String()
	if *doClean {
		clean(output, out)
	} else {
		fmt.Fprint(out, output)
	}
}

func runList(args []string) {
	fs := flag.NewFlagSet("list", flag.ExitOnError)
	refresh := fs.Bool("r", false, "refresh cached index")
	fs.Usage = func() {
		fmt.Println("usage: gitignorer list [-r]")
		fs.PrintDefaults()
	}
	fs.Parse(args)

	cache, err := loadOrFetch(*refresh)
	if err != nil {
		fatal(err)
	}

	for name := range cache {
		fmt.Println(name)
	}
}

func runClean(args []string) {
	fs := flag.NewFlagSet("clean", flag.ExitOnError)
	fs.Usage = func() {
		fmt.Println("usage: gitignorer clean [path]")
		fs.PrintDefaults()
	}
	fs.Parse(args)

	path := ".gitignore"
	if fs.NArg() > 0 {
		path = fs.Arg(0)
	}

	b, err := os.ReadFile(path)
	if err != nil {
		fatal(err)
	}

	clean(string(b), os.Stdout)
}

func runCache(args []string) {
    fmt.Println(cachePath())
}

func matchTemplates(args []string, cache Cache) []Template {
	var matches []Template
	for _, arg := range args {
		arg = strings.ToLower(arg)
		for name, t := range cache {
			if strings.ToLower(name) == arg {
				matches = append(matches, t)
				break
			}
		}
	}
	return matches
}

func globalUsage() {
	fmt.Println(`usage: gitignorer <subcommand> [flags] [args]

subcommands:
  fetch   fetch gitignore templates
  list    list available templates
  clean   deduplicate and reorganize a .gitignore`)
}

func fatal(err error) {
	fmt.Fprintln(os.Stderr, "error:", err)
	os.Exit(1)
}