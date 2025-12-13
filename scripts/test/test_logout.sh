#!/bin/bash

API_BASE="http://localhost:8080"

echo "=== 退出登录功能测试 ==="
echo ""

# 步骤1: 登录获取Token
echo "1. 登录获取Token..."
LOGIN_RESPONSE=$(curl -s -X POST "$API_BASE/auth/login" \
  -H "Content-Type: application/json" \
  -d '{"account":"Yw166332","password":"Aa123456"}')

echo "登录响应: $LOGIN_RESPONSE"

# 提取token
TOKEN=$(echo $LOGIN_RESPONSE | sed -n 's/.*"token":"\([^"]*\)".*/\1/p')
echo "提取的Token: $TOKEN"
echo ""

# 步骤2: 使用Token获取用户信息（验证Token有效）
echo "2. 验证Token有效..."
curl -s -X GET "$API_BASE/api/users/profile" \
  -H "Authorization: Bearer $TOKEN" | jq .
echo ""

# 步骤3: 退出登录
echo "3. 退出登录..."
LOGOUT_RESPONSE=$(curl -s -X POST "$API_BASE/api/auth/logout" \
  -H "Authorization: Bearer $TOKEN")
echo "退出响应: $LOGOUT_RESPONSE"
echo ""

# 步骤4: 验证Token已失效
echo "4. 验证Token已失效..."
curl -s -X GET "$API_BASE/api/users/profile" \
  -H "Authorization: Bearer $TOKEN" | jq .
echo ""

echo "=== 测试完成 ==="