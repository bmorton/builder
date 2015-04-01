package main

import (
	_ "expvar"
	"log"
	"net/http"
	_ "net/http/pprof"

	"github.com/fsouza/go-dockerclient"
	"github.com/rcrowley/go-tigertonic"
)

func main() {
	mux := tigertonic.NewTrieServeMux()
	mux.HandleNamespace("", http.DefaultServeMux)

	client, err := docker.NewTLSClient("tcp://192.168.59.103:2376", "/Users/bmorton/.boot2docker/certs/boot2docker-vm/cert.pem", "/Users/bmorton/.boot2docker/certs/boot2docker-vm/key.pem", "/Users/bmorton/.boot2docker/certs/boot2docker-vm/ca.pem")
	if err != nil {
		log.Fatal(err)
	}

	repo := NewJobRepository()
	buildQueue := NewBuildQueue(repo, client)
	githubHandler := &GithubHandler{queue: buildQueue}
	mux.Handle("POST", "/webhooks/github", tigertonic.ApacheLogged(tigertonic.Marshaled(githubHandler.Webhook)))
	buildsResource := &BuildsResource{builds: repo}
	mux.Handle("GET", "/builds/{id}", http.HandlerFunc(buildsResource.Show))

	go buildQueue.Run()

	server := tigertonic.NewServer(":3000", mux)
	server.ListenAndServe()
}
