package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/meixg/podcast-reader/internal/config"
	"github.com/meixg/podcast-reader/internal/downloader"
	"github.com/meixg/podcast-reader/internal/validator"
	"github.com/schollz/progressbar/v3"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:    "podcast-downloader",
		Usage:   "从小宇宙FM下载播客音频",
		Version: "1.0.0",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "output",
				Aliases: []string{"o"},
				Usage:   "下载文件保存目录",
				Value:   "./downloads",
			},
			&cli.BoolFlag{
				Name:    "overwrite",
				Aliases: []string{"f"},
				Usage:   "覆盖已存在的文件",
			},
			&cli.BoolFlag{
				Name:  "no-progress",
				Usage: "禁用下载进度条",
			},
			&cli.IntFlag{
				Name:  "retry",
				Usage: "最大重试次数",
				Value: 3,
			},
			&cli.DurationFlag{
				Name:  "timeout",
				Usage: "HTTP请求超时时间",
				Value: 30 * time.Second,
			},
		},
		Action: downloadPodcast,
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

// downloadPodcast is the main action that orchestrates the download process.
func downloadPodcast(ctx *cli.Context) error {
	// 1. Validate URL argument
	if ctx.NArg() < 1 {
		return cli.Exit("请提供小宇宙FM播客URL", 1)
	}

	url := ctx.Args().Get(0)

	// 2. Validate URL format
	urlValidator := validator.NewXiaoyuzhouURLValidator()
	valid, errMsg := urlValidator.ValidateURL(url)
	if !valid {
		return cli.Exit(fmt.Sprintf("URL格式错误: %s", errMsg), 1)
	}

	fmt.Printf("正在获取播客页面: %s\n", url)

	// 3. Create configuration
	cfg := createConfig(ctx)

	// 4. Initialize HTTP client and extractor
	httpClient := downloader.NewHTTPClient(cfg.Timeout)
	extractor := downloader.NewHTMLExtractor(httpClient)

	// 5. Extract episode metadata
	metadata, err := extractor.ExtractURL(context.Background(), url)
	if err != nil {
		return cli.Exit(fmt.Sprintf("获取音频链接失败: %v", err), 1)
	}

	if metadata.Title != "" {
		fmt.Printf("找到播客: %s\n", metadata.Title)
	}

	// 7. Generate file path with podcast title as subdirectory
	// Use sanitized podcast title as folder name
	podcastTitle := sanitizeDirectoryName(metadata.Title)
	podcastDir := filepath.Join(cfg.OutputDirectory, podcastTitle)

	// Create podcast directory if it doesn't exist
	if err := os.MkdirAll(podcastDir, 0755); err != nil {
		return cli.Exit(fmt.Sprintf("创建播客目录失败: %v", err), 1)
	}

	// Use simplified filenames since we already have the podcast title as directory name
	// Files: podcast.m4a, cover.jpg, shownotes.txt
	filename := "podcast.m4a"
	filePath := filepath.Join(podcastDir, filename)

	// 8. Validate file path
	pathValidator := validator.NewDefaultFilePathValidator()
	if err := pathValidator.ValidatePath(filePath, true); err != nil {
		return cli.Exit(fmt.Sprintf("文件路径验证失败: %v", err), 1)
	}

	// 9. Check for existing file
	if _, err := os.Stat(filePath); err == nil {
		if cfg.OverwriteExisting {
			fmt.Printf("文件已存在，将覆盖: %s\n", filePath)
		} else {
			return cli.Exit(fmt.Sprintf("文件已存在: %s\n使用 --overwrite 标志覆盖", filePath), 1)
		}
	}

	fmt.Printf("下载到: %s\n", filePath)

	// 10. Create progress bar
	var progressBar io.Writer
	if cfg.ShowProgress {
		progressBar = progressbar.DefaultBytes(
			-1, // Unknown size initially
			"下载中",
		)
	}

	// 11. Create downloader (use longer timeout for file downloads)
	// File downloads can take much longer than metadata fetching (30s default)
	// Use 1 hour timeout for large audio files
	downloadTimeout := 1 * time.Hour
	downloaderClient := &http.Client{
		Timeout: downloadTimeout,
	}
	fileDownloader := downloader.NewHTTPDownloader(downloaderClient, cfg.ShowProgress)

	// 12. Download file
	var progressWriter io.Writer
	if progressBar != nil {
		progressWriter = progressBar
	}

	fmt.Println() // Add newline before progress bar
	bytesWritten, err := fileDownloader.Download(context.Background(), metadata.AudioURL, filePath, progressWriter)
	if err != nil {
		return cli.Exit(fmt.Sprintf("下载失败: %v", err), 1)
	}

	// 13. Validate downloaded file
	if cfg.ValidateFiles {
		if err := fileDownloader.ValidateFile(filePath); err != nil {
			os.Remove(filePath) // Delete invalid file
			return cli.Exit(fmt.Sprintf("文件验证失败: %v", err), 1)
		}
	}

	// 14. Download cover image (if available)
	if metadata.CoverURL != "" {
		// Use simplified filename: cover.jpg
		coverPath := filepath.Join(podcastDir, "cover.jpg")

		// Create image downloader with separate client (images download quickly)
		imageHTTPClient := &http.Client{
			Timeout: 2 * time.Minute, // 2 minutes for images
		}
		imageDownloader := downloader.NewHTTPImageDownloader(imageHTTPClient, 10*1024*1024) // 10MB max

		// Try to download cover image with graceful degradation
		if _, err := imageDownloader.Download(context.Background(), metadata.CoverURL, coverPath, nil); err != nil {
			logWarning("Warning: Cover image download failed: %v. Audio download completed successfully.", err)
		} else {
			logSuccess("Cover image saved to: %s", coverPath)
		}
	}

	// 15. Save show notes (if available)
	if metadata.ShowNotes != "" {
		// Use simplified filename: shownotes.txt
		showNotesPath := filepath.Join(podcastDir, "shownotes.txt")

		// Create show notes saver
		showNotesSaver := downloader.NewPlainTextShowNotesSaver()

		// Try to save show notes with graceful degradation
		if err := showNotesSaver.Save(metadata.ShowNotes, showNotesPath); err != nil {
			logWarning("Warning: Show notes extraction failed: %v. Audio download completed successfully.", err)
		} else {
			logSuccess("Show notes saved to: %s", showNotesPath)
		}
	}

	// 16. Report success
	fmt.Printf("\n下载成功!\n")
	fmt.Printf("文件位置: %s\n", filePath)
	fmt.Printf("文件大小: %.2f MB\n", float64(bytesWritten)/(1024*1024))

	return nil
}

