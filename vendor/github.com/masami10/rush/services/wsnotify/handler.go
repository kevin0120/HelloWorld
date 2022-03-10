package wsnotify

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/kataras/iris/v12/websocket"
)

type WSRequestHandlerDiag interface {
	Error(msg string, err error)
	Debug(msg string)
}

type WSRequestHandlerMap map[string]WSRequestHandler
type WSRequestHandler func(c *websocket.NSConn, wsMsg *WSMsg)

type WSRequestHandlers struct {
	handlers WSRequestHandlerMap
	Diag     WSRequestHandlerDiag
}

func (h *WSRequestHandlers) SetupHandlers(handlers WSRequestHandlerMap) {
	if h.handlers == nil {
		h.handlers = WSRequestHandlerMap{}
	}

	h.handlers = handlers
}

func (h *WSRequestHandlers) getHandler(name string) (WSRequestHandler, error) {
	handler, exist := h.handlers[name]
	if !exist {
		return nil, errors.New("WSRequest Handler Not Found")
	}

	return handler, nil
}

func (h *WSRequestHandlers) HandleWSRequest(data interface{}) {

	if data == nil {
		h.Diag.Error("HandleWSRequest Error: Data is Nil", errors.New("HandleWSRequest Error: Data is Nil"))
		return
	}

	req := data.(*DispatcherNotifyPackage)
	var msg WSMsg
	err := json.Unmarshal(req.Data, &msg)
	if err != nil {
		h.Diag.Error(fmt.Sprintf("HandleWSRequest Error: Fail With Message %# 20X", req.Data), err)
		return
	}

	handler, err := h.getHandler(msg.Type)
	if err != nil {
		return
	}

	handler(req.C, &msg)
}
