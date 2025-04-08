package config

// Config holds the application configuration
type Config struct {
	URL           string
	Wordlist      string
	Threads       int
	Timeout       int
	Extensions    string
	UserAgent     string
	OutputFile    string
	Verbose       bool
	HideUnknownSize bool // New option to hide -1 content length
	NoColor       bool // Option to disable colored output
}
