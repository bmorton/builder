package main

import (
	"fmt"
	"log"
	"os"

	"github.com/docker/docker/pkg/archive"
	"github.com/fsouza/go-dockerclient"
	"gopkg.in/libgit2/git2go.v22"
)

type Build struct {
	RepositoryName string
	CloneURL       string
	CommitID       string
	GitRef         string
}

type BuildQueue struct {
	queue        chan *Build
	dockerClient *docker.Client
}

func (bq *BuildQueue) Add(build *Build) {
	bq.queue <- build
	return
}

func (bq *BuildQueue) Run() {
	for {
		log.Println("Waiting for builds...")
		build := <-bq.queue
		log.Println("Starting build...")
		bq.buildImage(build)
		log.Println("Pushing image...")
		bq.pushImage(build)
		log.Println("Build complete!")
	}
}

func (bq *BuildQueue) buildImage(build *Build) {
	repoPath := fmt.Sprintf("cache/%s", build.RepositoryName)
	repo, err := findOrClone(repoPath, build.CloneURL)
	handleError(err)

	remote, err := repo.LookupRemote("origin")
	handleError(err)
	err = remote.Fetch([]string{build.GitRef}, nil, "")
	handleError(err)
	oid, err := git.NewOid(build.CommitID)
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

	err = bq.dockerClient.BuildImage(docker.BuildImageOptions{
		Dockerfile:   "Dockerfile",
		Name:         fmt.Sprintf("%s/%s:%s", "192.168.59.103:5000", build.RepositoryName, build.CommitID[:6]),
		OutputStream: os.Stdout,
		InputStream:  context,
	})
	handleError(err)
}

func (bq *BuildQueue) pushImage(build *Build) {
	err := bq.dockerClient.PushImage(docker.PushImageOptions{
		Name:         fmt.Sprintf("%s/%s", "192.168.59.103:5000", build.RepositoryName),
		Tag:          build.CommitID[:6],
		OutputStream: os.Stdout,
	}, docker.AuthConfiguration{})
	handleError(err)
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
