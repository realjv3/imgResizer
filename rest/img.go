package rest

import (
	"fmt"
	"image/jpeg"
	"net/http"
	"os"
	"strconv"

	"github.com/nfnt/resize"
)

func Resize(w http.ResponseWriter, r *http.Request) {
	// validate e.g. get dimensions
	height := r.URL.Query().Get("height")
	width := r.URL.Query().Get("width")

	if height == "" || width == "" {
		http.Error(w, "missing dimensions", http.StatusBadRequest)
		return
	}

	err := r.ParseMultipartForm(0)
	if err != nil {
		http.Error(w, "error parsing multipart form", http.StatusInternalServerError)
		return
	}

	if r.MultipartForm == nil {
		http.Error(w, "multipart from is nil", http.StatusInternalServerError)
		return

	}
	for _, files := range r.MultipartForm.File {
		for _, file := range files {
			// imagick load image
			f, _ := file.Open()
			defer f.Close()

			img, err := jpeg.Decode(f)
			if err != nil {
				http.Error(w, fmt.Sprintf("error reading into imagick: %v"), http.StatusInternalServerError)
				return
			}

			// resize
			cols, _ := strconv.Atoi(height)
			rows, _ := strconv.Atoi(width)
			resized := resize.Resize(uint(rows), uint(cols), img, resize.Lanczos2)
			out, err := os.Create(file.Filename)
			if err != nil {
				http.Error(w, fmt.Sprintf("error creating file: %v"), http.StatusInternalServerError)
				return
			}
			err = jpeg.Encode(out, resized, nil)
			if err != nil {
				http.Error(w, fmt.Sprintf("error encoding jpeg: %v"), http.StatusInternalServerError)
				return
			}
		}
	}

	w.Write([]byte("Images resized ok."))

	return
}
