package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/DataDog/datadog-go/statsd"
)

var statsD *statsd.Client

type currency struct {
	Name string
}

var currencies = []currency{
	currency{
		Name: "USD",
	},
}

func main() {
	setupDependencies()
	statsD.Incr("golab2017.currency.start", []string{"golab2017"}, 1)

	http.DefaultServeMux.HandleFunc("/currency", handle)
	http.ListenAndServe(":9091", http.DefaultServeMux)
}

func handle(rw http.ResponseWriter, r *http.Request) {
	startTime := time.Now()

	err := json.NewEncoder(rw).Encode(currencies)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		statsD.Incr("golab2017.currency.called", []string{"golab2017"}, 1)
		return
	}

	time.Sleep(100 * time.Millisecond)

	statsD.Timing("golab2017.currency.timing", time.Now().Sub(startTime), []string{"golab2017"}, 1)
	statsD.Incr("golab2017.currency.called", []string{"golab2017"}, 1)
}

func setupDependencies() {
	var err error
	statsD, err = statsd.New("golab2017.demo.gs:9125")
	if err != nil {
		fmt.Println(err)
	}
}
