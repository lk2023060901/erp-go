-- ================================================================
-- 系统权限数据
-- ================================================================
-- 插入系统预定义的权限

\echo '插入系统管理权限...'

-- 系统管理模块权限
INSERT INTO permissions (id, name, code, resource, action, module, description, is_menu, is_button, is_api, menu_url, api_path, api_method, level, sort_order, created_at) VALUES

-- 系统管理菜单
(1, '系统管理', 'system', 'system', 'view', 'system', '系统管理菜单', true, false, false, '/system', NULL, NULL, 1, 100, CURRENT_TIMESTAMP),

-- 用户管理
(10, '用户管理', 'system.user', 'user', 'view', 'system', '用户管理菜单', true, false, false, '/system/user', NULL, NULL, 2, 110, CURRENT_TIMESTAMP),
(11, '查看用户', 'system.user.view', 'user', 'read', 'system', '查看用户列表和详情', false, true, true, NULL, '/v1/users', 'GET', 3, 111, CURRENT_TIMESTAMP),
(12, '创建用户', 'system.user.create', 'user', 'write', 'system', '创建新用户', false, true, true, NULL, '/v1/users', 'POST', 3, 112, CURRENT_TIMESTAMP),
(13, '编辑用户', 'system.user.edit', 'user', 'write', 'system', '编辑用户信息', false, true, true, NULL, '/v1/users/*', 'PUT', 3, 113, CURRENT_TIMESTAMP),
(14, '删除用户', 'system.user.delete', 'user', 'delete', 'system', '删除用户', false, true, true, NULL, '/v1/users/*', 'DELETE', 3, 114, CURRENT_TIMESTAMP),
(15, '重置密码', 'system.user.reset_password', 'user', 'write', 'system', '重置用户密码', false, true, true, NULL, '/v1/users/*/reset-password', 'POST', 3, 115, CURRENT_TIMESTAMP),

-- 角色管理
(20, '角色管理', 'system.role', 'role', 'view', 'system', '角色管理菜单', true, false, false, '/system/role', NULL, NULL, 2, 120, CURRENT_TIMESTAMP),
(21, '查看角色', 'system.role.view', 'role', 'read', 'system', '查看角色列表和详情', false, true, true, NULL, '/v1/roles', 'GET', 3, 121, CURRENT_TIMESTAMP),
(22, '创建角色', 'system.role.create', 'role', 'write', 'system', '创建新角色', false, true, true, NULL, '/v1/roles', 'POST', 3, 122, CURRENT_TIMESTAMP),
(23, '编辑角色', 'system.role.edit', 'role', 'write', 'system', '编辑角色信息', false, true, true, NULL, '/v1/roles/*', 'PUT', 3, 123, CURRENT_TIMESTAMP),
(24, '删除角色', 'system.role.delete', 'role', 'delete', 'system', '删除角色', false, true, true, NULL, '/v1/roles/*', 'DELETE', 3, 124, CURRENT_TIMESTAMP),
(25, '分配权限', 'system.role.assign_permission', 'role', 'write', 'system', '为角色分配权限', false, true, true, NULL, '/v1/roles/*/permissions', 'POST', 3, 125, CURRENT_TIMESTAMP),

-- 权限管理
(30, '权限管理', 'system.permission', 'permission', 'view', 'system', '权限管理菜单', true, false, false, '/system/permission', NULL, NULL, 2, 130, CURRENT_TIMESTAMP),
(31, '查看权限', 'system.permission.view', 'permission', 'read', 'system', '查看权限列表和详情', false, true, true, NULL, '/v1/permissions', 'GET', 3, 131, CURRENT_TIMESTAMP),
(32, '创建权限', 'system.permission.create', 'permission', 'write', 'system', '创建新权限', false, true, true, NULL, '/v1/permissions', 'POST', 3, 132, CURRENT_TIMESTAMP),
(33, '编辑权限', 'system.permission.edit', 'permission', 'write', 'system', '编辑权限信息', false, true, true, NULL, '/v1/permissions/*', 'PUT', 3, 133, CURRENT_TIMESTAMP),
(34, '删除权限', 'system.permission.delete', 'permission', 'delete', 'system', '删除权限', false, true, true, NULL, '/v1/permissions/*', 'DELETE', 3, 134, CURRENT_TIMESTAMP),

