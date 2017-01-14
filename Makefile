build:
	CGO_ENABLED=0 GOOS=linux go build -o ./api/service ./api/main.go
	CGO_ENABLED=0 GOOS=linux go build -o ./currency/service ./currency/main.go

run: build
	docker-compose up --build

run_slow: build
	SLEEP_TIME=0.03 docker-compose up --build

test:
	go run ./bench/main.go
