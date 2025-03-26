package util

import (
	"archive/zip"
	"io"
)

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
