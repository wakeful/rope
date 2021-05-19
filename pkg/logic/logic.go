package logic

import (
	"context"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/sirupsen/logrus"
	"rope/pkg/config"
	"rope/pkg/helpers"
)

//StopContainers stop containers from the config file
func StopContainers(ctx context.Context, dockerClient *client.Client, containers []types.Container) {
	for _, n := range containers {
		if stopErr := stopContainer(ctx, dockerClient, n.ID); stopErr != nil {
			logrus.WithField("error", stopErr).Error("stopping container")
			continue
		}
	}
}

func stopContainer(ctx context.Context, dockerClient *client.Client, id string) error {
	logrus.WithField("ID", id[:10]).Info("stopping container")
	if stopErr := dockerClient.ContainerStop(ctx, id, nil); stopErr != nil {
		return stopErr
	}
	return nil
}

//RunContainers start containers from the config file
func RunContainers(ctx context.Context, client *client.Client, config config.File) {
	hash := helpers.GetHash(helpers.ProjectDir)
	for imageName := range config.Services {
		if err := runContainer(ctx, client, imageName, hash); err != nil {
			logrus.WithField("error", err).Error("starting container")
		}
	}
}

func runContainer(ctx context.Context, client *client.Client, imageName string, hash string) error {
	pullErr := helpers.PullImage(ctx, client, imageName)
	if pullErr != nil {
		return pullErr
	}
	create, createErr := client.ContainerCreate(ctx, &container.Config{
		Image: imageName,
		Labels: map[string]string{
			"rope": hash,
		},
	}, nil, nil, nil, "")
	if createErr != nil {
		return createErr
	}
	if errStart := client.ContainerStart(ctx, create.ID, types.ContainerStartOptions{}); errStart != nil {
		return errStart
	}
	logrus.WithFields(logrus.Fields{
		"ID":    create.ID[:10],
		"image": imageName,
	}).Info("container started")
	return nil
}

//MngContainers the main logic loop
func MngContainers(ctx context.Context, config config.File, client *client.Client, running map[string][]string) {
	hash := helpers.GetHash(helpers.ProjectDir)
	for image, replicas := range config.Services {
		logrus.WithFields(logrus.Fields{
			"image": image,
			"have":  len(running[image]),
			"want":  replicas,
		}).Info("want container")
		if len(running[image]) < replicas {
			logrus.WithFields(logrus.Fields{
				"image": image,
				"have":  len(running[image]),
				"want":  replicas,
			}).Warning("starting container")
			for i := len(running[image]); i < replicas; i++ {
				if errRun := runContainer(ctx, client, image, hash); errRun != nil {
					logrus.WithField("error", errRun).Error("run container")
				}
			}
		} else if len(running[image]) > replicas {
			logrus.WithFields(logrus.Fields{
				"image": image,
				"have":  len(running[image]),
				"want":  replicas,
			}).Warning("stopping container")
			for i := len(running[image]); i > replicas; i-- {
				if errStop := stopContainer(ctx, client, running[image][i-1]); errStop != nil {
					logrus.WithField("error", errStop).Error("run container")
				}
			}
		}
	}
}

//CountContainers create map of images with a slice of ID
func CountContainers(input []types.Container) map[string][]string {
	var output = map[string][]string{}
	for _, n := range input {
		logrus.WithFields(logrus.Fields{
			"ID":     n.ID[:10],
			"image":  n.Image,
			"status": n.Status,
		}).Debug("found container")
		output[n.Image] = append(output[n.Image], n.ID)
	}
	return output
}
