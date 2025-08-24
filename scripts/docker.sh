#!/bin/bash

# ================================================================
# ERP系统后端Docker容器管理脚本
# ================================================================
# 功能：后端容器的构建、启动、停止等操作
# 使用：./docker.sh {build|start|stop|restart|logs|shell}
# 作者：ERP系统开发组
# 版本：1.0.0

set -e

# ================================================================
# 配置参数
# ================================================================
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "${SCRIPT_DIR}/.." && pwd)"
BACKEND_DIR="${PROJECT_ROOT}/backend"
FRONTEND_DIR="${PROJECT_ROOT}/frontend"

# Docker 配置
COMPOSE_FILE="${PROJECT_ROOT}/docker-compose.yml"
ENV_FILE="${PROJECT_ROOT}/.env"
DOCKER_NETWORK="erp-network"

# 容器名称
POSTGRES_CONTAINER="erp-postgres"
REDIS_CONTAINER="erp-redis"
NGINX_CONTAINER="erp-nginx"
API_CONTAINER="erp-api"
FRONTEND_CONTAINER="erp-frontend"

# 端口配置
POSTGRES_PORT=58001
REDIS_PORT=58002
FRONTEND_PORT=58003
API_PORT=58080
NGINX_PORT=58000

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
PURPLE='\033[0;35m'
NC='\033[0m' # No Color

# ================================================================
# 辅助函数
# ================================================================

print_header() {
    echo -e "${BLUE}=================================="
    echo -e "ERP系统Docker容器管理"
    echo -e "==================================${NC}"
    echo ""
}

print_success() {
    echo -e "${GREEN}✓ $1${NC}"
}

print_error() {
    echo -e "${RED}✗ $1${NC}"
}

print_warning() {
    echo -e "${YELLOW}⚠ $1${NC}"
}

print_info() {
    echo -e "${BLUE}ℹ $1${NC}"
}

print_section() {
    echo -e "${PURPLE}>> $1${NC}"
}

# 检查Docker是否安装
check_docker() {
    if ! command -v docker &> /dev/null; then
        print_error "Docker 未安装，请先安装 Docker"
        print_info "安装地址: https://docs.docker.com/get-docker/"
        exit 1
    fi
    
    if ! command -v docker-compose &> /dev/null && ! docker compose version &> /dev/null; then
        print_error "Docker Compose 未安装"
        print_info "请安装 Docker Compose 或使用新版本的 Docker (内置 compose)"
        exit 1
    fi
    
    # 检查Docker守护进程是否运行
    if ! docker info &> /dev/null; then
        print_error "Docker 守护进程未运行，请启动 Docker"
        exit 1
    fi
}

# 获取compose命令
get_compose_cmd() {
    if docker compose version &> /dev/null; then
        echo "docker compose"
    else
        echo "docker-compose"
    fi
}

# 检查容器状态
check_container_status() {
    local container_name=$1
    if docker ps --format "table {{.Names}}" | grep -q "^${container_name}$"; then
        echo -e "${GREEN}● $container_name${NC} - 运行中"
        return 0
    elif docker ps -a --format "table {{.Names}}" | grep -q "^${container_name}$"; then
        echo -e "${YELLOW}● $container_name${NC} - 已停止"
        return 1
    else
        echo -e "${RED}○ $container_name${NC} - 不存在"
        return 2
    fi
}

# 获取容器PID
get_container_pid() {
    local container_name=$1
    local pid=$(docker inspect -f '{{.State.Pid}}' "$container_name" 2>/dev/null)
    if [ "$pid" != "0" ] && [ -n "$pid" ]; then
        echo "$pid"
    else
        echo "N/A"
    fi
}

# 获取容器端口映射
get_container_ports() {
    local container_name=$1
    local port_mapping=$(docker port "$container_name" 2>/dev/null | head -1)
    if [ -n "$port_mapping" ]; then
        # 提取端口号，去除 0.0.0.0: 前缀
        echo "$port_mapping" | sed 's/.*-> 0.0.0.0://' | sed 's/.*-> //'
    else
        echo "N/A"
    fi
}

