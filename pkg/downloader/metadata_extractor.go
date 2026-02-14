package downloader

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/meixg/podcast-reader/pkg/models"
)

// MetadataExtractor extracts podcast metadata (duration and publish time) from HTML pages
type MetadataExtractor struct {
	client Doer
}

// NewMetadataExtractor creates a new metadata extractor
func NewMetadataExtractor(client Doer) *MetadataExtractor {
	return &MetadataExtractor{
		client: client,
	}
}

// ExtractMetadata extracts duration and publish time from the podcast page
// Returns nil metadata if extraction fails (graceful degradation)
func (e *MetadataExtractor) ExtractMetadata(ctx context.Context, pageURL string) (*models.PodcastMetadata, error) {
	// Fetch the page
	doc, err := e.client.Get(pageURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch page: %w", err)
	}

	metadata := models.NewPodcastMetadata()

	// Extract combined info text and parse it
	combinedInfo := e.extractCombinedInfoText(doc)
	if combinedInfo != "" {
		// Parse the combined text to extract duration and publish time
		metadata.Duration, metadata.PublishTime = e.parseInfoText(combinedInfo)
	}

	// Extract episode title
	metadata.EpisodeTitle = e.extractEpisodeTitle(doc)

	// Extract podcast name
	metadata.PodcastName = e.extractPodcastName(doc)

	return metadata, nil
}

// extractCombinedInfoText extracts the combined info text from elements with class containing "info"
func (e *MetadataExtractor) extractCombinedInfoText(doc *goquery.Document) string {
	// Find elements with class containing "info"
	var infoTexts []string
	doc.Find("[class*='info']").Each(func(i int, s *goquery.Selection) {
		text := strings.TrimSpace(s.Text())
		if text != "" && containsNumber(text) {
			infoTexts = append(infoTexts, text)
		}
	})

	// Return the first info text that contains both duration and time patterns
	for _, text := range infoTexts {
		if (strings.Contains(text, "分钟") || strings.Contains(text, "小时")) && isTimePattern(text) {
			return text
		}
	}

	// Fallback: return first info text with numbers
	if len(infoTexts) > 0 {
		return infoTexts[0]
	}

	return ""
}

// parseInfoText parses the combined info text to extract duration and publish time
// Expected format: "103分钟 ·2个月前35609·415" or similar
func (e *MetadataExtractor) parseInfoText(text string) (duration, publishTime string) {
	// Use regex to extract duration (e.g., "103分钟", "1小时15分钟")
	durationRegex := regexp.MustCompile(`\d+\s*[小时分分钟]+\s*[分钟分]?`)
	durationMatch := durationRegex.FindString(text)
	if durationMatch != "" {
		duration = strings.TrimSpace(durationMatch)
	}

	// Use regex to extract publish time (e.g., "2个月前", "刚刚发布", "3天前")
	publishTimeRegex := regexp.MustCompile(`(\d+\s*[个]?[周月年天小时分]+前|刚刚发布|[昨前]天)`)
	publishTimeMatch := publishTimeRegex.FindString(text)
	if publishTimeMatch != "" {
		publishTime = strings.TrimSpace(publishTimeMatch)
	}

	return duration, publishTime
}

// extractEpisodeTitle extracts the episode title from the page
func (e *MetadataExtractor) extractEpisodeTitle(doc *goquery.Document) string {
	// Try h1 first (usually episode title)
	title := strings.TrimSpace(doc.Find("h1").First().Text())
	if title != "" {
		return title
	}

	// Try meta tags
	if title, exists := doc.Find("meta[property='og:title']").Attr("content"); exists {
		return strings.TrimSpace(title)
	}

	// Try title tag
	title = strings.TrimSpace(doc.Find("title").First().Text())
	return title
}

// extractPodcastName extracts the podcast/series name from the page
func (e *MetadataExtractor) extractPodcastName(doc *goquery.Document) string {
	// Try to find podcast name in common locations
	// First try: look for links or headers that might contain podcast name
	podcastName := ""

	// Try meta tag for site name
	if name, exists := doc.Find("meta[property='og:site_name']").Attr("content"); exists {
		podcastName = strings.TrimSpace(name)
		if podcastName != "" {
			return podcastName
		}
	}

	// Try to extract from title (format often: "Episode Title | Podcast Name")
	title := doc.Find("title").First().Text()
	if idx := strings.LastIndex(title, "|"); idx != -1 {
		podcastName = strings.TrimSpace(title[idx+1:])
		if podcastName != "" {
			return podcastName
		}
	}

	return podcastName
}

// isTimePattern checks if text matches common Chinese relative time patterns
func isTimePattern(text string) bool {
	timePatterns := []string{
		"刚刚发布",
		"前", // X小时前, X分钟前, X天前
		"昨天",
		"前天",
		"天前",  // X天前
		"小时前", // X小时前
		"分钟前", // X分钟前
		"周前",  // X周前
		"月前",  // X月前
		"年前",  // X年前
	}

	for _, pattern := range timePatterns {
		if strings.Contains(text, pattern) {
			return true
		}
	}
	return false
}

// containsNumber checks if string contains any digit
func containsNumber(s string) bool {
	for _, c := range s {
		if c >= '0' && c <= '9' {
			return true
		}
	}
	return false
}
