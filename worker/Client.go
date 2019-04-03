package worker

import (
	"context"
	"fmt"
	"github.com/Deansquirrel/goWebSocketDemo/object"
	"github.com/gorilla/websocket"
	"sync"
)

import log "github.com/Deansquirrel/goToolLog"

type Client struct {
	id        string
	socket    *websocket.Conn
	ChReceive chan *object.SocketMessage
	ChSend    chan *object.SocketMessage
	lock      sync.Mutex

	ctx    context.Context
	cancel context.CancelFunc
}

func NewClient(id string, socket *websocket.Conn) *Client {
	ctx, cancel := context.WithCancel(context.Background())
	return &Client{
		id:        id,
		socket:    socket,
		ChReceive: make(chan *object.SocketMessage),
		ChSend:    make(chan *object.SocketMessage),
		ctx:       ctx,
		cancel:    cancel,
	}
}

func (c *Client) Start() {
	go c.read()
	go c.write()
}

func (c *Client) GetId() string {
	return c.id
}

func (c *Client) Close() {
	c.cancel()
	_ = c.socket.Close()
	close(c.ChReceive)
	close(c.ChSend)
}

func (c *Client) read() {
	log.Debug(fmt.Sprintf("Client read start,id:%s", c.id))
	defer log.Debug(fmt.Sprintf("Client read exit,id:%s", c.id))
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
		m := &object.SocketMessage{
			ClientId:    c.id,
			MessageType: t,
			Data:        d,
		}
		c.ChReceive <- m
	}
}

func (c *Client) write() {
	log.Debug(fmt.Sprintf("Client write start,id:%s", c.id))
	defer log.Debug(fmt.Sprintf("Client write exit,id:%s", c.id))
	for {
		select {
		case msg := <-c.ChSend:
			c.writeData(msg)
		case <-c.ctx.Done():
			return
		}
	}
}

func (c *Client) writeData(msg *object.SocketMessage) {
	log.Debug(fmt.Sprintf("Client write data start,id:%s,dateLength:%d", c.id, len(msg.Data)))
	defer log.Debug(fmt.Sprintf("Client write data exit,id:%s", c.id))
	err := c.socket.WriteMessage(msg.MessageType, msg.Data)
	if err != nil {
		log.Error(err.Error())
	}
}
