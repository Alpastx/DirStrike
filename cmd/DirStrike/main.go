package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/Alpastx/DirStrike/pkg/config"
	"github.com/Alpastx/DirStrike/pkg/scanner"
)

func main() {
	// Define command line flags
	url := flag.String("u", "", "Target URL (required)")
	wordlist := flag.String("w", "", "Path to wordlist (required)")
	threads := flag.Int("t", 10, "Number of concurrent threads")
	timeout := flag.Int("timeout", 10, "Timeout for HTTP requests in seconds")
	extensions := flag.String("x", "", "File extensions to search (comma separated)")
	userAgent := flag.String("ua", "Dirbuster/1.0", "User-Agent string")
	outputFile := flag.String("o", "", "Output file to write results")
	verbose := flag.Bool("v", false, "Verbose output")

	flag.Parse()

	// Validate required parameters
	if *url == "" || *wordlist == "" {
		fmt.Println("Error: URL (-u) and wordlist (-w) are required")
		flag.Usage()
		os.Exit(1)
	}

	// Create configuration
	cfg := &config.Config{
		URL:        *url,
		Wordlist:   *wordlist,
		Threads:    *threads,
		Timeout:    *timeout,
		Extensions: *extensions,
		UserAgent:  *userAgent,
		OutputFile: *outputFile,
		Verbose:    *verbose,
	}

	// Print banner
	printBanner()

	// Print configuration
	fmt.Printf("Target URL: %s\n", cfg.URL)
	fmt.Printf("Wordlist: %s\n", cfg.Wordlist)
	fmt.Printf("Threads: %d\n", cfg.Threads)
	fmt.Println("------------------------------------------------")

	// Create and run scanner
	s := scanner.NewScanner(cfg)
	s.Run()
}

func printBanner() {
	banner := `
	██████╗ ██╗██████╗ ██████╗ ██╗   ██╗███████╗████████╗███████╗██████╗ 
	██╔══██╗██║██╔══██╗██╔══██╗██║   ██║██╔════╝╚══██╔══╝██╔════╝██╔══██╗
	██║  ██║██║██████╔╝██████╔╝██║   ██║███████╗   ██║   █████╗  ██████╔╝
	██║  ██║██║██╔══██╗██╔══██╗██║   ██║╚════██║   ██║   ██╔══╝  ██╔══██╗
	██████╔╝██║██║  ██║██████╔╝╚██████╔╝███████║   ██║   ███████╗██║  ██║
	╚═════╝ ╚═╝╚═╝  ╚═╝╚═════╝  ╚═════╝ ╚══════╝   ╚═╝   ╚══════╝╚═╝  ╚═╝
	                                                                     
	Directory Busting Tool v1.0
	`
	fmt.Println(banner)
}
