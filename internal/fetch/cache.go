package fetch

import (
	"os"
	"path/filepath"
)

func getCachePath() (string, error) {
	if dir := os.Getenv("XDG_CACHE_HOME"); dir != "" {
		return filepath.Join(dir, "gitignorer", "index.json"), nil
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(home, ".cache", "gitignorer", "index.json"), nil
}

// TODO: impl LoadIndexCached (+ writeCache, loadCache)
