package wsnotify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"

	gorillaWs "github.com/gorilla/websocket"
	"github.com/kataras/iris/v12/websocket"
	"github.com/kataras/neffos"
	"github.com/kataras/neffos/gorilla"
	"github.com/masami10/rush/services/dispatcherbus"
	"github.com/masami10/rush/services/httpd"
	"github.com/masami10/rush/utils"
	"github.com/pkg/errors"
	"go.uber.org/atomic"
)

const (
	WS_EVENT_TIGHTENING  = "tightening"
	WS_EVENT_SCANNER     = "scanner"
	WS_EVENT_INPUT       = "input"
	WS_EVENT_IO          = "io"
	WS_EVENT_MAINTENANCE = "maintenance"
	WS_EVENT_READER      = "reader"
	WS_EVENT_REPLY       = "reply"
	WS_EVENT_DEVICE      = "device"
	WS_EVENT_ORDER       = "order"
	WS_EVENT_SERVICE     = "service"

	WS_EVENT_ERROR    = "err"
	WS_EVENT_REGISTER = "register"
	WS_EVENT_NOTIFY   = "notify"
)

type BaseDiag interface {
	Debug(msg string)
}

type Diagnostic interface {
	BaseDiag
	Error(msg string, err error)
	Disconnect(id string)
	OnMessage(msg string)
	Close()
	Closed()
}

type Service struct {
	configValue   atomic.Value
	diag          Diagnostic
	ws            *neffos.Server
	serverEvents  neffos.Events
	httpd         HTTPService
	clientManager *WSClientManager
	dispatcherBus Dispatcher
}

func (s *Service) Config() Config {
	return s.configValue.Load().(Config)
}

func (s *Service) NewWebSocketRecvHandler(handler func(interface{})) {
	fn := utils.CreateDispatchHandlerStruct(handler)
	s.dispatcherBus.Register(dispatcherbus.DispatcherWsNotify, fn)
}

// const namespace = "default"

func (s *Service) onWebsocketConnect(nsConn *websocket.Conn) error {
	// with `websocket.GetContext` you can retrieve the Iris' `Context`.
	ctx := websocket.GetContext(nsConn)

	s.diag.Debug(fmt.Sprintf("[%s] connected with IP [%s]",
		nsConn,
		ctx.RemoteAddr()))
	return nil
}

func (s *Service) OnWebsocketDisconnect(nsConn *websocket.Conn) {
	errInfo := "No Error"
	s.diag.Disconnect(fmt.Sprintf("ClientID: %s ErrorInfo:%s", nsConn.String(), errInfo))
	s.diag.Debug(fmt.Sprintf("Rest WS Clients: %+v", s.ws.GetConnections()))
}

func (s *Service) OnMessage(nsConn *websocket.NSConn, msg websocket.Message) error {
	ctx := websocket.GetContext(nsConn.Conn)

	s.diag.Debug(fmt.Sprintf("Receive Message From: [%s], Body: %v", ctx.RemoteAddr(), string(msg.Body)))
	var message WSMsg
	err := json.Unmarshal(msg.Body, &message)
	if err != nil {
		s.diag.Error("WSMsg Payload Error", err)
		return err
	}

	if message.Type == WS_REG {
		s.handleRegister(&message, nsConn)
	}

	s.postNotify(&DispatcherNotifyPackage{
		C:    nsConn,
		Data: msg.Body,
	})

	return nil
}

func (s *Service) initServerHandler() {
	if s.serverEvents == nil {
		return
	}
	s.serverEvents.On(websocket.OnNativeMessage, s.OnMessage)
}

func NewService(c Config, d Diagnostic, dp Dispatcher, httpd HTTPService) *Service {
	//defaultPongTimeout := 1 * time.Second
	s := &Service{
		diag:          d,
		dispatcherBus: dp,
		serverEvents:  neffos.Events{},
		clientManager: &WSClientManager{
			diag: d,
		},
		httpd: httpd,
	}

	s.initServerHandler()

	s.ws = websocket.New(
		gorilla.Upgrader(gorillaWs.Upgrader{CheckOrigin: func(*http.Request) bool { return true }}),
		s.serverEvents)

	s.ws.OnConnect = s.onWebsocketConnect

	s.ws.OnDisconnect = s.OnWebsocketDisconnect

	s.clientManager.Init(s.ws)

	s.configValue.Store(c)

	return s

}

