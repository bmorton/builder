package main

import (
	"fmt"
	"log"
	"os"

	"code.google.com/p/go-uuid/uuid"
	"github.com/bmorton/builder/streams"
	"github.com/docker/docker/pkg/archive"
	"github.com/fsouza/go-dockerclient"
	"gopkg.in/libgit2/git2go.v22"
)

type Build struct {
	JobID          string
	RepositoryName string
	CloneURL       string
	CommitID       string
	GitRef         string
}

type BuildQueue struct {
	queue        chan *Build
	dockerClient *docker.Client
	jobs         *JobRepository
}

func NewBuildQueue(repo *JobRepository, dockerClient *docker.Client) *BuildQueue {
	return &BuildQueue{
		queue:        make(chan *Build, 100),
		dockerClient: dockerClient,
		jobs:         repo,
	}
}

func (bq *BuildQueue) Add(build *Build) string {
	build.JobID = uuid.New()
	bq.queue <- build
	return build.JobID
}

func (bq *BuildQueue) Run() {
	for {
		log.Println("Waiting for builds...")
		build := <-bq.queue

		log.Printf("[%s] Starting job...\n", build.JobID)
		writer := streams.NewOutput()
		bq.jobs.Save(build.JobID, writer)

		log.Printf("[%s] Building image...\n", build.JobID)
		bq.buildImage(build)
		log.Printf("[%s] Pushing image...\n", build.JobID)
		bq.pushImage(build)
		log.Printf("[%s] Build complete!", build.JobID)

		writer.Close()
		bq.jobs.Destroy(build.JobID)
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
	stream, _ := bq.jobs.Find(build.JobID)
	err = bq.dockerClient.BuildImage(docker.BuildImageOptions{
		Dockerfile:   "Dockerfile",
		Name:         fmt.Sprintf("%s/%s:%s", "192.168.59.103:5000", build.RepositoryName, build.CommitID[:7]),
		OutputStream: stream,
		InputStream:  context,
	})
	handleError(err)
}

func (bq *BuildQueue) pushImage(build *Build) {
	stream, _ := bq.jobs.Find(build.JobID)
	err := bq.dockerClient.PushImage(docker.PushImageOptions{
		Name:         fmt.Sprintf("%s/%s", "192.168.59.103:5000", build.RepositoryName),
		Tag:          build.CommitID[:7],
		OutputStream: stream,
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
