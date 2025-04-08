package config

type Config struct {
	URL        string
	Wordlist   string
	Threads    int
	Timeout    int
	Extensions string
	UserAgent  string
	OutputFile string
	Verbose    bool
}
