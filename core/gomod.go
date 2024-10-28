package core

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/go-logr/logr"
	"golang.org/x/mod/modfile"
)

func ParseGoMod(logger logr.Logger) {
	goModPath, err := findGoMod(".")
	if err != nil {
		logger.Error(err, "Failed to find go.mod")
		return
	}

	content, err := os.ReadFile(goModPath)
	if err != nil {
		logger.Error(err, "Failed to read go.mod")
		return
	}

	f, err := modfile.Parse("go.mod", content, nil)
	if err != nil {
		logger.Error(err, "Failed to parse go.mod")
		return
	}

	logger.Info(fmt.Sprintf("Module: %s", f.Module.Mod.Path))
	logger.Info(fmt.Sprintf("Go Version: %s", f.Go.Version))

	logger.Info("Direct Dependencies:")
	for _, req := range f.Require {
		if req.Indirect {
			continue
		}
		logger.Info(fmt.Sprintf("%s %s", req.Mod.Path, req.Mod.Version))
	}
}

func findGoMod(dir string) (string, error) {
	goModPath := filepath.Join(dir, "go.mod")
	if _, err := os.Stat(goModPath); err == nil {
		return goModPath, nil
	}

	parent := filepath.Dir(dir)
	if parent == dir {
		return "", fmt.Errorf("go.mod not found")
	}

	return findGoMod(parent)
}
