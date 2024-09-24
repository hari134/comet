package transport

import (
	"github.com/google/uuid"
)

type CorrelationID uuid.UUID

func (c CorrelationID) ToString() string{
	return uuid.UUID(c).String()
}

type Receiver interface{
	Receive(correlationID CorrelationID,data interface{}) error
}


type Sender interface{
	Send(correlationID CorrelationID, data interface{}) error
}