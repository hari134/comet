package reactvitenode20

import (
	"archive/tar"
	"bytes"
	"context"
	"fmt"
	"github.com/hari134/comet/core/cdn"
	"io"
	"log/slog"
	"strings"

	"github.com/hari134/comet/builder/container"
	"github.com/hari134/comet/builder/pipeline"
)

func PullProjectFiles() pipeline.Stage {
	return pipeline.Stage{
		Name: "Pull Project Files From Storage",
		Execute: func(config *pipeline.PipelineConfig) error {
			slog.Debug("In stage : Pull Project Files From Storage")
			pipelineOptions := (*config).(*PipelineOptions)
			store := pipelineOptions.Store
			storageConfig := pipelineOptions.ProjectStorageConfig

			projectStorageKey := storageConfig.ProjectFilesKey
			projectStorageBucket := storageConfig.ProjectFilesBucket

			projectTarFile, err := store.Get(context.Background(), projectStorageBucket, projectStorageKey)
			if err != nil {
				return err
			}
			pipelineOptions.ProjectFileData = ProjectFileData{ProjectTarFile: projectTarFile, DirName: "/"}
			slog.Debug("Files pulled from storage successfully")
			return nil
		},
	}
}

// CopyProjectFilesToContainer creates a stage to copy project files to the container.
func CopyProjectFilesToContainer() pipeline.Stage {
	return pipeline.Stage{
		Name: "Copy Project Files to Container",
		Execute: func(config *pipeline.PipelineConfig) error {
			slog.Debug("In stage : Copy Project Files to Container")
			pipelineOptions := (*config).(*PipelineOptions)
			buildContainer := pipelineOptions.BuildContainer
			projectTarFileData := pipelineOptions.ProjectFileData
			err := buildContainer.CopyToContainer(projectTarFileData.ProjectTarFile, projectTarFileData.DirName)
			if err != nil {
				slog.Debug(err.Error())
				return err
			}
			slog.Debug("Files copied successfully")
			return nil
		},
	}
}

// InstallNpmDependencies creates a stage for installing npm dependencies in the container.
func InstallNpmDependencies() pipeline.Stage {
	return pipeline.Stage{
		Name: "Install Dependencies",
		Execute: func(config *pipeline.PipelineConfig) error {
			slog.Debug("In stage : Install Dependencies")
			pipelineOptions := (*config).(*PipelineOptions)
			buildContainer := pipelineOptions.BuildContainer
			execOpts := container.DefaultDockerExecOptions().
				WithCommand("npm install")

			execOpts, err := execOpts.WithStreamOptions(container.DockerStreamOptions{
				IsStreamingEnabled: pipelineOptions.StreamConfig.Enabled,
				Channel:            pipelineOptions.StreamConfig.Output,
			})

			if err != nil {
				return err
			}
			out, err := buildContainer.ExecCmd(execOpts)
			if err != nil {
				return err
			}
			slog.Debug(out)
			slog.Debug("Dependencies installed successfully.")
			return nil
		},
	}
}

// NpmBuild creates a stage for building the project using npm in the container.
func NpmBuild() pipeline.Stage {
	return pipeline.Stage{
		Name: "Build Project",
		Execute: func(config *pipeline.PipelineConfig) error {
			slog.Debug("In stage : Build Project")
			pipelineOptions := (*config).(*PipelineOptions)
			buildContainer := pipelineOptions.BuildContainer

			execOpts := container.DefaultDockerExecOptions().
				WithCommand("npm run build")

			execOpts, err := execOpts.WithStreamOptions(container.DockerStreamOptions{
				IsStreamingEnabled: pipelineOptions.StreamConfig.Enabled,
				Channel:            pipelineOptions.StreamConfig.Output,
			})
			if err != nil {
				return err
			}
			out, err := buildContainer.ExecCmd(execOpts)
			if err != nil {
				return err
			}
			slog.Debug(out)
			slog.Debug("Build completed successfully.")
			return nil
		},
	}
}

// CopyBuildFilesFromContainer creates a stage to copy build files from the container.
func CopyBuildFilesFromContainer() pipeline.Stage {
	return pipeline.Stage{
		Name: "Copy Build Files From Container",
		Execute: func(config *pipeline.PipelineConfig) error {
			slog.Debug("In stage : Copy Build Files From Container")
			cfg := (*config).(*PipelineOptions)
			buildContainer := cfg.BuildContainer

			tarReader, err := buildContainer.CopyFromContainer("/dist")
			if err != nil {
				return err
			}

			buffer := &bytes.Buffer{}
			_, err = io.Copy(buffer, tarReader)
			if err != nil {
				return err
			}

			cfg.BuildOutputFilesData = BuildOutputFilesData{
				ProjectTarFile: buffer,
				DirName:        "dist",
			}

			slog.Debug("Build files copied successfully.")
			return nil
		},
	}
}

