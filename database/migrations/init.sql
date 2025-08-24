-- ================================================================================================
-- ERP系统用户权限管理数据库完整初始化脚本
-- 基于RBAC (Role-Based Access Control) 权限控制模型
-- 支持：用户多角色、无限层级组织、权限继承、双重验证、审计日志
-- 
-- 版本: v1.0.0
-- 创建时间: 2024-08-22
-- 适用数据库: PostgreSQL 12+
-- ================================================================================================

-- 设置字符编码和时区
SET client_encoding = 'UTF8';
SET timezone = 'UTC';

-- 开启事务
BEGIN;

-- ================================================================================================
-- 扩展和函数
-- ================================================================================================

-- 启用必要的PostgreSQL扩展
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- 创建更新时间触发器函数
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

-- ================================================================================================
-- 1. 用户表 (users) - 存储用户基本信息
-- 支持分库分表：基于user_id进行分片
-- ================================================================================================
DROP TABLE IF EXISTS users CASCADE;

CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    username VARCHAR(50) NOT NULL UNIQUE,                      -- 用户名，唯一
    email VARCHAR(100) NOT NULL UNIQUE,                        -- 邮箱，唯一
    password_hash VARCHAR(255) NOT NULL,                       -- 密码哈希值
    salt VARCHAR(32) NOT NULL,                                 -- 密码盐值
    
    -- 个人基本信息
    first_name VARCHAR(50) NOT NULL,                           -- 姓
    last_name VARCHAR(50) NOT NULL,                            -- 名
    phone VARCHAR(20),                                         -- 手机号
    gender SMALLINT CHECK (gender IN (0, 1, 2)) DEFAULT 0,    -- 性别: 0-未知, 1-男, 2-女
    birth_date DATE,                                           -- 出生年月
    avatar_url VARCHAR(500),                                   -- 头像URL
    
    -- 登录相关信息
    last_login_time TIMESTAMP WITH TIME ZONE,                 -- 最近登录时间
    last_login_ip INET,                                        -- 最近登录IP
    last_logout_time TIMESTAMP WITH TIME ZONE,                -- 最近退出时间  
    last_logout_ip INET,                                       -- 最近退出IP
    login_count INTEGER DEFAULT 0,                            -- 登录次数
    failed_login_count INTEGER DEFAULT 0,                     -- 失败登录次数
    locked_until TIMESTAMP WITH TIME ZONE,                    -- 账户锁定到期时间
    
    -- 双重验证
    two_factor_enabled BOOLEAN DEFAULT FALSE,                 -- 是否启用双重验证
    two_factor_secret VARCHAR(32),                            -- Google Authenticator密钥
    backup_codes TEXT[],                                       -- 备用验证码
    
    -- 账户状态
    is_enabled BOOLEAN DEFAULT TRUE,                          -- 账户是否启用
    email_verify_token VARCHAR(64),                           -- 邮箱验证令牌
    email_verify_expires TIMESTAMP WITH TIME ZONE,            -- 邮箱验证过期时间
    password_reset_token VARCHAR(64),                         -- 密码重置令牌
    password_reset_expires TIMESTAMP WITH TIME ZONE,          -- 密码重置过期时间
    
    -- 审计字段
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    created_by BIGINT,                                         -- 创建人ID
    updated_by BIGINT,                                         -- 更新人ID
    
    -- 软删除
    deleted_at TIMESTAMP WITH TIME ZONE,
    
    CONSTRAINT fk_users_created_by FOREIGN KEY (created_by) REFERENCES users(id),
    CONSTRAINT fk_users_updated_by FOREIGN KEY (updated_by) REFERENCES users(id)
);

-- 用户表索引
CREATE UNIQUE INDEX idx_users_username ON users(username) WHERE deleted_at IS NULL;
CREATE UNIQUE INDEX idx_users_email ON users(email) WHERE deleted_at IS NULL;
CREATE INDEX idx_users_phone ON users(phone) WHERE deleted_at IS NULL;
CREATE INDEX idx_users_created_at ON users(created_at);
CREATE INDEX idx_users_last_login ON users(last_login_time);
CREATE INDEX idx_users_is_enabled ON users(is_enabled) WHERE deleted_at IS NULL;

-- 用户表触发器
CREATE TRIGGER update_users_updated_at BEFORE UPDATE ON users FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- ================================================================================================
-- 2. 角色表 (roles) - 存储角色定义
-- ================================================================================================
DROP TABLE IF EXISTS roles CASCADE;

