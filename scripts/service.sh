#!/bin/bash

# ================================================================
# ERP系统后端服务管理脚本 
# ================================================================
# 功能：启动/停止/重启后端服务
# 使用：./service.sh {start|stop|restart|status|logs}
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

# 服务配置
POSTGRES_PORT=5432
REDIS_PORT=6379
API_PORT=58080
FRONTEND_PORT=58000

# 服务名称
POSTGRES_SERVICE="postgresql"
REDIS_SERVICE="redis-server"
API_SERVICE="erp-api"
FRONTEND_SERVICE="erp-frontend"

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
    echo -e "ERP系统服务管理"
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

# 检查端口是否被占用
check_port() {
    local port=$1
    local service_name=$2
    
    if lsof -Pi :$port -sTCP:LISTEN -t >/dev/null 2>&1; then
        local pid=$(lsof -Pi :$port -sTCP:LISTEN -t)
        local process_name=$(ps -p $pid -o comm= 2>/dev/null || echo "unknown")
        print_warning "端口 $port 已被进程 $process_name (PID: $pid) 占用"
        
        echo -n "是否要终止该进程并继续? [y/N]: "
        read -r response
        if [[ "$response" =~ ^([yY][eE][sS]|[yY])$ ]]; then
            kill -9 $pid 2>/dev/null || true
            sleep 1
            print_success "已终止进程 $pid"
            return 0
        else
            print_error "无法启动 $service_name，端口 $port 被占用"
            return 1
        fi
    fi
    return 0
}

# 检查服务状态
check_service_status() {
    local service_name=$1
    local port=$2
    
    if lsof -Pi :$port -sTCP:LISTEN -t >/dev/null 2>&1; then
        local pid=$(lsof -Pi :$port -sTCP:LISTEN -t)
        echo -e "${GREEN}● $service_name${NC} - 运行中 (PID: $pid, Port: $port)"
        return 0
    else
        echo -e "${RED}○ $service_name${NC} - 未运行 (Port: $port)"
        return 1
    fi
}

# 等待服务启动
wait_for_service() {
    local port=$1
    local service_name=$2
    local timeout=${3:-30}
    
    print_info "等待 $service_name 服务启动..."
    
    for i in $(seq 1 $timeout); do
        if lsof -Pi :$port -sTCP:LISTEN -t >/dev/null 2>&1; then
            print_success "$service_name 服务已启动 (端口: $port)"
            return 0
        fi
        sleep 1
        echo -n "."
    done
    
    echo ""
    print_error "$service_name 服务启动超时"
    return 1
}

# ================================================================
# PostgreSQL 服务管理
# ================================================================

start_postgres() {
    print_section "启动 PostgreSQL 数据库服务"
    
    # 检查是否已安装
    if ! command -v psql &> /dev/null; then
        print_error "PostgreSQL 未安装，请先运行: ./service.sh install"
        return 1
    fi
    
    # 检查端口
    if ! check_port $POSTGRES_PORT "PostgreSQL"; then
        return 1
    fi
    
    # macOS 系统
    if [[ "$OSTYPE" == "darwin"* ]]; then
        # 使用 Homebrew 启动
        if command -v brew &> /dev/null; then
            brew services start postgresql@14 2>/dev/null || brew services start postgresql
            if wait_for_service $POSTGRES_PORT "PostgreSQL" 10; then
                return 0
            fi
        fi
        
        # 手动启动
        if [ -d "/usr/local/var/postgresql@14" ]; then
            pg_ctl -D /usr/local/var/postgresql@14 start
        elif [ -d "/usr/local/var/postgres" ]; then
            pg_ctl -D /usr/local/var/postgres start
        elif [ -d "/opt/homebrew/var/postgresql@14" ]; then
            pg_ctl -D /opt/homebrew/var/postgresql@14 start
        elif [ -d "/opt/homebrew/var/postgres" ]; then
            pg_ctl -D /opt/homebrew/var/postgres start
        fi
        
    # Linux 系统
    else
        sudo systemctl start postgresql
    fi
    
    wait_for_service $POSTGRES_PORT "PostgreSQL"
}

stop_postgres() {
    print_section "停止 PostgreSQL 数据库服务"
    
    # macOS 系统
    if [[ "$OSTYPE" == "darwin"* ]]; then
        if command -v brew &> /dev/null; then
            brew services stop postgresql@14 2>/dev/null || brew services stop postgresql
        fi
    # Linux 系统
    else
        sudo systemctl stop postgresql
    fi
    
    # 强制终止进程
    pkill -f postgres 2>/dev/null || true
    
    print_success "PostgreSQL 服务已停止"
}

