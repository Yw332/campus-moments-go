#!/bin/bash

# Campus Moments API 测试脚本
# 测试所有接口（除了好友相关）

BASE_URL="http://localhost:8080"
TOKEN=""

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 测试函数
test_endpoint() {
    local method=$1
    local endpoint=$2
    local data=$3
    local description=$4
    local headers=$5
    
    echo -e "\n${BLUE}=== 测试: $description ===${NC}"
    echo -e "${YELLOW}请求: $method $endpoint${NC}"
    
    if [ -n "$data" ]; then
        echo -e "${YELLOW}数据: $data${NC}"
    fi
    
    if [ -n "$headers" ]; then
        curl -X $method \
             -H "Content-Type: application/json" \
             $headers \
             -d "$data" \
             "$BASE_URL$endpoint" \
             -w "\n状态码: %{http_code}\n" \
             -s
    else
        curl -X $method \
             -H "Content-Type: application/json" \
             -d "$data" \
             "$BASE_URL$endpoint" \
             -w "\n状态码: %{http_code}\n" \
             -s
    fi
    
    echo -e "\n${GREEN}--- 测试完成 ---${NC}"
}

# 等待服务器启动
echo -e "${YELLOW}等待服务器启动...${NC}"
sleep 5

# 1. 测试健康检查
test_endpoint "GET" "/health" "" "健康检查"

# 2. 用户注册
echo -e "\n${BLUE}=== 开始用户注册测试 ===${NC}"
REGISTER_DATA='{
    "phone": "17875242006",
    "password": "test123456",
    "username": "testuser002",
    "nickname": "测试用户002"
}'
response=$(curl -X POST -H "Content-Type: application/json" -d "$REGISTER_DATA" "$BASE_URL/auth/register" -s)
echo $response

# 提取token
TOKEN=$(echo $response | grep -o '"token":"[^"]*' | cut -d'"' -f4)
if [ -n "$TOKEN" ]; then
    echo -e "\n${GREEN}Token获取成功: $TOKEN${NC}"
else
    echo -e "\n${RED}Token获取失败，尝试登录${NC}"
    # 3. 用户登录
    LOGIN_DATA='{
        "account": "17875242006",
        "password": "test123456"
    }'
    response=$(curl -X POST -H "Content-Type: application/json" -d "$LOGIN_DATA" "$BASE_URL/auth/login" -s)
    echo $response
    TOKEN=$(echo $response | grep -o '"token":"[^"]*' | cut -d'"' -f4)
fi

AUTH_HEADER="-H \"Authorization: Bearer $TOKEN\""

# 4. 发送验证码
test_endpoint "POST" "/auth/send-verification" '{
    "phone": "17875242006"
}' "发送验证码"

# 5. 验证并重置密码
test_endpoint "POST" "/auth/verify-and-reset" '{
    "phone": "17875242006",
    "verificationCode": "123456",
    "newPassword": "newtest123456"
}' "验证并重置密码"

# 6. 获取用户资料
test_endpoint "GET" "/api/users/profile" "" "获取用户资料" "-H \"Authorization: Bearer $TOKEN\""

# 7. 更新用户资料
test_endpoint "PUT" "/api/users/profile" '{
    "nickname": "更新的昵称",
    "bio": "这是我的个人简介",
    "avatar": "http://example.com/avatar.jpg"
}' "更新用户资料" "-H \"Authorization: Bearer $TOKEN\""

# 8. 修改密码
test_endpoint "PUT" "/api/users/password" '{
    "oldPassword": "test123456",
    "newPassword": "test123456"
}' "修改密码" "-H \"Authorization: Bearer $TOKEN\""

# 9. 发布动态
test_endpoint "POST" "/api/moments" '{
    "content": "这是我的测试动态内容",
    "tags": ["测试", "API", "校园"],
    "visibility": 0
}' "发布动态" "-H \"Authorization: Bearer $TOKEN\""

# 10. 获取动态列表
test_endpoint "GET" "/moments?page=1&pageSize=10" "" "获取动态列表"

# 11. 获取我的动态列表
test_endpoint "GET" "/api/moments/my?page=1&pageSize=10" "" "获取我的动态列表" "-H \"Authorization: Bearer $TOKEN\""

# 12. 搜索内容
test_endpoint "GET" "/api/search?keyword=测试" "" "搜索内容" "-H \"Authorization: Bearer $TOKEN\""

# 13. 获取热门关键词
test_endpoint "GET" "/api/search/hot-words" "" "获取热门关键词" "-H \"Authorization: Bearer $TOKEN\""

# 14. 获取搜索历史
test_endpoint "GET" "/api/search/history" "" "获取搜索历史" "-H \"Authorization: Bearer $TOKEN\""

echo -e "\n${GREEN}=== 所有接口测试完成 ===${NC}"