// createConfig creates the application configuration from CLI flags.
func createConfig(ctx *cli.Context) *config.Config {
	return &config.Config{
		OutputDirectory:   ctx.String("output"),
		OverwriteExisting: ctx.Bool("overwrite"),
		Timeout:           ctx.Duration("timeout"),
		MaxRetries:        ctx.Int("retry"),
		RetryDelay:        1 * time.Second,
		ShowProgress:      !ctx.Bool("no-progress"),
		ValidateFiles:     true,
	}
}

// extractEpisodeID extracts the episode ID from a Xiaoyuzhou FM URL.
func extractEpisodeID(url string) string {
	// URL format: https://www.xiaoyuzhoufm.com/episode/{episode_id}
	// Extract the last part after "/episode/"
	parts := strings.Split(url, "/")
	if len(parts) > 0 {
		return parts[len(parts)-1]
	}
	return "unknown"
}

// logWarning prints a warning message in yellow.
func logWarning(format string, args ...interface{}) {
	yellow := color.New(color.FgYellow).SprintFunc()
	msg := fmt.Sprintf(format, args...)
	fmt.Println(yellow(msg))
}

// logSuccess prints a success message in green.
func logSuccess(format string, args ...interface{}) {
	green := color.New(color.FgGreen).SprintFunc()
	msg := fmt.Sprintf(format, args...)
	fmt.Println(green(msg))
}

// sanitizeDirectoryName sanitizes a podcast title to be used as a directory name.
// Removes characters that are invalid in directory names on most filesystems.
func sanitizeDirectoryName(name string) string {
	// Replace invalid characters with underscore
	// Invalid: < > : " / \ | ? *
	invalidChars := []string{"<", ">", ":", "\"", "/", "\\", "|", "?", "*"}
	result := name
	for _, char := range invalidChars {
		result = strings.ReplaceAll(result, char, "_")
	}

	// Remove leading/trailing spaces and dots
	result = strings.Trim(result, " .")

	// Limit length to avoid filesystem issues (255 chars max for most filesystems)
	if len(result) > 200 {
		result = result[:200]
	}

	// If result is empty, use a default name
	if result == "" {
		result = "Unknown Podcast"
	}

	return result
}
