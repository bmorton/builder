package api

import (
	"io"
	"net/http"

	"github.com/bmorton/builder/builds"
	"github.com/bmorton/builder/streams"
	"github.com/bmorton/flushwriter"
	"github.com/gin-gonic/gin"
)

type StreamsResource struct {
	buildRepo  BuildRepository
	streamRepo StreamRepository
}

func NewStreamsResource(buildRepo BuildRepository, streamRepo StreamRepository) *StreamsResource {
	return &StreamsResource{buildRepo: buildRepo, streamRepo: streamRepo}
}

func (sr *StreamsResource) Show(c *gin.Context) {
	buildID := c.Params.ByName("id")
	streamType := c.Params.ByName("type")

	build, err := sr.buildRepo.Find(buildID)
	if err == builds.ErrNotFound {
		c.String(http.StatusNotFound, "")
		return
	} else if err != nil {
		c.String(http.StatusInternalServerError, "")
		return
	}

	stream, err := sr.streamRepo.Find(build.ID)
	if err == streams.ErrNotFound {
		c.String(http.StatusNotFound, "")
		return
	} else if err != nil {
		c.String(http.StatusInternalServerError, "")
		return
	}

	switch streamType {
	case "build":
		streamOutput(c, stream.BuildOutput, true)
	case "push":
		streamOutput(c, stream.PushOutput, false)
	default:
		c.String(http.StatusNotFound, "")
	}
}

func streamOutput(c *gin.Context, output *streams.Output, replay bool) {
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

	waitChan := make(chan bool, 1)
	notify := c.Writer.CloseNotify()

	go func() {
		<-notify
		waitChan <- true
	}()

	if replay {
		output.Replay(w)
	}
	output.Add(w, waitChan)
	<-waitChan
	output.Remove(w)
}
