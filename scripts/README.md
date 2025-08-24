# 项目管理脚本

本目录是ERP系统的唯一脚本目录，包含所有项目管理和构建脚本。

## 脚本说明

### build.sh - 后端构建脚本
构建Go二进制文件、生成protobuf代码、运行测试等。

```bash
# 构建后端
./scripts/build.sh build

# 生成protobuf代码
./scripts/build.sh proto

# 运行测试
./scripts/build.sh test

# 清理构建文件
./scripts/build.sh clean

# 安装依赖
./scripts/build.sh install
```

### docker.sh - Docker容器管理
管理整个项目的Docker容器，包括前端、后端、数据库和缓存服务。

```bash
# 启动所有服务
./scripts/docker.sh start

# 停止所有服务
./scripts/docker.sh stop

# 重新构建并启动
./scripts/docker.sh restart

# 查看服务状态
./scripts/docker.sh ps

# 查看日志
./scripts/docker.sh logs
```

### service.sh - 本地开发服务管理
管理本地开发环境的服务启动和停止。

```bash
# 启动开发环境
./scripts/service.sh start

# 停止开发环境
./scripts/service.sh stop

# 查看服务状态
./scripts/service.sh status
```

## 前端命令

前端相关命令直接在frontend目录下执行：

```bash
cd frontend
npm run dev      # 开发服务器
npm run build    # 生产构建
npm run preview  # 预览构建结果
```

## 使用建议

1. **开发阶段**: 使用 `./scripts/service.sh start` 启动本地开发环境
2. **构建测试**: 使用 `./scripts/build.sh build` 构建后端代码
3. **容器部署**: 使用 `./scripts/docker.sh start` 启动完整的容器环境
4. **生产部署**: 使用 `./scripts/docker.sh` 管理生产环境容器

## 项目结构

```
scripts/                    # 唯一的脚本目录
├── build.sh               # 后端构建和测试
├── docker.sh              # Docker容器管理
├── service.sh             # 服务管理
└── README.md              # 本文档

backend/                   # 后端Go项目（无scripts子目录）
frontend/                  # 前端React项目（使用npm scripts）
```