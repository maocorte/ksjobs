package dockerutils

import (
	"context"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/francescofrontera/ks-job-upload/utils"
	"io"
	"log"
	"os"
)

/* Docker Client Initialization */
type DockerClientResult struct {
	dockerClient *client.Client
	ctx context.Context
}

func InitClient() *DockerClientResult {
	cli, error := client.NewClientWithOpts(client.WithVersion("1.38")); if error != nil {
		log.Fatal(error)
	}
	ctx := context.Background()

	return &DockerClientResult{
		dockerClient: cli,
		ctx: ctx,
	}
}

/* DockerClient utils */
func getDockerFileCtx() (*os.File, error) {
	ctx, error := os.Open(utils.GetDockerFilePath())
	return ctx, error
}

/* Docker Client Result methods */
func (dcb *DockerClientResult) BuildImage(fileName string) string {
	dockerBuildContext, _ := getDockerFileCtx()
	defer dockerBuildContext.Close()

	cli := dcb.dockerClient
	ctx := dcb.ctx

	imageName := utils.NormalizeJarName(fileName)

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

func (dcb *DockerClientResult) RunContainer(imageName string) string {
	cli := dcb.dockerClient
	ctx := dcb.ctx

	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image: imageName,
		Tty:   true,
	}, nil, nil, ""); if err != nil {
		panic(err)
	}

	containerId := resp.ID

	if err := cli.ContainerStart(ctx, containerId, types.ContainerStartOptions{}); err != nil {
		panic(err)
	}

	return containerId
}

