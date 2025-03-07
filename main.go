package main

import (
	"log"
	"net/http"
	"time"

	"freshPaint/rest"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Recoverer)
	r.Use(middleware.Logger)
	r.Use(middleware.Timeout(time.Second * 60))

	r.Post("/", rest.Resize)

	log.Fatal(http.ListenAndServe(":8080", r))
}
