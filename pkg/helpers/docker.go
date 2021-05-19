package helpers

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"github.com/sirupsen/logrus"
)

//ProjectDir holds CWD used to generate docker label filter
var ProjectDir = ""

//ListContainers list docker containers that match our filter
func ListContainers(ctx context.Context, client *client.Client) ([]types.Container, error) {
	filter := GetFilter(ProjectDir)
	containers, err := client.ContainerList(ctx, types.ContainerListOptions{
		Filters: filter,
	})
	if err != nil {
		return nil, err
	}
	return containers, nil
}

//GetFilter we use the input to create a docker label filter
func GetFilter(input string) filters.Args {

	newArgs := filters.NewArgs()
	newArgs.Add("label", "rope="+GetHash(input))

	logrus.WithField("filter", newArgs).Debug("container list using")
	return newArgs
}

//GetHash calculate MD5 sum of input string
func GetHash(input string) string {
	sum := md5.Sum([]byte(input))
	hash := hex.EncodeToString(sum[:])
	return hash
}

//NewDockerClient create docker connection
func NewDockerClient() (*client.Client, error) {
	dockerClient, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation()) // @todo add remote host access
	if err != nil {
		return nil, err
	}
	return dockerClient, nil
}

//PullImage pull docker image from registry
func PullImage(ctx context.Context, dockerClient *client.Client, imageName string) error {
	_, err := dockerClient.ImagePull(ctx, imageName, types.ImagePullOptions{})
	if err != nil {
		return err
	}
	return nil
}
