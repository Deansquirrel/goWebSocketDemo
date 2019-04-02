package object

import (
	"fmt"
	"github.com/gorilla/websocket"
)

import log "github.com/Deansquirrel/goToolLog"

type Client struct {
	id        string
	socket    *websocket.Conn
	ChReceive chan *Message
	ChSend    chan *Message
	ChClose   chan struct{}
	ChErr     chan *error
}

func NewClient(id string, socket *websocket.Conn) *Client {
	return &Client{
		id:     id,
		socket: socket,
		//接收
		ChReceive: make(chan *Message),
		//发送
		ChSend:  make(chan *Message),
		ChClose: make(chan struct{}),
		ChErr:   make(chan *error),
	}
}

func (c *Client) Start(r func(), s func(), close func(), err func()) {
	go c.read()
	go c.write()
}

func (c *Client) GetId() string {
	return c.id
}

func (c *Client) Close() {
	_ = c.socket.Close()
	close(c.ChReceive)
	close(c.ChSend)
	close(c.ChClose)
	close(c.ChErr)
}

func (c *Client) read() {
	log.Debug(fmt.Sprintf("read start,id:%s", c.id))
	defer func() {
		log.Debug(fmt.Sprintf("read exit,id:%s", c.id))
	}()
	for {
		t, d, err := c.socket.ReadMessage()
		if err != nil {
			if websocket.IsCloseError(err) {
				c.ChClose <- struct{}{}
				return
			} else {
				c.ChErr <- &err
			}
		}
		c.ChReceive <- &Message{MessageType: t, Data: d}
	}
}

func (c *Client) write() {
	log.Debug(fmt.Sprintf("write start,id:%s", c.id))
	defer func() {
		log.Debug(fmt.Sprintf("write exit,id:%s", c.id))
	}()
	for {
		select {
		case msg, ok := <-c.ChSend:
			if !ok {
				_ = c.socket.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			_ = c.socket.WriteMessage(msg.MessageType, msg.Data)
		}
	}
}