-- 组织管理
(40, '组织管理', 'system.organization', 'organization', 'view', 'system', '组织管理菜单', true, false, false, '/system/organization', NULL, NULL, 2, 140, CURRENT_TIMESTAMP),
(41, '查看组织', 'system.organization.view', 'organization', 'read', 'system', '查看组织架构', false, true, true, NULL, '/v1/organizations', 'GET', 3, 141, CURRENT_TIMESTAMP),
(42, '创建组织', 'system.organization.create', 'organization', 'write', 'system', '创建新组织', false, true, true, NULL, '/v1/organizations', 'POST', 3, 142, CURRENT_TIMESTAMP),
(43, '编辑组织', 'system.organization.edit', 'organization', 'write', 'system', '编辑组织信息', false, true, true, NULL, '/v1/organizations/*', 'PUT', 3, 143, CURRENT_TIMESTAMP),
(44, '删除组织', 'system.organization.delete', 'organization', 'delete', 'system', '删除组织', false, true, true, NULL, '/v1/organizations/*', 'DELETE', 3, 144, CURRENT_TIMESTAMP),

-- 系统配置
(50, '系统配置', 'system.config', 'config', 'view', 'system', '系统配置菜单', true, false, false, '/system/config', NULL, NULL, 2, 150, CURRENT_TIMESTAMP),
(51, '查看配置', 'system.config.view', 'config', 'read', 'system', '查看系统配置', false, true, true, NULL, '/v1/configs', 'GET', 3, 151, CURRENT_TIMESTAMP),
(52, '修改配置', 'system.config.edit', 'config', 'write', 'system', '修改系统配置', false, true, true, NULL, '/v1/configs/*', 'PUT', 3, 152, CURRENT_TIMESTAMP),

-- 操作日志
(60, '操作日志', 'system.log', 'log', 'view', 'system', '操作日志菜单', true, false, false, '/system/log', NULL, NULL, 2, 160, CURRENT_TIMESTAMP),
(61, '查看日志', 'system.log.view', 'log', 'read', 'system', '查看操作日志', false, true, true, NULL, '/v1/logs', 'GET', 3, 161, CURRENT_TIMESTAMP);

-- 设置权限父子关系
UPDATE permissions SET parent_id = 1 WHERE id IN (10, 20, 30, 40, 50, 60);
UPDATE permissions SET parent_id = 10 WHERE id IN (11, 12, 13, 14, 15);
UPDATE permissions SET parent_id = 20 WHERE id IN (21, 22, 23, 24, 25);
UPDATE permissions SET parent_id = 30 WHERE id IN (31, 32, 33, 34);
UPDATE permissions SET parent_id = 40 WHERE id IN (41, 42, 43, 44);
UPDATE permissions SET parent_id = 50 WHERE id IN (51, 52);
UPDATE permissions SET parent_id = 60 WHERE id IN (61);

\echo '插入个人中心权限...'

-- 个人中心权限
INSERT INTO permissions (id, name, code, resource, action, module, description, is_menu, is_button, is_api, menu_url, api_path, api_method, level, sort_order, created_at) VALUES

-- 个人中心菜单
(100, '个人中心', 'profile', 'profile', 'view', 'profile', '个人中心菜单', true, false, false, '/profile', NULL, NULL, 1, 200, CURRENT_TIMESTAMP),
(101, '查看个人信息', 'profile.view', 'profile', 'read', 'profile', '查看个人信息', false, true, true, NULL, '/v1/profile', 'GET', 2, 201, CURRENT_TIMESTAMP),
(102, '修改个人信息', 'profile.edit', 'profile', 'write', 'profile', '修改个人信息', false, true, true, NULL, '/v1/profile', 'PUT', 2, 202, CURRENT_TIMESTAMP),
(103, '修改密码', 'profile.change_password', 'profile', 'write', 'profile', '修改登录密码', false, true, true, NULL, '/v1/profile/password', 'PUT', 2, 203, CURRENT_TIMESTAMP),
(104, '上传头像', 'profile.upload_avatar', 'profile', 'write', 'profile', '上传个人头像', false, true, true, NULL, '/v1/profile/avatar', 'POST', 2, 204, CURRENT_TIMESTAMP);

-- 设置权限父子关系
UPDATE permissions SET parent_id = 100 WHERE id IN (101, 102, 103, 104);

-- 重置权限表序列
SELECT setval('permissions_id_seq', (SELECT MAX(id) FROM permissions));

\echo '系统权限插入完成！'