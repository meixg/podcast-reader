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

	"github.com/meixg/podcast-reader/internal/config"
	"github.com/meixg/podcast-reader/internal/downloader"
	"github.com/meixg/podcast-reader/internal/models"
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

	// 5. Extract episode metadata (title, audio URL)
	audioURL, title, err := extractor.ExtractURL(context.Background(), url)
	if err != nil {
		return cli.Exit(fmt.Sprintf("获取音频链接失败: %v", err), 1)
	}

	if title != "" {
		fmt.Printf("找到播客: %s\n", title)
	}

	// 6. Create episode model
	episode := &models.Episode{
		ID:       extractEpisodeID(url),
		Title:    title,
		AudioURL: audioURL,
		PageURL:  url,
	}

	// 7. Generate file path
	filename := episode.GenerateFilename()
	filePath := filepath.Join(cfg.OutputDirectory, filename)

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

	// 11. Create downloader
	downloaderClient := &http.Client{
		Timeout: cfg.Timeout,
	}
	fileDownloader := downloader.NewHTTPDownloader(downloaderClient, cfg.ShowProgress)

	// 12. Download file
	var progressWriter io.Writer
	if progressBar != nil {
		progressWriter = progressBar
	}

	fmt.Println() // Add newline before progress bar
	bytesWritten, err := fileDownloader.Download(context.Background(), audioURL, filePath, progressWriter)
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

	// 14. Report success
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
