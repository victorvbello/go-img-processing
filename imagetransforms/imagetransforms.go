package imagetransforms

import (
	"image"
	"image/color"
	"math"

	"golang.org/x/image/draw"
	"golang.org/x/image/math/f64"

	utilsImage "github.com/victorvbello/img-processing/utils/image"
)

func calculateNewBounds(aff f64.Aff3, b image.Rectangle) image.Rectangle {
	points := [...]image.Point{
		{b.Min.X, b.Min.Y},
		{b.Max.X - 1, b.Min.Y},
		{b.Min.X, b.Max.Y - 1},
		{b.Max.X - 1, b.Max.Y - 1},
	}
	var min, max image.Point
	for i, p := range points {
		x0 := float64(p.X) + 0.5
		y0 := float64(p.Y) + 0.5
		x := aff[0]*x0 + aff[1]*y0 + aff[2]
		y := aff[3]*x0 + aff[4]*y0 + aff[5]
		pMin := image.Point{int(math.Floor(x)), int(math.Floor(y))}
		pMax := image.Point{int(math.Ceil(x)), int(math.Ceil(y))}
		if i == 0 {
			min = pMin
			max = pMax
			continue
		}
		if min.X > pMin.X {
			min.X = pMin.X
		}
		if min.Y > pMin.Y {
			min.Y = pMin.Y
		}
		if max.X < pMax.X {
			max.X = pMax.X
		}
		if max.Y < pMax.Y {
			max.Y = pMax.Y
		}

	}
	return image.Rectangle{Min: min, Max: max}
}

func Transparency(img image.Image, alpha uint8) image.Image {
	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y
	resultImg := image.NewRGBA(bounds)
	mask := image.NewAlpha(bounds)
	bg := image.NewRGBA(bounds)
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			bg.Set(x, y, image.White)
			mask.SetAlpha(x, y, color.Alpha{alpha})
		}
	}
	draw.Draw(resultImg, bounds, bg, image.ZP, draw.Src)
	draw.DrawMask(resultImg, bounds, img, image.ZP, mask, image.ZP, draw.Over)
	return resultImg
}

func Rotate(baseImg image.Image, imgToRotate image.Image, angle float64) image.Image {
	angleSin, angleCos := math.Sincos(math.Pi * angle / 180)

	xf, yf := float64(1), float64(1)

	matrix := f64.Aff3{
		angleCos, -angleSin, xf - xf*angleCos + yf*angleSin,
		angleSin, angleCos, yf - xf*angleSin - yf*angleCos,
	}

	copyBaseImg := image.NewRGBA(baseImg.Bounds())
	baseImgCenter := utilsImage.GetCenterPoint(copyBaseImg.Bounds())
	draw.Copy(copyBaseImg, image.ZP, baseImg, baseImg.Bounds(), draw.Over, nil)

	resizeBounds := calculateNewBounds(matrix, imgToRotate.Bounds())
	rotateImg := image.NewRGBA(resizeBounds)

	centerBounds := resizeBounds.Add(baseImgCenter.Sub(utilsImage.GetCenterPoint(resizeBounds)))

	draw.ApproxBiLinear.Transform(rotateImg, matrix, imgToRotate, imgToRotate.Bounds(), draw.Src, nil)
	draw.Draw(copyBaseImg, centerBounds, rotateImg, rotateImg.Bounds().Min, draw.Over)
	return copyBaseImg
}

func Resize(img image.Image, scale int) image.Image {
	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y
	newWidth := width - (scale*width)/100
	newHeight := height - (scale*height)/100
	newImg := image.NewRGBA(image.Rect(0, 0, newWidth, newHeight))
	draw.ApproxBiLinear.Scale(newImg, newImg.Rect, img, bounds, draw.Over, nil)
	return newImg
}
