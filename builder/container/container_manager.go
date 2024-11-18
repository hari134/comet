package container

import (
	"errors"
	"log"

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

func (dcm *DockerContainerManager) WithDefaultClient() *DockerContainerManager{
	client , err := client.NewClientWithOpts()
	if err != nil{
		log.Fatal("Failed to initialize default client for  container manager")
	}
	dcm.client = client
	return dcm
}

func (cm *DockerContainerManager) NewBuildContainer(buildType string) (BuildContainer,error) {
	switch buildType {
	case "reactvitenode20":
		dockerContainer,err := NewDockerBuildContainer().
			WithImage("comet-react-node20:v1.0").
			WithClient(cm.client).
			Create()

		if err != nil{
			return nil,err
		}
		if err := dockerContainer.Start(); err != nil{
			return nil, err
		}
		return dockerContainer, nil
	default:
		return nil,errors.New("container for given build type not found")
	}
}
