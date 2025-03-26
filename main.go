package main

import (
	"log"
	"net/http"

	"github.com/realjv3/imgResizer/rest"
)

func main() {
	log.Println("starting server...")

	http.HandleFunc("POST /", rest.ResizeHandler)

	log.Println("server is up and running on port 8080")
	log.Println("post images to root path with the desired height & width as url query params...")

	log.Fatal(http.ListenAndServe(":8080", nil))
}
