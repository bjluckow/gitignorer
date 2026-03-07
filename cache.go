package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

const cacheFile = "gitignorer/cache.json"

func cachePath() string {
	dir, err := os.UserCacheDir()
	if err != nil {
		return ".gitignorer_cache.json"
	}
	return filepath.Join(dir, cacheFile)
}

func loadCache() (Cache, error) {
	f, err := os.Open(cachePath())
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var c Cache
	if err := json.NewDecoder(f).Decode(&c); err != nil {
		return nil, err
	}
	return c, nil
}

func saveCache(c Cache) error {
	p := cachePath()
	if err := os.MkdirAll(filepath.Dir(p), 0755); err != nil {
		return err
	}
	f, err := os.Create(p)
	if err != nil {
		return err
	}
	defer f.Close()
	return json.NewEncoder(f).Encode(c)
}

func loadOrFetch(refresh bool) (Cache, error) {
	if !refresh {
		if c, err := loadCache(); err == nil && len(c) > 0 {
			return c, nil
		}
	}

	templates, err := fetchIndex()
	if err != nil {
		return nil, err
	}

	c := make(Cache)
	for _, t := range templates {
		c[t.Name] = t
	}
	if err := saveCache(c); err != nil {
		fmt.Fprintln(os.Stderr, "warning: could not save cache:", err)
	}
	return c, nil
}