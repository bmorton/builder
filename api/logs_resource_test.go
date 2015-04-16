package api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bmorton/builder/api/mocks"
	"github.com/bmorton/builder/builds"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func logsResourceWithMocks() (*LogsResource, *mocks.BuildRepository, *mocks.BuildLogRepository) {
	buildRepo := new(mocks.BuildRepository)
	logRepo := new(mocks.BuildLogRepository)
	b := NewLogsResource(buildRepo, logRepo)
	return b, buildRepo, logRepo
}

func TestLogsShow(t *testing.T) {
	b, buildRepo, logRepo := logsResourceWithMocks()
	buildRepo.On("Find", "abc123").Return(&builds.Build{ID: "abc123"}, nil)
	logRepo.On("FindByBuildID", "abc123", "build").Return(&builds.BuildLog{Data: "logged!"}, nil)

	req, _ := http.NewRequest("GET", "/builds/abc123/logs/build", nil)
	w := httptest.NewRecorder()

	r := gin.New()
	r.GET("/builds/:id/logs/:type", b.Show)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "logged!")
	buildRepo.AssertExpectations(t)
	logRepo.AssertExpectations(t)
}

func TestLogsShowNotFound(t *testing.T) {
	b, buildRepo, logRepo := logsResourceWithMocks()
	buildRepo.On("Find", "abc123").Return(&builds.Build{ID: "abc123"}, nil)
	logRepo.On("FindByBuildID", "abc123", "build").Return(&builds.BuildLog{}, builds.ErrNotFound)

	req, _ := http.NewRequest("GET", "/builds/abc123/logs/build", nil)
	w := httptest.NewRecorder()

	r := gin.New()
	r.GET("/builds/:id/logs/:type", b.Show)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	buildRepo.AssertExpectations(t)
	logRepo.AssertExpectations(t)
}
