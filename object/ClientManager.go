package object

import (
	"fmt"
	"sync"
)

import log "github.com/Deansquirrel/goToolLog"

var clients map[string]*client
var ChBroadcast chan *Message
var ChRegister chan *client
var ChUnregister chan string
var lock sync.Mutex

func init() {
	clients = make(map[string]*client)
	ChBroadcast = make(chan *Message)
	ChRegister = make(chan *client)
	ChUnregister = make(chan string)
	c := clientManager{}
	c.start()
}

type clientManager struct {
}

func (manager *clientManager) start() {
	go func() {
		for {
			select {
			case c := <-ChRegister:
				manager.register(c)
			case c := <-ChUnregister:
				manager.unregister(c)
			case m := <-ChBroadcast:
				manager.broad(m)
			}
		}
	}()
}

func (manager *clientManager) register(c *client) {
	lock.Lock()
	defer lock.Unlock()
	log.Info(fmt.Sprintf("Client Register: %s，CurrClientNum: %d", c.GetId(), len(clients)))
	clients[c.GetId()] = c
}

func (manager *clientManager) unregister(id string) {
	lock.Lock()
	defer lock.Unlock()
	log.Info(fmt.Sprintf("Client Unregister: %s", id))
	c, ok := clients[id]
	if ok {
		c.Close()
		delete(clients, c.GetId())
	}
}

func (manager *clientManager) broad(msg *Message) {
	go func() {
		for _, c := range clients {
			c.ChSend <- msg
		}
	}()
}
