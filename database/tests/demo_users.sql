-- ================================================================
-- 演示用户测试数据
-- 用于开发和测试环境的用户数据
-- ================================================================

-- 插入演示用户
INSERT OR IGNORE INTO users (username, email, password, first_name, last_name, phone, is_active) VALUES 
('demo_user1', 'demo1@erp.com', '$2a$10$rZF4vL4Td1rY8vZNzF.HruFjFcF5ZzXJKWX3yYfZ9YuY5KyTk2.Zi', '张', '三', '13800138001', true),
('demo_user2', 'demo2@erp.com', '$2a$10$rZF4vL4Td1rY8vZNzF.HruFjFcF5ZzXJKWX3yYfZ9YuY5KyTk2.Zi', '李', '四', '13800138002', true),
('demo_user3', 'demo3@erp.com', '$2a$10$rZF4vL4Td1rY8vZNzF.HruFjFcF5ZzXJKWX3yYfZ9YuY5KyTk2.Zi', '王', '五', '13800138003', true),
('demo_manager', 'manager@erp.com', '$2a$10$rZF4vL4Td1rY8vZNzF.HruFjFcF5ZzXJKWX3yYfZ9YuY5KyTk2.Zi', '经理', '用户', '13800138999', true);

-- 为演示用户分配角色
INSERT OR IGNORE INTO user_roles (user_id, role_id) 
SELECT u.id, r.id FROM users u, roles r 
WHERE u.username IN ('demo_user1', 'demo_user2', 'demo_user3') AND r.code = 'USER';

-- 为演示经理分配管理员角色
INSERT OR IGNORE INTO user_roles (user_id, role_id) 
SELECT u.id, r.id FROM users u, roles r 
WHERE u.username = 'demo_manager' AND r.code = 'ADMIN';

-- 插入更多测试组织
INSERT OR IGNORE INTO organizations (name, code, description, level, path) VALUES 
('研发部', 'RD', '研发部门', 2, '/1/2'),
('销售部', 'SALES', '销售部门', 2, '/1/3'),
('市场部', 'MARKETING', '市场部门', 2, '/1/4');

-- 为用户分配组织
INSERT OR IGNORE INTO user_organizations (user_id, organization_id, is_primary) 
SELECT u.id, o.id, true FROM users u, organizations o 
WHERE u.username = 'demo_user1' AND o.code = 'RD';

INSERT OR IGNORE INTO user_organizations (user_id, organization_id, is_primary) 
SELECT u.id, o.id, true FROM users u, organizations o 
WHERE u.username = 'demo_user2' AND o.code = 'SALES';

INSERT OR IGNORE INTO user_organizations (user_id, organization_id, is_primary) 
SELECT u.id, o.id, true FROM users u, organizations o 
WHERE u.username = 'demo_user3' AND o.code = 'MARKETING';