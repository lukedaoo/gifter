package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"io"
	"log"
	"math"
)

func testGraphics(w io.Writer) error {
	img := image.NewRGBA(image.Rect(0, 0, 10, 10))
	for x := range 10 {
		for y := range 10 {
			img.Set(x, y, color.RGBA{R: 255, G: 0, B: 0, A: 255})
		}
	}

	var buf bytes.Buffer
	err := png.Encode(&buf, img)
	if err != nil {
		return fmt.Errorf("Error encoding test image to PNG: %v", err)
	}

	log.Printf("Encoded test PNG: %d bytes", buf.Len())

	fmt.Fprint(w, "\033[H")
	err = kittyWritePNG(w, buf.Bytes(), 10, 10, 1)
	if err != nil {
		return fmt.Errorf("Error writing test PNG: %v", err)
	}

	return nil
}

func kittyWritePNG(w io.Writer, pngBytes []byte, width, height, imageID int) error {
	const chunkSize = 4096
	encoded := base64.StdEncoding.EncodeToString(pngBytes)

	for off := 0; off < len(encoded); off += chunkSize {
		end := off + chunkSize
		end = int(math.Min(float64(end), float64(len(encoded))))
		chunk := encoded[off:end]

		var err error
		switch {
		case off == 0 && end == len(encoded):
			_, err = fmt.Fprintf(w, "\033_Ga=T,f=100,s=%d,v=%d,i=%d,q=2;%s\033\\", width, height, imageID, chunk)
		case off == 0:
			_, err = fmt.Fprintf(w, "\033_Ga=T,f=100,s=%d,v=%d,i=%d,m=1,q=2;%s\033\\", width, height, imageID, chunk)
		case end == len(encoded):
			_, err = fmt.Fprintf(w, "\033_Gm=0,i=%d,q=2;%s\033\\", imageID, chunk)
		default:
			_, err = fmt.Fprintf(w, "\033_Gm=1,i=%d,q=2;%s\033\\", imageID, chunk)
		}
		if err != nil {
			log.Printf("Error writing PNG chunk %d-%d (%d bytes): %v", off, end, len(chunk), err)
			return err
		}
		// log.Printf("Wrote PNG chunk %d-%d (%d bytes, imageID=%d)", off, end, len(chunk), imageID)
	}
	return nil
}

func displayImage(w io.Writer, img image.Image, width, height, imageID int) error {
	var buf bytes.Buffer
	err := png.Encode(&buf, img)
	if err != nil {
		return fmt.Errorf("Error encoding frame to PNG: %v", err)
	}
	if buf.Len() == 0 {
		return fmt.Errorf("Encoded PNG frame is empty")
	}
	fmt.Fprint(w, "\033[H")
	err = kittyWritePNG(w, buf.Bytes(), width, height, imageID)
	if err != nil {
		return fmt.Errorf("Error writing PNG: %v", err)
	}

	return nil
}
