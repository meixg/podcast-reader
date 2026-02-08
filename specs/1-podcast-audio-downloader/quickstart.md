# Quickstart Guide: Podcast Audio Downloader

**Feature**: Podcast Audio Downloader (CLI tool)
**Date**: 2026-02-08
**Purpose**: Get the podcast downloader running from scratch in under 5 minutes

---

## Prerequisites

- **Go 1.21+** installed
- **Internet connection** for accessing Xiaoyuzhou FM
- **200MB+ free disk space** for downloaded podcast files
- **Write permissions** in the download directory

---

## Installation

### Step 1: Clone Repository (if applicable)

```bash
git clone https://github.com/your-org/podcast-reader.git
cd podcast-reader
```

### Step 2: Build the CLI Tool

```bash
cd cmd/podcast-downloader
go build -o podcast-downloader
```

This creates a standalone binary `podcast-downloader` in the current directory.

### Step 3: (Optional) Install Globally

```bash
# Linux/macOS
sudo mv podcast-downloader /usr/local/bin/

# Or add to PATH
export PATH=$PATH:$(pwd)
```

### Alternative: Install Directly from Source

```bash
go install github.com/your-org/podcast-reader/cmd/podcast-downloader@latest
```

---

## Basic Usage

### Download a Podcast Episode

```bash
podcast-downloader "https://www.xiaoyuzhoufm.com/episode/69392768281939cce65925d3"
```

**What happens**:
1. Tool validates the URL
2. Fetches the episode page
3. Extracts the audio file URL
4. Downloads the audio to `./downloads/` directory
5. Shows progress bar during download
6. Validates the downloaded file
7. Reports success or failure

**Expected output**:
```
Fetching episode page...
Found: "技术电台第42期: Go语言实战"
Downloading to: ./downloads/技术电台第42期_ Go语言实战_69392768281939cce65925d3.m4a

downloading                                      45.2 MiB / 50.0 MiB ( 90%) [=========>     ] 1.9 MiB/s 2s

Download complete: ./downloads/技术电台第42期_ Go语言实战_69392768281939cce65925d3.m4a
File size: 52,428,800 bytes
Duration: 28 seconds
```

---

## Advanced Usage

### Specify Output Directory

```bash
podcast-downloader --output ~/music/podcasts "https://www.xiaoyuzhoufm.com/episode/..."
```

### Overwrite Existing Files

```bash
podcast-downloader --overwrite "https://www.xiaoyuzhoufm.com/episode/..."
```

### Disable Progress Bar

```bash
podcast-downloader --no-progress "https://www.xiaoyuzhoufm.com/episode/..."
```

### Custom Retry Count

```bash
podcast-downloader --retry 5 "https://www.xiaoyuzhoufm.com/episode/..."
```

### Custom Timeout

```bash
podcast-downloader --timeout 60s "https://www.xiaoyuzhoufm.com/episode/..."
```

---

## Command-Line Options

| Flag | Short | Description | Default |
|------|-------|-------------|---------|
| `--output` | `-o` | Output directory for downloads | `./downloads` |
| `--overwrite` | `-f` | Overwrite existing files | `false` |
| `--no-progress` | | Disable progress bar | `false` |
| `--retry` | | Maximum retry attempts | `3` |
| `--timeout` | | HTTP request timeout | `30s` |
| `--help` | `-h` | Show help message | - |
| `--version` | `-v` | Show version | - |

---

## Common Workflows

### Download Multiple Episodes

Create a bash script to download multiple episodes:

```bash
#!/bin/bash
# download_episodes.sh

urls=(
    "https://www.xiaoyuzhoufm.com/episode/69392768281939cce65925d3"
    "https://www.xiaoyuzhoufm.com/episode/1234567890abcdef"
    "https://www.xiaoyuzhoufm.com/episode/abcdef1234567890"
)

for url in "${urls[@]}"; do
    podcast-downloader "$url"
done
```

### Download from File

Create a text file with URLs (one per line):

