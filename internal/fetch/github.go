package fetch

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/bjluckow/gitignorer/internal/model"
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

func LoadIndex() ([]model.Template, error) {
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

	var out []model.Template
	for _, node := range tr.Tree {
		if node.Type == "blob" && strings.HasSuffix(node.Path, ".gitignore") {
			name := strings.TrimSuffix(node.Path, ".gitignore")
			name = filepathBase(name)

			out = append(out, model.Template{
				Name: name,
				Path: node.Path,
			})
		}
	}

	return out, nil
}

func filepathBase(p string) string {
	parts := strings.Split(p, "/")
	return parts[len(parts)-1]
}

func FetchTemplate(t model.Template) (string, error) {
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
