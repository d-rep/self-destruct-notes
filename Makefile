run:
	go run main.go

docker-start:
	docker run -p 6379:6379 redis:latest