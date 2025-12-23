#!/bin/bash

echo "=== 搜索功能测试 ==="

# 1. 登录获取token
echo "1. 登录..."
LOGIN_RESPONSE=$(curl -s -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{"account":"testuser1766460501","password":"TestPass123"}')

echo "Login response: $LOGIN_RESPONSE"

TOKEN=$(echo $LOGIN_RESPONSE | jq -r '.data.token')
echo "Token: $TOKEN"

# 2. 测试热词（无需认证）
echo -e "\n2. 测试热词..."
curl -s http://localhost:8080/search/hot-words | jq .

# 3. 测试搜索
echo -e "\n3. 测试搜索..."
curl -s "http://localhost:8080/api/search?keyword=test" \
  -H "Authorization: Bearer $TOKEN" | jq .

echo -e "\n=== 测试完成 ==="