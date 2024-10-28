package gomod

type Dependency struct {
	Path    string
	Version string
}

type ModuleInfo struct {
	Module     string
	GoVersion  string
	DirectDeps []Dependency
}
