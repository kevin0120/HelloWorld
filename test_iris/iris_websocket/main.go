package main

import (
	"fmt"
	"github.com/kataras/iris"
	"github.com/kataras/iris/websocket"
)

func main() {
	ws := websocket.New(websocket.Config{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	})

	ws.OnConnection(handleConnection)

	app := iris.New()
	// 在端点上注册一个服务.
	// 请参阅 websockets.html 中的内联 JavaScrip t代码，此端点用于连接到服务器。
	app.Get("/echo", ws.Handler())

	// 提供内置客户端库的 javascript，
	//请参阅 weboskcets.html 脚本标记，使用此路径。
	app.Any("/iris-ws.js", func(ctx iris.Context) {
		ctx.Write(websocket.ClientSource)
	})

	app.Run(iris.Addr(":8080"))
}

func handleConnection(c websocket.Connection) {
	// 从浏览器中读取事件
	c.On("chat", func(msg string) {
		// 将消息打印到控制台，c .Context（）是iris的http上下文。
		fmt.Printf("%s sent: %s\n", c.Context().RemoteAddr(), msg)
		// 将消息写回客户端消息所有者：
		// c.Emit("chat", msg)
		c.To(websocket.Broadcast).Emit("chat", msg)
	})
}
