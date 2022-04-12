package main

import (
	"github.com/kataras/iris"
	"github.com/kataras/iris/context"
)

func main() {
	//1.创建app结构体对象
	app := iris.New()

	//url: http://localhost:8000/getRequest
	// type：GET请求
	app.Get("/getRequest", func(context context.Context) {
		path := context.Path()
		app.Logger().Info(path)
		context.WriteString("HelloWorld")
		context.Params().GetInt("id")
	})
	//url: http://localhost:/user/info //type：POST请求
	app.Handle("POST", "/user/info", func(context context.Context) {
		context.WriteString(" User Info is Post Request , Deal is in handle func ")
	})

	//2.端口监听
	app.Run(iris.Addr("127.0.0.1:8080"), iris.WithoutServerError(iris.ErrServerClosed))
	////application.Run(iris.Addr(":8080"))//第一种
	// application.Run(iris.Addr(":8080"), iris.WithoutServerError(iris.ErrServerClosed)) //第二种
}
