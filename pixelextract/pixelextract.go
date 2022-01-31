package pixelextract

import (
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"os"

	colors "gopkg.in/go-playground/colors.v1"
)

type PixelColor struct {
	X         int
	Y         int
	IsDark    bool
	IsLight   bool
	ColorRGBA color.RGBA
}

func ExtractPixelFromFile(imgFilepath string) ([]PixelColor, image.Image) {
	var xp []PixelColor
	imgFile, err := os.Open(imgFilepath)
	defer imgFile.Close()
	if err != nil {
		fmt.Println(err.Error())
	}

	img, err := jpeg.Decode(imgFile)
	if err != nil {
		fmt.Println(err.Error())
	}

	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			pColor := img.At(x, y)
			cr, cg, cb, ca := pColor.RGBA()
			RGBAColor := colors.FromStdColor(pColor)
			xp = append(xp, PixelColor{
				X:         x,
				Y:         y,
				ColorRGBA: color.RGBA{uint8(cr), uint8(cg), uint8(cb), uint8(ca)},
				IsDark:    RGBAColor.IsDark(),
				IsLight:   RGBAColor.IsLight(),
			})
		}
	}
	return xp, img
}

func ExtractPixelFromImg(img image.Image) []PixelColor {
	var xp []PixelColor
	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			pColor := img.At(x, y)
			cr, cg, cb, ca := pColor.RGBA()
			RGBAColor := colors.FromStdColor(pColor)
			xp = append(xp, PixelColor{
				X:         x,
				Y:         y,
				ColorRGBA: color.RGBA{uint8(cr), uint8(cg), uint8(cb), uint8(ca)},
				IsDark:    RGBAColor.IsDark(),
				IsLight:   RGBAColor.IsLight(),
			})
		}
	}
	return xp
}