func UploadBuildFilesToS3() pipeline.Stage {
	return pipeline.Stage{
		Name: "Upload Build Files to S3",
		Execute: func(config *pipeline.PipelineConfig) error {
			slog.Debug("In stage: Upload Build Files to S3")

			// Cast to PipelineOptions
			cfg := (*config).(*PipelineOptions)

			// Ensure BuildOutputFilesData is not nil
			if cfg.BuildOutputFilesData.ProjectTarFile == nil {
				return fmt.Errorf("no build output files to upload")
			}

			// Read the tar archive
			tarReader := tar.NewReader(cfg.BuildOutputFilesData.ProjectTarFile)
			for {
				header, err := tarReader.Next()
				if err == io.EOF {
					break // End of tar archive
				}
				if err != nil {
					return fmt.Errorf("failed to read tar archive: %w", err)
				}

				// Process only regular files
				if header.Typeflag == tar.TypeReg {
					var buf bytes.Buffer
					if _, err := io.Copy(&buf, tarReader); err != nil {
						return fmt.Errorf("failed to read file %s: %w", header.Name, err)
					}

					// Remove the dist folder prefix from the file path
					relativePath := strings.TrimPrefix(header.Name, cfg.BuildOutputFilesData.DirName+"/")
					if relativePath == header.Name {
						// This ensures that files directly in the dist directory are not affected
						relativePath = header.Name
					}

					// Construct the S3 key using subdomain and relative file path
					key := fmt.Sprintf("%s/%s", cfg.DeploymentConfig.Subdomain, relativePath)
					slog.Debug(key)
					// Upload each file to S3
					err = cfg.Store.Put(context.Background(), &buf, cfg.ProjectStorageConfig.BuildFilesBucket, key)
					if err != nil {
						return fmt.Errorf("failed to upload file %s to S3: %w", key, err)
					}

					slog.Debug("Uploaded file: ", key)
				}
			}

			// Ensure index.html is available for SPA routing
			defaultIndexKey := fmt.Sprintf("%s/index.html", cfg.DeploymentConfig.Subdomain)
			_, indexCheckErr := cfg.Store.Get(context.Background(), cfg.ProjectStorageConfig.BuildFilesBucket, defaultIndexKey)
			if indexCheckErr != nil {
				return fmt.Errorf("index.html not found in uploaded files: %w", indexCheckErr)
			}

			slog.Debug("All build files uploaded successfully for SPA.")
			return nil
		},
	}
}

func UploadToCloudfront() pipeline.Stage {
	return pipeline.Stage{
		Name: "Upload to CloudFront",
		Execute: func(config *pipeline.PipelineConfig) error {
			slog.Debug("In stage: Upload to CloudFront")

			// Cast to PipelineOptions
			cfg := (*config).(*PipelineOptions)

			// Create CloudFront distribution
			cfDomain, err := cfg.CDNProvider.CreateDistribution(cfg.Context, cdn.DistributionConfig{
				OriginDomain:      fmt.Sprintf("%s.s3.%s.amazonaws.com", cfg.ProjectStorageConfig.BuildFilesBucket, cfg.DeploymentConfig.CloudFrontConfig.Region),
				DomainName:        cfg.DeploymentConfig.Subdomain,
				CertificateARN:    cfg.DeploymentConfig.CloudFrontConfig.CertificateARN,
				DefaultRootObject: "index.html",
				OAI:               cfg.DeploymentConfig.CloudFrontConfig.OAI,
			})
			if err != nil {
				return fmt.Errorf("failed to create CloudFront distribution: %w", err)
			}

			cfg.DeploymentConfig.CloudFrontConfig.OAI = cfDomain
			slog.Debug("CloudFront distribution created successfully. Domain: ", cfDomain)
			return nil
		},
	}
}

func CustomSubdomainRouting() pipeline.Stage {
	return pipeline.Stage{
		Name: "Custom Subdomain Routing",
		Execute: func(config *pipeline.PipelineConfig) error {
			slog.Debug("In stage: Custom Subdomain Routing")

			// Cast to PipelineOptions
			cfg := (*config).(*PipelineOptions)

			// Update Route 53 DNS record
			err := cfg.DNSProvider.CreateOrUpdateRecord(cfg.Context, cfg.DeploymentConfig.Subdomain, cfg.DeploymentConfig.CloudFrontConfig.OAI, "CNAME")
			if err != nil {
				return fmt.Errorf("failed to update DNS record: %w", err)
			}

			slog.Debug("Subdomain routing configured successfully.")
			return nil
		},
	}
}
