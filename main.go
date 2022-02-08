package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sync"
	"time"

	"github.com/urfave/cli/v2"
	"github.com/victorvbello/img-processing/experiment"
	"github.com/victorvbello/img-processing/imagefilter"
	"github.com/victorvbello/img-processing/pixelextract"
)

//factor := uint32(id) * 40
//imagefilter.RandomColor(id, xp, img, c, factor)
//imagefilter.RandomRed(id, xp, img, c, factor)
//imagefilter.RandomGreen(id, xp, img, c, factor)
//imagefilter.RandomBlue(id, xp, img, c, factor)
//filterImg.GreyScale(id)

func main() {
	app := &cli.App{
		Commands: []*cli.Command{
			{
				Name:  "byte-img-processing",
				Usage: "Make new img using a byte filter",
				Action: func(c *cli.Context) error {
					byteImgProcessing("o_art", "./files/original/o.jpg")
					return nil
				},
			},
			{
				Name:  "grayscale-img-processing",
				Usage: "Make new img using a greyScale filter",
				Action: func(c *cli.Context) error {
					grayScaleImgProcessing("rhino_art", "./files/original/rhino.jpg")
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
				Name:  "character-pixel-color-replace",
				Usage: "Character pixel color replace",
				Action: func(c *cli.Context) error {
					matrixImgProcessing("rhino_art", "./rhino_art_grey.jpg")
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
	s := time.Now()
	xp, img, err := pixelextract.ExtractPixelFromFileJPEG(imgFile)
	if err != nil {
		log.Fatal(fmt.Errorf("extract-pixel-from-file %w", err))
	}
	fmt.Println("total open ", time.Since(s))

	c := make(chan string)

	t := 2

	filterImg := imagefilter.NewFilterImg(alias, img, xp, 0, c)

	for i := 1; i < t; i++ {
		wg.Add(1)
		go func(id int) {
			filterImg.MakeFromTxtFile(
				id,
				"byte",
				filterImg.Resize(id, 85).ByteScaleTxtFile(id),
				11,
				"Times-Roman",
			)
			wg.Done()
		}(i)
	}

	go func() {
		wg.Wait()
		close(c)
	}()

	for l := range c {
		fmt.Println(l)
	}
}

func grayScaleImgProcessing(alias string, imgFile string) {
	var wg sync.WaitGroup
	s := time.Now()
	xp, img, err := pixelextract.ExtractPixelFromFileJPEG(imgFile)
	if err != nil {
		log.Fatal(fmt.Errorf("extract-pixel-from-file %w", err))
	}
	fmt.Println("total open ", time.Since(s))

	c := make(chan string)

	t := 2

	filterImg := imagefilter.NewFilterImg(alias, img, xp, 0, c)

	for i := 1; i < t; i++ {
		wg.Add(1)
		go func(id int) {
			filterImg.GreyScale(id)
			wg.Done()
		}(i)
	}

	go func() {
		wg.Wait()
		close(c)
	}()

	for l := range c {
		fmt.Println(l)
	}
}

func matrixImgProcessing(alias string, imgFile string) {
	var wg sync.WaitGroup

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
	xp, img, err := pixelextract.ExtractPixelFromFileJPEG(imgFile)
	if err != nil {
		log.Fatal(fmt.Errorf("extract-pixel-from-file %w", err))
	}
	fmt.Println("total open ", time.Since(s))

	c := make(chan string)

	t := 2

	filterImg := imagefilter.NewFilterImg(alias, img, xp, 0, c)

	for i := 1; i < t; i++ {
		wg.Add(1)
		go func(id int) {
			filterImg.MakeFromTxtFile(
				id,
				"character",
				filterImg.Resize(id, 85).CharacterScaleTxtFile(id, charInfo),
				5,
				"Times-Roman",
			)
			wg.Done()
		}(i)
	}

	go func() {
		wg.Wait()
		close(c)
	}()

	for l := range c {
		fmt.Println(l)
	}
}
