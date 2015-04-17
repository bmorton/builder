package api

import (
	"net/http"

	"github.com/bmorton/builder/builds"
	"github.com/gin-gonic/gin"
)

type WebhookHandler struct {
	buildRepo BuildRepository
	queue     BuildQueue
}

func NewWebhookHandler(buildRepo BuildRepository, queue BuildQueue) *WebhookHandler {
	return &WebhookHandler{buildRepo: buildRepo, queue: queue}
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

	build := builds.New(event.Repository.Name, cloneURL, event.HeadCommit.ID)
	wh.buildRepo.Create(build)
	wh.queue.Add(build)

	c.JSON(http.StatusOK, build)
}
