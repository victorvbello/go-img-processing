package imagefilter

import (
	"bufio"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"log"
	"math"
	"os"
	"path/filepath"
	"time"

	"github.com/nfnt/resize"
	"github.com/victorvbello/img-processing/pixelextract"
	"golang.org/x/image/font"
	"golang.org/x/image/font/inconsolata"
	"golang.org/x/image/math/fixed"
)

type FilterImg struct {
	alias string
	image.Image
	custom map[image.Point]color.Color
	xp     []pixelextract.PixelColor
	Factor uint32
	log    chan<- string
}

func changeColorRgba(r, g, b, a, f uint32) color.RGBA {
	var safeR uint32 = 0
	var safeG uint32 = 0
	var safeB uint32 = 0
	var safeA uint32 = 0
	if r == 0 {
		safeR = 255
	}
	if g == 0 {
		safeG = 255
	}
	if b == 0 {
		safeB = 255
	}
	if newR := r - f; newR > 0 {
		safeR = newR
	}
	if newG := g - f; newG > 0 {
		safeG = newG
	}
	if newB := b - f; newB > 0 {
		safeB = newB
	}
	if newA := a - f; newA > 0 {
		safeA = newA
	}
	return color.RGBA{uint8(safeR), uint8(safeG), uint8(safeB), uint8(safeA)}
}

func NewFilterImg(a string, img image.Image, xp []pixelextract.PixelColor, f uint32, log chan<- string) *FilterImg {
	return &FilterImg{a, img, map[image.Point]color.Color{}, xp, f, log}
}

func (fi *FilterImg) Set(x, y int, c color.Color) {
	fi.custom[image.Point{x, y}] = c
}

func (fi *FilterImg) At(x, y int) color.Color {
	if c := fi.custom[image.Point{x, y}]; c != nil {
		return c
	}
	return fi.Image.At(x, y)
}

func (fi *FilterImg) GetAlias() string {
	return fi.alias
}

func (fi *FilterImg) GetXp() []pixelextract.PixelColor {
	return fi.xp
}

func (fi *FilterImg) AddLog(l string) {
	fi.log <- l
}

func (fi *FilterImg) RandomColor(id int) image.Image {
	for _, p := range fi.xp {
		cr, cg, cb, ca := p.ColorRGBA.RGBA()
		fi.Set(p.X, p.Y, changeColorRgba(cr, cg, cb, ca, fi.Factor))
	}

	return fi
}

func (fi *FilterImg) RandomRed(id int) image.Image {
	for _, p := range fi.xp {
		_, cg, cb, ca := p.ColorRGBA.RGBA()
		if (p.X+p.Y)%2 == 0 {
			fi.Set(p.X, p.Y, changeColorRgba(0, cg, cb, ca, fi.Factor))
		}
	}

	return fi
}

func (fi *FilterImg) RandomGreen(id int) image.Image {
	for _, p := range fi.xp {
		cr, _, cb, ca := p.ColorRGBA.RGBA()
		if (p.X+p.Y)%2 == 0 {
			fi.Set(p.X, p.Y, changeColorRgba(cr, 0, cb, ca, fi.Factor))
		}
	}

	return fi
}

func (fi *FilterImg) RandomBlue(id int) image.Image {
	for _, p := range fi.xp {
		cr, cg, _, ca := p.ColorRGBA.RGBA()
		if (p.X+p.Y)%2 == 0 {
			fi.Set(p.X, p.Y, changeColorRgba(cr, cg, 0, ca, fi.Factor))
		}
	}

	return fi
}

func (fi *FilterImg) GreyScale(id int) image.Image {
	for _, p := range fi.xp {
		cr, cg, cb, _ := p.ColorRGBA.RGBA()
		avg := 0.2125*float64(cr) + 0.7154*float64(cg) + 0.0721*float64(cb)
		grayColor := color.Gray{uint8(math.Ceil(avg))}
		fi.Set(p.X, p.Y, grayColor)
	}
	return fi
}

func (fi *FilterImg) ByteScaleTxtFile(id int) string {
	var currentY int = -1
	outFile, _ := os.Create(fi.alias + "_byte.txt")
	s := time.Now()
	for _, p := range fi.xp {
		currentValue := "0"
		if p.IsLight {
			currentValue = "1"
		}
		if currentY != p.Y {
			currentY = p.Y
			if _, err := outFile.Write([]byte("\n" + currentValue)); err != nil {
				log.Fatal(err)
				fi.log <- err.Error()
			}
			continue
		}
		if _, err := outFile.Write([]byte(currentValue)); err != nil {
			log.Fatal(err)
			fi.log <- err.Error()
		}
	}
	outFile.Close()
	e := time.Since(s)
	fi.log <- fmt.Sprintf("byte-scale, task: %d total create txt => %v", id, e)
	path, err := os.Getwd()
	if err != nil {
		fi.log <- fmt.Sprintf("byte-scale, task: %d error %v", id, err)
	}
	return filepath.Join(path, outFile.Name())
}

func (fi *FilterImg) Resize(id int, scale uint) *FilterImg {
	bounds := fi.Bounds()
	width := uint(bounds.Max.X)
	newWidth := width - (scale*width)/100
	s := time.Now()
	newImg := resize.Resize(newWidth, 0, fi, resize.Lanczos3)
	newXp := pixelextract.ExtractPixelFromImg(newImg)
	newFi := NewFilterImg(fi.alias, newImg, newXp, fi.Factor, fi.log)
	e := time.Since(s)
	fi.log <- fmt.Sprintf("resize, task: %d total resize => %v", id, e)
	return newFi
}

func (fi *FilterImg) MakeFromTxtFile(txtFilePath string) (image.Image, error) {
	bounds := fi.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y
	img := image.NewRGBA(image.Rect(0, 0, width+((20*width)/100), height+((24*height)/100)))

	draw.Draw(img, img.Bounds(), &image.Uniform{color.White}, image.ZP, draw.Src)

	x, y := 10, 10
	col := color.Black

	f, err := os.Open(txtFilePath)
	if err != nil {
		return nil, err
	}

	scanner := bufio.NewScanner(f)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		point := fixed.Point26_6{X: fixed.Int26_6(x), Y: fixed.Int26_6(y * 44)}
		d := font.Drawer{
			Dst:  img,
			Src:  image.NewUniform(col),
			Face: inconsolata.Bold8x16,
			Dot:  point,
		}
		d.DrawString(scanner.Text())
		y += 12
	}
	f.Close()
	err = os.Remove(txtFilePath)
	if err != nil {
		return nil, err
	}
	return img, nil
}

func (fi *FilterImg) Transparency(alpha uint8) image.Image {
	bounds := fi.Bounds()
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
	draw.DrawMask(resultImg, bounds, fi, image.ZP, mask, image.ZP, draw.Over)
	fi.Image = resultImg
	return resultImg
}
