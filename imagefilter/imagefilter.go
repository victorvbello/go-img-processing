package imagefilter

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"log"
	"math"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/nfnt/resize"
	"github.com/victorvbello/img-processing/pixelextract"
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

func (fi *FilterImg) RandomColor(id int) {
	for _, p := range fi.xp {
		cr, cg, cb, ca := p.ColorRGBA.RGBA()
		fi.Set(p.Y, p.X, changeColorRgba(cr, cg, cb, ca, fi.Factor))
	}

	sc := time.Now()
	outFile, _ := os.Create(fmt.Sprintf("%s_random_%d.jpg", fi.alias, fi.Factor))
	ec := time.Since(sc)
	defer outFile.Close()
	s := time.Now()
	jpeg.Encode(outFile, fi, nil)
	e := time.Since(s)
	fi.log <- fmt.Sprintf("random-color, task: %d total create new => %v \t total encode => %v", id, ec, e)
}

func (fi *FilterImg) RandomRed(id int) {
	for _, p := range fi.xp {
		_, cg, cb, ca := p.ColorRGBA.RGBA()
		if (p.X+p.Y)%2 == 0 {
			fi.Set(p.X, p.Y, changeColorRgba(0, cg, cb, ca, fi.Factor))
		}
	}

	sc := time.Now()
	outFile, _ := os.Create(fmt.Sprintf("%s_red_%d.jpg", fi.alias, fi.Factor))
	ec := time.Since(sc)
	defer outFile.Close()
	s := time.Now()
	jpeg.Encode(outFile, fi, nil)
	e := time.Since(s)
	fi.log <- fmt.Sprintf("random-red, task: %d total create new => %v \t total encode => %v", id, ec, e)
}

func (fi *FilterImg) RandomGreen(id int) {
	for _, p := range fi.xp {
		cr, _, cb, ca := p.ColorRGBA.RGBA()
		if (p.X+p.Y)%2 == 0 {
			fi.Set(p.X, p.Y, changeColorRgba(cr, 0, cb, ca, fi.Factor))
		}
	}

	sc := time.Now()
	outFile, _ := os.Create(fmt.Sprintf("%s_green_%d.jpg", fi.alias, fi.Factor))
	ec := time.Since(sc)
	defer outFile.Close()
	s := time.Now()
	jpeg.Encode(outFile, fi, nil)
	e := time.Since(s)
	fi.log <- fmt.Sprintf("random-green, task: %d total create new => %v \t total encode => %v", id, ec, e)
}

func (fi *FilterImg) RandomBlue(id int) {
	for _, p := range fi.xp {
		cr, cg, _, ca := p.ColorRGBA.RGBA()
		if (p.X+p.Y)%2 == 0 {
			fi.Set(p.X, p.Y, changeColorRgba(cr, cg, 0, ca, fi.Factor))
		}
	}

	sc := time.Now()
	outFile, _ := os.Create(fmt.Sprintf("%s_blue_%d.jpg", fi.alias, fi.Factor))
	ec := time.Since(sc)
	defer outFile.Close()
	s := time.Now()
	jpeg.Encode(outFile, fi, nil)
	e := time.Since(s)
	fi.log <- fmt.Sprintf("random-blue, task: %d total create new => %v \t total encode => %v", id, ec, e)
}

func (fi *FilterImg) GreyScale(id int) {
	for _, p := range fi.xp {
		cr, cg, cb, _ := p.ColorRGBA.RGBA()
		avg := 0.2125*float64(cr) + 0.7154*float64(cg) + 0.0721*float64(cb)
		grayColor := color.Gray{uint8(math.Ceil(avg))}
		fi.Set(p.X, p.Y, grayColor)
	}

	sc := time.Now()
	outFile, _ := os.Create(fi.alias + "_grey.jpg")
	ec := time.Since(sc)
	defer outFile.Close()
	s := time.Now()
	jpeg.Encode(outFile, fi, nil)
	e := time.Since(s)
	fi.log <- fmt.Sprintf("random-grey, task: %d total create new => %v \t total encode => %v", id, ec, e)
}

func (fi *FilterImg) ByteScaleTxtFile(id int) string {
	var currentY int
	outFile, _ := os.Create(fi.alias + "_byte.txt")
	s := time.Now()
	for _, p := range fi.xp {
		currentValue := " 0 "
		if p.IsLight {
			currentValue = " 1 "
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

func (fi *FilterImg) MakeFromTxtFile(id int, txtFilePath string) {
	s := time.Now()
	bounds := fi.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y
	basePath := filepath.Dir(txtFilePath)
	txtFileName := filepath.Base(txtFilePath)
	resultFileName := fi.alias + "_byte.jpg"
	cmd := exec.Command("convert",
		"-size", fmt.Sprintf("%dx%d", width*2, (height*2)-(20*height)/100), "xc:white",
		"-font", "Times-Roman",
		"-gravity", "Center",
		"-pointsize", "11",
		"-fill", "black",
		"-annotate", "+15+15",
		"@"+txtFileName,
		resultFileName,
	)

	cmd.Dir = basePath
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		fi.log <- fmt.Sprintf("byte-scale-make, task: %d error %v", id, err)
		return
	}
	e := time.Since(s)
	fi.log <- fmt.Sprintf("byte-scale-make, task: %d total create img => %v", id, e)
	fi.log <- fmt.Sprintf("byte-scale-make, task: %d flow output => [%v]", id, out.String())
	err = os.Remove(txtFileName)
	if err != nil {
		fi.log <- fmt.Sprintf("byte-scale-make, task: %d error %v", id, err)
	}
}
