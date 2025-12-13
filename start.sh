#!/bin/bash

echo "=== CNB Environment Startup ==="

# 设置环境变量
export PORT=${PORT:-8080}
export GIN_MODE=release

echo "Starting application on port: $PORT"

# 构建应用
echo "Building application..."
go build -o bin/api cmd/api/main.go

if [ $? -eq 0 ]; then
    echo "Build successful, starting application..."
    ./bin/api
else
    echo "Build failed"
    exit 1
fi