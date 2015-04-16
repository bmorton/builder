package api

import (
	"net/http"

	"github.com/bmorton/builder/builds"
	"github.com/bmorton/builder/streams"
	"github.com/gin-gonic/gin"
)

type LogsResource struct {
	buildRepo *builds.Repository
	logRepo   *builds.LogRepository
}

func NewLogsResource(buildRepo *builds.Repository, logRepo *builds.LogRepository) *LogsResource {
	return &LogsResource{buildRepo: buildRepo, logRepo: logRepo}
}

func (r *LogsResource) Show(c *gin.Context) {
	buildID := c.Params.ByName("id")
	logType := c.Params.ByName("type")

	build, err := r.buildRepo.Find(buildID)
	if err == builds.ErrNotFound {
		c.String(http.StatusNotFound, "")
		return
	} else if err != nil {
		c.String(http.StatusInternalServerError, "")
		return
	}

	buildLog, err := r.logRepo.FindByBuildID(build.ID, logType)
	if err == streams.ErrNotFound {
		c.String(http.StatusNotFound, "")
		return
	} else if err != nil {
		c.String(http.StatusInternalServerError, "")
		return
	}

	c.Writer.WriteHeader(http.StatusOK)
	c.Writer.Header().Set("Content-Type", "text/plain")

	c.String(http.StatusOK, buildLog.Data)
}
