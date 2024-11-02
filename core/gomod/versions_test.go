package gomod

import (
	"testing"
)

func TestGoListVersionFinder_GetPreviousVersion(t *testing.T) {
	finder := &GoListVersionFinder{}

	tests := []struct {
		name     string
		current  string
		versions []string
		want     string
	}{
		{
			name:     "multiple versions",
			current:  "v0.3.0",
			versions: []string{"v0.1.0", "v0.2.0", "v0.3.0", "v0.4.0"},
			want:     "v0.2.0",
		},
		{
			name:     "current is oldest version",
			current:  "v0.1.0",
			versions: []string{"v0.1.0", "v0.2.0"},
			want:     "",
		},
		{
			name:     "current version not found",
			current:  "v9.9.9",
			versions: []string{"v0.1.0", "v0.2.0"},
			want:     "",
		},
		{
			name:     "empty versions",
			current:  "v1.0.0",
			versions: []string{},
			want:     "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := finder.GetPreviousVersion(tt.current, tt.versions); got != tt.want {
				t.Errorf("GetPreviousVersion() = %v, want %v", got, tt.want)
			}
		})
	}
}
