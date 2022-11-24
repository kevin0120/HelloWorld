package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"
)

func Hello(w http.ResponseWriter, r *http.Request) {
	t1, err := template.ParseFiles("./testHtml/template/hello.html")
	if err != nil {
		panic(err)
	}
	err = t1.Execute(w, "hello world")
	if err != nil {
		return
	}
}

func report(w http.ResponseWriter, r *http.Request) {
	t1, err := template.ParseFiles("./testHtml/template/report.html")
	if err != nil {
		panic(err)
	}

	s := strings.Split(r.RequestURI, "/")
	par := "hello world"
	if len(s) >= 3 {
		par = s[2]
	}

	err = t1.Execute(w, par)
	if err != nil {
		return
	}
}

func main() {
	fmt.Println(os.Getwd())
	//http.Handle("/template/", http.StripPrefix("/template/", http.FileServer(http.Dir("./template"))))
	http.HandleFunc("/hello", Hello)

	http.HandleFunc("/report/", report)
	http.Handle("/", http.FileServer(http.Dir("./testHtml/template")))

	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err.Error())
	}
}
