package builds

import "github.com/bmorton/builder/streams"

type Build struct {
	ID             string          `json:"id"`
	RepositoryName string          `json:"repository_name"`
	CloneURL       string          `json:"clone_url"`
	CommitID       string          `json:"commit_id"`
	GitRef         string          `json:"git_ref"`
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
