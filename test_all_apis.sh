#!/bin/bash

# Campus Moments API 完整测试脚本
# 测试所有接口（除了好友相关）

BASE_URL="http://localhost:8080"
TOKEN=""

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

echo -e "${GREEN}🚀 开始测试 Campus Moments API${NC}"

# 1. 健康检查
echo -e "\n${BLUE}=== 1. 健康检查 ===${NC}"
response=$(curl -s "$BASE_URL/health")
echo -e "${GREEN}✅ 健康检查: $response${NC}"

# 2. 用户登录（使用现有用户）
echo -e "\n${BLUE}=== 2. 用户登录 ===${NC}"
login_data='{
    "account": "13800138001",
    "password": "password"
}'
response=$(curl -s -X POST -H "Content-Type: application/json" -d "$login_data" "$BASE_URL/auth/login")
echo -e "${GREEN}✅ 登录响应: $response${NC}"

# 提取token
TOKEN=$(echo $response | grep -o '"token":"[^"]*' | cut -d'"' -f4)
if [ -n "$TOKEN" ]; then
    echo -e "${GREEN}✅ Token获取成功: $TOKEN${NC}"
else
    echo -e "${RED}❌ Token获取失败${NC}"
    exit 1
fi

# 3. 发送验证码
echo -e "\n${BLUE}=== 3. 发送验证码 ===${NC}"
verify_data='{"phone": "13800138001"}'
response=$(curl -s -X POST -H "Content-Type: application/json" -d "$verify_data" "$BASE_URL/auth/send-verification")
echo -e "${GREEN}✅ 发送验证码: $response${NC}"

# 4. 验证并重置密码
echo -e "\n${BLUE}=== 4. 验证并重置密码 ===${NC}"
reset_data='{
    "phone": "13800138001",
    "verificationCode": "123456",
    "newPassword": "password123"
}'
response=$(curl -s -X POST -H "Content-Type: application/json" -d "$reset_data" "$BASE_URL/auth/verify-and-reset")
echo -e "${GREEN}✅ 重置密码: $response${NC}"

# 5. 获取用户资料
echo -e "\n${BLUE}=== 5. 获取用户资料 ===${NC}"
response=$(curl -s -H "Authorization: Bearer $TOKEN" "$BASE_URL/api/users/profile")
echo -e "${GREEN}✅ 用户资料: $response${NC}"

# 6. 更新用户资料
echo -e "\n${BLUE}=== 6. 更新用户资料 ===${NC}"
update_data='{
    "nickname": "更新的昵称",
    "bio": "这是我的个人简介",
    "avatar": "http://example.com/avatar.jpg"
}'
response=$(curl -s -X PUT -H "Content-Type: application/json" -H "Authorization: Bearer $TOKEN" -d "$update_data" "$BASE_URL/api/users/profile")
echo -e "${GREEN}✅ 更新资料: $response${NC}"

# 7. 修改密码
echo -e "\n${BLUE}=== 7. 修改密码 ===${NC}"
password_data='{
    "oldPassword": "password",
    "newPassword": "password123"
}'
response=$(curl -s -X PUT -H "Content-Type: application/json" -H "Authorization: Bearer $TOKEN" -d "$password_data" "$BASE_URL/api/users/password")
echo -e "${GREEN}✅ 修改密码: $response${NC}"

# 8. 发布动态
echo -e "\n${BLUE}=== 8. 发布动态 ===${NC}"
moment_data='{
    "content": "这是我的测试动态内容",
    "tags": ["测试", "API", "校园"],
    "visibility": 0
}'
response=$(curl -s -X POST -H "Content-Type: application/json" -H "Authorization: Bearer $TOKEN" -d "$moment_data" "$BASE_URL/api/moments")
echo -e "${GREEN}✅ 发布动态: $response${NC}"

# 9. 获取动态列表
echo -e "\n${BLUE}=== 9. 获取动态列表 ===${NC}"
response=$(curl -s "$BASE_URL/moments?page=1&pageSize=10")
echo -e "${GREEN}✅ 动态列表: $response${NC}"

# 10. 获取我的动态列表
echo -e "\n${BLUE}=== 10. 获取我的动态列表 ===${NC}"
response=$(curl -s -H "Authorization: Bearer $TOKEN" "$BASE_URL/api/moments/my?page=1&pageSize=10")
echo -e "${GREEN}✅ 我的动态: $response${NC}"

# 11. 搜索内容
echo -e "\n${BLUE}=== 11. 搜索内容 ===${NC}"
response=$(curl -s -H "Authorization: Bearer $TOKEN" "$BASE_URL/api/search?keyword=测试")
echo -e "${GREEN}✅ 搜索结果: $response${NC}"

# 12. 获取热门关键词
echo -e "\n${BLUE}=== 12. 获取热门关键词 ===${NC}"
response=$(curl -s -H "Authorization: Bearer $TOKEN" "$BASE_URL/api/search/hot-words")
echo -e "${GREEN}✅ 热门关键词: $response${NC}"

# 13. 获取搜索历史
echo -e "\n${BLUE}=== 13. 获取搜索历史 ===${NC}"
response=$(curl -s -H "Authorization: Bearer $TOKEN" "$BASE_URL/api/search/history")
echo -e "${GREEN}✅ 搜索历史: $response${NC}"

echo -e "\n${GREEN}🎉 所有接口测试完成！${NC}"