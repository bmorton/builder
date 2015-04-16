package builds

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDefaultName(t *testing.T) {
	b := New("", "https://github.com/bmorton/deployster", "")
	b.SetDefaultName()
	assert.Equal(t, "deployster", b.RepositoryName)
}

func TestStateMarshaling(t *testing.T) {
	s := Waiting
	actual, err := s.MarshalText()
	assert.Nil(t, err)
	assert.Equal(t, []byte("Waiting"), actual)
}
