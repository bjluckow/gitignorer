package main

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
)

const (
	repoAPI = "https://api.github.com/repos/github/gitignore/contents"
	rawBase = "https://raw.githubusercontent.com/github/gitignore/main/"
)

type Template struct {
	Name string
	Path string
}

type ghEntry struct {
	Name string `json:"name"`
	Path string `json:"path"`
	Type string `json:"type"`
	URL  string `json:"url"`
}

func LoadTemplateIndex() ([]Template, error) {
	var out []Template
	if err := walkDir(repoAPI, &out); err != nil {
		return nil, err
	}
	return out, nil
}
func walkDir(url string, out *[]Template) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var entries []ghEntry
	if err := json.NewDecoder(resp.Body).Decode(&entries); err != nil {
		return err
	}

	for _, e := range entries {
		switch e.Type {
		case "file":
			if strings.HasSuffix(e.Name, ".gitignore") {
				*out = append(*out, Template{
					Name: strings.TrimSuffix(e.Name, ".gitignore"),
					Path: e.Path,
				})
			}
		case "dir":
			if err := walkDir(e.URL, out); err != nil {
				return err
			}
		}
	}
	return nil
}

func FetchTemplate(t Template) (string, error) {
	resp, err := http.Get(rawBase + t.Path)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(b), nil
}
