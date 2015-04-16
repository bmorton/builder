package api

import (
	"github.com/bmorton/builder/builds"
	"github.com/bmorton/builder/streams"
)

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

type StreamRepository interface {
	Find(string) (*streams.BuildStream, error)
}
