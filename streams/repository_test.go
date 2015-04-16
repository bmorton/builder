package streams

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRepositoryCreate(t *testing.T) {
	r := NewRepository()
	expected := &BuildStream{BuildID: "abc123", BuildOutput: NewOutput(), PushOutput: NewOutput()}
	r.Create(expected)

	actual, err := r.Find("abc123")
	assert.Nil(t, err)
	assert.Equal(t, expected, actual)
}

func TestRepositoryFindNotFound(t *testing.T) {
	r := NewRepository()
	_, err := r.Find("abc123")
	assert.Equal(t, ErrNotFound, err)
}
