package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

type Bolt struct {
	Name      string
	Precision float64
	Recall    float64
	Number    int
}

var Bolts = map[string]Bolt{
	"OP1000-0-0": {
		Name:      "OP1000-0-0",
		Precision: 100.00,
		Recall:    99.68,
		Number:    403,
	},

	"H170-0-0": {
		Name:      "H170-0-0",
		Precision: 92.31,
		Recall:    99.46,
		Number:    746,
	},
}

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
	par := "H170-0-0"
	if len(s) >= 3 {
		par = s[2]
	}

	expiration := time.Now()
	expiration = expiration.AddDate(1, 0, 0)
	cookie := http.Cookie{
		Name:    "username",
		Value:   "hello cookie",
		Expires: expiration,
	}
	http.SetCookie(w, &cookie)

	var bolt Bolt
	bolt = Bolts["H170-0-0"]
	if _, ok := Bolts[par]; ok {
		bolt = Bolts[par]
	}

	err = t1.Execute(w, bolt)
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
