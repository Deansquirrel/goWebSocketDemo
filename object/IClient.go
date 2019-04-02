package object

type IClient interface {
	GetId() string
	Start(receiveFunc func(*Message), errFunc func(*error))
	Send(msg *Message)
	Close()
}
