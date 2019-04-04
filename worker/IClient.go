package worker

import "github.com/Deansquirrel/goWebSocketDemo/object"

type IClient interface {
	GetId() string
	SetId(id string)
	GetChReceive() <-chan *object.SocketMessage
	GetChSend() chan<- *object.SocketMessage
	GetChClose() <-chan struct{}
	Close()
}
