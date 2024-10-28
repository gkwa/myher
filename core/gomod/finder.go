package gomod

import (
	"fmt"
	"os"
	"path/filepath"
)

type Finder interface {
	Find(startDir string) (string, error)
}

type FileFinder struct{}

func NewFinder() Finder {
	return &FileFinder{}
}

func (f *FileFinder) Find(startDir string) (string, error) {
	goModPath := filepath.Join(startDir, "go.mod")
	if _, err := os.Stat(goModPath); err == nil {
		return goModPath, nil
	}

	parent := filepath.Dir(startDir)
	if parent == startDir {
		return "", fmt.Errorf("go.mod not found")
	}

	return f.Find(parent)
}
