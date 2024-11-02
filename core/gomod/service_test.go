package gomod

import (
	"testing"

	"github.com/go-logr/logr/testr"
)

type mockFinder struct {
	path string
	err  error
}

func (m *mockFinder) Find(string) (string, error) {
	return m.path, m.err
}

type mockReader struct {
	content []byte
	err     error
}

func (m *mockReader) Read(string) ([]byte, error) {
	return m.content, m.err
}

type mockVersionFinder struct {
	versions map[string][]string
	err      error
}

func (m *mockVersionFinder) FindVersions(deps []Dependency, concurrency int) (map[string][]string, error) {
	return m.versions, m.err
}

func (m *mockVersionFinder) GetPreviousVersion(current string, versions []string) string {
	for i := len(versions) - 1; i > 0; i-- {
		if versions[i] == current {
			return versions[i-1]
		}
	}
	return ""
}

func TestService_GenerateDowngradeCommands(t *testing.T) {
	testLogger := testr.New(t)

	modContent := []byte(`
module test
go 1.21
require (
  github.com/pkg/errors v0.9.1
  github.com/stretchr/testify v1.8.4
)`)

	mockVersions := map[string][]string{
		"github.com/pkg/errors":       {"v0.8.0", "v0.8.1", "v0.9.0", "v0.9.1"},
		"github.com/stretchr/testify": {"v1.7.0", "v1.8.0", "v1.8.3", "v1.8.4"},
	}

	svc := &ModuleService{
		finder:   &mockFinder{path: "go.mod"},
		reader:   &mockReader{content: modContent},
		parser:   NewParser(),
		versions: &mockVersionFinder{versions: mockVersions},
		logger:   testLogger,
	}

	tests := []struct {
		name        string
		concurrency int
		alternating bool
		want        []string
	}{
		{
			name:        "without alternating comments",
			concurrency: 5,
			alternating: false,
			want: []string{
				"go get github.com/pkg/errors@v0.9.0",
				"go get github.com/stretchr/testify@v1.8.3",
			},
		},
		{
			name:        "with alternating comments",
			concurrency: 5,
			alternating: true,
			want: []string{
				"go get github.com/pkg/errors@v0.9.0",
				"# go get github.com/stretchr/testify@v1.8.3",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := svc.GenerateDowngradeCommands(tt.concurrency, tt.alternating)
			if err != nil {
				t.Fatalf("GenerateDowngradeCommands() error = %v", err)
			}

			if len(got) != len(tt.want) {
				t.Errorf("GenerateDowngradeCommands() returned %d commands, want %d", len(got), len(tt.want))
				return
			}

			for i, cmd := range got {
				if cmd != tt.want[i] {
					t.Errorf("Command[%d] = %v, want %v", i, cmd, tt.want[i])
				}
			}
		})
	}
}
