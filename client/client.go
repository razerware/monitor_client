package client

import (
	"fmt"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
	"strconv"
	"encoding/json"
	"github.com/golang/glog"
	"gopkg.in/go-playground/pool.v3"
	"time"
)

var dbUrl = "http://10.109.252.172:8087"
var db = "lzy"
var db_user = "admin"
var db_user_password = "admin"

type HostInfo struct {
	NodeID  string
	HostIP  string
	SwarmID string
	Role string
}

type containerMonitorStats struct {
	HostInfo
	CpuPercent  float64
	MemPercent  float64
	Name        string
	serviceID   string
	serviceName string
}
type vmMonitorStats struct {
	HostInfo
	CpuPercent float64
	MemPercent float64
}

type ContainerStats struct {
	CPUStats struct {
		CPUUsage struct {
			TotalUsage int64 `json:"total_usage"`
		} `json:"cpu_usage"`
		SystemCPUUsage int64 `json:"system_cpu_usage"`
	} `json:"cpu_stats"`
	PrecpuStats struct {
		CPUUsage struct {
			TotalUsage int64 `json:"total_usage"`
		} `json:"cpu_usage"`
		SystemCPUUsage int64 `json:"system_cpu_usage"`
	} `json:"precpu_stats"`
	MemoryStats struct {
		Usage int   `json:"usage"`
		Limit int64 `json:"limit"`
	} `json:"memory_stats"`
}

type Container struct {
	ID    string   `json:"Id"`
	Names []string `json:"Names"`
	Labels struct {
		ComDockerSwarmServiceID   string `json:"com.docker.swarm.service.id"`
		ComDockerSwarmServiceName string `json:"com.docker.swarm.service.name"`
	} `json:"Labels"`
}

func CollectData(info HostInfo) {
	CollectVm(info)
	CollectContainer(info)
	glog.Info("wait 30s")
}

func CollectVm(info HostInfo) {
	glog.Info("Collect data for ",info.HostIP)
	v, _ := mem.VirtualMemory()
	c, _ := cpu.Percent(0, false)
	// almost every return value is a struct
	glog.Info(fmt.Sprintf("VM Data: Total: %v, Free:%v, UsedPercent:%f%%\n", v.Total/1024/1024/1024, v.Free, v.UsedPercent))
	result := vmMonitorStats{info, c[0], v.UsedPercent}
	// convert to JSON. String() is also implemented
	sendVmInfo("vm", result)
	glog.Info("CollectVm successed")
}
func getContainers(info HostInfo) []Container {
	url := "http://" + info.HostIP + ":2375/containers/json"
	code, body, err := MyGet(url, nil)
	if err != nil {
		//	// handle error
		glog.Error("error occur when GET containers,code is", code)
	}
	c := []Container{}
	json.Unmarshal(body, &c)
	return c
}


func CollectContainer(info HostInfo) {
	c := getContainers(info)
	p := pool.NewLimited(100)
	defer p.Close()

	batch := p.Batch()
	go func() {
		for _, i := range c {
			if i.Labels.ComDockerSwarmServiceID == "" {
				continue
			} else {
				batch.Queue(countAndSend(i, info))
			}

		}
		// DO NOT FORGET THIS OR GOROUTINES WILL DEADLOCK
		// if calling Cancel() it calles QueueComplete() internally
		batch.QueueComplete()
	}()
	batch.WaitAll()
	glog.Info("CollectContainer successed")
}

