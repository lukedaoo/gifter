package main

import (
	"fmt"
	"golang.org/x/image/draw"
	"image"
	// "image/color"
	"image/gif"
	"log"
	"math"
	"os"
	"strings"
	"time"
)

const (
	GRAD_NORMAL   = " .,:;ilwW#@$%"
	GRAD_ASCII2   = "$@B%8&WM#*oahkbdpqwmZO0QLCJUYXzcvunxrjft/\\|()1{}[]?-_+~<>i!lI;:,\"^`'. "
	GRAD_SHADED   = " .:;░▒▓█"
	GRAD_BORDERED = " .:-┼┤├┴┬┘└┐┌│"
	GRAD_BLOCKY   = " .:;▋▊▉█"
)

type Context struct {
	Width  int
	Height int
	Gamma  float64
	Color  bool
	Styles string
	Grad   []rune // gradient for unicode character
}

func ContextFromOptions(opts Options) Context {
	var grad []rune
	gamma := 1.0

	switch opts.Styles {
	case "normal":
		grad = []rune(GRAD_NORMAL)
	case "shaded":
		grad = []rune(GRAD_SHADED)
	case "bordered":
		grad = []rune(GRAD_BORDERED)
	case "blocky":
		grad = []rune(GRAD_BLOCKY)
	case "ascii2":
		grad = []rune(GRAD_ASCII2)
		gamma = 2.2
	default:
		fmt.Printf("Warning: unknown style %q, defaulting to ascii2\n", opts.Styles)
		grad = []rune(GRAD_ASCII2)
		gamma = 2.2
	}

	return Context{
		Width:  opts.Width,
		Height: opts.Height,
		Color:  opts.Color,
		Styles: opts.Styles,
		Grad:   grad,
		Gamma:  gamma,
	}
}

func execute(opts Options) {
	fmt.Printf("Processing file: %s\n", opts.FilePath)
	fmt.Printf("Input dimensions: width=%d, height=%d\n", opts.Width, opts.Height)
	fmt.Printf("Color output: %v\n", opts.Color)

	if opts.Styles != "" {
		fmt.Printf("Using style: %s\n", opts.Styles)
	} else {
		fmt.Println("No style provided, using default")
	}

	fmt.Printf("Processing started for file: %s with style: %s\n", opts.FilePath, opts.Styles)

	file, err := os.Open(opts.FilePath)
	if err != nil {
		log.Fatalf("Error opening file %s: %v", opts.FilePath, err)
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		log.Fatalf("Error getting file info for %s: %v", opts.FilePath, err)
	}

	fileSize := fileInfo.Size()
	fmt.Printf("File size of %s: %d bytes\n", opts.FilePath, fileSize)

	g, err := gif.DecodeAll(file)
	if err != nil {
		log.Fatalf("Error decoding GIF %s: %v", opts.FilePath, err)
	}

	ctx := ContextFromOptions(opts)
	displayGIF(g, ctx)
}

func imageToASCII(img image.Image, context Context) string {
	b := img.Bounds()
	imgW, imgH := b.Dx(), b.Dy()

	// Adjust for ASCII aspect ratio
	width := context.Width / 2
	height := context.Height / 2

	xScale := float64(imgW) / float64(width)
	yScale := float64(imgH) / float64(height)

	var sb strings.Builder
	sb.Grow(height * (width + 1)) // +1 for newline

	for y := range height {
		for x := range width {
			srcX := int(float64(x) * xScale)
			srcY := int(float64(y) * yScale)
			if srcX >= imgW || srcY >= imgH {
				sb.WriteByte(' ')
				continue
			}

			// r, g, b, _ := img.At(srcX, srcY).RGBA()
			// r8, g8, b8 := uint8(r>>8), uint8(g>>8), uint8(b>>8)
			//
			// avg := float64(r8+g8+b8) / 3.0

			// c := color.GrayModel.Convert(img.At(srcX, srcY)).(color.Gray)
			// intensity := math.Pow(float64(c.Y)/255.0, gamma)

			r, g, b, _ := img.At(srcX, srcY).RGBA()
			r8, g8, b8 := float64(r>>8), float64(g>>8), float64(b>>8)
			gray := 0.3*r8 + 0.59*g8 + 0.11*b8 // https://www.sciencedirect.com/topics/computer-science/luma-coefficient
			gray = math.Pow(gray/255.0, context.Gamma)

			idx := int(gray*float64(len(context.Grad)-1) + 0.5) // round
			runes := []rune(context.Grad)
			ch := runes[idx]

			if context.Color {
				sb.WriteString(fmt.Sprintf("\033[38;2;%d;%d;%dm%c\033[0m",
					uint8(r>>8), uint8(g>>8), uint8(b>>8), ch))
			} else {
				sb.WriteString(fmt.Sprintf("%c", ch))
			}
		}
		sb.WriteByte('\n')
	}
	return sb.String()
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

func displayGIF(g *gif.GIF, context Context) {

	// enter alternate screen buffer
	fmt.Print("\033[?1049h")
	defer fmt.Print("\033[?1049l")

	loops := g.LoopCount
	switch loops {
	case -1:
		loops = 0
	case 0:
		loops = -1
	default:
	}

	frames := make([]image.Image, len(g.Image))
	for i, frame := range g.Image {
		frames[i] = resizeImage(frame, context.Width, context.Height)
	}

	for loop := 0; loops == -1 || loop < loops; loop++ {
		for i, frame := range frames {
			asciiArt := imageToASCII(frame, context)
			clearTerminal()
			fmt.Print(asciiArt)
			delay := time.Duration(math.Max(50, float64(g.Delay[i])*10)) * time.Millisecond
			time.Sleep(delay)
		}
	}
}
