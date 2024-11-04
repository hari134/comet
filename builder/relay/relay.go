package relay

type StreamData interface{
  GetData() string
  GetSource() string
}

type OutputRelay interface {
    Start() error
    Stop()
    Send(data StreamData)
    Receive() <-chan StreamData
    Close()
}

