package experiment

import (
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"log"
	"os"
	"sort"
	"sync"

	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"

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
	Char            string
	Filename        string
	WhitePercentage float32
}

type byWhitePercentage []CharacterMetadata

func (wp byWhitePercentage) Len() int           { return len(wp) }
func (wp byWhitePercentage) Swap(i, j int)      { wp[i], wp[j] = wp[j], wp[i] }
func (wp byWhitePercentage) Less(i, j int) bool { return wp[i].WhitePercentage > wp[j].WhitePercentage } // DESC

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
	maxColorValue := 255 * baseWidth * baseHeight //R+G+B+A * (wight * height)
	colorFactor := float32(maxColorValue / len(s))
	for i := 0; i < len(s); i++ {
		fileData := <-fileNameList
		xp, _, err := pixelextract.ExtractPixelFromFilePNG(fileData.Filename)
		if err != nil {
			log.Fatal(fmt.Errorf("extract-pixel-from-file %w", err))
		}
		fileData.WhitePercentage = float32(characterColorWeight(xp)*100) / float32(maxColorValue)
		characterWeight = append(characterWeight, fileData)
	}
	sort.Sort(byWhitePercentage(characterWeight))

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

	f, err := os.Create(finalImgName)
	if err != nil {
		return "", err
	}
	defer f.Close()
	if err := png.Encode(f, img); err != nil {
		return "", err
	}
	return finalImgName, nil
}

func characterColorGrayScaleWeight(xp []pixelextract.PixelColor) uint32 {
	var weight uint32
	for _, p := range xp {
		weight += p.ColorGrayScaleWeight()
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
