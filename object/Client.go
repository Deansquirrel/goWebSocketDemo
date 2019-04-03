package object

import (
	"context"
	"fmt"
	"github.com/gorilla/websocket"
	"sync"
)

import log "github.com/Deansquirrel/goToolLog"

type client struct {
	id        string
	socket    *websocket.Conn
	ChReceive chan *Message
	ChSend    chan *Message
	lock      sync.Mutex

	ctx    context.Context
	cancel context.CancelFunc
}

func NewClient(id string, socket *websocket.Conn) *client {
	ctx, cancel := context.WithCancel(context.Background())
	return &client{
		id:        id,
		socket:    socket,
		ChReceive: make(chan *Message),
		ChSend:    make(chan *Message),
		ctx:       ctx,
		cancel:    cancel,
	}
}

func (c *client) GetId() string {
	return c.id
}

func (c *client) Start() {
	go c.read()
	go c.write()
}

func (c *client) read() {
	log.Debug(fmt.Sprintf("read start,id:%s", c.id))
	defer func() {
		log.Debug(fmt.Sprintf("read exit,id:%s", c.id))
	}()
	for {
		t, d, err := c.socket.ReadMessage()
		if err != nil {
			if websocket.IsCloseError(err) {
				c.Close()
			} else {
				log.Error(fmt.Sprintf("Read error:%s,ClientID:%s", err.Error(), c.GetId()))
			}
			return
		}
		c.ChReceive <- &Message{t, d}
	}
}

func (c *client) write() {
	log.Debug(fmt.Sprintf("write start,id:%s", c.id))
	defer func() {
		log.Debug(fmt.Sprintf("write exit,id:%s", c.id))
	}()
	for {
		select {
		case m := <-c.ChSend:
			c.writeWorker(m)
		case <-c.ctx.Done():
			return
		}
	}
}

func (c *client) writeWorker(msg *Message) {
	log.Debug(fmt.Sprintf("write msg start,id:%s,msg:%s", c.id, msg.Data))
	defer func() {
		log.Debug(fmt.Sprintf("write msg exit,id:%s", c.id))
	}()
	_ = c.socket.WriteMessage(msg.MessageType, msg.Data)
}

func (c *client) Close() {
	_ = c.socket.Close()
}
