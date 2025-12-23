#!/bin/bash

echo "=== 搜索功能完整测试 ==="

# 1. 登录
echo "1. 登录..."
LOGIN_RESPONSE=$(curl -s -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{"account":"testuser1766460501","password":"TestPass123"}')

TOKEN=$(echo $LOGIN_RESPONSE | jq -r '.data.token')
echo "Token获取成功: ${TOKEN:0:50}..."

# 2. 测试搜索不同排序
echo -e "\n2. 测试搜索功能..."

# 2.1 测试空结果
echo "2.1 搜索不存在的关键词..."
curl -s "http://localhost:8080/api/search?keyword=test" \
  -H "Authorization: Bearer $TOKEN" | jq '.data.pagination'

# 2.2 测试排序选项
echo "2.2 测试排序选项..."
for sort in "latest" "hottest" "comprehensive"; do
  echo "  - 测试 $sort 排序..."
  curl -s "http://localhost:8080/api/search?keyword=学习&sortBy=$sort" \
    -H "Authorization: Bearer $TOKEN" | jq '.code'
done

# 3. 测试搜索历史
echo -e "\n3. 测试搜索历史..."
echo "3.1 获取搜索历史..."
curl -s "http://localhost:8080/api/search/history" \
  -H "Authorization: Bearer $TOKEN" | jq '.data'

echo "3.2 保存搜索历史..."
curl -s -X POST "http://localhost:8080/api/search/history" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"keyword":"测试搜索"}' | jq '.code'

# 4. 测试热词
echo -e "\n4. 测试热词..."
curl -s "http://localhost:8080/search/hot-words" | jq '.data[0:5]'

echo -e "\n=== 测试完成 ==="