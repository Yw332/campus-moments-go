#!/bin/bash

echo "测试服务器状态..."
echo "=================="

echo "1. 测试服务器8080端口连接:"
timeout 5 bash -c "</dev/tcp/106.52.165.122/8080" && echo "✅ 端口可连接" || echo "❌ 端口不可连接"

echo ""
echo "2. 测试HTTP响应:"
response=$(curl -s --connect-timeout 3 --max-time 5 http://106.52.165.122:8080/health 2>/dev/null)
if [ -n "$response" ]; then
    echo "✅ 有响应: $response"
else
    echo "❌ 无响应或空响应"
fi

echo ""
echo "3. 测试本地对比:"
local_response=$(curl -s --connect-timeout 3 http://localhost:8080/health 2>/dev/null)
if [ -n "$local_response" ]; then
    echo "✅ 本地正常: $local_response"
else
    echo "❌ 本地也有问题"
fi

echo ""
echo "=================="
echo "如果端口可连接但无HTTP响应，说明服务在运行但处理请求有问题。"