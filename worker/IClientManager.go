package worker

import "github.com/Deansquirrel/goWebSocketDemo/object"

type IClientManager interface {
	Start()
	GetChRegister() chan *Client
	GetChUnregister() chan string
	GetChBroadcast() chan *object.SocketMessage
}
