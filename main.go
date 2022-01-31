package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/victorvbello/img-processing/imagefilter"
	"github.com/victorvbello/img-processing/pixelextract"
)

func main() {

	var wg sync.WaitGroup

	s := time.Now()
	xp, img := pixelextract.ExtractPixelFromFile("./files/original/rhino.jpg")
	fmt.Println("total open ", time.Since(s))

	c := make(chan string)

	t := 2

	filterImg := imagefilter.NewFilterImg("rhino_art", img, xp, 0, c)

	for i := 1; i < t; i++ {
		wg.Add(1)
		go func(id int) {
			//factor := uint32(id) * 40
			//imagefilter.RandomColor(id, xp, img, c, factor)
			//imagefilter.RandomRed(id, xp, img, c, factor)
			//imagefilter.RandomGreen(id, xp, img, c, factor)
			//imagefilter.RandomBlue(id, xp, img, c, factor)
			//filterImg.GreyScale(id)
			filterImg.MakeFromTxtFile(id,
				filterImg.Resize(id, 85).ByteScaleTxtFile(id),
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
