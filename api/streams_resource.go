package api

import (
	"fmt"
	"io"
	"net/http"

	"github.com/bmorton/builder/builds"
	"github.com/bmorton/flushwriter"
	"github.com/gin-gonic/gin"
)

type StreamsResource struct {
	buildRepo *builds.Repository
}

func NewStreamsResource(buildRepo *builds.Repository) *StreamsResource {
	return &StreamsResource{buildRepo: buildRepo}
}

func (sr *StreamsResource) Build(c *gin.Context) {
	buildID := c.Params.ByName("id")
	build := sr.buildRepo.Find(buildID)

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

	if build.IsFinished() {
		log := sr.buildRepo.FindBuildLog(build.ID)
		fmt.Fprint(w, log.Data)
	} else {
		build.BuildStream.Replay(w)
		build.BuildStream.Add(w, waitChan)
		<-waitChan
		build.BuildStream.Remove(w)
	}
}

func (sr *StreamsResource) Push(c *gin.Context) {
	buildID := c.Params.ByName("id")
	build := sr.buildRepo.Find(buildID)

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

	build.PushStream.Add(w, waitChan)
	<-waitChan
	build.PushStream.Remove(w)
}
