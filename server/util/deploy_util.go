package util

import (
	"bytes"
	"github.com/docker/docker/pkg/namesgenerator"
	"io"
	"log"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func GetFileBytesBuffer(c *fiber.Ctx, fileKey string) (*bytes.Buffer, error) {
	file, err := c.FormFile(fileKey)
	if err != nil {
		return nil, err
	}

	fileReader, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer fileReader.Close()

	buf := new(bytes.Buffer)
	if _, err := io.Copy(buf, fileReader); err != nil {
		return nil, err
	}

	return buf, nil
}

func GetRandomName() string {
	return namesgenerator.GetRandomName(0)
}

func GetSubdomain(hostname string) string {
	log.Println(hostname)
	parts := strings.Split(hostname, ".")
	if len(parts) < 2 {
		return ""
	}
	return parts[0]
}

func GetContentType(key string) string {
	switch {
	case strings.HasSuffix(key, ".html"):
		return "text/html"
	case strings.HasSuffix(key, ".css"):
		return "text/css"
	case strings.HasSuffix(key, ".js"):
		return "application/javascript"
	case strings.HasSuffix(key, ".svg"):
		return "image/svg+xml"
	case strings.HasSuffix(key, ".png"):
		return "image/png"
	case strings.HasSuffix(key, ".jpg"), strings.HasSuffix(key, ".jpeg"):
		return "image/jpeg"
	case strings.HasSuffix(key, ".gif"):
		return "image/gif"
	case strings.HasSuffix(key, ".ico"):
		return "image/x-icon"
	case strings.HasSuffix(key, ".json"):
		return "application/json"
	default:
		return "application/octet-stream"
	}
}
