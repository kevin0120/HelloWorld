package wsnotify

import (
	"github.com/kataras/iris/v12/websocket"
)

type DispatcherNotifyPackage struct {
	C    *websocket.NSConn
	Data []byte
}
