package streams

import (
	"bytes"
	"io"
	"sync"
)

type Output struct {
	sync.Mutex
	buffer *bytes.Buffer
	dests  []*Destination
	used   bool
}

type Destination struct {
	writer   io.Writer
	waitChan chan bool
}

// NewOutput returns a new Output object with no destinations attached.
// Writing to an empty Output will cause the written data to be discarded.
func NewOutput() *Output {
	buf := bytes.NewBuffer(nil)
	o := &Output{buffer: buf}
	o.Add(buf, make(chan bool, 1))
	return o
}

// Return true if something was written on this output
func (o *Output) Used() bool {
	o.Lock()
	defer o.Unlock()
	return o.used
}

// Add attaches a new destination to the Output. Any data subsequently written
// to the output will be written to the new destination in addition to all the others.
// This method is thread-safe.
func (o *Output) Add(dst io.Writer, waitChan chan bool) {
	o.Lock()
	defer o.Unlock()
	dest := &Destination{writer: dst, waitChan: waitChan}
	o.dests = append(o.dests, dest)
}

func (o *Output) Remove(dst io.Writer) {
	o.Lock()
	defer o.Unlock()
	for i, d := range o.dests {
		if d.writer == dst {
			o.dests = append(o.dests[:i], o.dests[i+1:]...)
		}
	}
}

// Write writes the same data to all registered destinations.
// This method is thread-safe.
func (o *Output) Write(p []byte) (n int, err error) {
	o.Lock()
	defer o.Unlock()
	o.used = true
	var firstErr error
	for _, dst := range o.dests {
		_, err := dst.writer.Write(p)
		if err != nil && firstErr == nil {
			firstErr = err
		}
	}
	return len(p), firstErr
}

// Close unregisters all destinations and waits for all background
// AddTail and AddString tasks to complete.
// The Close method of each destination is called if it exists.
func (o *Output) Close() error {
	o.Lock()
	defer o.Unlock()
	o.Notify()
	var firstErr error
	for _, dst := range o.dests {
		if closer, ok := dst.writer.(io.Closer); ok {
			err := closer.Close()
			if err != nil && firstErr == nil {
				firstErr = err
			}
		}
	}
	o.dests = nil
	return firstErr
}

func (o *Output) Notify() {
	for _, dst := range o.dests {
		dst.waitChan <- true
	}
}

func (o *Output) Replay(dst io.Writer) {
	o.Lock()
	defer o.Unlock()
	dst.Write(o.buffer.Bytes())
	o.buffer.UnreadByte()
}
