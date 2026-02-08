package downloader

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
)

// Define file download error types
var (
	ErrNetworkTimeout    = fmt.Errorf("网络超时")
	ErrConnectionRefused = fmt.Errorf("连接被拒绝")
	ErrDiskFull          = fmt.Errorf("磁盘空间不足")
	ErrPermissionDenied  = fmt.Errorf("权限被拒绝")
	ErrInvalidAudio      = fmt.Errorf("音频文件无效")
)

// FileDownloader defines the interface for downloading files with progress tracking.
type FileDownloader interface {
	// Download fetches the audio file and writes it to the local filesystem.
	Download(ctx context.Context, audioURL, filePath string, progress io.Writer) (bytesWritten int64, err error)

	// ValidateFile checks if the downloaded file is a valid audio file.
	ValidateFile(filePath string) error
}

// HTTPDownloader implements FileDownloader with HTTP client and progress tracking.
type HTTPDownloader struct {
	// client is the HTTP client to use for downloads
	client *http.Client
	// showProgress controls whether progress should be displayed
	showProgress bool
}

// NewHTTPDownloader creates a new HTTP downloader.
func NewHTTPDownloader(client *http.Client, showProgress bool) *HTTPDownloader {
	return &HTTPDownloader{
		client:       client,
		showProgress: showProgress,
	}
}

// Download fetches the audio file and writes it to the local filesystem.
func (d *HTTPDownloader) Download(ctx context.Context, audioURL, filePath string, progress io.Writer) (int64, error) {
	// Create HTTP request
	req, err := http.NewRequestWithContext(ctx, "GET", audioURL, nil)
	if err != nil {
		return 0, fmt.Errorf("%w: %v", ErrConnectionRefused, err)
	}

	// Execute request
	resp, err := d.client.Do(req)
	if err != nil {
		return 0, fmt.Errorf("%w: %v", ErrNetworkTimeout, err)
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("下载失败: HTTP %d", resp.StatusCode)
	}

	// Create output file
	out, err := os.Create(filePath)
	if err != nil {
		return 0, fmt.Errorf("%w: %v", ErrPermissionDenied, err)
	}
	defer out.Close()

	// Track bytes written
	var bytesWritten int64

	// Use progress writer if provided
	var writer io.Writer = out
	if progress != nil {
		writer = io.MultiWriter(out, progress)
	}

	// Copy with progress tracking
	bytesWritten, err = io.Copy(writer, resp.Body)
	if err != nil {
		return 0, fmt.Errorf("下载中断: %w", err)
	}

	return bytesWritten, nil
}

// ValidateFile checks if the downloaded file is a valid audio file.
func (d *HTTPDownloader) ValidateFile(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("无法打开文件: %w", err)
	}
	defer file.Close()

	// Read first 12 bytes for magic byte check
	header := make([]byte, 12)
	_, err = file.Read(header)
	if err != nil {
		return fmt.Errorf("无法读取文件: %w", err)
	}

	// Check for M4A/MP4 magic bytes (ftyp)
	// M4A files start with: 00 00 00 xx 66 74 79 70
	// Where "ftyp" is at offset 4
	if len(header) < 8 {
		return ErrInvalidAudio
	}

	if string(header[4:8]) != "ftyp" {
		return fmt.Errorf("%w: 下载的文件不是有效的M4A音频", ErrInvalidAudio)
	}

	return nil
}
