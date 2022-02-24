package pixelextract

import (
	"image"
	"image/color"

	colors "gopkg.in/go-playground/colors.v1"
)

type PixelColor struct {
	X             int
	Y             int
	IsDark        bool
	IsLight       bool
	OriginalColor color.Color
	ColorRGBA     color.RGBA
}

func (p PixelColor) ColorWeight() uint32 {
	return uint32(p.ColorRGBA.R + p.ColorRGBA.G + p.ColorRGBA.B + p.ColorRGBA.A)
}

func (p PixelColor) ColorGrayScale() uint32 {
	var fr uint32 = 19595
	var fg uint32 = 38470
	var fb uint32 = 7471
	var cR uint32 = uint32(p.ColorRGBA.R)
	var cG uint32 = uint32(p.ColorRGBA.G)
	var cB uint32 = uint32(p.ColorRGBA.B)
	y := (fr*cR + fg*cG + fb*cB + 1<<15) >> 16 //greyScale 16bit
	return y
}

func ExtractPixelFromImg(img image.Image) []PixelColor {
	var xp []PixelColor
	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			pColor := img.At(x, y)
			RGBAColor := colors.FromStdColor(pColor)
			originalRGBA := color.RGBAModel.Convert(pColor).(color.RGBA)
			xp = append(xp, PixelColor{
				X:             x,
				Y:             y,
				IsDark:        RGBAColor.IsDark(),
				IsLight:       RGBAColor.IsLight(),
				OriginalColor: pColor,
				ColorRGBA:     originalRGBA,
			})
		}
	}
	return xp
}
