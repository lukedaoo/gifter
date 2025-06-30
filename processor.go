package main

import (
	"fmt"
	"golang.org/x/image/draw"
	"image"
	"image/color"
	"image/gif"
	"log"
	"math"
	"os"
	"time"
)

const ASCII_CHARS = "$@%&#*+/\\|()1{}[]?-_+~<>!;:,\"^`'.                                                 "

func execute(filePath string, styles string) {
	fmt.Printf("Processing file: %s\n", filePath)

	if styles != "" {
		fmt.Printf("Using styles: %s\n", styles)
	} else {
		fmt.Println("No styles provided, using default")
		styles = "normal"
	}

	fmt.Printf("Processing complete for file: %s with styles: %s\n", filePath, styles)

	file, err := os.Open(filePath)
	if err != nil {
		log.Fatalf("Error opening file: %v", err)
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		log.Fatalf("Error getting file info from opened file: %v", err)
	}

	fileSize := fileInfo.Size()

	fmt.Printf("The size of %s is %d bytes.\n", filePath, fileSize)

	g, err := gif.DecodeAll(file)
	if err != nil {
		panic(err)
	}
	displayGIF(g)
}

func imageToASCII(img image.Image, width, height int) string {
	bounds := img.Bounds()
	imgWidth, imgHeight := bounds.Max.X, bounds.Max.Y

	xScale := float64(imgWidth) / float64(width)
	yScale := float64(imgHeight) / float64(height)

	output := ""

	const gamma = 2.2

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			srcX := int(float64(x) * xScale)
			srcY := int(float64(y) * yScale)

			if srcX >= imgWidth || srcY >= imgHeight {
				output += " "
				continue
			}

			c := color.GrayModel.Convert(img.At(srcX, srcY)).(color.Gray)
			intensity := math.Pow(float64(c.Y)/255.0, gamma)
			charIndex := int(intensity * float64(len(ASCII_CHARS)-1))
			output += string(ASCII_CHARS[charIndex])
		}
		output += "\n"
	}
	return output
}

func clearTerminal() {
	fmt.Print("\033[H\033[2J") // move cursor to top-left and clear screen
	fmt.Print("\033[3J")       // clear scrollback buffer (Kitty-specific)
	fmt.Print("\033c")         // reset terminal
	os.Stdout.Sync()           // force flush output buffer
}

func resizeImage(img image.Image, width, height int) image.Image {
	resized := image.NewRGBA(image.Rect(0, 0, width, height))
	draw.NearestNeighbor.Scale(resized, resized.Bounds(), img, img.Bounds(), draw.Src, nil)
	return resized
}

func displayGIF(g *gif.GIF) {
	// enter alternate screen buffer
	fmt.Print("\033[?1049h")
	defer fmt.Print("\033[?1049l")

	loops := g.LoopCount
	if loops == -1 {
		loops = 1
	} else if loops == 0 {
		loops = -1
	}

	frames := make([]image.Image, len(g.Image))
	for i, frame := range g.Image {
		frames[i] = resizeImage(frame, 160, 80)
	}

	for loop := 0; loops == -1 || loop < loops; loop++ {
		for i, frame := range frames {
			asciiArt := imageToASCII(frame, 80, 40)
			clearTerminal()
			fmt.Print(asciiArt)
			delay := time.Duration(g.Delay[i]) * 10 * time.Millisecond
			if delay < 50*time.Millisecond {
				delay = 50 * time.Millisecond
			}
			time.Sleep(delay)
		}
	}
}
