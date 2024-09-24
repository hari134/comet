package stream

import (
	"context"
	"log"

	"github.com/hari134/comet/builder/transport"
)

type StreamManager struct{
	Sender transport.Sender
}

func NewStreamManager(sender transport.Sender) *StreamManager{
	return &StreamManager{sender}
}


func (sm *StreamManager) SendStream(ctx context.Context, dataChan <- chan Stream){
	correlationID := ctx.Value("correlationID").(transport.CorrelationID)

	for data := range dataChan{
		err := sm.Sender.Send(correlationID,data)
		if err != nil{
			log.Printf("Failed to send stream with correlationID : %s",correlationID.ToString())
		}
	}
}