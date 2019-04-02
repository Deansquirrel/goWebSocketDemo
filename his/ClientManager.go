package object

import (
	"fmt"
	log "github.com/Deansquirrel/goToolLog"
	"sync"
)

type clientManager struct {
	clients    map[string]*Client
	broadcast  chan *Message
	register   chan *Client
	unregister chan *Client

	lock sync.Mutex
}

func NewClientManager() *clientManager {
	return &clientManager{
		clients:    make(map[string]*Client, 0),
		broadcast:  make(chan *Message),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

func (manager *clientManager) Start() {
	log.Info("clientManager Start")
	defer func() {
		log.Info("clientManager Exit")
	}()
	for {
		select {
		case conn := <-manager.register:
			manager.addClient(conn)
		case conn := <-manager.unregister:
			manager.delClient(conn)
		case msg := <-manager.broadcast:
			manager.broad(msg)
		}
	}
}

func (manager *clientManager) Register(c *Client) {
	manager.register <- c
}

func (manager *clientManager) Unregister(c *Client) {
	manager.unregister <- c
}

func (manager *clientManager) Broadcast(msg *Message) {
	manager.broadcast <- msg
}

func (manager *clientManager) Send(id string, msg *Message) {
	c, ok := manager.clients[id]
	if !ok {
		log.Error(fmt.Sprintf("Client[%s]is not exists", id))
		return
	}
	c.ChSend <- msg
}

func (manager *clientManager) Close() {
	close(manager.broadcast)
	close(manager.register)
	close(manager.unregister)
}

func (manager *clientManager) addClient(c *Client) {
	manager.lock.Lock()
	defer manager.lock.Unlock()
	manager.clients[c.id] = c
}

func (manager *clientManager) delClient(c *Client) {
	manager.lock.Lock()
	defer manager.lock.Unlock()
	_, ok := manager.clients[c.id]
	if ok {
		c.Close()
		delete(manager.clients, c.id)
	}
}

func (manager *clientManager) broad(msg *Message) {
	go func() {
		for _, conn := range manager.clients {
			select {
			case conn.ChSend <- msg:
			default:
				manager.delClient(conn)
			}
		}
	}()
}
