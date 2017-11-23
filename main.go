package main

import (
	"fmt"
	"os"
	"os/signal"
	"rpc_client/client"
)

func main(){
	c := make(chan os.Signal)
	signal.Notify(c)
	client:=client.ConnServer("127.0.0.1:4200")
	client.CallRemoteFunc()
	for {
		select {
		case <-c:
			fmt.Println("get signal:",c)
			os.Exit(1)
		}
	}
}