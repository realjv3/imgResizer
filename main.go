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

	log.Println("server is up and running on port 8080")
	log.Println("post images to root path with the desired height & width as url query params...")

	log.Fatal(http.ListenAndServe(":8080", r))
}
