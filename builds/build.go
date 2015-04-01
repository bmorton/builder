package builds

import "github.com/bmorton/builder/streams"

type Build struct {
	ID             string
	RepositoryName string
	CloneURL       string
	CommitID       string
	GitRef         string
	OutputStream   *streams.Output `json:"-"`
}

func New(name, cloneURL, commitID, gitRef string) *Build {
	return &Build{
		RepositoryName: name,
		CloneURL:       cloneURL,
		CommitID:       commitID,
		GitRef:         gitRef,
	}
}
