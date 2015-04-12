package api

import (
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWrites(t *testing.T) {
	w := httptest.NewRecorder()
	subject := NewSSEWriter(w)
	subject.Write([]byte("I wanna be hugged by you, by you"))

	assert.Equal(t, "data: I wanna be hugged by you, by you\n\n", w.Body.String())
}

func TestWritesMultipleLines(t *testing.T) {
	w := httptest.NewRecorder()
	subject := NewSSEWriter(w)
	subject.Write([]byte("Life's an endless party, not a punchcard\nI don't understand some people's drive\nLets just fuckin' drink and be alive\nNot just survive"))

	assert.Equal(t, `data: Life's an endless party, not a punchcard

data: I don't understand some people's drive

data: Lets just fuckin' drink and be alive

data: Not just survive

`, w.Body.String())
}

func TestWriteFlushesAutomatically(t *testing.T) {
	w := httptest.NewRecorder()
	subject := NewSSEWriter(w)
	subject.Write([]byte("test"))

	assert.True(t, w.Flushed)
}
