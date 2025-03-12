package util

import (
	"archive/zip"
	"bytes"
	"fmt"
	"image/jpeg"
	"io"
	"mime/multipart"

	"github.com/nfnt/resize"
)

func ResizeFile(file multipart.File, height int, width int) (io.Reader, error) {
	img, err := jpeg.Decode(file)
	if err != nil {
		return nil, fmt.Errorf("error reading into imagick: %w", err)
	}

	resized := resize.Resize(uint(width), uint(height), img, resize.Lanczos2)

	var out bytes.Buffer
	err = jpeg.Encode(&out, resized, nil)
	if err != nil {
		return nil, fmt.Errorf("error encoding jpeg: %w", err)
	}

	return &out, nil
}

func ZipFile(filename string, imgBytes io.Reader, zipWriter *zip.Writer) error {
	var header zip.FileHeader
	header.Name = filename
	header.Method = zip.Deflate

	writer, err := zipWriter.CreateHeader(&header)
	if err != nil {
		return err
	}

	_, err = io.Copy(writer, imgBytes)
	if err != nil {
		return err
	}

	return nil
}
