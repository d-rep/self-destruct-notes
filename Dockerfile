FROM golang:alpine as buildgo
WORKDIR /tmp
RUN mkdir -v ./compile
COPY ./ ./compile
RUN cd ./compile && CGO_ENABLED=0 go build -o app main.go

FROM scratch
COPY --from=buildgo ./tmp/compile/app /app
ENTRYPOINT ["/app"]
