#!/bin/bash

echo "=== 服务器诊断脚本 ==="
echo "服务器地址: 106.52.165.122:8080"
echo "测试时间: $(date)"
echo ""

# 测试1: 端口连通性
echo "1. 测试端口连通性..."
if command -v nc >/dev/null 2>&1; then
    echo -n "   nc 测试: "
    nc -zv 106.52.165.122 8080 2>&1
elif command -v telnet >/dev/null 2>&1; then
    echo "   telnet 测试: 请运行 'telnet 106.52.165.122 8080'"
else
    echo "   缺少 nc 和 telnet 命令"
fi

# 测试2: HTTP连接
echo ""
echo "2. 测试HTTP连接..."
echo "   curl 测试:"
curl -v --connect-timeout 10 http://106.52.165.122:8080/health 2>&1

# 测试3: 与本地对比
echo ""
echo "3. 本地服务测试..."
echo "   curl 本地:"
curl -s --connect-timeout 3 http://localhost:8080/health 2>/dev/null | head -3

echo ""
echo "=== 诊断完成 ==="
echo ""
echo "如果服务器访问失败，请检查："
echo "1. SSH登录服务器: ssh root@106.52.165.122"
echo "2. 检查服务状态: ps aux | grep campus-moments"
echo "3. 检查端口监听: netstat -tlnp | grep 8080"
echo "4. 检查防火墙: ufw status 或 firewall-cmd --list-ports"
echo "5. 重启服务: ./start.sh"