package client

import (
	"net"
	"net/rpc"
	"log"
	"fmt"
	"docker-beego/client"
)

type MonitorStats struct {
	CpuPercent string
	MemPercent string
}

func ConnServer(ip string){
	address, err := net.ResolveTCPAddr("tcp", ip)
	if err != nil {
		panic(err)
	}
	conn, _ := net.DialTCP("tcp", nil, address)
	defer conn.Close()
	client := rpc.NewClient(conn)
	defer client.Close()

	CallRemoteFunc(client)
}

func CallRemoteFunc(client *rpc.Client){

	args:=CollectDatas()
	res:=MonitorStats{}
	funcName:="MonitorStats.Collect"
	call:=client.Go(funcName,args,&res,nil)
	call_res:=<-call.Done
	if call_res.Error != nil {
		log.Fatal("error:", funcName,call_res.Error)
	}
	fmt.Printf("ok:",funcName,res)
}

func CollectDatas() MonitorStats{

	a:=MonitorStats{"7%","8%"}
	return a
}
