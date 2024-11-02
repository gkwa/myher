package gomod

import (
	"fmt"
	"os/exec"
	"strings"
	"sync"
)

type VersionFinder interface {
	FindVersions(modules []Dependency, concurrency int) (map[string][]string, error)
	GetPreviousVersion(current string, versions []string) string
}

type GoListVersionFinder struct{}

func NewVersionFinder() VersionFinder {
	return &GoListVersionFinder{}
}

func (f *GoListVersionFinder) FindVersions(modules []Dependency, concurrency int) (map[string][]string, error) {
	results := make(map[string][]string)
	resultsMux := sync.Mutex{}
	errorsChan := make(chan error, len(modules))

	semaphore := make(chan struct{}, concurrency)
	var wg sync.WaitGroup

	for _, mod := range modules {
		wg.Add(1)
		go func(module Dependency) {
			defer wg.Done()
			semaphore <- struct{}{}        // Acquire
			defer func() { <-semaphore }() // Release

			cmd := exec.Command("go", "list", "-m", "-versions", module.Path)
			output, err := cmd.Output()
			if err != nil {
				errorsChan <- fmt.Errorf("failed to get versions for %s: %v", module.Path, err)
				return
			}

			parts := strings.Fields(string(output))
			if len(parts) <= 1 {
				errorsChan <- fmt.Errorf("no versions found for %s", module.Path)
				return
			}

			resultsMux.Lock()
			results[module.Path] = parts[1:]
			resultsMux.Unlock()
		}(mod)
	}

	wg.Wait()
	close(errorsChan)

	var errors []error
	for err := range errorsChan {
		if err != nil {
			errors = append(errors, err)
		}
	}

	if len(errors) > 0 {
		return results, fmt.Errorf("errors occurred while fetching versions: %v", errors)
	}

	return results, nil
}

func (f *GoListVersionFinder) GetPreviousVersion(current string, versions []string) string {
	for i := len(versions) - 1; i > 0; i-- {
		if versions[i] == current {
			return versions[i-1]
		}
	}
	return ""
}
