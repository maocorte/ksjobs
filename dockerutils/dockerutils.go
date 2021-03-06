package dockerutils

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
	"github.com/francescofrontera/ksjobs/utils"
	"log"
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

func (dcb *DockerClientResult) RunContainer(jarToMount, mainClass string) (string, error) {
	cli := dcb.dockerClient
	ctx := dcb.ctx

	containerConfig := &container.Config{
		Image: "base_image_jar",
		Tty:   true,
		Env: []string{
			fmt.Sprintf("JAR_TO_EXECUTE=%s", jarToMount),
			fmt.Sprintf("MAIN_CLASS=%s", mainClass),
		},
	}

	sourcePath, targetPath := utils.GetPathToJar(jarToMount)
	hostConfig := &container.HostConfig{
		Mounts: []mount.Mount{
			{
				Type: mount.TypeBind,
				Source: sourcePath,
				Target: targetPath,
			},
		},
	}

	resp, err := cli.ContainerCreate(ctx, containerConfig, hostConfig, nil, ""); if err != nil {
		return "", err
	}

	containerId := resp.ID

	if err := cli.ContainerStart(ctx, containerId, types.ContainerStartOptions{}); err != nil {
		return "", err
	}

	return containerId, nil
}

