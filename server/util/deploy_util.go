package util

import (
	"bytes"
	"io"
	"github.com/docker/docker/pkg/namesgenerator"

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