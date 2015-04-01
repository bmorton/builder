package main

import (
	"net/http"
	"net/url"
)

type GithubHandler struct {
	queue *BuildQueue
}

func (gh *GithubHandler) Webhook(u *url.URL, h http.Header, req *GithubPushEvent) (int, http.Header, *Build, error) {
	build := &Build{
		RepositoryName: req.Repository.Name,
		CloneURL:       req.Repository.CloneURL,
		CommitID:       req.HeadCommit.ID,
		GitRef:         req.Ref,
	}
	gh.queue.Add(build)

	return http.StatusOK, nil, build, nil
}
