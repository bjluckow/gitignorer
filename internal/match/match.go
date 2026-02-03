package match

import (
	"regexp"
	"strings"

	"github.com/bjluckow/gitignorer/internal/model"
)

func Templates(args []string, templates []model.Template) []model.Template {
	seen := map[string]bool{}
	var out []model.Template

	for _, arg := range args {
		rx := regexp.MustCompile("(?i)" + regexp.QuoteMeta(arg))

		for _, t := range templates {
			if rx.MatchString(normalize(t.Name)) {
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
	s = strings.ReplaceAll(s, "_", " ")
	s = strings.ReplaceAll(s, "-", " ")
	return s
}
