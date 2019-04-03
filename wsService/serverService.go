package wsService

import (
	"fmt"
	"github.com/Deansquirrel/goToolCommon"
	"github.com/Deansquirrel/goWebSocketDemo/object"
	"github.com/Deansquirrel/goWebSocketDemo/service"
	"github.com/gorilla/websocket"
	"net/http"
)

import log "github.com/Deansquirrel/goToolLog"

type serverService struct {
}

func (s *serverService) Start() {
	log.Info("Starting application")
	http.HandleFunc("/websocket", s.wsPage)
	_ = http.ListenAndServe(":1234", nil)
}

func (s *serverService) wsPage(res http.ResponseWriter, req *http.Request) {
	conn, err := (&websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}).Upgrade(res, req, nil)
	if err != nil {
		http.NotFound(res, req)
		return
	}

	c := object.NewClient(goToolCommon.Guid(), conn)
	c.Start(s.msgHandler, s.errHandler)

	msg := fmt.Sprintf("hello, %s", c.GetId())
	hMessage := &object.Message{
		MessageType: websocket.TextMessage,
		Data:        []byte(msg),
	}
	service.ChRegister <- c
	service.ChBroadcast <- hMessage
}

func (s *serverService) msgHandler(msg *object.Message) {
	log.Info(fmt.Sprintf("Received msg: %s", msg.Data))
}

func (s *serverService) errHandler(err error) {
	log.Error(fmt.Sprintf("Received err: %s", err.Error()))

}
