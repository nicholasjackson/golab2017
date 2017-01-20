# GoLab 2017
This respository contains the example code and simulation from my talk at GoLab2017

## Branches
* master - Simple service with no timeouts, circuit breaking, etc
* circuitbreaker - Implements many of the patterns such as timeouts, backoffs, circuit breaking

## Dependencies
To run the examples in this package you must have Go 1.7 or greater and Docker

## Setup
Running `make run` will start the build the source code and start a docker stack which includes a Grafana server and Prometheus datastore.  
Next we need to configure Grafana to add the Prometheus data source, unfortunately there is no way of doing this from configuration.  Point your webbrowser at http://your-docker-ip:3000 and log in with the default credentials "admin", "admin".  Then you need to click on datasources and add a new source for Prometheus which is located at `http://prometheus:9090`.  Save and test and then we can add the dashboard.
To import the dashboard you can click on dashboards and choose import, select the dashboard which is located at ./golab/grafana/dashboards/golab.json, make sure you select the "Prometheus" data source that you selected earlier.

## Running the simulation
To run the simulation you can run the command `DOCKER=192.168.165.129 make test` from the root of the source code replacing the ip address in the environment variable docker to your docker host's ip address.

If you switch back to Grafana in your browser and look at your dashboard you will see a benchmark output for a normally operating server.

## Simulating failure
To see how the system will operate when there is a failing downstream system we can start our server stack in a different mode.  Stop the existing docker stack you started with `make run` and restart it with `make run_slow`.  This will start the server with a deliberate flaw which makes the currency downstream service run slowly.

If again we run our `make test` and look at our grafana dashboard you will see that the service is now experiencing much less throughput due to the slow dependency in the currency service.

## Fixing things
Once you have seen how things operate in a failure state it is time to check out the `circuitbreaker`branch and re-run these steps to see an example of our simulation which implements circuit breaking.
