package api

import (
	"net/http"

	"github.com/bmorton/builder/builds"
	"github.com/bmorton/flushwriter"
	"github.com/gin-gonic/gin"
)

type BuildsResource struct {
	buildRepo *builds.Repository
}

func NewBuildsResource(buildRepo *builds.Repository) *BuildsResource {
	return &BuildsResource{buildRepo: buildRepo}
}

func (br *BuildsResource) Show(c *gin.Context) {
	c.Writer.Header().Set("Content-Type", "text/plain")
	jobID := c.Params.ByName("id")

	build, ok := br.buildRepo.Find(jobID)
	if !ok || build == nil {
		c.String(http.StatusNotFound, "Not Found\n")
		return
	}

	waitChan := make(chan bool, 1)
	notify := c.Writer.CloseNotify()

	go func() {
		<-notify
		waitChan <- true
	}()

	c.Writer.WriteHeader(http.StatusOK)
	fw := flushwriter.New(c.Writer)
	build.OutputStream.Replay(&fw)
	build.OutputStream.Add(&fw, waitChan)
	<-waitChan
	build.OutputStream.Remove(&fw)
}
