package dockerutils

import (
	"context"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func BuildContainer(fileName string) string {
	cli, _ := client.NewClientWithOpts(client.WithVersion("1.38"))
	ctx := context.Background()

	dockerBuildContext, err := os.Open("/Users/francescofrontera/go/src/github.com/francescofrontera/ks-job-upload/Dockerfile.tar.gz")
	defer dockerBuildContext.Close()

	imageName := "internal_image/"+strings.ToLower(strings.TrimSpace(strings.TrimSuffix(fileName, filepath.Ext(fileName))))
	log.Print(imageName)

	buildOptions := types.ImageBuildOptions{
		Tags: []string{imageName},
		Dockerfile: "docker/DockerFile",
		BuildArgs: map[string]*string{
			"JAR_TO_EXECUTE": &fileName,
		},
	}

	response, err := cli.ImageBuild(ctx, dockerBuildContext, buildOptions); if err != nil {
		log.Fatal(err)
	}

	io.Copy(os.Stdout, response.Body)

	defer response.Body.Close()
	return imageName
}