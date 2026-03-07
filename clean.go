package main

import (
	"fmt"
	"io"
	"strings"
)

func clean(input string, out io.Writer) {
	sections := parseSections(input)

	// count how many sections each line appears in
	lineCounts := make(map[string]int)
	for _, lines := range sections {
		seen := make(map[string]bool)
		for _, line := range lines {
			if line == "" || strings.HasPrefix(line, "#") {
				continue
			}
			if !seen[line] {
				lineCounts[line]++
				seen[line] = true
			}
		}
	}

	var common, custom []string
	commonSeen := make(map[string]bool)

	for _, lines := range sections {
		for _, line := range lines {
			if line == "" || strings.HasPrefix(line, "#") {
				continue
			}
			if lineCounts[line] > 1 && !commonSeen[line] {
				common = append(common, line)
				commonSeen[line] = true
			}
		}
	}

	// custom: lines unique to their section and not already in common
	customSeen := make(map[string]bool)
	for _, lines := range sections {
		for _, line := range lines {
			if line == "" || strings.HasPrefix(line, "#") {
				continue
			}
			if lineCounts[line] == 1 && !customSeen[line] {
				custom = append(custom, line)
				customSeen[line] = true
			}
		}
	}

	if len(common) > 0 {
		fmt.Fprintln(out, "# === common ===")
		for _, line := range common {
			fmt.Fprintln(out, line)
		}
		fmt.Fprintln(out)
	}

	if len(custom) > 0 {
		fmt.Fprintln(out, "# === custom ===")
		for _, line := range custom {
			fmt.Fprintln(out, line)
		}
		fmt.Fprintln(out)
	}
}

// parseSections splits a gitignore into per-section line slices, keyed by header
func parseSections(input string) map[string][]string {
	sections := make(map[string][]string)
	current := "custom"
	for _, line := range strings.Split(input, "\n") {
		line = strings.TrimRight(line, "\r")
		if strings.HasPrefix(line, "# ===") {
			current = line
		} else {
			sections[current] = append(sections[current], line)
		}
	}
	return sections
}