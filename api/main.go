package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/DataDog/datadog-go/statsd"
	"github.com/nicholasjackson/loadbalancer"
)

var statsD *statsd.Client
var client *loadbalancer.Client
var urls []url.URL

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
	http.DefaultTransport.(*http.Transport).MaxIdleConnsPerHost = 500

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
	var err error
	startTime := time.Now()

	if os.Getenv("MODE") == "breaker" {
		err = getCurrencyLB(rw)
	} else {
		err = getCurrency(rw, urls[1])
		fmt.Println(err)
	}

	if err != nil {
		statsD.Incr("golab2017.api.detail.error", nil, 1)
	} else {
		statsD.Incr("golab2017.api.detail.success", nil, 1)
	}

	statsD.Timing("golab2017.api.detail.timing", time.Now().Sub(startTime), nil, 1)
}

func getCurrencyLB(rw http.ResponseWriter) error {
	c := client.Clone() // clone the client

	//get currency
	err := c.Do(func(endpoint url.URL) error {
		return getCurrency(rw, endpoint)
	})

	return err
}

func getCurrency(rw http.ResponseWriter, endpoint url.URL) error {
	res, err := http.Get(endpoint.String())
	if err != nil {
		return err
	}
	defer res.Body.Close()
	_, _ = ioutil.ReadAll(res.Body)

	return json.NewEncoder(rw).Encode(kittens[0])
}

func getURL(uri string) url.URL {
	u, _ := url.Parse(uri)
	return *u
}

func setupDependencies() {
	var err error
	statsD, err = statsd.New("statsd:9125")
	if err != nil {
		fmt.Println(err)
	}

	urls = []url.URL{
		getURL("http://currency:9091/currency"),
		getURL("http://currencyslow:9091/currency"),
	}

	client = loadbalancer.NewClient(
		loadbalancer.Config{
			Retries:                5,
			RetryDelay:             100 * time.Millisecond,
			Timeout:                600 * time.Millisecond,
			MaxConcurrentRequests:  500,
			ErrorPercentThreshold:  50,
			DefaultVolumeThreshold: 1000,
			Endpoints:              urls,
			StatsD: loadbalancer.StatsD{
				Prefix: "golab2017.api.detail.currency",
			},
		},
		&loadbalancer.RoundRobinStrategy{},
		&loadbalancer.ExponentialBackoff{},
	)

	lbStats, _ := loadbalancer.NewDogStatsD(url.URL{Host: "statsd:9125"})
	client.RegisterStats(lbStats)
}
