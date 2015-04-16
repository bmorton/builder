package builds

import (
	"log"

	"github.com/bmorton/builder/streams"
)

type Builder interface {
	BuildImage(*Build, *streams.Output) error
	PushImage(*Build, *streams.Output) error
}

type BuildSaver interface {
	Save(*Build)
}

type StreamCreateDestroyer interface {
	Create(*streams.BuildStream)
	Destroy(string)
}

type LogCreator interface {
	CreateFromOutput(*streams.BuildStream) (*BuildLog, *BuildLog)
}

type Queue struct {
	queue   chan *Build
	builds  BuildSaver
	streams StreamCreateDestroyer
	logs    LogCreator
	builder Builder
}

func NewQueue(buildRepo BuildSaver, streamRepo StreamCreateDestroyer, logRepo LogCreator, builder Builder) *Queue {
	return &Queue{
		queue:   make(chan *Build, 100),
		builds:  buildRepo,
		streams: streamRepo,
		logs:    logRepo,
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
		stream.Close()
		q.logs.CreateFromOutput(stream)
		q.streams.Destroy(stream.BuildID)
		return
	}

	log.Printf("[%s] Pushing image...\n", build.ID)
	build.State = Pushing
	q.builds.Save(build)
	q.builder.PushImage(build, stream.PushOutput)
	build.State = Complete
	q.builds.Save(build)
	log.Printf("[%s] Build complete!", build.ID)

	stream.Close()
	q.logs.CreateFromOutput(stream)
	q.streams.Destroy(stream.BuildID)
}
