package docker_client

import (
	"github.com/easy-oj/common/logs"
	"github.com/easy-oj/common/settings"
	"github.com/fsouza/go-dockerclient"
)

var (
	Client *docker.Client
)

func InitDockerClient() {
	Client = dial(settings.Judger.Docker)
}

func dial(endpoint string) *docker.Client {
	client, err := docker.NewClient(endpoint)
	if err != nil {
		panic(err)
	}
	logs.Info("[DockerClient] dial %s", endpoint)
	return client
}
