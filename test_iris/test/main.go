package main

import (
	"github.com/kataras/iris"
)

func main() {
	app := iris.Default()

	// Query string parameters are parsed using the existing underlying request object.
	// The request responds to a url matching:  /welcome?firstname=Jane&lastname=Doe.
	app.Get("/welcome", func(ctx iris.Context) {
		firstname := ctx.URLParamDefault("firstname", "Guest")
		// shortcut for ctx.Request().URL.Query().Get("lastname").
		lastname := ctx.URLParam("lastname")

		ctx.Writef("Hello %s %s", firstname, lastname)
	})

	//只能用127.0.0.1访问
	//app.Run(iris.Addr("127.0.0.1:8080"))

	//只能用192.168.4.188访问
	//app.Run(iris.Addr("192.168.4.188:8080"))

	//即能用192.168.4.188访问,也能用127.0.0.1访问,等价与0.0.0.0
	//app.Run(iris.Addr(":8080"))

	app.Run(iris.Addr(":8080"))
}
