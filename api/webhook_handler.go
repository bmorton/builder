package api

import (
	"net/http"

	"github.com/bmorton/builder/builds"
	"github.com/gin-gonic/gin"
)

type WebhookHandler struct {
	queue *builds.Queue
}

func NewWebhookHandler(queue *builds.Queue) *WebhookHandler {
	return &WebhookHandler{queue: queue}
}

func (wh *WebhookHandler) Github(c *gin.Context) {
	var event GithubPushEvent
	var cloneURL string

	c.Bind(&event)

	if event.Ref != "refs/heads/master" {
		c.String(http.StatusBadRequest, "Only builds of the master branch are currently supported.\n")
		return
	}

	if event.Repository.CloneURL != "" {
		cloneURL = event.Repository.CloneURL
	} else {
		cloneURL = event.Repository.URL
	}

	build := &builds.Build{
		RepositoryName: event.Repository.Name,
		CloneURL:       cloneURL,
		CommitID:       event.HeadCommit.ID,
		GitRef:         event.Ref,
	}
	wh.queue.Add(build)

	c.JSON(http.StatusOK, build)
}
