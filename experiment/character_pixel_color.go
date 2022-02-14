package experiment

import (
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	"log"
	"os"
	"path/filepath"
	"sort"
	"sync"
	"time"

	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"

	"github.com/victorvbello/img-processing/imagefilter"
	"github.com/victorvbello/img-processing/pixelextract"
)

type CharacterInfo struct {
	StrBase       string
	ColorFactor   float32
	MaxColorValue int
	BaseWidth     int
	BaseHeight    int
	CharacterData []CharacterMetadata
}
type CharacterMetadata struct {
	Char           string
	Filename       string
	GrayPercentage float32
}

type byGrayPercentage []CharacterMetadata

func (wp byGrayPercentage) Len() int           { return len(wp) }
func (wp byGrayPercentage) Swap(i, j int)      { wp[i], wp[j] = wp[j], wp[i] }
func (wp byGrayPercentage) Less(i, j int) bool { return wp[i].GrayPercentage > wp[j].GrayPercentage } // DESC

const MAX_COLOR_VALUE = 255

func makeCharanterImg(basePath string, alias string, char string, width int, height int) (string, error) {
	finalImgName := basePath + alias + "_" + char + ".png"
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	x, y := 2, 10
	col := color.Black

	point := fixed.Point26_6{X: fixed.Int26_6(x * 64), Y: fixed.Int26_6(y * 64)}

	d := &font.Drawer{
		Dst:  img,
		Src:  image.NewUniform(col),
		Face: basicfont.Face7x13,
		Dot:  point,
	}
	d.DrawString(char)

	_, err := imagefilter.EncodeIMG(img, finalImgName)
	if err != nil {
		return "", err
	}
	return finalImgName, nil
}

func characterColorGrayScaleWeight(xp []pixelextract.PixelColor) uint32 {
	var weight uint32
	for _, p := range xp {
		weight += p.ColorGrayScale()
	}
	return weight
}

func characterColorWeight(xp []pixelextract.PixelColor) uint32 {
	var weight uint32
	for _, p := range xp {
		weight += p.ColorWeight()
	}
	return weight
}

func CharacterPixelTypeCountMakeFile(s string) {
	var wg sync.WaitGroup
	characterWeight := []CharacterMetadata{}
	logger := make(chan string)
	fileNameList := make(chan CharacterMetadata)
	basePath := "./files/unpublished/matrix/"
	baseWidth := 10
	baseHeight := 12
	for taskID, c := range s {
		wg.Add(1)
		go func(id int, char string) {
			fileName, err := makeCharanterImg(basePath, "character", char, baseWidth, baseHeight)
			if err != nil {
				logger <- fmt.Sprintf("MakeCharanterImg, task: %d error %v", id, err)
				return
			}
			go func() {
				fileNameList <- CharacterMetadata{char, fileName, 0}
			}()
			wg.Done()
		}(taskID, string(c))
	}

	go func() {
		wg.Wait()
		close(logger)
	}()

	for l := range logger {
		fmt.Println(l)
	}
	maxColorValue := MAX_COLOR_VALUE * baseWidth * baseHeight
	colorFactor := float32(maxColorValue / len(s))
	for i := 0; i < len(s); i++ {
		fileData := <-fileNameList
		img, err := imagefilter.DecodeImg(fileData.Filename)
		if err != nil {
			log.Fatal(fmt.Errorf("decode-png-file %w", err))
		}
		xp := pixelextract.ExtractPixelFromImg(img)
		fileData.GrayPercentage = float32(characterColorWeight(xp)*100) / float32(maxColorValue)
		characterWeight = append(characterWeight, fileData)
	}
	sort.Sort(byGrayPercentage(characterWeight))

	fileBody, err := json.Marshal(CharacterInfo{
		StrBase:       s,
		ColorFactor:   colorFactor,
		MaxColorValue: maxColorValue,
		BaseWidth:     baseWidth,
		BaseHeight:    baseHeight,
		CharacterData: characterWeight})
	if err != nil {
		log.Fatal(fmt.Errorf("json-marshal %w", err))
	}
	outFileName := basePath + "charts_weight.txt"
	outFile, _ := os.Create(outFileName)
	if _, err := outFile.Write(fileBody); err != nil {
		log.Fatal(fmt.Errorf("create-file %w", err))
	}
	outFile.Close()
}

func CharacterScaleTxtFile(fi *imagefilter.FilterImg, id int, chartInfo CharacterInfo) string {
	var currentY int = -1
	outFile, _ := os.Create(fi.GetAlias() + "_character.txt")
	s := time.Now()
	for _, p := range fi.GetXp() {
		pixelGrayScaleWight := p.ColorGrayScale()
		charLent := len(chartInfo.CharacterData)
		colorRageIndex := float32(pixelGrayScaleWight*uint32(charLent-1)) / float32(MAX_COLOR_VALUE)
		charIndex := int(colorRageIndex)

		if charIndex >= charLent {
			err := fmt.Errorf("charIndex: %d, not fount in characterData", charIndex)
			log.Fatal(err)
			fi.AddLog(err.Error())
		}
		defaultSpace := ""
		currentValue := defaultSpace + chartInfo.CharacterData[charIndex].Char + defaultSpace
		if currentY != p.Y {
			currentY = p.Y
			if _, err := outFile.Write([]byte("\n" + currentValue)); err != nil {
				log.Fatal(err)
				fi.AddLog(err.Error())
			}
			continue
		}
		if _, err := outFile.Write([]byte(currentValue)); err != nil {
			log.Fatal(err)
			fi.AddLog(err.Error())
		}
	}
	outFile.Close()
	e := time.Since(s)
	fi.AddLog(fmt.Sprintf("character-scale, task: %d total create txt => %v", id, e))
	path, err := os.Getwd()
	if err != nil {
		fi.AddLog(fmt.Sprintf("character-scale, task: %d error %v", id, err))
	}
	return filepath.Join(path, outFile.Name())
}
