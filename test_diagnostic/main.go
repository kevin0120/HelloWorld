package main

import (
	"HelloWorld/test_diagnostic/diagnostic"
	"fmt"
	"os"
	"time"

	"github.com/masami10/rush/keyvalue"
)

var (
	version string
	commit  string
	branch  string
)

type Service struct {
	Tag string
	//Deprecated
	Abandon string
}

type Service1 struct {
	Tag string
	//Deprecated
	Abandon string
}

//Deprecated
func (s Service) say() {
	fmt.Println(s)
}
func main() {

	c := diagnostic.Config{
		File:   "STDOUT",
		//File:   "/home/kevin/Downloads/gopath/src/HelloWorld/test_diagnostic/log/%Y%m%d.log",
		Level:  "DEBUG",
		MaxAge: time.Duration(3000 * time.Hour),
		Rotate: time.Duration(24 * time.Hour),
	}


	s := Service{
		Tag:     "hello service",
		Abandon: "ssss",
	}

	s.say()


	fmt.Println("############################", "diaService")
	diagService := diagnostic.NewService(c, os.Stdout, os.Stderr)

	if err := diagService.Open(); err != nil {
		fmt.Println(fmt.Errorf("failed to open diagnostic service: %v", err))
	}

	fmt.Println("##############################3", "server")
	sever_Diag := diagService.NewServerHandler()
	sever_Diag.Debug("opening service", keyvalue.KV("service", fmt.Sprintf("%T", s)))
	err := fmt.Errorf("%s", "hello err!")
	sever_Diag.Error("error closing service", err, keyvalue.KV("service", fmt.Sprintf("%T", s)))
	sever_Diag.Info("opening service", keyvalue.KV("service", fmt.Sprintf("%T", s)))

	fmt.Println("##############################3", "diatightening")
	tightening_Diag := diagService.NewTighteningDeviceHandler()
	tightening_Diag.Error("Create Controller Failed", err)
	tightening_Diag.Debug("hello debug")

	fmt.Println("hello world!")

}
