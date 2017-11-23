package client

import (
	"net"
	"net/rpc"
	"log"
	"fmt"
	"time"
	"os"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
)

type MonitorStats struct {
	Hostid	int
	Hostip	string
	CpuPercent float64
	MemPercent float64
}

func ConnServer(ip string, stopch chan os.Signal) {
	address, err := net.ResolveTCPAddr("tcp", ip)
	if err != nil {
		panic(err)
	}
	conn, _ := net.DialTCP("tcp", nil, address)
	client := rpc.NewClient(conn)
	CallRemoteFunc(conn, client)
}

func CallRemoteFunc(conn *net.TCPConn, client *rpc.Client) {
	defer conn.Close()
	defer client.Close()
	for {
		args := CollectData()
		var res bool
		funcName := "MonitorStats.Collect"
		err:=client.Call(funcName, args, &res)
		if err!=nil{
			log.Fatal("error:", funcName, err)
		}else{
			fmt.Printf("ok:%s%v \n", funcName, res)
		}

		time.Sleep(1*time.Second)
	}

}

func CollectData() MonitorStats {
	v, _ := mem.VirtualMemory()
	c, _ := cpu.Percent(0, false)
	// almost every return value is a struct
	fmt.Printf("Total: %v, Free:%v, UsedPercent:%f%%\n", v.Total/1024/1024/1024, v.Free, v.UsedPercent)
	result := MonitorStats{0,"127.0.0.1",c[0],v.UsedPercent}
	// convert to JSON. String() is also implemented
	fmt.Println(result)
	return result
}
