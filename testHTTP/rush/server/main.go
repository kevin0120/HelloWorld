package main

import (
	"fmt"
	"github.com/iris-contrib/middleware/cors"
	"github.com/kataras/iris"
	"github.com/kataras/iris/context"
)

const (
	// Root path for the API
	BasePath        = "/rush/v1"
	ROUTE_TYPE_HTTP = "http"
	ROUTE_TYPE_WS   = "websocket"
	IP              = "localhost:8082"
)

type Route struct {
	RouteType   string
	Method      string
	Pattern     string
	HandlerFunc context.Handler
}

type Handler struct {
	party   *iris.Party
	service *iris.Application
}

func (h *Handler) AddRoute(r Route) error {
	if len(r.Pattern) > 0 && r.Pattern[0] != '/' {
		return fmt.Errorf("route patterns must begin with a '/' %s", r.Pattern)
	}
	if r.RouteType == ROUTE_TYPE_HTTP {
		(*h.party).Handle(r.Method, r.Pattern, r.HandlerFunc)
	} else {
		h.service.Get(r.Pattern, r.HandlerFunc)
	}

	return nil
}

func main() {

	server := iris.New()

	crs := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
		AllowedMethods:   []string{"GET", "HEAD", "POST", "PUT", "PATCH", "OPTIONS"},
	})

	p := server.Party(BasePath, crs).AllowMethods(iris.MethodOptions)

	party := &p

	h := &Handler{
		party:   party,
		service: server,
	}

	var r Route

	///http://127.0.0.1:8082/rush/v1/hello/3?q=1
	r = Route{
		RouteType: ROUTE_TYPE_HTTP,
		Method:    "GET",
		Pattern:   "/hello/{id:int}",
		HandlerFunc: func(i context.Context) {
			fmt.Println(i.String())
			i.Write([]byte("hello world!"))
			fmt.Println(i.URLParam("q"))
			fmt.Println(i.Params().Get("id"))
		},
	}
	h.AddRoute(r)

	err := server.Run(iris.Addr(IP), iris.WithoutInterruptHandler)

	if err != nil {
		fmt.Println(err)

	}

}
