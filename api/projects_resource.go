package api

import (
	"net/http"

	"github.com/bmorton/builder/projects"
	"github.com/gin-gonic/gin"
)

type ProjectsResource struct {
	projectRepo ProjectRepository
}

func NewProjectsResource(projectRepo ProjectRepository) *ProjectsResource {
	return &ProjectsResource{
		projectRepo: projectRepo,
	}
}

func (r *ProjectsResource) Index(c *gin.Context) {
	records := r.projectRepo.All()
	c.JSON(http.StatusOK, records)
}

func (r *ProjectsResource) Show(c *gin.Context) {
	projectID := c.Params.ByName("id")
	record, err := r.projectRepo.Find(projectID)
	if err == projects.ErrNotFound {
		c.String(http.StatusNotFound, "")
		return
	} else if err != nil {
		c.String(http.StatusInternalServerError, "")
		return
	}

	c.JSON(http.StatusOK, record)
}
