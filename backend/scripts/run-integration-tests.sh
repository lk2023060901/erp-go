#!/bin/bash

# 集成测试运行脚本
set -e

echo "🚀 启动 Frappe 权限系统集成测试..."

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 检查Docker是否运行
check_docker() {
    if ! docker info > /dev/null 2>&1; then
        echo -e "${RED}❌ Docker 未运行，请先启动 Docker${NC}"
        exit 1
    fi
}

# 清理旧容器和数据
cleanup() {
    echo -e "${YELLOW}🧹 清理旧的测试环境...${NC}"
    docker-compose -f docker-compose.test.yml down -v --remove-orphans || true
    docker system prune -f > /dev/null 2>&1 || true
}

# 启动测试环境
start_test_env() {
    echo -e "${BLUE}🔧 启动测试环境...${NC}"
    docker-compose -f docker-compose.test.yml up -d postgres-test redis-test
    
    echo -e "${BLUE}⏳ 等待数据库就绪...${NC}"
    timeout=60
    while [ $timeout -gt 0 ]; do
        if docker-compose -f docker-compose.test.yml exec -T postgres-test pg_isready -U erp_user -d erp_test > /dev/null 2>&1; then
            echo -e "${GREEN}✅ PostgreSQL 已就绪${NC}"
            break
        fi
        sleep 2
        timeout=$((timeout-2))
    done
    
    if [ $timeout -le 0 ]; then
        echo -e "${RED}❌ PostgreSQL 启动超时${NC}"
        exit 1
    fi
    
    echo -e "${BLUE}⏳ 等待 Redis 就绪...${NC}"
    timeout=30
    while [ $timeout -gt 0 ]; do
        if docker-compose -f docker-compose.test.yml exec -T redis-test redis-cli ping > /dev/null 2>&1; then
            echo -e "${GREEN}✅ Redis 已就绪${NC}"
            break
        fi
        sleep 2
        timeout=$((timeout-2))
    done
    
    if [ $timeout -le 0 ]; then
        echo -e "${RED}❌ Redis 启动超时${NC}"
        exit 1
    fi
}

# 运行单元测试
run_unit_tests() {
    echo -e "${BLUE}🧪 运行单元测试...${NC}"
    
    # Data层测试
    echo -e "${YELLOW}📊 Testing Data Layer...${NC}"
    go test -v ./internal/data -timeout=30s
    
    # Business层测试
    echo -e "${YELLOW}💼 Testing Business Layer...${NC}"
    go test -v ./internal/biz -timeout=30s
    
    # Service层测试
    echo -e "${YELLOW}⚙️ Testing Service Layer...${NC}"
    go test -v ./internal/service -timeout=30s
    
    echo -e "${GREEN}✅ 单元测试通过${NC}"
}

# 运行集成测试
run_integration_tests() {
    echo -e "${BLUE}🔗 运行集成测试...${NC}"
    
    # 设置测试环境变量
    export DB_HOST=localhost
    export DB_PORT=15432
    export DB_NAME=erp_test
    export DB_USER=erp_user
    export DB_PASSWORD=erp_pass
    export REDIS_HOST=localhost
    export REDIS_PORT=16379
    export REDIS_DB=1
    export LOG_LEVEL=debug
    
    # 运行集成测试
    go test -v -tags=integration ./integration_test.go -timeout=10m
    
    echo -e "${GREEN}✅ 集成测试通过${NC}"
}

# 运行负载测试
run_load_tests() {
    echo -e "${BLUE}⚡ 运行负载测试...${NC}"
    
    # 设置环境变量
    export DB_HOST=localhost
    export DB_PORT=15432
    export DB_NAME=erp_test
    export DB_USER=erp_user
    export DB_PASSWORD=erp_pass
    export REDIS_HOST=localhost
    export REDIS_PORT=16379
    export REDIS_DB=1
    
    # 运行基准测试
    go test -v -bench=. -benchmem ./integration_test.go -timeout=5m
    
    echo -e "${GREEN}✅ 负载测试完成${NC}"
}

# 生成测试报告
generate_test_report() {
    echo -e "${BLUE}📊 生成测试报告...${NC}"
    
    # 创建报告目录
    mkdir -p test-reports
    
    # 生成覆盖率报告
    go test -coverprofile=test-reports/coverage.out ./internal/...
    go tool cover -html=test-reports/coverage.out -o test-reports/coverage.html
    go tool cover -func=test-reports/coverage.out > test-reports/coverage.txt
    
    # 显示覆盖率摘要
    echo -e "${YELLOW}📈 代码覆盖率摘要:${NC}"
    tail -1 test-reports/coverage.txt
    
    echo -e "${GREEN}✅ 测试报告生成完成: test-reports/${NC}"
}

# 主函数
main() {
    echo -e "${GREEN}🎯 Frappe 权限系统集成测试套件${NC}"
    echo "======================================="
    
    # 检查参数
    case "${1:-all}" in
        "unit")
            echo -e "${BLUE}📝 仅运行单元测试${NC}"
            run_unit_tests
            ;;
        "integration") 
            echo -e "${BLUE}🔗 运行集成测试${NC}"
            check_docker
            cleanup
            start_test_env
            run_integration_tests
            cleanup
            ;;
        "load")
            echo -e "${BLUE}⚡ 运行负载测试${NC}"
            check_docker
            cleanup
            start_test_env
            run_load_tests
            cleanup
            ;;
        "all"|"")
            echo -e "${BLUE}🚀 运行完整测试套件${NC}"
            
            # 运行单元测试
            run_unit_tests
            
            # 运行集成测试
            check_docker
            cleanup
            start_test_env
            run_integration_tests
            
            # 运行负载测试
            run_load_tests
            
            # 生成报告
            generate_test_report
            
            # 清理
            cleanup
            ;;
        "help"|"-h"|"--help")
            echo -e "${YELLOW}用法: $0 [unit|integration|load|all]${NC}"
            echo ""
            echo -e "${YELLOW}选项:${NC}"
            echo -e "  unit        仅运行单元测试"
            echo -e "  integration 仅运行集成测试"
            echo -e "  load        仅运行负载测试"
            echo -e "  all         运行完整测试套件 (默认)"
            echo -e "  help        显示此帮助信息"
            exit 0
            ;;
        *)
            echo -e "${RED}❌ 未知选项: $1${NC}"
            echo -e "${YELLOW}使用 '$0 help' 查看使用说明${NC}"
            exit 1
            ;;
    esac
    
    echo ""
    echo -e "${GREEN}🎉 测试完成! 所有测试都已通过${NC}"
    echo "======================================="
}

# 捕获中断信号进行清理
trap cleanup EXIT INT TERM

# 运行主函数
main "$@"