-- ================================================================
-- 核心系统种子数据
-- 包含：基础角色、权限、组织、管理员用户等必要数据
-- ================================================================

-- 插入默认角色
INSERT OR IGNORE INTO roles (name, code, description, is_system_role) VALUES 
('超级管理员', 'SUPER_ADMIN', '系统超级管理员', true),
('管理员', 'ADMIN', '系统管理员', true),
('普通用户', 'USER', '普通用户', true);

-- 插入默认权限
INSERT OR IGNORE INTO permissions (name, code, resource, action, module, description, is_menu, is_api) VALUES
('用户管理', 'user.list', 'user', 'list', 'user', '查看用户列表', true, true),
('创建用户', 'user.create', 'user', 'create', 'user', '创建新用户', false, true),
('编辑用户', 'user.edit', 'user', 'edit', 'user', '编辑用户信息', false, true),
('删除用户', 'user.delete', 'user', 'delete', 'user', '删除用户', false, true),
('角色管理', 'role.list', 'role', 'list', 'role', '查看角色列表', true, true),
('创建角色', 'role.create', 'role', 'create', 'role', '创建新角色', false, true),
('编辑角色', 'role.edit', 'role', 'edit', 'role', '编辑角色信息', false, true),
('删除角色', 'role.delete', 'role', 'delete', 'role', '删除角色', false, true);

-- 插入默认组织
INSERT OR IGNORE INTO organizations (name, code, description, level, path) VALUES 
('总部', 'HQ', '公司总部', 1, '/1');

-- 创建默认管理员用户 (密码: admin123)
INSERT OR IGNORE INTO users (username, email, password, first_name, last_name) VALUES 
('admin', 'admin@erp.com', '$2a$10$rZF4vL4Td1rY8vZNzF.HruFjFcF5ZzXJKWX3yYfZ9YuY5KyTk2.Zi', 'Admin', 'User');

-- 为管理员分配超级管理员角色
INSERT OR IGNORE INTO user_roles (user_id, role_id) 
SELECT u.id, r.id FROM users u, roles r 
WHERE u.username = 'admin' AND r.code = 'SUPER_ADMIN';

-- 为超级管理员角色分配所有权限
INSERT OR IGNORE INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id FROM roles r, permissions p 
WHERE r.code = 'SUPER_ADMIN';