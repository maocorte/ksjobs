package main

import (
	"github.com/francescofrontera/ks-job-uploader/api"
	"github.com/gin-gonic/gin"
)

func main() {
	route := gin.Default()

	v1 := route.Group("/v1/api")
	{
		v1.POST("/upload", api.UploadHandler)
		v1.POST("/run", api.RunKSJob)
	}

	route.Run(":8080")
}