-- ================================================================
-- 默认管理员用户数据
-- ================================================================
-- 创建默认的超级管理员账户

\echo '创建默认超级管理员用户...'

-- 插入默认超级管理员用户
INSERT INTO users (
    id, 
    username, 
    email, 
    first_name, 
    last_name, 
    password_hash, 
    salt, 
    status, 
    is_enabled, 
    created_at
) VALUES (
    1,
    'admin',
    'admin@example.com',
    'Super',
    'Admin',
    '$2a$10$p9dDHTgjk8FOFRK9Ju/yIeRICQ8XOdlFaFvjRM2WoOhHKQ2ILKCAi', -- 密码: Admin123@
    'random_salt_string',
    1,
    true,
    CURRENT_TIMESTAMP
);

-- 为管理员分配超级管理员角色
INSERT INTO user_roles (user_id, role_id, is_active, granted_at) VALUES
(1, 1, true, CURRENT_TIMESTAMP);

-- 将管理员加入总公司组织
INSERT INTO user_organizations (user_id, org_id, position, is_primary, is_leader, is_active, created_at) VALUES
(1, 1, '系统管理员', true, true, true, CURRENT_TIMESTAMP);

-- 重置用户表序列
SELECT setval('users_id_seq', (SELECT MAX(id) FROM users));

\echo '创建测试用户...'

-- 插入测试普通用户
INSERT INTO users (
    id, 
    username, 
    email, 
    first_name, 
    last_name, 
    phone,
    gender,
    password_hash, 
    salt, 
    status, 
    is_enabled, 
    created_at
) VALUES 
(2, 'user001', 'user001@example.com', '张', '三', '13800000001', 'M', '$2a$10$p9dDHTgjk8FOFRK9Ju/yIeRICQ8XOdlFaFvjRM2WoOhHKQ2ILKCAi', 'salt_user001', 1, true, CURRENT_TIMESTAMP),
(3, 'user002', 'user002@example.com', '李', '四', '13800000002', 'F', '$2a$10$p9dDHTgjk8FOFRK9Ju/yIeRICQ8XOdlFaFvjRM2WoOhHKQ2ILKCAi', 'salt_user002', 1, true, CURRENT_TIMESTAMP),
(4, 'manager01', 'manager01@example.com', '王', '五', '13800000003', 'M', '$2a$10$p9dDHTgjk8FOFRK9Ju/yIeRICQ8XOdlFaFvjRM2WoOhHKQ2ILKCAi', 'salt_manager01', 1, true, CURRENT_TIMESTAMP);

-- 为测试用户分配角色
INSERT INTO user_roles (user_id, role_id, is_active, granted_at) VALUES
(2, 4, true, CURRENT_TIMESTAMP), -- 普通用户角色
(3, 4, true, CURRENT_TIMESTAMP), -- 普通用户角色
(4, 3, true, CURRENT_TIMESTAMP); -- 部门管理员角色

-- 为测试用户分配组织
INSERT INTO user_organizations (user_id, org_id, position, is_primary, is_active, created_at) VALUES
(2, 5, '前端开发工程师', true, true, CURRENT_TIMESTAMP),
(3, 6, '后端开发工程师', true, true, CURRENT_TIMESTAMP),
(4, 2, '技术部经理', true, true, CURRENT_TIMESTAMP);

-- 更新组织负责人
UPDATE organizations SET leader_id = 4 WHERE id = 2; -- 技术部经理

-- 重置用户表序列
SELECT setval('users_id_seq', (SELECT MAX(id) FROM users));

\echo '默认用户数据创建完成！'
\echo ''
\echo '默认登录信息:'
\echo '超级管理员 - 用户名: admin, 密码: Admin123@'
\echo '普通用户 - 用户名: user001, 密码: Admin123@'
\echo '部门经理 - 用户名: manager01, 密码: Admin123@'