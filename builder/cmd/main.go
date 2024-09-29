package main

import (
	"log"
	"os"
	"strconv"

	"github.com/docker/docker/client"
	"github.com/hari134/comet/builder/container"
	"github.com/hari134/comet/builder/pipeline/pipelines"
	"github.com/hari134/comet/core/storage"
	"github.com/hari134/comet/core/transport"
	"github.com/joho/godotenv"
)

func init() {
	pipelines.InitializePipelines()
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Initialize dependencies
	capacity,err := strconv.Atoi(os.Getenv("capacity"))
	if err != nil{
		log.Fatal(err)
	}
	dockerClient, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		log.Fatal(err)
	}
	containerManager := container.NewDockerContainerManager().WithCapacity(capacity).WithClient(dockerClient)

	// Create AWS credentials from environment variables
	awsCreds := storage.AWSCredentials{
		AccessKey:       os.Getenv("AWS_ACCESS_KEY_ID"),
		SecretAccessKey: os.Getenv("AWS_SECRET_ACCESS_KEY"),
		BucketName:      os.Getenv("AWS_BUCKET_NAME"),
		Region:          os.Getenv("AWS_REGION"),
	}
	store ,err:= storage.NewS3Store(awsCreds)
	if err != nil{
		log.Fatal(err)
	}

	// Start the receiver in a goroutine
	receiver := transport.NewRestReceiver().WithEndpoint(":8080")

	eventHandler := transport.NewRestReceiverEventHandler().
		WithContainerManager(containerManager).
		WithStorage(store)

	go func() {
		log.Println("Starting receiver on port 8080...")
		err := receiver.StartReceiving(eventHandler, transport.Event{})
		if err != nil {
			log.Fatalf("Receiver failed to start: %v", err)
		}
	}()
}
