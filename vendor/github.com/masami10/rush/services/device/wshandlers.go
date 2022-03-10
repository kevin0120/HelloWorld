package device

import (
	"encoding/json"
	"github.com/kataras/iris/v12/websocket"
	"github.com/masami10/rush/services/wsnotify"
)

func (s *Service) OnWSDeviceStatus(c *websocket.NSConn, msg *wsnotify.WSMsg) {
	devices := s.fetchAllDevices()
	body, _ := json.Marshal(wsnotify.GenerateWSMsg(msg.SeqNumber, msg.Type, devices))

	_ = wsnotify.WSClientSend(c, wsnotify.WS_EVENT_DEVICE, string(body), s.diag)
}
