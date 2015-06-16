package projects

import (
	"fmt"
	"net/url"
	"path"
	"path/filepath"
	"strings"
	"time"
)

type Project struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	FullName  string    `json:"full_name" sql:"unique_index"`
	CloneURL  string    `json:"clone_url" sql:"unique_index"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func New(cloneURL string) *Project {
	p := &Project{
		CloneURL: cloneURL,
	}
	p.SetDefaultName()
	return p
}

func (p *Project) SetDefaultName() {
	parsed, err := url.Parse(p.CloneURL)
	if err != nil {
		return
	}
	basename := path.Base(parsed.Path)
	p.Name = strings.TrimSuffix(strings.TrimLeft(basename, "/"), filepath.Ext(basename))

	if parsed.Host == "github.com" {
		p.FullName = strings.TrimSuffix(strings.TrimLeft(parsed.Path, "/"), filepath.Ext(basename))
	} else {
		p.FullName = fmt.Sprintf("%s/%s", parsed.Host, strings.TrimSuffix(strings.TrimLeft(parsed.Path, "/"), filepath.Ext(basename)))
	}
}
