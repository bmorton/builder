package api

import (
	"net/http"

	"github.com/bmorton/builder/builds"
	"github.com/gin-gonic/gin"
)

type BuildRepository interface {
	All() []*builds.Build
	Create(*builds.Build)
	Save(*builds.Build)
	Find(string) *builds.Build
}

type BuildQueue interface {
	Add(*builds.Build) string
}

type BuildsResource struct {
	buildRepo  BuildRepository
	buildQueue BuildQueue
}

func NewBuildsResource(buildRepo BuildRepository, buildQueue BuildQueue) *BuildsResource {
	return &BuildsResource{
		buildRepo:  buildRepo,
		buildQueue: buildQueue,
	}
}

func (br *BuildsResource) Index(c *gin.Context) {
	builds := br.buildRepo.All()
	c.JSON(http.StatusOK, builds)
}

func (br *BuildsResource) Create(c *gin.Context) {
	build := &builds.Build{}

	c.Bind(build)
	if build.RepositoryName == "" {
		build.SetDefaultName()
	}

	br.buildRepo.Create(build)
	br.buildQueue.Add(build)

	c.JSON(http.StatusCreated, build)
}

func (br *BuildsResource) Show(c *gin.Context) {
	buildID := c.Params.ByName("id")
	build := br.buildRepo.Find(buildID)

	c.JSON(http.StatusOK, build)
}
