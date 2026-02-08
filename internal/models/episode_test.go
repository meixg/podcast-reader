package models

import (
	"strings"
	"testing"
)

func TestEpisode_SanitizedTitle_ValidTitle(t *testing.T) {
	episode := &Episode{
		ID:    "12345",
		Title: "Test Episode Title",
	}

	sanitized := episode.SanitizedTitle()

	if sanitized != "Test Episode Title" {
		t.Errorf("SanitizedTitle() = %q, want %q", sanitized, "Test Episode Title")
	}
}

func TestEpisode_SanitizedTitle_InvalidCharacters(t *testing.T) {
	tests := []struct {
		name     string
		title    string
		expected string
	}{
		{
			name:     "forward slash",
			title:    "Episode/1/2024",
			expected: "Episode_1_2024",
		},
		{
			name:     "backslash",
			title:    "Episode\\1\\2024",
			expected: "Episode_1_2024",
		},
		{
			name:     "colon",
			title:    "Episode: 2024-01-01",
			expected: "Episode_ 2024-01-01",
		},
		{
			name:     "asterisk",
			title:    "Episode*Important*",
			expected: "Episode_Important_",
		},
		{
			name:     "question mark",
			title:    "What is this?",
			expected: "What is this_",
		},
		{
			name:     "quotes",
			title:    `Episode "Test"`,
			expected: "Episode _Test_",
		},
		{
			name:     "less than and greater than",
			title:    "Episode <Test>",
			expected: "Episode _Test_",
		},
		{
			name:     "pipe",
			title:    "Episode | Test",
			expected: "Episode _ Test",
		},
		{
			name:     "mixed invalid characters",
			title:    `Episode: "Test/2024" <New> | Review?`,
			expected: "Episode_ _Test_2024_ _New_ _ Review_",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			episode := &Episode{
				ID:    "12345",
				Title: tc.title,
			}

			sanitized := episode.SanitizedTitle()

			if sanitized != tc.expected {
				t.Errorf("SanitizedTitle() = %q, want %q", sanitized, tc.expected)
			}
		})
	}
}

func TestEpisode_SanitizedTitle_TooLong(t *testing.T) {
	longTitle := strings.Repeat("A", 250)
	episode := &Episode{
		ID:    "12345",
		Title: longTitle,
	}

	sanitized := episode.SanitizedTitle()

	if len(sanitized) != 200 {
		t.Errorf("SanitizedTitle() length = %d, want 200", len(sanitized))
	}

	// Should be all As
	if strings.Trim(sanitized, "A") != "" {
		t.Error("SanitizedTitle() should contain only As")
	}
}

func TestEpisode_SanitizedTitle_EmptyTitle(t *testing.T) {
	episode := &Episode{
		ID:    "12345",
		Title: "",
	}

	sanitized := episode.SanitizedTitle()

	if sanitized != "12345" {
		t.Errorf("SanitizedTitle() = %q, want %q", sanitized, "12345")
	}
}

func TestEpisode_SanitizedTitle_WhitespaceOnly(t *testing.T) {
	episode := &Episode{
		ID:    "12345",
		Title: "   ",
	}

	sanitized := episode.SanitizedTitle()

	if sanitized != "12345" {
		t.Errorf("SanitizedTitle() = %q, want %q", sanitized, "12345")
	}
}

func TestEpisode_SanitizedTitle_TrimsWhitespace(t *testing.T) {
	episode := &Episode{
		ID:    "12345",
		Title: "  Test Episode  ",
	}

	sanitized := episode.SanitizedTitle()

	if sanitized != "Test Episode" {
		t.Errorf("SanitizedTitle() = %q, want %q", sanitized, "Test Episode")
	}
}

func TestEpisode_GenerateFilename(t *testing.T) {
	tests := []struct {
		name     string
		episode  *Episode
		expected string
	}{
		{
			name: "normal title",
			episode: &Episode{
				ID:    "12345",
				Title: "Test Episode",
			},
			expected: "Test Episode_12345.m4a",
		},
		{
			name: "title with invalid characters",
			episode: &Episode{
				ID:    "12345",
				Title: "Episode/1",
			},
			expected: "Episode_1_12345.m4a",
		},
		{
			name: "empty title",
			episode: &Episode{
				ID:    "12345",
				Title: "",
			},
			expected: "12345_12345.m4a",
		},
		{
			name: "title with spaces",
			episode: &Episode{
				ID:    "12345",
				Title: "  Test Episode  ",
			},
			expected: "Test Episode_12345.m4a",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			filename := tc.episode.GenerateFilename()

			if filename != tc.expected {
				t.Errorf("GenerateFilename() = %q, want %q", filename, tc.expected)
			}
		})
	}
}

