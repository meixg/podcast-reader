# Podcast API Server 测试指南

## 启动服务器

```bash
# 基本启动（使用默认配置）
./build/podcast-server

# 自定义配置
./build/podcast-server -port 8080 -downloads ./downloads -verbose
```

## 自动化测试

使用提供的测试脚本：

```bash
./test-api.sh
```

这会自动测试所有API端点和各种场景。

## 手动测试各个端点

### 1. 提交下载任务 (POST /tasks)

```bash
curl -X POST http://localhost:8080/tasks \
  -H "Content-Type: application/json" \
  -d '{"url": "https://www.xiaoyuzhoufm.com/episode/69806618073030367acec13b"}'
```

**预期响应（202 Accepted）：**
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "url": "https://www.xiaoyuzhoufm.com/episode/69806618073030367acec13b",
  "status": "pending",
  "created_at": "2026-02-08T10:30:00Z"
}
```

保存返回的 `id`，后续查询会用到。

### 2. 查询任务状态 (GET /tasks/{id})

```bash
# 将 {id} 替换为实际的任务ID
curl http://localhost:8080/tasks/550e8400-e29b-41d4-a716-446655440000
```

**响应示例（进行中）：**
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "url": "https://www.xiaoyuzhoufm.com/episode/69806618073030367acec13b",
  "status": "in_progress",
  "progress": 45,
  "created_at": "2026-02-08T10:30:00Z",
  "started_at": "2026-02-08T10:30:01Z"
}
```

**响应示例（已完成）：**
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "url": "https://www.xiaoyuzhoufm.com/episode/69806618073030367acec13b",
  "status": "completed",
  "progress": 100,
  "created_at": "2026-02-08T10:30:00Z",
  "started_at": "2026-02-08T10:30:01Z",
  "completed_at": "2026-02-08T10:32:15Z",
  "podcast": {
    "title": "播客标题",
    "source_url": "https://www.xiaoyuzhoufm.com/episode/69806618073030367acec13b",
    "audio_path": "/path/to/downloads/播客标题/podcast.m4a",
    "cover_path": "/path/to/downloads/播客标题/cover.jpg",
    "shownotes_path": "/path/to/downloads/播客标题/shownotes.txt",
    "downloaded_at": "2026-02-08T10:32:15Z"
  }
}
```

**响应示例（失败）：**
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "url": "https://www.xiaoyuzhoufm.com/episode/69806618073030367acec13b",
  "status": "failed",
  "error": "下载失败: 网络超时"
}
```

### 3. 列出已下载的播客 (GET /podcasts)

```bash
# 基本查询（默认：limit=100, offset=0）
curl http://localhost:8080/podcasts

# 自定义分页
curl "http://localhost:8080/podcasts?limit=10&offset=0"

# 只查询前5条
curl "http://localhost:8080/podcasts?limit=5"
```

**响应示例：**
```json
{
  "podcasts": [
    {
      "url": "https://www.xiaoyuzhoufm.com/episode/69806618073030367acec13b",
      "title": "播客标题",
      "directory": "/path/to/downloads/播客标题",
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

## 错误场景测试

### 1. 重复提交（应返回409 Conflict）

```bash
# 第一次提交
curl -X POST http://localhost:8080/tasks \
  -H "Content-Type: application/json" \
  -d '{"url": "https://www.xiaoyuzhoufm.com/episode/69806618073030367acec13b"}'

# 立即再次提交相同URL
curl -X POST http://localhost:8080/tasks \
  -H "Content-Type: application/json" \
  -d '{"url": "https://www.xiaoyuzhoufm.com/episode/69806618073030367acec13b"}'
```

**预期响应（409 Conflict）：**
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "url": "https://www.xiaoyuzhoufm.com/episode/69806618073030367acec13b",
  "status": "in_progress"
}
```

### 2. 无效URL（应返回400 Bad Request）

```bash
curl -X POST http://localhost:8080/tasks \
  -H "Content-Type: application/json" \
  -d '{"url": "https://example.com/not-valid"}'
```

**预期响应（400 Bad Request）：**
```json
{
  "error": {
    "code": "INVALID_URL",
    "message": "The provided URL is not valid",
    "details": "invalid URL: URL格式不正确..."
  }
}
```

### 3. 任务不存在（应返回404 Not Found）

```bash
curl http://localhost:8080/tasks/00000000-0000-0000-0000-000000000000
```

**预期响应（404 Not Found）：**
```json
{
  "error": {
    "code": "TASK_NOT_FOUND",
    "message": "Task not found"
  }
}
```

### 4. 无效的分页参数（应返回400 Bad Request）

```bash
# limit > 1000
curl "http://localhost:8080/podcasts?limit=2000"

# offset < 0
curl "http://localhost:8080/podcasts?offset=-1"
```

## 监控下载进度

### 实时查询任务状态

```bash
# 保存任务ID
TASK_ID="550e8400-e29b-41d4-a716-446655440000"

# 循环查询状态
watch -n 2 "curl -s http://localhost:8080/tasks/$TASK_ID | jq '.status, .progress'"
```

### 使用脚本监控

```bash
#!/bin/bash
TASK_ID="your-task-id-here"
while true; do
  RESPONSE=$(curl -s "http://localhost:8080/tasks/$TASK_ID")
  STATUS=$(echo "$RESPONSE" | jq -r '.status')
  PROGRESS=$(echo "$RESPONSE" | jq -r '.progress')
  echo "[$(date +%H:%M:%S)] 状态: $STATUS, 进度: $PROGRESS%"

  if [ "$STATUS" = "completed" ] || [ "$STATUS" = "failed" ]; then
    break
  fi

  sleep 2
done
```

## 测试文件位置

下载完成后，文件会保存在 downloads 目录：

```
downloads/
├── 播客标题/
│   ├── podcast.m4a       # 音频文件
│   ├── cover.jpg         # 封面图片
│   ├── shownotes.txt     # 节目笔记
│   └── .metadata.json    # 元数据
```

## HTTP状态码参考

| 状态码 | 场景 |
|--------|------|
| 200 | OK - 请求成功（如查询已完成任务、列表播客） |
| 202 | Accepted - 任务已接受，正在处理 |
| 400 | Bad Request - 请求参数错误（无效URL、无效分页参数） |
| 404 | Not Found - 任务不存在 |
| 409 | Conflict - 任务已存在（重复提交） |
| 500 | Internal Server Error - 服务器内部错误 |

## 常见问题

**Q: 任务一直处于 pending 状态？**
A: 检查服务器日志，可能是网络问题或URL无法访问。

**Q: 下载失败？**
A: 使用 `GET /tasks/{id}` 查看错误信息，检查网络连接和URL有效性。

**Q: 如何查看下载的文件？**
A: 使用 `ls -la ./downloads/` 查看下载目录。

**Q: 如何清空任务列表？**
A: 重启服务器会清空内存中的任务列表，但已下载的文件和目录扫描的播客目录会保留。

## 使用示例URL

测试时可以使用 `examples/xiaoyuzhou_urls` 文件中的URL：

```bash
# 随机选择一个URL测试
TEST_URL=$(shuf -n 1 examples/xiaoyuzhou_urls)
curl -X POST http://localhost:8080/tasks \
  -H "Content-Type: application/json" \
  -d "{\"url\": \"$TEST_URL\"}"
```