# ================================================================
# Redis 服务管理
# ================================================================

start_redis() {
    print_section "启动 Redis 缓存服务"
    
    # 检查是否已安装
    if ! command -v redis-server &> /dev/null; then
        print_error "Redis 未安装，请先运行: ./service.sh install"
        return 1
    fi
    
    # 检查端口
    if ! check_port $REDIS_PORT "Redis"; then
        return 1
    fi
    
    # macOS 系统
    if [[ "$OSTYPE" == "darwin"* ]]; then
        if command -v brew &> /dev/null; then
            brew services start redis
        else
            nohup redis-server > /tmp/redis.log 2>&1 &
        fi
    # Linux 系统
    else
        sudo systemctl start redis-server
    fi
    
    wait_for_service $REDIS_PORT "Redis"
}

stop_redis() {
    print_section "停止 Redis 缓存服务"
    
    # macOS 系统
    if [[ "$OSTYPE" == "darwin"* ]]; then
        if command -v brew &> /dev/null; then
            brew services stop redis
        fi
    # Linux 系统
    else
        sudo systemctl stop redis-server
    fi
    
    # 强制终止进程
    pkill -f redis-server 2>/dev/null || true
    
    print_success "Redis 服务已停止"
}

# ================================================================
# API 服务管理
# ================================================================

start_api() {
    print_section "启动 Kratos API 服务"
    
    # 检查端口
    if ! check_port $API_PORT "API Server"; then
        return 1
    fi
    
    cd "$BACKEND_DIR"
    
    # 检查是否已编译
    if [ ! -f "./bin/erp-api" ]; then
        print_info "未找到编译的 API 服务，正在构建..."
        if ! build_api; then
            print_error "API 服务构建失败"
            return 1
        fi
    fi
    
    # 启动服务
    nohup ./bin/erp-api -conf ./configs/config.yaml > ./logs/api.log 2>&1 &
    local pid=$!
    echo $pid > ./tmp/api.pid
    
    wait_for_service $API_PORT "API Server"
}

stop_api() {
    print_section "停止 Kratos API 服务"
    
    cd "$BACKEND_DIR"
    
    # 从 PID 文件停止
    if [ -f "./tmp/api.pid" ]; then
        local pid=$(cat ./tmp/api.pid)
        if kill -0 $pid 2>/dev/null; then
            kill -TERM $pid
            sleep 2
            if kill -0 $pid 2>/dev/null; then
                kill -9 $pid
            fi
        fi
        rm -f ./tmp/api.pid
    fi
    
    # 强制终止进程
    pkill -f erp-api 2>/dev/null || true
    
    print_success "API 服务已停止"
}

# ================================================================
# 前端服务管理
# ================================================================

start_frontend() {
    print_section "启动 React 前端服务"
    
    # 检查前端目录
    if [ ! -d "$FRONTEND_DIR" ]; then
        print_warning "前端目录不存在: $FRONTEND_DIR"
        print_info "跳过前端服务启动"
        return 0
    fi
    
    # 检查端口
    if ! check_port $FRONTEND_PORT "Frontend Server"; then
        return 1
    fi
    
    cd "$FRONTEND_DIR"
    
    # 检查依赖
    if [ ! -d "node_modules" ]; then
        print_info "正在安装前端依赖..."
        npm install || {
            print_error "前端依赖安装失败"
            return 1
        }
    fi
    
    # 启动开发服务器
    nohup npm run dev > "$BACKEND_DIR/logs/frontend.log" 2>&1 &
    local pid=$!
    echo $pid > "$BACKEND_DIR/tmp/frontend.pid"
    
    wait_for_service $FRONTEND_PORT "Frontend Server"
}

stop_frontend() {
    print_section "停止 React 前端服务"
    
    cd "$BACKEND_DIR"
    
    # 从 PID 文件停止
    if [ -f "./tmp/frontend.pid" ]; then
        local pid=$(cat ./tmp/frontend.pid)
        if kill -0 $pid 2>/dev/null; then
            kill -TERM $pid
            sleep 2
            if kill -0 $pid 2>/dev/null; then
                kill -9 $pid
            fi
        fi
        rm -f ./tmp/frontend.pid
    fi
    
    # 强制终止进程
    pkill -f "npm.*dev" 2>/dev/null || true
    pkill -f "vite" 2>/dev/null || true
    
    print_success "前端服务已停止"
}

