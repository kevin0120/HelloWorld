package io

import (
	"encoding/json"
	"github.com/kataras/iris/v12/websocket"
	"github.com/masami10/rush/services/device"
	"github.com/masami10/rush/services/wsnotify"
)

// 获取连接状态
func (s *Service) OnWSIOStatus(c *websocket.NSConn, msg *wsnotify.WSMsg) {
	byteData, _ := json.Marshal(msg.Data)

	ioStatus := device.Status{}
	err := json.Unmarshal(byteData, &ioStatus)
	if err != nil {
		_ = wsnotify.WSClientSend(c, wsnotify.WS_EVENT_REPLY, wsnotify.GenerateReply(msg.SeqNumber, msg.Type, -1, err.Error()), s.diag)
		return
	}

	m, err := s.getIO(ioStatus.SN)
	if err != nil {
		_ = wsnotify.WSClientSend(c, wsnotify.WS_EVENT_REPLY, wsnotify.GenerateReply(msg.SeqNumber, msg.Type, -2, err.Error()), s.diag)
		return
	}

	wsMsg := wsnotify.GenerateWSMsg(msg.SeqNumber, wsnotify.WS_IO_STATUS, []device.Status{
		{
			SN:     ioStatus.SN,
			Type:   device.BaseDeviceTypeIO,
			Status: m.Status(),
		},
	})

	_ = wsnotify.WSClientSend(c, wsnotify.WS_EVENT_IO, wsMsg, s.diag)
}

// 获取io状态
func (s *Service) OnWSIOContact(c *websocket.NSConn, msg *wsnotify.WSMsg) {
	byteData, _ := json.Marshal(msg.Data)

	ioContact := IoContact{}
	err := json.Unmarshal(byteData, &ioContact)
	if err != nil {
		_ = wsnotify.WSClientSend(c, wsnotify.WS_EVENT_REPLY, wsnotify.GenerateReply(msg.SeqNumber, msg.Type, -1, err.Error()), s.diag)
		return
	}

	inputs, outputs, err := s.Read(ioContact.SN)
	if err != nil {
		_ = wsnotify.WSClientSend(c, wsnotify.WS_EVENT_REPLY, wsnotify.GenerateReply(msg.SeqNumber, msg.Type, -2, err.Error()), s.diag)
		return
	}

	wsMsg := wsnotify.GenerateWSMsg(msg.SeqNumber, wsnotify.WS_IO_CONTACT, IoContact{
		Src:     device.BaseDeviceTypeIO,
		SN:      ioContact.SN,
		Inputs:  inputs,
		Outputs: outputs,
	})

	_ = wsnotify.WSClientSend(c, wsnotify.WS_EVENT_REPLY, wsMsg, s.diag)
	_ = wsnotify.WSClientSend(c, wsnotify.WS_EVENT_IO, wsMsg, s.diag)
}

// 控制输出
func (s *Service) OnWSIOSet(c *websocket.NSConn, msg *wsnotify.WSMsg) {
	byteData, _ := json.Marshal(msg.Data)

	ioSet := IoSet{}
	err := json.Unmarshal(byteData, &ioSet)
	if err != nil {
		_ = wsnotify.WSClientSend(c, wsnotify.WS_EVENT_REPLY, wsnotify.GenerateReply(msg.SeqNumber, msg.Type, -1, err.Error()), s.diag)
		return
	}

	err = s.Write(ioSet.SN, ioSet.Index, ioSet.Status)
	if err != nil {
		_ = wsnotify.WSClientSend(c, wsnotify.WS_EVENT_REPLY, wsnotify.GenerateReply(msg.SeqNumber, msg.Type, -2, err.Error()), s.diag)
		return
	}

	_ = wsnotify.WSClientSend(c, wsnotify.WS_EVENT_REPLY, wsnotify.GenerateReply(msg.SeqNumber, msg.Type, 0, ""), s.diag)
}