func (s *Service) handleRegister(msg *WSMsg, c *websocket.NSConn) {
	result := 0
	resultMsg := ""

	// 将客户端加入列表
	if err := s.clientManager.CheckClient(c); err != nil {
		result = -1
		resultMsg = err.Error()
		s.diag.Error(fmt.Sprintf("CheckClient Failed: %s", c.String()), err)
	}

	// 注册成功
	_ = WSClientSend(c, WS_EVENT_REPLY, GenerateReply(msg.SeqNumber, msg.Type, result, resultMsg), s.diag)
}

func (s *Service) createAndStartWebSocketNotifyDispatcher() error {
	if err := s.dispatcherBus.Create(dispatcherbus.DispatcherWsNotify, utils.DefaultDispatcherBufLen); err != nil {
		return err
	} else {
		return s.dispatcherBus.Start(dispatcherbus.DispatcherWsNotify)
	}
}

func (s *Service) postNotify(msg *DispatcherNotifyPackage) {
	s.diag.Debug(fmt.Sprintf("WS REQ: %s", string(msg.Data)))
	if err := s.dispatcherBus.Dispatch(dispatcherbus.DispatcherWsNotify, msg); err != nil {
		s.diag.Error("notify", err)
	}
}

func (s *Service) Open() error {

	c := s.Config()

	if err := s.createAndStartWebSocketNotifyDispatcher(); err != nil {
		s.diag.Error("createAndStartWebSocketNotifyDispatcher Error", err)
		return nil
	}

	r := httpd.WebsocketRoute{
		Pattern: c.Route,
		Server:  s.ws,
	}
	_ = s.httpd.AddNewWebsocketHandler(r)

	return nil
}

func (s *Service) Close() error {
	s.diag.Close()

	if s.clientManager != nil {
		s.clientManager.CloseAll()
		s.clientManager = nil
	}

	if s.ws != nil {
		s.ws.Close()
		s.ws = nil
	}

	s.diag.Closed()

	return nil
}

func (s *Service) NotifyAll(evt string, payload string) {
	if s == nil || s.clientManager == nil {
		return
	}
	s.clientManager.NotifyALL(evt, payload)
	s.diag.Debug(fmt.Sprintf("WS NOTIFY: %s", payload))
}

func GenerateReply(sn uint64, wsType string, result int, msg string) *WSMsg {
	return GenerateWSMsg(sn, wsType, WSReply{
		Result: result,
		Msg:    msg,
	})
}

func GenerateWSMsg(sn uint64, wsType string, data interface{}) *WSMsg {
	return &WSMsg{
		SeqNumber: sn,
		Type:      wsType,
		Data:      data,
	}
}

func WSConnSend(c *websocket.Conn, event string, payload interface{}, diag BaseDiag) error {
	if c == nil {
		return errors.New("conn is nil")
	}
	var bt bytes.Buffer
	bt.WriteString(fmt.Sprintf("event:%s;", event))

	t := reflect.TypeOf(payload)
	switch t.Kind() { //nolint exhaustive
	case reflect.String:
		bt.WriteString(payload.(string))
	case reflect.Struct:
		body, _ := json.Marshal(payload)
		bt.Write(body)
	default:
		body, _ := json.Marshal(payload)
		bt.Write(body)
	}

	ss := bt.String()
	diag.Debug(fmt.Sprintf("WS RESP: %s", ss))
	m := neffos.Message{Namespace: "", Event: event, Body: bt.Bytes()}
	m.IsNative = true // 发送原始数据包
	m.SetBinary = false
	success := c.Write(m)
	if !success {
		return errors.Errorf("Emit Event: %s Message Error", event)
	}
	return nil
}

func WSClientSend(c *websocket.NSConn, event string, payload interface{}, diag BaseDiag) error {
	if c == nil {
		return errors.New("conn is nil")
	}
	return WSConnSend(c.Conn, event, payload, diag)
}

func WSNotifyInfo(c *websocket.NSConn, payload string, diag BaseDiag) {
	data, _ := json.Marshal(WSMsg{
		Type: WS_NOTIFY,
		Data: map[string]interface{}{
			"message": payload,
			"variant": "Info",
			"config": map[string]interface{}{
				"autoHideDuration": 6000,
			},
		},
	})
	_ = WSClientSend(c, WS_EVENT_NOTIFY, string(data), diag)
}
