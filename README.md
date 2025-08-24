# ERP系统 - 企业资源规划系统

基于Go + Kratos框架的现代化ERP系统，支持完整的用户权限管理（RBAC）、多角色、无限层级组织架构等功能。

## 系统架构

- **后端**: Go + Kratos框架 + PostgreSQL + Redis
- **前端**: React 18 + TypeScript + Vite + Ant Design
- **部署**: Docker + Docker Compose + Nginx

## 功能特性

### 用户管理
- 支持用户名/邮箱登录
- 完整的用户信息（姓名、手机、邮箱、性别、生日等）
- 登录时间和IP地址跟踪
- 双因素认证（Google Authenticator）

### 权限系统（RBAC）
- 基于角色的访问控制
- 用户多角色支持
- 三级权限控制：菜单权限、按钮权限、API权限
- 权限继承和覆盖机制

### 组织架构
- 无限层级组织结构
- 用户多组织归属
- 组织权限继承

### 系统功能
- JWT认证 + 刷新令牌机制
- 三级缓存（内存 + 本地缓存 + Redis）
- 完整的操作审计日志
- RESTful API设计
- 实时权限验证

## 项目结构

```
├── api/                    # API定义（protobuf）
├── backend/
│   └── scripts/           # 后端管理脚本
├── cmd/
│   └── server/           # 服务器入口
├── configs/              # 配置文件
├── database/             # 数据库相关
│   ├── migrations/       # 数据库迁移
│   ├── tables/          # 表结构
│   ├── seed_data/       # 种子数据
│   └── triggers/        # 触发器
├── frontend/            # 前端应用
├── internal/           # 内部包
│   ├── biz/           # 业务逻辑
│   ├── data/          # 数据层
│   ├── middleware/    # 中间件
│   └── service/       # 服务层
└── docker-compose.yml  # Docker编排
```

## 快速开始

### 环境要求

- Go 1.21+
- Node.js 18+
- PostgreSQL 12+
- Redis 6+
- Docker & Docker Compose

### 使用Docker快速启动

1. 克隆项目并配置环境变量：
```bash
cp .env.example .env
# 编辑 .env 文件配置数据库和其他参数
```

2. 启动所有服务：
```bash
# 使用管理脚本
./backend/scripts/docker.sh start

# 或直接使用docker-compose
docker-compose up -d
```

3. 访问应用：
- 前端应用: http://localhost:3000
- 后端API: http://localhost:58080
- API文档: http://localhost:58080/swagger

### 本地开发

1. 安装后端依赖：
```bash
go mod tidy
```

2. 初始化数据库：
```bash
# 启动PostgreSQL和Redis
docker-compose up -d postgres redis

# 运行数据库初始化脚本
psql -h localhost -U erp_user -d erp_db -f database/migrations/init.sql
```

3. 启动后端服务：
```bash
go run cmd/server/main.go -conf configs/config.yaml
```

4. 安装和启动前端：
```bash
cd frontend
npm install
npm run dev
```

## API文档

### 认证相关
- `POST /v1/auth/login` - 用户登录
- `POST /v1/auth/logout` - 用户登出
- `POST /v1/auth/refresh` - 刷新令牌
- `GET /v1/auth/profile` - 获取用户信息

### 用户管理
- `GET /v1/users` - 获取用户列表
- `POST /v1/users` - 创建用户
- `PUT /v1/users/{id}` - 更新用户
- `DELETE /v1/users/{id}` - 删除用户

### 角色权限
- `GET /v1/roles` - 获取角色列表
- `GET /v1/permissions` - 获取权限列表
- `POST /v1/permissions/check` - 检查权限

## 默认账户

系统初始化后会创建以下默认账户：

- **超级管理员**: admin / admin123
- **普通用户**: user / user123

## 技术栈详情

### 后端技术栈
- **Web框架**: Kratos v2.7+
- **数据库**: PostgreSQL 15
- **缓存**: Redis 7
- **认证**: JWT + 双因素认证
- **API**: gRPC + HTTP
- **日志**: 结构化日志
- **监控**: Prometheus + Grafana（可选）

### 前端技术栈
- **框架**: React 18 + TypeScript
- **构建工具**: Vite
- **UI库**: Ant Design 5
- **路由**: React Router v6
- **状态管理**: Zustand
- **HTTP客户端**: Axios

## 部署说明

### Docker部署

```bash
# 生产环境部署
./backend/scripts/docker.sh deploy

# 停止服务
./backend/scripts/docker.sh stop

# 查看日志
./backend/scripts/docker.sh logs
```

### 数据库备份

```bash
# 备份数据库
./backend/scripts/service.sh backup

# 恢复数据库
./backend/scripts/service.sh restore backup_file.sql
```

## 开发指南

### 添加新的API接口

1. 在 `api/` 目录下定义protobuf文件
2. 在 `internal/service/` 实现服务逻辑
3. 在 `internal/biz/` 实现业务逻辑
4. 在 `internal/data/` 实现数据访问层

### 权限控制

使用中间件进行权限控制：

```go
// 需要登录
server.Use(middleware.Auth())

// 需要特定权限
server.Use(middleware.Permission("user.create"))

// 需要角色
server.Use(middleware.Role("ADMIN"))
```

## 贡献指南

1. Fork项目
2. 创建功能分支
3. 提交代码
4. 发起Pull Request

## 许可证

MIT License