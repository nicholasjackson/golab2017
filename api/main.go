package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/DataDog/datadog-go/statsd"
)

var statsD *statsd.Client

type kitten struct {
	Name          string `json:"name"`
	Weight        int    `json:"weight"`
	FavouriteFood string `json:"favouriteFood,omitempty"`
}

var kittens = []kitten{
	kitten{
		Name:   "Benny",
		Weight: 23,
	},
	kitten{
		Name:   "Fat Freddy's Cat",
		Weight: 50,
	},
}

func main() {
	setupDependencies()
	statsD.Incr("golab2017.api.start", []string{"golab2017"}, 1)

	http.DefaultServeMux.HandleFunc("/list", handleList)
	http.DefaultServeMux.HandleFunc("/detail", handleDetail)
	http.ListenAndServe(":9090", http.DefaultServeMux)
}

func handleList(rw http.ResponseWriter, r *http.Request) {
	startTime := time.Now()

	err := json.NewEncoder(rw).Encode(kittens)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		statsD.Incr("golab2017.api.list.error", []string{"golab2017"}, 1)
		return
	}

	statsD.Timing("golab2017.api.list.timing", time.Now().Sub(startTime), []string{"golab2017"}, 1)
	statsD.Incr("golab2017.api.list.called", []string{"golab2017"}, 1)
}

func handleDetail(rw http.ResponseWriter, r *http.Request) {
	startTime := time.Now()

	//get currency
	res, err := http.Get("http://currency:9091/currency")
	if err != nil {
		fmt.Println(err)

		rw.WriteHeader(http.StatusInternalServerError)
		statsD.Incr("golab2017.api.detail.called", []string{"golab2017"}, 1)
		return
	}
	defer res.Body.Close()

	err = json.NewEncoder(rw).Encode(kittens[0])
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		statsD.Incr("golab2017.api.detail.called", []string{"golab2017"}, 1)
		return
	}

	statsD.Timing("golab2017.api.detail.timing", time.Now().Sub(startTime), []string{"golab2017"}, 1)
	statsD.Incr("golab2017.api.detail.called", []string{"golab2017"}, 1)
}

func setupDependencies() {
	var err error
	statsD, err = statsd.New("golab2017.demo.gs:9125")
	if err != nil {
		fmt.Println(err)
	}
}
