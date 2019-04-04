package worker

import "github.com/Deansquirrel/goWebSocketDemo/object"

type IClientManager interface {
	GetChRegister() chan<- IClient
	GetChUnregister() chan<- string
	GetChBroadcast() chan<- *object.SocketMessage
	GetClient(id string) IClient
}
