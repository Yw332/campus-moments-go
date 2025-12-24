# 使用 Go 1.25.3 官方镜像作为基础镜像
FROM golang:1.25.3-alpine AS builder

# 设置工作目录
WORKDIR /app

# 安装必要的包
RUN apk add --no-cache git tzdata ca-certificates

# 复制 go mod 文件
COPY go.mod go.sum ./

# 下载依赖
RUN go mod download

# 复制源代码
COPY . .

# 构建应用
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main cmd/main.go

# 使用轻量级镜像作为最终镜像
FROM alpine:latest

# 安装 ca-certificates 和 tzdata
RUN apk --no-cache add ca-certificates tzdata

# 设置工作目录
WORKDIR /root/

# 从构建阶段复制二进制文件
COPY --from=builder /app/main .
COPY --from=builder /app/start.sh .
COPY --from=builder /app/.env.example .env

# 设置权限
RUN chmod +x start.sh

# 暴露端口
EXPOSE 8080

# 启动命令
CMD ["./start.sh"]