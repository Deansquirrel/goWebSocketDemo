package service

import (
	"fmt"
	"github.com/Deansquirrel/goToolCommon"
	"github.com/Deansquirrel/goWebSocketDemo/global"
	"github.com/Deansquirrel/goWebSocketDemo/object"
	"github.com/gorilla/websocket"
	"net/http"
)

import log "github.com/Deansquirrel/goToolLog"

type serverService struct {
}

func (s *serverService) Start() {
	log.Info("Starting application")
	go global.ClientManager.Start()
	http.HandleFunc("/websocket", s.wsPage)
	_ = http.ListenAndServe(":1234", nil)
	//http.Handle("/websocket",websocket.Handler(s.echo))
	//err := http.ListenAndServe(":1234",nil)
	//if err != nil {
	//	log.Error("Start error:" + err.Error())
	//}
}

func (s *serverService) wsPage(res http.ResponseWriter, req *http.Request) {
	conn, err := (&websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}).Upgrade(res, req, nil)
	if err != nil {
		http.NotFound(res, req)
		return
	}

	c := object.NewClient(goToolCommon.Guid(), conn)
	global.ClientManager.Register(c)

	msg := fmt.Sprintf("new client %s", c.GetId())
	mData := &object.Message{
		MessageType: websocket.TextMessage,
		Data:        []byte(msg),
	}
	global.ClientManager.Broadcast(mData)
}

//func (s *serverService) echo (ws *websocket.Conn) {
//var err error
//for {
//	var reply string
//	err = websocket.Message.Receive(ws,&reply)
//	if err != nil {
//		log.Error(fmt.Sprintf("Reveived error: %s",err.Error()))
//		break
//	}
//	log.Debug(fmt.Sprintf("Reveived from client: %s", reply))
//	msg := fmt.Sprintf("Reveived: %s",reply)
//
//	err = websocket.Message.Send(ws,msg)
//	if err != nil {
//		log.Error(fmt.Sprintf("Send error: %s",err.Error()))
//		break
//	}
//}
//}
