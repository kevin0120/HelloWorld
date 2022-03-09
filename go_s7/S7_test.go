package hello

import (
	"strconv"
	"testing"
	"time"

	"fmt"
	"github.com/robinson/gos7"
)

const (
	tcpDevice = "192.168.5.99"
	rack      = 0
	slot      = 1
)

// slot 2 for 300/400, slot 1 for 1200/1500
// conn.initiateConnection({port: 102, host: '192.168.0.2', localTSAP: 0x0100, remoteTSAP: 0x0200, timeout: 8000, doNotOptimize: true}, connected);
// local and remote TSAP can also be directly specified instead. The timeout option specifies the TCP timeout.

func TestS7(t *testing.T) {
	handler := gos7.NewTCPClientHandler(tcpDevice, rack, slot)
	handler.Timeout = 200 * time.Second
	handler.IdleTimeout = 200 * time.Second
	//handler.Logger = log.New(os.Stdout, "tcp: ", log.LstdFlags)
	handler.Connect()
	defer handler.Close()
	client := gos7.NewClient(handler)
	i := 1

	for {
		i++
		//读取plc的状态
		ClientPLCGetStatus(client)
		//读取某一数据区的数据
		ClientReadIntDB(client)
		//修改其中的一部分数据
		ClientWriteIntDB(client, int16(i))
		//读取plc的硬件bool输入
		ClientReadIO(client)
		//控制plc的硬件bool输出
		ClientWriteIO(client,byte(i%256))
	}

}

//ClientReadIO client test read int
func ClientReadIO(client gos7.Client) {
	time.Sleep(1 * time.Second)
	buf := make([]byte, 255)
	err := client.AGReadEB(0,1,buf)
	if err != nil {
		fmt.Println(err)
	}

	var s7 gos7.Helper
	var result byte
	s7.GetValueAt(buf, 0, &result)

	fmt.Printf("读到PLC I0.0-0.7的某段数据为%08b\n", result)
}

//ClientWriteIO client test ClientWriteIO int
func ClientWriteIO(client gos7.Client, value byte) {
	time.Sleep(1 * time.Second)

	buf := make([]byte, 255)

	var helper gos7.Helper
	helper.SetValueAt(buf, 0, value)


	err := client.AGWriteAB(0,1,buf)
	if err != nil {
		fmt.Println(err)
		return
	}


	// result := binary.BigEndian.Uint16(results)

	fmt.Printf("输出PLC Q0.0-Q0.7的某段数据为%08b\n", value)
}



//ClientTestReadIntDB client test read int
func ClientReadIntDB(client gos7.Client) {
	time.Sleep(1 * time.Second)
	address := 2710
	start := 1
	size := 8
	buf := make([]byte, 255)
	err := client.AGReadDB(address, start, size, buf)
	if err != nil {
		fmt.Println(err)
		return
	}
	// result := binary.BigEndian.Uint16(results)
	var s7 gos7.Helper
	var result [8]byte
	s7.GetValueAt(buf, 0, &result)

	fmt.Printf("读到DB2710中的某段数据为%s\n", result)
}

//ClientTestWriteIntDB client test write int
func ClientWriteIntDB(client gos7.Client, value int16) {
	time.Sleep(1 * time.Second)
	address := 2710
	start := 7
	size := 2
	buffer := make([]byte, 255)

	//binary.BigEndian.PutUint16(buffer[0:], uint16(value))
	var helper gos7.Helper
	helper.SetValueAt(buffer, 0, []byte(strconv.Itoa(int(value))))
	err := client.AGWriteDB(address, start, size, buffer)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("写入DB2710中的某段数据为%v\n", value)
}


//ClientPLCGetStatus get PLC status
func ClientPLCGetStatus(client gos7.Client) {
	status, err := client.PLCGetStatus()
	if err != nil {
		fmt.Println(err)
	}
	switch status {
	case 8:
		fmt.Println("PLC running")
	case 4:
		fmt.Println("PLC stop")
	default:
		fmt.Println("PLC unknown")
	} //8=running, 4=stop, 0=unknown
}

