package wsService

import (
	"fmt"
	"github.com/Deansquirrel/goToolCommon"
	log "github.com/Deansquirrel/goToolLog"
	"github.com/Deansquirrel/goWebSocketDemo/object"
	"github.com/gorilla/websocket"
	"github.com/kataras/iris/core/errors"
	"net/url"
	"os"
	"time"
)

type clientService struct {
	client object.IClient
}

func (s *clientService) Start() {
	err := s.restart()
	if err != nil {
		log.Error(err.Error())
	}
}

func (s *clientService) restart() error {
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
		return errors.New(fmt.Sprintf("WebSocket Dial error: %s", err.Error()))
	}
	c := object.NewClient(goToolCommon.Guid(), conn)
	c.Start(s.msgHandler, s.errHandler)
	s.client = c
	return nil
}

func (s *clientService) msgHandler(msg *object.Message) {
	log.Info(fmt.Sprintf("Received msg: %s", msg.Data))
}

func (s *clientService) errHandler(err error) {
	log.Error(fmt.Sprintf("Received err: %s", err.Error()))
	os.Exit(-1)
}
