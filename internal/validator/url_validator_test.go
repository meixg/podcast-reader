package validator

import (
	"strings"
	"testing"
)

func TestNewXiaoyuzhouURLValidator(t *testing.T) {
	validator := NewXiaoyuzhouURLValidator()

	if validator == nil {
		t.Fatal("NewXiaoyuzhouURLValidator returned nil")
	}

	if validator.pattern == nil {
		t.Error("validator.pattern is nil")
	}
}

func TestXiaoyuzhouURLValidator_ValidateURL_ValidURLs(t *testing.T) {
	validator := NewXiaoyuzhouURLValidator()

	validURLs := []string{
		"https://www.xiaoyuzhoufm.com/episode/12345678",
		"http://www.xiaoyuzhoufm.com/episode/12345678",
		"https://xiaoyuzhoufm.com/episode/12345678",
		"http://xiaoyuzhoufm.com/episode/abcd1234-5678-90ef-ghij",
		"https://www.xiaoyuzhoufm.com/episode/12345678?param=value",
	}

	for _, url := range validURLs {
		t.Run(url, func(t *testing.T) {
			valid, errMsg := validator.ValidateURL(url)
			if !valid {
				t.Errorf("URL should be valid: %s, error: %s", url, errMsg)
			}
			if errMsg != "" {
				t.Errorf("Error message should be empty for valid URL: %s", errMsg)
			}
		})
	}
}

func TestXiaoyuzhouURLValidator_ValidateURL_InvalidURLs(t *testing.T) {
	validator := NewXiaoyuzhouURLValidator()

	invalidURLs := []struct {
		url      string
		expected string
	}{
		{
			url:      "not-a-url",
			expected: "URL必须使用HTTP或HTTPS协议",
		},
		{
			url:      "ftp://www.xiaoyuzhoufm.com/episode/12345678",
			expected: "URL必须使用HTTP或HTTPS协议",
		},
		{
			url:      "https://www.example.com/episode/12345678",
			expected: "URL格式不正确",
		},
		{
			url:      "https://www.xiaoyuzhoufm.com/podcast/12345678",
			expected: "URL格式不正确",
		},
		{
			url:      "https://www.xiaoyuzhoufm.com/episode/",
			expected: "URL格式不正确",
		},
		{
			url:      "",
			expected: "URL必须使用HTTP或HTTPS协议",
		},
	}

	for _, tc := range invalidURLs {
		t.Run(tc.url, func(t *testing.T) {
			valid, errMsg := validator.ValidateURL(tc.url)
			if valid {
				t.Errorf("URL should be invalid: %s", tc.url)
			}
			if errMsg == "" {
				t.Errorf("Error message should not be empty for invalid URL: %s", tc.url)
			}
			// Check if error message contains expected substring
			if len(tc.expected) > 0 && !strings.Contains(errMsg, tc.expected) {
				t.Errorf("Error message should contain %q, got: %s", tc.expected, errMsg)
			}
		})
	}
}

func TestXiaoyuzhouURLValidator_ValidateURL_EdgeCases(t *testing.T) {
	validator := NewXiaoyuzhouURLValidator()

	tests := []struct {
		name     string
		url      string
		valid    bool
		contains string
	}{
		{
			name:  "URL with fragment",
			url:   "https://www.xiaoyuzhoufm.com/episode/12345678#timestamp",
			valid: true,
		},
		{
			name:     "URL with port",
			url:      "https://www.xiaoyuzhoufm.com:443/episode/12345678",
			valid:    false, // Regex doesn't match port numbers
			contains: "URL格式不正确",
		},
		{
			name:  "URL with subdomain",
			url:   "https://podcast.xiaoyuzhoufm.com/episode/12345678",
			valid: true,
		},
		{
			name:  "URL with multiple subdomains",
			url:   "https://api.www.xiaoyuzhoufm.com/episode/12345678",
			valid: true,
		},
		{
			name:     "missing protocol",
			url:      "www.xiaoyuzhoufm.com/episode/12345678",
			valid:    false,
			contains: "URL必须使用HTTP或HTTPS协议",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			valid, errMsg := validator.ValidateURL(tc.url)
			if valid != tc.valid {
				t.Errorf("ValidateURL(%s) validity = %v, want %v", tc.url, valid, tc.valid)
			}
			if !tc.valid && tc.contains != "" && !strings.Contains(errMsg, tc.contains) {
				t.Errorf("Error message should contain %q, got: %s", tc.contains, errMsg)
			}
		})
	}
}
