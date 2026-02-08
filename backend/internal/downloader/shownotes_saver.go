package downloader

import (
	"fmt"
	"os"
	"strings"
	"unicode/utf8"

	"github.com/PuerkitoBio/goquery"
)

// ShowNotesSaver defines the interface for saving show notes.
type ShowNotesSaver interface {
	// Save saves show notes content to a file with UTF-8-BOM encoding.
	//
	// Parameters:
	//   content - The HTML content to format and save
	//   destPath - The destination file path
	//
	// Returns:
	//   error - Error if saving fails
	Save(content string, destPath string) error

	// FormatHTMLToText converts HTML show notes to plain text.
	//
	// Parameters:
	//   html - The HTML content to convert
	//
	// Returns:
	//   string - Plain text with preserved structure
	FormatHTMLToText(html string) string
}

// PlainTextShowNotesSaver implements ShowNotesSaver for plain text files.
type PlainTextShowNotesSaver struct{}

// NewPlainTextShowNotesSaver creates a new plain text show notes saver.
func NewPlainTextShowNotesSaver() *PlainTextShowNotesSaver {
	return &PlainTextShowNotesSaver{}
}

// Save saves show notes to a file with UTF-8-BOM encoding.
func (s *PlainTextShowNotesSaver) Save(content string, destPath string) error {
	// Validate UTF-8 encoding
	if !utf8.ValidString(content) {
		return fmt.Errorf("%w: show notes content contains invalid UTF-8 sequences", ErrInvalidEncoding)
	}

	// Convert HTML to plain text
	text := s.FormatHTMLToText(content)

	// Add UTF-8 BOM for Windows compatibility
	bom := []byte{0xEF, 0xBB, 0xBF}
	fileContent := string(bom) + text

	// Write to file
	if err := os.WriteFile(destPath, []byte(fileContent), 0644); err != nil {
		return fmt.Errorf("写入show notes文件失败: %w", err)
	}

	return nil
}

// FormatHTMLToText converts HTML show notes to readable plain text.
func (s *PlainTextShowNotesSaver) FormatHTMLToText(html string) string {
	// Parse HTML
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		// If parsing fails, return original content
		return html
	}

	var result strings.Builder

	// Process each node
	doc.Find("body").Contents().Each(func(i int, sel *goquery.Selection) {
		result.WriteString(s.processNode(sel))
	})

	// Clean up excessive whitespace
	text := result.String()
	text = s.cleanupWhitespace(text)

	return text
}

// processNode converts a goquery node to plain text.
func (s *PlainTextShowNotesSaver) processNode(sel *goquery.Selection) string {
	if sel.Length() == 0 {
		return ""
	}

	// Get node name
	nodeName := goquery.NodeName(sel)
	if nodeName == "" {
		return sel.Text()
	}

	// Handle different element types
	switch nodeName {
	case "a":
		// Convert links to: text (URL: url)
		text := sel.Text()
		if href, exists := sel.Attr("href"); exists && href != "" {
			return fmt.Sprintf("%s (URL: %s)", text, href)
		}
		return text

	case "ul":
		var items []string
		sel.Find("li").Each(func(i int, li *goquery.Selection) {
			items = append(items, fmt.Sprintf("• %s", strings.TrimSpace(li.Text())))
		})
		return strings.Join(items, "\n") + "\n"

	case "ol":
		var items []string
		sel.Find("li").Each(func(i int, li *goquery.Selection) {
			items = append(items, fmt.Sprintf("%d. %s", i+1, strings.TrimSpace(li.Text())))
		})
		return strings.Join(items, "\n") + "\n"

	case "h1", "h2", "h3", "h4", "h5", "h6":
		text := strings.ToUpper(strings.TrimSpace(sel.Text()))
		underline := strings.Repeat("=", len(text))
		return fmt.Sprintf("%s\n%s\n\n", text, underline)

	case "p":
		return sel.Text() + "\n\n"

	case "br":
		return "\n"

	case "blockquote":
		lines := strings.Split(sel.Text(), "\n")
		var quoted []string
		for _, line := range lines {
			quoted = append(quoted, "> "+line)
		}
		return strings.Join(quoted, "\n") + "\n"

	case "strong", "b":
		return fmt.Sprintf("**%s**", sel.Text())

	case "em", "i":
		return fmt.Sprintf("*%s*", sel.Text())

	case "code":
		return fmt.Sprintf("`%s`", sel.Text())

	case "pre":
		return fmt.Sprintf("```\n%s\n```\n", sel.Text())

	case "div", "span", "section", "article":
		// Process children
		var result strings.Builder
		sel.Contents().Each(func(i int, child *goquery.Selection) {
			result.WriteString(s.processNode(child))
		})
		return result.String()

	default:
		// Default: just get text content
		return sel.Text()
	}
}

// cleanupWhitespace removes excessive newlines and trims spaces.
func (s *PlainTextShowNotesSaver) cleanupWhitespace(text string) string {
	// Replace multiple consecutive newlines with at most 2
	lines := strings.Split(text, "\n")
	var cleaned []string
	emptyCount := 0

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			emptyCount++
			if emptyCount <= 2 { // Keep up to 2 consecutive empty lines
				cleaned = append(cleaned, "")
			}
		} else {
			emptyCount = 0
			cleaned = append(cleaned, trimmed)
		}
	}

	// Join and trim
	result := strings.Join(cleaned, "\n")
	result = strings.Trim(result, "\n")

	return result
}
