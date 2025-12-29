#!/bin/bash

API_BASE="http://localhost:8080"

echo "🧪 进行更多接口详细测试"
echo "=================================="

# 创建测试用户并获取token
echo "📝 创建测试用户..."
response=$(curl -s -X POST -H "Content-Type: application/json" \
    -d '{"username":"testuser2","phone":"13800138001","password":"test123456"}' \
    "$API_BASE/auth/register")

echo "注册响应: $response"

# 登录获取token
response=$(curl -s -X POST -H "Content-Type: application/json" \
    -d '{"account":"testuser2","password":"test123456"}' \
    "$API_BASE/auth/login")

TOKEN=$(echo "$response" | grep -o '"token":"[^"]*"' | cut -d'"' -f4)
echo "🔑 获取Token: ${TOKEN:0:30}..."

if [ ! -z "$TOKEN" ]; then
    AUTH_HEADER="Authorization: Bearer $TOKEN"
    echo "✅ 认证成功，开始测试需要认证的接口"
    
    echo ""
    echo "📋 测试帖子创建..."
    response=$(curl -s -w "%{http_code}" -o /tmp/post_response.json \
        -X POST -H "Content-Type: application/json" -H "$AUTH_HEADER" \
        -d '{"title":"API测试帖子","content":"这是一个API测试帖子内容","images":[],"video":"","visibility":0,"tags":["API测试"]}' \
        "$API_BASE/api/posts")
    echo "创建帖子: $response"
    
    # 获取创建的帖子ID
    if [ -f /tmp/post_response.json ]; then
        POST_ID=$(grep -o '"id":[0-9]*' /tmp/post_response.json | cut -d':' -f2)
        echo "📄 帖子ID: $POST_ID"
        
        if [ ! -z "$POST_ID" ]; then
            echo ""
            echo "📋 测试帖子详情获取..."
            response=$(curl -s -w "%{http_code}" -o /tmp/detail_response.json \
                "$API_BASE/public/posts/$POST_ID")
            echo "帖子详情: $response"
            
            echo ""
            echo "📋 测试帖子点赞..."
            response=$(curl -s -w "%{http_code}" -o /tmp/like_response.json \
                -X POST -H "$AUTH_HEADER" \
                "$API_BASE/api/likes/post/$POST_ID")
            echo "帖子点赞: $response"
            
            echo ""
            echo "📋 测试评论创建..."
            response=$(curl -s -w "%{http_code}" -o /tmp/comment_response.json \
                -X POST -H "Content-Type: application/json" -H "$AUTH_HEADER" \
                -d '{"content":"这是一个API测试评论"}' \
                "$API_BASE/api/comments/post/$POST_ID")
            echo "创建评论: $response"
        fi
    fi
    
    echo ""
    echo "📋 测试我的帖子列表..."
    response=$(curl -s -w "%{http_code}" -o /tmp/myposts_response.json \
        -H "$AUTH_HEADER" \
        "$API_BASE/api/posts/my")
    echo "我的帖子: $response"
    
    echo ""
    echo "📋 测试标签创建..."
    response=$(curl -s -w "%{http_code}" -o /tmp/tag_response.json \
        -X POST -H "Content-Type: application/json" -H "$AUTH_HEADER" \
        -d '{"name":"API测试标签","color":"#FF6B6B","icon":"test","description":"通过API创建的测试标签"}' \
        "$API_BASE/api/tags")
    echo "创建标签: $response"
    
    echo ""
    echo "📋 测试好友列表..."
    response=$(curl -s -w "%{http_code}" -o /tmp/friends_response.json \
        -H "$AUTH_HEADER" \
        "$API_BASE/api/friends")
    echo "好友列表: $response"
    
    echo ""
    echo "📋 测试会话列表..."
    response=$(curl -s -w "%{http_code}" -o /tmp/conversations_response.json \
        -H "$AUTH_HEADER" \
        "$API_BASE/api/conversations")
    echo "会话列表: $response"
    
else
    echo "❌ 无法获取Token，跳过认证接口测试"
fi

echo ""
echo "=================================="
echo "🎯 详细测试完成！"

# 清理临时文件
rm -f /tmp/post_response.json /tmp/detail_response.json /tmp/like_response.json \
      /tmp/comment_response.json /tmp/myposts_response.json /tmp/tag_response.json \
      /tmp/friends_response.json /tmp/conversations_response.json