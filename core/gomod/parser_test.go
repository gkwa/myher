// core/gomod/parser_test.go
package gomod

import (
	"os"
	"testing"
)

func TestParser_Parse(t *testing.T) {
	content, err := os.ReadFile("testdata/go.mod")
	if err != nil {
		t.Fatalf("Failed to read test file: %v", err)
	}

	parser := NewParser()
	got, err := parser.Parse(content)
	if err != nil {
		t.Fatalf("Parser.Parse() error = %v", err)
	}

	want := &ModuleInfo{
		Module:    "github.com/gkwa/myher",
		GoVersion: "1.21",
		DirectDeps: []Dependency{
			{Path: "github.com/fatih/color", Version: "v1.18.0"},
			{Path: "github.com/go-git/go-git/v5", Version: "v5.11.0"},
			{Path: "github.com/go-logr/logr", Version: "v1.3.0"},
			{Path: "github.com/spf13/cobra", Version: "v1.8.1"},
		},
	}

	if got.Module != want.Module {
		t.Errorf("Module = %v, want %v", got.Module, want.Module)
	}
	if got.GoVersion != want.GoVersion {
		t.Errorf("GoVersion = %v, want %v", got.GoVersion, want.GoVersion)
	}
	if len(got.DirectDeps) != len(want.DirectDeps) {
		t.Errorf("DirectDeps length = %v, want %v", len(got.DirectDeps), len(want.DirectDeps))
		return
	}

	for i, dep := range got.DirectDeps {
		if dep.Path != want.DirectDeps[i].Path {
			t.Errorf("DirectDeps[%d].Path = %v, want %v", i, dep.Path, want.DirectDeps[i].Path)
		}
		if dep.Version != want.DirectDeps[i].Version {
			t.Errorf("DirectDeps[%d].Version = %v, want %v", i, dep.Version, want.DirectDeps[i].Version)
		}
	}
}