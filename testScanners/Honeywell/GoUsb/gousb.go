package GoUsb

import (
	"errors"
	"fmt"
	"github.com/google/gousb"
	"strconv"
	"strings"
	"time"
)

type GoUsb struct {
	Name string
}

//sudo chmod 777 /dev/bus/usb/001/008
//修改权限为可读可写可执行，但是这种设置电脑重启后，又会出现这种问题，还要重新设置．因此查询资料，可以用下面这条指令：
//
//　　sudo usermod -aG　dialout kevin


func (s1 *GoUsb) search() error {
	ctx := gousb.NewContext()
	defer ctx.Close()
	var vid, pid int64
	ls := strings.Split(s1.Name, ":")
	if len(ls) != 2 {
		return errors.New(fmt.Sprintf("配置文件的pidvid错误:%s", s1.Name))
	}
	vid, _ = strconv.ParseInt(ls[0], 10, 16)
	pid, _ = strconv.ParseInt(ls[1], 10, 16)

	d, err := ctx.OpenDeviceWithVIDPID(gousb.ID(vid), gousb.ID(pid))
	if err == nil && d != nil {
		fmt.Println(d)
	} else if err != nil {
		fmt.Printf("gousb错误：%v\n", err)
		return err
	} else {
		err = errors.New(fmt.Sprintf("Open Fail VID:%d, PID:%d", vid, pid))
		fmt.Printf("自定义错误：%v\n", err)
		return err
	}
	return nil
}

func (s1 *GoUsb) Connect() {
	//打开串口
	for {
		err := s1.search()
		if err != nil {
			time.Sleep(1 * time.Second)
			continue
		}
		break
	}
	fmt.Printf("%s 连接成功\n", s1.Name)
}

func (s1 *GoUsb) Read() (string, error) {
	//打开串口
	time.Sleep(5 * time.Second)
	return "hello world", nil
}
