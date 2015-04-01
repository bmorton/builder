package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type WebhookHandler struct {
	queue *BuildQueue
}

func (wh *WebhookHandler) Github(c *gin.Context) {
	var event GithubPushEvent

	c.Bind(&event)
	build := &Build{
		RepositoryName: event.Repository.Name,
		CloneURL:       event.Repository.CloneURL,
		CommitID:       event.HeadCommit.ID,
		GitRef:         event.Ref,
	}
	wh.queue.Add(build)

	c.JSON(http.StatusOK, build)
}
