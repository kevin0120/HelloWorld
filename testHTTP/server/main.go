package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	_ "net/http/pprof"
)

func helloHandler(w http.ResponseWriter, r *http.Request) {

	fmt.Println("method=", r.Method)
	_, err := w.Write([]byte(homeTemplate))
	if err != nil {
		fmt.Println(err)
	}
}

func echoHandler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Hello, world!\n")
	io.WriteString(w, r.URL.Path)
}
func Hello(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Hello World.")

	fmt.Fprintf(w, "Hello World!!!!.\n")
	time.Sleep(1 * time.Second)
	io.WriteString(w, "Hello World.\n")

}

//w, 给客户端回复数据
//r, 读取客户端发送的数据
func HandConn(w http.ResponseWriter, r *http.Request) {
	fmt.Println("r.Method = ", r.Method)
	fmt.Println("r.URL = ", r.URL)
	fmt.Println("r.Header = ", r.Header)
	fmt.Println("r.Body = ", r.Body)
	a, _ := json.Marshal(map[string]string{"info": "a", "master": "b", "controllers": "c", "gun": "d"})
	c := `{"controllers":"c","gun":"d","info":"a","master":"b"}`
	//c这种写法是不可以的
	b, _ := json.Marshal(c)
	fmt.Println(string(b))
	fmt.Println(string(a))
	w.Write(a) //给客户端回复数据
}

func test(w http.ResponseWriter, r *http.Request) {

	w.Write([]byte("该处理函数延时30妙注册")) //给客户端回复数据
	fmt.Println(r.RemoteAddr)
}

func mytestHandle() {
	time.Sleep(30 * time.Second)
	http.HandleFunc("/api/v1/test", test)
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/hello/{workcenter:string}", helloHandler)
	mux.HandleFunc("/home", echoHandler)
	http.HandleFunc("/api/v1/healthz", Hello)
	//注册处理函数，用户连接，自动调用指定的处理函数
	http.HandleFunc("/api/v1/workcenter", HandConn)
	go mytestHandle()
	go func() {
		log.Println(http.ListenAndServe("0.0.0.0:8068", nil))
	}()
	fmt.Println("ldfdl")
	//
	//http://127.0.0.1:8070/debug/pprof/   可以查看攜程信息...
	go func() {
		log.Println(http.ListenAndServe("0.0.0.0:8070", nil))
	}()
	http.ListenAndServe(":12345", mux)

}

var homeTemplate = `
<!DOCTYPE html>
<html>
<head>
<meta charset="utf-8">
<script>  
window.addEventListener("load", function(evt) {
    var output = document.getElementById("output");
    var input = document.getElementById("input");
    var ws;
    var print = function(message) {
        var d = document.createElement("div");
        d.innerHTML = message;
        output.appendChild(d);
    };
    document.getElementById("open").onclick = function(evt) {
        if (ws) {
            return false;
        }
        ws = new WebSocket("{{.}}");
        ws.onopen = function(evt) {
            print("OPEN");
        }
        ws.onclose = function(evt) {
            print("CLOSE");
            ws = null;
        }
        ws.onmessage = function(evt) {
            print("RESPONSE: " + evt.data);
        }
        ws.onerror = function(evt) {
            print("ERROR: " + evt.data);
        }
        return false;
    };
    document.getElementById("send").onclick = function(evt) {
        if (!ws) {
            return false;
        }
        print("SEND: " + input.value);
        ws.send(input.value);
        return false;
    };
    document.getElementById("close").onclick = function(evt) {
        if (!ws) {
            return false;
        }
        ws.close();
        return false;
    };
});
</script>
</head>
<body>
<table>
<tr><td valign="top" width="50%">
<p>Click "Open" to create a connection to the server, 
"Send" to send a message to the server and "Close" to close the connection. 
You can change the message and send multiple times.
<p>
<form>
<button id="open">Open</button>
<button id="close">Close</button>
<p><input id="input" type="text" value="Hello world!">
<button id="send">Send</button>
</form>
</td><td valign="top" width="50%">
<div id="output"></div>
</td></tr></table>
</body>
</html>
`
