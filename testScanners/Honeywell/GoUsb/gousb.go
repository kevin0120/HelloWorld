package GoUsb

import (
	"errors"
	"fmt"
	"github.com/google/gousb"
	"log"
	"strconv"
	"strings"
	"time"
)

type GoUsb struct {
	Name       string
	InEndpoint *gousb.InEndpoint
	Finish     chan string
	Scanner
}

type Scanner struct {
	Debounced       func(f func())
	debounceTrigger bool
	init            bool
}

//sudo chmod 777 /dev/bus/usb/001/008
//修改权限为可读可写可执行，但是这种设置电脑重启后，又会出现这种问题，还要重新设置．因此查询资料，可以用下面这条指令：
// udev device permissions
//　　sudo usermod -aG　input kevin

//linux work before
//sudo nano /etc/udev/rules.d/51-blink1.rules
//SUBSYSTEM=="input",GROUP="input",MODE="0666"
//SUBSYSTEM=="usb", ATTRS{idVendor}=="0c2e",ATTRS{idProduct}=="0901",MODE="0666",GROUP="plugdev"
//KERNEL=="hidraw*", ATTRS{idVendor}=="0c2e",ATTRS{idProduct}=="0901",MODE="0666",GROUP="plugdev"

//sudo udevadm control --reload-rules

func (s1 *GoUsb) search() error {
	ctx := gousb.NewContext()
	//defer ctx.Close()

	//ctx.Debug(3)
	var vid, pid int64
	ls := strings.Split(s1.Name, ":")
	if len(ls) != 2 {
		return errors.New(fmt.Sprintf("配置文件的pidvid错误:%s", s1.Name))
	}
	vid, _ = strconv.ParseInt(ls[0], 10, 16)
	pid, _ = strconv.ParseInt(ls[1], 10, 16)

	dev, err := ctx.OpenDeviceWithVIDPID(gousb.ID(vid), gousb.ID(pid))

	//devs, err := ctx.OpenDevices(func(desc *gousb.DeviceDesc) bool {
	//	switch {
	//	case gousb.ID(vid) == desc.Vendor && gousb.ID(pid) == desc.Product:
	//		return true
	//	}
	//	return false
	//})
	//
	//switch {
	//case len(devs) == 0:
	//	log.Fatal("No matching devices found.")
	//case len(devs) > 1:
	//	log.Printf("Warning: multiple devices found. Using bus %d, addr %d.", devs[0].Desc.Bus, devs[0].Desc.Address)
	//	for _, d := range devs[1:] {
	//		d.Close()
	//	}
	//	devs = devs[:1]
	//}
	//dev := devs[0]

	if err == nil && dev != nil {

		fmt.Println(dev)
	} else if err != nil {
		fmt.Printf("gousb错误：%v\n", err)
		return err
	} else {
		err = errors.New(fmt.Sprintf("Open Fail VID:%d, PID:%d", vid, pid))
		fmt.Printf("自定义错误：%v\n", err)
		return err
	}

	log.Print("Enabling autodetach")
	dev.SetAutoDetach(true)

	log.Printf("Setting configuration %d...", 1)
	cfg, err := dev.Config(1)
	if err != nil {
		log.Fatalf("dev.Config(%d): %v", 1, err)
	}
	log.Printf("Claiming interface %d (alt setting %d)...", 0, 0)
	intf, err := cfg.Interface(0, 0)
	if err != nil {
		log.Fatalf("cfg.Interface(%d, %d): %v", 0, 0, err)
	}

	log.Printf("Using endpoint %d...", 4)
	ep, err := intf.InEndpoint(4)
	if err != nil {
		log.Fatalf("dev.InEndpoint(): %s", err)
	}
	log.Printf("Found endpoint: %s", ep)

	s1.InEndpoint = ep

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

func (s1 *GoUsb) read() {
	buf := make([]byte, 1024)
	//log.Print("Reading...")
	strRecv := ""
	for {
		num, err := s1.InEndpoint.Read(buf)
		if err != nil {
			log.Fatalf("Reading from device failed: %v", err)
		}
		if num > 0 {
			s1.TriggerDebounce()
			s, e := CommonParse(buf[0:num])
			if e == nil {
				strRecv += s

			}
			s1.Debounced(func() {
				if strRecv != "" {
					//fmt.Printf("receive:%s", strRecv)
					s1.ResetDebounce()
					s1.Finish <- strRecv
					strRecv = ""
				}
			})
		}
		//os.Stdout.Write(buf[:num])
	}
}

func (s1 *GoUsb) Read() (string, error) {

	if s1.Finish == nil {
		s1.Finish = make(chan string, 1)
		go s1.read()
	}
	s := <-s1.Finish

	return s, nil
	//打开串口

}
