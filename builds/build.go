package builds

import (
	"net/url"
	"path"
	"time"

	"github.com/bmorton/builder/streams"
)

type Build struct {
	ID             string          `json:"id"`
	RepositoryName string          `json:"repository_name"`
	CloneURL       string          `json:"clone_url"`
	CommitID       string          `json:"commit_id"`
	GitRef         string          `json:"git_ref"`
	ImageTag       string          `json:"image_tag"`
	State          State           `json:"state"`
	CreatedAt      time.Time       `json:"created_at"`
	UpdatedAt      time.Time       `json:"updated_at"`
	BuildStream    *streams.Output `json:"-", sql:"-"`
	PushStream     *streams.Output `json:"-", sql:"-"`
}

type BuildLog struct {
	ID        string    `json:"id"`
	Data      string    `json:"data" sql:"type:blob"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type PushLog struct {
	ID        string    `json:"id"`
	Data      string    `json:"data" sql:"type:blob"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func New(name, cloneURL, commitID, gitRef string) *Build {
	return &Build{
		RepositoryName: name,
		CloneURL:       cloneURL,
		CommitID:       commitID,
		GitRef:         gitRef,
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

func (b *Build) IsFinished() bool {
	switch b.State {
	case Complete:
		return true
	case Failed:
		return true
	default:
		return false
	}
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
