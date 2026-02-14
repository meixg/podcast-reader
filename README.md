# Podcast Reader

Podcast下载工具和API服务器，支持从小宇宙FM下载播客音频。

Podcast download tool and API server for downloading podcast audio from Xiaoyuzhou FM.

## 功能特性 (Features)

- 🎧 从小宇宙FM下载播客音频
- 📝 自动提取播客标题并生成文件名
- ✨ 支持中文文件名（自动清理特殊字符）
- 📊 实时下载进度显示（速度、大小、时间）
- 🔒 文件验证（M4A格式检查）
- 🔄 自动重试机制（可配置）
- ⚡ 覆盖或跳过已存在的文件
- 🛠️ 完整的错误提示（中文）
- 🖼️ 自动下载封面图片
- 📄 保存节目笔记（show notes）
- 🌐 HTTP API服务器接口
- 📋 任务状态查询和播客列表

## 安装 (Installation)

### 使用 Docker (推荐)

```bash
# 从 GitHub Container Registry 拉取镜像
docker pull ghcr.io/meixg/podcast-reader:latest

# 运行容器
docker run -d \
  -p 8080:8080 \
  -v $(pwd)/downloads:/app/downloads \
  --name podcast-reader \
  ghcr.io/meixg/podcast-reader:latest
```

### 使用 Docker Compose

```bash
# 启动服务
docker-compose up -d

# 查看日志
docker-compose logs -f
```

### 从源码编译 (Build from source)

```bash
# 克隆仓库
git clone https://github.com/meixg/podcast-reader.git
cd podcast-reader

# 下载依赖
go mod download

# 编译 CLI 工具
go build -o podcast-downloader cmd/downloader/main.go

# 编译 API 服务器
go build -o podcast-server cmd/server/main.go

### 构建 Docker 镜像

```bash
# 构建镜像
docker build -t podcast-reader:latest .

# 查看镜像大小
docker images podcast-reader:latest
```
```

### 系统要求 (Requirements)

- Go 1.21 或更高版本（从源码编译）
- Docker 20.10+（使用 Docker 镜像）

## 使用方法 (Usage)

### CLI 工具 (CLI Tool)

#### 基本用法 (Basic Usage)

```bash
./podcast-downloader "https://www.xiaoyuzhoufm.com/episode/69392768281939cce65925d3"
```

#### 命令行选项 (Options)

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

#### 使用示例 (Examples)

```bash
# 指定输出目录
./podcast-downloader -o ~/podcasts "https://www.xiaoyuzhoufm.com/episode/69392768281939cce65925d3"

# 覆盖已存在的文件
./podcast-downloader --overwrite "https://www.xiaoyuzhoufm.com/episode/69392768281939cce65925d3"

# 调整超时和重试次数
./podcast-downloader --timeout 60s --retry 5 "https://www.xiaoyuzhoufm.com/episode/69392768281939cce65925d3"
```

### API 服务器 (API Server)

#### 启动服务器 (Start Server)

```bash
# 使用默认配置（端口8080，下载目录./downloads）
./podcast-server

# 自定义配置
./podcast-server -port 3000 -downloads ~/podcasts -verbose
```

#### 服务器选项 (Server Options)

```
OPTIONS:
   --host value       服务器绑定地址 (default: "0.0.0.0")
   --port value       HTTP服务器端口 (default: 8080)
   --downloads value  下载文件保存目录 (default: "./downloads")
   --verbose          启用详细日志 (default: false)
```

#### API 端点 (API Endpoints)

**1. 提交下载任务 (Submit Download Task)**

```bash
POST /tasks
Content-Type: application/json

{
  "url": "https://www.xiaoyuzhoufm.com/episode/69392768281939cce65925d3"
}
```

响应示例：
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "url": "https://www.xiaoyuzhoufm.com/episode/69392768281939cce65925d3",
  "status": "accepted",
  "created_at": "2026-02-08T10:30:00Z"
}
```

**2. 查询任务状态 (Query Task Status)**

```bash
GET /tasks/{id}
```

响应示例：
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "url": "https://www.xiaoyuzhoufm.com/episode/69392768281939cce65925d3",
  "status": "completed",
  "progress": 100,
  "created_at": "2026-02-08T10:30:00Z",
  "started_at": "2026-02-08T10:30:01Z",
  "completed_at": "2026-02-08T10:32:15Z",
  "podcast": {
    "title": "罗永浩的十字路口",
    "audio_path": "/path/to/downloads/罗永浩的十字路口/podcast.m4a",
    "cover_path": "/path/to/downloads/罗永浩的十字路口/cover.jpg",
    "shownotes_path": "/path/to/downloads/罗永浩的十字路口/shownotes.txt"
  }
}
```

**3. 列出已下载的播客 (List Downloaded Podcasts)**

```bash
GET /podcasts?limit=100&offset=0
```

响应示例：
```json
{
  "podcasts": [
    {
      "url": "https://www.xiaoyuzhoufm.com/episode/69392768281939cce65925d3",
      "title": "罗永浩的十字路口",
      "directory": "/path/to/downloads/罗永浩的十字路口",
      "audio_file": "podcast.m4a",
      "has_cover": true,
      "has_shownotes": true
    }
  ],
  "total": 1,
  "limit": 100,
  "offset": 0
}
```

#### 使用 curl 测试 API (Test API with curl)

```bash
# 提交下载任务
curl -X POST http://localhost:8080/tasks \
  -H "Content-Type: application/json" \
  -d '{"url": "https://www.xiaoyuzhoufm.com/episode/69392768281939cce65925d3"}'

# 查询任务状态
curl http://localhost:8080/tasks/{task_id}

# 列出已下载的播客
curl http://localhost:8080/podcasts
```

## 文件名格式 (Filename Format)

下载的文件组织结构：

```
downloads/
├── Podcast Title/
│   ├── podcast.m4a       # 音频文件
│   ├── cover.jpg         # 封面图片
│   ├── shownotes.txt     # 节目笔记
│   └── .metadata.json    # 元数据（包含原始URL）
```

## 项目结构 (Project Structure)

```
podcast-reader/
├── cmd/
│   ├── podcast-downloader/    # CLI工具入口
│   └── podcast-server/        # API服务器入口
├── internal/
│   ├── config/                # 配置管理
│   ├── downloader/            # 下载器和URL提取器
│   ├── models/                # 数据模型
│   ├── server/                # HTTP服务器和处理器
│   ├── taskmanager/           # 任务管理和目录扫描
│   └── validator/             # URL和文件路径验证
├── pkg/
│   └── httpclient/            # HTTP客户端（带重试）
├── specs/                     # 规格文档
├── go.mod
├── go.sum
└── README.md
```

## API 文档 (API Documentation)

详细的 OpenAPI 规范文档请参阅：[specs/3-podcast-api-server/contracts/openapi.yaml](specs/3-podcast-api-server/contracts/openapi.yaml)

## 依赖项 (Dependencies)

- [goquery](https://github.com/PuerkitoBio/goquery) - HTML解析
- [urfave/cli](https://github.com/urfave/cli) - CLI框架
- [progressbar/v3](https://github.com/schollz/progressbar) - 进度条显示
- [google/uuid](https://github.com/google/uuid) - UUID生成

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
