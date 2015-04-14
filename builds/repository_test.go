package builds

import (
	"database/sql"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
)

func repository() *Repository {
	db, _ := sql.Open("sqlite3", ":memory:")
	repo := NewRepository("sqlite3", db)
	repo.Migrate()
	return repo
}

func TestRepositoryCreateFind(t *testing.T) {
	expected := New("deployster", "https://github.com/bmorton/deployster", "abc123", "refs/heads/master")
	r := repository()
	r.Create(expected)

	actual := r.Find(expected.ID)
	assert.Equal(t, expected.ID, actual.ID)
}

func TestRepositoryAll(t *testing.T) {
	expected := New("deployster", "https://github.com/bmorton/deployster", "abc123", "refs/heads/master")
	r := repository()
	r.Create(expected)

	assert.Equal(t, expected.ID, r.All()[0].ID)
}
