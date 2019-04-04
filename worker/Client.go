package worker

import (
	"context"
	"fmt"
	"github.com/Deansquirrel/goWebSocketDemo/object"
	"github.com/gorilla/websocket"
	"sync"
)

import log "github.com/Deansquirrel/goToolLog"

type client struct {
	id        string
	socket    *websocket.Conn
	chReceive chan *object.SocketMessage
	chSend    chan *object.SocketMessage
	lock      sync.Mutex

	ctx    context.Context
	cancel context.CancelFunc
}

func NewClient(id string, socket *websocket.Conn) *client {
	ctx, cancel := context.WithCancel(context.Background())
	c := client{
		id:        id,
		socket:    socket,
		chReceive: make(chan *object.SocketMessage),
		chSend:    make(chan *object.SocketMessage),
		ctx:       ctx,
		cancel:    cancel,
	}
	c.start()
	return &c
}

func (c *client) GetChReceive() <-chan *object.SocketMessage {
	return c.chReceive
}

func (c *client) GetChSend() chan<- *object.SocketMessage {
	return c.chSend
}

func (c *client) GetChClose() <-chan struct{} {
	return c.ctx.Done()
}

func (c *client) start() {
	go c.read()
	go c.write()
}

func (c *client) GetId() string {
	return c.id
}

func (c *client) SetId(id string) {
	c.id = id
}

func (c *client) Close() {
	c.cancel()
	_ = c.socket.Close()
	close(c.chReceive)
	close(c.chSend)
}

func (c *client) read() {
	log.Debug(fmt.Sprintf("Client read start,id:%s", c.id))
	defer log.Debug(fmt.Sprintf("Client read exit,id:%s", c.id))
	for {
		t, d, err := c.socket.ReadMessage()
		if err != nil {
			log.Error(fmt.Sprintf("Read error:%s,ClientID:%s", err.Error(), c.GetId()))
			c.Close()
			//if websocket.IsCloseError(err) {
			//	c.Close()
			//} else {
			//	log.Error(fmt.Sprintf("Read error:%s,ClientID:%s", err.Error(), c.GetId()))
			//
			//}
			return
		}
		m := &object.SocketMessage{
			ClientId:    c.id,
			MessageType: t,
			Data:        d,
		}
		c.chReceive <- m
	}
}

func (c *client) write() {
	log.Debug(fmt.Sprintf("Client write start,id:%s", c.id))
	defer log.Debug(fmt.Sprintf("Client write exit,id:%s", c.id))
	for {
		select {
		case msg, ok := <-c.chSend:
			if ok {
				c.writeData(msg)
			}
		case <-c.ctx.Done():
			return
		}
	}
}

func (c *client) writeData(msg *object.SocketMessage) {
	log.Debug(fmt.Sprintf("Client write data start,id:%s,dateLength:%d", c.id, len(msg.Data)))
	defer log.Debug(fmt.Sprintf("Client write data exit,id:%s", c.id))
	err := c.socket.WriteMessage(msg.MessageType, msg.Data)
	if err != nil {
		log.Error(err.Error())
	}
}
