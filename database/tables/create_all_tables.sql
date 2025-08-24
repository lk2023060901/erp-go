-- ================================================================
-- 创建所有数据表的主脚本
-- ================================================================
-- 按照正确的依赖顺序执行所有表创建脚本
-- 执行方式: psql -U postgres -d erp_system -f create_all_tables.sql

\echo '开始创建ERP系统数据表...'
\echo ''

-- 启用pg_trgm扩展（用于模糊搜索）
CREATE EXTENSION IF NOT EXISTS pg_trgm;

-- 设置时区
SET timezone = 'Asia/Shanghai';

\echo '1. 创建用户表...'
\i 01_users.sql
\echo '✓ 用户表创建完成'
\echo ''

\echo '2. 创建角色表...'
\i 02_roles.sql
\echo '✓ 角色表创建完成'
\echo ''

\echo '3. 创建权限表...'
\i 03_permissions.sql
\echo '✓ 权限表创建完成'
\echo ''

\echo '4. 创建组织表...'
\i 04_organizations.sql
\echo '✓ 组织表创建完成'
\echo ''

\echo '5. 创建用户角色关联表...'
\i 05_user_roles.sql
\echo '✓ 用户角色关联表创建完成'
\echo ''

\echo '6. 创建角色权限关联表...'
\i 06_role_permissions.sql
\echo '✓ 角色权限关联表创建完成'
\echo ''

\echo '7. 创建用户组织关联表...'
\i 07_user_organizations.sql
\echo '✓ 用户组织关联表创建完成'
\echo ''

\echo '8. 创建系统配置表...'
\i 08_system_configs.sql
\echo '✓ 系统配置表创建完成'
\echo ''

\echo '9. 创建操作日志表...'
\i 09_operation_logs.sql
\echo '✓ 操作日志表创建完成'
\echo ''

\echo '10. 创建用户会话表...'
\i 10_user_sessions.sql
\echo '✓ 用户会话表创建完成'
\echo ''

\echo '11. 创建用户过滤器表...'
\i 11_user_filters.sql
\echo '✓ 用户过滤器表创建完成'
\echo ''

\echo '所有数据表创建完成！'
\echo '下一步请执行索引、触发器和初始化数据脚本。'