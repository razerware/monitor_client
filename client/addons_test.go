package client

import (
	"testing"
	"flag"
	"os"
)

func TestMain(m *testing.M) {
	flag.Set("alsologtostderr", "true")
	//flag.Set("log_dir", "/tmp")
	flag.Set("v", "1")
	flag.Parse()
	MysqlConnectTest()
	ret := m.Run()
	os.Exit(ret)

}
func TestGetInternal(t *testing.T) {
	GetInternal()
}
