package helpers

import (
	"context"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"os"
	"reflect"
	"testing"
)

func TestGetFilter(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  filters.Args
	}{
		{
			input: "/projects/fizzBuzz",
			want: filters.NewArgs(filters.KeyValuePair{
				Key:   "label",
				Value: "rope=afcae8ea1db6e31eaae4643b2325fd2c",
			}),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetFilter(tt.input); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetFilter() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestListContainersWithNoContainersRunning(t *testing.T) {
	dockerClient, errNewClient := NewDockerClient()
	ProjectDir = "/non_existing_path/42"
	if errNewClient != nil {
		t.Error("unable to connect to docker")
	}

	type args struct {
		ctx    context.Context
		client *client.Client
	}
	tests := []struct {
		name    string
		args    args
		want    []types.Container
		wantErr bool
	}{
		{
			args: args{
				ctx:    context.TODO(),
				client: dockerClient,
			},
			want:    []types.Container{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ListContainers(tt.args.ctx, tt.args.client)
			if (err != nil) != tt.wantErr {
				t.Errorf("ListContainers() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ListContainers() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewDockerClientError(t *testing.T) {
	_ = os.Setenv("DOCKER_HOST", "1")
	tests := []struct {
		name    string
		want    *client.Client
		wantErr bool
	}{
		{
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewDockerClient()
			if (err != nil) != tt.wantErr {
				t.Errorf("NewDockerClient() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewDockerClient() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPullImage(t *testing.T) {
	_ = os.Setenv("DOCKER_HOST", "")
	dockerClient, errClient := NewDockerClient()
	if errClient != nil {
		t.Error(errClient)
	}
	type args struct {
		ctx       context.Context
		client    *client.Client
		imageName string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			args: args{
				ctx:       context.TODO(),
				client:    dockerClient,
				imageName: "nginx",
			},
			wantErr: false,
		},
		{
			args: args{
				ctx:       context.TODO(),
				client:    dockerClient,
				imageName: "non.existing.registry.48k.io/fake",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := PullImage(tt.args.ctx, tt.args.client, tt.args.imageName); (err != nil) != tt.wantErr {
				t.Errorf("PullImage() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
