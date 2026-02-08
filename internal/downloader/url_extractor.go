package downloader

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// Define error types
var (
	ErrInvalidURL        = errors.New("无效的URL")
	ErrPageNotFound      = errors.New("页面不存在")
	ErrAudioNotFound     = errors.New("未找到音频文件")
	ErrAccessDenied      = errors.New("访问被拒绝")
	ErrCoverNotFound     = errors.New("未找到封面图片")
	ErrShowNotesNotFound = errors.New("未找到节目show notes")
	ErrInvalidImage      = errors.New("无效的图片文件")
	ErrImageTooLarge     = errors.New("图片文件过大")
	ErrInvalidEncoding   = errors.New("无效的字符编码")
)

// URLExtractor defines the interface for extracting metadata from podcast pages.
type URLExtractor interface {
	// ExtractURL fetches the episode page and extracts metadata.
	//
	// Parameters:
	//   ctx - Context for cancellation and timeout
	//   pageURL - The episode page URL to scrape
	//
	// Returns:
	//   *EpisodeMetadata - Contains audio URL, cover URL, show notes, title, etc.
	//   error - Err if page cannot be fetched or required data not found
	ExtractURL(ctx context.Context, pageURL string) (*EpisodeMetadata, error)
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

// ExtractURL fetches the episode page and extracts metadata.
func (e *HTMLExtractor) ExtractURL(ctx context.Context, pageURL string) (*EpisodeMetadata, error) {
	// Fetch the page
	doc, err := e.client.Get(pageURL)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrPageNotFound, err)
	}

	// Create metadata struct
	metadata := &EpisodeMetadata{}

	// Extract title (required)
	metadata.Title = e.extractTitle(doc)

	// Extract audio URL (required)
	audioURL, err := e.extractAudioURL(doc)
	if err != nil {
		return nil, err
	}
	metadata.AudioURL = audioURL

	// Cover URL and show notes will be added in later phases
	// For now, they remain empty strings

	// Extract cover URL (optional)
	if coverURL := e.extractCoverURL(doc); coverURL != "" {
		metadata.CoverURL = coverURL
	}

	// Extract show notes (optional)
	if showNotesHTML := e.extractShowNotes(doc); showNotesHTML != "" {
		metadata.ShowNotes = showNotesHTML
	}

	return metadata, nil
}

// extractCoverURL extracts the cover image URL from the HTML document.
// Selects the first <img> from .avater-container element (Xiaoyuzhou FM specific).
// IMPORTANT: .avater-container contains TWO images:
//   - First <img>: Episode cover (single episode artwork) ✅ This is what we want
//   - Second <img>: Podcast account cover (channel/series artwork) ❌ Skip this
func (e *HTMLExtractor) extractCoverURL(doc *goquery.Document) string {
	// Find .avater-container and select first img element
	if selection := doc.Find(".avater-container img").First(); selection.Length() > 0 {
		if src, exists := selection.Attr("src"); exists && src != "" {
			return src
		}
	}
	return ""
}

// extractShowNotes extracts show notes content using multi-fallback strategy.
// Strategy: (1) aria-label="节目show notes", (2) aria-label containing "show notes",
// (3) semantic selectors, (4) log failure if not found.
func (e *HTMLExtractor) extractShowNotes(doc *goquery.Document) string {
	// Strategy 1: Search for <section aria-label="节目show notes"> (exact match)
	if selection := doc.Find("section[aria-label=\"节目show notes\"]"); selection.Length() > 0 {
		if html, err := selection.Html(); err == nil && html != "" {
			return html
		}
	}

	// Strategy 2: Search for any element with aria-label containing "show notes" (case-insensitive)
	var foundSelection *goquery.Selection
	doc.Find("*[aria-label]").Each(func(i int, s *goquery.Selection) {
		if ariaLabel, exists := s.Attr("aria-label"); exists {
			if stringsContains(strings.ToLower(ariaLabel), "show notes") {
				foundSelection = s
				return // Stop iteration
			}
		}
	})
	if foundSelection != nil && foundSelection.Length() > 0 {
		if html, err := foundSelection.Html(); err == nil && html != "" {
			return html
		}
	}

	// Strategy 3: Use semantic HTML selectors
	selectors := []string{
		"article.show-notes",
		"section.description",
		"div.description",
		"article p",
		".episode-description",
	}

	for _, selector := range selectors {
		if selection := doc.Find(selector).First(); selection.Length() > 0 {
			if html, err := selection.Html(); err == nil && html != "" {
				return html
			}
		}
	}

	// Strategy 4: All strategies failed - return empty string (no failure needed, show notes are optional)
	return ""
}

// stringsContains is a simple string contains check.
func stringsContains(s, substr string) bool {
	return len(s) >= len(substr) && s != substr && findSubstring(s, substr) != -1 || s == substr
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
