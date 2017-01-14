package main

import (
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/nicholasjackson/bench"
	"github.com/nicholasjackson/bench/output"
)

var mutex sync.Mutex
var requestType = 0

func main() {
	mutex = sync.Mutex{}

	fmt.Println("Benchmarking application")

	if _, err := os.Open("./errors.log"); err == nil {
		os.Remove("./errors.log")
	}
	file, _ := os.Create("./errors.log")

	b := bench.New(50, 300*time.Second, 30*time.Second, 30*time.Second)
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
	resp, err := http.Get("http://192.168.165.129:9090/list")
	if resp != nil && resp.Body != nil {
		defer resp.Body.Close()
	}

	if err != nil || resp.StatusCode != 200 {
		return fmt.Errorf("Oops")
	}

	return nil
}

func requestDetail() error {

	resp, err := http.Get("http://192.168.165.129:9090/detail")
	if resp != nil && resp.Body != nil {
		defer resp.Body.Close()
	}

	if err != nil || resp.StatusCode != 200 {
		return fmt.Errorf("Oops")
	}

	return nil
}
