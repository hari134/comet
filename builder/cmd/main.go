package main

import (
	"log"
	"os"
	"strconv"

	"github.com/docker/docker/client"
	"github.com/hari134/comet/builder"
	"github.com/hari134/comet/builder/container"
	"github.com/hari134/comet/builder/modules"
	"github.com/hari134/comet/builder/transport"
	"github.com/hari134/comet/core/storage"
	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Initialize dependencies
	capacity, err := strconv.Atoi(os.Getenv("CONTAINER_CONCURRENCY"))
	if err != nil {
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
		Region:          os.Getenv("AWS_REGION"),
	}
	store, err := storage.NewS3Store(awsCreds)
	if err != nil {
		log.Fatal(err)
	}

	// Start the receiver in a goroutine
	receiver := transport.NewRestReceiver().WithEndpoint("127.0.0.1:8080")
	pipelineFactory := modules.NewFactory()
	builder := builder.NewBuilder(store, containerManager, pipelineFactory)
	eventHandler := transport.NewRestReceiverEventHandler().
		WithBuilder(builder)

	go func() {
		log.Println("Starting receiver on port 8080...")
		err := receiver.StartReceiving(eventHandler)
		if err != nil {
			log.Fatalf("Receiver failed to start: %v", err)
		}
	}()
	select {}
}
