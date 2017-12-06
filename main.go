package main

import (
	"fmt"
	"os"
	"os/signal"
	"rpc_client/client"
)

func main() {
	c := make(chan os.Signal)
	signal.Notify(c)
	info:=client.HostInfo{1,"10.109.252.172"}
	client.CollectData(info)
	for {
		select {
		case <-c:
			fmt.Println("get signal:", c)
			os.Exit(1)
		}
	}
}
