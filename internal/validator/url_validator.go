package validator

import (
	"net/url"
	"regexp"
)

// URLValidator defines the interface for validating URLs.
type URLValidator interface {
	// ValidateURL checks if a URL is valid and meets requirements.
	//
	// Parameters:
	//   url - The URL to validate
	//
	// Returns:
	//   bool - True if URL is valid
	//   string - Error message if invalid, empty if valid
	ValidateURL(url string) (isValid bool, errMsg string)
}

// XiaoyuzhouURLValidator validates Xiaoyuzhou FM episode URLs.
type XiaoyuzhouURLValidator struct {
	// Pattern matches: *.xiaoyuzhoufm.com/episode/{episode_id}
	pattern *regexp.Regexp
}

// NewXiaoyuzhouURLValidator creates a new validator for Xiaoyuzhou FM URLs.
func NewXiaoyuzhouURLValidator() *XiaoyuzhouURLValidator {
	// Pattern matches: http(s)://*.xiaoyuzhoufm.com/episode/{episode_id}
	pattern := regexp.MustCompile(`^https?://[a-z0-9\.]*xiaoyuzhoufm\.com/episode/[a-z0-9]+$`)
	return &XiaoyuzhouURLValidator{
		pattern: pattern,
	}
}

// ValidateURL implements URLValidator for Xiaoyuzhou FM URLs.
func (v *XiaoyuzhouURLValidator) ValidateURL(urlStr string) (bool, string) {
	// Check if URL is well-formed
	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		return false, "URL格式无效: " + err.Error()
	}

	// Check protocol
	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		return false, "URL必须使用HTTP或HTTPS协议"
	}

	// Check domain matches Xiaoyuzhou FM
	if !v.pattern.MatchString(urlStr) {
		return false, "URL格式不正确，应为: https://www.xiaoyuzhoufm.com/episode/{episode_id}"
	}

	return true, ""
}
