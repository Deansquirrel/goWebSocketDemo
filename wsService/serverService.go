package wsService

import (
	"github.com/Deansquirrel/goToolCommon"
	"github.com/Deansquirrel/goWebSocketDemo/object"
	"github.com/Deansquirrel/goWebSocketDemo/worker"
	"github.com/gorilla/websocket"
	"net/http"
)

import log "github.com/Deansquirrel/goToolLog"

type serverService struct {
	manager worker.IClientManager
}

func NewServerService() *serverService {
	return &serverService{}
}

func (s *serverService) Start() {
	log.Info("Starting application")

	s.manager = worker.NewClientManager()
	s.manager.Start()

	http.HandleFunc("/websocket", s.wsPage)
	_ = http.ListenAndServe(":1234", nil)
}

func (s *serverService) wsPage(res http.ResponseWriter, req *http.Request) {
	conn, err := (&websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}).Upgrade(res, req, nil)
	if err != nil {
		http.NotFound(res, req)
		return
	}

	c := worker.NewClient(goToolCommon.Guid(), conn)
	//c.Start()
	//s.manager.ChRegister <- c
	s.manager.GetChRegister() <- c

	m := object.SocketMessage{
		MessageType: websocket.TextMessage,
		Data:        []byte("OK"),
	}
	s.manager.GetChBroadcast() <- &m
}

//
//func (s *serverService) msgHandler(client *worker.Client) {
//	for{
//		select{
//		case msg := <- client.GetReceiveChan():
//			log.Debug(fmt.Sprintf("ID:%s,Type:%d,Msg:%s",msg.ClientId,msg.MessageType,msg.Data))
//		}
//	}
//}
