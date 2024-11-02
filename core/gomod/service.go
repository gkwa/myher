package gomod

import (
	"fmt"
	"strings"
)

type Logger interface {
	Info(msg string, keysAndValues ...interface{})
	Error(err error, msg string, keysAndValues ...interface{})
}

var NewService = func(logger Logger) Service {
	return &ModuleService{
		finder:   NewFinder(),
		reader:   NewReader(),
		parser:   NewParser(),
		versions: NewVersionFinder(),
		logger:   logger,
	}
}

type Service interface {
	GetModuleInfo() (*ModuleInfo, error)
	PrettyPrint(info *ModuleInfo)
	GenerateDowngradeCommands(concurrency int, alternatingComments bool) ([]string, error)
}

type ModuleService struct {
	finder   Finder
	reader   Reader
	parser   Parser
	versions VersionFinder
	logger   Logger
}

func (s *ModuleService) GetModuleInfo() (*ModuleInfo, error) {
	path, err := s.finder.Find(".")
	if err != nil {
		s.logger.Error(err, "Failed to find go.mod")
		return nil, err
	}

	content, err := s.reader.Read(path)
	if err != nil {
		s.logger.Error(err, "Failed to read go.mod")
		return nil, err
	}

	info, err := s.parser.Parse(content)
	if err != nil {
		s.logger.Error(err, "Failed to parse go.mod")
		return nil, err
	}

	return info, nil
}

func (s *ModuleService) PrettyPrint(info *ModuleInfo) {
	fmt.Printf("Module: %s\n", info.Module)
	fmt.Printf("Go Version: %s\n\n", info.GoVersion)
	fmt.Println("Direct Dependencies:")
	fmt.Println("-------------------")
	for _, dep := range info.DirectDeps {
		fmt.Printf("%s %s\n", dep.Path, dep.Version)
	}
}

func (s *ModuleService) GenerateDowngradeCommands(concurrency int, alternatingComments bool) ([]string, error) {
	info, err := s.GetModuleInfo()
	if err != nil {
		return nil, err
	}

	versionMap, err := s.versions.FindVersions(info.DirectDeps, concurrency)
	if err != nil {
		s.logger.Error(err, "Some errors occurred while fetching versions")
	}

	var commands []string
	for i, dep := range info.DirectDeps {
		versions, ok := versionMap[dep.Path]
		if !ok {
			continue
		}

		if previousVersion := s.versions.GetPreviousVersion(dep.Version, versions); previousVersion != "" {
			version := strings.TrimPrefix(previousVersion, "v")
			cmd := fmt.Sprintf("go get %s@v%s", dep.Path, version)
			if alternatingComments && i%2 == 1 {
				cmd = "# " + cmd
			}
			commands = append(commands, cmd)
		}
	}

	return commands, nil
}
