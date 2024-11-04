package transport

import (
	"errors"

	"github.com/google/uuid"
)

type CorrelationID uuid.UUID

func (c CorrelationID) ToString() string {
	return uuid.UUID(c).String()
}

type Event struct {
	Type          string
	CorrelationID CorrelationID
	Payload       Payload
}

type Payload struct {
	Data map[string]interface{}
}

func NewPayload() Payload {
	return Payload{}
}

func (p Payload) SetData(key string, value interface{}) {
	p.Data[key] = value
}

func (p Payload) GetData(key string) (interface{}, error) {
	val, ok := p.Data[key]
	if !ok {
		return nil, errors.New("value not found for corresponsing key")
	}
	return val, nil
}

func NewEvent(eventType string, correlationID CorrelationID, payload Payload) Event {
	return Event{
		eventType,
		correlationID,
		payload,
	}
}

type TransportError struct {
	Message string
	Err     error
}

func (t *TransportError) Error() string {
	if t.Err != nil {
		return t.Message + ": " + t.Err.Error()
	}
	return t.Message
}

// NewTransportError creates a new structured error.
func NewTransportError(message string, err error) error {
	return &TransportError{
		Message: message,
		Err:     err,
	}
}

// Receiver defines the behavior for receiving and processing incoming events.
type Receiver interface {
	// StartReceiving listens for events and executes the callback for each incoming event.
	StartReceiving(eventHandler EventHandler) error
}

// Sender defines the behavior for sending events to another system.
type Sender interface {
	// Send sends an event along with its correlation ID and metadata.
	Send(event Event) error
}
