package api

import (
	"fmt"
	"github.com/francescofrontera/ksjobs/dockerutils"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"strings"
)

type (
	RunJar struct {
		JarName string `json:jarName`
		MainClass string `json:mainClass`
	}
)

func UploadHandler(c *gin.Context){
	file, err := c.FormFile("uploadFile")
	if err != nil {
		panic(err)
	}

	workdir, err := os.Getwd(); if err != nil {
		panic(err)
	}

	dst := strings.Join([]string{workdir, "jars", file.Filename}, "/")
	if err := c.SaveUploadedFile(file, dst); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("upload file error %s", err.Error())})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{"fileName": file.Filename, "status": http.StatusAccepted})
}

func RunKSJob(c *gin.Context){
	initDockerClient := dockerutils.InitClient()

	runJar := &RunJar{}
	if error := c.BindJSON(runJar); error == nil {
		containerId, dockerClientError := initDockerClient.RunContainer(runJar.JarName, runJar.MainClass); if dockerClientError != nil {
			responseStatus := http.StatusInternalServerError
			c.JSON(responseStatus, gin.H{"error": dockerClientError.Error(), "status": responseStatus})
			return
		}

		responseStatus := http.StatusAccepted
		c.JSON(responseStatus, gin.H{"containerId": containerId, "status": responseStatus})
	}
}