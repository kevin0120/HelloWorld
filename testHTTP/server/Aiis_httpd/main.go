package main

import (
	"fmt"
	"github.com/kataras/iris"
	"os"
	"time"

	//	"github.com/kataras/iris/hero"
	"strconv"
)

func posting(ctx iris.Context) {
	fmt.Println(ctx.Request())
}
func getting(ctx iris.Context) {
	fmt.Println(ctx.Write([]byte("hello world!!!!")))

}

func main2() {
	app := iris.Default()

	// Simple group: v1.
	v1 := app.Party("/v1")
	{
		v1.Post("/login", posting)
		v1.Post("/submit", posting)
		v1.Post("/read", posting)
	}

	// Simple group: v2.
	v2 := app.Party("/v2")
	{
		v2.Post("/login", posting)
		v2.Post("/submit", posting)
		v2.Post("/read", posting)
	}

	app.Run(iris.Addr(":8081"))
}
func main() {
	// Creates an application with default middleware:
	// logger and recovery (crash-free) middleware.
	app := iris.Default()

	app.Get("/someGet", getting)
	app.Post("/somePost", posting)

	// This handler will match /users/42
	// but will not match /users/-1 because uint should be bigger than zero
	// neither /users or /users/.
	app.Get("/users/{id:uint64}", func(ctx iris.Context) {
		//若浏览器断输入id则id==id否则id=0
		id := ctx.Params().GetUint64Default("id", 1)

		fmt.Println(ctx.Write([]byte(strconv.Itoa(int(id)))))
		// [...]
	})
	// This handler will match /user/john but will not match neither /user/ or /user.
	app.Get("/user/{name}", func(ctx iris.Context) {
		name := ctx.Params().Get("name")
		ctx.Writef("Hello %s", name)
	})

	// However, this one will match /user/john/send and also /user/john/everything/else/here
	// but will not match /user/john neither /user/john/.
	app.Post("/user/{name:string}/{action:path}", func(ctx iris.Context) {
		name := ctx.Params().Get("name")
		action := ctx.Params().Get("action")
		message := name + " is " + action
		ctx.WriteString(message)
	})

	//hero.Handler(hello)

	go main2()

	//可以将信息打印到文件中去
	/*
		f := newLogFile()
		defer f.Close()
		app.Logger().SetOutput(f)
	*/
	//app := iris.New()
	// Attach the file as logger, remember, iris' app logger is just an io.Writer.
	// Use the following code if you need to write the logs to file and console at the same time.
	// app.Logger().SetOutput(io.MultiWriter(f, os.Stdout))

	//app.Logger().SetOutput(f)

	app.Get("/ping", func(ctx iris.Context) {
		// for the sake of simplicity, in order see the logs at the ./_today_.txt
		ctx.Application().Logger().Infof("Request path: %s", ctx.Path())
		ctx.WriteString("pong")
	})

	// Navigate to http://localhost:8080/ping
	// and open the ./logs{TODAY}.txt file.
	app.Run(iris.Addr(":8080"))

}

// Get a filename based on the date, just for the sugar.
func todayFilename() string {
	today := time.Now().Format("Jan 02 2006")
	return today + ".txt"
}

func newLogFile() *os.File {
	filename := todayFilename()
	// Open the file, this will append to the today's file if server restarted.
	f, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}

	return f
}
