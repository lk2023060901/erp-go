-- ================================================================
-- 性能优化索引脚本
-- ================================================================
-- 为权限验证和常用查询创建高性能索引

\echo '创建权限验证核心索引...'

-- 权限验证最关键的复合索引
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_user_roles_permission_verify 
ON user_roles(user_id, is_active, expires_at) 
WHERE is_active = true;

CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_role_permissions_verify 
ON role_permissions(role_id, is_granted, permission_id) 
WHERE is_granted = true;

-- 权限验证视图索引
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_permissions_code_resource_action 
ON permissions(code, resource, action) 
WHERE deleted_at IS NULL;

\echo '创建用户查询优化索引...'

-- 用户登录索引
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_users_login_username 
ON users(username, password_hash, is_enabled) 
WHERE deleted_at IS NULL;

CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_users_login_email 
ON users(email, password_hash, is_enabled) 
WHERE deleted_at IS NULL;

-- 用户状态索引
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_users_status_enabled 
ON users(status, is_enabled, created_at DESC) 
WHERE deleted_at IS NULL;

\echo '创建组织查询索引...'

-- 组织层级查询索引
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_organizations_path_level 
ON organizations USING GIN (path gin_trgm_ops) 
WHERE deleted_at IS NULL;

-- 用户主组织索引
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_user_orgs_primary 
ON user_organizations(user_id) 
WHERE is_primary = true AND is_active = true;

\echo '创建日志查询索引...'

-- 操作日志时间序列索引
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_operation_logs_user_time 
ON operation_logs(user_id, created_at DESC, operation_type);

-- 操作日志错误查询索引
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_operation_logs_errors 
ON operation_logs(is_success, created_at DESC, error_message) 
WHERE is_success = false;

\echo '创建会话管理索引...'

-- 活跃会话查询索引
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_user_sessions_active 
ON user_sessions(user_id, is_active, last_activity_at DESC) 
WHERE is_active = true;

-- Token查找索引
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_user_sessions_token_active 
ON user_sessions(token_hash, is_active, expires_at) 
WHERE is_active = true;

\echo '创建系统配置索引...'

-- 配置查询索引
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_system_configs_key_enabled 
ON system_configs(config_key, is_enabled) 
WHERE is_enabled = true;

CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_system_configs_group_type 
ON system_configs(group_name, config_type, sort_order);

\echo '创建角色权限聚合索引...'

-- 角色有效权限索引
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_roles_enabled_system 
ON roles(is_enabled, is_system_role, code) 
WHERE deleted_at IS NULL;

-- 权限类型筛选索引
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_permissions_type_module 
ON permissions(is_menu, is_button, is_api, module) 
WHERE deleted_at IS NULL;

\echo '所有性能索引创建完成！'