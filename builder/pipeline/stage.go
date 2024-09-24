package pipeline

import (
	"github.com/hari134/comet/builder/container"
)

type Stage struct {
	command string
}

func NewStage(command string) Stage {
	return Stage{
		command,
	}
}

func (stage Stage) Execute(container container.BuildContainer) (string,error){
	resp, err := container.ExecCmd(stage.command)
	if err != nil{
		return "",err
	}
	return resp,nil
}
