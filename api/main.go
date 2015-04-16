package api

import "github.com/bmorton/builder/builds"

type BuildRepository interface {
	All() []*builds.Build
	Create(*builds.Build)
	Save(*builds.Build)
	Find(string) (*builds.Build, error)
}

type BuildQueue interface {
	Add(*builds.Build) string
}

type BuildLogRepository interface {
	FindByBuildID(string, string) (*builds.BuildLog, error)
}
