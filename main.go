package main

import (
	"os"
	"os/signal"
	"rpc_client/client"
	"time"
	"flag"
	"fmt"
)

func main() {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, os.Kill)
	hostid := flag.Int("hostid", 0, "number")
	hostip := flag.String("hostip", "", "string")
	flag.Parse()
	info := client.HostInfo{*hostid, *hostip}
	fmt.Println(info)
	t:=make(chan int)
	go func(t chan int ) {
		t<-1
		time.Sleep(10*time.Second)
	}(t)
	for {
		select {
		case <-c:
			os.Exit(1)
		case <-t:
			client.CollectData(info)
		}
	}
}