# ================================================================
# 构建服务
# ================================================================

# 构建后端API服务
build_api() {
    print_section "构建 Kratos API 服务"
    
    cd "$BACKEND_DIR"
    
    # 创建必要的目录
    mkdir -p bin logs tmp
    
    # 检查Go环境
    if ! command -v go &> /dev/null; then
        print_error "Go 未安装，请先运行: ./service.sh install"
        return 1
    fi
    
    print_info "正在下载依赖..."
    go mod download
    
    print_info "正在编译 API 服务..."
    if go build -o bin/erp-api cmd/server/main.go; then
        print_success "API 服务编译成功"
        chmod +x bin/erp-api
        return 0
    else
        print_error "API 服务编译失败"
        return 1
    fi
}

# 构建前端服务
build_frontend() {
    print_section "构建 React 前端服务"
    
    # 检查前端目录
    if [ ! -d "$FRONTEND_DIR" ]; then
        print_warning "前端目录不存在: $FRONTEND_DIR"
        print_info "跳过前端服务构建"
        return 0
    fi
    
    cd "$FRONTEND_DIR"
    
    # 检查Node.js环境
    if ! command -v node &> /dev/null; then
        print_error "Node.js 未安装，请先运行: ./service.sh install"
        return 1
    fi
    
    # 检查并安装依赖
    if [ ! -d "node_modules" ]; then
        print_info "正在安装前端依赖..."
        if ! npm install; then
            print_error "前端依赖安装失败"
            return 1
        fi
    fi
    
    print_info "正在构建前端应用..."
    if npm run build; then
        print_success "前端构建成功"
        print_info "构建文件位于: $FRONTEND_DIR/dist"
        return 0
    else
        print_error "前端构建失败，请修复 TypeScript 错误后重新构建"
        return 1
    fi
}

# 构建所有服务
build_all() {
    print_header
    print_info "构建所有服务..."
    echo ""
    
    build_api || return 1
    build_frontend || return 1
    
    echo ""
    print_success "所有服务构建完成！"
    print_info "可以运行 ./service.sh start 启动服务"
}

# ================================================================
# 服务安装
# ================================================================

install_services() {
    print_section "安装系统依赖"
    
    # macOS 系统
    if [[ "$OSTYPE" == "darwin"* ]]; then
        # 检查 Homebrew
        if ! command -v brew &> /dev/null; then
            print_error "请先安装 Homebrew: https://brew.sh/"
            return 1
        fi
        
        print_info "安装 PostgreSQL..."
        brew install postgresql@14 || brew install postgresql
        
        print_info "安装 Redis..."
        brew install redis
        
        print_info "安装 Go..."
        brew install go
        
        print_info "安装 Node.js..."
        brew install node
        
    # Linux 系统 (Ubuntu/Debian)
    elif command -v apt-get &> /dev/null; then
        print_info "更新软件包列表..."
        sudo apt-get update
        
        print_info "安装 PostgreSQL..."
        sudo apt-get install -y postgresql postgresql-contrib
        
        print_info "安装 Redis..."
        sudo apt-get install -y redis-server
        
        print_info "安装 Go..."
        sudo apt-get install -y golang-go
        
        print_info "安装 Node.js..."
        curl -fsSL https://deb.nodesource.com/setup_18.x | sudo -E bash -
        sudo apt-get install -y nodejs
        
    # Linux 系统 (RHEL/CentOS)
    elif command -v yum &> /dev/null; then
        print_info "安装 PostgreSQL..."
        sudo yum install -y postgresql-server postgresql-contrib
        sudo postgresql-setup initdb
        
        print_info "安装 Redis..."
        sudo yum install -y redis
        
        print_info "安装 Go..."
        sudo yum install -y golang
        
        print_info "安装 Node.js..."
        curl -fsSL https://rpm.nodesource.com/setup_18.x | sudo bash -
        sudo yum install -y nodejs
        
    else
        print_error "不支持的操作系统"
        return 1
    fi
    
    print_success "依赖安装完成"
    
    # 创建必要的目录
    mkdir -p "$BACKEND_DIR"/{logs,tmp,bin}
    
    print_info "请运行以下命令初始化数据库:"
    echo "  createdb erp_system"
    echo "  psql -U postgres -d erp_system -f database/setup_database.sql"
}

