package main

import (
	_ "expvar"
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"

	"github.com/bmorton/builder/api"
	"github.com/bmorton/builder/builds"
	"github.com/fsouza/go-dockerclient"
	"github.com/gin-gonic/gin"
	"github.com/namsral/flag"
)

func main() {
	var listen string
	var dockerHost string
	var dockerTLSVerify bool
	var dockerCertPath string
	var debugMode bool
	var debugListen string
	var registryURL string

	flag.StringVar(&listen, "listen", ":3000", "host:port to listen on")
	flag.StringVar(&dockerHost, "docker-host", "", "address of Docker host")
	flag.BoolVar(&dockerTLSVerify, "docker-tls-verify", false, "use TLS client for Docker")
	flag.StringVar(&dockerCertPath, "docker-cert-path", "", "path to the cert.pem, key.pem, and ca.pem for authenticating to Docker")
	flag.BoolVar(&debugMode, "debug", false, "enable /debug endpoints on DEBUG_LISTEN")
	flag.StringVar(&debugListen, "debug-listen", ":3001", "host:port to listen on for debug requests")
	flag.StringVar(&registryURL, "registry-url", "192.168.59.103:5000", "host:port of the registry for pushing images")
	flag.Parse()

	router := gin.Default()
	client := dockerClient(dockerHost, dockerTLSVerify, dockerCertPath)
	repo := builds.NewRepository()
	builder := builds.NewBuilder(registryURL, client)
	buildQueue := builds.NewQueue(repo, builder)

	webhookHandler := api.NewWebhookHandler(buildQueue)
	router.POST("/webhooks/github", webhookHandler.Github)

	buildsResource := api.NewBuildsResource(repo)
	router.GET("/builds", buildsResource.Index)
	router.GET("/builds/:id", buildsResource.Show)

	go buildQueue.Run()

	if debugMode {
		log.Println("Starting debug server on :3001")
		go http.ListenAndServe(debugListen, http.DefaultServeMux)
	}

	router.Run(listen)
}

func dockerClient(host string, tls bool, certPath string) *docker.Client {
	var client *docker.Client
	var err error

	if tls {
		cert := fmt.Sprintf("%s/cert.pem", certPath)
		key := fmt.Sprintf("%s/key.pem", certPath)
		ca := fmt.Sprintf("%s/ca.pem", certPath)
		client, err = docker.NewTLSClient(host, cert, key, ca)
	} else {
		client, err = docker.NewClient(host)
	}

	if err != nil {
		log.Fatal(err)
	}

	return client
}
