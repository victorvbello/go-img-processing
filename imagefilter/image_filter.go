package imagefilter

import (
	"bufio"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"log"
	"os"
	"path/filepath"
	"time"

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
	Factor uint8
	log    chan<- string
}

func changeColorRgba(r, g, b, a, f uint8) color.RGBA {
	var safeR uint8 = 0
	var safeG uint8 = 0
	var safeB uint8 = 0
	var safeA uint8 = 0
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
	return color.RGBA{safeR, safeG, safeB, safeA}
}

func NewFilterImg(a string, img image.Image, xp []pixelextract.PixelColor, f uint8, log chan<- string) *FilterImg {
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

func (fi *FilterImg) SetXp(nXp []pixelextract.PixelColor) {
	fi.xp = nXp
}

func (fi *FilterImg) SetImg(nImg image.Image) {
	fi.Image = nImg
}

func (fi *FilterImg) AddLog(l string) {
	fi.log <- l
}

func (fi *FilterImg) RandomColor(id int) image.Image {
	for _, p := range fi.xp {
		original, ok := color.RGBAModel.Convert(p.OriginalColor).(color.RGBA)
		if ok {
			fi.Set(p.X, p.Y, changeColorRgba(original.A, original.B, original.G, original.A, fi.Factor))
		}
	}

	return fi
}

func (fi *FilterImg) RandomRed(id int) image.Image {
	for _, p := range fi.xp {
		if (p.X+p.Y)%2 == 0 {
			fi.Set(p.X, p.Y, changeColorRgba(0, p.ColorRGBA.G, p.ColorRGBA.B, p.ColorRGBA.A, fi.Factor))
		}
	}

	return fi
}

func (fi *FilterImg) RandomGreen(id int) image.Image {
	for _, p := range fi.xp {
		if (p.X+p.Y)%2 == 0 {
			fi.Set(p.X, p.Y, changeColorRgba(p.ColorRGBA.R, 0, p.ColorRGBA.B, p.ColorRGBA.A, fi.Factor))
		}
	}

	return fi
}

func (fi *FilterImg) RandomBlue(id int) image.Image {
	for _, p := range fi.xp {
		if (p.X+p.Y)%2 == 0 {
			fi.Set(p.X, p.Y, changeColorRgba(p.ColorRGBA.R, p.ColorRGBA.G, 0, p.ColorRGBA.A, fi.Factor))
		}
	}

	return fi
}

func (fi *FilterImg) GreyScale(id int) image.Image {
	for _, p := range fi.xp {
		originalColor, ok := color.RGBAModel.Convert(p.OriginalColor).(color.RGBA)
		if !ok {
			fmt.Println("type conversion went wrong")
		}
		grey := uint8(float64(originalColor.R)*0.21 + float64(originalColor.G)*0.72 + float64(originalColor.B)*0.07)
		grayColor := color.RGBA{
			grey,
			grey,
			grey,
			originalColor.A,
		}
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
