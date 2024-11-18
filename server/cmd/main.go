package main

import (
	"fmt"
	"log"
	"log/slog"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/hari134/comet/builder"
	"github.com/hari134/comet/builder/container"
	"github.com/hari134/comet/builder/modules"
	"github.com/hari134/comet/core/storage"
	"github.com/hari134/comet/server/rest/handlers"
	"github.com/joho/godotenv"
)

func initLogger(level slog.Level) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: level,
	}))

	slog.SetDefault(logger)
}

func main() {
	initLogger(slog.LevelDebug)

	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}

	// Initialize all dependenciegos

	// Storage dependency (S3)
	awsCreds := storage.AWSCredentials{
		AccessKey:       os.Getenv("AWS_ACCESS_KEY_ID"),
		SecretAccessKey: os.Getenv("AWS_SECRET_ACCESS_KEY"),
		Region:          os.Getenv("AWS_REGION"),
	}
	projectFilesBucketName := os.Getenv("AWS_PROJECT_FILES_BUCKET")
	s3Store, err := storage.NewS3Store(awsCreds)
	if err != nil {
		log.Fatal(err)
	}

	// Container manager
	containerManager := container.NewDockerContainerManager().
		WithCapacity(20).
		WithDefaultClient()

	// Pipeline factor
	pipelineFactory := modules.NewFactory()

	// Builder
	builderService := builder.NewBuilder(s3Store, containerManager, pipelineFactory)

	// Start http server
	app := fiber.New()

	deploymentHandler := handlers.NewDeployHandler(builderService, s3Store, projectFilesBucketName)

	deploymentRoutes := app.Group("/deployments")
	deploymentRoutes.Post("/create-deployment", deploymentHandler.CreateDeployment)

	port := 8080
	slog.Debug(fmt.Sprintf("Starting server on http://localhost:%d\n", port))
	log.Fatal(app.Listen(fmt.Sprintf(":%d", port)))
}
