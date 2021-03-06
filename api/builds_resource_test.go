package api

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/bmorton/builder/api/mocks"
	"github.com/bmorton/builder/builds"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func resourceWithMocks() (*BuildsResource, *mocks.BuildRepository, *mocks.BuildQueue) {
	repo := new(mocks.BuildRepository)
	queue := new(mocks.BuildQueue)
	b := NewBuildsResource(repo, queue)
	return b, repo, queue
}

func TestIndex(t *testing.T) {
	b, repo, _ := resourceWithMocks()
	repo.On("All").Return([]*builds.Build{
		&builds.Build{ID: "abc123"},
	})

	req, _ := http.NewRequest("GET", "/builds", nil)
	w := httptest.NewRecorder()

	r := gin.New()
	r.GET("/builds", b.Index)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "abc123")
	repo.AssertExpectations(t)
}

func TestCreate(t *testing.T) {
	b, repo, queue := resourceWithMocks()
	queue.On("Add", &builds.Build{
		RepositoryName: "deployster",
		CloneURL:       "https://github.com/bmorton/deployster",
	}).Return("abc123")
	repo.On("Create", &builds.Build{
		RepositoryName: "deployster",
		CloneURL:       "https://github.com/bmorton/deployster",
	}).Return()

	payload := `{"clone_url":"https://github.com/bmorton/deployster"}`
	req, _ := http.NewRequest("POST", "/builds", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r := gin.New()
	r.POST("/builds", b.Create)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	queue.AssertExpectations(t)
	repo.AssertExpectations(t)
}

func TestShow(t *testing.T) {
	b, repo, _ := resourceWithMocks()
	repo.On("Find", "abc123").Return(&builds.Build{ID: "abc123"}, nil)

	req, _ := http.NewRequest("GET", "/builds/abc123", nil)
	w := httptest.NewRecorder()

	r := gin.New()
	r.GET("/builds/:id", b.Show)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "abc123")
	repo.AssertExpectations(t)
}

func TestShowNotFound(t *testing.T) {
	b, repo, _ := resourceWithMocks()
	repo.On("Find", "abc123").Return(&builds.Build{}, builds.ErrNotFound)

	req, _ := http.NewRequest("GET", "/builds/abc123", nil)
	w := httptest.NewRecorder()

	r := gin.New()
	r.GET("/builds/:id", b.Show)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	repo.AssertExpectations(t)
}
