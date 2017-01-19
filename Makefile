build:
	CGO_ENABLED=0 GOOS=linux go build -o ./api/service ./api/main.go
	CGO_ENABLED=0 GOOS=linux go build -o ./currency/service ./currency/main.go

run: build
	SLEEP_TIME=0.0 docker-compose up --build > /dev/null

run_slow: build
	SLEEP_TIME=0.015 docker-compose up --build > /dev/null

test:
	go run ./bench/main.go
