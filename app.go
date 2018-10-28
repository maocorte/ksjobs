package main

import (
	"github.com/francescofrontera/ks-job-uploader/api"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func main() {
	route := chi.NewRouter()
	route.Use(middleware.Logger)

	route.Route("/v1", func(r chi.Router) {
		r.Mount("/api", api.UploaderRoute())
	})

	http.ListenAndServe(":3000", route)
}