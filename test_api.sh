#!/bin/bash

API_BASE="http://localhost:8080"

echo "🚀 开始测试 Campus Moments Go API 接口"
echo "=================================="

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 测试函数
test_api() {
    local method=$1
    local url=$2
    local data=$3
    local description=$4
    
    echo -n "测试 $description: "
    
    if [ "$method" = "GET" ]; then
        response=$(curl -s -w "%{http_code}" -o /tmp/api_response.json "$API_BASE$url")
    elif [ "$method" = "POST" ]; then
        response=$(curl -s -w "%{http_code}" -o /tmp/api_response.json -X POST -H "Content-Type: application/json" -d "$data" "$API_BASE$url")
    elif [ "$method" = "PUT" ]; then
        response=$(curl -s -w "%{http_code}" -o /tmp/api_response.json -X PUT -H "Content-Type: application/json" -d "$data" "$API_BASE$url")
    fi
    
    if [ "$response" = "200" ] || [ "$response" = "201" ]; then
        echo -e "${GREEN}✅ $response${NC}"
        # 显示响应内容（前200字符）
        if [ -f /tmp/api_response.json ]; then
            content=$(head -c 200 /tmp/api_response.json)
            echo "   响应: $content..."
        fi
    else
        echo -e "${RED}❌ $response${NC}"
        if [ -f /tmp/api_response.json ]; then
            echo "   错误: $(cat /tmp/api_response.json)"
        fi
    fi
    echo ""
}

echo "📋 1. 系统基础接口"
test_api "GET" "/" "" "首页"
test_api "GET" "/health" "" "健康检查"

echo ""
echo "📋 2. 公开接口"
test_api "GET" "/public/posts" "" "公开帖子列表"
test_api "GET" "/public/tags" "" "公开标签列表"
test_api "GET" "/public/tags/hot" "" "热门标签"

echo ""
echo "📋 3. 认证接口"
test_api "POST" '/auth/register' '{"username":"testuser","phone":"13800138000","password":"test123456"}' "用户注册"
test_api "POST" '/auth/login' '{"account":"testuser","password":"test123456"}' "用户登录"

# 提取token用于后续测试
if [ -f /tmp/api_response.json ]; then
    TOKEN=$(grep -o '"token":"[^"]*"' /tmp/api_response.json | cut -d'"' -f4)
    echo "🔑 获取到Token: ${TOKEN:0:20}..."
    AUTH_HEADER="Authorization: Bearer $TOKEN"
else
    echo "❌ 无法获取Token，跳过需要认证的接口测试"
    AUTH_HEADER=""
fi

echo ""
echo "📋 4. 帖子接口 (需要认证)"
if [ ! -z "$AUTH_HEADER" ]; then
    test_api "POST" '/api/posts' '{"title":"测试帖子","content":"这是一个测试帖子内容","images":[],"video":"","visibility":0,"tags":["测试"]}' "创建帖子"
    test_api "GET" "/api/posts/my" "" "我的帖子"
    test_api "GET" "/api/posts/user/testuser" "" "用户帖子"
else
    echo "⚠️  跳过帖子接口测试(需要认证)"
fi

echo ""
echo "📋 5. 标签接口"
test_api "GET" "/public/tags/search?keyword=测试" "" "搜索标签"
test_api "GET" "/public/tags/by-name/测试/posts" "" "标签相关帖子"

echo ""
echo "📋 6. 用户接口 (需要认证)"
if [ ! -z "$AUTH_HEADER" ]; then
    test_api "GET" "/api/users/profile" "" "用户信息"
else
    echo "⚠️  跳过用户接口测试(需要认证)"
fi

echo ""
echo "📋 7. 搜索接口 (需要认证)"
if [ ! -z "$AUTH_HEADER" ]; then
    test_api "GET" "/api/search?keyword=测试" "" "搜索内容"
    test_api "GET" "/api/search/history" "" "搜索历史"
else
    echo "⚠️  跳过搜索接口测试(需要认证)"
fi

echo ""
echo "📋 8. 好友接口 (需要认证)"
if [ ! -z "$AUTH_HEADER" ]; then
    test_api "GET" "/api/friends" "" "好友列表"
    test_api "GET" "/api/friends/requests" "" "好友请求"
else
    echo "⚠️  跳过好友接口测试(需要认证)"
fi

echo ""
echo "📋 9. 消息接口 (需要认证)"
if [ ! -z "$AUTH_HEADER" ]; then
    test_api "GET" "/api/conversations" "" "会话列表"
    test_api "GET" "/api/conversations/unread" "" "未读消息数"
else
    echo "⚠️  跳过消息接口测试(需要认证)"
fi

echo ""
echo "📋 10. 上传接口 (需要认证)"
if [ ! -z "$AUTH_HEADER" ]; then
    echo "上传接口需要文件，跳过文件上传测试"
else
    echo "⚠️  跳过上传接口测试(需要认证)"
fi

echo ""
echo "=================================="
echo "🎯 测试完成！请查看上方结果"

# 清理临时文件
rm -f /tmp/api_response.json