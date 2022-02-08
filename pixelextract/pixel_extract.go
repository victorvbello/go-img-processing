package pixelextract

import (
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
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

func (p PixelColor) ColorWeight() uint32 {
	return uint32(p.ColorRGBA.R + p.ColorRGBA.G + p.ColorRGBA.B + p.ColorRGBA.A)
}

func (p PixelColor) ColorGrayScaleWeight() uint32 {
	var fr uint32 = 19595
	var fg uint32 = 38470
	var fb uint32 = 7471
	var cR uint32 = uint32(p.ColorRGBA.R)
	var cG uint32 = uint32(p.ColorRGBA.G)
	var cB uint32 = uint32(p.ColorRGBA.B)
	y := (fr*cR + fg*cG + fb*cB + 1<<15) >> 16 //greyScale 16bit
	return y
}

func ExtractPixelFromFileJPEG(imgFilepath string) ([]PixelColor, image.Image, error) {
	var xp []PixelColor
	imgFile, err := os.Open(imgFilepath)
	if err != nil {
		fmt.Println(err.Error())
		return nil, nil, err
	}
	defer imgFile.Close()

	img, err := jpeg.Decode(imgFile)
	if err != nil {
		fmt.Println(err.Error())
		return nil, nil, err
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
	return xp, img, nil
}

func ExtractPixelFromFilePNG(imgFilepath string) ([]PixelColor, image.Image, error) {
	var xp []PixelColor
	imgFile, err := os.Open(imgFilepath)
	if err != nil {
		fmt.Println(err.Error())
		return nil, nil, err
	}
	defer imgFile.Close()

	img, err := png.Decode(imgFile)
	if err != nil {
		fmt.Println(err.Error())
		return nil, nil, err
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
	return xp, img, nil
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
