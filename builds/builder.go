package builds

import (
	"fmt"
	"log"
	"os"

	"github.com/docker/docker/pkg/archive"
	"github.com/fsouza/go-dockerclient"
	"gopkg.in/libgit2/git2go.v22"
)

type Builder struct {
	registryURL  string
	dockerClient *docker.Client
	cachePath    string
}

func NewBuilder(registryURL string, dockerClient *docker.Client, cachePath string) *Builder {
	return &Builder{registryURL: registryURL, dockerClient: dockerClient, cachePath: cachePath}
}

func (b *Builder) BuildImage(build *Build) error {
	repoPath := fmt.Sprintf("%s/%s", b.cachePath, build.RepositoryName)
	repo, err := findOrClone(repoPath, build.CloneURL)
	if err != nil {
		return err
	}

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
	err = b.dockerClient.BuildImage(docker.BuildImageOptions{
		Dockerfile:   "Dockerfile",
		Name:         fmt.Sprintf("%s/%s:%s", b.registryURL, build.RepositoryName, build.CommitID[:7]),
		OutputStream: build.OutputStream,
		InputStream:  context,
	})
	handleError(err)

	return nil
}

func (b *Builder) PushImage(build *Build) {
	err := b.dockerClient.PushImage(docker.PushImageOptions{
		Name:         fmt.Sprintf("%s/%s", b.registryURL, build.RepositoryName),
		Tag:          build.CommitID[:7],
		OutputStream: build.OutputStream,
	}, docker.AuthConfiguration{})
	handleError(err)
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
