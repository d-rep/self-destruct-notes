# Self Destruct Notes app

This is from [Building a secure note sharing service in Go](https://dusted.codes/building-a-secure-note-sharing-service-in-go)
and source code from https://github.com/dustinmoris/self-destruct-notes

## Setup

You will need [Go](https://go.dev/dl/) and [Docker](https://www.docker.com/) installed to compile and run the project.


## Run

Use docker to start up a Redis instance locally:
```shell
make docker-start
```

Then run the backend web service:
```shell
make run
```

Then the app can be viewed at http://localhost:3000/
