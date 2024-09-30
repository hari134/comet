package stream

import (
	"context"
	"log"

	"github.com/hari134/comet/core/transport"
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
		payload := transport.NewPayload()
		payload.SetData("data",data.Data)
		event := transport.NewEvent("builder.stream",correlationID,payload)
		err := sm.Sender.Send(event)
		if err != nil{
			log.Printf("Failed to send stream with correlationID : %s",correlationID.ToString())
		}
	}
}