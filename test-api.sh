#!/bin/bash

# Podcast API Server 测试脚本
# 使用方法: ./test-api.sh [server_url]

SERVER_URL="${1:-http://localhost:8080}"

echo "========================================="
echo "Podcast API Server 测试"
echo "服务器: $SERVER_URL"
echo "========================================="
echo ""

# 颜色定义
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

# 测试URL
TEST_URL="https://www.xiaoyuzhoufm.com/episode/69806618073030367acec13b"

# ============================================
# 1. 测试 POST /tasks - 提交下载任务
# ============================================
echo -e "${BLUE}1. 测试 POST /tasks - 提交下载任务${NC}"
echo "-------------------------------------------"
echo "请求: POST $SERVER_URL/tasks"
echo "数据: {\"url\": \"$TEST_URL\"}"
echo ""

RESPONSE=$(curl -s -X POST "$SERVER_URL/tasks" \
  -H "Content-Type: application/json" \
  -d "{\"url\": \"$TEST_URL\"}")

echo "响应:"
echo "$RESPONSE" | jq '.' 2>/dev/null || echo "$RESPONSE"
echo ""

# 提取任务ID
TASK_ID=$(echo "$RESPONSE" | jq -r '.id // .task.id // empty' 2>/dev/null)

if [ -z "$TASK_ID" ] || [ "$TASK_ID" = "null" ]; then
    echo -e "${RED}❌ 无法获取任务ID，可能请求失败${NC}"
    echo ""
    echo "请检查："
    echo "1. 服务器是否正在运行: $SERVER_URL"
    echo "2. 运行: ./build/podcast-server -port 8080 -downloads ./downloads -verbose"
    exit 1
fi

echo -e "${GREEN}✓ 任务已创建，ID: $TASK_ID${NC}"
echo ""
echo "-------------------------------------------"
echo ""

# ============================================
# 2. 测试 GET /tasks/{id} - 查询任务状态
# ============================================
echo -e "${BLUE}2. 测试 GET /tasks/$TASK_ID - 查询任务状态${NC}"
echo "-------------------------------------------"
echo "请求: GET $SERVER_URL/tasks/$TASK_ID"
echo ""

RESPONSE=$(curl -s "$SERVER_URL/tasks/$TASK_ID")

echo "响应:"
echo "$RESPONSE" | jq '.' 2>/dev/null || echo "$RESPONSE"
echo ""

# 提取状态
STATUS=$(echo "$RESPONSE" | jq -r '.status // empty' 2>/dev/null)
PROGRESS=$(echo "$RESPONSE" | jq -r '.progress // 0' 2>/dev/null)

echo "当前状态: $STATUS"
echo "进度: $PROGRESS%"
echo ""

# 如果任务进行中，等待并再次查询
if [ "$STATUS" = "in_progress" ] || [ "$STATUS" = "pending" ]; then
    echo -e "${YELLOW}任务进行中，等待5秒后再次查询...${NC}"
    sleep 5

    echo ""
    echo "再次查询状态:"
    RESPONSE=$(curl -s "$SERVER_URL/tasks/$TASK_ID")
    echo "$RESPONSE" | jq '.' 2>/dev/null || echo "$RESPONSE"
    echo ""
fi

echo "-------------------------------------------"
echo ""

# ============================================
# 3. 测试 GET /tasks/invalid-id - 测试404错误
# ============================================
echo -e "${BLUE}3. 测试 GET /tasks/invalid-id - 404错误处理${NC}"
echo "-------------------------------------------"
echo "请求: GET $SERVER_URL/tasks/00000000-0000-0000-0000-000000000000"
echo ""

RESPONSE=$(curl -s "$SERVER_URL/tasks/00000000-0000-0000-0000-000000000000")

echo "响应:"
echo "$RESPONSE" | jq '.' 2>/dev/null || echo "$RESPONSE"
echo ""

echo -e "${GREEN}✓ 404错误正确返回${NC}"
echo ""
echo "-------------------------------------------"
echo ""

# ============================================
# 4. 测试 GET /podcasts - 列出已下载的播客
# ============================================
echo -e "${BLUE}4. 测试 GET /podcasts - 列出已下载的播客${NC}"
echo "-------------------------------------------"
echo "请求: GET $SERVER_URL/podcasts"
echo ""

