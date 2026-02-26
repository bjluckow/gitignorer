package match

import (
	"strings"

	"github.com/bjluckow/gitignorer/internal/model"
)

func Templates(args []string, templates []model.Template) []model.Template {
	seen := map[string]bool{}
	var out []model.Template

	for _, arg := range args {
		argNorm := normalize(arg)

		for _, t := range templates {
			if normalize(t.Name) == argNorm {
				if !seen[t.Path] {
					seen[t.Path] = true
					out = append(out, t)
				}
			}
		}
	}

	return out
}

func normalize(s string) string {
	s = strings.ToLower(s)
	s = strings.ReplaceAll(s, "_", "")
	s = strings.ReplaceAll(s, "-", "")
	s = strings.ReplaceAll(s, " ", "")
	return s
}
