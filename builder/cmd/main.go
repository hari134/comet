package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/docker/docker/client"
	"github.com/hari134/comet/builder/container"
	"github.com/hari134/comet/builder/pipeline/pipelines"
	"github.com/hari134/comet/builder/transport"
	"github.com/hari134/comet/core/storage"
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
	capacity := os.Getenv("capacity")
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

	// Prepare the sender to send an event to the receiver
	sender := transport.RestSender{
		Endpoint: "http://localhost:8080/",
	}

	// Create an event
	event := transport.Event{
		CorrelationID: "12345",
		Type:          "project.uploaded",
		Payload: transport.EventPayload{
			Data: map[string]interface{}{
				"buildType": "ReactViteBuilder",
			},
		},
	}

	// Send the event
	err := sender.Send(event)
	if err != nil {
		log.Fatalf("Failed to send event: %v", err)
	}

	// Wait for a while before stopping the receiver
	time.Sleep(5 * time.Second)

	// Stop the receiver
	err = receiver.StopReceiving()
	if err != nil {
		log.Fatalf("Failed to stop receiver: %v", err)
	}

	fmt.Println("Receiver stopped")
}
