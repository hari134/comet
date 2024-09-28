package container

import (
	"errors"

	"github.com/docker/docker/client"
)

type ContainerManager interface {
	NewBuildContainer(buildType string) (BuildContainer,error)
}

type DockerContainerManager struct {
	capacity int // concurrency limit the number of container to run concurrently
	client   *client.Client
}

func NewDockerContainerManager() *DockerContainerManager{
	return &DockerContainerManager{}
}

func (dcm *DockerContainerManager) WithCapacity(capacity int) *DockerContainerManager{
	dcm.capacity = capacity
	return dcm
}

func (dcm *DockerContainerManager) WithClient(client *client.Client) *DockerContainerManager{
	dcm.client = client
	return dcm
}

func (cm *DockerContainerManager) NewBuildContainer(buildType string) (BuildContainer,error) {
	switch buildType {
	case "ReactViteNode20":
		dockerContainer,err := NewDockerBuildContainer().
			WithImage("node:20").
			WithClient(cm.client).
			Create()

		if err != nil{
			return nil,err
		}
		return dockerContainer, nil
	default:
		return nil,errors.New("container for given build type not found")
	}
}
