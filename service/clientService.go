package service

import (
	"github.com/Deansquirrel/goToolCommon"
	log "github.com/Deansquirrel/goToolLog"
	"github.com/Deansquirrel/goWebSocketDemo/object"
	"github.com/gorilla/websocket"
	"time"
)

type clientService struct {
	client *object.Client
}

func (s *clientService) Start() {
	log.Debug("Start")
	if s.client == nil {
		conn, err := s.getClient()
		if err != nil {
			log.Error(err.Error())
			return
		}
		s.client = object.NewClient(goToolCommon.Guid(), conn)
	}

}

func (s *clientService) getClient() (*websocket.Conn, error) {
	var dialer = &websocket.Dialer{
		HandshakeTimeout: time.Second * 30,
	}
	var url = "ws://127.0.0.1:1234/websocket"
	conn, _, err := dialer.Dial(url, nil)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func (s *clientService) startWork(c *object.Client) {

}

//func (s *clientService) sendMsg(msg string) {
//	var origin = "http://127.0.0.1:1234/"
//	var url = "ws://127.0.0.1:1234/websocket"
//	ws,err := websocket.Dial(url,"",origin)
//	if err != nil {
//		log.Error(fmt.Sprintf("Dial error:%s，msg:%s",err.Error(),msg))
//		return
//	}
//	defer func(){
//		log.Debug("Close")
//		_ = ws.Close()
//	}()
//	_,err = ws.Write([]byte(msg))
//	if err != nil {
//		log.Error(fmt.Sprintf("Write error:%s，msg:%s",err.Error(),msg))
//		return
//	}
//}
