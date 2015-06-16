package api

import (
	"net/http"

	"github.com/bmorton/builder/builds"
	"github.com/gin-gonic/gin"
)

type WebhookHandler struct {
	buildRepo   BuildRepository
	projectRepo ProjectRepository
	queue       BuildQueue
}

func NewWebhookHandler(buildRepo BuildRepository, projectRepo ProjectRepository, queue BuildQueue) *WebhookHandler {
	return &WebhookHandler{buildRepo: buildRepo, projectRepo: projectRepo, queue: queue}
}

func (wh *WebhookHandler) Github(c *gin.Context) {
	var event GithubPushEvent
	var cloneURL string

	c.BindJSON(&event)

	if event.Ref != "refs/heads/master" {
		c.String(http.StatusBadRequest, "Only builds of the master branch are currently supported.\n")
		return
	}

	if event.Repository.CloneURL != "" {
		cloneURL = event.Repository.CloneURL
	} else {
		cloneURL = event.Repository.URL
	}

	project := wh.projectRepo.FindOrCreateByCloneURL(cloneURL)
	build := builds.New(event.Repository.Name, cloneURL, event.HeadCommit.ID)
	build.ProjectID = project.ID

	wh.buildRepo.Create(build)
	wh.queue.Add(build)

	c.JSON(http.StatusOK, build)
}
