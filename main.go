package main

import (
	"fmt"
	"os"
	"os/signal"
	"rpc_client/client"
	"time"
)

func main() {
	c := make(chan os.Signal)
	signal.Notify(c)
	info:=client.HostInfo{1,"10.109.252.172"}
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
