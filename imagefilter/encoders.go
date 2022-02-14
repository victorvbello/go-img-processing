package imagefilter

import (
	"errors"
	"image"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
	"strings"
)

func EncodeIMG(img image.Image, fileName string) (*os.File, error) {
	var f *os.File
	var err error

	dir := filepath.Dir(fileName)

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err := os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			return nil, err
		}
	}

	switch strings.TrimPrefix(filepath.Ext(fileName), ".") {
	case "jpeg", "jpg":
		f, err = encodeJPEG(img, fileName)
	case "png":
		f, err = encodePNG(img, fileName)
	default:
		f = nil
		err = errors.New("content type not available")
	}
	return f, err
}

func encodeJPEG(img image.Image, fileName string) (*os.File, error) {
	outFile, err := os.Create(fileName)
	if err != nil {
		return nil, err
	}
	defer outFile.Close()

	if err := jpeg.Encode(outFile, img, nil); err != nil {
		return nil, err
	}
	return outFile, nil
}

func encodePNG(img image.Image, fileName string) (*os.File, error) {
	outFile, err := os.Create(fileName)
	if err != nil {
		return nil, err
	}
	defer outFile.Close()

	if err := png.Encode(outFile, img); err != nil {
		return nil, err
	}
	return outFile, nil
}
