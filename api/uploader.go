package api

import (
	"github.com/francescofrontera/ks-job-upload/dockerutils"
	"github.com/go-chi/chi"
	"io"
	"log"
	"net/http"
	"os"
)

var dockerClient = dockerutils.InitClient()

func UploadHandler(w http.ResponseWriter, r *http.Request) {
	file, handler, err := r.FormFile("uploadFile")
	if err != nil {
		panic(err) //dont do this
	}
	defer file.Close()

	f, err := os.OpenFile("/Users/francescofrontera/go/src/github.com/francescofrontera/ks-job-upload/docker/jars/"+handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	io.Copy(f, file)
}

func BuilderHandler(w http.ResponseWriter, r *http.Request) {
	// Extract json for request {time: Long, jarName: String}
	result := dockerClient.BuildImage("vertx-start-project-1.0-SNAPSHOT-fat.jar")
	io.WriteString(w, result)
}

func RunHandler(w http.ResponseWriter, r *http.Request) {
	containerId := dockerClient.RunContainer("internal_image/vertx-start-project-1.0-snapshot-fat")
	io.WriteString(w, containerId)
}

func UploaderRoute() *chi.Mux  {
	uploadRoute := chi.NewRouter()
	uploadRoute.Post("/upload", UploadHandler)
	uploadRoute.Post("/build", BuilderHandler)
	uploadRoute.Post("/run", RunHandler)
	return uploadRoute
}
