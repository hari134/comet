package transport

type EventHandler interface{
	HandleEvent(event Event) error
}