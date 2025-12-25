#!/bin/bash

BASE_URL="http://localhost:8080"
TOKEN=""

# 颜色定义
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m'

echo -e "${GREEN}🚀 测试动态相关接口${NC}"

# 1. 登录获取token
echo -e "\n${BLUE}=== 1. 用户登录 ===${NC}"
login_data='{
    "account": "13800138001",
    "password": "password"
}'
response=$(curl -s -X POST -H "Content-Type: application/json" -d "$login_data" "$BASE_URL/auth/login")
echo -e "${GREEN}登录响应: $response${NC}"

TOKEN=$(echo $response | grep -o '"token":"[^"]*' | cut -d'"' -f4)
if [ -n "$TOKEN" ]; then
    echo -e "${GREEN}✅ Token获取成功${NC}"
else
    echo -e "${YELLOW}❌ Token获取失败${NC}"
    exit 1
fi

# 2. 发布带图片的动态
echo -e "\n${BLUE}=== 2. 发布带图片的动态 ===${NC}"
moment_data='{
    "content": "这是修复后的测试动态，包含图片",
    "tags": ["测试", "修复"],
    "media": [
        {
            "type": "image",
            "url": "http://106.52.165.122:8080/static/files/20251225013944_076749e3.png"
        }
    ],
    "visibility": 0
}'
response=$(curl -s -X POST -H "Content-Type: application/json" -H "Authorization: Bearer $TOKEN" -d "$moment_data" "$BASE_URL/api/moments")
echo -e "${GREEN}发布响应: $response${NC}"

# 3. 获取动态列表
echo -e "\n${BLUE}=== 3. 获取动态列表 ===${NC}"
response=$(curl -s "$BASE_URL/moments?page=1&pageSize=5")
echo -e "${GREEN}动态列表: $response${NC}"

# 4. 获取我的动态
echo -e "\n${BLUE}=== 4. 获取我的动态 ===${NC}"
response=$(curl -s -H "Authorization: Bearer $TOKEN" "$BASE_URL/api/moments/my?page=1&pageSize=5")
echo -e "${GREEN}我的动态: $response${NC}"

# 5. 测试搜索功能
echo -e "\n${BLUE}=== 5. 搜索动态 ===${NC}"
response=$(curl -s "$BASE_URL/search?keyword=修复")
echo -e "${GREEN}搜索结果: $response${NC}"

echo -e "\n${GREEN}🎉 动态接口测试完成！${NC}"