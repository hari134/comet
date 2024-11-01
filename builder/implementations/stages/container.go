package stages

import (
	"bytes"
	"log/slog"

	"github.com/hari134/comet/builder/container"
	"github.com/hari134/comet/builder/pipeline"
)

// CopyProjectFilesToContainer creates a stage to copy project files to the container.
func CopyProjectFilesToContainer() pipeline.Stage {
	return pipeline.Stage{
		Name: "Copy Project Files to Container",
		Execute: func(dm pipeline.DependencyManager) error {
			var buildContainer container.BuildContainer
			var projectTarFile *bytes.Buffer
			err := dm.Invoke(func(bc container.BuildContainer,pf *bytes.Buffer) {
				buildContainer = bc
				projectTarFile = pf
			})
			if err != nil {
				return err
			}

			buildContainer.CopyToContainer(projectTarFile, "/app")
			slog.Debug("Files copied successfully")
			return nil
		},
	}
}

// ExtractProject creates a stage for extracting the project archive in the container.
func ExtractProject() pipeline.Stage {
	return pipeline.Stage{
		Name: "Extract Archive",
		Execute: func(dm pipeline.DependencyManager) error {
			var buildContainer container.BuildContainer
			err := dm.Invoke(func(bc container.BuildContainer) {
				buildContainer = bc
			})
			if err != nil {
				return err
			}

			_, err = buildContainer.ExecCmd("tar -xvf /app/full.tar -C /app")
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
		Execute: func(dm pipeline.DependencyManager) error {
			var buildContainer container.BuildContainer
			err := dm.Invoke(func(bc container.BuildContainer) {
				buildContainer = bc
			})
			if err != nil {
				return err
			}

			_, err = buildContainer.ExecCmd("cd /app && npm install")
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
		Execute: func(dm pipeline.DependencyManager) error {
			var buildContainer container.BuildContainer
			err := dm.Invoke(func(bc container.BuildContainer) {
				buildContainer = bc
			})
			if err != nil {
				return err
			}

			_, err = buildContainer.ExecCmd("cd /app && npm run build")
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
		Execute: func(dm pipeline.DependencyManager) error {
			var buildContainer container.BuildContainer
			err := dm.Invoke(func(bc container.BuildContainer) {
				buildContainer = bc
			})
			if err != nil {
				return err
			}

			_, err = buildContainer.CopyFromContainer("/app/dist")
			if err != nil {
				return err
			}
			slog.Debug("Build files copied successfully.")
			return nil
		},
	}
}
