package reactvitenode20

import (
	"context"
	"log/slog"

	"github.com/hari134/comet/builder/container"
	"github.com/hari134/comet/builder/pipeline"
)

func PullProjectFiles() pipeline.Stage {
	return pipeline.Stage{
		Name: "Pull Project Files From Storage",
		Execute: func(config pipeline.PipelineConfig) error {
			cfg := config.(*PipelineConfig)
			store := cfg.Store
			storageConfig := cfg.ProjectStorageConfig

			projectStorageKey := storageConfig.ProjectStorageKey
			projectStorageBucket := storageConfig.ProjectStorageBucket

			projectTarFile, err := store.Get(context.Background(), projectStorageBucket, projectStorageKey)
			if err != nil {
				return err
			}
			cfg.ProjectFileData = ProjectFileData{ProjectTarFile: projectTarFile, DirName: "app"}
			slog.Debug("Files pulled from storage successfully")
			return nil
		},
	}
}

// CopyProjectFilesToContainer creates a stage to copy project files to the container.
func CopyProjectFilesToContainer() pipeline.Stage {
	return pipeline.Stage{
		Name: "Copy Project Files to Container",
		Execute: func(config pipeline.PipelineConfig) error {
			cfg := config.(*PipelineConfig)
			buildContainer := cfg.BuildContainer
			projectTarFileData := cfg.ProjectFileData

			buildContainer.CopyToContainer(projectTarFileData.ProjectTarFile, projectTarFileData.DirName)
			slog.Debug("Files copied successfully")
			return nil
		},
	}
}

// ExtractProject creates a stage for extracting the project archive in the container.
func ExtractProject() pipeline.Stage {
	return pipeline.Stage{
		Name: "Extract Archive",
		Execute: func(config pipeline.PipelineConfig) error {
			cfg := config.(*PipelineConfig)
			buildContainer := cfg.BuildContainer

			execOpts := container.DefaultDockerExecOptions().
				WithCommand("tar -xvf /app/full.tar -C /app")

			execOpts, err := execOpts.WithStreamOptions(container.DockerStreamOptions{
				IsStreamingEnabled: cfg.StreamConfig.StreamingEnabled,
				Channel:            cfg.StreamConfig.Output,
			})

			if err != nil {
				return err
			}
			_, err = buildContainer.ExecCmd(execOpts)
			if err != nil {
				return err
			}
			slog.Debug("Extraction completed successfully.")
			return nil
		},
	}
}

// InstallNpmDependencies creates a stage for installing npm dependencies in the container.
func InstallNpmDependencies() pipeline.Stage {
	return pipeline.Stage{
		Name: "Install Dependencies",
		Execute: func(config pipeline.PipelineConfig) error {
			cfg := config.(*PipelineConfig)
			buildContainer := cfg.BuildContainer
			execOpts := container.DefaultDockerExecOptions().
				WithCommand("cd /app && npm install")

			execOpts, err := execOpts.WithStreamOptions(container.DockerStreamOptions{
				IsStreamingEnabled: cfg.StreamConfig.StreamingEnabled,
				Channel:            cfg.StreamConfig.Output,
			})

			if err != nil{
				return err
			}
			_, err = buildContainer.ExecCmd(execOpts)
			if err != nil {
				return err
			}
			slog.Debug("Dependencies installed successfully.")
			return nil
		},
	}
}

// NpmBuild creates a stage for building the project using npm in the container.
func NpmBuild() pipeline.Stage {
	return pipeline.Stage{
		Name: "Build Project",
		Execute: func(config pipeline.PipelineConfig) error {
			cfg := config.(*PipelineConfig)
			buildContainer := cfg.BuildContainer

			execOpts := container.DefaultDockerExecOptions().
				WithCommand("cd /app && npm run build")

			execOpts, err := execOpts.WithStreamOptions(container.DockerStreamOptions{
				IsStreamingEnabled: cfg.StreamConfig.StreamingEnabled,
				Channel:            cfg.StreamConfig.Output,
			})
			if err != nil{
				return err
			}
			_, err = buildContainer.ExecCmd(execOpts)
			if err != nil {
				return err
			}
			slog.Debug("Build completed successfully.")
			return nil
		},
	}
}

// CopyBuildFilesFromContainer creates a stage to copy build files from the container.
func CopyBuildFilesFromContainer() pipeline.Stage {
	return pipeline.Stage{
		Name: "Copy Build Files From Container",
		Execute: func(config pipeline.PipelineConfig) error {
			cfg := config.(*PipelineConfig)
			buildContainer := cfg.BuildContainer

			_, err := buildContainer.CopyFromContainer("/app/dist")
			if err != nil {
				return err
			}
			slog.Debug("Build files copied successfully.")
			return nil
		},
	}
}
