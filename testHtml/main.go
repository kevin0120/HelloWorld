package main

import (
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
)

func Hello(w http.ResponseWriter, r *http.Request) {
	hp := `<html>
    <head>
    <title>okkkkkk</title>
    <link rel="stylesheet" href="/template/css/main.css" type="text/css" /> 
    </head>
    <body>
        <h2>this is a test for golang.</h2>
    </body>
    </html>`
	io.WriteString(w, hp)
}

//func StaticServer(w http.ResponseWriter, r *http.Request) {
//	w.Header().Set("content-type", "text/html")
//	staticHandler := http.FileServer(http.Dir("./template/"))
//	staticHandler.ServeHTTP(w, r)
//	return
//}

func tmpl(w http.ResponseWriter, r *http.Request) {
	t1, err := template.ParseFiles("./template/hello.html")
	if err != nil {
		panic(err)
	}
	err = t1.Execute(w, "hello world")
	if err != nil {
		return
	}
}

func main() {
	fmt.Println(os.Getwd())
	//http.Handle("/template/", http.StripPrefix("/template/", http.FileServer(http.Dir("./template"))))
	////http.HandleFunc("/", Hello)

	http.HandleFunc("/tmpl", tmpl)
	http.Handle("/template/", http.FileServer(http.Dir("./")))

	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err.Error())
	}
}
