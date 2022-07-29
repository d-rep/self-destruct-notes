package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
)

type Server struct{}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet || r.Method == http.MethodHead {
		noteID := strings.TrimPrefix(r.URL.Path, "/")
		if noteID != "" {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(fmt.Sprintf("You requested the note with ID %q", noteID)))
			return
		}
	}

	if r.Method == http.MethodPost && r.URL.Path == "/" {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("You posted to /"))
		return
	}

	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte("Not Found"))
}

func main() {
	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "3000"
	}
	addr := "localhost:" + port
	log.Printf("Starting web server, listening on %s\n", addr)
	err := http.ListenAndServe(addr, &Server{})
	if err != nil {
		panic(err)
	}
}