CREATE TABLE roles (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL,                                -- 角色名称
    code VARCHAR(50) NOT NULL UNIQUE,                         -- 角色编码，唯一
    description TEXT,                                          -- 角色描述
    
    -- 角色属性
    is_system_role BOOLEAN DEFAULT FALSE,                     -- 是否为系统预定义角色
    is_enabled BOOLEAN DEFAULT TRUE,                          -- 角色是否启用
    sort_order INTEGER DEFAULT 0,                             -- 排序顺序
    
    -- 审计字段
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    created_by BIGINT,
    updated_by BIGINT,
    
    -- 软删除
    deleted_at TIMESTAMP WITH TIME ZONE,
    
    CONSTRAINT fk_roles_created_by FOREIGN KEY (created_by) REFERENCES users(id),
    CONSTRAINT fk_roles_updated_by FOREIGN KEY (updated_by) REFERENCES users(id)
);

-- 角色表索引
CREATE UNIQUE INDEX idx_roles_code ON roles(code) WHERE deleted_at IS NULL;
CREATE INDEX idx_roles_name ON roles(name) WHERE deleted_at IS NULL;
CREATE INDEX idx_roles_system ON roles(is_system_role);
CREATE INDEX idx_roles_enabled ON roles(is_enabled) WHERE deleted_at IS NULL;

-- 角色表触发器
CREATE TRIGGER update_roles_updated_at BEFORE UPDATE ON roles FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- ================================================================================================
-- 3. 权限表 (permissions) - 存储权限定义
-- 功能级权限控制，细化到按钮级别
-- ================================================================================================
DROP TABLE IF EXISTS permissions CASCADE;

CREATE TABLE permissions (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,                               -- 权限名称
    code VARCHAR(100) NOT NULL UNIQUE,                        -- 权限编码，唯一
    resource VARCHAR(50) NOT NULL,                            -- 资源名称 (如: user, role, order)
    action VARCHAR(20) NOT NULL,                              -- 操作类型 (如: read, write, delete)
    description TEXT,                                          -- 权限描述
    
    -- 权限层级和分类
    module VARCHAR(50) NOT NULL,                              -- 所属模块 (如: system, finance, hr)
    parent_id BIGINT,                                          -- 父权限ID (支持权限树形结构)
    level INTEGER DEFAULT 1,                                  -- 权限层级
    path VARCHAR(500),                                         -- 权限路径 (如: /system/user/read)
    
    -- 权限属性  
    is_menu BOOLEAN DEFAULT FALSE,                            -- 是否为菜单权限
    is_button BOOLEAN DEFAULT FALSE,                          -- 是否为按钮权限
    is_api BOOLEAN DEFAULT FALSE,                             -- 是否为API权限
    sort_order INTEGER DEFAULT 0,                             -- 排序顺序
    
    -- 审计字段
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    
    CONSTRAINT fk_permissions_parent FOREIGN KEY (parent_id) REFERENCES permissions(id)
);

-- 权限表索引
CREATE UNIQUE INDEX idx_permissions_code ON permissions(code);
CREATE INDEX idx_permissions_resource_action ON permissions(resource, action);
CREATE INDEX idx_permissions_module ON permissions(module);
CREATE INDEX idx_permissions_parent ON permissions(parent_id);
CREATE INDEX idx_permissions_path ON permissions(path);
CREATE INDEX idx_permissions_is_menu ON permissions(is_menu) WHERE is_menu = true;
CREATE INDEX idx_permissions_is_button ON permissions(is_button) WHERE is_button = true;
CREATE INDEX idx_permissions_is_api ON permissions(is_api) WHERE is_api = true;

-- ================================================================================================
-- 4. 用户角色关联表 (user_roles) - 多对多关系
-- 支持用户绑定多个角色
-- ================================================================================================
DROP TABLE IF EXISTS user_roles CASCADE;

CREATE TABLE user_roles (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,                                  -- 用户ID
    role_id BIGINT NOT NULL,                                  -- 角色ID
    
    -- 授权信息
    assigned_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP, -- 分配时间
    assigned_by BIGINT,                                       -- 分配人ID
    expires_at TIMESTAMP WITH TIME ZONE,                      -- 过期时间 (临时权限)
    is_active BOOLEAN DEFAULT TRUE,                           -- 是否激活
    
    -- 审计字段
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    
    CONSTRAINT fk_user_roles_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT fk_user_roles_role FOREIGN KEY (role_id) REFERENCES roles(id) ON DELETE CASCADE,
    CONSTRAINT fk_user_roles_assigned_by FOREIGN KEY (assigned_by) REFERENCES users(id),
    CONSTRAINT uk_user_roles UNIQUE(user_id, role_id)
);

