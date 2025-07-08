package main

import (
	"fmt"
	"image"
	"math"
	"strings"
)

const (
	GRAD_NORMAL   = " .,:;ilwW#@$%"
	GRAD_ASCII2   = "$@B%8&WM#*oahkbdpqwmZO0QLCJUYXzcvunxrjft/\\|()1{}[]?-_+~<>i!lI;:,\"^`'. "
	GRAD_SHADED   = " .:;░▒▓█"
	GRAD_BORDERED = " .:-┼┤├┴┬┘└┐┌│"
	GRAD_BLOCKY   = " .:;▋▊▉█"
)

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

			r, g, b, _ := img.At(srcX, srcY).RGBA()
			r8, g8, b8 := float64(r>>8), float64(g>>8), float64(b>>8)
			gray := 0.3*r8 + 0.59*g8 + 0.11*b8 // https://www.sciencedirect.com/topics/computer-science/luma-coefficient
			gray = math.Pow(gray/255.0, context.Gamma)

			idx := int(gray*float64(len(context.Grad)-1) + 0.5) // round
			ch := context.Grad[idx]

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
