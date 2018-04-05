package main

import (
	"os"
	"os/signal"
	"rpc_client/client"
	"time"
	"flag"
	"github.com/golang/glog"
)

func main() {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, os.Kill)
	hostid := flag.Int("hostid", 0, "number")
	hostip := flag.String("hostip", "", "string")
	swarmid := flag.String("swarmid", "", "string")
	flag.Parse()
	glog.Flush()
	info := client.HostInfo{*hostid, *hostip, *swarmid}
	glog.Info("Data collecting start...", info)
	t := make(chan int)
	go func(t chan int) {
		for {
			t <- 1
			duration := 10 * time.Second
			time.Sleep(duration)
		}

	}(t)
	for {
		select {
		case <-c:
			os.Exit(1)
		case <-t:
			glog.Info("Collecting data...")
			client.CollectData(info)
		}
	}
}