# 显示服务详细状态和地址
show_service_details() {
    local container_name=$1
    local service_name=$2
    local default_port=$3
    local service_url_path=${4:-""}
    
    if docker ps --format "table {{.Names}}" | grep -q "^${container_name}$"; then
        local pid=$(get_container_pid "$container_name")
        local port=$(get_container_ports "$container_name")
        local status="${GREEN}运行中${NC}"
        
        # 构造服务访问地址
        local access_url=""
        case "$service_name" in
            "PostgreSQL")
                access_url="postgresql://localhost:${port}/erp_system"
                ;;
            "Redis")
                access_url="redis://localhost:${port}"
                ;;
            "API服务"|"后端API")
                access_url="http://localhost:${port}${service_url_path}"
                ;;
            "前端服务")
                access_url="http://localhost:${port}${service_url_path}"
                ;;
            "Nginx")
                access_url="http://localhost:${port}${service_url_path}"
                ;;
            *)
                access_url="localhost:${port}"
                ;;
        esac
        
        printf "%-12s ${status} | PID: %-8s | 端口: %-12s | 访问地址: %s\n" \
            "$service_name" "$pid" "$port" "$access_url"
    else
        local status="${RED}未运行${NC}"
        printf "%-12s ${status} | PID: %-8s | 端口: %-12s | 访问地址: %s\n" \
            "$service_name" "N/A" "N/A" "N/A"
    fi
}

# 显示所有服务详细信息
show_all_services() {
    print_section "服务状态详情"
    
    echo -e "${BLUE}服务名称     状态     | 进程PID    | 端口映射     | 访问地址${NC}"
    echo "=================================================================================="
    
    show_service_details "$POSTGRES_CONTAINER" "PostgreSQL" "$POSTGRES_PORT"
    show_service_details "$REDIS_CONTAINER" "Redis" "$REDIS_PORT"
    show_service_details "$API_CONTAINER" "后端API" "$API_PORT" "/api/v1"
    show_service_details "$FRONTEND_CONTAINER" "前端服务" "$FRONTEND_PORT" ""
    show_service_details "$NGINX_CONTAINER" "Nginx" "$NGINX_PORT" ""
}

# ================================================================
# Docker Compose文件创建
# ================================================================

create_docker_compose() {
    print_section "创建 Docker Compose 配置文件"
    
    cat > "$COMPOSE_FILE" << 'EOF'
version: '3.8'

services:
  # PostgreSQL 数据库
  postgres:
    container_name: erp-postgres
    image: postgres:14-alpine
    environment:
      POSTGRES_DB: ${POSTGRES_DB:-erp_system}
      POSTGRES_USER: ${POSTGRES_USER:-postgres}
      POSTGRES_PASSWORD: ${POSTGRES_PASSWORD:-postgres123}
      POSTGRES_INITDB_ARGS: "--encoding=UTF-8 --locale=C"
    ports:
      - "${POSTGRES_PORT:-5432}:5432"
    volumes:
      - postgres_data:/var/lib/postgresql/data
      - ./database/docker-init.sql:/docker-entrypoint-initdb.d/01-init.sql:ro
    networks:
      - erp-network
    restart: unless-stopped
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U ${POSTGRES_USER:-postgres} -d ${POSTGRES_DB:-erp_system}"]
      interval: 10s
      timeout: 5s
      retries: 5
      start_period: 30s

  # Redis 缓存
  redis:
    container_name: erp-redis
    image: redis:7-alpine
    command: redis-server --appendonly yes --requirepass ${REDIS_PASSWORD:-redis123}
    ports:
      - "${REDIS_PORT:-6379}:6379"
    volumes:
      - redis_data:/data
    networks:
      - erp-network
    restart: unless-stopped
    healthcheck:
      test: ["CMD", "redis-cli", "--raw", "incr", "ping"]
      interval: 10s
      timeout: 3s
      retries: 5
      start_period: 30s

  # Kratos API 后端服务
  api:
    container_name: erp-api
    build:
      context: .
      dockerfile: Dockerfile
      target: production
    environment:
      - APP_ENV=production
      - DATABASE_URL=postgresql://${POSTGRES_USER:-postgres}:${POSTGRES_PASSWORD:-postgres123}@postgres:5432/${POSTGRES_DB:-erp_system}?sslmode=disable
      - REDIS_URL=redis://:${REDIS_PASSWORD:-redis123}@redis:6379/0
      - JWT_SECRET=${JWT_SECRET:-your-super-secret-jwt-key-change-in-production}
      - API_PORT=${API_PORT:-58080}
    ports:
      - "${API_PORT:-58080}:58080"
    volumes:
      - ./configs:/app/configs:ro
      - api_logs:/app/logs
      - api_uploads:/app/uploads
    depends_on:
      postgres:
        condition: service_healthy
      redis:
        condition: service_healthy
    networks:
      - erp-network
    restart: unless-stopped
    healthcheck:
      test: ["CMD", "wget", "--quiet", "--tries=1", "--spider", "http://localhost:58080/health"]
      interval: 30s
      timeout: 10s
      retries: 3
      start_period: 40s

  # React 前端服务 (开发环境)
  frontend-dev:
    container_name: erp-frontend-dev
    build:
      context: ../frontend
      dockerfile: Dockerfile.dev
    environment:
      - REACT_APP_API_URL=http://localhost:${API_PORT:-58080}
      - CHOKIDAR_USEPOLLING=true
    ports:
      - "${FRONTEND_PORT:-3000}:3000"
    volumes:
      - ../frontend/src:/app/src:ro
      - ../frontend/public:/app/public:ro
      - frontend_node_modules:/app/node_modules
    depends_on:
      - api
    networks:
      - erp-network
    restart: unless-stopped
    profiles:
      - dev

  # React 前端服务 (生产环境)
  frontend:
    container_name: erp-frontend
    build:
      context: ../frontend
      dockerfile: Dockerfile
      target: production
    environment:
      - REACT_APP_API_URL=http://api.xx.com:${API_PORT:-58080}
    expose:
      - "80"
    depends_on:
      - api
    networks:
      - erp-network
    restart: unless-stopped
    profiles:
      - prod

  # Nginx 反向代理
  nginx:
    container_name: erp-nginx
    image: nginx:alpine
    ports:
      - "${NGINX_PORT:-80}:80"
      - "443:443"
    volumes:
      - ./docker/nginx/nginx.conf:/etc/nginx/nginx.conf:ro
      - ./docker/nginx/conf.d:/etc/nginx/conf.d:ro
      - ./docker/ssl:/etc/ssl/certs:ro
      - nginx_logs:/var/log/nginx
    depends_on:
      - api
      - frontend
    networks:
      - erp-network
    restart: unless-stopped
    profiles:
      - prod

volumes:
  postgres_data:
    driver: local
  redis_data:
    driver: local
  api_logs:
    driver: local
  api_uploads:
    driver: local
  frontend_node_modules:
    driver: local
  nginx_logs:
    driver: local

networks:
  erp-network:
    driver: bridge
    ipam:
      config:
        - subnet: 172.20.0.0/16
EOF

    print_success "Docker Compose 文件创建完成"
}

