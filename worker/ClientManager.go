package worker

import (
	"fmt"
	"github.com/Deansquirrel/goWebSocketDemo/object"
	"sync"
)

import log "github.com/Deansquirrel/goToolLog"

type clientManager struct {
	clients      map[string]*Client
	chRegister   chan *Client
	chUnregister chan string
	chBroadcast  chan *object.SocketMessage
	lock         sync.Mutex
}

func NewClientManager() *clientManager {
	return &clientManager{
		clients:      make(map[string]*Client),
		chRegister:   make(chan *Client),
		chUnregister: make(chan string),
		chBroadcast:  make(chan *object.SocketMessage),
	}
}

func (manager *clientManager) GetChRegister() chan *Client {
	return manager.chRegister
}

func (manager *clientManager) GetChUnregister() chan string {
	return manager.chUnregister
}

func (manager *clientManager) GetChBroadcast() chan *object.SocketMessage {
	return manager.chBroadcast
}

func (manager *clientManager) Start() {
	go func() {
		for {
			select {
			case c := <-manager.chRegister:
				manager.register(c)
			case c := <-manager.chUnregister:
				manager.unregister(c)
			case m := <-manager.chBroadcast:
				manager.broad(m)
			}
		}
	}()
}

func (manager *clientManager) register(c *Client) {
	manager.lock.Lock()
	defer manager.lock.Unlock()
	log.Info(fmt.Sprintf("Client Register: %s，CurrClientNum: %d", c.GetId(), len(manager.clients)))
	manager.clients[c.GetId()] = c
	log.Info(fmt.Sprintf("Client Register: %s，CurrClientNum: %d", c.GetId(), len(manager.clients)))
}

func (manager *clientManager) unregister(id string) {
	manager.lock.Lock()
	defer manager.lock.Unlock()
	log.Info(fmt.Sprintf("Client Unregister: %s", id))
	c, ok := manager.clients[id]
	if ok {
		c.Close()
		delete(manager.clients, c.GetId())
	}
}

func (manager *clientManager) broad(msg *object.SocketMessage) {
	go func() {
		for _, c := range manager.clients {
			log.Debug(fmt.Sprintf("Board msg: Type:%d,Msg:%s", msg.MessageType, msg.Data))
			c.ChSend <- msg
		}
	}()
}
