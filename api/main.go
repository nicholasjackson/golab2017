package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/DataDog/datadog-go/statsd"
)

var statsD *statsd.Client

type Kitten struct {
	Name          string `json:"name"`
	Weight        int    `json:"weight"`
	FavouriteFood string `json:"favouriteFood,omitempty"`
}

var kittens []Kitten = []Kitten{
	Kitten{
		Name:   "Benny",
		Weight: 23,
	},
	Kitten{
		Name:   "Fat Freddy's Cat",
		Weight: 50,
	},
}

func main() {
	setupDependencies()
	statsD.Incr("golab2017.api.start", []string{"golab2017"}, 1)

	http.DefaultServeMux.HandleFunc("/hello", handle)
	http.ListenAndServe(":9090", http.DefaultServeMux)
}

func handle(rw http.ResponseWriter, r *http.Request) {
	startTime := time.Now()

	err := json.NewEncoder(rw).Encode(kittens)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		statsD.Incr("golab2017.api.called", []string{"golab2017"}, 1)
		return
	}

	statsD.Timing("golab2017.api.timing", time.Now().Sub(startTime), []string{"golab2017"}, 1)
	statsD.Incr("golab2017.api.called", []string{"golab2017"}, 1)
}

func setupDependencies() {
	var err error
	statsD, err = statsd.New("golab2017.demo.gs:9125")
	if err != nil {
		fmt.Println(err)
	}
}
