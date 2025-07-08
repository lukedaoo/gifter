package main

import (
	"fmt"
	"golang.org/x/image/draw"
	"image"
	"image/gif"
	"math"
	"os"
	"strings"
	"time"
)

type Context struct {
	Width  int
	Height int
	Gamma  float64
	Color  bool
	Styles string
	Grad   []rune // gradient for unicode character
	Mode   string
}

func ContextFromOptions(opts Options) Context {
	var grad []rune
	gamma := 1.0

	if opts.Mode == "ascii" {
		switch opts.Styles {
		case "shaded":
			grad = []rune(GRAD_SHADED)
		case "bordered":
			grad = []rune(GRAD_BORDERED)
		case "blocky":
			grad = []rune(GRAD_BLOCKY)
		case "ascii2":
			grad = []rune(GRAD_ASCII2)
			gamma = 2.2
		case "normal":
			grad = []rune(GRAD_NORMAL)
		default:
			fmt.Printf("Warning: unknown style %q, defaulting to ascii2\n", opts.Styles)
			grad = []rune(GRAD_ASCII2)
			gamma = 2.2
		}
	}

	return Context{
		Width:  opts.Width,
		Height: opts.Height,
		Color:  opts.Color,
		Styles: opts.Styles,
		Mode:   opts.Mode,
		Grad:   grad,
		Gamma:  gamma,
	}
}

func execute(opts Options) error {
	fmt.Printf("Processing file: %s\n", opts.FilePath)
	fmt.Printf("Input dimensions: width=%d, height=%d\n", opts.Width, opts.Height)
	fmt.Printf("Color output: %v\n", opts.Color)
	fmt.Printf("Mode: %s\n", opts.Mode)

	if opts.Styles != "" && opts.Mode == "ascii" {
		fmt.Printf("Using style: %s\n", opts.Styles)
	} else if opts.Mode == "ascii" {
		fmt.Println("No style provided, using default")
	}

	fmt.Printf("Processing started for file: %s with style: %s\n", opts.FilePath, opts.Styles)

	if opts.Mode == "graphic" {
		err := testGraphics(os.Stdout)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error in Kitty graphics test: %v\n", err)
			fmt.Fprintf(os.Stderr, "Warning: Kitty graphics test failed. Ensure you're using Kitty terminal.\n")
			os.Exit(1)
		}
	}

	var g *gif.GIF
	var err error

	if strings.HasPrefix(opts.FilePath, "http://") ||
		strings.HasPrefix(opts.FilePath, "https://") {
		g, err = downloadGIF(opts.FilePath)
		if err != nil {
			return fmt.Errorf("error downloading GIF from %s: %w", opts.FilePath, err)
		}
		fmt.Printf("Successfully downloaded GIF from %s\n", opts.FilePath)
	} else {
		file, err := os.Open(opts.FilePath)
		if err != nil {
			return fmt.Errorf("error opening file %s: %w", opts.FilePath, err)
		}
		defer file.Close()

		fileInfo, err := file.Stat()
		if err != nil {
			return fmt.Errorf("error getting file info for %s: %w", opts.FilePath, err)
		}

		fileSize := fileInfo.Size()
		fmt.Printf("File size of %s: %d bytes\n", opts.FilePath, fileSize)

		g, err = gif.DecodeAll(file)
		if err != nil {
			return fmt.Errorf("error decoding GIF %s: %w", opts.FilePath, err)
		}
	}

	ctx := ContextFromOptions(opts)
	displayGIF(g, ctx)
	return nil
}

func displayGIF(g *gif.GIF, context Context) {
	fmt.Println("Note: Displaying GIF in graphical mode. Use a terminal like Kitty or iTerm2 for best results.")
	if context.Mode == "graphic" {
		fmt.Fprint(os.Stdout, "\033[H") // Clear test image
		clearTerminal()
	} else {
		fmt.Println("Displaying GIF in ASCII mode.")
	}

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

	const animationID = 1

	for loop := 0; loops == -1 || loop < loops; loop++ {
		for i, frame := range frames {
			if context.Mode == "graphic" {
				err := displayImage(os.Stdout, frame, context.Width, context.Height, animationID)
				if err != nil {
					fmt.Fprintf(os.Stderr, "Error displaying frame %d: %v\n", i, err)
					continue
				}
			} else {
				asciiArt := imageToASCII(frame, context)
				clearTerminal()
				fmt.Print(asciiArt)
			}
			delay := time.Duration(math.Max(50, float64(g.Delay[i])*10)) * time.Millisecond
			time.Sleep(delay)
		}
	}
}

// helper
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
