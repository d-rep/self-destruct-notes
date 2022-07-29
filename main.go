package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"
)

type Server struct{}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet || r.Method == http.MethodHead {
		s.handleGET(w, r)
		return
	}
	if r.Method == http.MethodPost && r.URL.Path == "/" {
		s.handlePOST(w, r)
		return
	}
	http.NotFound(w, r)
}

func (s *Server) handlePOST(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("You posted to /"))
}

func (s *Server) handleGET(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		s.renderTemplate(w, r, nil, "layout", "dist/layout.html", "dist/index.html")
		return
	}
	noteID := strings.TrimPrefix(r.URL.Path, "/")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("You requested the note with the ID '%s'.", noteID)))
}

func (s *Server) renderTemplate(w http.ResponseWriter, r *http.Request, data interface{}, name string, files ...string) {
	t := template.Must(template.ParseFiles(files...))
	err := t.ExecuteTemplate(w, name, data)
	if err != nil {
		panic(err)
	}
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
