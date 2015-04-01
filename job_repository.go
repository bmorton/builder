package main

import "github.com/bmorton/builder/streams"

type JobRepository struct {
	jobs map[string]*streams.Output
}

func NewJobRepository() *JobRepository {
	return &JobRepository{
		jobs: make(map[string]*streams.Output),
	}
}

func (jr *JobRepository) Find(key string) (*streams.Output, bool) {
	output, ok := jr.jobs[key]
	return output, ok
}

func (jr *JobRepository) Save(key string, output *streams.Output) {
	jr.jobs[key] = output
}

func (jr *JobRepository) Destroy(key string) {
	jr.jobs[key] = nil
}

func (jr *JobRepository) Keys() []string {
	keys := make([]string, 0, len(jr.jobs))
	for k := range jr.jobs {
		keys = append(keys, k)
	}
	return keys
}
