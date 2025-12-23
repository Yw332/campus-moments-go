#!/bin/bash

# 测试搜索功能的完整流程

echo "🔍 开始测试搜索功能..."

# 1. 获取热词（无需认证）
echo "1. 测试获取热词..."
curl -s -X GET http://localhost:8080/search/hot-words | jq .

# 2. 用户登录
echo "2. 用户登录..."
LOGIN_RESPONSE=$(curl -s -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{"account":"testuser1766460501","password":"TestPass123"}')

TOKEN=$(echo $LOGIN_RESPONSE | jq -r '.data.token')
echo "获取到的Token: $TOKEN"

# 3. 测试搜索关键词（不同排序方式）
echo "3. 测试搜索功能..."

# 3.1 按最新排序
echo "3.1 按最新排序搜索..."
curl -s -X GET "http://localhost:8080/api/search?keyword=学习&sortBy=latest" \
  -H "Authorization: Bearer $TOKEN" | jq .

# 3.2 按最热排序
echo "3.2 按最热排序搜索..."
curl -s -X GET "http://localhost:8080/api/search?keyword=学习&sortBy=hottest" \
  -H "Authorization: Bearer $TOKEN" | jq .

# 3.3 按综合排序
echo "3.3 按综合排序搜索..."
curl -s -X GET "http://localhost:8080/api/search?keyword=学习&sortBy=comprehensive" \
  -H "Authorization: Bearer $TOKEN" | jq .

# 4. 等待一下让异步保存完成
sleep 2

# 5. 获取搜索历史
echo "4. 获取搜索历史..."
curl -s -X GET http://localhost:8080/api/search/history \
  -H "Authorization: Bearer $TOKEN" | jq .

# 6. 保存搜索历史（手动）
echo "5. 手动保存搜索历史..."
curl -s -X POST http://localhost:8080/api/search/history \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"keyword":"校园活动"}' | jq .

# 7. 再次获取搜索历史验证保存
echo "6. 验证搜索历史保存..."
curl -s -X GET http://localhost:8080/api/search/history \
  -H "Authorization: Bearer $TOKEN" | jq .

echo "✅ 搜索功能测试完成！"