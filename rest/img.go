package rest

import (
	"archive/zip"
	"fmt"
	"image/jpeg"
	"io"
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

	output, err := os.Create("output.zip")
	if err != nil {
		http.Error(w, "error creating output.zip", http.StatusInternalServerError)
		return
	}
	defer output.Close()

	zipWriter := zip.NewWriter(output)
	defer zipWriter.Close()

	for _, files := range r.MultipartForm.File {
		for _, file := range files {
			// imagick load image
			f, _ := file.Open()
			defer f.Close()

			img, err := jpeg.Decode(f)
			if err != nil {
				http.Error(w, fmt.Sprintf("error reading into imagick: %w", err), http.StatusInternalServerError)
				return
			}

			// resize
			cols, _ := strconv.Atoi(height)
			rows, _ := strconv.Atoi(width)
			resized := resize.Resize(uint(rows), uint(cols), img, resize.Lanczos2)
			out, err := os.Create(file.Filename)
			if err != nil {
				http.Error(w, fmt.Sprintf("error creating file: %w", err), http.StatusInternalServerError)
				return
			}

			err = jpeg.Encode(out, resized, nil)
			if err != nil {
				http.Error(w, fmt.Sprintf("error encoding jpeg: %w", err), http.StatusInternalServerError)
				return
			}

			// out needs to be closed and reopened to reset the file pointer back to the beginning of the file
			err = out.Close()
			if err != nil {
				http.Error(w, fmt.Sprintf("error closing temporary jpeg: %w", err), http.StatusInternalServerError)
				return
			}

			out, err = os.Open(file.Filename)
			if err != nil {
				http.Error(w, fmt.Sprintf("error opening temporary jpeg: %w", err), http.StatusInternalServerError)
				return
			}

			err = zipFile(file.Filename, out, zipWriter)
			if err != nil {
				http.Error(w, fmt.Sprintf("error zipping jpeg: %w", err), http.StatusInternalServerError)
				return
			}

			err = os.Remove(file.Filename)
			if err != nil {
				http.Error(w, fmt.Sprintf("error deleting temporary jpeg: %w", err), http.StatusInternalServerError)
				return
			}
		}
	}

	w.Header().Set("Content-Type", "application/zip")
	w.Header().Set("Content-Disposition", "attachment; filename=\"output.zip\"")
	io.Copy(w, output)

	return
}

func zipFile(filename string, file io.Reader, zipWriter *zip.Writer) error {
	fileInfo, err := os.Stat(filename)
	if err != nil {
		return err
	}

	header, err := zip.FileInfoHeader(fileInfo)
	if err != nil {
		return err
	}
	header.Name = filename
	header.Method = zip.Deflate

	writer, err := zipWriter.CreateHeader(header)
	if err != nil {
		return err
	}

	n, err := io.Copy(writer, file)
	if err != nil {
		return err
	}

	if n != fileInfo.Size() {
		return fmt.Errorf("wrote %d bytes instead of expected %d for %s", n, fileInfo.Size(), filename)
	}

	return nil
}
