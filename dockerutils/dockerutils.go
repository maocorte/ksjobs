package dockerutils

import (
	"context"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/francescofrontera/ks-job-upload/utils"
	"io"
	"log"
	"os"
)

type DockerClientResult struct {
	dockerClient *client.Client
	ctx context.Context
	ImageName string
	ContainerName string
}

type DockerClientBuilder struct {
	DockerClientResult
}

func (dcb *DockerClientBuilder) InitClient() *DockerClientBuilder {
	cli, error := client.NewClientWithOpts(client.WithVersion("1.38")); if error != nil {
		log.Panicf("Error during docker client initialization..", error)
	}
	ctx := context.Background()

	dcb.dockerClient = cli
	dcb.ctx = ctx

	return dcb
}

func (dcb *DockerClientBuilder) BuildImage(fileName string) *DockerClientBuilder {
	dockerBuildContext, _ := os.Open(utils.GetDockerFilePath())
	defer dockerBuildContext.Close()

	imageName := utils.NormalizeJarName(fileName)

	buildOptions := types.ImageBuildOptions{
		Tags: []string{imageName},
		Dockerfile: "docker/DockerFile",
		BuildArgs: map[string]*string{
			"JAR_TO_EXECUTE": &fileName,
		},
	}

	response, err := dcb.dockerClient.ImageBuild(dcb.ctx, dockerBuildContext, buildOptions); if err != nil {
		log.Fatal(err)
	}

	io.Copy(os.Stdout, response.Body)

	defer response.Body.Close()

	dcb.ImageName = imageName

	return dcb
}

func (dcb *DockerClientBuilder) Build() DockerClientResult {
	return dcb.DockerClientResult
}