func countAndSend(i Container, info HostInfo) pool.WorkFunc{
	return func(wu pool.WorkUnit) (interface{}, error) {

		//// simulate waiting for something, like TCP connection to be established
		//// or connection from pool grabbed
		//time.Sleep(time.Second * 1)
		glog.Info("Collect data for ",i.Names)
		url := "http://" + info.HostIP + ":2375/containers/" + i.ID + "/stats?stream=false"
		glog.V(1).Info("CollectContainer url is ", url)
		code, body, err := MyGet(url, nil)

		if err != nil {
			//	// handle error
			glog.Error("error occur when GET containers,code is", code)
		}
		cs := ContainerStats{}
		json.Unmarshal(body, &cs)
		json_cs, ok := json.Marshal(cs)
		if ok != nil {
			glog.Error(ok)
		} else {
			glog.V(1).Info("Container stat: ", string(json_cs))
		}

		cpu_percent := float64((cs.CPUStats.CPUUsage.TotalUsage -
			cs.PrecpuStats.CPUUsage.TotalUsage)) /
			float64((cs.CPUStats.SystemCPUUsage - cs.PrecpuStats.SystemCPUUsage)) * 100

		mem_percent := float64(cs.MemoryStats.Usage) / float64(cs.MemoryStats.Limit) * 100

		service_id := i.Labels.ComDockerSwarmServiceID

		service_name := i.Labels.ComDockerSwarmServiceName

		ms := containerMonitorStats{info, cpu_percent,
		mem_percent, i.Names[0], service_id, service_name}

		sendContainerInfo("container", ms)
		json_ms, ok := json.Marshal(ms)
		if ok != nil {
			glog.Error(ok)
		} else {
			glog.V(1).Info("Monitor stat: ", string(json_ms))
		}

		if wu.IsCancelled() {
			// return values not used
			glog.Info("cancelled")
			return nil, nil
		}
		glog.V(1).Info("not cancelled")
		// ready for processing...

		return true, nil // everything ok, send nil, error if not
	}
}
func sendContainerInfo(field string, stat containerMonitorStats) {
	glog.V(1).Info("sendContainerInfo...")
	url := fmt.Sprintf("%s/write?db=%s&u=%s&p=%s", dbUrl, db, db_user, db_user_password)

	//url := dbUrl + "/write?db=" + db + "&u=" + db_user + "&p=" + db_user_password
	tags := fmt.Sprintf("node_id=%s,service_id=%s,service_name=%s",
		stat.NodeID, stat.serviceID, stat.serviceName)

	stat_string := fmt.Sprintf("%s,%s cpu=%s,mem=%s", field, tags,
		strconv.FormatFloat(stat.CpuPercent, 'f', 2, 64),
		strconv.FormatFloat(stat.MemPercent, 'f', 2, 64))

	stat_byte := []byte(stat_string)
	glog.V(1).Info("Container info is:", stat_string)

	code, body, err := MyPost(url, stat_byte)
	if err != nil || code>204{
		// handle error
		glog.Error(err, "Container info send fail,",string(body),"code is ", code)
	} else {
		glog.V(1).Info("Container info send successed ", code)
	}

}
func sendVmInfo(field string, stat vmMonitorStats) {
	glog.V(1).Info("sendVmInfo...")

	url := fmt.Sprintf("%s/write?db=%s&u=%s&p=%s", dbUrl, db, db_user, db_user_password)

	tags := fmt.Sprintf("node_id=%s,swarm_id=%s,role=%s",
		stat.NodeID, stat.SwarmID,stat.Role)
	//tags := "hostid=" + strconv.Itoa(stat.Hostid) + ",swarmid=" + stat.swarmId

	stat_string := fmt.Sprintf("%s,%s cpu=%s,mem=%s", field, tags,
		strconv.FormatFloat(stat.CpuPercent, 'f', 2, 64),
		strconv.FormatFloat(stat.MemPercent, 'f', 2, 64))

	stat_byte := []byte(stat_string)
	glog.V(1).Info("VM info is:", stat_string)
	code, body, err := MyPost(url, stat_byte)

	if err != nil ||code>204{
		// handle error
		glog.Error(err, "Vm info send failed ",string(body),"code is ", code)
	} else {
		glog.V(1).Info("Vm info send successed ", code)
	}
}

func SendEmail(email string) pool.WorkFunc {
	return func(wu pool.WorkUnit) (interface{}, error) {

		// simulate waiting for something, like TCP connection to be established
		// or connection from pool grabbed
		time.Sleep(time.Second * 1)

		if wu.IsCancelled() {
			// return values not used
			return nil, nil
		}

		// ready for processing...
		fmt.Println("gggg")

		return true, nil // everything ok, send nil, error if not
	}
}