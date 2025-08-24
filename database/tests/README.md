# 测试数据目录

该目录包含用于开发和测试环境的数据文件。

## 文件说明

- **demo_users.sql**: 演示用户数据，包含不同角色的测试用户
- **permission_test_data.sql**: 权限系统测试数据，包含复杂的权限配置示例

## 使用方法

### SQLite 环境
```bash
# 导入演示用户数据
sqlite3 database.db < database/tests/demo_users.sql

# 导入权限测试数据
sqlite3 database.db < database/tests/permission_test_data.sql
```

### 集成测试
在运行集成测试前，确保加载了相应的测试数据。

## 测试用户说明

| 用户名 | 密码 | 角色 | 部门 |
|--------|------|------|------|
| demo_user1 | admin123 | 普通用户 | 研发部 |
| demo_user2 | admin123 | 普通用户 | 销售部 |
| demo_user3 | admin123 | 普通用户 | 市场部 |
| demo_manager | admin123 | 管理员 | - |

## 注意事项

- 这些数据仅用于开发和测试环境
- 生产环境不应使用这些测试数据
- 所有测试用户的密码都是 `admin123`