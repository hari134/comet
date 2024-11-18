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
}

func NewDeployHandler(builderService *builder.Builder,store storage.Store,projectFilesBucketName string) *DeployHandler {
	return &DeployHandler{
		builderService: builderService,
		store : store,
		storageConfig: DeployStorageConfig{
			ProjectFilesBucketName: projectFilesBucketName,
		},
	}
}

func (dh *DeployHandler) CreateDeployment(c *fiber.Ctx) error {
	projectTarFile, err := util.GetFileBytesBuffer(c, "file")
	if err != nil {
		return c.Status(400).SendString("project tar file not uploaded")
	}
	projectName := util.GetRandomName() + ".tar"
	slog.Debug(projectName)
	err = dh.store.Put(context.Background(), projectTarFile, dh.storageConfig.ProjectFilesBucketName, projectName)
	if err != nil {
		slog.Debug(err.Error())
		return c.Status(400).SendString("error uploading project tar file")
	}
	projectBuildType, err := deployment.DetectProjectType(projectTarFile)

	slog.Debug(fmt.Sprintf("detected type of project as %s",projectBuildType))

	if err != nil {
		slog.Debug(fmt.Sprintf("Detected type is %s", projectBuildType))
		return c.Status(400).SendString("error detect project type")
	}
	err = dh.builderService.DeployProject(projectName, dh.storageConfig.ProjectFilesBucketName, projectBuildType)
	if err != nil {
		slog.Debug(fmt.Sprintf("Error in deploying project in handler %s", err.Error()))
		return c.Status(400).SendString("error deploying project")
	}
	return c.Status(201).JSON(map[string]interface{}{
		"message": "successfully created project deployment",
	})
}
