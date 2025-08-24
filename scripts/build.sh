#!/bin/bash

# ================================================================
# 后端构建脚本
# ================================================================
# 功能：构建Go二进制文件、生成protobuf、运行测试等
# 使用：./build.sh {build|proto|test|clean|install}
# 作者：ERP系统开发组
# 版本：1.0.0

set -e

# ================================================================
# 配置参数
# ================================================================
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
BACKEND_DIR="$(cd "${SCRIPT_DIR}/../backend" && pwd)"

# 构建配置
BINARY_NAME="erp-server"
BUILD_DIR="${BACKEND_DIR}/build"
CMD_DIR="${BACKEND_DIR}/cmd/server"

# ================================================================
# 函数定义
# ================================================================

# 显示帮助信息
show_help() {
    echo "ERP系统后端构建脚本"
    echo ""
    echo "使用方法:"
    echo "  $0 {build|proto|test|clean|install|help}"
    echo ""
    echo "命令说明:"
    echo "  build    - 构建后端二进制文件"
    echo "  proto    - 生成protobuf代码"
    echo "  test     - 运行测试"
    echo "  clean    - 清理构建文件"
    echo "  install  - 安装依赖"
    echo "  help     - 显示此帮助信息"
}

# 构建二进制文件
build_binary() {
    echo "🔨 构建后端二进制文件..."
    
    cd "${BACKEND_DIR}"
    
    # 创建构建目录
    mkdir -p "${BUILD_DIR}"
    
    # 设置构建参数
    GOOS=${GOOS:-linux}
    GOARCH=${GOARCH:-amd64}
    CGO_ENABLED=${CGO_ENABLED:-0}
    
    # 构建信息
    VERSION=$(git describe --tags --always --dirty 2>/dev/null || echo "dev")
    COMMIT=$(git rev-parse HEAD 2>/dev/null || echo "unknown")
    BUILD_TIME=$(date +%Y-%m-%dT%H:%M:%S)
    
    # 构建标志
    LDFLAGS="-s -w"
    LDFLAGS="${LDFLAGS} -X main.Version=${VERSION}"
    LDFLAGS="${LDFLAGS} -X main.Commit=${COMMIT}"
    LDFLAGS="${LDFLAGS} -X main.BuildTime=${BUILD_TIME}"
    
    echo "版本: ${VERSION}"
    echo "提交: ${COMMIT}"
    echo "构建时间: ${BUILD_TIME}"
    echo "目标平台: ${GOOS}/${GOARCH}"
    
    # 执行构建
    CGO_ENABLED=${CGO_ENABLED} GOOS=${GOOS} GOARCH=${GOARCH} \
        go build -a -installsuffix cgo -ldflags "${LDFLAGS}" \
        -o "${BUILD_DIR}/${BINARY_NAME}" "${CMD_DIR}/main.go"
    
    echo "✅ 构建完成: ${BUILD_DIR}/${BINARY_NAME}"
}

# 生成protobuf代码
generate_proto() {
    echo "🔄 生成protobuf代码..."
    
    cd "${BACKEND_DIR}"
    
    # 检查protoc是否安装
    if ! command -v protoc &> /dev/null; then
        echo "❌ protoc未安装，请先安装Protocol Buffers编译器"
        exit 1
    fi
    
    # 生成API代码
    find api -name "*.proto" -exec protoc \
        --proto_path=. \
        --proto_path=third_party \
        --go_out=paths=source_relative:. \
        --go-http_out=paths=source_relative:. \
        --go-grpc_out=paths=source_relative:. \
        --validate_out=paths=source_relative,lang=go:. {} \;
    
    # 生成配置代码
    protoc --proto_path=. \
        --proto_path=third_party \
        --go_out=paths=source_relative:. \
        internal/conf/conf.proto
    
    echo "✅ protobuf代码生成完成"
}

# 运行测试
run_tests() {
    echo "🧪 运行测试..."
    
    cd "${BACKEND_DIR}"
    
    # 运行所有测试
    go test -v -race -coverprofile=coverage.out ./...
    
    # 生成覆盖率报告
    if [ -f coverage.out ]; then
        go tool cover -html=coverage.out -o coverage.html
        echo "📊 测试覆盖率报告已生成: coverage.html"
    fi
    
    echo "✅ 测试完成"
}

# 清理构建文件
clean_build() {
    echo "🧹 清理构建文件..."
    
    cd "${BACKEND_DIR}"
    
    # 清理构建目录
    if [ -d "${BUILD_DIR}" ]; then
        rm -rf "${BUILD_DIR}"
        echo "已清理构建目录: ${BUILD_DIR}"
    fi
    
    # 清理测试文件
    rm -f coverage.out coverage.html
    
    # 清理依赖缓存
    go clean -cache -modcache -i -r
    
    echo "✅ 清理完成"
}

# 安装依赖
install_deps() {
    echo "📦 安装依赖..."
    
    cd "${BACKEND_DIR}"
    
    # 下载依赖
    go mod download
    go mod verify
    go mod tidy
    
    echo "✅ 依赖安装完成"
}

# ================================================================
# 主逻辑
# ================================================================

case "${1}" in
    build)
        build_binary
        ;;
    proto)
        generate_proto
        ;;
    test)
        run_tests
        ;;
    clean)
        clean_build
        ;;
    install)
        install_deps
        ;;
    help|--help|-h)
        show_help
        ;;
    *)
        echo "❌ 未知命令: ${1}"
        show_help
        exit 1
        ;;
esac