build:
	CGO_ENABLED=0 GOOS=linux go build -o ./api/service ./api/main.go
	CGO_ENABLED=0 GOOS=linux go build -o ./currency/service ./currency/main.go

run_single: build
	SLEEP_TIME=0.0 MODE=normal SERVERS=1 docker-compose up --build > /dev/null

run_single_slow: build
	SLEEP_TIME=0.015 MODE=normal SERVERS=1 docker-compose up --build > /dev/null

run_single_cb: build
	SLEEP_TIME=0.0 MODE=breaker SERVERS=1 docker-compose up --build > /dev/null

run_single_cb_slow: build
	SLEEP_TIME=0.015 MODE=breaker SERVERS=1 docker-compose up --build > /dev/null

run_cluster_cb: build
	SLEEP_TIME=0.0 MODE=breaker SERVERS=2 docker-compose up --build > /dev/null

run_cluster_cb_slow: build
	SLEEP_TIME=0.015 MODE=breaker SERVERS=2 docker-compose up --build > /dev/null

test:
	go run ./bench/main.go
