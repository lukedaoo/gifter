package main

import (
	"fmt"
)

func execute(filePath string, styles string) {
	fmt.Printf("Processing file: %s\n", filePath)
	if styles != "" {
		fmt.Printf("Using styles: %s\n", styles)
	} else {
		fmt.Println("No styles provided, using default")
		styles = "normal"
	}

	fmt.Printf("Processing complete for file: %s with styles: %s\n", filePath, styles)
}
