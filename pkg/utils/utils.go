package utils

import (
	"strings"
)

// JoinURL joins base URL with path properly
func JoinURL(baseURL, path string) string {
	baseURL = strings.TrimSuffix(baseURL, "/")
	path = strings.TrimPrefix(path, "/")
	return baseURL + "/" + path
}
