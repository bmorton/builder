package streams

type BuildStream struct {
	BuildID     string
	BuildOutput *Output
	PushOutput  *Output
}

func NewBuildStream(buildID string) *BuildStream {
	return &BuildStream{
		BuildID:     buildID,
		BuildOutput: NewOutput(),
		PushOutput:  NewOutput(),
	}
}

func (b *BuildStream) Close() {
	b.BuildOutput.Close()
	b.PushOutput.Close()
}
