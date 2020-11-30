package main

import (
	"fmt"
	"github.com/kataras/iris"
)
//https://learnku.com/docs/iris-go/10/mvc_2/3774
func main() {
	app := iris.Default()

	party1 := app.Party("/rush/v1")

	party2 := app.Party("/rush/v2")

	party1.Handle("GET", "/hello", func(ctx iris.Context) {
		ctx.Writef("Hello /rush/v1")
	})

	party2.Get("/world", func(ctx iris.Context) {
		ctx.Writef("Hello /rush/v2")
	})
	// Query string parameters are parsed using the existing underlying request object.
	// The request responds to a url matching:  /welcome?firstname=Jane&lastname=Doe.
	app.Get("/welcome/{key}", func(ctx iris.Context) {
		firstname := ctx.URLParamDefault("firstname", "Guest")
		// shortcut for ctx.Request().URL.Query().Get("lastname").
		lastname := ctx.URLParam("lastname")

		a := ctx.Params()
		fmt.Println(a.Get("key"))
		ctx.Writef("Hello %s %s", firstname, lastname)
	})

	//只能用127.0.0.1访问
	//app.Run(iris.Addr("127.0.0.1:8080"))

	//只能用192.168.4.188访问
	//app.Run(iris.Addr("192.168.4.188:8080"))

	//即能用192.168.4.188访问,也能用127.0.0.1访问,等价与0.0.0.0
	//app.Run(iris.Addr(":8080"))

	app.Get("/luyou", func(ctx iris.Context) { ctx.Writef("Hello"); ctx.Next() },
	func(ctx iris.Context) { ctx.Writef("world!"); ctx.Next() },
		func(ctx iris.Context) { ctx.Writef("123!"); ctx.Next() })

	app.Run(iris.Addr(":8080"))
}
