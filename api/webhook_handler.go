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

	c.Bind(&event)
	build := &builds.Build{
		RepositoryName: event.Repository.Name,
		CloneURL:       event.Repository.CloneURL,
		CommitID:       event.HeadCommit.ID,
		GitRef:         event.Ref,
	}
	wh.queue.Add(build)

	c.JSON(http.StatusOK, build)
}
