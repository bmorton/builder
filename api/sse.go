package api

import (
	"fmt"
	"io"
	"strings"

	"github.com/bmorton/flushwriter"
)

// SSEWriter takes a normal io.Writer and pads it with the proper format so that
// any normal stream of data can be sent, albeit unstructured, to the client as
// SSE.
type SSEWriter struct {
	writer flushwriter.FlushWriter
}

// Write prepends "data: " to every line of the byte slice, leaving an empty
// line inbetween each data event.
func (sw SSEWriter) Write(p []byte) (int, error) {
	fixed := strings.Replace(string(p), "\n", "\n\ndata: ", -1)
	event := fmt.Sprintf("data: %s\n\n", fixed)
	return sw.writer.Write([]byte(event))
}

func NewSSEWriter(w io.Writer) SSEWriter {
	return SSEWriter{writer: flushwriter.New(w)}
}