func TestEpisode_Validate_ValidEpisode(t *testing.T) {
	episode := &Episode{
		ID:       "12345",
		Title:    "Test Episode",
		AudioURL: "https://example.com/audio.m4a",
		PageURL:  "https://example.com/episode/12345",
	}

	err := episode.Validate()

	if err != nil {
		t.Errorf("Validate() error = %v", err)
	}
}

func TestEpisode_Validate_MissingID(t *testing.T) {
	episode := &Episode{
		Title:    "Test Episode",
		AudioURL: "https://example.com/audio.m4a",
		PageURL:  "https://example.com/episode/12345",
	}

	err := episode.Validate()

	if err == nil {
		t.Error("Validate() should return error for missing ID")
	}

	expected := "episode ID is required"
	if err.Error() != expected {
		t.Errorf("Error message = %q, want %q", err.Error(), expected)
	}
}

func TestEpisode_Validate_MissingPageURL(t *testing.T) {
	episode := &Episode{
		ID:       "12345",
		Title:    "Test Episode",
		AudioURL: "https://example.com/audio.m4a",
	}

	err := episode.Validate()

	if err == nil {
		t.Error("Validate() should return error for missing PageURL")
	}

	expected := "page URL is required"
	if err.Error() != expected {
		t.Errorf("Error message = %q, want %q", err.Error(), expected)
	}
}

func TestEpisode_Validate_MissingAudioURL(t *testing.T) {
	episode := &Episode{
		ID:      "12345",
		Title:   "Test Episode",
		PageURL: "https://example.com/episode/12345",
	}

	err := episode.Validate()

	if err == nil {
		t.Error("Validate() should return error for missing AudioURL")
	}

	expected := "audio URL is required"
	if err.Error() != expected {
		t.Errorf("Error message = %q, want %q", err.Error(), expected)
	}
}

func TestEpisode_Validate_InvalidAudioFormat(t *testing.T) {
	tests := []struct {
		name     string
		audioURL string
		expected string
	}{
		{
			name:     "mp3 format",
			audioURL: "https://example.com/audio.mp3",
			expected: "audio URL must be .m4a format",
		},
		{
			name:     "wav format",
			audioURL: "https://example.com/audio.wav",
			expected: "audio URL must be .m4a format",
		},
		{
			name:     "no extension",
			audioURL: "https://example.com/audio",
			expected: "audio URL must be .m4a format",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			episode := &Episode{
				ID:       "12345",
				Title:    "Test Episode",
				AudioURL: tc.audioURL,
				PageURL:  "https://example.com/episode/12345",
			}

			err := episode.Validate()

			if err == nil {
				t.Error("Validate() should return error for invalid audio format")
			}

			if err.Error() != tc.expected {
				t.Errorf("Error message = %q, want %q", err.Error(), tc.expected)
			}
		})
	}
}

func TestEpisode_Validate_AllFieldsMissing(t *testing.T) {
	episode := &Episode{}

	err := episode.Validate()

	if err == nil {
		t.Error("Validate() should return error for empty episode")
	}

	// Should error on ID first
	expected := "episode ID is required"
	if err.Error() != expected {
		t.Errorf("Error message = %q, want %q", err.Error(), expected)
	}
}

func TestEpisode_Validate_ValidM4aFormat(t *testing.T) {
	validURLs := []string{
		"https://example.com/audio.m4a",
		"https://cdn.example.com/podcasts/episode.m4a",
		"https://example.com/audio.m4a?param=value",
		"http://example.com/audio.m4a",
	}

	for _, audioURL := range validURLs {
		t.Run(audioURL, func(t *testing.T) {
			episode := &Episode{
				ID:       "12345",
				Title:    "Test Episode",
				AudioURL: audioURL,
				PageURL:  "https://example.com/episode/12345",
			}

			err := episode.Validate()

			if err != nil {
				t.Errorf("Validate() should not return error for valid .m4a URL: %v", err)
			}
		})
	}
}
