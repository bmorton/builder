package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/docker/docker/pkg/archive"
	"github.com/fsouza/go-dockerclient"
	"gopkg.in/libgit2/git2go.v22"
)

type GithubHandler struct {
	dockerClient *docker.Client
}

func (gh *GithubHandler) Webhook(u *url.URL, h http.Header, req *GithubPushEvent) (int, http.Header, interface{}, error) {
	repoPath := fmt.Sprintf("cache/%s", req.Repository.Name)
	repo, err := findOrClone(repoPath, req.Repository.CloneURL)
	handleError(err)

	remote, err := repo.LookupRemote("origin")
	handleError(err)
	err = remote.Fetch([]string{req.Ref}, nil, "")
	handleError(err)
	oid, err := git.NewOid(req.HeadCommit.ID)
	handleError(err)
	commit, err := repo.LookupCommit(oid)
	handleError(err)
	tree, err := commit.Tree()
	handleError(err)
	err = repo.CheckoutTree(tree, &git.CheckoutOpts{Strategy: git.CheckoutForce})
	handleError(err)
	err = repo.SetHeadDetached(oid, nil, "")
	handleError(err)

	options := &archive.TarOptions{
		Compression:     archive.Uncompressed,
		ExcludePatterns: []string{".git"},
		IncludeFiles:    []string{"."},
	}
	context, err := archive.TarWithOptions(repoPath, options)

	err = gh.dockerClient.BuildImage(docker.BuildImageOptions{
		Dockerfile:   "Dockerfile",
		Name:         fmt.Sprintf("%s/%s:%s", "192.168.59.103:5000", req.Repository.Name, req.HeadCommit.ID[:6]),
		OutputStream: os.Stdout,
		InputStream:  context,
	})
	handleError(err)

	err = gh.dockerClient.PushImage(docker.PushImageOptions{
		Name:         fmt.Sprintf("%s/%s", "192.168.59.103:5000", req.Repository.Name),
		Tag:          req.HeadCommit.ID[:6],
		OutputStream: os.Stdout,
	}, docker.AuthConfiguration{})
	handleError(err)

	return http.StatusOK, nil, nil, nil
}

func findOrClone(path string, cloneURL string) (*git.Repository, error) {
	var repo *git.Repository
	var err error

	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			repo, err = git.Clone(cloneURL, path, &git.CloneOptions{})
		} else {
			return &git.Repository{}, err
		}
	} else {
		repo, err = git.OpenRepository(path)
	}

	return repo, err
}

func handleError(err error) {
	if err != nil {
		log.Println(err)
	}
}
