package handlers

import (
	"context"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/hari134/comet/core/storage"
	"github.com/hari134/comet/server/util"
	"log/slog"
)

type ServeHandler struct {
	store         storage.Store
	storageConfig ServeStorageConfig
}

type ServeStorageConfig struct {
	BuildFilesBucketName string
}

func NewServeHandler(
	store storage.Store,
	buildFilesBucketName string,
) *ServeHandler {
	return &ServeHandler{
		store: store,
		storageConfig: ServeStorageConfig{
			BuildFilesBucketName: buildFilesBucketName,
		},
	}
}
func (dh *ServeHandler) ServeSPA(c *fiber.Ctx) error {
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
