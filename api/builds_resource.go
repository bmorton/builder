package api

import (
	"io"
	"net/http"

	"github.com/bmorton/builder/builds"
	"github.com/bmorton/flushwriter"
	"github.com/gin-gonic/gin"
)

type BuildsResource struct {
	buildRepo  *builds.Repository
	buildQueue *builds.Queue
}

func NewBuildsResource(buildRepo *builds.Repository, buildQueue *builds.Queue) *BuildsResource {
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
	br.buildQueue.Add(build)

	c.JSON(http.StatusOK, build)
}

func (br *BuildsResource) Show(c *gin.Context) {
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
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")

	var w io.Writer
	if c.Request.Header.Get("Accept") == "text/event-stream" {
		c.Writer.Header().Set("Content-Type", "text/event-stream")
		w = NewSSEWriter(c.Writer)
	} else {
		c.Writer.Header().Set("Content-Type", "text/plain")
		w = flushwriter.New(c.Writer)
	}

	build.OutputStream.Replay(w)
	build.OutputStream.Add(w, waitChan)
	<-waitChan
	build.OutputStream.Remove(w)
}
