package pipeline

import (
	"fmt"
	cont "github.com/hari134/comet/builder/container"
)

// Stage defines an interface for pipeline stages
type Stage interface {
	Execute(ctx *PipelineContext) error
}

// CommandStage is a stage that runs a command inside a container
type CommandStage struct {
	command string
}

func NewCommandStage(command string) *CommandStage {
	return &CommandStage{command: command}
}

func (s *CommandStage) Execute(ctx *PipelineContext) error {
	container, err := ctx.Get("container")
	if err != nil {
		return err
	}
	_, err = container.(cont.BuildContainer).ExecCmd(s.command)
	if err != nil {
		return fmt.Errorf("command stage failed: %w", err)
	}
	return nil
}

// FunctionStage is a stage that runs a custom function
type FunctionStage struct {
	fn func(ctx *PipelineContext) error
}

func NewFunctionStage(fn func(ctx *PipelineContext) error) *FunctionStage {
	return &FunctionStage{fn: fn}
}

func (s *FunctionStage) Execute(ctx *PipelineContext) error {
	return s.fn(ctx)
}
