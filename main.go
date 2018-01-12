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

func main() {
	fmt.Println("\nstart time: ", time.Now())
	var cmd *exec.Cmd
	if runtime.GOOS == MacOS {
		cmd = exec.Command("/bin/sh", "-c", "adb shell screencap -p /sdcard/screen.png")
	} else if runtime.GOOS == Linux {
		cmd = exec.Command("/bin/sh", "-c", "adb shell screencap -p | sed 's/\r$//' > screen.png")
	}
	if err := cmd.Run(); err != nil {
		log.Printf("Error %v executing command!", err)
		os.Exit(1)
	}

	if runtime.GOOS == MacOS {
		cmd = exec.Command("/bin/sh", "-c", "adb pull /sdcard/screen.png screen.png")
		if err := cmd.Run(); err != nil {
			log.Printf("Error %v executing command!", err)
			os.Exit(1)
		}
	}

	originFile, err := os.Open("screen.png")
	if err != nil {
		log.Printf("Error %v open image file.", err)
		os.Exit(1)
	}
	defer originFile.Close()

	origin, _, err := image.Decode(originFile)
	if err != nil {
		log.Printf("Error %v decode image file.", err)
		os.Exit(1)
	}

	for i :=0; i < 4; i++ {
		wg.Add(1)
		ocrTask(origin, i)
	}

	wg.Wait()

	quiz := results[0]
	quiz = pickQuiz(quiz)
	fmt.Println("quiz:", color.GreenString(quiz))

	doc, err := goquery.NewDocument("http://www.baidu.com/s?wd="+quiz)
	if err != nil {
		log.Printf("Error %v while HTTP reqeust.", err)
		panic(err)
	}

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

	fmt.Println("end time: ", time.Now())
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

		if err != nil {
			log.Printf("Error %v crop %d image file", err, index)
			os.Exit(1)
		}

		croppedFile, err := os.Create(strconv.Itoa(index) + ".png")
		if err != nil {
			log.Printf("Error %v create cropped %d image file", err, index)
			os.Exit(1)
		}

		err = png.Encode(croppedFile, croppedImg)
		if err != nil {
			log.Printf("Error %v encode cropped %d image file", err, index)
			os.Exit(1)
		}
		croppedFile.Close()

		client := gosseract.NewClient()
		defer client.Close()
		client.SetImage(strconv.Itoa(index) + ".png")
		client.SetLanguage("chi_sim")
		text, _ := client.Text()
		results[index] = text
	}()
}

func pickQuiz(str string) string {
	rs := []rune(str)
	length := len(rs)
	return strings.Replace(strings.Replace(string(rs[3:length]), " ", "", -1), "\n", "", -1)
}