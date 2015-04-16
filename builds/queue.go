package builds

import (
	"log"

	"github.com/bmorton/builder/streams"
)

type Builder interface {
	BuildImage(*Build, *streams.Output) error
	PushImage(*Build, *streams.Output) error
}

type Queue struct {
	queue   chan *Build
	builds  *Repository
	streams *streams.Repository
	builder Builder
}

func NewQueue(buildRepo *Repository, streamRepo *streams.Repository, builder Builder) *Queue {
	return &Queue{
		queue:   make(chan *Build, 100),
		builds:  buildRepo,
		streams: streamRepo,
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
	stream := streams.NewBuildStream(build.ID)
	q.streams.Create(stream)
	build.State = Building
	q.builds.Save(build)

	log.Printf("[%s] Building image...\n", build.ID)
	err := q.builder.BuildImage(build, stream.BuildOutput)
	if err != nil {
		log.Println(err)
		stream.BuildOutput.Write([]byte(err.Error()))
		log.Printf("[%s] Build failed!", build.ID)
		build.State = Failed
		q.builds.Save(build)
		stream.BuildOutput.Close()
		stream.PushOutput.Close()
		return
	}

	log.Printf("[%s] Pushing image...\n", build.ID)
	build.State = Pushing
	q.builds.Save(build)
	q.builder.PushImage(build, stream.PushOutput)
	build.State = Complete
	q.builds.Save(build)
	log.Printf("[%s] Build complete!", build.ID)

	stream.BuildOutput.Close()
	stream.PushOutput.Close()
}
