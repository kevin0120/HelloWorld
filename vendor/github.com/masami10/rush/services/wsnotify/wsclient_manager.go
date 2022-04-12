package wsnotify

import (
	"context"
	"fmt"
	"time"

	"github.com/kataras/iris/v12/websocket"
	"github.com/kataras/neffos"
	"github.com/masami10/rush/utils"
	"github.com/pkg/errors"
)

var MaxClient = int(utils.GetEnvInt32("ENV_WEBSOCKET_MAX_CONNECTIONS", 1))

type WSClientManager struct {
	diag Diagnostic

	ws *neffos.Server
}

func (s *WSClientManager) Init(ws *neffos.Server) {
	s.ws = ws

	s.diag.Debug(fmt.Sprintf("Init Websocket Clients: %+v", s.ws.GetConnections()))
}

func (s *WSClientManager) CheckClient(c *websocket.NSConn) error {

	s.diag.Debug(fmt.Sprintf("Receive Conn Req: %s", c.String()))
	if !utils.CheckRuntimeEnvIsDev() {
		conns := s.ws.GetConnections()
		if len(conns) > MaxClient {
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			if err := c.Disconnect(ctx); err != nil {
				s.diag.Error(fmt.Sprintf("CheckClient Disconnect Failed: %s", c.String()), err)
			}
			return errors.New("Max Client Reached ")
		}
	}

	s.diag.Debug(fmt.Sprintf("Add WS Client Success: %s", c.String()))
	s.diag.Debug(fmt.Sprintf("Rest WS Clients: %+v", s.ws.GetConnections()))

	return nil
}

func (s *WSClientManager) NotifyALL(evt string, payload string) {
	if s.ws == nil {
		s.diag.Error("NotifyALL", errors.New("websocket Server Is Empty!!!"))
	}
	cons := s.ws.GetConnections()
	for _, v := range cons {
		if err := WSConnSend(v, evt, payload, s.diag); err != nil {
			s.diag.Error("NotifyALL", err)
		}
	}
}

func (s *WSClientManager) CloseAll() {

	cons := s.ws.GetConnections()
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()
	for _, v := range cons {
		if err := v.DisconnectAll(ctx); err != nil {
			s.diag.Error(fmt.Sprintf("Close Websocket Client:%s Failed", v.String()), err)
		}
	}
}
