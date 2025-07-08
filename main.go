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
	Mode     string
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
		Mode:   "ascii",
	}
}

func main() {
	opts := DefaultOptions()

	help := flag.Bool("help", false, "Show this help message")
	flag.StringVar(&opts.Styles, "s", opts.Styles, "Styles to use for ASCII mode: normal, ascii2, shaded, bordered, blocky (optional, default: ascii2)")
	flag.StringVar(&opts.Mode, "m", opts.Mode, "Output mode: ascii or graphic (optional, default: ascii)")
	flag.IntVar(&opts.Width, "w", 90, "Set output width (optional)")
	flag.IntVar(&opts.Height, "h", 90, "Set output height (optional)")
	flag.BoolVar(&opts.Color, "c", false, "Enable color output for ASCII mode (optional)")
	flag.IntVar(&opts.Width, "width", 90, "Set output width (optional)")
	flag.IntVar(&opts.Height, "height", 90, "Set output height (optional)")
	flag.BoolVar(&opts.Color, "color", false, "Enable color output for ASCII mode (optional)")
	flag.StringVar(&opts.Mode, "mode", opts.Mode, "Output mode: ascii or graphic (optional, default: ascii)")
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

	fmt.Printf("Processing file: %s\n", opts.FilePath)

	if !strings.HasSuffix(strings.ToLower(opts.FilePath), ".gif") &&
		!strings.HasPrefix(strings.ToLower(opts.FilePath), "http://") &&
		!strings.HasPrefix(strings.ToLower(opts.FilePath), "https://") {
		fmt.Println("Error: Input must be a .gif file or a URL starting with http:// or https://")
		os.Exit(1)
	}

	if opts.Mode != "ascii" && opts.Mode != "graphic" {
		fmt.Printf("Error: Invalid mode %q. Must be 'ascii' or 'graphic'\n", opts.Mode)
		os.Exit(1)
	}

	if opts.Styles == "" && opts.Mode == "ascii" {
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
	fmt.Println("  gifter [options] <file_path>.gif | <URL>")
	fmt.Println()
	fmt.Println("Options:")
	fmt.Println("  -m, --mode=MODE        Output mode: ascii or graphic (optional, default: ascii)")
	fmt.Println("  -s, --styles=STYLE     Styles to use for ASCII mode: normal, ascii2, shaded, bordered, blocky (optional, default: ascii2)")
	fmt.Println("  -w, --width=WIDTH      Set output width (optional, default: 90)")
	fmt.Println("  -h, --height=HEIGHT    Set output height (optional, default: 90)")
	fmt.Println("  -c, --color            Enable color output for ASCII mode (optional, default: false)")
	fmt.Println("  --help                 Show this help message")
	fmt.Println()
	fmt.Println("Example:")
	fmt.Println("  gifter --mode=ascii --styles=ascii2 --width=80 --height=40 path/to/file.gif")
	fmt.Println("  gifter --mode=graphic https://media.giphy.com/media/DvyLQztQwmyAM/giphy.gif")
	fmt.Println("  gifter --mode=graphic --height=120 --width=90 https://media.giphy.com/media/DvyLQztQwmyAM/giphy.gif")
}
