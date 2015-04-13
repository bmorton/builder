package streams

import (
	"bytes"
	"fmt"
	"io"
	"testing"
)

type sentinelWriteCloser struct {
	calledWrite bool
	calledClose bool
}

func (w *sentinelWriteCloser) Write(p []byte) (int, error) {
	w.calledWrite = true
	return len(p), nil
}

func (w *sentinelWriteCloser) Close() error {
	w.calledClose = true
	return nil
}

func TestOutputAddClose(t *testing.T) {
	o := NewOutput()
	var s sentinelWriteCloser
	o.Add(&s, make(chan bool, 1))
	if err := o.Close(); err != nil {
		t.Fatal(err)
	}
	// Write data after the output is closed.
	// Write should succeed, but no destination should receive it.
	if _, err := o.Write([]byte("foo bar")); err != nil {
		t.Fatal(err)
	}
	if !s.calledClose {
		t.Fatal("Output.Close() didn't close the destination")
	}
}

func TestOutputAdd(t *testing.T) {
	o := NewOutput()
	b := &bytes.Buffer{}
	o.Add(b, make(chan bool, 1))
	input := "hello, world!"
	if n, err := o.Write([]byte(input)); err != nil {
		t.Fatal(err)
	} else if n != len(input) {
		t.Fatalf("Expected %d, got %d", len(input), n)
	}
	if output := b.String(); output != input {
		t.Fatalf("Received wrong data from Add.\nExpected: '%s'\nGot:     '%s'", input, output)
	}
}

func TestOutputWriteError(t *testing.T) {
	o := NewOutput()
	buf := &bytes.Buffer{}
	o.Add(buf, make(chan bool, 1))
	r, w := io.Pipe()
	input := "Hello there"
	expectedErr := fmt.Errorf("This is an error")
	r.CloseWithError(expectedErr)
	o.Add(w, make(chan bool, 1))
	n, err := o.Write([]byte(input))
	if err != expectedErr {
		t.Fatalf("Output.Write() should return the first error encountered, if any")
	}
	if buf.String() != input {
		t.Fatalf("Output.Write() should attempt write on all destinations, even after encountering an error")
	}
	if n != len(input) {
		t.Fatalf("Output.Write() should return the size of the input if it successfully writes to at least one destination")
	}
}
