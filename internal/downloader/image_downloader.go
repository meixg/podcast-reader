package downloader

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/schollz/progressbar/v3"
)

// ImageDownloader defines the interface for downloading cover images.
type ImageDownloader interface {
	// Download downloads an image from the given URL and saves it to the specified path.
	//
	// Parameters:
	//   ctx - Context for cancellation and timeout
	//   imageURL - The URL of the image to download
	//   destPath - The destination file path
	//   progress - Optional progress writer (can be nil)
	//
	// Returns:
	//   int64 - Number of bytes written
	//   error - Error if download or validation fails
	Download(ctx context.Context, imageURL string, destPath string, progress io.Writer) (int64, error)

	// ValidateImage validates that a file is a valid image format.
	//
	// Parameters:
	//   filePath - Path to the file to validate
	//
	// Returns:
	//   error - Error if file is not a valid image
	ValidateImage(filePath string) error
}

// HTTPImageDownloader implements ImageDownloader using HTTP client.
type HTTPImageDownloader struct {
	client       *http.Client
	maxFileSize  int64
	retryCount   int
	retryDelay   time.Duration
	showProgress bool
}

// NewHTTPImageDownloader creates a new HTTP image downloader.
func NewHTTPImageDownloader(client *http.Client, maxFileSize int64) *HTTPImageDownloader {
	return &HTTPImageDownloader{
		client:       client,
		maxFileSize:  maxFileSize,
		retryCount:   3,
		retryDelay:   1 * time.Second,
		showProgress: true,
	}
}

// Download downloads an image with retry logic and progress tracking.
func (d *HTTPImageDownloader) Download(ctx context.Context, imageURL string, destPath string, progress io.Writer) (int64, error) {
	var lastErr error
	var bytesWritten int64

	// Retry logic
	for attempt := 0; attempt <= d.retryCount; attempt++ {
		if attempt > 0 {
			time.Sleep(d.retryDelay * time.Duration(attempt))
		}

		// Create request
		req, err := http.NewRequestWithContext(ctx, "GET", imageURL, nil)
		if err != nil {
			lastErr = fmt.Errorf("创建请求失败: %w", err)
			continue
		}

		// Execute request
		resp, err := d.client.Do(req)
		if err != nil {
			lastErr = fmt.Errorf("下载失败: %w", err)
			continue
		}

		// Check response status
		if resp.StatusCode != http.StatusOK {
			resp.Body.Close()
			lastErr = fmt.Errorf("HTTP状态码: %d", resp.StatusCode)
			continue
		}

		// Check content length
		contentLength := resp.ContentLength
		if d.maxFileSize > 0 && contentLength > d.maxFileSize {
			resp.Body.Close()
			return 0, fmt.Errorf("%w: 图片大小 %d 字节超过限制 %d 字节", ErrImageTooLarge, contentLength, d.maxFileSize)
		}

		// Create destination directory
		destDir := filepath.Dir(destPath)
		if err := os.MkdirAll(destDir, 0755); err != nil {
			resp.Body.Close()
			return 0, fmt.Errorf("创建目录失败: %w", err)
		}

		// Create destination file
		destFile, err := os.Create(destPath)
		if err != nil {
			resp.Body.Close()
			return 0, fmt.Errorf("创建文件失败: %w", err)
		}

		// Setup progress tracking
		var writer io.Writer = destFile
		if progress != nil && d.showProgress && contentLength > 1024*1024 { // Only show progress for images > 1MB
			bar := progressbar.DefaultBytes(
				contentLength,
				"下载封面中",
			)
			writer = io.MultiWriter(destFile, bar)
		}

		// Copy data
		bytesWritten, err = io.Copy(writer, resp.Body)
		resp.Body.Close()
		destFile.Close()

		if err != nil {
			os.Remove(destPath) // Clean up partial file
			lastErr = fmt.Errorf("写入文件失败: %w", err)
			continue
		}

		// Validate downloaded file
		if err := d.ValidateImage(destPath); err != nil {
			os.Remove(destPath) // Remove invalid file
			return 0, fmt.Errorf("%w: %v", ErrInvalidImage, err)
		}

		// Success
		return bytesWritten, nil
	}

	// All retries failed
	return 0, lastErr
}

// ValidateImage validates that a file is a valid image using magic byte detection.
func (d *HTTPImageDownloader) ValidateImage(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("打开文件失败: %w", err)
	}
	defer file.Close()

	// Read first 12 bytes for magic byte detection
	header := make([]byte, 12)
	_, err = file.Read(header)
	if err != nil {
		return fmt.Errorf("读取文件头失败: %w", err)
	}

	// Detect format
	format := d.detectFormat(header)
	if format == "" {
		return fmt.Errorf("无法识别的图片格式")
	}

	return nil
}

// detectFormat identifies image format from binary header using magic bytes.
func (d *HTTPImageDownloader) detectFormat(header []byte) string {
	// JPEG: FF D8 FF
	if len(header) >= 3 && header[0] == 0xFF && header[1] == 0xD8 && header[2] == 0xFF {
		return "jpeg"
	}

	// PNG: 89 50 4E 47
	if len(header) >= 4 &&
		header[0] == 0x89 && header[1] == 0x50 && header[2] == 0x4E && header[3] == 0x47 {
		return "png"
	}

	// WebP: RIFF....WEBP
	if len(header) >= 12 &&
		header[0] == 0x52 && header[1] == 0x49 && header[2] == 0x46 && header[3] == 0x46 && // "RIFF"
		header[8] == 0x57 && header[9] == 0x45 && header[10] == 0x42 && header[11] == 0x50 { // "WEBP"
		return "webp"
	}

	// GIF: 47 49 46 38
	if len(header) >= 4 &&
		header[0] == 0x47 && header[1] == 0x49 && header[2] == 0x46 && header[3] == 0x38 {
		return "gif"
	}

	return ""
}
