docker-start:
	docker run -p 6379:6379 redis:latest

run:
	go run main.go

build:
	go build .
