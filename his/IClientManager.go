package his

type IClientManager interface {
	Start()
	Register(c IClient)
	Unregister(c IClient)
	Broadcast(msg *Message)
	Send(id string, msg *Message)
	Close()
}
