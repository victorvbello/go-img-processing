package experiment

import (
	"image"
	"image/draw"
	"sync"

	"github.com/victorvbello/img-processing/imagetransforms"
)

type resizeImgItem struct {
	index int
	img   image.Image
}

func center(rect image.Rectangle) image.Point {
	return image.Point{(rect.Max.X - rect.Min.X) / 2, (rect.Max.Y - rect.Min.Y) / 2}
}

func ImgInfinite(img image.Image, percentage int) image.Image {
	bounds := img.Bounds()
	resultImg := image.NewRGBA(bounds)
	centerBounds := center(bounds)

	var imgX []image.Image
	var wg sync.WaitGroup
	var countImg int

	imgX = make([]image.Image, int(100/percentage))

	resizeChan := make(chan resizeImgItem)
	draw.Draw(resultImg, bounds, img, image.ZP, draw.Src)
	for i := 1; i < 100; {
		wg.Add(1)
		go func(fact int, imgIndex int) {
			resizeImg := imagetransforms.Resize(img, fact)
			resizeChan <- resizeImgItem{imgIndex, resizeImg}
			wg.Done()
		}(i, countImg)
		i += percentage
		countImg++
	}

	go func() {
		wg.Wait()
		close(resizeChan)
	}()

	for ri := range resizeChan {
		imgX[ri.index] = ri.img
	}
	for _, img := range imgX {
		resizeBounds := img.Bounds()
		newRect := image.Rect(0, 0, resizeBounds.Dx(), resizeBounds.Dy())
		newRect = newRect.Add(centerBounds.Sub(image.Pt(newRect.Dx()/2, newRect.Dy()/2)))

		draw.Draw(resultImg, newRect, img, image.ZP, draw.Src)
	}
	return resultImg
}

func ImgInfiniteSpiral(img image.Image, angle int) image.Image {
	var imgX []image.Image
	var wg sync.WaitGroup
	var count int

	imgX = make([]image.Image, int(100/angle))

	bounds := img.Bounds()
	resultImg := image.NewRGBA(bounds)

	draw.Draw(resultImg, bounds, img, image.ZP, draw.Src)

	resizeChan := make(chan resizeImgItem)

	for i := 1; i < 100; {
		wg.Add(1)
		go func(fact int, xIndex int) {
			resizeImg := imagetransforms.Resize(img, fact)
			resizeChan <- resizeImgItem{xIndex, resizeImg}
			wg.Done()
		}(i, count)
		count++
		i += angle
	}
	go func() {
		wg.Wait()
		close(resizeChan)
	}()

	for ri := range resizeChan {
		imgX[ri.index] = ri.img
	}

	for i, img := range imgX {
		resultImg = imagetransforms.Rotate(resultImg, img, float64(i+angle)).(*image.RGBA)
	}

	return resultImg
}
