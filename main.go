package main

import (
	_ "expvar"
	"log"
	"net/http"
	_ "net/http/pprof"
	"net/url"

	"github.com/fsouza/go-dockerclient"
	"github.com/rcrowley/go-tigertonic"
)

type RootHandler struct{}

func main() {
	mux := tigertonic.NewTrieServeMux()
	mux.HandleNamespace("", http.DefaultServeMux)

	client, err := docker.NewTLSClient("tcp://192.168.59.103:2376", "/Users/bmorton/.boot2docker/certs/boot2docker-vm/cert.pem", "/Users/bmorton/.boot2docker/certs/boot2docker-vm/key.pem", "/Users/bmorton/.boot2docker/certs/boot2docker-vm/ca.pem")
	if err != nil {
		log.Fatal(err)
	}
	githubHandler := &GithubHandler{dockerClient: client}
	mux.Handle("POST", "/webhooks/github", tigertonic.Marshaled(githubHandler.Webhook))

	server := tigertonic.NewServer(":3000", tigertonic.ApacheLogged(mux))
	server.ListenAndServe()
}

func (r *RootHandler) Index(u *url.URL, h http.Header, req interface{}) (int, http.Header, interface{}, error) {
	return http.StatusOK, nil, nil, nil
}
