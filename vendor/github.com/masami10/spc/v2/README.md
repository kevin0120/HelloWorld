# spc
SPC库

1. 获取子模块
```bash
1. git submodule update --init --recursive
```

2. 安装boost
```bash
1. pacman -S mingw-w64-x86_64-boost // mingw-w64-x86_64 环境
2. apt-get  install -y libboost-all-dev  // 安装boost
```

3. install go wrapper
```bash
go get github.com/masami10/spc/v2
```

3.  go demo
```bash
package main

import (
	"fmt"
	_ "github.com/masami10/spc/v2"
	"github.com/masami10/spc/v2/wrapper/golang/spc"
)

func main() {
	var a = []float64{5., 5., 10., 12., 5., 5., 10., 12., 5., 5.}

	c, _ := spc.Cpk(a, 10, 14, 4)
	fmt.Printf("cpk: %f", c)
}
```