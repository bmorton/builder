package api

import (
	"net/http"

	"github.com/bmorton/builder/builds"
	"github.com/gin-gonic/gin"
)

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
	build, err := br.buildRepo.Find(buildID)
	if err == builds.ErrNotFound {
		c.String(http.StatusNotFound, "")
		return
	} else if err != nil {
		c.String(http.StatusInternalServerError, "")
		return
	}

	c.JSON(http.StatusOK, build)
}
