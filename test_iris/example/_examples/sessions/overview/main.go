package main

import (
	"github.com/kataras/iris/v12"

	"github.com/kataras/iris/v12/sessions"
)

var (
	cookieNameForSessionID = "mycookiesessionnameid"
	sess                   = sessions.New(sessions.Config{Cookie: cookieNameForSessionID, AllowReclaim: true})
)

func secret(ctx iris.Context) {

	//ctx.SetCookieKV("mycookie1","hello world12")
	// Check if user is authenticated
	session := sess.Start(ctx)
	if auth, _ := session.GetBoolean("authenticated"); !auth {
		ctx.StatusCode(iris.StatusForbidden)
		return
	}

	// Print secret message
	ctx.WriteString("The cake is a lie!")
}

func login(ctx iris.Context) {
	//ctx.RemoveCookie("session_id")
	session := sess.Start(ctx)

	// Authentication goes here
	// ...

	// Set user as authenticated
	session.Set("authenticated", true)
}

func logout(ctx iris.Context) {
	session := sess.Start(ctx)

	// Revoke users authentication
	session.Set("authenticated", false)
}

func main() {
	app := iris.New()

	app.Get("/secret", secret)
	app.Get("/login", login)
	app.Get("/logout", logout)

	app.Run(iris.Addr(":8080"))
}