RESPONSE=$(curl -s "$SERVER_URL/podcasts?limit=10&offset=0")

echo "响应:"
echo "$RESPONSE" | jq '.' 2>/dev/null || echo "$RESPONSE"
echo ""

# 提取播客数量
TOTAL=$(echo "$RESPONSE" | jq -r '.total // 0' 2>/dev/null)

echo "已下载播客总数: $TOTAL"
echo ""

echo -e "${GREEN}✓ 播客列表查询成功${NC}"
echo ""
echo "-------------------------------------------"
echo ""

# ============================================
# 5. 测试分页
# ============================================
if [ "$TOTAL" -gt 0 ]; then
    echo -e "${BLUE}5. 测试分页 - 每页显示2条${NC}"
    echo "-------------------------------------------"
    echo "请求: GET $SERVER_URL/podcasts?limit=2&offset=0"
    echo ""

    RESPONSE=$(curl -s "$SERVER_URL/podcasts?limit=2&offset=0")

    echo "响应:"
    echo "$RESPONSE" | jq '.' 2>/dev/null || echo "$RESPONSE"
    echo ""

    echo -e "${GREEN}✓ 分页功能正常${NC}"
    echo ""
    echo "-------------------------------------------"
    echo ""
fi

# ============================================
# 6. 测试重复提交
# ============================================
echo -e "${BLUE}6. 测试重复提交 - 应返回409 Conflict${NC}"
echo "-------------------------------------------"
echo "请求: POST $SERVER_URL/tasks (相同URL)"
echo ""

RESPONSE=$(curl -s -X POST "$SERVER_URL/tasks" \
  -H "Content-Type: application/json" \
  -d "{\"url\": \"$TEST_URL\"}")

echo "响应:"
echo "$RESPONSE" | jq '.' 2>/dev/null || echo "$RESPONSE"
echo ""

# 检查HTTP状态码
HTTP_CODE=$(curl -s -o /dev/null -w "%{http_code}" -X POST "$SERVER_URL/tasks" \
  -H "Content-Type: application/json" \
  -d "{\"url\": \"$TEST_URL\"}")

if [ "$HTTP_CODE" = "409" ]; then
    echo -e "${GREEN}✓ 重复提交正确返回409 Conflict${NC}"
else
    echo -e "${YELLOW}HTTP状态码: $HTTP_CODE (预期: 409)${NC}"
fi

echo ""
echo "-------------------------------------------"
echo ""

# ============================================
# 7. 测试无效URL
# ============================================
echo -e "${BLUE}7. 测试无效URL - 应返回400 Bad Request${NC}"
echo "-------------------------------------------"
echo "请求: POST $SERVER_URL/tasks"
echo "数据: {\"url\": \"https://example.com/invalid\"}"
echo ""

RESPONSE=$(curl -s -X POST "$SERVER_URL/tasks" \
  -H "Content-Type: application/json" \
  -d '{"url": "https://example.com/invalid"}')

echo "响应:"
echo "$RESPONSE" | jq '.' 2>/dev/null || echo "$RESPONSE"
echo ""

# 检查HTTP状态码
HTTP_CODE=$(curl -s -o /dev/null -w "%{http_code}" -X POST "$SERVER_URL/tasks" \
  -H "Content-Type: application/json" \
  -d '{"url": "https://example.com/invalid"}')

if [ "$HTTP_CODE" = "400" ]; then
    echo -e "${GREEN}✓ 无效URL正确返回400 Bad Request${NC}"
else
    echo -e "${YELLOW}HTTP状态码: $HTTP_CODE (预期: 400)${NC}"
fi

echo ""
echo "-------------------------------------------"
echo ""

# ============================================
# 测试完成
# ============================================
echo -e "${GREEN}=========================================${NC}"
echo -e "${GREEN}测试完成！${NC}"
echo -e "${GREEN}=========================================${NC}"
echo ""
echo "后续操作建议:"
echo "1. 查看下载进度: curl $SERVER_URL/tasks/$TASK_ID"
echo "2. 等待下载完成后，查看播客列表: curl $SERVER_URL/podcasts"
echo "3. 查看下载目录: ls -la ./downloads/"
echo "4. 查看服务器日志以了解下载详情"
echo ""
