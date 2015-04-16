package builds

import (
	"testing"

	"github.com/bmorton/builder/builds/mocks"
	"github.com/bmorton/builder/streams"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockBuilder struct {
	calledBuild bool
	calledPush  bool
}

func (m *mockBuilder) BuildImage(b *Build, o *streams.Output) error {
	m.calledBuild = true
	return nil
}

func (m *mockBuilder) PushImage(b *Build, o *streams.Output) error {
	m.calledPush = true
	return nil
}

func queue() (BuildSaver, *mocks.StreamCreateDestroyer, LogCreator, *mockBuilder, *Queue) {
	buildRepo := repository()
	streamRepo := new(mocks.StreamCreateDestroyer)
	logRepo := logRepository()
	builder := &mockBuilder{}
	q := NewQueue(buildRepo, streamRepo, logRepo, builder)

	return buildRepo, streamRepo, logRepo, builder, q
}

func TestQueueAdd(t *testing.T) {
	_, _, _, _, q := queue()
	build := New("deployster", "https://github.com/bmorton/deployster", "abc123", "refs/heads/master")
	q.Add(build)

	assert.Equal(t, build, <-q.queue)
}

func TestQueueSingleRun(t *testing.T) {
	_, streamRepo, _, builder, q := queue()
	streamRepo.On("Create", mock.Anything).Return(nil)
	streamRepo.On("Destroy", "test").Return(nil)

	build := New("deployster", "https://github.com/bmorton/deployster", "abc123", "refs/heads/master")
	build.ID = "test"
	q.PerformTask(build)

	assert.True(t, builder.calledBuild, builder.calledPush)
}
