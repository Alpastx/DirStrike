package scanner

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/Alpastx/DirStrike/pkg/config"
	"github.com/Alpastx/DirStrike/pkg/utils"
	"github.com/fatih/color"
)

// Scanner represents a directory scanner
type Scanner struct {
	config  *config.Config
	client  *http.Client
	results []Result
	mutex   sync.Mutex
}

// Result represents a scan result
type Result struct {
	URL        string
	StatusCode int
	Size       int64
}

// NewScanner creates a new scanner
func NewScanner(cfg *config.Config) *Scanner {
	client := &http.Client{
		Timeout: time.Duration(cfg.Timeout) * time.Second,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	return &Scanner{
		config:  cfg,
		client:  client,
		results: []Result{},
	}
}

// Run starts the scanning process
func (s *Scanner) Run() {
	// Open wordlist file
	file, err := os.Open(s.config.Wordlist)
	if err != nil {
		color.Red("Error opening wordlist: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	// Read wordlist into a slice
	var paths []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		path := scanner.Text()
		if path != "" {
			paths = append(paths, path)
			
			// Add extensions if specified
			if s.config.Extensions != "" {
				exts := strings.Split(s.config.Extensions, ",")
				for _, ext := range exts {
					if !strings.HasPrefix(ext, ".") {
						ext = "." + ext
					}
					paths = append(paths, path+ext)
				}
			}
		}
	}

	if err := scanner.Err(); err != nil {
		color.Red("Error reading wordlist: %v\n", err)
		os.Exit(1)
	}

	// Create a channel to receive paths to check
	pathChan := make(chan string)

	// Create a wait group to wait for all goroutines to finish
	var wg sync.WaitGroup

	// Start worker goroutines
	for i := 0; i < s.config.Threads; i++ {
		wg.Add(1)
		go s.worker(pathChan, &wg)
	}

	// Send paths to the channel
	for _, path := range paths {
		pathChan <- path
	}
	close(pathChan)

	// Wait for all goroutines to finish
	wg.Wait()
	
	// Write results to file if specified
	if s.config.OutputFile != "" {
		s.writeResults()
	}
	
	color.Green("Directory busting completed!")
}

func (s *Scanner) worker(paths <-chan string, wg *sync.WaitGroup) {
	defer wg.Done()

	for path := range paths {
		fullURL := utils.JoinURL(s.config.URL, path)
		
		req, err := http.NewRequest("GET", fullURL, nil)
		if err != nil {
			continue
		}
		
		req.Header.Set("User-Agent", s.config.UserAgent)
		
		resp, err := s.client.Do(req)
		if err != nil {
			if s.config.Verbose {
				color.Red("Error: %s - %v\n", fullURL, err)
			}
			continue
		}
		
		size := resp.ContentLength
		resp.Body.Close()

		// Check if the response status code indicates a successful or redirection response
		if resp.StatusCode >= 200 && resp.StatusCode < 404 {
			// Skip if we're hiding unknown size and size is -1
			if s.config.HideUnknownSize && size == -1 {
				continue
			}
			
			result := Result{
				URL:        fullURL,
				StatusCode: resp.StatusCode,
				Size:       size,
			}
			
			s.mutex.Lock()
			s.results = append(s.results, result)
			s.mutex.Unlock()
			
			// Print with color based on status code
			s.printResult(result)
		}
	}
}

func (s *Scanner) printResult(result Result) {
	statusStr := fmt.Sprintf("[%d]", result.StatusCode)
	sizeStr := fmt.Sprintf("%-8d", result.Size)
	
	// Color based on status code
	var statusColor *color.Color
	switch {
	case result.StatusCode >= 200 && result.StatusCode < 300:
		statusColor = color.New(color.FgGreen)
	case result.StatusCode >= 300 && result.StatusCode < 400:
		statusColor = color.New(color.FgCyan)
	case result.StatusCode >= 400 && result.StatusCode < 500:
		statusColor = color.New(color.FgYellow)
	default:
		statusColor = color.New(color.FgRed)
	}
	
	// Color for size
	sizeColor := color.New(color.FgWhite)
	if result.Size == -1 {
		sizeColor = color.New(color.FgMagenta)
	}
	
	// URL color
	urlColor := color.New(color.FgHiWhite)
	
	fmt.Printf("%s %s %s\n", 
		statusColor.Sprint(statusStr),
		sizeColor.Sprint(sizeStr),
		urlColor.Sprint(result.URL))
}

func (s *Scanner) writeResults() {
	file, err := os.Create(s.config.OutputFile)
	if err != nil {
		color.Red("Error creating output file: %v\n", err)
		return
	}
	defer file.Close()
	
	for _, result := range s.results {
		file.WriteString(fmt.Sprintf("[%d] %-8d %s\n", result.StatusCode, result.Size, result.URL))
	}
	
	color.Green("Results written to %s\n", s.config.OutputFile)
}
