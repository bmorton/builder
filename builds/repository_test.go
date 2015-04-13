package builds

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRepositorySaveFind(t *testing.T) {
	expected := New("deployster", "https://github.com/bmorton/deployster", "abc123", "refs/heads/master")
	r := NewRepository()
	r.Save("123", expected)

	actual, ok := r.Find("123")
	assert.True(t, ok)
	assert.Equal(t, expected, actual)
}

func TestRepositoryDestroy(t *testing.T) {
	expected := New("deployster", "https://github.com/bmorton/deployster", "abc123", "refs/heads/master")
	r := NewRepository()
	r.Save("123", expected)

	_, ok := r.Find("123")
	assert.True(t, ok)

	r.Destroy("123")
	_, ok = r.Find("123")
	assert.False(t, ok)
}

func TestRepositoryKeys(t *testing.T) {
	expected := New("deployster", "https://github.com/bmorton/deployster", "abc123", "refs/heads/master")
	r := NewRepository()
	r.Save("123", expected)

	assert.Equal(t, []string{"123"}, r.Keys())
}

func TestRepositoryAll(t *testing.T) {
	expected := New("deployster", "https://github.com/bmorton/deployster", "abc123", "refs/heads/master")
	r := NewRepository()
	r.Save("123", expected)

	assert.Equal(t, []*Build{expected}, r.All())
}
