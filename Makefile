# start up Redis for data storage
docker-start-redis:
	docker run -p 6379:6379 redis:latest

# compile webapp itself into a dedicated docker image, and then run it
docker-start-web:
	docker build --tag self-destruct-notes .
	docker run -p 3000:3000 -e REDIS_URL="redis://host.docker.internal:6379/1" self-destruct-notes

run:
	go run main.go

build:
	go build .
