package container

import (
	"bytes"
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"log/slog"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/client"
	"github.com/hari134/comet/builder/relay"
)

type Image string

type DockerExecOptions struct {
	cmd           string
	streamOptions DockerStreamOptions
}

type DockerStreamOptions struct {
	IsStreamingEnabled bool
	Channel            chan relay.StreamData
}

func DefaultDockerExecOptions() *DockerExecOptions {
	return &DockerExecOptions{}
}

func (dockerExecOptions *DockerExecOptions) IsStreamingEnabled() bool {
	return dockerExecOptions.streamOptions.IsStreamingEnabled
}

func (dockerExecOptions *DockerExecOptions) WithCommand(cmd string) *DockerExecOptions {
	dockerExecOptions.cmd = cmd
	return dockerExecOptions
}

func (dockerExecOptions *DockerExecOptions) WithStreamOptions(opts DockerStreamOptions) (*DockerExecOptions, error) {
	if opts.IsStreamingEnabled && opts.Channel == nil {
		return nil, errors.New("stream is enabled but no data channel is provided")
	}
	return dockerExecOptions, nil
}

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

	_, _, err := c.client.ImageInspectWithRaw(ctx, string(c.image))
	if err != nil {
		slog.Debug("image not found locally, attempting to pull", "image", c.image)
		_, err := c.client.ImagePull(ctx, string(c.image), image.PullOptions{})
		if err != nil {
			return nil, fmt.Errorf("failed to pull image: %v", err)
		}
	}

	containerConfig := &container.Config{
		Image: string(c.image),
	}
	resp, err := c.client.ContainerCreate(ctx, containerConfig, nil, nil, nil, "")
	if err != nil {
		return nil, fmt.Errorf("failed to create container: %v", err)
	}
	c.id = resp.ID
	return c, nil
}

// BuildContainer interface functions

func (c *DockerBuildContainer) CopyToContainer(content *bytes.Buffer, containerPath string) error {
	ctx := context.Background()
	if err := c.client.CopyToContainer(ctx, c.id, containerPath, content, types.CopyToContainerOptions{}); err != nil {
		return err
	}
	return nil
}

func (c *DockerBuildContainer) CopyFromContainer(containerPath string) (io.ReadCloser, error) {
	distData, _, err := c.client.CopyFromContainer(context.Background(), c.id, containerPath)
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

func (c *DockerBuildContainer) ExecCmd(opts ExecOptions) (string, error) {
	execOpts, ok := opts.(*DockerExecOptions)
	if !ok {
		return "", errors.New("invalid docker exec config")
	}
	execResp, err := c.client.ContainerExecCreate(context.Background(), c.id, types.ExecConfig{
		Cmd:          []string{"sh", "-c", execOpts.cmd},
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
	buffer := make([]byte, 4096)

	for {
		// Read the 8-byte Docker header
		header := make([]byte, 8)
		_, err := io.ReadFull(execAttachResp.Reader, header)
		if err != nil {
			if err == io.EOF {
				break
			}
			return "", fmt.Errorf("error reading header: %w", err)
		}

		// Determine the stream type (stdout or stderr)
		streamType := header[0]
		length := int(binary.BigEndian.Uint32(header[4:]))

		// Read the actual payload based on the length
		if length > 0 {
			n, err := execAttachResp.Reader.Read(buffer[:length])
			if n > 0 {
				outputBuf.Write(buffer[:n])

				if opts.IsStreamingEnabled() {
					logData := relay.DockerLogData{
						Data: string(buffer[:n]),
					}
					execOpts.streamOptions.Channel <- logData
				}
			}
			if err != nil {
				if err == io.EOF {
					break
				}
				return "", fmt.Errorf("error reading from log relay: %w", err)
			}
		}

		// Optional: Handle stream type if needed
		if streamType == 2 {
			// This is stderr, you can handle it separately if needed
		}
	}

	return outputBuf.String(), nil
}
