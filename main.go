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
	"github.com/victorvbello/img-processing/imagetransforms"
	"github.com/victorvbello/img-processing/pixelextract"
)

const (
	INPUT_DIR  = "./files/original/"
	OUTPUT_DIR = "./files/unpublished/"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	app := &cli.App{
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "alias",
				Aliases:  []string{"a"},
				Usage:    "Alias of file",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "file",
				Aliases:  []string{"f"},
				Usage:    "Original file path",
				Required: true,
			},
		},
		Commands: []*cli.Command{
			{
				Name:  "all",
				Usage: "Make new images using al filters",
				Action: func(c *cli.Context) error {
					s := time.Now()
					log.Println("----ALL INIT----")
					alias := c.String("alias")
					inputFile := c.String("file")
					byteImgProcessing(alias, inputFile)
					characterImgProcessing(alias, inputFile)
					grayScaleImgProcessing(alias, inputFile)
					randomColorImgProcessing(alias, inputFile)
					randomColorRedImgProcessing(alias, inputFile)
					randomColorGreenImgProcessing(alias, inputFile)
					randomColorBlueImgProcessing(alias, inputFile)
					log.Println("----ALL END----", time.Since(s))
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
				Name:  "infinite",
				Usage: "Mane new img infinite",
				Action: func(c *cli.Context) error {
					alias := c.String("alias")
					inputFile := c.String("file")
					infiniteImgProcessing(alias, inputFile)
					return nil
				},
			},
			{
				Name:  "infinite-spiral",
				Usage: "Mane new img infinite spiral",
				Action: func(c *cli.Context) error {
					alias := c.String("alias")
					inputFile := c.String("file")
					infiniteSpiralImgProcessing(alias, inputFile)
					return nil
				},
			},
			{
				Name:  "byte",
				Usage: "Make new img using a byte filter",
				Action: func(c *cli.Context) error {
					alias := c.String("alias")
					inputFile := c.String("file")
					byteImgProcessing(alias, inputFile)
					return nil
				},
			},
			{
				Name:  "character-pixel-color-replace",
				Usage: "Character pixel color replace",
				Action: func(c *cli.Context) error {
					alias := c.String("alias")
					inputFile := c.String("file")
					characterImgProcessing(alias, inputFile)
					return nil
				},
			},
			{
				Name:  "grayscale",
				Usage: "Make new img using a greyScale filter",
				Action: func(c *cli.Context) error {
					alias := c.String("alias")
					inputFile := c.String("file")
					grayScaleImgProcessing(alias, inputFile)
					return nil
				},
			},
			{
				Name:  "random-color",
				Usage: "Make new img using a randomColor filter",
				Action: func(c *cli.Context) error {
					alias := c.String("alias")
					inputFile := c.String("file")
					randomColorImgProcessing(alias, inputFile)
					return nil
				},
			},
			{
				Name:  "random-color-red",
				Usage: "Make new img using a random color red filter",
				Action: func(c *cli.Context) error {
					alias := c.String("alias")
					inputFile := c.String("file")
					randomColorRedImgProcessing(alias, inputFile)
					return nil
				},
			},
			{
				Name:  "random-color-green",
				Usage: "Make new img using a random color gree filter",
				Action: func(c *cli.Context) error {
					alias := c.String("alias")
					inputFile := c.String("file")
					randomColorGreenImgProcessing(alias, inputFile)
					return nil
				},
			},
			{
				Name:  "random-color-blue",
				Usage: "Make new img using a random color blue filter",
				Action: func(c *cli.Context) error {
					alias := c.String("alias")
					inputFile := c.String("file")
					randomColorBlueImgProcessing(alias, inputFile)
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
	log.Println("process", fileProcessFlag)
	s := time.Now()
	img, err := imagefilter.DecodeImg(imgFile)
	if err != nil {
		log.Fatal(fmt.Errorf("decode-file %w", err))
	}
	xp := pixelextract.ExtractPixelFromImg(img)
	log.Println("total open ", time.Since(s))

	c := make(chan string)

	t := 2

	filterImg := imagefilter.NewFilterImg(alias, img, xp, 0, c)
	for i := 1; i < t; i++ {
		wg.Add(1)
		go func(id int) {
			sr := time.Now()
			newImg := imagetransforms.Resize(filterImg, 85)
			filterImg.AddLog(fmt.Sprintf("resize, task: %d total resize => %v", id, time.Since(sr).String()))
			newXp := pixelextract.ExtractPixelFromImg(newImg)
			ss := time.Now()
			filterImg.SetXp(newXp)
			txtFileName := filterImg.ByteScaleTxtFile(id)
			img, err := filterImg.MakeFromTxtFile(txtFileName)
			if err != nil {
				filterImg.AddLog("Error txt to img " + err.Error())
			}
			_, err = imagefilter.EncodeIMG(img, OUTPUT_DIR+alias+"/"+alias+"_"+fileProcessFlag+filepath.Ext(imgFile))
			if err != nil {
				filterImg.AddLog("Error encode img " + err.Error())
				return
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
	log.Println("process", fileProcessFlag)
	s := time.Now()
	img, err := imagefilter.DecodeImg(imgFile)
	if err != nil {
		log.Fatal(fmt.Errorf("decode-file %w", err))
	}
	xp := pixelextract.ExtractPixelFromImg(img)
	log.Println("total open ", time.Since(s))

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
				return
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
	log.Println("process", fileProcessFlag)
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
	log.Println("total open charts_weight file", time.Since(sw))

	s := time.Now()

	img, err := imagefilter.DecodeImg(imgFile)
	if err != nil {
		log.Fatal(fmt.Errorf("decode-file %w", err))
	}
	xp := pixelextract.ExtractPixelFromImg(img)
	log.Println("total open ", time.Since(s))

	c := make(chan string)

	t := 2

	filterImg := imagefilter.NewFilterImg(alias, img, xp, 0, c)
	newImg := imagetransforms.Transparency(filterImg, 2)
	for i := 1; i < t; i++ {
		wg.Add(1)
		go func(id int) {
			sr := time.Now()
			newImg := imagetransforms.Resize(newImg, 85)
			newXp := pixelextract.ExtractPixelFromImg(newImg)
			filterImg.AddLog(fmt.Sprintf("resize, task: %d total resize => %v", id, time.Since(sr).String()))
			ss := time.Now()
			filterImg.SetXp(newXp)
			txtFileName := experiment.CharacterScaleTxtFile(filterImg, id, charInfo)
			img, err := filterImg.MakeFromTxtFile(txtFileName)
			if err != nil {
				filterImg.AddLog("Error txt to img " + err.Error())
			}

			_, err = imagefilter.EncodeIMG(img, OUTPUT_DIR+alias+"/"+alias+"_"+fileProcessFlag+filepath.Ext(imgFile))
			if err != nil {
				filterImg.AddLog("Error encode img " + err.Error())
				return
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
	log.Println("process", fileProcessFlag)
	s := time.Now()
	img, err := imagefilter.DecodeImg(imgFile)
	if err != nil {
		log.Fatal(fmt.Errorf("decode-file %w", err))
	}
	xp := pixelextract.ExtractPixelFromImg(img)
	log.Println("total open ", time.Since(s))

	c := make(chan string)

	t := 2

	filterImg := imagefilter.NewFilterImg(alias, img, xp, 0, c)

	for i := 1; i < t; i++ {
		wg.Add(1)
		go func(id int) {
			filterImg.Factor = uint8(rand.Intn(255))
			filterImg.AddLog(fmt.Sprintf("factor %d", filterImg.Factor))
			ss := time.Now()
			img := filterImg.RandomColor(id)
			_, err = imagefilter.EncodeIMG(img, OUTPUT_DIR+alias+"/"+alias+"_"+fileProcessFlag+filepath.Ext(imgFile))
			if err != nil {
				filterImg.AddLog("Error encode img " + err.Error())
				return
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
	log.Println("process", fileProcessFlag)
	s := time.Now()
	img, err := imagefilter.DecodeImg(imgFile)
	if err != nil {
		log.Fatal(fmt.Errorf("decode-file %w", err))
	}
	xp := pixelextract.ExtractPixelFromImg(img)
	log.Println("total open ", time.Since(s))

	c := make(chan string)

	t := 2

	filterImg := imagefilter.NewFilterImg(alias, img, xp, 0, c)

	for i := 1; i < t; i++ {
		wg.Add(1)
		go func(id int) {
			filterImg.Factor = uint8(rand.Intn(255))
			filterImg.AddLog(fmt.Sprintf("factor %d", filterImg.Factor))
			ss := time.Now()
			img := filterImg.RandomRed(id)
			_, err = imagefilter.EncodeIMG(img, OUTPUT_DIR+alias+"/"+alias+"_"+fileProcessFlag+filepath.Ext(imgFile))
			if err != nil {
				filterImg.AddLog("Error encode img " + err.Error())
				return
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
	log.Println("process", fileProcessFlag)
	s := time.Now()
	img, err := imagefilter.DecodeImg(imgFile)
	if err != nil {
		log.Fatal(fmt.Errorf("decode-file %w", err))
	}
	xp := pixelextract.ExtractPixelFromImg(img)
	log.Println("total open ", time.Since(s))

	c := make(chan string)

	t := 2

	filterImg := imagefilter.NewFilterImg(alias, img, xp, 0, c)

	for i := 1; i < t; i++ {
		wg.Add(1)
		go func(id int) {
			filterImg.Factor = uint8(rand.Intn(255))
			filterImg.AddLog(fmt.Sprintf("factor %d", filterImg.Factor))
			ss := time.Now()
			img := filterImg.RandomBlue(id)
			_, err = imagefilter.EncodeIMG(img, OUTPUT_DIR+alias+"/"+alias+"_"+fileProcessFlag+filepath.Ext(imgFile))
			if err != nil {
				filterImg.AddLog("Error encode img " + err.Error())
				return
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
	log.Println("process", fileProcessFlag)
	s := time.Now()
	img, err := imagefilter.DecodeImg(imgFile)
	if err != nil {
		log.Fatal(fmt.Errorf("decode-file %w", err))
	}
	xp := pixelextract.ExtractPixelFromImg(img)
	log.Println("total open ", time.Since(s))

	c := make(chan string)

	t := 2

	filterImg := imagefilter.NewFilterImg(alias, img, xp, 0, c)

	for i := 1; i < t; i++ {
		wg.Add(1)
		go func(id int) {
			filterImg.Factor = uint8(rand.Intn(255))
			filterImg.AddLog(fmt.Sprintf("factor %d", filterImg.Factor))
			ss := time.Now()
			img := filterImg.RandomGreen(id)
			_, err = imagefilter.EncodeIMG(img, OUTPUT_DIR+alias+"/"+alias+"_"+fileProcessFlag+filepath.Ext(imgFile))
			if err != nil {
				filterImg.AddLog("Error encode img " + err.Error())
				return
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

func infiniteImgProcessing(alias string, imgFile string) {
	var wg sync.WaitGroup
	fileProcessFlag := "infinite"
	log.Println("process", fileProcessFlag)
	s := time.Now()
	img, err := imagefilter.DecodeImg(imgFile)
	if err != nil {
		log.Fatal(fmt.Errorf("decode-file %w", err))
	}
	log.Println("total open ", time.Since(s))

	c := make(chan string)

	t := 2

	for i := 1; i < t; i++ {
		wg.Add(1)
		go func(id int) {
			ss := time.Now()
			newImg := experiment.ImgInfinite(img, 5)
			c <- fmt.Sprintf("percentage %d", 5)
			_, err = imagefilter.EncodeIMG(newImg, OUTPUT_DIR+alias+"/"+alias+"_"+fileProcessFlag+filepath.Ext(imgFile))
			if err != nil {
				c <- "Error encode img " + err.Error()
				return
			}
			c <- "total process " + time.Since(ss).String()
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

func infiniteSpiralImgProcessing(alias string, imgFile string) {
	var wg sync.WaitGroup
	fileProcessFlag := "infinite_spiral"
	log.Println("process", fileProcessFlag)
	s := time.Now()
	img, err := imagefilter.DecodeImg(imgFile)
	if err != nil {
		log.Fatal(fmt.Errorf("decode-file %w", err))
	}
	log.Println("total open ", time.Since(s))

	c := make(chan string)

	t := 2

	for i := 1; i < t; i++ {
		wg.Add(1)
		go func(id int) {
			ss := time.Now()
			newImg := experiment.ImgInfiniteSpiral(img, 5)
			_, err = imagefilter.EncodeIMG(newImg, fmt.Sprintf("%s%s/%s_%s%s", OUTPUT_DIR, alias, alias, fileProcessFlag, ".png"))
			if err != nil {
				c <- "Error encode img " + err.Error()
				return
			}
			c <- "total process " + time.Since(ss).String()
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
