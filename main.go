package main

import (
	"flag"
	"fmt"
	"os"
	"slices"
	"strings"
)

type Options struct {
	FilePath string
	Styles   string
	Width    int
	Height   int
	Color    bool
}

func DefaultOptions() Options {
	return Options{
		Width:  90,
		Height: 90,
		Color:  false,
		Styles: "ascii2",
	}
}

func main() {
	opts := DefaultOptions()

	help := flag.Bool("help", false, "Show this help message")
	flag.StringVar(&opts.Styles, "s", "", "Styles to use (optional)")
	flag.IntVar(&opts.Width, "w", 90, "Set output width (optional)")
	flag.IntVar(&opts.Height, "h", 90, "Set output height (optional)")
	flag.BoolVar(&opts.Color, "c", false, "Enable color output (optional)")
	flag.IntVar(&opts.Width, "width", 90, "Set output width (optional)")
	flag.IntVar(&opts.Height, "height", 90, "Set output height (optional)")
	flag.BoolVar(&opts.Color, "color", false, "Enable color output (optional)")
	flag.Parse()

	if slices.Contains(os.Args[1:], "--help") {
		*help = true
	}

	if *help {
		printHelp()
		return
	}

	opts.FilePath = flag.Arg(0)
	if opts.FilePath == "" {
		fmt.Println("Error: <file_path>.gif is required")
		printHelp()
		os.Exit(1)
	}

	if !strings.HasSuffix(strings.ToLower(opts.FilePath), ".gif") {
		fmt.Println("Error: File must have a .gif extension")
		os.Exit(1)
	}

	if opts.Styles == "" {
		fmt.Println("Using default styles: ascii2")
		opts.Styles = "ascii2"
	}

	if opts.Width < 0 || opts.Height < 0 {
		fmt.Println("Error: Width and height must be positive")
		os.Exit(1)
	}

	execute(opts)
}

func printHelp() {
	fmt.Println("Usage:")
	fmt.Println("  gifter [options] <file_path>.gif")
	fmt.Println()
	fmt.Println("Options:")
	fmt.Println("  -s, --styles=\"styles\" Styles to use (optional)")
	fmt.Println("  -w, --width=WIDTH      Set output width (optional, default: 90)")
	fmt.Println("  -h, --height=HEIGHT    Set output height (optional, default: 90)")
	fmt.Println("  -c, --color            Enable color output (optional, default: true)")
	fmt.Println("  --help                 Show this help message")
	fmt.Println()
	fmt.Println("Example:")
	fmt.Println("  gifter --styles=ascii2 --width=800 --height=600 --color path/to/file.gif")
}
