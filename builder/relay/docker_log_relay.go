package relay

import (
	"errors"
	"sync"
)

type DockerLogData struct {
	Data  string
	Source string
}

func (logData DockerLogData) GetData() string{
	return logData.Data
}


func (logData DockerLogData) GetSource() string{
	return logData.Source
}


type DockerLogRelay struct {
	outputChannel chan StreamData
	isActive      bool
	mu            sync.Mutex
}

func NewDockerLogRelay(bufferSize int) *DockerLogRelay {
	return &DockerLogRelay{
		outputChannel: make(chan StreamData, bufferSize),
		isActive:      false,
	}
}

func (d *DockerLogRelay) Start() error {
	d.mu.Lock()
	defer d.mu.Unlock()
	if d.isActive {
		return errors.New("relay is already active")
	}
	d.isActive = true
	return nil
}

func (d *DockerLogRelay) Stop() {
	d.mu.Lock()
	defer d.mu.Unlock()
	if d.isActive {
		d.isActive = false
		close(d.outputChannel)
	}
}

func (d *DockerLogRelay) Send(data StreamData) {
	d.mu.Lock()
	defer d.mu.Unlock()
	if d.isActive {
		d.outputChannel <- data
	}
}

func (d *DockerLogRelay) Receive() <-chan StreamData {
	return d.outputChannel
}

func (d *DockerLogRelay) Close() {
	d.Stop()
}





