package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"strings"
)

const (
	treeAPI = "https://api.github.com/repos/github/gitignore/git/trees/main?recursive=1"
	rawBase = "https://raw.githubusercontent.com/github/gitignore/main/"
)

type treeResponse struct {
	Tree []struct {
		Path string `json:"path"`
		Type string `json:"type"`
	} `json:"tree"`
}

func fetchIndex() ([]Template, error) {
	resp, err := http.Get(treeAPI)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("github API error: %s\n%s", resp.Status, string(body))
	}

	var tr treeResponse
	if err := json.NewDecoder(resp.Body).Decode(&tr); err != nil {
		return nil, err
	}

	var out []Template
	for _, node := range tr.Tree {
		if node.Type == "blob" && strings.HasSuffix(node.Path, ".gitignore") {
			name := strings.TrimSuffix(node.Path, ".gitignore")
			name = filepath.Base(name)
			out = append(out, Template{Name: name, Path: node.Path})
		}
	}
	return out, nil
}

func fetchTemplate(t Template) (Template, error) {
	resp, err := http.Get(rawBase + t.Path)
	if err != nil {
		return Template{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return Template{}, fmt.Errorf("failed to fetch %s: %s", t.Name, resp.Status)
	}

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return Template{}, err
	}

	t.Content = string(b)
	return t, nil
}