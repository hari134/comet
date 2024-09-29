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
	"github.com/hari134/comet/core/storage"
	"github.com/hari134/comet/core/transport"
)

// RestSender implements the Sender interface, allowing events to be sent via REST.
type RestSender struct {
	Endpoint string
}

// Send sends an event via an HTTP POST request to the specified REST endpoint.
func (r *RestSender) Send(event transport.Event) error {
	// Serialize the event into JSON
	eventJSON, err := json.Marshal(event)
	if err != nil {
		return transport.NewTransportError("failed to marshal event", err)
	}

	resp, err := http.Post(r.Endpoint, "application/json", bytes.NewBuffer(eventJSON))
	if err != nil {
		return transport.NewTransportError("failed to send event", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return transport.NewTransportError(fmt.Sprintf("failed to send event, status code: %d", resp.StatusCode), errors.New(string(body)))
	}

	return nil
}

// RestReceiver implements the Receiver interface, allowing events to be received via REST.
type RestReceiver struct {
	Endpoint string       // REST endpoint where this service listens for incoming events
	server   *http.Server // HTTP server for handling incoming requests
}

func NewRestReceiver() *RestReceiver{
	return &RestReceiver{}
}

func (r *RestReceiver) WithEndpoint(endpoint string) *RestReceiver{
	r.Endpoint = endpoint
	return r
}

// StartReceiving listens for incoming events on the specified endpoint and executes the provided handler.
func (r *RestReceiver) StartReceiving(eventHandler transport.EventHandler,event transport.Event) error {
	if r.Endpoint == "" {
		return transport.NewTransportError("no endpoint provided for receiving events", nil)
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

		var event transport.Event
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
		return transport.NewTransportError("failed to start REST receiver", err)
	}

	return nil
}

// StopReceiver stops the HTTP server for the RestReceiver.
func (r *RestReceiver) StopReceiving() error {
	if r.server != nil {
		err := r.server.Close()
		if err != nil {
			return transport.NewTransportError("failed to stop REST receiver", err)
		}
	}
	return nil
}

type RestReceiverEventHandler struct {
	containerManager container.ContainerManager
	store storage.Store
}

func NewRestReceiverEventHandler() *RestReceiverEventHandler {
	return &RestReceiverEventHandler{}
}

func (restReceiverEH *RestReceiverEventHandler) WithContainerManager(containerManager container.ContainerManager) *RestReceiverEventHandler{
	restReceiverEH.containerManager = containerManager
	return restReceiverEH
}

func (restReceiverEH *RestReceiverEventHandler) WithStorage(store storage.Store) *RestReceiverEventHandler{
	restReceiverEH.store = store
	return restReceiverEH
}

func (rh *RestReceiverEventHandler) HandleEvent(event transport.Event) error {
	correlationId := event.CorrelationID
	payload := event.Payload
	eventType := event.Type

	switch eventType {
	case "project.uploaded":
		buildType, err := payload.Get
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