```bash
# episodes.txt
https://www.xiaoyuzhoufm.com/episode/69392768281939cce65925d3
https://www.xiaoyuzhoufm.com/episode/1234567890abcdef
```

Download all episodes:

```bash
cat episodes.txt | xargs -I {} podcast-downloader {}
```

---

## Troubleshooting

### Error: "Invalid URL format"

**Cause**: URL does not match Xiaoyuzhou FM pattern

**Solution**: Verify URL format is correct:
```
✅ https://www.xiaoyuzhoufm.com/episode/69392768281939cce65925d3
❌ https://xiaoyuzhoufm.com/podcast/12345 (wrong path)
❌ www.xiaoyuzhoufm.com/episode/... (missing protocol)
```

### Error: "Episode not found (404)"

**Cause**: Episode does not exist or was deleted

**Solution**: Verify the episode exists by opening the URL in a browser

### Error: "No audio link found"

**Cause**: Website structure changed or episode has no audio

**Solution**: Open the episode page in browser and check if audio player exists. If issue persists, report as bug.

### Error: "Permission denied"

**Cause**: Cannot write to output directory

**Solution**:
```bash
# Create directory with proper permissions
mkdir -p downloads
chmod 755 downloads

# Or specify a different output directory
podcast-downloader --output ~/downloads "..."
```

### Error: "Download failed: timeout"

**Cause**: Network is slow or server is not responding

**Solution**: Increase timeout:
```bash
podcast-downloader --timeout 60s "..."
```

### Error: "Downloaded file is not valid audio"

**Cause**: Server returned error page instead of audio file

**Solution**: Check if audio URL is accessible by testing in browser. May need to update URL extraction logic.

---

## Development Setup

### Prerequisites for Development

```bash
# Install dependencies
go mod tidy

# Run tests
go test ./...

# Run with race detector
go run -race main.go
```

### Project Structure

```
cmd/podcast-downloader/
├── main.go           # Application entry point
└── root.go           # CLI command definitions

internal/
├── downloader/       # Download logic
│   ├── downloader.go
│   ├── url_extractor.go
│   └── progress.go
├── models/           # Data structures
│   ├── episode.go
│   └── download_session.go
├── validator/        # Input validation
│   └── url_validator.go
└── config/           # Configuration
    └── config.go

pkg/
└── httpclient/       # Reusable HTTP client
    └── client.go
```

### Running Tests

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests for specific package
go test ./internal/downloader

# Run tests with verbose output
go test -v ./internal/validator
```

---

## Performance Tips

### Download Speed

If downloads are slow:
1. Check your internet connection
2. Try increasing timeout: `--timeout 60s`
3. Verify audio URL is accessible (some URLs may be rate-limited)

### Disk Space

Check available disk space before downloading large batches:
```bash
df -h .
```

### Concurrent Downloads

For concurrent downloads, run multiple instances:
```bash
podcast-downloader "url1" &
podcast-downloader "url2" &
wait
```

---

## FAQ

### Q: Can I download from other podcast websites?

**A**: No, this tool is specifically designed for Xiaoyuzhou FM. Support for other websites would require updating the URL extraction logic.

### Q: What happens if the download is interrupted?

**A**: Partial downloads are not automatically resumed. You need to restart the download. Use `--overwrite` flag to replace the partial file.

### Q: Can I download private/premium episodes?

**A**: No, the tool assumes audio files are publicly accessible (Assumption #2). Premium content requiring authentication is not supported.

### Q: How do I report bugs or request features?

**A**: Please open an issue on the GitHub repository with:
- Command used
- Error message (if any)
- Episode URL (if applicable)
- Expected vs actual behavior

---

## Next Steps

1. **Try it out**: Download your first podcast episode
2. **Customize**: Adjust output directory and flags to your preference
3. **Batch download**: Create scripts to download multiple episodes
4. **Contribute**: Check out the source code and submit pull requests

---

## Support

For issues, questions, or contributions:
- **GitHub Issues**: [github.com/your-org/podcast-reader/issues]
- **Documentation**: [docs/]
- **Example Usage**: See `examples/` directory

---

**Enjoy your podcasts! 🎧**
