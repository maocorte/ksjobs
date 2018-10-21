package api

import (
	"github.com/francescofrontera/ks-job-upload/dockerutils"
	"github.com/go-chi/chi"
	"io"
	"log"
	"net/http"
	"os"
)

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
	dockerClient := &dockerutils.DockerClientBuilder{}
	result := dockerClient.InitClient().BuildImage("vertx-start-project-1.0-SNAPSHOT-fat.jar").Build()
	io.WriteString(w, result.ImageName)
}

func UploaderRoute() *chi.Mux  {
	uploadRoute := chi.NewRouter()
	uploadRoute.Post("/upload", UploadHandler)
	uploadRoute.Post("/build", BuilderHandler)
	return uploadRoute
}
