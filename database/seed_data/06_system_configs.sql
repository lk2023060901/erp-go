-- ================================================================
-- 系统配置数据
-- ================================================================
-- 插入系统默认配置参数

\echo '插入系统配置数据...'

INSERT INTO system_configs (config_key, config_value, config_type, description, is_system, group_name, sort_order, is_enabled) VALUES

-- 系统基础配置
('system.name', 'ERP管理系统', 'string', '系统名称', true, 'system', 1, true),
('system.version', '1.0.0', 'string', '系统版本', true, 'system', 2, true),
('system.company', '某某科技有限公司', 'string', '公司名称', false, 'system', 3, true),
('system.logo_url', '/static/images/logo.png', 'string', '系统Logo URL', false, 'system', 4, true),
('system.copyright', '© 2024 某某科技有限公司', 'string', '版权信息', false, 'system', 5, true),

-- 用户认证配置
('auth.login_max_attempts', '5', 'integer', '登录最大尝试次数', true, 'auth', 10, true),
('auth.login_lock_duration', '30', 'integer', '登录锁定时间（分钟）', true, 'auth', 11, true),
('auth.password_min_length', '8', 'integer', '密码最小长度', true, 'auth', 12, true),
('auth.password_require_uppercase', 'true', 'boolean', '密码必须包含大写字母', true, 'auth', 13, true),
('auth.password_require_lowercase', 'true', 'boolean', '密码必须包含小写字母', true, 'auth', 14, true),
('auth.password_require_number', 'true', 'boolean', '密码必须包含数字', true, 'auth', 15, true),
('auth.password_require_special', 'false', 'boolean', '密码必须包含特殊字符', true, 'auth', 16, true),
('auth.password_expiry_days', '90', 'integer', '密码过期天数（0为不过期）', true, 'auth', 17, true),
('auth.session_timeout', '1440', 'integer', '会话超时时间（分钟）', true, 'auth', 18, true),
('auth.two_factor_enabled', 'false', 'boolean', '是否启用双重认证', true, 'auth', 19, true),

-- JWT配置
('jwt.secret_key', 'your-secret-key-change-in-production', 'string', 'JWT签名密钥', true, 'jwt', 20, true),
('jwt.access_token_expire', '1440', 'integer', '访问令牌过期时间（分钟）', true, 'jwt', 21, true),
('jwt.refresh_token_expire', '10080', 'integer', '刷新令牌过期时间（分钟）', true, 'jwt', 22, true),
('jwt.issuer', 'erp-system', 'string', 'JWT签发者', true, 'jwt', 23, true),

-- 邮件配置
('email.smtp_host', 'smtp.example.com', 'string', 'SMTP服务器地址', false, 'email', 30, true),
('email.smtp_port', '587', 'integer', 'SMTP端口', false, 'email', 31, true),
('email.smtp_username', '', 'string', 'SMTP用户名', false, 'email', 32, true),
('email.smtp_password', '', 'string', 'SMTP密码', true, 'email', 33, true),
('email.from_address', 'noreply@example.com', 'string', '发送邮箱地址', false, 'email', 34, true),
('email.from_name', 'ERP系统', 'string', '发送者名称', false, 'email', 35, true),

-- 文件上传配置
('upload.max_file_size', '10485760', 'integer', '最大文件大小（字节）', true, 'upload', 40, true),
('upload.allowed_extensions', 'jpg,jpeg,png,gif,pdf,doc,docx,xls,xlsx', 'string', '允许的文件扩展名', true, 'upload', 41, true),
('upload.storage_path', '/uploads', 'string', '文件存储路径', true, 'upload', 42, true),
('upload.avatar_max_size', '2097152', 'integer', '头像最大大小（字节）', true, 'upload', 43, true),

-- 日志配置
('log.level', 'info', 'string', '日志级别', true, 'log', 50, true),
('log.retention_days', '90', 'integer', '日志保留天数', true, 'log', 51, true),
('log.max_file_size', '100', 'integer', '单个日志文件最大大小（MB）', true, 'log', 52, true),

-- 缓存配置
('cache.default_expire', '3600', 'integer', '默认缓存过期时间（秒）', true, 'cache', 60, true),
('cache.user_permission_expire', '1800', 'integer', '用户权限缓存过期时间（秒）', true, 'cache', 61, true),
('cache.system_config_expire', '86400', 'integer', '系统配置缓存过期时间（秒）', true, 'cache', 62, true),

-- 接口配置
('api.rate_limit_enabled', 'true', 'boolean', '是否启用接口限流', true, 'api', 70, true),
('api.rate_limit_requests', '1000', 'integer', '每分钟请求限制', true, 'api', 71, true),
('api.cors_enabled', 'true', 'boolean', '是否启用跨域请求', true, 'api', 72, true),
('api.cors_origins', '*', 'string', '允许的跨域来源', true, 'api', 73, true),

-- 业务配置
('business.user_registration_enabled', 'false', 'boolean', '是否允许用户自助注册', false, 'business', 80, true),
('business.email_verification_required', 'true', 'boolean', '是否需要邮箱验证', false, 'business', 81, true),
('business.default_user_role', 'USER', 'string', '默认用户角色', false, 'business', 82, true),
('business.default_organization', '1', 'string', '默认组织ID', false, 'business', 83, true);

\echo '系统配置数据插入完成！'