package tests

import (
	"testing"

	"github.com/Alpastx/DirStrike/pkg/utils"
)

func TestJoinURL(t *testing.T) {
	tests := []struct {
		baseURL  string
		path     string
		expected string
	}{
		{"http://example.com", "test", "http://example.com/test"},
		{"http://example.com/", "test", "http://example.com/test"},
		{"http://example.com", "/test", "http://example.com/test"},
		{"http://example.com/", "/test", "http://example.com/test"},
	}

	for _, test := range tests {
		result := utils.JoinURL(test.baseURL, test.path)
		if result != test.expected {
			t.Errorf("JoinURL(%s, %s) = %s; want %s",
				test.baseURL, test.path, result, test.expected)
		}
	}
}
