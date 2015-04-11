package builds

import (
	"net/url"
	"path"

	"github.com/bmorton/builder/streams"
)

type Build struct {
	ID             string          `json:"id"`
	RepositoryName string          `json:"repository_name"`
	CloneURL       string          `json:"clone_url"`
	CommitID       string          `json:"commit_id"`
	GitRef         string          `json:"git_ref"`
	ImageTag       string          `json:"image_tag"`
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

func (b *Build) SetDefaultName() {
	parsed, err := url.Parse(b.CloneURL)
	if err != nil {
		return
	}
	b.RepositoryName = path.Base(parsed.Path)
}
