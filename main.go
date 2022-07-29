package main

import (
	"context"
	"embed"
	_ "embed"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/go-redis/cache/v8"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
)

type Server struct {
	RedisCache *cache.Cache
	BaseURL    string
}

var (
	//go:embed dist/*.html
	layoutFiles embed.FS
)

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

type Note struct {
	Data     []byte
	Destruct bool
}

func (s *Server) renderMessage(
	w http.ResponseWriter,
	r *http.Request,
	title string,
	paragraphs ...interface{},
) {
	s.renderTemplate(
		w, r,
		struct {
			Title      string
			Paragraphs []interface{}
		}{
			Title:      title,
			Paragraphs: paragraphs,
		},
		"layout",
		"dist/layout.html",
		"dist/message.html",
	)
}

func (s *Server) handlePOST(w http.ResponseWriter, r *http.Request) {
	mediaType := r.Header.Get("Content-Type")
	if mediaType != "application/x-www-form-urlencoded" {
		s.badRequest(
			w, r,
			http.StatusUnsupportedMediaType,
			"Invalid media type posted.")
		return
	}
	err := r.ParseForm()
	if err != nil {
		s.badRequest(
			w, r,
			http.StatusBadRequest,
			"Invalid form data posted.")
		return
	}
	form := r.PostForm

	message := form.Get("message")
	destruct := false
	ttl := time.Hour * 24
	if form.Get("ttl") == "untilRead" {
		destruct = true
		ttl = ttl * 365
	}
	note := &Note{
		Data:     []byte(message),
		Destruct: destruct,
	}
	key := uuid.NewString()
	err = s.RedisCache.Set(
		&cache.Item{
			Ctx:            r.Context(),
			Key:            key,
			Value:          note,
			TTL:            ttl,
			SkipLocalCache: true,
		})
	if err != nil {
		log.Printf("could not write to redis cache: %s", err)
		s.serverError(w, r)
		return
	}

	noteURL := fmt.Sprintf("%s/%s", s.BaseURL, key)
	w.WriteHeader(http.StatusOK)
	s.renderMessage(
		w, r,
		"Note was successfully created",
		template.HTML(
			fmt.Sprintf("<a href='%s'>%s</a>", noteURL, noteURL)))
}

func (s *Server) handleGET(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		s.renderTemplate(w, r, nil, "layout", "dist/layout.html", "dist/index.html")
		return
	}
	noteID := strings.TrimPrefix(r.URL.Path, "/")
	ctx := r.Context()
	note := &Note{}
	err := s.RedisCache.GetSkippingLocalCache(
		ctx,
		noteID,
		note)
	if err != nil {
		s.badRequest(
			w, r,
			http.StatusNotFound,
			fmt.Sprintf("Note with ID %s does not exist.", noteID))
		return
	}
	if note.Destruct {
		err := s.RedisCache.Delete(ctx, noteID)
		if err != nil {
			log.Printf("could not delete noteID %q from redis: %s", noteID, err)
			s.serverError(w, r)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
	log.Printf("You requested the note with the ID '%s'.", noteID)

	w.Write(note.Data)
}

func (s *Server) badRequest(
	w http.ResponseWriter,
	r *http.Request,
	statusCode int,
	message string,
) {
	w.WriteHeader(statusCode)
	w.Write([]byte(message))
}

func (s *Server) serverError(
	w http.ResponseWriter,
	r *http.Request,
) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte("Oops something went wrong. Please check the server logs."))
}

func (s *Server) renderTemplate(w http.ResponseWriter, r *http.Request, data interface{}, name string, files ...string) {
	t := template.Must(template.ParseFS(layoutFiles, files...))
	err := t.ExecuteTemplate(w, name, data)
	if err != nil {
		panic(err)
	}
}

const (
	envRedis        = "REDIS_URL"
	envPort         = "PORT"
	defaultRedisURL = "redis://:6379/1"
	defaultPort     = "3000"
)

func main() {
	redisURL := os.Getenv(envRedis)
	if redisURL == "" {
		redisURL = defaultRedisURL
	}
	redisOptions, err := redis.ParseURL(redisURL)
	if err != nil {
		log.Fatalf("could not parse Redis URL: %s", err)
	}
	redisClient := redis.NewClient(redisOptions)
	defer redisClient.Close()
	redisCache := cache.New(&cache.Options{Redis: redisClient})
	port := os.Getenv(envPort)
	if len(port) == 0 {
		port = defaultPort
	}
	addr := ":" + port
	baseURL := os.Getenv("BASE_URL")
	if len(baseURL) == 0 {
		baseURL = fmt.Sprintf("http://localhost:%s", port)
	}
	server := &Server{
		RedisCache: redisCache,
		BaseURL:    baseURL,
	}
	log.Printf("Starting web server, listening on %s\n", baseURL)
	pong, err := redisClient.Ping(context.Background()).Result()
	if err != nil {
		log.Printf("could not ping redis server %q: %s", redisURL, err)
	} else {
		log.Printf("redis client replied to ping with %q", pong)
	}
	err = http.ListenAndServe(addr, server)
	if err != nil {
		panic(err)
	}
}
