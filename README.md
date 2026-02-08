# Podcast Downloader

一个简单的命令行工具，用于从小宇宙FM下载播客音频文件。

A simple CLI tool to download podcast audio files from Xiaoyuzhou FM.

## 功能特性 (Features)

- 🎧 从小宇宙FM下载播客音频
- 📝 自动提取播客标题并生成文件名
- ✨ 支持中文文件名（自动清理特殊字符）
- 📊 实时下载进度显示（速度、大小、时间）
- 🔒 文件验证（M4A格式检查）
- 🔄 自动重试机制（可配置）
- ⚡ 覆盖或跳过已存在的文件
- 🛠️ 完整的错误提示（中文）

## 安装 (Installation)

### 从源码编译 (Build from source)

```bash
# 克隆仓库
git clone https://github.com/meixg/podcast-reader.git
cd podcast-reader

# 下载依赖
go mod download

# 编译
go build -o podcast-downloader cmd/podcast-downloader/main.go
```

### 系统要求 (Requirements)

- Go 1.21 或更高版本

## 使用方法 (Usage)

### 基本用法 (Basic Usage)

```bash
./podcast-downloader "https://www.xiaoyuzhoufm.com/episode/69392768281939cce65925d3"
```

### 命令行选项 (Options)

```
OPTIONS:
   --output value, -o value  下载文件保存目录 (default: "./downloads")
   --overwrite, -f           覆盖已存在的文件 (default: false)
   --no-progress             禁用下载进度条 (default: false)
   --retry value             最大重试次数 (default: 3)
   --timeout value           HTTP请求超时时间 (default: 30s)
   --help, -h                显示帮助信息
   --version, -v             显示版本号
```

### 使用示例 (Examples)

#### 指定输出目录 (Specify output directory)

```bash
./podcast-downloader -o ~/podcasts "https://www.xiaoyuzhoufm.com/episode/69392768281939cce65925d3"
```

#### 覆盖已存在的文件 (Overwrite existing files)

```bash
./podcast-downloader --overwrite "https://www.xiaoyuzhoufm.com/episode/69392768281939cce65925d3"
```

#### 禁用进度条 (Disable progress bar)

```bash
./podcast-downloader --no-progress "https://www.xiaoyuzhoufm.com/episode/69392768281939cce65925d3"
```

#### 调整超时和重试次数 (Adjust timeout and retries)

```bash
./podcast-downloader --timeout 60s --retry 5 "https://www.xiaoyuzhoufm.com/episode/69392768281939cce65925d3"
```

### 批量下载 (Batch Download)

你可以使用shell脚本批量下载多集播客：

```bash
#!/bin/bash
# batch_download.sh

URLS=(
  "https://www.xiaoyuzhoufm.com/episode/69392768281939cce65925d3"
  "https://www.xiaoyuzhoufm.com/episode/another-episode-id"
  "https://www.xiaoyuzhoufm.com/episode/yet-another-episode-id"
)

for url in "${URLS[@]}"; do
  ./podcast-downloader "$url"
done
```

## 文件名格式 (Filename Format)

下载的文件使用以下命名格式：

```
{清理后的标题}_{集数ID}.m4a
```

例如：
```
罗永浩的十字路口_Episode01_69392768281939cce65925d3.m4a
```

标题中的特殊字符（`< > : " / \ | ? *`）会被自动替换为下划线。

## 项目结构 (Project Structure)

```
podcast-reader/
├── cmd/
│   ├── podcast-downloader/    # 主程序入口
│   └── inspect/               # HTML检查工具（调试用）
├── internal/
│   ├── config/                # 配置管理
│   ├── downloader/            # 下载器和URL提取器
│   ├── models/                # 数据模型
│   └── validator/             # URL和文件路径验证
├── pkg/
│   └── httpclient/            # HTTP客户端（带重试）
├── specs/                     # 规格文档
├── go.mod
├── go.sum
└── README.md
```

## 依赖项 (Dependencies)

- [goquery](https://github.com/PuerkitoBio/goquery) - HTML解析
- [urfave/cli](https://github.com/urfave/cli) - CLI框架
- [progressbar/v3](https://github.com/schollz/progressbar) - 进度条显示

## 开发 (Development)

### 运行测试 (Run tests)

```bash
go test ./...
```

### 代码格式化 (Format code)

```bash
gofmt -w .
```

### 代码检查 (Lint code)

```bash
go vet ./...
```

## 注意事项 (Notes)

- 本工具仅供个人学习使用
- 请遵守小宇宙FM的服务条款
- 下载的音频文件仅供个人使用，不得用于商业目的
- 建议使用官方客户端支持播客创作者

## 许可证 (License)

MIT License

## 贡献 (Contributing)

欢迎提交 Issue 和 Pull Request！

## 相关链接 (Links)

- [小宇宙FM](https://www.xiaoyuzhoufm.com)
- [Go语言官方网站](https://golang.org)

## 致谢 (Acknowledgments)

感谢小宇宙FM提供的优质播客平台。
