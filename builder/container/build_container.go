package container

import (
	"bytes"
	"io"
)

/*
	Build container is a stateless container using technology such as docker and its implementation is further
	used to build pipelines.
	Build container interface contains function signatures for basic container functionalities such as :
   1. Write directories/files to container
   2. Get direcotories/files from container
   3. Start the container
   4. Stop the container
   5. Execute commands in the container
*/

type ExecOptions interface{
	IsStreamingEnabled() bool
}

type BuildContainer interface {
	CopyToContainer(tarFile *bytes.Buffer, containerPath string) error
	CopyFromContainer(containerPath string) (io.ReadCloser, error)
	Start() error
	Stop() error
	Remove() error
	ExecCmd(opts ExecOptions) (string, error)
}
