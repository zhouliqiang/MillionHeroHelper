package main

import (
	"log"
	"time"
	"os/exec"
	"os"
	"runtime"
	"image"
	"github.com/oliamb/cutter"
	"image/png"
	"github.com/otiai10/gosseract"
	"fmt"
	"sync"
	"strconv"
	"strings"
	"github.com/fatih/color"
	"github.com/PuerkitoBio/goquery"
)

const (
	MacOS = "darwin"
	Linux = "linux"
)

var wg sync.WaitGroup
var results = make(map[int]string, 4)
var startTime time.Time

func main() {
	startTime = time.Now()
	var cmd *exec.Cmd
	if runtime.GOOS == MacOS {
		cmd = exec.Command("/bin/sh", "-c", "adb shell screencap -p /sdcard/screen.png")
	} else if runtime.GOOS == Linux {
		cmd = exec.Command("/bin/sh", "-c", "adb shell screencap -p | sed 's/\r$//' > screen.png")
	}
	err := cmd.Run()
	handleError(err)

	if runtime.GOOS == MacOS {
		cmd = exec.Command("/bin/sh", "-c", "adb pull /sdcard/screen.png screen.png")
		err := cmd.Run()
		handleError(err)
	}

	originFile, err := os.Open("screen.png")
	handleError(err)

	defer originFile.Close()

	origin, _, err := image.Decode(originFile)
	handleError(err)

	for i := 0; i < 4; i++ {
		wg.Add(1)
		ocrTask(origin, i)
	}

	wg.Wait()

	filterResults()

	quiz := results[0]
	fmt.Println("quiz:", color.GreenString(quiz))

	doc, err := goquery.NewDocument("http://www.baidu.com/s?wd="+quiz)
	handleError(err)

	doc.Find("div.result").Each(func(i int, s *goquery.Selection) {
		content := s.Find("div.c-abstract").Text()
		if strings.Contains(content, results[1]) {
			content = strings.Replace(content, results[1], color.CyanString(results[1]), -1)
		}

		if strings.Contains(content, results[2]){
			content = strings.Replace(content, results[2], color.CyanString(results[2]), -1)
		}

		if strings.Contains(content, results[3]) {
			content = strings.Replace(content, results[3], color.CyanString(results[3]), -1)
		}
		fmt.Println(content)
	})

	fmt.Println("cost time: ", time.Since(startTime))
	os.Exit(0)
}

func ocrTask(origin image.Image, index int) {
	go func() {
		defer wg.Done()
		var croppedImg image.Image
		var err error
		if index == 0 {
			croppedImg, _ = cutter.Crop(origin, cutter.Config{
				Width:  900,
				Height: 235,
				Anchor: image.Point{90, 330},
				Mode:   cutter.TopLeft, // optional, default value
			})
		} else {
			croppedImg, _ = cutter.Crop(origin, cutter.Config{
				Width:  760,
				Height: 160,
				Anchor: image.Point{160, 655 + (index-1) * 193},
				Mode:   cutter.TopLeft, // optional, default value
			})
		}

		handleError(err)

		croppedFile, err := os.Create(strconv.Itoa(index) + ".png")
		handleError(err)

		err = png.Encode(croppedFile, croppedImg)
		handleError(err)
		croppedFile.Close()

		client := gosseract.NewClient()
		defer client.Close()
		client.SetImage(strconv.Itoa(index) + ".png")
		client.SetLanguage("chi_sim")
		text, err := client.Text()
		handleError(err)
		results[index] = text
	}()
}

func filterResults() {
	for k, v := range results {
		switch k {
		case 0:
			rs := []rune(v)
			length := len(rs)
			results[k] = strings.Replace(strings.Replace(string(rs[3:length]), " ", "", -1), "\n", "", -1)
			break
		default:
			results[k] = strings.Replace(v, " ", "", -1)
			break
		}
	}
}

func handleError(err error) {
	if err != nil {
		log.Printf("Error %v .", err)
		os.Exit(1)
	}
}