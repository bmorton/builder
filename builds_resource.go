package main

import (
	"net/http"

	"github.com/bmorton/flushwriter"
)

type BuildsResource struct {
	builds *JobRepository
}

func (br *BuildsResource) Show(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")

	jobID := r.URL.Query().Get("id")

	writer, ok := br.builds.Find(jobID)
	if !ok || writer == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	waitChan := make(chan bool, 1)
	notify := w.(http.CloseNotifier).CloseNotify()

	go func() {
		<-notify
		waitChan <- true
	}()

	w.WriteHeader(http.StatusOK)
	fw := flushwriter.New(w)
	writer.Replay(&fw)
	writer.Add(&fw, waitChan)
	<-waitChan
	writer.Remove(&fw)
}
