package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/DataDog/datadog-go/statsd"
)

var totalSleepTime = 0.0
var sleepTime = 0.0

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
	var err error
	sleepTime, err = strconv.ParseFloat(os.Getenv("SLEEP_TIME"), 32)
	if err != nil {
		sleepTime = 0
	}

	fmt.Println("Starting with delay: ", sleepTime)

	setupDependencies()
	statsD.Incr("golab2017.currency.start", nil, 1)

	http.DefaultServeMux.HandleFunc("/currency", handle)
	http.ListenAndServe(":9091", http.DefaultServeMux)
}

func handle(rw http.ResponseWriter, r *http.Request) {
	startTime := time.Now()

	err := json.NewEncoder(rw).Encode(currencies)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		statsD.Incr("golab2017.currency.error", nil, 1)
		return
	}

	time.Sleep(time.Duration(totalSleepTime) * time.Millisecond)

	statsD.Timing("golab2017.currency.timing", time.Now().Sub(startTime), nil, 1)
	statsD.Incr("golab2017.currency.success", nil, 1)

	totalSleepTime += sleepTime
}

func setupDependencies() {
	var err error
	statsD, err = statsd.New("statsd:9125")
	if err != nil {
		fmt.Println(err)
	}
}
