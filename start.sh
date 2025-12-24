#!/bin/bash

echo "=== Campus Moments Server Startup ==="

# 加载环境变量
if [ -f .env ]; then
    export $(cat .env | grep -v '^#' | xargs)
    echo "Environment variables loaded from .env file"
fi

# 设置默认环境变量
export PORT=${PORT:-8080}
export APP_ENV=${APP_ENV:-production}
export GIN_MODE=${GIN_MODE:-release}

echo "Starting application on port: $PORT"
echo "Environment: $APP_ENV"
echo "Gin Mode: $GIN_MODE"

# 构建应用
echo "Building application..."
go build -o bin/main cmd/api/main.go

if [ $? -eq 0 ]; then
    echo "Build successful, starting application..."
    ./bin/main
else
    echo "Build failed"
    exit 1
fi