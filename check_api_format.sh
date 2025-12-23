#!/bin/bash

echo "=== API响应格式一致性检查 ==="

BASE_URL="http://localhost:8080"
TOKEN=""

# 登录获取token
echo "1. 登录获取token..."
LOGIN_RESPONSE=$(curl -s -X POST $BASE_URL/auth/login \
  -H "Content-Type: application/json" \
  -d '{"account":"testuser1766460501","password":"TestPass123"}')

TOKEN=$(echo $LOGIN_RESPONSE | jq -r '.data.token')
echo "登录接口: $(echo $LOGIN_RESPONSE | jq -c '{code: .code, message: .message}')"

# 测试各个接口的响应格式
echo -e "\n2. 检查各接口响应格式..."

# 公共接口
echo "=== 公共接口 ==="
echo "2.1 首页接口:"
curl -s $BASE_URL/ | jq -c '{code: .code, message: .message}'

echo "2.2 健康检查:"
curl -s $BASE_URL/health | jq -c '{code: .code, message: .message}'

echo "2.3 热词接口:"
curl -s $BASE_URL/search/hot-words | jq -c '{code: .code, message: .message}'

# 需要认证的接口
echo -e "\n=== 认证接口 ==="
echo "3.1 搜索接口:"
curl -s "$BASE_URL/api/search?keyword=test" \
  -H "Authorization: Bearer $TOKEN" | jq -c '{code: .code, message: .message}'

echo "3.2 搜索历史:"
curl -s "$BASE_URL/api/search/history" \
  -H "Authorization: Bearer $TOKEN" | jq -c '{code: .code, message: .message}'

echo "3.3 用户资料:"
curl -s "$BASE_URL/api/users/profile" \
  -H "Authorization: Bearer $TOKEN" | jq -c '{code: .code, message: .message}'

echo "3.4 动态列表:"
curl -s "$BASE_URL/moments" | jq -c '{code: .code, message: .message}'

echo -e "\n=== 响应格式分析 ==="
echo "成功响应code应该是: 200"
echo "错误响应code示例: 400, 401, 404, 500"

echo -e "\n=== 检查完成 ==="