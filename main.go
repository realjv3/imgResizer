package main

import (
	"log"
	"net/http"
	"time"

	"github.com/realjv3/imgResizer/rest"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	log.Println("starting server...")

	r := chi.NewRouter()
	r.Use(middleware.Recoverer)
	r.Use(middleware.Logger)
	r.Use(middleware.Timeout(time.Second * 60))

	r.Post("/", rest.ResizeHandler)

	log.Fatal(http.ListenAndServe(":8080", r))
}
