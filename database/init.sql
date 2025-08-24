-- ================================================================
-- ERP系统数据库初始化脚本
-- ================================================================
-- Docker容器启动时自动执行的数据库初始化脚本
-- 按正确顺序创建表结构和加载初始数据

\echo '开始初始化ERP系统数据库...'
\echo ''

-- 设置时区和编码
SET timezone = 'Asia/Shanghai';
SET client_encoding = 'UTF8';

-- 启用必要的扩展
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

\echo '1. 创建数据库表结构...'

-- 创建用户表
\i tables/01_users.sql

-- 创建角色表
\i tables/02_roles.sql

-- 创建权限表
\i tables/03_permissions.sql

-- 创建组织表
\i tables/04_organizations.sql

-- 创建关联表
\i tables/05_user_roles.sql
\i tables/06_role_permissions.sql
\i tables/07_user_organizations.sql

-- 创建系统表
\i tables/08_system_configs.sql
\i tables/09_operation_logs.sql
\i tables/10_user_sessions.sql
\i tables/11_user_filters.sql

\echo '2. 创建索引和触发器...'

-- 创建性能索引
\i indexes/performance_indexes.sql

-- 创建审计触发器
\i triggers/audit_triggers.sql

\echo '3. 加载初始数据...'

-- 加载所有种子数据
\i seed_data/load_all_seed_data.sql

\echo ''
\echo '========================================='
\echo 'ERP系统数据库初始化完成！'
\echo ''
\echo '系统访问信息：'
\echo '  数据库: erp_system'
\echo '  用户名: admin'
\echo '  密码: Admin123@'
\echo ''
\echo '请在生产环境中及时修改默认密码！'
\echo '========================================='