package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/urfave/cli/v2"
	"github.com/victorvbello/img-processing/experiment"
	"github.com/victorvbello/img-processing/imagefilter"
	"github.com/victorvbello/img-processing/pixelextract"
)

const (
	INPUT_DIR  = "./files/original/"
	OUTPUT_DIR = "./files/unpublished/"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	app := &cli.App{
		Commands: []*cli.Command{
			{
				Name:  "all",
				Usage: "Make new images using al filters",
				Action: func(c *cli.Context) error {
					s := time.Now()
					fmt.Println("----ALL INIT----")
					byteImgProcessing("rhino_art", INPUT_DIR+"rhino.jpg")
					characterImgProcessing("rhino_art", INPUT_DIR+"rhino.jpg")
					grayScaleImgProcessing("rhino_art", INPUT_DIR+"rhino.jpg")
					randomColorImgProcessing("rhino_art", INPUT_DIR+"rhino.jpg")
					randomColorRedImgProcessing("rhino_art", INPUT_DIR+"rhino.jpg")
					randomColorGreenImgProcessing("rhino_art", INPUT_DIR+"rhino.jpg")
					randomColorBlueImgProcessing("rhino_art", INPUT_DIR+"rhino.jpg")
					fmt.Println("----ALL END----", time.Since(s))
					return nil
				},
			},
			{
				Name:  "character-pixel-weight",
				Usage: "Character pixel weight to file",
				Action: func(c *cli.Context) error {
					experiment.CharacterPixelTypeCountMakeFile("ABCDEFGHIJKLMNOPQRSTUVWXYZ")
					return nil
				},
			},
			{
				Name:  "byte",
				Usage: "Make new img using a byte filter",
				Action: func(c *cli.Context) error {
					byteImgProcessing("rhino_art", INPUT_DIR+"rhino.jpg")
					return nil
				},
			},
			{
				Name:  "character-pixel-color-replace",
				Usage: "Character pixel color replace",
				Action: func(c *cli.Context) error {
					characterImgProcessing("rhino_art", INPUT_DIR+"rhino.jpg")
					return nil
				},
			},
			{
				Name:  "grayscale",
				Usage: "Make new img using a greyScale filter",
				Action: func(c *cli.Context) error {
					grayScaleImgProcessing("rhino_art", INPUT_DIR+"rhino.jpg")
					return nil
				},
			},
			{
				Name:  "random-color",
				Usage: "Make new img using a randomColor filter",
				Action: func(c *cli.Context) error {
					randomColorImgProcessing("rhino_art", INPUT_DIR+"rhino.jpg")
					return nil
				},
			},
			{
				Name:  "random-color-red",
				Usage: "Make new img using a random color red filter",
				Action: func(c *cli.Context) error {
					randomColorRedImgProcessing("rhino_art", INPUT_DIR+"rhino.jpg")
					return nil
				},
			},
			{
				Name:  "random-color-green",
				Usage: "Make new img using a random color gree filter",
				Action: func(c *cli.Context) error {
					randomColorGreenImgProcessing("rhino_art", INPUT_DIR+"rhino.jpg")
					return nil
				},
			},
			{
				Name:  "random-color-blue",
				Usage: "Make new img using a random color blue filter",
				Action: func(c *cli.Context) error {
					randomColorBlueImgProcessing("rhino_art", INPUT_DIR+"rhino.jpg")
					return nil
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func byteImgProcessing(alias string, imgFile string) {
	var wg sync.WaitGroup
	fileProcessFlag := "byte"
	fmt.Println("process", fileProcessFlag)
	s := time.Now()
	img, err := imagefilter.DecodeImg(imgFile)
	if err != nil {
		log.Fatal(fmt.Errorf("decode-file %w", err))
	}
	xp := pixelextract.ExtractPixelFromImg(img)
	fmt.Println("total open ", time.Since(s))

	c := make(chan string)

	t := 2

	filterImg := imagefilter.NewFilterImg(alias, img, xp, 0, c)

	for i := 1; i < t; i++ {
		wg.Add(1)
		go func(id int) {
			ss := time.Now()
			txtFileName := filterImg.Resize(id, 85).ByteScaleTxtFile(id)
			img, err := filterImg.MakeFromTxtFile(txtFileName)
			if err != nil {
				filterImg.AddLog("Error txt to img " + err.Error())
			}
			_, err = imagefilter.EncodeIMG(img, OUTPUT_DIR+alias+"/"+alias+"_"+fileProcessFlag+filepath.Ext(imgFile))
			if err != nil {
				filterImg.AddLog("Error encode img " + err.Error())
			}
			filterImg.AddLog("total process " + time.Since(ss).String())
			wg.Done()
		}(i)
	}

	go func() {
		wg.Wait()
		close(c)
	}()

	for l := range c {
		fmt.Printf("\t%s\n", l)
	}
}

func grayScaleImgProcessing(alias string, imgFile string) {
	var wg sync.WaitGroup
	fileProcessFlag := "grayscale"
	fmt.Println("process", fileProcessFlag)
	s := time.Now()
	img, err := imagefilter.DecodeImg(imgFile)
	if err != nil {
		log.Fatal(fmt.Errorf("decode-file %w", err))
	}
	xp := pixelextract.ExtractPixelFromImg(img)
	fmt.Println("total open ", time.Since(s))

	c := make(chan string)

	t := 2

	filterImg := imagefilter.NewFilterImg(alias, img, xp, 0, c)

	for i := 1; i < t; i++ {
		wg.Add(1)
		go func(id int) {
			ss := time.Now()
			img := filterImg.GreyScale(id)
			_, err = imagefilter.EncodeIMG(img, OUTPUT_DIR+alias+"/"+alias+"_"+fileProcessFlag+filepath.Ext(imgFile))
			if err != nil {
				filterImg.AddLog("Error encode img " + err.Error())
			}
			filterImg.AddLog("total process " + time.Since(ss).String())
			wg.Done()
		}(i)
	}

	go func() {
		wg.Wait()
		close(c)
	}()

	for l := range c {
		fmt.Printf("\t%s\n", l)
	}
}

func characterImgProcessing(alias string, imgFile string) {
	var wg sync.WaitGroup
	fileProcessFlag := "character"
	fmt.Println("process", fileProcessFlag)
	sw := time.Now()
	jsonCharWeight, err := os.Open("./files/unpublished/matrix/charts_weight.txt")
	if err != nil {
		log.Fatal(fmt.Errorf("open-weight-file %w", err))
	}
	byteValue, err := ioutil.ReadAll(jsonCharWeight)
	if err != nil {
		log.Fatal(fmt.Errorf("read-weight-file %w", err))
	}
	var charInfo experiment.CharacterInfo
	err = json.Unmarshal(byteValue, &charInfo)
	if err != nil {
		log.Fatal(fmt.Errorf("json-unmarshal-weight %w", err))
	}
	jsonCharWeight.Close()
	fmt.Println("total open charts_weight file", time.Since(sw))

	s := time.Now()

	img, err := imagefilter.DecodeImg(imgFile)
	if err != nil {
		log.Fatal(fmt.Errorf("decode-file %w", err))
	}
	xp := pixelextract.ExtractPixelFromImg(img)
	fmt.Println("total open ", time.Since(s))

	c := make(chan string)

	t := 2

	filterImg := imagefilter.NewFilterImg(alias, img, xp, 0, c)
	filterImg.Transparency(2)
	for i := 1; i < t; i++ {
		wg.Add(1)
		go func(id int) {
			ss := time.Now()
			txtFileName := experiment.CharacterScaleTxtFile(filterImg.Resize(id, 85), id, charInfo)
			img, err := filterImg.MakeFromTxtFile(txtFileName)
			if err != nil {
				filterImg.AddLog("Error txt to img " + err.Error())
			}

			_, err = imagefilter.EncodeIMG(img, OUTPUT_DIR+alias+"/"+alias+"_"+fileProcessFlag+filepath.Ext(imgFile))
			if err != nil {
				filterImg.AddLog("Error encode img " + err.Error())
			}
			filterImg.AddLog("total process " + time.Since(ss).String())
			wg.Done()
		}(i)
	}

	go func() {
		wg.Wait()
		close(c)
	}()

	for l := range c {
		fmt.Printf("\t%s\n", l)
	}
}

func randomColorImgProcessing(alias string, imgFile string) {
	var wg sync.WaitGroup
	fileProcessFlag := "random_color"
	fmt.Println("process", fileProcessFlag)
	s := time.Now()
	img, err := imagefilter.DecodeImg(imgFile)
	if err != nil {
		log.Fatal(fmt.Errorf("decode-file %w", err))
	}
	xp := pixelextract.ExtractPixelFromImg(img)
	fmt.Println("total open ", time.Since(s))

	c := make(chan string)

	t := 2

	filterImg := imagefilter.NewFilterImg(alias, img, xp, 0, c)

	for i := 1; i < t; i++ {
		wg.Add(1)
		go func(id int) {
			filterImg.Factor = uint32(rand.Intn(255))
			filterImg.AddLog(fmt.Sprintf("factor %d", filterImg.Factor))
			ss := time.Now()
			img := filterImg.RandomColor(id)
			_, err = imagefilter.EncodeIMG(img, OUTPUT_DIR+alias+"/"+alias+"_"+fileProcessFlag+filepath.Ext(imgFile))
			if err != nil {
				filterImg.AddLog("Error encode img " + err.Error())
			}
			filterImg.AddLog("total process " + time.Since(ss).String())
			wg.Done()
		}(i)
	}

	go func() {
		wg.Wait()
		close(c)
	}()

	for l := range c {
		fmt.Printf("\t%s\n", l)
	}
}

func randomColorRedImgProcessing(alias string, imgFile string) {
	var wg sync.WaitGroup
	fileProcessFlag := "random_color_red"
	fmt.Println("process", fileProcessFlag)
	s := time.Now()
	img, err := imagefilter.DecodeImg(imgFile)
	if err != nil {
		log.Fatal(fmt.Errorf("decode-file %w", err))
	}
	xp := pixelextract.ExtractPixelFromImg(img)
	fmt.Println("total open ", time.Since(s))

	c := make(chan string)

	t := 2

	filterImg := imagefilter.NewFilterImg(alias, img, xp, 0, c)

	for i := 1; i < t; i++ {
		wg.Add(1)
		go func(id int) {
			filterImg.Factor = uint32(rand.Intn(255))
			filterImg.AddLog(fmt.Sprintf("factor %d", filterImg.Factor))
			ss := time.Now()
			img := filterImg.RandomRed(id)
			_, err = imagefilter.EncodeIMG(img, OUTPUT_DIR+alias+"/"+alias+"_"+fileProcessFlag+filepath.Ext(imgFile))
			if err != nil {
				filterImg.AddLog("Error encode img " + err.Error())
			}
			filterImg.AddLog("total process " + time.Since(ss).String())
			wg.Done()
		}(i)
	}

	go func() {
		wg.Wait()
		close(c)
	}()

	for l := range c {
		fmt.Printf("\t%s\n", l)
	}
}

func randomColorBlueImgProcessing(alias string, imgFile string) {
	var wg sync.WaitGroup
	fileProcessFlag := "random_color_blue"
	fmt.Println("process", fileProcessFlag)
	s := time.Now()
	img, err := imagefilter.DecodeImg(imgFile)
	if err != nil {
		log.Fatal(fmt.Errorf("decode-file %w", err))
	}
	xp := pixelextract.ExtractPixelFromImg(img)
	fmt.Println("total open ", time.Since(s))

	c := make(chan string)

	t := 2

	filterImg := imagefilter.NewFilterImg(alias, img, xp, 0, c)

	for i := 1; i < t; i++ {
		wg.Add(1)
		go func(id int) {
			filterImg.Factor = uint32(rand.Intn(255))
			filterImg.AddLog(fmt.Sprintf("factor %d", filterImg.Factor))
			ss := time.Now()
			img := filterImg.RandomBlue(id)
			_, err = imagefilter.EncodeIMG(img, OUTPUT_DIR+alias+"/"+alias+"_"+fileProcessFlag+filepath.Ext(imgFile))
			if err != nil {
				filterImg.AddLog("Error encode img " + err.Error())
			}
			filterImg.AddLog("total process " + time.Since(ss).String())
			wg.Done()
		}(i)
	}

	go func() {
		wg.Wait()
		close(c)
	}()

	for l := range c {
		fmt.Printf("\t%s\n", l)
	}
}
func randomColorGreenImgProcessing(alias string, imgFile string) {
	var wg sync.WaitGroup
	fileProcessFlag := "random_color_green"
	fmt.Println("process", fileProcessFlag)
	s := time.Now()
	img, err := imagefilter.DecodeImg(imgFile)
	if err != nil {
		log.Fatal(fmt.Errorf("decode-file %w", err))
	}
	xp := pixelextract.ExtractPixelFromImg(img)
	fmt.Println("total open ", time.Since(s))

	c := make(chan string)

	t := 2

	filterImg := imagefilter.NewFilterImg(alias, img, xp, 0, c)

	for i := 1; i < t; i++ {
		wg.Add(1)
		go func(id int) {
			filterImg.Factor = uint32(rand.Intn(255))
			filterImg.AddLog(fmt.Sprintf("factor %d", filterImg.Factor))
			ss := time.Now()
			img := filterImg.RandomGreen(id)
			_, err = imagefilter.EncodeIMG(img, OUTPUT_DIR+alias+"/"+alias+"_"+fileProcessFlag+filepath.Ext(imgFile))
			if err != nil {
				filterImg.AddLog("Error encode img " + err.Error())
			}
			filterImg.AddLog("total process " + time.Since(ss).String())
			wg.Done()
		}(i)
	}

	go func() {
		wg.Wait()
		close(c)
	}()

	for l := range c {
		fmt.Printf("\t%s\n", l)
	}
}