-- 用户角色关联表索引 (重要：权限验证核心索引)
CREATE INDEX idx_user_roles_user ON user_roles(user_id, is_active);
CREATE INDEX idx_user_roles_role ON user_roles(role_id, is_active);
CREATE INDEX idx_user_roles_expires ON user_roles(expires_at) WHERE expires_at IS NOT NULL;
CREATE INDEX idx_user_roles_assigned_at ON user_roles(assigned_at);

-- 用户角色表触发器
CREATE TRIGGER update_user_roles_updated_at BEFORE UPDATE ON user_roles FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- ================================================================================================
-- 5. 角色权限关联表 (role_permissions) - 多对多关系
-- 角色与权限的关联关系
-- ================================================================================================
DROP TABLE IF EXISTS role_permissions CASCADE;

CREATE TABLE role_permissions (
    id BIGSERIAL PRIMARY KEY,
    role_id BIGINT NOT NULL,                                  -- 角色ID
    permission_id BIGINT NOT NULL,                            -- 权限ID
    
    -- 授权信息
    assigned_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP, -- 分配时间
    assigned_by BIGINT,                                       -- 分配人ID
    is_granted BOOLEAN DEFAULT TRUE,                          -- 是否授权 (支持拒绝权限)
    
    -- 审计字段
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    
    CONSTRAINT fk_role_permissions_role FOREIGN KEY (role_id) REFERENCES roles(id) ON DELETE CASCADE,
    CONSTRAINT fk_role_permissions_permission FOREIGN KEY (permission_id) REFERENCES permissions(id) ON DELETE CASCADE,
    CONSTRAINT fk_role_permissions_assigned_by FOREIGN KEY (assigned_by) REFERENCES users(id),
    CONSTRAINT uk_role_permissions UNIQUE(role_id, permission_id)
);

-- 角色权限关联表索引 (重要：权限验证核心索引)
CREATE INDEX idx_role_permissions_role ON role_permissions(role_id, is_granted);
CREATE INDEX idx_role_permissions_permission ON role_permissions(permission_id);
CREATE INDEX idx_role_permissions_assigned_at ON role_permissions(assigned_at);

-- 角色权限表触发器
CREATE TRIGGER update_role_permissions_updated_at BEFORE UPDATE ON role_permissions FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- ================================================================================================
-- 6. 组织表 (organizations) - 支持无限层级的组织架构
-- ================================================================================================
DROP TABLE IF EXISTS organizations CASCADE;

CREATE TABLE organizations (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,                               -- 组织名称
    code VARCHAR(50) UNIQUE,                                  -- 组织编码
    description TEXT,                                          -- 组织描述
    
    -- 层级结构 (支持无限层级)
    parent_id BIGINT,                                          -- 父组织ID
    level INTEGER DEFAULT 1,                                  -- 组织层级 (1为顶级)
    path VARCHAR(500),                                         -- 组织路径 (如: /1/2/3/)
    sort_order INTEGER DEFAULT 0,                             -- 同级排序
    
    -- 组织属性
    org_type VARCHAR(20) DEFAULT 'department',                -- 组织类型: company, department, team
    is_enabled BOOLEAN DEFAULT TRUE,                          -- 是否启用
    contact_person VARCHAR(50),                               -- 负责人
    contact_phone VARCHAR(20),                                -- 联系电话
    contact_email VARCHAR(100),                               -- 联系邮箱
    address TEXT,                                              -- 组织地址
    
    -- 权限隔离设置
    data_isolation BOOLEAN DEFAULT FALSE,                     -- 是否启用数据隔离
    
    -- 审计字段
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    created_by BIGINT,
    updated_by BIGINT,
    
    -- 软删除
    deleted_at TIMESTAMP WITH TIME ZONE,
    
    CONSTRAINT fk_organizations_parent FOREIGN KEY (parent_id) REFERENCES organizations(id),
    CONSTRAINT fk_organizations_created_by FOREIGN KEY (created_by) REFERENCES users(id),
    CONSTRAINT fk_organizations_updated_by FOREIGN KEY (updated_by) REFERENCES users(id)
);