# ================================================================
# Dockerfile 创建
# ================================================================

create_backend_dockerfile() {
    print_section "创建后端 Dockerfile"
    
    cat > "$BACKEND_DIR/Dockerfile" << 'EOF'
# ================================================================
# 多阶段构建 Dockerfile - Kratos API 后端服务
# ================================================================

# 构建阶段
FROM golang:1.21-alpine AS builder

# 安装必要工具
RUN apk add --no-cache git ca-certificates tzdata

# 设置工作目录
WORKDIR /build

# 复制 go mod 文件
COPY go.mod go.sum ./

# 下载依赖（利用Docker缓存层）
RUN go mod download

# 复制源代码
COPY . .

# 构建应用
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-w -s" -o erp-api ./cmd/server

# 生产阶段
FROM alpine:latest AS production

# 安装必要的运行时依赖
RUN apk --no-cache add ca-certificates tzdata curl

# 创建非root用户
RUN addgroup -g 1001 app && \
    adduser -u 1001 -G app -s /bin/sh -D app

# 设置工作目录
WORKDIR /app

# 从构建阶段复制二进制文件
COPY --from=builder /build/erp-api .

# 复制配置文件和静态资源
COPY --chown=app:app configs ./configs
COPY --chown=app:app migrations ./migrations

# 创建必要的目录
RUN mkdir -p logs uploads tmp && \
    chown -R app:app logs uploads tmp

# 切换到非root用户
USER app

# 暴露端口
EXPOSE 58080

# 健康检查
HEALTHCHECK --interval=30s --timeout=10s --start-period=30s --retries=3 \
    CMD curl -f http://localhost:58080/health || exit 1

# 启动命令
CMD ["./erp-api", "-conf", "./configs/config.yaml"]
EOF

    print_success "后端 Dockerfile 创建完成"
}

