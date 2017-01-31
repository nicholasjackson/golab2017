package loadbalancer

import (
	"net/url"
	"time"

	"github.com/DataDog/datadog-go/statsd"
)

// DogStatsD implements a StatsD metrics endpoint for which implements
// the DataDog statsd protocol
type DogStatsD struct {
	client *statsd.Client
}

// NewDogStatsD creates a new implementation of the DogStatsD metrics client
func NewDogStatsD(server url.URL) (*DogStatsD, error) {
	var err error
	c := &DogStatsD{}
	c.client, err = statsd.New(server.String())

	return c, err
}

// Increment sends a statsd message to increment a bucket to datadog
func (d *DogStatsD) Increment(name string, tags []string, rate float64) {
	d.client.Incr(name, tags, rate)
}

// Timing sends the execution time for the given function to statsd
func (d *DogStatsD) Timing(name string, tags []string, duration time.Duration, rate float64) {
	d.client.Timing(name, duration, tags, rate)
}
