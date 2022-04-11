package wsnotify

import (
	"github.com/masami10/rush/services/httpd"
	"github.com/masami10/rush/utils"
)

type HTTPService interface {
	AddNewWebsocketHandler(r httpd.WebsocketRoute) error
}

type Dispatcher interface {
	Create(name string, len int) error
	Start(name string) error
	Dispatch(name string, data interface{}) error
	Register(name string, handler *utils.DispatchHandlerStruct)
	Release(name string, handler string) error
}
