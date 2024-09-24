package stream


import (
	"github.com/hari134/comet/builder/transport"
)

type Stream struct {
	CorrelationID transport.CorrelationID
	Data          string
}

func NewStream(correlationID transport.CorrelationID, data string) Stream {
	return Stream{
		correlationID,
		data,
	}
}
