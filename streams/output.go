// This package is mostly ripped out of docker/engine/streams.go:
// https://github.com/docker/docker/blob/d045b9776b5dc16e12b3d7c7558a24cdc5d1aba7/engine/streams.go
// At the time, it couldn't stand on its own, so it was added here.
// Additionally, for our purposes we only care about the `Output` part of that
// file.
package streams

import (
	"bytes"
	"io"
	"sync"
)

// Output is responsible for distributing a single stream of writes to multiple
// destinations.  It also contains the buffer of all bytes written from the
// beginning so that new destinations can replay what they haven't yet seen.
type Output struct {
	sync.Mutex
	buffer *bytes.Buffer
	dests  []*Destination
	used   bool
}

// Destination is a wrapper for a single output stream along with a channel that
// is used to signal when the output stream has been closed and no more data
// will be written.
type Destination struct {
	writer   io.Writer
	waitChan chan bool
}

// NewOutput returns a new Output object with no destinations attached.
func NewOutput() *Output {
	buf := bytes.NewBuffer(nil)
	o := &Output{buffer: buf}
	o.Add(buf, make(chan bool, 1))
	return o
}

// Used return true if something was written on this output.
func (o *Output) Used() bool {
	o.Lock()
	defer o.Unlock()
	return o.used
}

// Add attaches a new destination to the Output. Any data subsequently written
// to the output will be written to the new destination in addition to all the
// others.  This method is thread-safe.
func (o *Output) Add(dst io.Writer, waitChan chan bool) {
	o.Lock()
	defer o.Unlock()
	dest := &Destination{writer: dst, waitChan: waitChan}
	o.dests = append(o.dests, dest)
}

// Remove allows a destination to be removed from the output stream.
// This method is thread-safe.
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

// Close notifies and unregisters all destinations. The Close method of each
// destination is called if it exists.
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

// Notify informs all destinations that no more data will be written to the
// stream.
func (o *Output) Notify() {
	for _, dst := range o.dests {
		dst.waitChan <- true
	}
}

// Replay quickly plays back any data that has been written to the stream since
// the beginning.  It blocks all other Output actions until it completes.
func (o *Output) Replay(dst io.Writer) {
	o.Lock()
	defer o.Unlock()
	dst.Write(o.buffer.Bytes())
	o.buffer.UnreadByte()
}
