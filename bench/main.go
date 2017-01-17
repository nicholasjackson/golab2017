package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/nicholasjackson/bench"
	"github.com/nicholasjackson/bench/output"
)

var mutex sync.Mutex
var requestType = 0

const (
	connections = 250
	duration    = 300 * time.Second
	rampUpTime  = 30 * time.Second
	timeout     = 30 * time.Second
)

func main() {
	mutex = sync.Mutex{}

	fmt.Println("Benchmarking application")

	if _, err := os.Open("./errors.log"); err == nil {
		os.Remove("./errors.log")
	}
	file, _ := os.Create("./errors.log")

	// set max idle connections to be equal to threads
	http.DefaultTransport.(*http.Transport).MaxIdleConnsPerHost = connections

	b := bench.New(connections, duration, rampUpTime, timeout)
	b.AddOutput(0*time.Second, file, output.WriteErrorLogs)
	b.AddOutput(0*time.Second, os.Stdout, output.WriteTabularData)
	b.RunBenchmarks(request)
}

func request() error {
	mutex.Lock()
	local := requestType
	if requestType == 0 {
		requestType = 1
	} else {
		requestType = 0
	}
	mutex.Unlock()

	if local == 0 {
		return requestList()
	}

	return requestDetail()
}

func requestList() error {
	return httpGet("http://192.168.165.129:9090/list")
}

func requestDetail() error {
	return httpGet("http://192.168.165.129:9090/detail")
}

func httpGet(uri string) error {
	resp, err := http.Get(uri)
	if resp != nil && resp.Body != nil {
		defer resp.Body.Close()
		_, _ = ioutil.ReadAll(resp.Body)
	}

	if err != nil {
		return fmt.Errorf("Error in requestDetail: %v", err)
	} else if resp != nil && resp.StatusCode != 200 {
		return fmt.Errorf("Error in requestDetail Status code: %v", resp.StatusCode)
	}

	return nil
}