create_frontend_dockerfile() {
    print_section "创建前端 Dockerfile"
    
    # 检查前端目录
    if [ ! -d "$FRONTEND_DIR" ]; then
        print_warning "前端目录不存在: $FRONTEND_DIR"
        print_info "跳过前端 Dockerfile 创建"
        return 0
    fi
    
    # 生产环境 Dockerfile
    cat > "$FRONTEND_DIR/Dockerfile" << 'EOF'
# ================================================================
# 多阶段构建 Dockerfile - React 前端应用
# ================================================================

# 构建阶段
FROM node:18-alpine AS builder

# 安装必要工具
RUN apk add --no-cache git

# 设置工作目录
WORKDIR /build

# 复制 package 文件
COPY package*.json ./

# 安装依赖
RUN npm ci --only=production

# 复制源代码
COPY . .

# 构建应用
RUN npm run build

# 生产阶段
FROM nginx:alpine AS production

# 复制自定义nginx配置
COPY docker/nginx.conf /etc/nginx/conf.d/default.conf

# 复制构建产物
COPY --from=builder /build/dist /usr/share/nginx/html

# 创建非root用户
RUN addgroup -g 1001 app && \
    adduser -u 1001 -G app -s /bin/sh -D app

# 暴露端口
EXPOSE 80

# 健康检查
HEALTHCHECK --interval=30s --timeout=3s --start-period=30s --retries=3 \
    CMD wget --quiet --tries=1 --spider http://localhost/ || exit 1

# 启动命令
CMD ["nginx", "-g", "daemon off;"]
EOF

    # 开发环境 Dockerfile
    cat > "$FRONTEND_DIR/Dockerfile.dev" << 'EOF'
# ================================================================
# 开发环境 Dockerfile - React 前端应用
# ================================================================

FROM node:18-alpine

# 安装必要工具
RUN apk add --no-cache git

# 设置工作目录
WORKDIR /app

# 复制 package 文件
COPY package*.json ./

# 安装依赖（包括开发依赖）
RUN npm install

# 复制源代码
COPY . .

# 暴露端口
EXPOSE 3000

# 健康检查
HEALTHCHECK --interval=30s --timeout=3s --start-period=30s --retries=3 \
    CMD wget --quiet --tries=1 --spider http://localhost:3000/ || exit 1

# 启动开发服务器
CMD ["npm", "run", "dev"]
EOF

    # 创建前端nginx配置
    mkdir -p "$FRONTEND_DIR/docker"
    cat > "$FRONTEND_DIR/docker/nginx.conf" << 'EOF'
server {
    listen 80;
    server_name localhost;
    root /usr/share/nginx/html;
    index index.html index.htm;

    # 启用gzip压缩
    gzip on;
    gzip_vary on;
    gzip_min_length 1024;
    gzip_types
        text/plain
        text/css
        text/xml
        text/javascript
        application/javascript
        application/xml+rss
        application/json;

    # 处理前端路由
    location / {
        try_files $uri $uri/ /index.html;
        add_header Cache-Control "no-cache, no-store, must-revalidate";
        add_header Pragma "no-cache";
        add_header Expires "0";
    }

    # 静态资源缓存
    location ~* \.(js|css|png|jpg|jpeg|gif|ico|svg)$ {
        expires 1y;
        add_header Cache-Control "public, immutable";
    }

    # API代理
    location /api/ {
        proxy_pass http://erp-api:58080/;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }

    # 健康检查
    location /health {
        access_log off;
        return 200 "healthy\n";
        add_header Content-Type text/plain;
    }
}
EOF

    print_success "前端 Dockerfile 创建完成"
}

# ================================================================
# Nginx 配置创建
# ================================================================

create_nginx_config() {
    print_section "创建 Nginx 配置文件"
    
    mkdir -p "$BACKEND_DIR/docker/nginx/conf.d"
    
    # 主配置文件
    cat > "$BACKEND_DIR/docker/nginx/nginx.conf" << 'EOF'
user nginx;
worker_processes auto;
error_log /var/log/nginx/error.log notice;
pid /var/run/nginx.pid;

events {
    worker_connections 1024;
}

http {
    include /etc/nginx/mime.types;
    default_type application/octet-stream;

    # 日志格式
    log_format main '$remote_addr - $remote_user [$time_local] "$request" '
                    '$status $body_bytes_sent "$http_referer" '
                    '"$http_user_agent" "$http_x_forwarded_for"';

    access_log /var/log/nginx/access.log main;

    # 基本设置
    sendfile on;
    tcp_nopush on;
    tcp_nodelay on;
    keepalive_timeout 65;
    types_hash_max_size 2048;

    # Gzip 压缩
    gzip on;
    gzip_vary on;
    gzip_min_length 1024;
    gzip_comp_level 6;
    gzip_types
        text/plain
        text/css
        text/xml
        text/javascript
        application/javascript
        application/xml+rss
        application/json
        application/xml
        image/svg+xml;

    # 包含站点配置
    include /etc/nginx/conf.d/*.conf;
}
EOF

    # 站点配置文件
    cat > "$BACKEND_DIR/docker/nginx/conf.d/default.conf" << 'EOF'
# ERP系统反向代理配置
upstream api_backend {
    server erp-api:58080;
    keepalive 32;
}

upstream frontend_backend {
    server erp-frontend:80;
    keepalive 32;
}

# HTTP 配置
server {
    listen 80;
    server_name api.xx.com localhost;

    # 安全头
    add_header X-Frame-Options "SAMEORIGIN" always;
    add_header X-XSS-Protection "1; mode=block" always;
    add_header X-Content-Type-Options "nosniff" always;
    add_header Referrer-Policy "no-referrer-when-downgrade" always;
    add_header Content-Security-Policy "default-src 'self' http: https: data: blob: 'unsafe-inline'" always;

    # API 服务代理
    location /v1/ {
        proxy_pass http://api_backend/v1/;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        proxy_cache_bypass $http_upgrade;
        
        # 超时设置
        proxy_connect_timeout 30s;
        proxy_send_timeout 30s;
        proxy_read_timeout 30s;
    }

    # 健康检查
    location /health {
        proxy_pass http://api_backend/health;
        proxy_set_header Host $host;
        access_log off;
    }

    # 前端静态文件
    location / {
        proxy_pass http://frontend_backend/;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        proxy_cache_bypass $http_upgrade;
    }

    # 文件上传大小限制
    client_max_body_size 50M;

    # 错误页面
    error_page 500 502 503 504 /50x.html;
    location = /50x.html {
        root /usr/share/nginx/html;
    }
}

# HTTPS 配置 (需要SSL证书)
# server {
#     listen 443 ssl http2;
#     server_name api.xx.com;
#
#     ssl_certificate /etc/ssl/certs/server.crt;
#     ssl_certificate_key /etc/ssl/certs/server.key;
#     ssl_session_timeout 1d;
#     ssl_session_cache shared:SSL:50m;
#     ssl_session_tickets off;
#
#     ssl_protocols TLSv1.2 TLSv1.3;
#     ssl_ciphers ECDHE-RSA-AES256-GCM-SHA512:DHE-RSA-AES256-GCM-SHA512:ECDHE-RSA-AES256-GCM-SHA384:DHE-RSA-AES256-GCM-SHA384;
#     ssl_prefer_server_ciphers off;
#
#     # HSTS
#     add_header Strict-Transport-Security "max-age=63072000" always;
#
#     # 其他配置与HTTP相同...
# }
EOF

    print_success "Nginx 配置文件创建完成"
}

# ================================================================
# Docker初始化
# ================================================================

docker_init() {
    print_section "初始化 Docker 环境"
    
    check_docker
    
    # 创建必要的目录结构
    mkdir -p "$BACKEND_DIR"/{logs,tmp,uploads,docker/nginx/{conf.d,ssl}}
    mkdir -p "$PROJECT_ROOT/volumes/postgres"
    mkdir -p "$PROJECT_ROOT/volumes/redis"
    
    # 创建 Docker 配置文件
    create_docker_compose
    create_backend_dockerfile
    create_frontend_dockerfile
    create_nginx_config
    
    # 创建 .dockerignore 文件
    cat > "$BACKEND_DIR/.dockerignore" << 'EOF'
# Git
.git
.gitignore

# Docker
Dockerfile*
docker-compose*.yml
.dockerignore

# Documentation
README.md
docs/
*.md

# Logs and temporary files
logs/
tmp/
*.log

# Development files
.air.toml
Makefile

# IDE
.vscode/
.idea/
*.swp
*.swo

# OS
.DS_Store
Thumbs.db

# Dependencies that will be downloaded
vendor/

# Test files
*_test.go
test/
coverage.out
EOF

    # 创建 Docker 网络
    if ! docker network ls | grep -q "$DOCKER_NETWORK"; then
        docker network create "$DOCKER_NETWORK"
        print_success "Docker 网络创建完成"
    fi
    
    print_success "Docker 环境初始化完成"
    print_info "配置文件已创建:"
    print_info "  - docker-compose.yml"
    print_info "  - Dockerfile (后端)"
    print_info "  - Dockerfile (前端)"
    print_info "  - Nginx 配置文件"
}

# ================================================================
# 容器管理
# ================================================================

docker_remove() {
    print_section "删除所有容器和数据"
    
    check_docker
    local compose_cmd=$(get_compose_cmd)
    
    print_warning "此操作将删除所有容器、镜像和数据卷！"
    echo -n "确定要继续吗? [y/N]: "
    read -r response
    
    if [[ ! "$response" =~ ^([yY][eE][sS]|[yY])$ ]]; then
        print_info "操作已取消"
        return 0
    fi
    
    cd "$PROJECT_ROOT"
    
    # 停止并删除容器
    $compose_cmd down -v --remove-orphans 2>/dev/null || true
    
    # 删除相关镜像
    docker images | grep -E "erp-|postgres|redis|nginx" | awk '{print $3}' | xargs -r docker rmi -f 2>/dev/null || true
    
    # 删除数据卷
    docker volume ls | grep -E "erp|postgres|redis" | awk '{print $2}' | xargs -r docker volume rm 2>/dev/null || true
    
    # 清理未使用的资源
    docker system prune -f
    
    print_success "所有容器和数据已删除"
}

docker_dev() {
    print_section "启动开发环境（不使用外部 Nginx）"
    
    check_docker
    local compose_cmd=$(get_compose_cmd)
    
    cd "$PROJECT_ROOT"
    
    # 检查环境文件
    if [ ! -f "$ENV_FILE" ]; then
        print_warning "环境配置文件不存在，使用默认配置"
    fi
    
    # 启动基础服务
    print_info "启动基础服务 (PostgreSQL, Redis)..."
    $compose_cmd up -d postgres redis
    
    # 等待数据库启动
    print_info "等待数据库服务启动..."
    sleep 10
    
    # 启动后端服务
    print_info "启动后端 API 服务..."
    $compose_cmd up -d api-server
    
    # 启动前端服务（开发模式）
    print_info "启动前端服务..."
    $compose_cmd up -d frontend
    
    # 等待服务完全启动
    print_info "等待服务完全启动..."
    sleep 8
    
    # 显示详细的服务状态信息
    show_all_services
    
    echo ""
    print_success "开发环境启动完成！"
    echo ""
    print_info "开发环境访问链接："
    print_info "  🌐 前端开发服务: http://localhost:$FRONTEND_PORT"
    print_info "  🚀 API服务: http://localhost:$API_PORT/v1"
    print_info "  💾 数据库: postgresql://localhost:$POSTGRES_PORT/erp_db"
    print_info "  📦 Redis: redis://localhost:$REDIS_PORT"
}

docker_prod() {
    print_section "启动生产环境（使用外部 Nginx）"
    
    check_docker
    local compose_cmd=$(get_compose_cmd)
    
    cd "$PROJECT_ROOT"
    
    # 检查环境文件
    if [ ! -f "$ENV_FILE" ]; then
        print_warning "环境配置文件不存在，使用默认配置"
    fi
    
    # 检查nginx配置文件
    if [ ! -f "$PROJECT_ROOT/docker/nginx/default.conf" ]; then
        print_error "docker/nginx/default.conf 配置文件不存在，无法启动生产环境"
        print_info "请先创建 docker/nginx/default.conf 配置文件"
        return 1
    fi
    
    # 启动基础服务
    print_info "启动基础服务 (PostgreSQL, Redis)..."
    $compose_cmd up -d postgres redis
    
    # 等待数据库启动
    print_info "等待数据库服务启动..."
    sleep 10
    
    # 启动后端服务
    print_info "启动后端 API 服务..."
    $compose_cmd up -d api-server
    
    # 启动前端服务
    print_info "启动前端服务..."
    $compose_cmd up -d frontend
    
    # 启动Nginx代理服务
    print_info "启动 Nginx 代理服务..."
    $compose_cmd up -d nginx
    
    # 等待服务完全启动
    print_info "等待服务完全启动..."
    sleep 8
    
    # 显示详细的服务状态信息
    show_all_services
    
    echo ""
    print_success "生产环境启动完成！"
    echo ""
    print_info "生产环境访问链接："
    print_info "  🌐 前端应用 (通过Nginx): http://localhost:$NGINX_PORT"
    print_info "  🚀 API服务: http://localhost:$API_PORT/v1"
    print_info "  💾 数据库: postgresql://localhost:$POSTGRES_PORT/erp_db"
    print_info "  📦 Redis: redis://localhost:$REDIS_PORT"
}

docker_start() {
    print_section "启动 Docker 容器"
    
    check_docker
    local compose_cmd=$(get_compose_cmd)
    
    cd "$PROJECT_ROOT"
    
    # 检查环境文件
    if [ ! -f "$ENV_FILE" ]; then
        print_warning "环境配置文件不存在，使用默认配置"
    fi
    
    # 启动服务
    print_info "启动基础服务 (PostgreSQL, Redis)..."
    $compose_cmd up -d postgres redis
    
    # 等待数据库启动
    print_info "等待数据库服务启动..."
    sleep 10
    
    # 启动应用服务
    print_info "启动应用服务..."
    $compose_cmd up -d api-server
    
    # 启动前端服务
    print_info "启动前端服务..."
    $compose_cmd up -d frontend
    
    # 启动Nginx代理服务（如果使用生产配置且配置文件存在）
    if $compose_cmd config --services | grep -q nginx && [ -f "$PROJECT_ROOT/nginx.conf" ]; then
        print_info "启动Nginx代理服务..."
        $compose_cmd up -d nginx
    elif $compose_cmd config --services | grep -q nginx; then
        print_warning "跳过Nginx服务启动 - nginx.conf 配置文件不存在"
    fi
    
    # 等待服务完全启动
    print_info "等待服务完全启动..."
    sleep 8
    
    # 显示详细的服务状态信息
    show_all_services
    
    echo ""
    print_success "所有容器启动完成！"
    echo ""
    print_info "快速访问链接："
    print_info "  🌐 前端应用: http://localhost:$FRONTEND_PORT"
    print_info "  🚀 API服务: http://localhost:$API_PORT/api/v1"
    print_info "  💾 数据库: postgresql://localhost:$POSTGRES_PORT/erp_system"
    print_info "  📦 Redis: redis://localhost:$REDIS_PORT"
}

docker_stop() {
    print_section "停止 Docker 容器"
    
    check_docker
    local compose_cmd=$(get_compose_cmd)
    
    cd "$PROJECT_ROOT"
    
    # 在停止前显示当前运行的服务状态
    print_info "停止前服务状态："
    show_all_services
    
    echo ""
    print_info "正在停止所有容器..."
    
    # 逐个停止服务以便观察
    if docker ps --format "table {{.Names}}" | grep -q "^${NGINX_CONTAINER}$"; then
        print_info "停止 Nginx 服务..."
        $compose_cmd stop nginx
    fi
    
    if docker ps --format "table {{.Names}}" | grep -q "^${FRONTEND_CONTAINER}$"; then
        print_info "停止前端服务..."
        $compose_cmd stop frontend
    fi
    
    if docker ps --format "table {{.Names}}" | grep -q "^${API_CONTAINER}$"; then
        print_info "停止API服务..."
        $compose_cmd stop api-server
    fi
    
    if docker ps --format "table {{.Names}}" | grep -q "^${REDIS_CONTAINER}$"; then
        print_info "停止Redis服务..."
        $compose_cmd stop redis
    fi
    
    if docker ps --format "table {{.Names}}" | grep -q "^${POSTGRES_CONTAINER}$"; then
        print_info "停止PostgreSQL服务..."
        $compose_cmd stop postgres
    fi
    
    # 确保所有容器都已停止
    $compose_cmd down
    
    echo ""
    print_success "所有容器已停止"
    
    # 显示停止后的状态
    echo ""
    print_info "停止后服务状态："
    show_all_services
}

docker_restart() {
    print_section "重启 Docker 容器"
    
    docker_stop
    sleep 3
    docker_start
}

docker_status() {
    print_section "检查容器状态"
    
    check_docker
    
    # 显示详细的服务状态
    show_all_services
    
    echo ""
    echo -e "${BLUE}网络和端口信息:${NC}"
    echo "============================================================================="
    
    # 显示Docker网络状态
    if docker network ls | grep -q "$DOCKER_NETWORK"; then
        echo -e "Docker网络   ${GREEN}● $DOCKER_NETWORK${NC} - 已创建"
    else
        echo -e "Docker网络   ${RED}○ $DOCKER_NETWORK${NC} - 不存在"
    fi
    
    # 显示端口配置
    echo ""
    echo "配置的端口映射："
    echo "  PostgreSQL: $POSTGRES_PORT"
    echo "  Redis:      $REDIS_PORT"
    echo "  API:        $API_PORT"
    echo "  Frontend:   $FRONTEND_PORT"
    echo "  Nginx:      $NGINX_PORT"
    
    # 检查端口是否被占用
    echo ""
    echo "端口占用检查："
    for port in $POSTGRES_PORT $REDIS_PORT $API_PORT $FRONTEND_PORT $NGINX_PORT; do
        if lsof -i ":$port" >/dev/null 2>&1; then
            echo -e "  端口 $port: ${YELLOW}占用中${NC}"
        else
            echo -e "  端口 $port: ${GREEN}可用${NC}"
        fi
    done
}

docker_logs() {
    print_section "查看容器日志"
    
    check_docker
    local compose_cmd=$(get_compose_cmd)
    
    cd "$PROJECT_ROOT"
    
    echo "选择要查看的日志:"
    echo "1) API 服务日志"
    echo "2) 前端服务日志"
    echo "3) PostgreSQL 日志"
    echo "4) Redis 日志"
    echo "5) Nginx 日志"
    echo "6) 所有服务日志"
    echo -n "请选择 [1-6]: "
    
    read -r choice
    
    case $choice in
        1) $compose_cmd logs -f api-server ;;
        2) $compose_cmd logs -f frontend ;;
        3) $compose_cmd logs -f postgres ;;
        4) $compose_cmd logs -f redis ;;
        5) $compose_cmd logs -f nginx ;;
        6) $compose_cmd logs -f ;;
        *) print_error "无效的选择" ;;
    esac
}

docker_ps() {
    print_section "显示容器进程"
    
    check_docker
    
    echo "运行中的容器："
    docker ps --format "table {{.Names}}\t{{.Status}}\t{{.Ports}}\t{{.Image}}"
    
    echo ""
    echo "所有容器："
    docker ps -a --format "table {{.Names}}\t{{.Status}}\t{{.Ports}}\t{{.Image}}"
}

# ================================================================
# 主程序
# ================================================================

main() {
    case "${1:-}" in
        init)
            print_header
            docker_init
            ;;
        rm)
            print_header
            docker_remove
            ;;
        dev)
            print_header
            docker_dev
            ;;
        prod)
            print_header
            docker_prod
            ;;
        start)
            print_header
            docker_start
            ;;
        stop)
            print_header
            docker_stop
            ;;
        restart)
            print_header
            docker_restart
            ;;
        status)
            print_header
            docker_status
            ;;
        info)
            print_header
            show_all_services
            ;;
        logs)
            print_header
            docker_logs
            ;;
        ps)
            print_header
            docker_ps
            ;;
        *)
            print_header
            echo "用法: $0 {init|rm|dev|prod|start|stop|restart|status|info|logs|ps}"
            echo ""
            echo "初始化管理:"
            echo "  init            - 初始化Docker环境（创建配置文件）"
            echo "  rm              - 删除所有容器和数据"
            echo ""
            echo "环境管理:"
            echo "  dev             - 启动开发环境（不使用外部Nginx）"
            echo "  prod            - 启动生产环境（使用外部Nginx）"
            echo ""
            echo "容器管理:"
            echo "  start           - 启动所有容器（显示PID和服务地址）"
            echo "  stop            - 停止所有容器（显示PID和服务地址）"
            echo "  restart         - 重启所有容器"
            echo "  status          - 查看完整容器状态和端口信息"
            echo "  info            - 显示服务详细信息（PID、端口、访问地址）"
            echo ""
            echo "监控工具:"
            echo "  logs            - 查看容器日志"
            echo "  ps              - 显示容器进程"
            echo ""
            echo "示例:"
            echo "  $0 init           # 初始化Docker环境"
            echo "  $0 dev            # 启动开发环境"
            echo "  $0 prod           # 启动生产环境"
            echo "  $0 start          # 启动系统（显示服务PID和访问地址）"
            echo "  $0 info           # 查看所有服务详细状态"
            echo "  $0 status         # 查看完整状态和端口信息"
            echo "  $0 stop           # 停止系统（显示停止前后状态）"
            echo "  $0 logs           # 查看日志"
            exit 1
            ;;
    esac
}

# 执行主程序
main "$@"