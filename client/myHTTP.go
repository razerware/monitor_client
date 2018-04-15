package client

import (
	"io/ioutil"
	"net/http"
	"bytes"
	"github.com/golang/glog"
)

func MyPost(url string, send_bytes []byte) (int, []byte, error) {
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(send_bytes))
	if err != nil {
		// handle error
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		glog.Error(err)
		// handle error
	} else {
		glog.V(1).Info("MyPost successed!")
	}
	code := resp.StatusCode
	return code, body, err
}

func MyGet(url string, queries map[string]string) (int, []byte, error) {
	client := &http.Client{}
	glog.V(1).Info(url)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		// handle error
		glog.Error(err)
	}
	//add queries
	q := req.URL.Query()
	if queries != nil {
		for k, v := range queries {
			q.Add(k, v)
		}
		req.URL.RawQuery = q.Encode()
		glog.Info(req.URL.String())
	}

	//q.Add("fromImage","hello-world")
	//q.Add("tag","latest")

	resp, err := client.Do(req)
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		// handle error
		glog.Error(err)
	} else {
		glog.V(1).Info("MyGet successed!")
	}
	code := resp.StatusCode
	return code, body, err
}
