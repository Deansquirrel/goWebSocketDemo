package worker

import (
	"fmt"
	"github.com/Deansquirrel/goWebSocketDemo/object"
	"sync"
)

import log "github.com/Deansquirrel/goToolLog"

type clientManager struct {
	clients      map[string]IClient
	chRegister   chan IClient
	chUnregister chan string
	chBroadcast  chan *object.SocketMessage
	lock         sync.Mutex
}

func NewClientManager() *clientManager {
	cm := clientManager{
		clients:      make(map[string]IClient),
		chRegister:   make(chan IClient),
		chUnregister: make(chan string),
		chBroadcast:  make(chan *object.SocketMessage),
	}
	cm.start()
	return &cm
}

func (manager *clientManager) GetClient(id string) IClient {
	client, ok := manager.clients[id]
	if ok {
		return client
	}
	return nil
}

func (manager *clientManager) GetChRegister() chan<- IClient {
	return manager.chRegister
}

func (manager *clientManager) GetChUnregister() chan<- string {
	return manager.chUnregister
}

func (manager *clientManager) GetChBroadcast() chan<- *object.SocketMessage {
	return manager.chBroadcast
}

func (manager *clientManager) start() {
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

func (manager *clientManager) register(c IClient) {
	manager.lock.Lock()
	defer manager.lock.Unlock()
	_, ok := manager.clients[c.GetId()]
	if ok {
		log.Error(fmt.Sprintf("ClientId %s is already exist", c.GetId()))
		return
	}
	manager.clients[c.GetId()] = c
	log.Info(fmt.Sprintf("Client Register: %s，CurrClientNum: %d", c.GetId(), len(manager.clients)))
	return
}

func (manager *clientManager) unregister(id string) {
	manager.lock.Lock()
	defer manager.lock.Unlock()
	log.Info(fmt.Sprintf("Client Unregister: %s", id))
	_, ok := manager.clients[id]
	if ok {
		delete(manager.clients, id)
	}
}

func (manager *clientManager) broad(msg *object.SocketMessage) {
	go func() {
		for _, c := range manager.clients {
			log.Debug(fmt.Sprintf("Board msg: Type:%d,Msg:%s", msg.MessageType, msg.Data))
			c.GetChSend() <- msg
		}
	}()
}
