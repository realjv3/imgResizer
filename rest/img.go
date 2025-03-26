package rest

import (
	"archive/zip"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"

	"github.com/realjv3/imgResizer/util"
)

func ResizeHandler(w http.ResponseWriter, r *http.Request) {
	// validate e.g. get dimensions
	height := r.URL.Query().Get("height")
	width := r.URL.Query().Get("width")

	if height == "" || width == "" {
		http.Error(w, "missing dimensions", http.StatusBadRequest)
		return
	}

	cols, err := strconv.Atoi(height)
	if err != nil {
		http.Error(w, "error parsing desired height", http.StatusInternalServerError)
		return
	}

	rows, err := strconv.Atoi(width)
	if err != nil {
		http.Error(w, "error parsing desired width", http.StatusInternalServerError)
		return
	}

	err = r.ParseMultipartForm(0)
	if err != nil || r.MultipartForm == nil {
		http.Error(w, "it appears no images were passed in the request", http.StatusBadRequest)
		return
	}

	var resizedFile io.Reader
	var zipWriter *zip.Writer
	var zipOutput *os.File
	zipFiles := len(r.MultipartForm.File) > 1

	if zipFiles {
		zipOutput, err = os.Create("output.zip")
		if err != nil {
			http.Error(w, "error creating output.zip", http.StatusInternalServerError)
			return
		}
		defer zipOutput.Close()

		zipWriter = zip.NewWriter(zipOutput)
		defer zipWriter.Close()
	}

	for _, files := range r.MultipartForm.File {
		for _, file := range files {
			f, _ := file.Open()
			defer f.Close()

			resizedFile, err = util.ResizeFile(f, cols, rows)
			if err != nil {
				http.Error(w, fmt.Sprintf("error resizing file: %w", err), http.StatusInternalServerError)
				return
			}

			if zipFiles {
				err = util.ZipFile(file.Filename, resizedFile, zipWriter)
				if err != nil {
					http.Error(w, fmt.Sprintf("error zipping jpeg: %w", err), http.StatusInternalServerError)
					return
				}
			}
		}
	}

	if zipFiles {
		fi, err := os.Stat("output.zip")
		if err != nil {
			http.Error(w, "error getting output.zip file info", http.StatusInternalServerError)
			return
		}

		if fi.Size() > 0 {
			w.Header().Set("Content-Type", "application/zip")
			w.Header().Set("Content-Disposition", "attachment; filename=\"output.zip\"")
		}
		_, err = io.Copy(w, zipOutput)
		if err != nil {
			http.Error(w, fmt.Sprintf("error returning jpeg: %w", err), http.StatusInternalServerError)
			return
		}
	} else {
		w.Header().Set("Content-Type", "image/jpeg")
		w.Header().Set("Content-Disposition", "attachment; filename=\"resized.jpg\"")
		_, err = io.Copy(w, resizedFile)
	}

	return
}
