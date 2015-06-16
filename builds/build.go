package builds

import (
	"net/url"
	"path"
	"time"
)

type Build struct {
	ID             string    `json:"id"`
	RepositoryName string    `json:"repository_name"`
	CloneURL       string    `json:"clone_url"`
	ProjectID      string    `json:"project_id" sql:"index"`
	CommitID       string    `json:"commit_id"`
	ImageTag       string    `json:"image_tag"`
	State          State     `json:"state"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

func New(name, cloneURL, commitID string) *Build {
	return &Build{
		RepositoryName: name,
		CloneURL:       cloneURL,
		CommitID:       commitID,
		State:          Waiting,
	}
}

func (b *Build) SetDefaultName() {
	parsed, err := url.Parse(b.CloneURL)
	if err != nil {
		return
	}
	b.RepositoryName = path.Base(parsed.Path)
}

type State int

const (
	Waiting State = iota
	Building
	Pushing
	Complete
	Failed
)

func (s State) MarshalText() ([]byte, error) {
	return []byte(s.String()), nil
}
