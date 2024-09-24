package container

import (
	"archive/tar"
	"archive/zip"
	"bytes"
	"context"
	"fmt"
	"io"
	"path/filepath"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

type Image string

// Implements the BuildContainer interface
type DockerBuildContainer struct {
	id     string
	image  Image
	client *client.Client
}

func NewDockerBuildContainer() *DockerBuildContainer {
	return &DockerBuildContainer{}
}

// Builder functions

func (c *DockerBuildContainer) WithImage(image Image) *DockerBuildContainer {
	c.image = image
	return c
}

func (c *DockerBuildContainer) WithClient(client *client.Client) *DockerBuildContainer {
	c.client = client
	return c
}

func (c *DockerBuildContainer) Create() (*DockerBuildContainer, error) {
	ctx := context.Background()
	containerConfig := &container.Config{
		Image: string(c.image),
	}

	resp, err := c.client.ContainerCreate(ctx, containerConfig, nil, nil, nil, "")
	if err != nil {
		return nil, err
	}
	c.id = resp.ID
	return c, nil
}

// BuildContainer interface functions

func (c *DockerBuildContainer) CopyToContainer(content io.Reader, containerPath string) error {
	ctx := context.Background()
	if err := c.client.CopyToContainer(ctx, c.id, "/", content, types.CopyToContainerOptions{}); err != nil {
		return err
	}
	return nil
}

func (c *DockerBuildContainer) CopyFromContainer() (io.ReadCloser, error) {
	distData, _, err := c.client.CopyFromContainer(context.Background(), c.id, "/dist/")
	if err != nil {
		return nil, err
	}
	return distData, nil
}

func (c *DockerBuildContainer) Start() error {
	return c.client.ContainerStart(context.Background(), c.id, container.StartOptions{})
}

func (c *DockerBuildContainer) Stop() error {
	err := c.client.ContainerStop(context.Background(), c.id, container.StopOptions{})
	if err != nil {
		return fmt.Errorf("container stop error: %v", err)
	}
	return nil
}

func (c *DockerBuildContainer) Remove() error {
	return c.client.ContainerRemove(context.Background(), c.id, container.RemoveOptions{})
}


func (c *DockerBuildContainer) ExecCmd(cmd string) (string, error) {
	execResp, err := c.client.ContainerExecCreate(context.Background(), c.id, types.ExecConfig{
		Cmd:          []string{"sh", "-c", cmd},
		AttachStdout: true,
		AttachStderr: true,
	})
	if err != nil {
		return "", err
	}

	execAttachResp, err := c.client.ContainerExecAttach(context.Background(), execResp.ID, types.ExecStartCheck{})
	if err != nil {
		return "", err
	}
	defer execAttachResp.Close()

	var outputBuf bytes.Buffer
	if _, err := io.Copy(&outputBuf, execAttachResp.Reader); err != nil {
		return "", err
	}

	return outputBuf.String(), nil
}

func (c *DockerBuildContainer) unzipFile(filePath string) (string, error) {
	cmd := fmt.Sprintf("unzip %s -d %s", filePath, filepath.Dir(filePath))
	return c.ExecCmd(cmd)
}

func (c *DockerBuildContainer) createDirectoryInContainer(directoryPath string) error {
	execOptions := types.ExecConfig{
		Cmd:          []string{"mkdir", "-p", directoryPath},
		AttachStdout: true,
		AttachStderr: true,
	}

	resp, err := c.client.ContainerExecCreate(context.Background(), c.id, execOptions)
	if err != nil {
		return err
	}

	err = c.client.ContainerExecStart(context.Background(), resp.ID, types.ExecStartCheck{})
	if err != nil {
		return err
	}

	execResult, err := c.client.ContainerExecInspect(context.Background(), resp.ID)
	if err != nil {
		return err
	}

	if execResult.ExitCode != 0 {
		return fmt.Errorf("command failed with exit code %d", execResult.ExitCode)
	}

	return nil
}

// Utility functions

func (c *DockerBuildContainer) createFileInContainer(filePath string) error {
	execOptions := types.ExecConfig{
		Cmd:          []string{"touch", filePath},
		AttachStdout: true,
		AttachStderr: true,
	}

	resp, err := c.client.ContainerExecCreate(context.Background(), c.id, execOptions)
	if err != nil {
		return err
	}

	err = c.client.ContainerExecStart(context.Background(), resp.ID, types.ExecStartCheck{})
	if err != nil {
		return err
	}

	execResult, err := c.client.ContainerExecInspect(context.Background(), resp.ID)
	if err != nil {
		return err
	}

	if execResult.ExitCode != 0 {
		return fmt.Errorf("command failed with exit code %d", execResult.ExitCode)
	}

	return nil
}

func (c *DockerBuildContainer) writeDataToContainer(data []byte, filePath string) error {
	// Create a reader for the data buffer
	var buf bytes.Buffer

	// Create a new TAR writer
	tw := tar.NewWriter(&buf)

	tarHeader := &tar.Header{
		Name: filePath,
		Size: int64(len(data)),
	}
	if err := tw.WriteHeader(tarHeader); err != nil {
		return err
	}
	c.createFileInContainer(filePath)
	// Write the file content to the TAR archive
	if _, err := tw.Write(data); err != nil {
		return err
	}

	if err := tw.Close(); err != nil {
		return err
	}
	err := c.client.CopyToContainer(context.Background(), c.id, filePath, &buf, types.CopyToContainerOptions{})
	if err != nil {
		return err
	}

	return nil
}


func addFileToZip(zipWriter *zip.Writer, reader io.Reader, filename string) error {
	writer, err := zipWriter.Create(filename)
	if err != nil {
		return err
	}

	_, err = io.Copy(writer, reader)
	return err
}

