package main

import (
	"fmt"
	"os"
	"os/signal"
	"rpc_client/client"
	"time"
	"flag"
)

func main() {
	c := make(chan os.Signal)
	signal.Notify(c)
	hostid:=flag.Int("hostid", 0, "number")
	hostip:=flag.String("hostip", "", "string")
	info:=client.HostInfo{*hostid,*hostip}
	for {
		select {
		case <-c:
			fmt.Println("get signal:", c)
			os.Exit(1)
		default:
			client.CollectData(info)
			time.Sleep(10*time.Second)
		}
	}
}
