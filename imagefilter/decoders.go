package imagefilter

import (
	"bufio"
	"errors"
	"image"
	"image/jpeg"
	"image/png"
	"net/http"
	"os"
)

func decodeJPEG(imgFile *os.File) (image.Image, error) {
	img, err := jpeg.Decode(imgFile)
	if err != nil {
		return nil, err
	}
	return img, nil
}

func decodePNG(imgFile *os.File) (image.Image, error) {
	img, err := png.Decode(imgFile)
	if err != nil {
		return nil, err
	}
	return img, nil
}

// DecodeImg open file and decode image using content type
func DecodeImg(imgFilepath string) (image.Image, error) {
	var finalImg image.Image
	imgFile, err := os.Open(imgFilepath)
	if err != nil {
		return nil, err
	}
	defer imgFile.Close()

	buff := bufio.NewReader(imgFile)
	buffType, err := buff.Peek(512)
	if err != nil {
		return nil, err
	}

	imgFile.Seek(0, 0)
	contentType := http.DetectContentType(buffType)
	switch contentType {
	case "image/jpeg":
		finalImg, err = decodeJPEG(imgFile)
	case "image/png":
		finalImg, err = decodePNG(imgFile)
	default:
		finalImg = nil
		err = errors.New("content type not available")
	}

	return finalImg, err

}

func DecodeJPEGByPath(imgFilepath string) (image.Image, error) {
	imgFile, err := os.Open(imgFilepath)
	if err != nil {
		return nil, err
	}
	defer imgFile.Close()

	return decodeJPEG(imgFile)
}

func DecodePNGByPath(imgFilepath string) (image.Image, error) {
	imgFile, err := os.Open(imgFilepath)
	if err != nil {
		return nil, err
	}
	defer imgFile.Close()

	return decodePNG(imgFile)
}
