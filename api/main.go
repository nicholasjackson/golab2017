package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/DataDog/datadog-go/statsd"
	"github.com/nicholasjackson/loadbalancer"
)

var statsD *statsd.Client
var client *loadbalancer.Client

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
	http.DefaultTransport.(*http.Transport).MaxIdleConnsPerHost = 250

	setupDependencies()
	statsD.Incr("golab2017.api.start", nil, 1)

	http.DefaultServeMux.HandleFunc("/list", handleList)
	http.DefaultServeMux.HandleFunc("/detail", handleDetail)
	http.ListenAndServe(":9090", http.DefaultServeMux)
}

func handleList(rw http.ResponseWriter, r *http.Request) {
	startTime := time.Now()

	err := json.NewEncoder(rw).Encode(kittens)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		statsD.Incr("golab2017.api.list.error", nil, 1)
		return
	}

	statsD.Timing("golab2017.api.list.timing", time.Now().Sub(startTime), nil, 1)
	statsD.Incr("golab2017.api.list.success", nil, 1)
}

func handleDetail(rw http.ResponseWriter, r *http.Request) {
	startTime := time.Now()

	//get currency
	err := client.Do(func(endpoint url.URL) error {
		res, err := http.Get("http://currency:9091/currency")
		if err != nil {
			return err
		}
		defer res.Body.Close()
		_, _ = ioutil.ReadAll(res.Body)

		return json.NewEncoder(rw).Encode(kittens[0])
	})

	if err != nil {
		errors := err.(loadbalancer.ClientError)
		fmt.Println(err)

		for _, e := range errors.Errors() {
			switch e.Error() {
			case loadbalancer.ErrorTimeout:
				statsD.Incr("golab2017.api.detail.currency.timeout", nil, 1)
			case loadbalancer.ErrorCircuitOpen:
				statsD.Incr("golab2017.api.detail.currency.circuitopen", nil, 1)
			default:
				statsD.Incr("golab2017.api.detail.currency.error", nil, 1)
			}
		}

		statsD.Incr("golab2017.api.detail.error", nil, 1)
	} else {
		statsD.Incr("golab2017.api.detail.success", nil, 1)
	}

	statsD.Timing("golab2017.api.detail.timing", time.Now().Sub(startTime), nil, 1)
}

func setupDependencies() {
	var err error
	statsD, err = statsd.New("statsd:9125")
	if err != nil {
		fmt.Println(err)
	}

	u, _ := url.Parse("http://currency:9091/currency")

	client = loadbalancer.NewClient(
		loadbalancer.Config{
			Timeout:                600 * time.Millisecond,
			MaxConcurrentRequests:  500,
			ErrorPercentThreshold:  50,
			DefaultVolumeThreshold: 20,
			Endpoints:              []url.URL{*u},
			StatsD: loadbalancer.StatsD{
				Enabled: true,
				Server:  "statsd:9125",
				Prefix:  "golab.hystrix",
			},
		},
		&loadbalancer.RandomStrategy{},
		&loadbalancer.ExponentialBackoff{},
	)

}
