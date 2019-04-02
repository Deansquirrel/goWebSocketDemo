package object

import (
	"fmt"
	"github.com/gorilla/websocket"
)

import log "github.com/Deansquirrel/goToolLog"

type client struct {
	id        string
	socket    *websocket.Conn
	chReceive chan *Message
	chSend    chan *Message
	chErr     chan *error
}

func (c *client) GetId() string {
	return c.id
}

func (c *client) Start(receiveFunc func(*Message), errFunc func(*error)) {
	go c.handler(receiveFunc, errFunc)
}

func (c *client) Send(msg *Message) {
	c.chSend <- msg
}

func (c *client) Close() {
	c.chSend <- c.getCloseMsg()
	_ = c.socket.Close()
	close(c.chReceive)
	close(c.chSend)
	close(c.chErr)
}

func (c *client) handler(receiveFunc func(*Message), errFunc func(*error)) {
	for {
		select {
		case msg := <-c.chReceive:
			receiveFunc(msg)
		case err := <-c.chErr:
			errFunc(err)
		}
	}
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
				c.chErr <- &err
			}
			return
		}
		c.chReceive <- &Message{MessageType: t, Data: d}
	}
}

func (c *client) write(msg *Message) {
	log.Debug(fmt.Sprintf("write start,id:%s", c.id))
	defer func() {
		log.Debug(fmt.Sprintf("write exit,id:%s", c.id))
	}()
	select {
	case msg, ok := <-c.chSend:
		if !ok {
			_ = c.socket.WriteMessage(websocket.CloseMessage, []byte{})
			return
		}
		_ = c.socket.WriteMessage(msg.MessageType, msg.Data)
	}
}

func (c *client) getCloseMsg() *Message {
	return &Message{
		MessageType: websocket.CloseMessage,
		Data:        []byte{},
	}
}
