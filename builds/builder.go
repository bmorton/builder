package builds

import (
	"fmt"
	"log"
	"os"

	"github.com/bmorton/builder/streams"
	"github.com/docker/docker/pkg/archive"
	"github.com/fsouza/go-dockerclient"
	"gopkg.in/libgit2/git2go.v22"
)

type DockerBuilder struct {
	registryURL  string
	dockerClient *docker.Client
	cachePath    string
}

func NewBuilder(registryURL string, dockerClient *docker.Client, cachePath string) *DockerBuilder {
	return &DockerBuilder{registryURL: registryURL, dockerClient: dockerClient, cachePath: cachePath}
}

func (b *DockerBuilder) BuildImage(build *Build, stream *streams.Output) error {
	repoPath := fmt.Sprintf("%s/%s", b.cachePath, build.RepositoryName)
	repo, err := findOrClone(repoPath, build.CloneURL)
	if err != nil {
		return err
	}

	remote, err := repo.LookupRemote("origin")
	if err != nil {
		return err
	}

	refSpecs, err := remote.FetchRefspecs()
	handleError(err)

	err = remote.Fetch(refSpecs, nil, "")
	if err != nil {
		return err
	}

	oid, err := git.NewOid(build.CommitID)
	if err != nil {
		return err
	}

	commit, err := repo.LookupCommit(oid)
	handleError(err)
	tree, err := commit.Tree()
	handleError(err)
	err = repo.CheckoutTree(tree, &git.CheckoutOpts{Strategy: git.CheckoutForce})
	handleError(err)
	err = repo.SetHeadDetached(oid, nil, "")
	handleError(err)

	build.ImageTag = build.CommitID[:7]
	options := &archive.TarOptions{
		Compression:     archive.Uncompressed,
		ExcludePatterns: []string{".git"},
		IncludeFiles:    []string{"."},
	}
	context, err := archive.TarWithOptions(repoPath, options)
	name := fmt.Sprintf("%s/%s:%s", b.registryURL, build.RepositoryName, build.ImageTag)
	err = b.dockerClient.BuildImage(docker.BuildImageOptions{
		Dockerfile:   "Dockerfile",
		Name:         name,
		OutputStream: stream,
		InputStream:  context,
	})
	if err != nil {
		return err
	}

	err = b.dockerClient.TagImage(name, docker.TagImageOptions{
		Repo: fmt.Sprintf("%s/%s", b.registryURL, build.RepositoryName),
		Tag:  "latest",
	})

	return err
}

func (b *DockerBuilder) PushImage(build *Build, stream *streams.Output) error {
	name := fmt.Sprintf("%s/%s", b.registryURL, build.RepositoryName)

	stream.Write([]byte(fmt.Sprintf("Pushing %s:%s...", name, build.ImageTag)))
	err := b.dockerClient.PushImage(docker.PushImageOptions{
		Name:         name,
		Tag:          build.ImageTag,
		OutputStream: stream,
	}, docker.AuthConfiguration{})

	if err != nil {
		return err
	}

	stream.Write([]byte(fmt.Sprintf("Pushing %s:latest...", name)))
	err = b.dockerClient.PushImage(docker.PushImageOptions{
		Name:         name,
		Tag:          "latest",
		OutputStream: stream,
	}, docker.AuthConfiguration{})

	return err
}

func findOrClone(path string, cloneURL string) (*git.Repository, error) {
	var repo *git.Repository
	var err error

	if _, err = os.Stat(path); err != nil {
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
