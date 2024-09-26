package transport

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/hari134/comet/builder/container"
	"github.com/hari134/comet/builder/pipeline"
	"github.com/hari134/comet/builder/pipeline/pipelines"
)

// RestSender implements the Sender interface, allowing events to be sent via REST.
type RestSender struct {
	Endpoint string
}

// Send sends an event via an HTTP POST request to the specified REST endpoint.
func (r *RestSender) Send(event Event) error {
	// Serialize the event into JSON
	eventJSON, err := json.Marshal(event)
	if err != nil {
		return NewTransportError("failed to marshal event", err)
	}

	resp, err := http.Post(r.Endpoint, "application/json", bytes.NewBuffer(eventJSON))
	if err != nil {
		return NewTransportError("failed to send event", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return NewTransportError(fmt.Sprintf("failed to send event, status code: %d", resp.StatusCode), errors.New(string(body)))
	}

	return nil
}

// RestReceiver implements the Receiver interface, allowing events to be received via REST.
type RestReceiver struct {
	Endpoint string       // REST endpoint where this service listens for incoming events
	server   *http.Server // HTTP server for handling incoming requests
}

// StartReceiving listens for incoming events on the specified endpoint and executes the provided handler.
func (r *RestReceiver) StartReceiving(eventHandler EventHandler,event Event) error {
	if r.Endpoint == "" {
		return NewTransportError("no endpoint provided for receiving events", nil)
	}

	http.HandleFunc(r.Endpoint, func(w http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodPost {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
			return
		}

		// Read the request body
		body, err := io.ReadAll(req.Body)
		if err != nil {
			http.Error(w, "Failed to read request body", http.StatusInternalServerError)
			return
		}

		var event Event
		if err := json.Unmarshal(body, &event); err != nil {
			http.Error(w, "Failed to unmarshal event", http.StatusBadRequest)
			return
		}

		if err := eventHandler.HandleEvent(event); err != nil {
			http.Error(w, "Failed to process event", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	})

	// Start the HTTP server
	r.server = &http.Server{Addr: r.Endpoint}
	err := r.server.ListenAndServe()
	if err != nil {
		return NewTransportError("failed to start REST receiver", err)
	}

	return nil
}

// StopReceiver stops the HTTP server for the RestReceiver.
func (r *RestReceiver) StopReceiving() error {
	if r.server != nil {
		err := r.server.Close()
		if err != nil {
			return NewTransportError("failed to stop REST receiver", err)
		}
	}
	return nil
}

type RestReceiverEventHandler struct {
	containerManager container.ContainerManager
}

func NewRestReceiverEventHandler(containerManager container.ContainerManager) *RestReceiverEventHandler {
	return &RestReceiverEventHandler{containerManager}
}

func (rh *RestReceiverEventHandler) HandleEvent(event Event) error {
	correlationId := event.CorrelationID
	payload := event.Payload
	eventType := event.Type

	switch eventType {
	case "project.uploaded":
		buildType, err := payload.GetData("buildType")
		if err != nil {
			return err
		}

		buildPipeline, err := pipelines.PipelineFactory(buildType.(string))
		if err != nil {
			return err
		}

		buildContainer, err := rh.containerManager.NewBuildContainer(buildType.(string))
		if err != nil {
			return err
		}

		ctx := pipeline.NewPipelineContext().WithContainer(buildContainer)
		ctx.Set("correlationId", correlationId)

		err = buildPipeline.Run(ctx)
		if err != nil {
			return err
		}
		return nil
	default:
		return errors.New("invalid event type")
	}
}
