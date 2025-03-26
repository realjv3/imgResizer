package util

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"io"
	"mime/multipart"

	"golang.org/x/image/draw"
)

func ResizeFile(file multipart.File, height int, width int) (io.Reader, error) {
	img, err := jpeg.Decode(file)
	if err != nil {
		return nil, fmt.Errorf("error reading into imagick: %w", err)
	}

	resized := image.NewRGBA(image.Rect(0, 0, width, height))
	draw.NearestNeighbor.Scale(resized, resized.Bounds(), img, img.Bounds(), draw.Src, nil)

	var out bytes.Buffer
	err = jpeg.Encode(&out, resized, nil)
	if err != nil {
		return nil, fmt.Errorf("error encoding jpeg: %w", err)
	}

	return &out, nil
}
