package handlers

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/gofiber/fiber/v2"
	"github.com/hari134/comet/builder"
	"github.com/hari134/comet/core/deployment"
	"github.com/hari134/comet/core/storage"
	"github.com/hari134/comet/server/util"
)

type DeployHandler struct {
	builderService *builder.Builder
	store          storage.Store
	storageConfig  DeployStorageConfig
}

type DeployStorageConfig struct {
	ProjectFilesBucketName string
	BuildFilesBucketName   string
}

func NewDeployHandler(builderService *builder.Builder, store storage.Store, projectFilesBucketName string, buildFilesBucketName string) *DeployHandler {
	return &DeployHandler{
		builderService: builderService,
		store:          store,
		storageConfig: DeployStorageConfig{
			ProjectFilesBucketName: projectFilesBucketName,
			BuildFilesBucketName:   buildFilesBucketName,
		},
	}
}

func (dh *DeployHandler) CreateDeployment(c *fiber.Ctx) error {
	projectTarFile, err := util.GetFileBytesBuffer(c, "file")
	if err != nil {
		return c.Status(400).SendString("project tar file not uploaded")
	}
	subdomain := util.GetRandomName()
	projectName := subdomain + ".tar"
	slog.Debug(fmt.Sprintf("deployment started for project : %s", projectName))
	err = dh.store.Put(context.Background(), projectTarFile, dh.storageConfig.ProjectFilesBucketName, projectName)
	if err != nil {
		slog.Debug(err.Error())
		return c.Status(400).SendString("error uploading project tar file")
	}
	projectBuildType, err := deployment.DetectProjectType(projectTarFile)

	if err != nil {
		slog.Debug(fmt.Sprintf("Detected type is %s", projectBuildType))
		return c.Status(400).SendString("error detect project type")
	}
	slog.Debug(fmt.Sprintf("detected type of project as %s", projectBuildType))

	projectDeploymentConfig := builder.ProjectDeploymentConfig{
		ProjectStorageKey:    projectName,
		ProjectStorageBucket: dh.storageConfig.ProjectFilesBucketName,
		BuildFilesBucket:     dh.storageConfig.BuildFilesBucketName,
		BuildEnvType:         projectBuildType,
		OriginDomain:         "cometinfra.live",
		SubDomain:            subdomain,
	}

	err = dh.builderService.DeployProject(projectDeploymentConfig)
	if err != nil {
		slog.Debug(fmt.Sprintf("Error in deploying project in handler %s", err.Error()))
		return c.Status(400).SendString("error deploying project")
	}
	return c.Status(201).JSON(map[string]interface{}{
		"message":        "successfully created project deployment",
		"deployment_url": "http://" + subdomain + ".localhost:8080",
	})
}

func (dh *DeployHandler) ServeSPA(c *fiber.Ctx) error {
	subdomain := util.GetSubdomain(c.Hostname()) // Extract subdomain from the hostname
	if subdomain == "" {
		return c.Status(400).SendString("invalid subdomain")
	}

	requestedPath := c.Path()
	if requestedPath == "/" {
		requestedPath = "/index.html" // Default to index.html for the root path
	}

	s3Key := fmt.Sprintf("%s%s", subdomain, requestedPath)

	// Fetch the file from S3
	fileData, err := dh.store.Get(context.Background(), dh.storageConfig.BuildFilesBucketName, s3Key)
	if err != nil {
		// Fallback to index.html for SPA routing
		s3Key = fmt.Sprintf("%s/index.html", subdomain)
		fileData, err = dh.store.Get(context.Background(), dh.storageConfig.BuildFilesBucketName, s3Key)
		if err != nil {
			slog.Debug("File not found", "subdomain", subdomain, "path", requestedPath, "error", err.Error())
			return c.Status(404).SendString("file not found")
		}
	}

	// Set the appropriate Content-Type header
	contentType := util.GetContentType(s3Key)
	c.Set("Content-Type", contentType)

	// Serve the file content
	return c.SendStream(fileData)
}
