package gomod

import (
	"testing"
)

func TestGoListVersionFinder_GetSecondLatest(t *testing.T) {
	finder := &GoListVersionFinder{}

	tests := []struct {
		name     string
		versions []string
		want     string
	}{
		{
			name:     "multiple versions",
			versions: []string{"v0.1.0", "v0.2.0", "v0.3.0", "v0.4.0"},
			want:     "v0.3.0",
		},
		{
			name:     "only one version",
			versions: []string{"v0.1.0"},
			want:     "",
		},
		{
			name:     "empty versions",
			versions: []string{},
			want:     "",
		},
		{
			name:     "two versions",
			versions: []string{"v0.1.0", "v0.2.0"},
			want:     "v0.1.0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := finder.GetSecondLatest(tt.versions); got != tt.want {
				t.Errorf("GetSecondLatest() = %v, want %v", got, tt.want)
			}
		})
	}
}
