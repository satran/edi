package main

import (
	"path/filepath"
)

func relativePath(parent, child string) string {
	if !filepath.IsAbs(parent) {
		parent, _ = filepath.Abs(parent)
	}
	relDir, _ := filepath.Rel(parent, child)

	return relDir
}
