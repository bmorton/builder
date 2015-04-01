package main

import (
	_ "expvar"
	"log"
	"net/http"
	_ "net/http/pprof"

	"github.com/fsouza/go-dockerclient"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	client, err := docker.NewTLSClient("tcp://192.168.59.103:2376", "/Users/bmorton/.boot2docker/certs/boot2docker-vm/cert.pem", "/Users/bmorton/.boot2docker/certs/boot2docker-vm/key.pem", "/Users/bmorton/.boot2docker/certs/boot2docker-vm/ca.pem")
	if err != nil {
		log.Fatal(err)
	}

	repo := NewJobRepository()
	buildQueue := NewBuildQueue(repo, client)

	webhookHandler := &WebhookHandler{queue: buildQueue}
	router.POST("/webhooks/github", webhookHandler.Github)

	buildsResource := &BuildsResource{builds: repo}
	router.GET("/builds/:id", buildsResource.Show)

	go buildQueue.Run()

	log.Println("Starting debug server on :3001")
	go http.ListenAndServe(":3001", http.DefaultServeMux)

	router.Run(":3000")
}
