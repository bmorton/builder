package builds

import (
	"database/sql"
	"testing"
	"time"

	"github.com/bmorton/builder/streams"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
)

func logRepository() *LogRepository {
	db, _ := sql.Open("sqlite3", ":memory:")
	repo := NewLogRepository("sqlite3", db)
	repo.Migrate()
	return repo
}

func TestLogRepositoryCreateBuildLog(t *testing.T) {
	r := logRepository()
	bs := streams.NewBuildStream("abc123")
	bs.BuildOutput.Write([]byte("logged"))
	buildLog, _ := r.CreateFromOutput(bs)

	timeStub := time.Now()
	buildLog.CreatedAt = timeStub
	buildLog.UpdatedAt = timeStub

	expected := &BuildLog{
		ID:        buildLog.ID,
		Type:      "build",
		BuildID:   "abc123",
		Data:      "logged",
		CreatedAt: timeStub,
		UpdatedAt: timeStub,
	}
	assert.Equal(t, expected, buildLog)
}

func TestLogRepositoryCreatePushLog(t *testing.T) {
	r := logRepository()
	bs := streams.NewBuildStream("abc123")
	bs.PushOutput.Write([]byte("logged"))
	_, pushLog := r.CreateFromOutput(bs)

	timeStub := time.Now()
	pushLog.CreatedAt = timeStub
	pushLog.UpdatedAt = timeStub

	expected := &BuildLog{
		ID:        pushLog.ID,
		Type:      "push",
		BuildID:   "abc123",
		Data:      "logged",
		CreatedAt: timeStub,
		UpdatedAt: timeStub,
	}
	assert.Equal(t, expected, pushLog)
}

func TestLogRepositoryFind(t *testing.T) {
	r := logRepository()
	bs := streams.NewBuildStream("abc123")
	buildLog, _ := r.CreateFromOutput(bs)

	timeStub := time.Now()
	buildLog.CreatedAt = timeStub
	buildLog.UpdatedAt = timeStub

	actual, err := r.FindByBuildID("abc123", "build")
	actual.CreatedAt = timeStub
	actual.UpdatedAt = timeStub

	assert.Nil(t, err)
	assert.Equal(t, buildLog, actual)
}

func TestLogRepositoryFindNotFound(t *testing.T) {
	r := logRepository()
	_, err := r.FindByBuildID("abc123", "build")
	assert.Equal(t, ErrNotFound, err)
}