-- 组织表索引
CREATE INDEX idx_organizations_parent ON organizations(parent_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_organizations_path ON organizations(path) WHERE deleted_at IS NULL;
CREATE INDEX idx_organizations_level ON organizations(level) WHERE deleted_at IS NULL;
CREATE INDEX idx_organizations_code ON organizations(code) WHERE deleted_at IS NULL;
CREATE INDEX idx_organizations_type ON organizations(org_type) WHERE deleted_at IS NULL;
CREATE INDEX idx_organizations_enabled ON organizations(is_enabled) WHERE deleted_at IS NULL;

-- 组织表触发器
CREATE TRIGGER update_organizations_updated_at BEFORE UPDATE ON organizations FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- ================================================================================================
-- 7. 用户组织关联表 (user_organizations) - 用户单组织归属
-- ================================================================================================
DROP TABLE IF EXISTS user_organizations CASCADE;

CREATE TABLE user_organizations (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,                                  -- 用户ID
    org_id BIGINT NOT NULL,                                   -- 组织ID
    position VARCHAR(50),                                      -- 职位
    is_primary BOOLEAN DEFAULT TRUE,                          -- 是否为主组织 (用户单组织)
    
    -- 任职信息
    joined_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP, -- 加入时间
    left_at TIMESTAMP WITH TIME ZONE,                         -- 离开时间
    is_active BOOLEAN DEFAULT TRUE,                           -- 是否活跃
    
    -- 审计字段
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    
    CONSTRAINT fk_user_organizations_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT fk_user_organizations_org FOREIGN KEY (org_id) REFERENCES organizations(id) ON DELETE CASCADE,
    CONSTRAINT uk_user_organizations UNIQUE(user_id, org_id)
);

-- 用户组织关联表索引
CREATE INDEX idx_user_organizations_user ON user_organizations(user_id, is_active);
CREATE INDEX idx_user_organizations_org ON user_organizations(org_id, is_active);
CREATE INDEX idx_user_organizations_primary ON user_organizations(user_id, is_primary) WHERE is_primary = TRUE;
CREATE INDEX idx_user_organizations_joined_at ON user_organizations(joined_at);

-- 用户组织表触发器
CREATE TRIGGER update_user_organizations_updated_at BEFORE UPDATE ON user_organizations FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- ================================================================================================
-- 8. 系统配置表 (system_configs) - 存储系统配置
-- ================================================================================================
DROP TABLE IF EXISTS system_configs CASCADE;

CREATE TABLE system_configs (
    id BIGSERIAL PRIMARY KEY,
    config_key VARCHAR(100) NOT NULL UNIQUE,                 -- 配置键
    config_value TEXT,                                         -- 配置值
    config_type VARCHAR(20) DEFAULT 'string',                 -- 配置类型: string, int, bool, json
    description TEXT,                                          -- 配置描述
    is_public BOOLEAN DEFAULT FALSE,                          -- 是否公开 (前端可访问)
    is_encrypted BOOLEAN DEFAULT FALSE,                       -- 是否加密存储
    
    -- 审计字段
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_by BIGINT,
    
    CONSTRAINT fk_system_configs_updated_by FOREIGN KEY (updated_by) REFERENCES users(id)
);

-- 系统配置表索引
CREATE UNIQUE INDEX idx_system_configs_key ON system_configs(config_key);
CREATE INDEX idx_system_configs_public ON system_configs(is_public);
CREATE INDEX idx_system_configs_type ON system_configs(config_type);

-- 系统配置表触发器
CREATE TRIGGER update_system_configs_updated_at BEFORE UPDATE ON system_configs FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- ================================================================================================
-- 9. 操作日志表 (operation_logs) - 记录用户操作日志
-- 支持分库分表：基于created_at按时间分片
-- ================================================================================================
DROP TABLE IF EXISTS operation_logs CASCADE;

CREATE TABLE operation_logs (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT,                                           -- 操作用户ID
    username VARCHAR(50),                                      -- 操作用户名
    operation VARCHAR(50) NOT NULL,                          -- 操作类型
    resource VARCHAR(50),                                      -- 操作资源
    resource_id VARCHAR(50),                                   -- 资源ID
    details JSONB,                                            -- 操作详情 (JSON格式)
    
    -- 请求信息
    ip_address INET,                                          -- 操作IP
    user_agent TEXT,                                          -- 用户代理
    request_method VARCHAR(10),                               -- 请求方法
    request_url VARCHAR(500),                                 -- 请求URL
    request_params JSONB,                                     -- 请求参数
    
    -- 结果信息
    status VARCHAR(20) DEFAULT 'success',                     -- 操作状态: success, failed, error
    error_message TEXT,                                       -- 错误信息
    response_time INTEGER,                                    -- 响应时间(ms)
    
    -- 审计字段
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    
    CONSTRAINT fk_operation_logs_user FOREIGN KEY (user_id) REFERENCES users(id)
);

-- 操作日志表索引
CREATE INDEX idx_operation_logs_user ON operation_logs(user_id, created_at);
CREATE INDEX idx_operation_logs_operation ON operation_logs(operation, created_at);
CREATE INDEX idx_operation_logs_resource ON operation_logs(resource, resource_id);
CREATE INDEX idx_operation_logs_ip ON operation_logs(ip_address, created_at);
CREATE INDEX idx_operation_logs_status ON operation_logs(status, created_at);
CREATE INDEX idx_operation_logs_created_at ON operation_logs(created_at);

-- ================================================================================================
-- 10. 会话表 (user_sessions) - 用户登录会话管理
-- ================================================================================================
DROP TABLE IF EXISTS user_sessions CASCADE;

CREATE TABLE user_sessions (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,                                  -- 用户ID
    session_token VARCHAR(64) NOT NULL UNIQUE,               -- 会话令牌
    refresh_token VARCHAR(64) NOT NULL UNIQUE,               -- 刷新令牌
    
    -- 会话信息
    ip_address INET,                                          -- 登录IP
    user_agent TEXT,                                          -- 用户代理
    device_type VARCHAR(20),                                  -- 设备类型: web, mobile, desktop
    location VARCHAR(100),                                    -- 登录地点
    
    -- 时间信息
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    last_accessed_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,            -- 过期时间
    is_active BOOLEAN DEFAULT TRUE,                          -- 是否活跃
    
    CONSTRAINT fk_user_sessions_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- 会话表索引
CREATE UNIQUE INDEX idx_user_sessions_session_token ON user_sessions(session_token);
CREATE UNIQUE INDEX idx_user_sessions_refresh_token ON user_sessions(refresh_token);
CREATE INDEX idx_user_sessions_user ON user_sessions(user_id, is_active);
CREATE INDEX idx_user_sessions_expires_at ON user_sessions(expires_at);
CREATE INDEX idx_user_sessions_last_accessed ON user_sessions(last_accessed_at);

-- ================================================================================================
-- 权限验证相关视图和函数
-- ================================================================================================

-- 创建用户权限视图 (优化查询性能)
CREATE OR REPLACE VIEW user_permissions_view AS
SELECT DISTINCT
    ur.user_id,
    p.id as permission_id,
    p.code as permission_code,
    p.name as permission_name,
    p.resource,
    p.action,
    p.module,
    r.id as role_id,
    r.code as role_code,
    r.name as role_name
FROM user_roles ur
JOIN role_permissions rp ON ur.role_id = rp.role_id AND rp.is_granted = true
JOIN permissions p ON rp.permission_id = p.id
JOIN roles r ON ur.role_id = r.id
WHERE ur.is_active = true 
    AND r.is_enabled = true
    AND (ur.expires_at IS NULL OR ur.expires_at > CURRENT_TIMESTAMP);

-- 创建权限验证函数
CREATE OR REPLACE FUNCTION check_user_permission(
    p_user_id BIGINT,
    p_permission_code VARCHAR(100)
) RETURNS BOOLEAN AS $$
DECLARE
    permission_count INTEGER;
BEGIN
    SELECT COUNT(*) INTO permission_count
    FROM user_permissions_view
    WHERE user_id = p_user_id 
        AND permission_code = p_permission_code;
    
    RETURN permission_count > 0;
END;
$$ LANGUAGE plpgsql;

-- 创建获取用户所有权限函数
CREATE OR REPLACE FUNCTION get_user_permissions(
    p_user_id BIGINT
) RETURNS TABLE(
    permission_code VARCHAR(100),
    permission_name VARCHAR(100),
    resource VARCHAR(50),
    action VARCHAR(20),
    module VARCHAR(50)
) AS $$
BEGIN
    RETURN QUERY
    SELECT DISTINCT
        upv.permission_code,
        upv.permission_name,
        upv.resource,
        upv.action,
        upv.module
    FROM user_permissions_view upv
    WHERE upv.user_id = p_user_id
    ORDER BY upv.module, upv.permission_code;
END;
$$ LANGUAGE plpgsql;

-- 创建组织路径更新函数
CREATE OR REPLACE FUNCTION update_organization_path()
RETURNS TRIGGER AS $$
DECLARE
    parent_path VARCHAR(500);
BEGIN
    IF NEW.parent_id IS NULL THEN
        NEW.path := '/' || NEW.id || '/';
        NEW.level := 1;
    ELSE
        SELECT path, level INTO parent_path, NEW.level FROM organizations WHERE id = NEW.parent_id;
        NEW.path := parent_path || NEW.id || '/';
        NEW.level := NEW.level + 1;
    END IF;
    
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- 创建组织路径触发器
CREATE TRIGGER update_org_path_trigger
    BEFORE INSERT OR UPDATE ON organizations
    FOR EACH ROW EXECUTE FUNCTION update_organization_path();

-- ================================================================================================
-- 初始化数据
-- ================================================================================================

-- 插入系统配置
INSERT INTO system_configs (config_key, config_value, config_type, description, is_public) VALUES
('user_registration_enabled', 'true', 'bool', '是否允许用户自主注册', true),
('email_verification_enabled', 'true', 'bool', '是否启用邮箱验证', true),
('welcome_email_enabled', 'true', 'bool', '是否发送欢迎邮件', false),
('two_factor_required', 'false', 'bool', '是否强制双重验证', true),
('password_min_length', '8', 'int', '密码最小长度', true),
('password_complexity', 'true', 'bool', '是否要求密码复杂度', true),
('session_timeout', '86400', 'int', '会话超时时间(秒)', false),
('max_login_attempts', '5', 'int', '最大登录尝试次数', false),
('account_lock_duration', '1800', 'int', '账户锁定时长(秒)', false),
('jwt_secret_key', 'your-jwt-secret-key-change-in-production', 'string', 'JWT密钥', false),
('jwt_expires_in', '3600', 'int', 'JWT过期时间(秒)', false),
('refresh_token_expires_in', '604800', 'int', '刷新令牌过期时间(秒)', false);

-- 创建默认组织 (总公司)
INSERT INTO organizations (name, code, description, org_type, level, path) VALUES
('总公司', 'ROOT', '系统默认根组织', 'company', 1, '/1/');

-- 创建系统默认权限
INSERT INTO permissions (name, code, resource, action, module, description, is_menu, is_button, is_api, path, sort_order) VALUES
-- 系统管理权限
('系统管理', 'system:manage', 'system', 'manage', 'system', '系统管理权限', true, false, false, '/system', 1000),
('用户管理', 'user:manage', 'user', 'manage', 'system', '用户管理权限', true, false, false, '/system/user', 1100),
('用户查看', 'user:read', 'user', 'read', 'system', '查看用户信息', false, true, true, '/system/user/read', 1101),
('用户创建', 'user:create', 'user', 'create', 'system', '创建用户', false, true, true, '/system/user/create', 1102),
('用户编辑', 'user:update', 'user', 'update', 'system', '编辑用户信息', false, true, true, '/system/user/update', 1103),
('用户删除', 'user:delete', 'user', 'delete', 'system', '删除用户', false, true, true, '/system/user/delete', 1104),
('角色管理', 'role:manage', 'role', 'manage', 'system', '角色管理权限', true, false, false, '/system/role', 1200),
('角色查看', 'role:read', 'role', 'read', 'system', '查看角色信息', false, true, true, '/system/role/read', 1201),
('角色创建', 'role:create', 'role', 'create', 'system', '创建角色', false, true, true, '/system/role/create', 1202),
('角色编辑', 'role:update', 'role', 'update', 'system', '编辑角色信息', false, true, true, '/system/role/update', 1203),
('角色删除', 'role:delete', 'role', 'delete', 'system', '删除角色', false, true, true, '/system/role/delete', 1204),
('权限分配', 'permission:assign', 'permission', 'assign', 'system', '分配权限', false, true, true, '/system/permission/assign', 1205),
('组织管理', 'organization:manage', 'organization', 'manage', 'system', '组织管理权限', true, false, false, '/system/organization', 1300),
('组织查看', 'organization:read', 'organization', 'read', 'system', '查看组织信息', false, true, true, '/system/organization/read', 1301),
('组织创建', 'organization:create', 'organization', 'create', 'system', '创建组织', false, true, true, '/system/organization/create', 1302),
('组织编辑', 'organization:update', 'organization', 'update', 'system', '编辑组织信息', false, true, true, '/system/organization/update', 1303),
('组织删除', 'organization:delete', 'organization', 'delete', 'system', '删除组织', false, true, true, '/system/organization/delete', 1304),
-- 仪表盘权限
('仪表盘', 'dashboard:read', 'dashboard', 'read', 'system', '访问系统仪表盘', true, false, true, '/dashboard', 100),
-- 系统日志权限
('系统日志', 'log:read', 'log', 'read', 'system', '查看系统日志', true, true, true, '/system/log', 1400),
-- 系统配置权限
('系统配置', 'config:manage', 'config', 'manage', 'system', '系统配置管理', true, false, false, '/system/config', 1500),
('配置查看', 'config:read', 'config', 'read', 'system', '查看系统配置', false, true, true, '/system/config/read', 1501),
('配置更新', 'config:update', 'config', 'update', 'system', '更新系统配置', false, true, true, '/system/config/update', 1502);

-- 创建系统默认角色
INSERT INTO roles (name, code, description, is_system_role, sort_order) VALUES
('超级管理员', 'SUPER_ADMIN', '系统超级管理员，拥有所有权限', true, 1),
('系统管理员', 'ADMIN', '系统管理员，拥有大部分管理权限', true, 2),
('普通用户', 'USER', '普通用户，基础权限', true, 10),
('游客', 'GUEST', '游客用户，只读权限', true, 99);

-- 为超级管理员角色分配所有权限
INSERT INTO role_permissions (role_id, permission_id, assigned_by) 
SELECT 1, id, NULL FROM permissions;

-- 为系统管理员角色分配管理权限(除了超级管理员专用权限)
INSERT INTO role_permissions (role_id, permission_id, assigned_by) 
SELECT 2, id, NULL FROM permissions 
WHERE code NOT IN ('system:manage', 'config:update');

-- 为普通用户角色分配基础权限  
INSERT INTO role_permissions (role_id, permission_id, assigned_by) 
SELECT 3, id, NULL FROM permissions 
WHERE code IN ('dashboard:read', 'user:read', 'organization:read');

-- 为游客角色分配只读权限
INSERT INTO role_permissions (role_id, permission_id, assigned_by) 
SELECT 4, id, NULL FROM permissions 
WHERE code IN ('dashboard:read');

-- 创建默认超级管理员用户
-- 密码: admin123 (请在生产环境中修改)
INSERT INTO users (username, email, password_hash, salt, first_name, last_name, is_enabled, created_by) 
VALUES (
    'admin', 
    'admin@example.com', 
    '$2b$12$LQv3c1yqBWVHxkd0LHAkCOYz6TtxMQJqhN8/LewLHoQEkPKz7gOoK',  -- bcrypt hash of 'admin123'
    'salt123', 
    'Super', 
    'Admin', 
    true, 
    NULL
);

-- 为超级管理员分配角色
INSERT INTO user_roles (user_id, role_id, assigned_by) VALUES (1, 1, NULL);

-- 将超级管理员加入默认组织
INSERT INTO user_organizations (user_id, org_id, position, is_primary) VALUES (1, 1, 'CEO', true);

-- ================================================================================================
-- 性能优化
-- ================================================================================================

-- 创建复合索引优化权限验证查询
CREATE INDEX idx_user_roles_permission_check ON user_roles(user_id, is_active, expires_at);
CREATE INDEX idx_role_permissions_check ON role_permissions(role_id, is_granted, permission_id);

-- 创建部分索引优化查询
CREATE INDEX idx_users_active ON users(id) WHERE is_enabled = true AND deleted_at IS NULL;
CREATE INDEX idx_roles_active ON roles(id) WHERE is_enabled = true AND deleted_at IS NULL;

-- 创建GIN索引用于JSONB字段
CREATE INDEX idx_operation_logs_details_gin ON operation_logs USING GIN (details);
CREATE INDEX idx_operation_logs_request_params_gin ON operation_logs USING GIN (request_params);

-- ================================================================================================
-- 数据完整性检查
-- ================================================================================================

-- 检查是否所有表都创建成功
DO $$
DECLARE
    table_count INTEGER;
BEGIN
    SELECT COUNT(*) INTO table_count 
    FROM information_schema.tables 
    WHERE table_schema = 'public' 
        AND table_type = 'BASE TABLE'
        AND table_name IN (
            'users', 'roles', 'permissions', 'user_roles', 'role_permissions',
            'organizations', 'user_organizations', 'system_configs', 
            'operation_logs', 'user_sessions'
        );
    
    IF table_count = 10 THEN
        RAISE NOTICE '✅ 所有数据表创建成功 (10/10)';
    ELSE
        RAISE WARNING '⚠️  数据表创建不完整 (%/10)', table_count;
    END IF;
END $$;

-- 检查初始化数据
DO $$
DECLARE
    config_count INTEGER;
    permission_count INTEGER;
    role_count INTEGER;
    user_count INTEGER;
BEGIN
    SELECT COUNT(*) INTO config_count FROM system_configs;
    SELECT COUNT(*) INTO permission_count FROM permissions;
    SELECT COUNT(*) INTO role_count FROM roles;
    SELECT COUNT(*) INTO user_count FROM users;
    
    RAISE NOTICE '📊 初始化数据统计:';
    RAISE NOTICE '   - 系统配置: % 条', config_count;
    RAISE NOTICE '   - 权限数据: % 条', permission_count;
    RAISE NOTICE '   - 系统角色: % 条', role_count;
    RAISE NOTICE '   - 用户账号: % 条', user_count;
END $$;

-- ================================================================================================
-- 提交事务
-- ================================================================================================
COMMIT;

-- ================================================================================================
-- 说明文档
-- ================================================================================================

COMMENT ON DATABASE erp_go IS 'ERP系统用户权限管理数据库';

-- 表注释
COMMENT ON TABLE users IS '用户表：存储用户基本信息，支持分库分表';
COMMENT ON TABLE roles IS '角色表：存储角色定义，支持系统预定义和自定义角色';
COMMENT ON TABLE permissions IS '权限表：存储权限定义，功能级粒度控制';
COMMENT ON TABLE user_roles IS '用户角色关联表：支持用户绑定多个角色，权限为并集';
COMMENT ON TABLE role_permissions IS '角色权限关联表：角色与权限的多对多关系';
COMMENT ON TABLE organizations IS '组织表：支持无限层级的组织架构';
COMMENT ON TABLE user_organizations IS '用户组织关联表：用户单组织归属';
COMMENT ON TABLE system_configs IS '系统配置表：存储系统运行时配置';
COMMENT ON TABLE operation_logs IS '操作日志表：记录用户操作审计日志，支持按时间分片';
COMMENT ON TABLE user_sessions IS '用户会话表：JWT令牌和登录会话管理';

-- 视图注释
COMMENT ON VIEW user_permissions_view IS '用户权限视图：优化权限验证查询性能';

-- 函数注释
COMMENT ON FUNCTION check_user_permission IS '权限验证函数：检查用户是否具有指定权限';
COMMENT ON FUNCTION get_user_permissions IS '获取用户权限函数：返回用户的所有权限列表';
COMMENT ON FUNCTION update_organization_path IS '组织路径更新函数：自动维护组织层级路径';

-- ================================================================================================
-- 初始化完成提示
-- ================================================================================================

DO $$
BEGIN
    RAISE NOTICE '';
    RAISE NOTICE '🎉 ERP用户权限管理系统数据库初始化完成！';
    RAISE NOTICE '';
    RAISE NOTICE '📋 系统信息:';
    RAISE NOTICE '   - 数据库版本: PostgreSQL %', version();
    RAISE NOTICE '   - 初始化时间: %', CURRENT_TIMESTAMP;
    RAISE NOTICE '   - 时区设置: %', CURRENT_SETTING('timezone');
    RAISE NOTICE '';
    RAISE NOTICE '🔐 默认管理员账号:';
    RAISE NOTICE '   - 用户名: admin';
    RAISE NOTICE '   - 密码: admin123';
    RAISE NOTICE '   - 邮箱: admin@example.com';
    RAISE NOTICE '   ⚠️  请在生产环境中修改默认密码！';
    RAISE NOTICE '';
    RAISE NOTICE '📚 使用说明:';
    RAISE NOTICE '   - 权限验证: SELECT check_user_permission(user_id, ''permission:code'');';
    RAISE NOTICE '   - 获取用户权限: SELECT * FROM get_user_permissions(user_id);';
    RAISE NOTICE '   - 权限视图查询: SELECT * FROM user_permissions_view WHERE user_id = ?;';
    RAISE NOTICE '';
    RAISE NOTICE '✨ 数据库初始化成功，可以开始使用ERP用户权限系统！';
    RAISE NOTICE '';
END $$;