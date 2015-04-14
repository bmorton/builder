package builds

import (
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
)

type mockBuilder struct {
	calledBuild bool
	calledPush  bool
}

func (m *mockBuilder) BuildImage(b *Build) error {
	m.calledBuild = true
	return nil
}

func (m *mockBuilder) PushImage(b *Build) error {
	m.calledPush = true
	return nil
}

func TestQueueAdd(t *testing.T) {
	repo := repository()
	builder := &mockBuilder{}
	q := NewQueue(repo, builder)

	build := New("deployster", "https://github.com/bmorton/deployster", "abc123", "refs/heads/master")
	q.Add(build)

	assert.Equal(t, build, <-q.queue)
}

func TestQueueSingleRun(t *testing.T) {
	repo := repository()
	builder := &mockBuilder{}
	q := NewQueue(repo, builder)

	build := New("deployster", "https://github.com/bmorton/deployster", "abc123", "refs/heads/master")
	build.ID = "test"
	q.PerformTask(build)

	assert.True(t, builder.calledBuild, builder.calledPush)
}
