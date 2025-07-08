package main

import (
	"bytes"
	"fmt"
	"image/gif"
	"io"
	"net/http"
	"strings"
)

func downloadGIF(url string) (*gif.GIF, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("Error fetching GIF from %s: %v", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Unexpected status code %d for %s", resp.StatusCode, url)
	}

	// Check Content-Type to ensure it's a GIF
	contentType := resp.Header.Get("Content-Type")
	if !strings.Contains(contentType, "image/gif") {
		return nil, fmt.Errorf("URL %s does not point to a GIF (Content-Type: %s)", url, contentType)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Error reading GIF data from %s: %v", url, err)
	}

	g, err := gif.DecodeAll(bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("Error decoding GIF from %s: %v", url, err)
	}

	return g, nil
}
