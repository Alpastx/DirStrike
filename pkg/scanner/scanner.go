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
)


type Scanner struct {
	config *config.Config
	client *http.Client
	results []Result
	mutex   sync.Mutex
}

type Result struct {
	URL        string
	StatusCode int
	Size       int64
}

func NewScanner(cfg *config.Config) *Scanner {
	client := &http.Client{
		Timeout: time.Duration(cfg.Timeout) * time.Second,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	return &Scanner{
		config: cfg,
		client: client,
		results: []Result{},
	}
}


func (s *Scanner) Run() {
	// Open wordlist file
	file, err := os.Open(s.config.Wordlist)
	if err != nil {
		fmt.Printf("Error opening wordlist: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	// Read 
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
		fmt.Printf("Error reading wordlist: %v\n", err)
		os.Exit(1)
	}

	
	pathChan := make(chan string)

	var wg sync.WaitGroup

	for i := 0; i < s.config.Threads; i++ {
		wg.Add(1)
		go s.worker(pathChan, &wg)
	}

	for _, path := range paths {
		pathChan <- path
	}
	close(pathChan)

	wg.Wait()
	
	if s.config.OutputFile != "" {
		s.writeResults()
	}
	
	fmt.Println("Directory busting completed!")
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
				fmt.Printf("Error: %s - %v\n", fullURL, err)
			}
			continue
		}
		
		size := resp.ContentLength
		resp.Body.Close()

		if resp.StatusCode >= 200 && resp.StatusCode < 404 {
			result := Result{
				URL:        fullURL,
				StatusCode: resp.StatusCode,
				Size:       size,
			}
			
			s.mutex.Lock()
			s.results = append(s.results, result)
			s.mutex.Unlock()
			
			fmt.Printf("[%d] %-8d %s\n", resp.StatusCode, size, fullURL)
		}
	}
}

func (s *Scanner) writeResults() {
	file, err := os.Create(s.config.OutputFile)
	if err != nil {
		fmt.Printf("Error creating output file: %v\n", err)
		return
	}
	defer file.Close()
	
	for _, result := range s.results {
		file.WriteString(fmt.Sprintf("[%d] %-8d %s\n", result.StatusCode, result.Size, result.URL))
	}
	
	fmt.Printf("Results written to %s\n", s.config.OutputFile)
}