# ================================================================
# 日志查看
# ================================================================

show_logs() {
    print_section "查看服务日志"
    
    echo "选择要查看的日志:"
    echo "1) API 服务日志"
    echo "2) 前端服务日志"
    echo "3) PostgreSQL 日志"
    echo "4) Redis 日志"
    echo "5) 所有日志"
    echo -n "请选择 [1-5]: "
    
    read -r choice
    
    cd "$BACKEND_DIR"
    
    case $choice in
        1)
            if [ -f "./logs/api.log" ]; then
                tail -f ./logs/api.log
            else
                print_error "API 日志文件不存在"
            fi
            ;;
        2)
            if [ -f "./logs/frontend.log" ]; then
                tail -f ./logs/frontend.log
            else
                print_error "前端日志文件不存在"
            fi
            ;;
        3)
            # PostgreSQL 日志位置因系统而异
            if [[ "$OSTYPE" == "darwin"* ]]; then
                if [ -f "/usr/local/var/log/postgresql@14.log" ]; then
                    tail -f /usr/local/var/log/postgresql@14.log
                elif [ -f "/usr/local/var/log/postgres.log" ]; then
                    tail -f /usr/local/var/log/postgres.log
                elif [ -f "/opt/homebrew/var/log/postgresql@14.log" ]; then
                    tail -f /opt/homebrew/var/log/postgresql@14.log
                else
                    print_error "PostgreSQL 日志文件未找到"
                fi
            else
                sudo journalctl -u postgresql -f
            fi
            ;;
        4)
            if [[ "$OSTYPE" == "darwin"* ]]; then
                if [ -f "/usr/local/var/log/redis.log" ]; then
                    tail -f /usr/local/var/log/redis.log
                elif [ -f "/tmp/redis.log" ]; then
                    tail -f /tmp/redis.log
                else
                    print_error "Redis 日志文件未找到"
                fi
            else
                sudo journalctl -u redis -f
            fi
            ;;
        5)
            print_info "显示所有可用日志..."
            ls -la ./logs/ 2>/dev/null || print_warning "日志目录为空"
            ;;
        *)
            print_error "无效的选择"
            ;;
    esac
}

# ================================================================
# 服务状态检查
# ================================================================

check_status() {
    print_section "检查服务状态"
    
    check_service_status "PostgreSQL" $POSTGRES_PORT
    check_service_status "Redis" $REDIS_PORT  
    check_service_status "API Server" $API_PORT
    check_service_status "Frontend Server" $FRONTEND_PORT
    
    echo ""
    print_info "系统信息:"
    echo "  操作系统: $(uname -s)"
    echo "  架构: $(uname -m)"
    echo "  Go 版本: $(go version 2>/dev/null | cut -d' ' -f3 || echo '未安装')"
    echo "  Node 版本: $(node --version 2>/dev/null || echo '未安装')"
    echo "  PostgreSQL: $(psql --version 2>/dev/null | cut -d' ' -f3 || echo '未安装')"
    echo "  Redis: $(redis-server --version 2>/dev/null | cut -d' ' -f1,3 || echo '未安装')"
}

# ================================================================
# 详细服务信息显示 (类似 docker.sh info)
# ================================================================

# 显示单个服务详细信息
show_service_details() {
    local service_name=$1
    local port=$2
    local service_url_path=${3:-""}
    
    if lsof -Pi :$port -sTCP:LISTEN -t >/dev/null 2>&1; then
        local pid=$(lsof -Pi :$port -sTCP:LISTEN -t)
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
            "后端API")
                access_url="http://localhost:${port}${service_url_path}"
                ;;
            "前端服务")
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
    
    show_service_details "PostgreSQL" $POSTGRES_PORT
    show_service_details "Redis" $REDIS_PORT
    show_service_details "后端API" $API_PORT "/api/v1"
    show_service_details "前端服务" $FRONTEND_PORT ""
    
    echo ""
    print_section "快速访问链接"
    echo -e "${BLUE}应用访问:${NC}"
    echo "  前端管理界面: http://localhost:$FRONTEND_PORT"
    echo "  API 接口文档: http://localhost:$API_PORT/api/v1/docs"
    echo "  API 健康检查: http://localhost:$API_PORT/api/v1/health"
    echo ""
    echo -e "${BLUE}数据库连接:${NC}"
    echo "  PostgreSQL: postgresql://localhost:$POSTGRES_PORT/erp_system"
    echo "  Redis: redis://localhost:$REDIS_PORT"
    echo ""
    echo -e "${BLUE}系统信息:${NC}"
    echo "  操作系统: $(uname -s) $(uname -m)"
    echo "  Go 版本: $(go version 2>/dev/null | cut -d' ' -f3 || echo '未安装')"
    echo "  Node 版本: $(node --version 2>/dev/null || echo '未安装')"
    echo "  时间: $(date '+%Y-%m-%d %H:%M:%S')"
}

