package main

import (
	"flag"
	"fmt"
	"os"
	"slices"
	"strings"
)

func main() {
	help := flag.Bool("h", false, "Show this help message")
	styles := flag.String("s", "", "Styles to use (optional)")

	flag.Parse()

	if slices.Contains(os.Args[1:], "--help") {
		*help = true
	}

	if *help {
		printHelp()
		return
	}

	filePath := flag.Arg(0)
	if filePath == "" {
		fmt.Println("Error: <file_path>.gif is required")
		printHelp()
		os.Exit(1)
	}

	if !strings.HasSuffix(strings.ToLower(filePath), ".gif") {
		fmt.Println("Error: File must have a .gif extension")
		os.Exit(1)
	}

	execute(filePath, *styles)
}

func printHelp() {
	fmt.Println("Usage:")
	fmt.Println("  gifter [options] <file_path>.gif")
	fmt.Println()
	fmt.Println("Options:")
	fmt.Println("  -s, --styles=\"styles\" Styles to use (optional)")
	fmt.Println("  -h, --help              Show this help message")
	fmt.Println("Example:")
	fmt.Println("  gifter example.gif --styles=ascii2")
}
