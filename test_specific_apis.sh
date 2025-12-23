#!/bin/bash

echo "=== 检查重要接口响应格式 ==="

BASE_URL="http://localhost:8080"

# 1. 登录
LOGIN_RESPONSE=$(curl -s -X POST $BASE_URL/auth/login \
  -H "Content-Type: application/json" \
  -d '{"account":"testuser1766460501","password":"TestPass123"}')

TOKEN=$(echo $LOGIN_RESPONSE | jq -r '.data.token')
echo "1. 登录接口:"
echo $LOGIN_RESPONSE | jq '.'

# 2. 获取动态列表
echo -e "\n2. 获取动态列表:"
MOMENTS_RESPONSE=$(curl -s $BASE_URL/moments)
echo $MOMENTS_RESPONSE | jq '.'

# 3. 获取用户信息
echo -e "\n3. 获取用户信息:"
PROFILE_RESPONSE=$(curl -s "$BASE_URL/api/users/profile" \
  -H "Authorization: Bearer $TOKEN")
echo $PROFILE_RESPONSE | jq '.'

# 4. 文件上传测试（创建一个测试文件）
echo -e "\n4. 创建测试文件并上传:"
echo "test content" > /tmp/test.txt
UPLOAD_RESPONSE=$(curl -s -X POST "$BASE_URL/api/upload/file" \
  -H "Authorization: Bearer $TOKEN" \
  -F "file=@/tmp/test.txt")
echo $UPLOAD_RESPONSE | jq '.'

# 5. 更新资料
echo -e "\n5. 更新用户资料:"
UPDATE_RESPONSE=$(curl -s -X PUT "$BASE_URL/api/users/profile" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"nickname":"测试昵称","bio":"测试简介"}')
echo $UPDATE_RESPONSE | jq '.'

echo -e "\n=== 响应格式分析 ==="
echo "✅ 所有成功响应应该都是: \"code\": 200"
echo "✅ 响应结构: {code: 200, message: \"消息\", data: {数据}}"