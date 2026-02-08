#!/bin/bash
# 批量下载小宇宙FM播客单集
# Batch download Xiaoyuzhou FM podcast episodes

# 确保编译好的程序存在
if [ ! -f "../podcast-downloader" ]; then
    echo "错误: 找不到 podcast-downloader 程序"
    echo "请先运行: go build -o podcast-downloader cmd/podcast-downloader/main.go"
    exit 1
fi

# 播客URL列表
# 将你要下载的播客URL添加到这里
URLS=(
  "https://www.xiaoyuzhoufm.com/episode/69392768281939cce65925d3"
  "https://www.xiaoyuzhoufm.com/episode/6774d5c315a5fd520e309381"
  # "https://www.xiaoyuzhoufm.com/episode/another-episode-id"
  # "https://www.xiaoyuzhoufm.com/episode/yet-another-episode-id"
)

# 下载配置
OUTPUT_DIR="./batch_downloads"
TIMEOUT=60s
MAX_RETRIES=3

# 创建输出目录
mkdir -p "$OUTPUT_DIR"

echo "开始批量下载..."
echo "输出目录: $OUTPUT_DIR"
echo "URL数量: ${#URLS[@]}"
echo ""

# 计数器
SUCCESS=0
FAILED=0
SKIPPED=0

# 遍历URL列表
for i in "${!URLS[@]}"; do
    url="${URLS[$i]}"
    num=$((i + 1))
    total=${#URLS[@]}
    
    echo "========================================="
    echo "正在处理 [$num/$total]: $url"
    echo "========================================="
    
    # 尝试下载
    if ../podcast-downloader -o "$OUTPUT_DIR" --timeout "$TIMEOUT" --retry "$MAX_RETRIES" "$url"; then
        SUCCESS=$((SUCCESS + 1))
        echo "✓ 下载成功"
    else
        exit_code=$?
        if [ $exit_code -eq 1 ]; then
            # 检查是否是文件已存在的错误
            if grep -q "文件已存在" <<< "$(../podcast-downloader -o "$OUTPUT_DIR" "$url" 2>&1)"; then
                SKIPPED=$((SKIPPED + 1))
                echo "⊘ 文件已存在，跳过"
            else
                FAILED=$((FAILED + 1))
                echo "✗ 下载失败"
            fi
        else
            FAILED=$((FAILED + 1))
            echo "✗ 下载失败 (退出码: $exit_code)"
        fi
    fi
    echo ""
done

# 打印汇总
echo "========================================="
echo "批量下载完成！"
echo "========================================="
echo "成功: $SUCCESS"
echo "跳过: $SKIPPED"
echo "失败: $FAILED"
echo "总计: $total"
echo "========================================="
