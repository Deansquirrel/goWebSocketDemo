package wsService

import (
	"fmt"
	"github.com/Deansquirrel/goToolCommon"
	"github.com/Deansquirrel/goWebSocketDemo/worker"
	"github.com/gorilla/websocket"
	"net/url"
	"time"
)

import log "github.com/Deansquirrel/goToolLog"

type clientService struct {
	client *worker.Client
}

func NewClientService() *clientService {
	return &clientService{}
}

func (s *clientService) Start() {
	if s.client != nil {
		s.client.Close()
	}
	s.client = nil
	u := url.URL{Scheme: "ws", Host: "127.0.0.1:1234", Path: "/websocket"}
	var dialer = &websocket.Dialer{
		HandshakeTimeout: 30 * time.Second,
	}
	conn, _, err := dialer.Dial(u.String(), nil)
	if err != nil {
		log.Error(fmt.Sprintf("WebSocket Dial error: %s", err.Error()))
		return
	}
	s.client = worker.NewClient(goToolCommon.Guid(), conn)
	go s.msgHandler()
}

func (s *clientService) msgHandler() {
	for {
		select {
		case msg := <-s.client.ChReceive:
			log.Debug(fmt.Sprintf("ID:%s,Type:%d,Msg:%s", msg.ClientId, msg.MessageType, msg.Data))
		}
	}
}
