package main

import (
	"net/http"

	"github.com/bmorton/flushwriter"
	"github.com/gin-gonic/gin"
)

type BuildsResource struct {
	builds *JobRepository
}

func (br *BuildsResource) Show(c *gin.Context) {
	c.Writer.Header().Set("Content-Type", "text/plain")
	jobID := c.Params.ByName("id")

	writer, ok := br.builds.Find(jobID)
	if !ok || writer == nil {
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
	writer.Replay(&fw)
	writer.Add(&fw, waitChan)
	<-waitChan
	writer.Remove(&fw)
}
