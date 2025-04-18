# Dirbuster Usage Guide

Dirbuster is a powerful directory busting tool written in Go that helps discover hidden directories and files on web servers.

## Basic Usage

```bash
./dirbuster -u http://example.com -w /path/to/wordlist.txt
```

## Available Options

| Option     | Description                                 | Default       |
| ---------- | ------------------------------------------- | ------------- |
| `-u`       | Target URL (required)                       | -             |
| `-w`       | Path to wordlist (required)                 | -             |
| `-t`       | Number of concurrent threads                | 10            |
| `-timeout` | Timeout for HTTP requests in seconds        | 10            |
| `-x`       | File extensions to search (comma separated) | -             |
| `-ua`      | User-Agent string                           | Dirbuster/1.0 |
| `-o`       | Output file to write results                | -             |
| `-v`       | Verbose output                              | -             |

## Examples

### Scan with 20 threads

```bash
./dirbuster -u http://example.com -w /path/to/wordlist.txt -t 20
```

### Scan with specific extensions

```bash
./dirbuster -u http://example.com -w /path/to/wordlist.txt -x php,html,txt
```

### Save results to a file

```bash
./dirbuster -u http://example.com -w /path/to/wordlist.txt -o results.txt
```
