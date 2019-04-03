package his

import (
	"fmt"
	log "github.com/Deansquirrel/goToolLog"
	"sync"
)

type clientManager struct {
	clients      map[string]IClient
	chBroadcast  chan *Message
	chRegister   chan IClient
	chUnregister chan IClient

	lock sync.Mutex
}

func NewClientManager() *clientManager {
	return &clientManager{
		clients:      make(map[string]IClient),
		chBroadcast:  make(chan *Message),
		chRegister:   make(chan IClient),
		chUnregister: make(chan IClient),
	}
}

func (manager *clientManager) Start() {
	log.Info("clientManager Start")
	defer log.Info("clientManager Exit")
	go func() {
		for {
			select {
			case c := <-manager.chRegister:
				log.Debug("R " + c.GetId())
				manager.register(c)
			case c := <-manager.chUnregister:
				manager.unregister(c)
			case m := <-manager.chBroadcast:
				manager.Broadcast(m)
			}
		}
	}()
}

func (manager *clientManager) Register(c IClient) {
	log.Debug("RR " + c.GetId())
	manager.chRegister <- c
}

func (manager *clientManager) Unregister(c IClient) {
	manager.chUnregister <- c
}

func (manager *clientManager) Broadcast(msg *Message) {
	manager.chBroadcast <- msg
}

func (manager *clientManager) Send(id string, msg *Message) {
	c, ok := manager.clients[id]
	if !ok {
		log.Error(fmt.Sprintf("Client[%s]is not exists", id))
		return
	}
	c.Send(msg)
}

func (manager *clientManager) Close() {
	close(manager.chBroadcast)
	close(manager.chRegister)
	close(manager.chUnregister)
}

func (manager *clientManager) register(c IClient) {
	manager.lock.Lock()
	defer manager.lock.Unlock()
	log.Info(fmt.Sprintf("Client Register: %s，CurrClientNum: %d", c.GetId(), len(manager.clients)))
	manager.clients[c.GetId()] = c
}

func (manager *clientManager) unregister(c IClient) {
	manager.lock.Lock()
	defer manager.lock.Unlock()
	log.Info(fmt.Sprintf("Client Unregister: %s", c.GetId()))
	_, ok := manager.clients[c.GetId()]
	if ok {
		c.Close()
		delete(manager.clients, c.GetId())
	}
}

func (manager *clientManager) broad(msg *Message) {
	go func() {
		for _, c := range manager.clients {
			c.Send(msg)
		}
	}()
}
