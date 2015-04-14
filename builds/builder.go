package builds

import (
	"fmt"
	"log"
	"os"

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

func (b *DockerBuilder) BuildImage(build *Build) error {
	repoPath := fmt.Sprintf("%s/%s", b.cachePath, build.RepositoryName)
	repo, err := findOrClone(repoPath, build.CloneURL)
	if err != nil {
		return err
	}

	remote, err := repo.LookupRemote("origin")
	if err != nil {
		return err
	}

	var refSpecs []string
	if build.GitRef == "" {
		refSpecs, err = remote.FetchRefspecs()
		handleError(err)
	} else {
		refSpecs = []string{build.GitRef}
	}

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
	err = b.dockerClient.BuildImage(docker.BuildImageOptions{
		Dockerfile:   "Dockerfile",
		Name:         fmt.Sprintf("%s/%s:%s", b.registryURL, build.RepositoryName, build.ImageTag),
		OutputStream: build.BuildStream,
		InputStream:  context,
	})

	return err
}

func (b *DockerBuilder) PushImage(build *Build) error {
	return b.dockerClient.PushImage(docker.PushImageOptions{
		Name:         fmt.Sprintf("%s/%s", b.registryURL, build.RepositoryName),
		Tag:          build.ImageTag,
		OutputStream: build.PushStream,
	}, docker.AuthConfiguration{})
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
