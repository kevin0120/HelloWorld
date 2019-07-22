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

	app.Run(iris.Addr(":8080"))
}
