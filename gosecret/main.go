package main

import (
	"fmt"
	"github.com/StackExchange/wmi"
	"net"
	"strings"
	"syscall"
	"time"
	"unsafe"
)

type cpuInfo struct {
	Name          string
	NumberOfCores uint32
	ThreadCount   uint32
}

func getCPUInfo() {

	var cpuinfo []cpuInfo

	err := wmi.Query("Select * from Win32_Processor", &cpuinfo)
	if err != nil {
		return
	}
	fmt.Printf("Cpu info =%+v \n", cpuinfo)
}

type operatingSystem struct {
	Name    string
	Version string
}

func getOSInfo() {
	var operatingsystem []operatingSystem
	err := wmi.Query("Select * from Win32_OperatingSystem", &operatingsystem)
	if err != nil {
		return
	}
	fmt.Printf("OS info =%+v \n", operatingsystem)
}

var kernel = syscall.NewLazyDLL("Kernel32.dll")

type memoryStatusEx struct {
	cbSize                  uint32
	dwMemoryLoad            uint32
	ullTotalPhys            uint64 // in bytes
	ullAvailPhys            uint64
	ullTotalPageFile        uint64
	ullAvailPageFile        uint64
	ullTotalVirtual         uint64
	ullAvailVirtual         uint64
	ullAvailExtendedVirtual uint64
}

func getMemoryInfo() {

	GlobalMemoryStatusEx := kernel.NewProc("GlobalMemoryStatusEx")
	var memInfo memoryStatusEx
	memInfo.cbSize = uint32(unsafe.Sizeof(memInfo))
	mem, _, _ := GlobalMemoryStatusEx.Call(uintptr(unsafe.Pointer(&memInfo)))
	if mem == 0 {
		return
	}
	fmt.Printf("total=:%+v ", memInfo.ullTotalPhys)
	fmt.Printf("free=:%+v \n", memInfo.ullAvailPhys)
}

type Network struct {
	Name       string
	IP         string
	MACAddress string
}

type intfInfo struct {
	Name       string
	MacAddress string
	Ipv4       []string
}

func getNetworkInfo() error {
	intf, err := net.Interfaces()
	if err != nil {
		fmt.Errorf("get network info failed: %v", err)
		return err
	}
	var is = make([]intfInfo, len(intf))
	for i, v := range intf {
		ips, err := v.Addrs()
		if err != nil {
			fmt.Errorf("get network addr failed: %v", err)
			return err
		}
		//此处过滤loopback（本地回环）和isatap（isatap隧道）
		if !strings.Contains(v.Name, "Loopback") && !strings.Contains(v.Name, "isatap") {
			var network Network
			is[i].Name = v.Name
			is[i].MacAddress = v.HardwareAddr.String()
			for _, ip := range ips {
				if strings.Contains(ip.String(), ".") {
					is[i].Ipv4 = append(is[i].Ipv4, ip.String())
				}
			}
			network.Name = is[i].Name
			network.MACAddress = is[i].MacAddress
			if len(is[i].Ipv4) > 0 {
				network.IP = is[i].Ipv4[0]
			}

			fmt.Printf("network:=%+v \n ", network)
		}

	}

	return nil
}

type DiskInfo struct {
	// 全部字段请看文档,这里我只取了硬盘名字,大小,序列号
	// https://docs.microsoft.com/zh-cn/windows/win32/cimwin32prov/win32-diskdrive#syntax
	Name         string
	Size         uint64
	SerialNumber string
}

func getDiskInfo() {
	var diskInfos []DiskInfo
	err := wmi.Query("Select * from Win32_DiskDrive", &diskInfos)
	if err != nil {
		return
	}
	for _, disk := range diskInfos {
		fmt.Printf("硬盘名称是%s,硬盘大小是%dG,硬盘序列号是%s\n", disk.Name, disk.Size/1024/1024/1024, disk.SerialNumber)
		fmt.Printf("disk:=%+v \n ", disk)

	}
}

type Storage struct {
	Name       string
	FileSystem string
	Total      uint64
	Free       uint64
}

type storageInfo struct {
	Name       string
	Size       uint64
	FreeSpace  uint64
	FileSystem string
}

func getStorageInfo() {
	var storageinfo []storageInfo
	var loaclStorages []Storage
	err := wmi.Query("Select * from Win32_LogicalDisk", &storageinfo)
	if err != nil {
		return
	}

	for _, storage := range storageinfo {
		info := Storage{
			Name:       storage.Name,
			FileSystem: storage.FileSystem,
			Total:      storage.Size / 1024 / 1024 / 1024,
			Free:       storage.FreeSpace / 1024 / 1024 / 1024,
		}
		loaclStorages = append(loaclStorages, info)
	}
	fmt.Printf("localStorages:=%+v \n", loaclStorages)
}

func main() {
	getCPUInfo()
	getOSInfo()
	getMemoryInfo()
	err := getNetworkInfo()
	if err != nil {
		return
	}
	getStorageInfo()
	getDiskInfo()
	for true {
		time.Sleep(1 * time.Second)
	}
}
