package wsService

import (
	"encoding/json"
	"fmt"
	"github.com/Deansquirrel/goToolCommon"
	"github.com/Deansquirrel/goWebSocketDemo/global"
	"github.com/Deansquirrel/goWebSocketDemo/object"
	"github.com/Deansquirrel/goWebSocketDemo/worker"
	"github.com/gorilla/websocket"
	"net/url"
	"time"
)

import log "github.com/Deansquirrel/goToolLog"

type clientService struct {
	client worker.IClient
}

func NewClientService() *clientService {
	return &clientService{}
}

func (s *clientService) Start() {
	s.client = nil
	u := url.URL{Scheme: "ws", Host: "127.0.0.1:1234", Path: "/websocket"}
	var dialer = &websocket.Dialer{
		HandshakeTimeout: 30 * time.Second,
	}
	conn, _, err := dialer.Dial(u.String(), nil)
	if err != nil {
		log.Error(fmt.Sprintf("WebSocket Dial error: %s", err.Error()))
		time.AfterFunc(global.ReConnectDuration, s.Start)
		return
	}
	s.client = worker.NewClient(goToolCommon.Guid(), conn)
	go s.msgHandler()

	time.AfterFunc(time.Second*8, func() {
		d := s.getUpdateId()
		if d != nil {
			s.client.GetChSend() <- d
		}
	})

	time.AfterFunc(time.Second*5, func() {
		d := s.getHelloData()
		if d != nil {
			s.client.GetChSend() <- d
		}
	})
}

func (s *clientService) msgHandler() {
	for {
		select {
		case msg, ok := <-s.client.GetChReceive():
			//log.Debug(fmt.Sprintf("ID:%s,Type:%d,Msg:%s", msg.ClientId, msg.MessageType, msg.Data))
			if ok {
				s.msgHandlerWorker(msg)
			}
		case <-s.client.GetChClose():
			time.AfterFunc(global.ReConnectDuration, s.Start)
		}
	}
}

func (s *clientService) msgHandlerWorker(msg *object.SocketMessage) {
	if msg.ClientId != s.client.GetId() {
		log.Error(fmt.Sprintf("Err Client,exp: %s,act %s", s.client.GetId(), msg.ClientId))
		s.printMsgData(msg)
	}
	if msg.MessageType != websocket.TextMessage {
		log.Error(fmt.Sprintf("Unexpected MessageType,exp: %d,act %d", websocket.TextMessage, msg.MessageType))
		s.printMsgData(msg)
	}
	var rm object.ReturnMessage
	err := json.Unmarshal(msg.Data, &rm)
	if err != nil {
		log.Error(fmt.Sprintf("Convert ReturnMessage error: %s", err.Error()))
		s.printMsgData(msg)
	} else {
		log.Info(fmt.Sprintf("code: %d,msg: %s", rm.ErrCode, rm.ErrMsg))
	}
}

func (s *clientService) printMsgData(msg *object.SocketMessage) {
	if len(msg.Data) < 1024 {
		log.Error(fmt.Sprintf("ClientId: %s,Type: %d,Data: %s", msg.ClientId, msg.MessageType, msg.Data))
	} else {
		log.Error(fmt.Sprintf("ClientId: %s,Type: %d, DataLength: %d,Data: %s",
			msg.ClientId,
			msg.MessageType,
			len(msg.Data),
			msg.Data[:1024]))
	}
}

func (s *clientService) getSocketMessage(key string, v interface{}) *object.SocketMessage {
	d, err := json.Marshal(v)
	if err != nil {
		log.Error(fmt.Sprintf("Get SocketMessage error : %s", err.Error()))
		return nil
	}
	log.Debug(string(d))
	m := object.Message{
		Id:   s.client.GetId(),
		Key:  key,
		Data: string(d),
	}
	sd, err := json.Marshal(m)
	if err != nil {
		log.Error(fmt.Sprintf("Get SocketMessage error : %s", err.Error()))
		return nil
	}
	log.Debug(string(sd))
	return &object.SocketMessage{
		ClientId:    s.client.GetId(),
		MessageType: websocket.TextMessage,
		Data:        sd,
	}
}

func (s *clientService) getHelloData() *object.SocketMessage {
	h := object.Hello{
		Msg: "Hello Server",
	}
	return s.getSocketMessage("hello", h)
}

func (s *clientService) getUpdateId() *object.SocketMessage {
	u := object.UpdateId{
		Id: s.client.GetId(),
	}
	return s.getSocketMessage("updateId", u)
}