# ================================================================
# 主要服务控制函数
# ================================================================

start_all() {
    print_header
    print_info "启动所有服务..."
    echo ""
    
    # 创建必要目录
    mkdir -p "$BACKEND_DIR"/{logs,tmp}
    
    # 按依赖顺序启动
    start_postgres || return 1
    start_redis || return 1
    
    # 等待数据库服务完全启动
    sleep 2
    
    start_api || return 1
    start_frontend || return 1
    
    echo ""
    print_success "所有服务启动完成！"
    print_info "API 服务: http://localhost:$API_PORT"
    print_info "前端服务: http://localhost:$FRONTEND_PORT"
}

stop_all() {
    print_header
    print_info "停止所有服务..."
    echo ""
    
    stop_frontend
    stop_api
    stop_redis
    stop_postgres
    
    echo ""
    print_success "所有服务已停止"
}

restart_all() {
    print_header
    print_info "重启所有服务..."
    echo ""
    
    stop_all
    sleep 3
    start_all
}

# ================================================================
# 主程序
# ================================================================

main() {
    case "${1:-}" in
        start)
            start_all
            ;;
        stop)
            stop_all
            ;;
        restart)
            restart_all
            ;;
        status)
            print_header
            check_status
            ;;
        info)
            print_header
            show_all_services
            ;;
        build)
            build_all
            ;;
        install)
            print_header
            install_services
            ;;
        logs)
            print_header
            show_logs
            ;;
        postgres)
            case "${2:-}" in
                start) start_postgres ;;
                stop) stop_postgres ;;
                restart) stop_postgres; sleep 2; start_postgres ;;
                *) print_error "用法: $0 postgres {start|stop|restart}" ;;
            esac
            ;;
        redis)
            case "${2:-}" in
                start) start_redis ;;
                stop) stop_redis ;;
                restart) stop_redis; sleep 2; start_redis ;;
                *) print_error "用法: $0 redis {start|stop|restart}" ;;
            esac
            ;;
        api)
            case "${2:-}" in
                start) start_api ;;
                stop) stop_api ;;
                restart) stop_api; sleep 2; start_api ;;
                build) build_api ;;
                *) print_error "用法: $0 api {start|stop|restart|build}" ;;
            esac
            ;;
        frontend)
            case "${2:-}" in
                start) start_frontend ;;
                stop) stop_frontend ;;
                restart) stop_frontend; sleep 2; start_frontend ;;
                build) build_frontend ;;
                *) print_error "用法: $0 frontend {start|stop|restart|build}" ;;
            esac
            ;;
        *)
            print_header
            echo "用法: $0 {start|stop|restart|status|info|build|install|logs}"
            echo ""
            echo "服务管理:"
            echo "  start     - 启动所有服务"
            echo "  stop      - 停止所有服务" 
            echo "  restart   - 重启所有服务"
            echo "  status    - 查看服务状态"
            echo "  info      - 显示详细服务信息和访问地址"
            echo "  build     - 构建所有服务"
            echo ""
            echo "单独管理:"
            echo "  postgres {start|stop|restart}        - PostgreSQL 数据库"
            echo "  redis {start|stop|restart}           - Redis 缓存"
            echo "  api {start|stop|restart|build}       - Kratos API 服务"
            echo "  frontend {start|stop|restart|build}  - React 前端"
            echo ""
            echo "工具功能:"
            echo "  install   - 安装系统依赖"
            echo "  logs      - 查看服务日志"
            echo ""
            echo "示例:"
            echo "  $0 build          # 构建所有服务"
            echo "  $0 start          # 启动所有服务"
            echo "  $0 api build      # 只构建API服务"
            echo "  $0 postgres start # 只启动数据库"
            echo "  $0 status         # 查看状态"
            echo "  $0 info           # 查看详细服务信息"
            exit 1
            ;;
    esac
}

# 执行主程序
main "$@"