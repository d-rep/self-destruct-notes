package main

import (
	"log"
	"net/http"
	"os"
)

func webServer(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hello, world!"))
}

func main() {
	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "3000"
	}
	addr := "localhost:" + port
	log.Printf("Starting web server, listening on %s\n", addr)
	err := http.ListenAndServe(addr, http.HandlerFunc(webServer))
	if err != nil {
		panic(err)
	}
}
