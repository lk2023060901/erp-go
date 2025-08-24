# 多阶段构建
FROM golang:1.21-alpine AS builder

WORKDIR /app

# 安装必要的包
RUN apk add --no-cache git

# 复制go mod文件
COPY go.mod go.sum ./

# 下载依赖
RUN go mod download

# 复制源码
COPY . .

# 构建应用
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main cmd/server/main.go

# 运行阶段
FROM alpine:latest

WORKDIR /app

# 安装ca证书
RUN apk --no-cache add ca-certificates tzdata

# 复制构建的二进制文件
COPY --from=builder /app/main .
COPY --from=builder /app/configs ./configs

# 暴露端口
EXPOSE 58080

# 运行应用
CMD ["./main", "-conf", "configs/config.yaml"]