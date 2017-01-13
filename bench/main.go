package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/nicholasjackson/bench"
	"github.com/nicholasjackson/bench/output"
)

func main() {

	fmt.Println("Benchmarking application")

	b := bench.New(700, 300*time.Second, 30*time.Second, 30*time.Second)
	b.AddOutput(0*time.Second, os.Stdout, output.WriteTabularData)
	b.RunBenchmarks(request)
}

func request() error {
	resp, err := http.Get("http://localhost:9090/hello")
	if resp != nil && resp.Body != nil {
		defer resp.Body.Close()
	}

	if err != nil || resp.StatusCode != 200 {
		fmt.Println(resp.StatusCode)
		return fmt.Errorf("Oops")
	}

	return nil
}
