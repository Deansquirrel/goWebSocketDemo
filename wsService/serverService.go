package wsService

import (
	"encoding/json"
	"fmt"
	"github.com/Deansquirrel/goToolCommon"
	"github.com/Deansquirrel/goWebSocketDemo/object"
	"github.com/Deansquirrel/goWebSocketDemo/worker"
	"github.com/gorilla/websocket"
	"github.com/kataras/iris/core/errors"
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
	//hm := object.SocketMessage{
	//	MessageType: websocket.TextMessage,
	//	Data:        []byte(fmt.Sprintf("Hello %s",c.GetId())),
	//}
	//c.GetChSend() <- &hm
	//bm := object.SocketMessage{
	//	MessageType: websocket.TextMessage,
	//	Data:        []byte(fmt.Sprintf("New Client[%s] Connect",c.GetId())),
	//}
	//s.manager.GetChBroadcast() <- &bm
	s.manager.GetChRegister() <- c
	go func() {
		select {
		case <-c.GetChClose():
			s.manager.GetChUnregister() <- c.GetId()
		}
	}()
	go func() {
		select {
		case msg := <-c.GetChReceive():
			rm := s.msgHandler(msg)
			if rm != nil {
				c.GetChSend() <- rm
			}
		}
	}()
}

func (s *serverService) msgHandler(msg *object.SocketMessage) *object.SocketMessage {
	//-------------------------------------------------------------------------------------------
	if msg.MessageType != websocket.TextMessage {
		errMsg := fmt.Sprintf("Unexpected MessageType : %d", msg.MessageType)
		log.Error(errMsg)
		return s.getRMessage(msg.ClientId, -1, errMsg)
	}
	log.Debug(fmt.Sprintf("Client: %s , Message %s", msg.ClientId, msg.Data))
	var m object.Message
	err := json.Unmarshal(msg.Data, &m)
	if err != nil {
		errMsg := fmt.Sprintf("Convert Message error: %s", err.Error())
		log.Error(errMsg)
		return s.getRMessage(msg.ClientId, -1, errMsg)
	}
	//-------------------------------------------------------------------------------------------
	log.Debug(fmt.Sprintf("Client: %s", msg.ClientId))
	log.Debug(fmt.Sprintf("MessageClient: %s", m.Id))
	log.Debug(fmt.Sprintf("MessageKey: %s", m.Key))
	log.Debug(fmt.Sprintf("MessageData: %s", m.Data))
	//-------------------------------------------------------------------------------------------
	switch m.Key {
	case "hello":
		var h object.Hello
		err := json.Unmarshal([]byte(m.Data), &h)
		if err != nil {
			log.Error(err.Error())
			return s.getRMessage(msg.ClientId, -1, err.Error())
		}
		log.Info(fmt.Sprintf("Client %s say hello, data : %s", msg.ClientId, h.Msg))
		return s.getRMessage(msg.ClientId, 0, fmt.Sprintf("Hello %s", msg.ClientId))
	case "updateId":
		var u object.UpdateId
		err := json.Unmarshal(msg.Data, &u)
		if err != nil {
			log.Error(err.Error())
			return s.getRMessage(msg.ClientId, -1, err.Error())
		}
		c := s.manager.GetClient(u.Id)
		if c != nil {
			errMsg := fmt.Sprintf("Client %s is already exist", u.Id)
			log.Error(errMsg)
			return s.getRMessage(msg.ClientId, -1, errMsg)
		}
		err = s.updateClientId(msg.ClientId, u.Id)
		if err != nil {
			log.Error(err.Error())
			return s.getRMessage(msg.ClientId, -1, err.Error())
		}
		return s.getRMessage(u.Id, 0, "update success")
	default:
		errMsg := fmt.Sprintf("Unexpected Command Key : %s", m.Key)
		log.Error(errMsg)
		return s.getRMessage(msg.ClientId, -1, errMsg)
	}
	//-------------------------------------------------------------------------------------------
}

func (s *serverService) getRMessage(clientId string, code int, data string) *object.SocketMessage {
	m := object.ReturnMessage{
		ErrCode: code,
		ErrMsg:  data,
	}
	d, err := json.Marshal(m)
	if err != nil {
		log.Error(fmt.Sprintf("getRMessage error: %s", err.Error()))
		return nil
	}
	return &object.SocketMessage{
		ClientId:    clientId,
		MessageType: websocket.TextMessage,
		Data:        d,
	}
}

func (s *serverService) updateClientId(old, new string) error {
	c := s.manager.GetClient(old)
	if c != nil {
		s.manager.GetChUnregister() <- old
		c.SetId(new)
		s.manager.GetChRegister() <- c
		return nil
	} else {
		return errors.New(fmt.Sprintf("UpdateClientId Error:Client is not exist[%s]", old))
	}
}
