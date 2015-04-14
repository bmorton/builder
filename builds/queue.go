package builds

import (
	"log"

	"github.com/bmorton/builder/streams"
)

type Builder interface {
	BuildImage(*Build) error
	PushImage(*Build) error
}

type Queue struct {
	queue   chan *Build
	builds  *Repository
	builder Builder
}

func NewQueue(repo *Repository, builder Builder) *Queue {
	return &Queue{
		queue:   make(chan *Build, 100),
		builds:  repo,
		builder: builder,
	}
}

func (q *Queue) Add(build *Build) string {
	q.queue <- build
	return build.ID
}

func (q *Queue) Run() {
	for {
		log.Println("Waiting for builds...")
		q.PerformTask(<-q.queue)
	}
}

func (q *Queue) PerformTask(build *Build) {
	log.Printf("[%s] Starting job...\n", build.ID)
	build.BuildStream = streams.NewOutput()
	build.PushStream = streams.NewOutput()
	build.State = Building
	q.builds.Save(build)

	log.Printf("[%s] Building image...\n", build.ID)
	err := q.builder.BuildImage(build)
	if err != nil {
		log.Println(err)
		build.BuildStream.Write([]byte(err.Error()))
		log.Printf("[%s] Build failed!", build.ID)
		build.State = Failed
		q.builds.Save(build)
		q.builds.PersistStreams(build.ID)
		build.BuildStream.Close()
		build.PushStream.Close()
		return
	}

	log.Printf("[%s] Pushing image...\n", build.ID)
	build.State = Pushing
	q.builds.Save(build)
	q.builder.PushImage(build)
	build.State = Complete
	q.builds.Save(build)
	q.builds.PersistStreams(build.ID)
	log.Printf("[%s] Build complete!", build.ID)

	build.BuildStream.Close()
	build.PushStream.Close()
}
