package api

import (
	"bufio"
	"bytes"
	"fmt"
	"io"

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
	if isSkippable(p) {
		return 0, nil
	}

	reader := bytes.NewReader(p)
	scanner := bufio.NewScanner(reader)
	var fixed string
	for scanner.Scan() {
		fixed = fmt.Sprintf("%sdata: %s\n\n", fixed, scanner.Text())
	}
	return sw.writer.Write([]byte(fixed))
}

func NewSSEWriter(w io.Writer) SSEWriter {
	return SSEWriter{writer: flushwriter.New(w)}
}

func isSkippable(output []byte) bool {
	s := string(output)
	if s == "Buffering to disk\n" || s == "Pushing\n" || s == "Downloading\n" {
		return true
	}
	return false
}
