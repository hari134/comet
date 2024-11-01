package pipeline

import (
    "fmt"
    "go.uber.org/dig"
)

// DependencyManager defines the interface for dependency management.
type DependencyManager interface {
    AddDependency(constructor interface{}) DependencyManager
    Invoke(function interface{}) error
}

// digDependencyManager is a concrete implementation of the DependencyManager interface.
type digDependencyManager struct {
    container *dig.Container
}

// NewDependencyManager creates a new instance of digDependencyManager.
func NewDefaultDependencyManager() DependencyManager {
    return &digDependencyManager{container: dig.New()}
}

// AddDependency registers a constructor in the container.
func (dm *digDependencyManager) AddDependency(constructor interface{}) DependencyManager {
    err := dm.container.Provide(constructor)
    if err != nil {
        fmt.Printf("Error adding dependency: %v\n", err)
    }
    return dm
}

// Invoke runs a function with dependencies injected from the container.
func (dm *digDependencyManager) Invoke(function interface{}) error {
    return dm.container.Invoke(function)
}
