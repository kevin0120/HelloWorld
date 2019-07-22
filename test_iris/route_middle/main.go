/*
当我们在 iris 中讨论中间件的时候，我们就是在讨论运行一个 HTTP 请求处理器前后的代码。
例如，日志中间件会把即将到来的请求明细写到日志中，然后在写响应日志详细之前，调用处理期代码。
比较酷的是这些中间件是松耦合，而且可重用。
中间件仅仅是 func(ctx iris.Context) 的处理器形式，中间件在前一个调用 ctx.Next() 时执行，
这个可以用于去认证，例如，如果登录了，调用  ctx.Next() 否则将触发一个错误响应。
*/
// 注册 "before"  处理器作为当前域名所有路由中第一个处理函数
// 或者使用  `UseGlobal`  去注册一个中间件，用于在所有子域名中使用
//app.Use(before)
// 注册  "after" ，在所有路由的处理程序之后调用
//app.Done(after)

// 注册路由
//app.Get("/", indexHandler)
///app.Get("/contact", contactHandler)

package main

import "github.com/kataras/iris"

func main() {
	app := iris.New()
	app.Get("/", before, mainHandler, after)
	app.Get("/q", before) //松耦合.可重用
	app.Run(iris.Addr(":8080"))
}

func before(ctx iris.Context) {
	shareInformation := "this is a sharable information between handlers"

	requestPath := ctx.Path()
	println("Before the mainHandler: " + requestPath)

	ctx.Values().Set("info", shareInformation)
	ctx.Next() // 执行下一个处理器。
}

func after(ctx iris.Context) {
	println("After the mainHandler")
}

func mainHandler(ctx iris.Context) {
	println("Inside mainHandler")

	// 获取 "before" 处理器中的设置的 "info" 值。
	info := ctx.Values().GetString("info")

	// 响应客户端
	ctx.HTML("<h1>Response</h1>")
	ctx.HTML("<br/> Info: " + info)

	ctx.Next() // execute the "after".
}
