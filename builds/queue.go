package builds

import (
	"log"

	"code.google.com/p/go-uuid/uuid"
	"github.com/bmorton/builder/streams"
)

type Queue struct {
	queue   chan *Build
	builds  *Repository
	builder *Builder
}

func NewQueue(repo *Repository, builder *Builder) *Queue {
	return &Queue{
		queue:   make(chan *Build, 100),
		builds:  repo,
		builder: builder,
	}
}

func (q *Queue) Add(build *Build) string {
	build.ID = uuid.New()
	q.queue <- build
	return build.ID
}

func (q *Queue) Run() {
	for {
		log.Println("Waiting for builds...")
		build := <-q.queue

		log.Printf("[%s] Starting job...\n", build.ID)
		build.OutputStream = streams.NewOutput()
		build.State = Building

		log.Printf("[%s] Building image...\n", build.ID)
		err := q.builder.BuildImage(build)
		if err != nil {
			log.Println(err)
			build.OutputStream.Write([]byte(err.Error()))
			log.Printf("[%s] Build failed!", build.ID)
			build.OutputStream.Close()
			build.State = Failed
			continue
		}

		log.Printf("[%s] Pushing image...\n", build.ID)
		build.State = Pushing
		q.builder.PushImage(build)
		build.State = Complete
		log.Printf("[%s] Build complete!", build.ID)

		build.OutputStream.Close()
		q.builds.Destroy(build.ID)
	}
}
