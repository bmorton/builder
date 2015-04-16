package api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bmorton/builder/api/mocks"
	"github.com/bmorton/builder/builds"
	"github.com/bmorton/builder/streams"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type closeRecorder struct {
	*httptest.ResponseRecorder
	closer chan bool
}

func newCloseRecorder() *closeRecorder {
	return &closeRecorder{
		httptest.NewRecorder(),
		make(chan bool, 1),
	}
}

func (c *closeRecorder) close() {
	c.closer <- true
}

func (c *closeRecorder) CloseNotify() <-chan bool {
	return c.closer
}

func streamsResourceWithMocks() (*StreamsResource, *mocks.BuildRepository, *mocks.StreamRepository) {
	buildRepo := new(mocks.BuildRepository)
	streamRepo := new(mocks.StreamRepository)
	b := NewStreamsResource(buildRepo, streamRepo)
	return b, buildRepo, streamRepo
}

func TestStreamsShowPlainText(t *testing.T) {
	b, buildRepo, streamRepo := streamsResourceWithMocks()
	buildRepo.On("Find", "abc123").Return(&builds.Build{ID: "abc123"}, nil)
	s := streams.NewBuildStream("abc123")
	streamRepo.On("Find", "abc123").Return(s, nil)

	req, _ := http.NewRequest("GET", "/builds/abc123/streams/build", nil)
	w := newCloseRecorder()
	w.close()

	r := gin.New()
	r.GET("/builds/:id/streams/:type", b.Show)

	s.BuildOutput.Write([]byte("logged!"))
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "text/plain", w.Header().Get("Content-Type"))
	assert.Equal(t, w.Body.String(), "logged!")
	buildRepo.AssertExpectations(t)
	streamRepo.AssertExpectations(t)
}

func TestStreamsShowEventStream(t *testing.T) {
	b, buildRepo, streamRepo := streamsResourceWithMocks()
	buildRepo.On("Find", "abc123").Return(&builds.Build{ID: "abc123"}, nil)
	s := streams.NewBuildStream("abc123")
	streamRepo.On("Find", "abc123").Return(s, nil)

	req, _ := http.NewRequest("GET", "/builds/abc123/streams/build", nil)
	req.Header.Set("Accept", "text/event-stream")
	w := newCloseRecorder()
	w.close()

	r := gin.New()
	r.GET("/builds/:id/streams/:type", b.Show)

	s.BuildOutput.Write([]byte("logged!"))
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "text/event-stream", w.Header().Get("Content-Type"))
	assert.Equal(t, w.Body.String(), "data: logged!\n\n")
	buildRepo.AssertExpectations(t)
	streamRepo.AssertExpectations(t)
}

func TestStreamsShowNotFound(t *testing.T) {
	b, buildRepo, streamRepo := streamsResourceWithMocks()
	buildRepo.On("Find", "abc123").Return(&builds.Build{ID: "abc123"}, nil)
	s := streams.NewBuildStream("abc123")
	streamRepo.On("Find", "abc123").Return(s, streams.ErrNotFound)

	req, _ := http.NewRequest("GET", "/builds/abc123/streams/build", nil)
	w := httptest.NewRecorder()

	r := gin.New()
	r.GET("/builds/:id/streams/:type", b.Show)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	buildRepo.AssertExpectations(t)
	streamRepo.AssertExpectations(t)
}
