package utils

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"
)

type Pkg struct {
	Key   string `json:"key"`
	Value string `json:"value"`
	K2    string `json:"k2"`
	K3    string `json:"k3"`
}

func Test_dispatcher(t *testing.T) {
	d := CreateDispatcher(DefaultDispatcherBufLen)
	d.Register("d", onData)
	d.Start()
	d.Dispatch(&Pkg{
		Key:   "13",
		Value: "12313",
	})

	time.Sleep(1 * time.Second)
}

func onData(data interface{}) {
	fmt.Println(data)
}

type PkgM struct {
	Key   string                 `json:"key"`
	Value string                 `json:"value"`
	Rest  map[string]interface{} `json:"-"`
}

func Test_common(t *testing.T) {
	pkg := Pkg{
		Key:   "1",
		Value: "2",
		K2:    "3",
		K3:    "4",
	}

	strPkg, _ := json.Marshal(pkg)

	pkgm := PkgM{}
	json.Unmarshal(strPkg, &pkgm)

	fmt.Println(pkgm)
}
