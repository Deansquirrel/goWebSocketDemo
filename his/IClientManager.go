package object

type IClientManager interface {
	Start()
	Register(c *client)
	Unregister(c *client)
	Broadcast(msg *Message)
	Send(id string, msg *Message)
	Close()
}
