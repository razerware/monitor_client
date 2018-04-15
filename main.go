package main

import (
	"os"
	"os/signal"
	"flag"
	"time"
	"github.com/razerware/monitor_client/client"
	"github.com/golang/glog"
)

func main() {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, os.Kill)
	//hostid := flag.Int("hostid", 0, "number")
	//hostip := flag.String("hostip", "", "string")
	//swarmid := flag.String("swarmid", "", "string")
	flag.Parse()
	glog.Flush()
	client.MysqlConnect()
	hostid,hostip,swarmid:=client.GetInternal()
	info := client.HostInfo{hostid, hostip, swarmid}
	glog.Info("Data collecting start...", info)
	t := make(chan int)
	go func(t chan int) {
		for {
			t <- 1
			duration := 30 * time.Second
			time.Sleep(duration)
		}

	}(t)
	for {
		select {
		case <-c:
			glog.Exit("program exit")
			os.Exit(1)
		case <-t:
			glog.Info("Collecting data...")
			//client.CollectData(info)
		}
	}
}
