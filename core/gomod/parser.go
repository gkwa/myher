package gomod

import (
	"golang.org/x/mod/modfile"
)

type Parser interface {
	Parse(content []byte) (*ModuleInfo, error)
}

type ModFileParser struct{}

func NewParser() Parser {
	return &ModFileParser{}
}

func (p *ModFileParser) Parse(content []byte) (*ModuleInfo, error) {
	f, err := modfile.Parse("go.mod", content, nil)
	if err != nil {
		return nil, err
	}

	info := &ModuleInfo{
		Module:    f.Module.Mod.Path,
		GoVersion: f.Go.Version,
	}

	for _, req := range f.Require {
		if !req.Indirect {
			info.DirectDeps = append(info.DirectDeps, Dependency{
				Path:    req.Mod.Path,
				Version: req.Mod.Version,
			})
		}
	}

	return info, nil
}
