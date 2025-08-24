# 数据表SQL脚本

本目录包含ERP系统数据库的单表创建脚本，便于维护和单独更新特定表结构。

## 脚本说明

### 核心表脚本
- `01_users.sql` - 用户表，存储用户基本信息和认证数据
- `02_roles.sql` - 角色表，存储角色定义和属性
- `03_permissions.sql` - 权限表，存储权限定义和层级结构
- `04_organizations.sql` - 组织表，支持无限层级组织架构

### 关联表脚本
- `05_user_roles.sql` - 用户角色关联表，支持多角色和临时权限
- `06_role_permissions.sql` - 角色权限关联表，支持权限授予和拒绝
- `07_user_organizations.sql` - 用户组织关联表，支持单主组织模式

### 辅助表脚本
- `08_system_configs.sql` - 系统配置表，存储运行时配置参数
- `09_operation_logs.sql` - 操作日志表，审计日志记录
- `10_user_sessions.sql` - 用户会话表，JWT Token管理

## 使用方法

### 执行单个表脚本
```bash
psql -U postgres -d erp_system -f 01_users.sql
```

### 执行所有表脚本
```bash
psql -U postgres -d erp_system -f create_all_tables.sql
```

### 脚本执行顺序
脚本按照数字前缀的顺序执行，确保外键依赖关系正确：
1. 先创建主表（users, roles, permissions, organizations）
2. 再创建关联表（user_roles, role_permissions, user_organizations）
3. 最后创建辅助表（system_configs, operation_logs, user_sessions）

## 注意事项

1. **依赖关系**：关联表依赖主表，必须按顺序执行
2. **删除表**：每个脚本都包含 `DROP TABLE IF EXISTS` 语句
3. **触发器**：每个表都包含自动更新 `updated_at` 字段的触发器
4. **索引**：包含了性能优化的关键索引
5. **约束**：包含了数据完整性的约束检查

## 扩展说明

如需修改单个表结构：
1. 编辑对应的SQL文件
2. 重新执行该文件（注意外键依赖）
3. 如有必要，同时更新完整的 `init.sql` 文件