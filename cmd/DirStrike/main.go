package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/Alpastx/DirStrike/pkg/config"
	"github.com/Alpastx/DirStrike/pkg/scanner"
	"github.com/fatih/color"
)

func main() {
	// Define command line flags
	url := flag.String("u", "", "Target URL (required)")
	wordlist := flag.String("w", "", "Path to wordlist (required)")
	threads := flag.Int("t", 10, "Number of concurrent threads")
	timeout := flag.Int("timeout", 10, "Timeout for HTTP requests in seconds")
	extensions := flag.String("x", "", "File extensions to search (comma separated)")
	userAgent := flag.String("ua", "DirStrike/1.0", "User-Agent string")
	outputFile := flag.String("o", "", "Output file to write results")
	verbose := flag.Bool("v", false, "Verbose output")
	hideUnknownSize := flag.Bool("hide-unknown", false, "Hide results with unknown content length (-1)")
	noColor := flag.Bool("no-color", false, "Disable colored output")

	flag.Parse()

	// Validate required parameters
	if *url == "" || *wordlist == "" {
		fmt.Println("Error: URL (-u) and wordlist (-w) are required")
		flag.Usage()
		os.Exit(1)
	}

	// Create configuration
	cfg := &config.Config{
		URL:             *url,
		Wordlist:        *wordlist,
		Threads:         *threads,
		Timeout:         *timeout,
		Extensions:      *extensions,
		UserAgent:       *userAgent,
		OutputFile:      *outputFile,
		Verbose:         *verbose,
		HideUnknownSize: *hideUnknownSize,
		NoColor:         *noColor,
	}

	// Disable color if requested
	if cfg.NoColor {
		color.NoColor = true
	}

	// Print banner
	printBanner()

	// Print configuration
	fmt.Printf("Target URL: %s\n", cfg.URL)
	fmt.Printf("Wordlist: %s\n", cfg.Wordlist)
	fmt.Printf("Threads: %d\n", cfg.Threads)
	if cfg.HideUnknownSize {
		fmt.Println("Hiding results with unknown content length (-1)")
	}
	fmt.Println("------------------------------------------------")

	// Create and run scanner
	s := scanner.NewScanner(cfg)
	s.Run()
}

func printBanner() {
	banner := color.CyanString(`
	██████╗ ██╗██████╗ ███████╗████████╗██████╗ ██╗██╗  ██╗███████╗
	██╔══██╗██║██╔══██╗██╔════╝╚══██╔══╝██╔══██╗██║██║ ██╔╝██╔════╝
	██║  ██║██║██████╔╝███████╗   ██║   ██████╔╝██║█████╔╝ █████╗  
	██║  ██║██║██╔══██╗╚════██║   ██║   ██╔══██╗██║██╔═██╗ ██╔══╝  
	██████╔╝██║██║  ██║███████║   ██║   ██║  ██║██║██║  ██╗███████╗
	╚═════╝ ╚═╝╚═╝  ╚═╝╚══════╝   ╚═╝   ╚═╝  ╚═╝╚═╝╚═╝  ╚═╝╚══════╝
	                                                                
	`) + color.YellowString("Directory Busting Tool v1.0")

	fmt.Println(banner)
}
