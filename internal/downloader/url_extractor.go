package downloader

import (
	"context"
	"errors"
	"fmt"

	"github.com/PuerkitoBio/goquery"
)

// Define error types
var (
	ErrInvalidURL    = errors.New("无效的URL")
	ErrPageNotFound  = errors.New("页面不存在")
	ErrAudioNotFound = errors.New("未找到音频文件")
	ErrAccessDenied  = errors.New("访问被拒绝")
)

// URLExtractor defines the interface for extracting audio URLs from web pages.
type URLExtractor interface {
	// ExtractURL fetches the episode page and extracts the direct audio file URL.
	//
	// Parameters:
	//   ctx - Context for cancellation and timeout
	//   pageURL - The episode page URL to scrape
	//
	// Returns:
	//   string - The direct audio file URL (.m4a)
	//   string - The episode title for filename generation
	//   error - Err if page cannot be fetched or audio URL not found
	ExtractURL(ctx context.Context, pageURL string) (audioURL string, title string, err error)
}

// HTMLExtractor implements URLExtractor using goquery for HTML parsing.
type HTMLExtractor struct {
	// client is the HTTP client to use for fetching pages
	client Doer
}

// Doer is the interface for HTTP GET requests.
type Doer interface {
	Get(url string) (*goquery.Document, error)
}

// NewHTMLExtractor creates a new HTML extractor.
func NewHTMLExtractor(client Doer) *HTMLExtractor {
	return &HTMLExtractor{
		client: client,
	}
}

// ExtractURL fetches the episode page and extracts the direct audio file URL.
func (e *HTMLExtractor) ExtractURL(ctx context.Context, pageURL string) (string, string, error) {
	// Fetch the page
	doc, err := e.client.Get(pageURL)
	if err != nil {
		return "", "", fmt.Errorf("%w: %v", ErrPageNotFound, err)
	}

	// Try to extract title
	title := e.extractTitle(doc)

	// Try to extract audio URL from multiple selectors
	audioURL, err := e.extractAudioURL(doc)
	if err != nil {
		return "", title, err
	}

	return audioURL, title, nil
}

// extractTitle extracts the episode title from the HTML document.
func (e *HTMLExtractor) extractTitle(doc *goquery.Document) string {
	// Try <title> tag first
	title := doc.Find("title").First().Text()
	if title != "" {
		return title
	}

	// Try common meta tags
	if title, exists := doc.Find("meta[property='og:title']").Attr("content"); exists {
		return title
	}

	if title, exists := doc.Find("meta[name='title']").Attr("content"); exists {
		return title
	}

	return ""
}

// extractAudioURL extracts the audio file URL from the HTML document.
func (e *HTMLExtractor) extractAudioURL(doc *goquery.Document) (string, error) {
	// Try meta tag with og:audio property (Xiaoyuzhou FM specific)
	if audioURL, exists := doc.Find("meta[property='og:audio']").Attr("content"); exists && len(audioURL) > 0 {
		return audioURL, nil
	}

	// Try JSON-LD structured data as fallback
	var jsonLDURL string
	doc.Find("script[type='application/ld+json']").Each(func(i int, s *goquery.Selection) {
		jsonText := s.Text()
		if len(jsonText) > 0 && jsonLDURL == "" {
			// Look for "contentUrl":"https://..." pattern
			// This is a simple regex-based approach for extraction
			if idx := findSubstring(jsonText, `"contentUrl":"`); idx != -1 {
				start := idx + len(`"contentUrl":"`)
				end := findSubstring(jsonText[start:], `"`)
				if end != -1 {
					jsonLDURL = jsonText[start : start+end]
				}
			}
		}
	})
	if jsonLDURL != "" {
		return jsonLDURL, nil
	}

	// Try common selectors for audio elements as fallback
	selectors := []string{
		"audio[src]",
		"source[src]",
		".audio-player audio",
		"[data-audio-url]",
	}

	for _, selector := range selectors {
		if audioURL, exists := doc.Find(selector).First().Attr("src"); exists && len(audioURL) > 0 {
			return audioURL, nil
		}
	}

	return "", ErrAudioNotFound
}

// findSubstring is a helper to find substring index (simple replacement for strings.Index)
func findSubstring(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}
