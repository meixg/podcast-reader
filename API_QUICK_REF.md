# API 快速参考

## 服务器命令

```bash
# 启动服务器
./build/podcast-server

# 自定义端口和下载目录
./build/podcast-server -port 3000 -downloads ~/podcasts -verbose

# 完整选项
./build/podcast-server -host 0.0.0.0 -port 8080 -downloads ./downloads -verbose
```

## API 端点

### POST /tasks - 提交下载任务
```bash
curl -X POST http://localhost:8080/tasks \
  -H "Content-Type: application/json" \
  -d '{"url": "https://www.xiaoyuzhoufm.com/episode/69806618073030367acec13b"}'
```

### GET /tasks/{id} - 查询任务状态
```bash
curl http://localhost:8080/tasks/{task_id}
```

### GET /podcasts - 列出播客
```bash
# 默认（limit=100, offset=0）
curl http://localhost:8080/podcasts

# 分页
curl "http://localhost:8080/podcasts?limit=10&offset=0"
```

## 快速测试

### 一键测试所有端点
```bash
./test-api.sh
```

### 提交并监控任务
```bash
# 1. 提交任务
RESPONSE=$(curl -s -X POST http://localhost:8080/tasks \
  -H "Content-Type: application/json" \
  -d '{"url": "https://www.xiaoyuzhoufm.com/episode/69806618073030367acec13b"}')

# 2. 提取任务ID
TASK_ID=$(echo $RESPONSE | jq -r '.id')
echo "任务ID: $TASK_ID"

# 3. 监控进度
watch -n 2 "curl -s http://localhost:8080/tasks/$TASK_ID | jq '{status, progress}'"
```

### 批量下载
```bash
# 从文件读取URL并批量提交
cat examples/xiaoyuzhou_urls | while read url; do
  curl -s -X POST http://localhost:8080/tasks \
    -H "Content-Type: application/json" \
    -d "{\"url\": \"$url\"}" | jq '.id, .status'
  sleep 1  # 避免请求过快
done
```

## 响应状态

| 状态 | 说明 |
|------|------|
| `pending` | 任务已创建，等待开始 |
| `in_progress` | 正在下载 |
| `completed` | 下载完成 |
| `failed` | 下载失败 |

## HTTP状态码

- `202` - 任务已接受
- `200` - 查询成功
- `400` - 请求错误
- `404` - 任务不存在
- `409` - 任务已存在

## 文件位置

```
downloads/
└── 播客标题/
    ├── podcast.m4a
    ├── cover.jpg
    ├── shownotes.txt
    └── .metadata.json
```

## 故障排查

```bash
# 查看服务器日志
./build/podcast-server -verbose

# 查看任务错误信息
curl http://localhost:8080/tasks/{id} | jq '.error'

# 查看下载的文件
ls -la ./downloads/

# 测试URL是否有效
curl -I "https://www.xiaoyuzhoufm.com/episode/..."
```
