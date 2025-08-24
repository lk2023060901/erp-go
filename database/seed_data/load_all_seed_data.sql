-- ================================================================
-- 加载所有初始化数据的主脚本
-- ================================================================
-- 按照正确的依赖顺序执行所有种子数据脚本
-- 执行方式: psql -U postgres -d erp_system -f load_all_seed_data.sql

\echo '开始加载ERP系统初始化数据...'
\echo ''

-- 设置时区
SET timezone = 'Asia/Shanghai';

\echo '1. 插入系统默认角色...'
\i 01_system_roles.sql
\echo ''

\echo '2. 插入系统权限数据...'
\i 02_system_permissions.sql
\echo ''

\echo '3. 配置角色权限关联...'
\i 03_role_permissions.sql
\echo ''

\echo '4. 创建默认组织架构...'
\i 04_default_organization.sql
\echo ''

\echo '5. 创建默认管理员用户...'
\i 05_admin_user.sql
\echo ''

\echo '6. 加载系统配置数据...'
\i 06_system_configs.sql
\echo ''

\echo '所有初始化数据加载完成！'
\echo ''
\echo '========================================='
\echo '系统初始化完成，以下是默认登录信息：'
\echo ''
\echo '超级管理员:'
\echo '  用户名: admin'
\echo '  密码: Admin123@'
\echo '  角色: 超级管理员'
\echo ''
\echo '测试用户:'
\echo '  用户名: user001 / 密码: Admin123@ (普通用户)'
\echo '  用户名: manager01 / 密码: Admin123@ (部门经理)'
\echo ''
\echo '请及时修改默认密码！'
\echo '========================================='