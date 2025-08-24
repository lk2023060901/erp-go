# Docker Compose 使用说明

## 概述
项目使用统一的docker-compose.yml文件，通过profiles区分不同环境。

## 环境配置

### 生产环境 (prod)
```bash
# 启动完整生产环境
docker-compose --profile prod up -d

# 包含服务：postgres, redis, api-server, frontend, nginx
```

### 开发环境 (dev)
```bash
# 启动开发环境（不含前端和nginx）
docker-compose --profile dev up -d

# 包含服务：postgres, redis, api-server
```

### 测试环境 (test)
```bash
# 启动测试环境
docker-compose --profile test up -d

# 包含服务：postgres-test, redis-test, api-server-test
```

### CI环境 (ci)
```bash
# 启动CI测试环境
docker-compose --profile ci up -d

# 包含服务：postgres-test, redis-test, api-server-test
# 使用临时内存存储，优化CI性能
```

## 环境变量配置

可以通过环境变量或.env文件自定义配置：

### 生产环境变量
```bash
# 数据库配置
POSTGRES_DB=erp_db
POSTGRES_USER=erp_user
POSTGRES_PASSWORD=erp_pass
POSTGRES_PORT=58001

# Redis配置
REDIS_PORT=58002

# API配置
API_PORT=58080
JWT_SECRET=your-jwt-secret-key
LOG_LEVEL=info
ENVIRONMENT=production

# 前端配置
FRONTEND_PORT=58003
FRONTEND_API_URL=http://localhost:58080
NGINX_PORT=58000
```

### 测试环境变量
```bash
# 测试数据库配置
POSTGRES_TEST_DB=erp_test_db
POSTGRES_TEST_USER=erp_test_user
POSTGRES_TEST_PASSWORD=erp_test_pass
POSTGRES_TEST_PORT=5433

# 测试Redis配置
REDIS_TEST_PORT=6380

# 测试API配置
API_TEST_PORT=58081
JWT_TEST_SECRET=test-jwt-secret-key
```

## 常用命令

```bash
# 查看当前运行的服务
docker-compose ps

# 查看特定profile的服务
docker-compose --profile prod config

# 停止所有服务
docker-compose down

# 停止并清理数据
docker-compose down -v

# 重新构建服务
docker-compose --profile dev up --build

# 查看日志
docker-compose logs -f api-server

# 进入服务容器
docker-compose exec api-server bash
```

## 数据库初始化

数据库会自动执行以下初始化脚本：
- `/database/schema/` - 表结构定义
- `/database/data/` - 种子数据

## 网络和存储

- **网络**: 所有服务使用`erp-network`桥接网络
- **持久存储**: 生产环境使用命名卷`postgres_data`和`redis_data`
- **临时存储**: 测试环境使用tmpfs，数据仅存在于内存中

## 环境切换示例

```bash
# 从开发环境切换到生产环境
docker-compose --profile dev down
docker-compose --profile prod up -d

# 运行集成测试
docker-compose --profile test up -d
go test ./tests/integration/...
docker-compose --profile test down
```