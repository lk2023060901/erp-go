-- ================================================================
-- 角色权限关联数据
-- ================================================================
-- 为系统默认角色分配权限

\echo '为超级管理员分配所有权限...'

-- 超级管理员拥有所有权限
INSERT INTO role_permissions (role_id, permission_id, is_granted, granted_at)
SELECT 1, id, true, CURRENT_TIMESTAMP
FROM permissions
WHERE deleted_at IS NULL;

\echo '为系统管理员分配管理权限...'

-- 系统管理员权限（除了删除用户和删除角色）
INSERT INTO role_permissions (role_id, permission_id, is_granted, granted_at) VALUES
-- 系统管理菜单
(2, 1, true, CURRENT_TIMESTAMP),
-- 用户管理（不包括删除）
(2, 10, true, CURRENT_TIMESTAMP),
(2, 11, true, CURRENT_TIMESTAMP),
(2, 12, true, CURRENT_TIMESTAMP),
(2, 13, true, CURRENT_TIMESTAMP),
(2, 15, true, CURRENT_TIMESTAMP),
-- 角色管理（不包括删除系统角色）
(2, 20, true, CURRENT_TIMESTAMP),
(2, 21, true, CURRENT_TIMESTAMP),
(2, 22, true, CURRENT_TIMESTAMP),
(2, 23, true, CURRENT_TIMESTAMP),
(2, 25, true, CURRENT_TIMESTAMP),
-- 权限管理（只读）
(2, 30, true, CURRENT_TIMESTAMP),
(2, 31, true, CURRENT_TIMESTAMP),
-- 组织管理
(2, 40, true, CURRENT_TIMESTAMP),
(2, 41, true, CURRENT_TIMESTAMP),
(2, 42, true, CURRENT_TIMESTAMP),
(2, 43, true, CURRENT_TIMESTAMP),
(2, 44, true, CURRENT_TIMESTAMP),
-- 系统配置
(2, 50, true, CURRENT_TIMESTAMP),
(2, 51, true, CURRENT_TIMESTAMP),
(2, 52, true, CURRENT_TIMESTAMP),
-- 操作日志
(2, 60, true, CURRENT_TIMESTAMP),
(2, 61, true, CURRENT_TIMESTAMP),
-- 个人中心
(2, 100, true, CURRENT_TIMESTAMP),
(2, 101, true, CURRENT_TIMESTAMP),
(2, 102, true, CURRENT_TIMESTAMP),
(2, 103, true, CURRENT_TIMESTAMP),
(2, 104, true, CURRENT_TIMESTAMP);

\echo '为部门管理员分配部门权限...'

-- 部门管理员权限（部分用户管理和基础功能）
INSERT INTO role_permissions (role_id, permission_id, is_granted, granted_at) VALUES
-- 用户管理（查看、编辑、重置密码）
(3, 11, true, CURRENT_TIMESTAMP),
(3, 13, true, CURRENT_TIMESTAMP),
(3, 15, true, CURRENT_TIMESTAMP),
-- 角色管理（只读）
(3, 21, true, CURRENT_TIMESTAMP),
-- 组织管理（查看、编辑）
(3, 41, true, CURRENT_TIMESTAMP),
(3, 43, true, CURRENT_TIMESTAMP),
-- 个人中心
(3, 100, true, CURRENT_TIMESTAMP),
(3, 101, true, CURRENT_TIMESTAMP),
(3, 102, true, CURRENT_TIMESTAMP),
(3, 103, true, CURRENT_TIMESTAMP),
(3, 104, true, CURRENT_TIMESTAMP);

\echo '为普通用户分配基础权限...'

-- 普通用户权限（个人中心）
INSERT INTO role_permissions (role_id, permission_id, is_granted, granted_at) VALUES
-- 个人中心
(4, 100, true, CURRENT_TIMESTAMP),
(4, 101, true, CURRENT_TIMESTAMP),
(4, 102, true, CURRENT_TIMESTAMP),
(4, 103, true, CURRENT_TIMESTAMP),
(4, 104, true, CURRENT_TIMESTAMP);

\echo '为访客用户分配只读权限...'

-- 访客用户权限（只读个人信息）
INSERT INTO role_permissions (role_id, permission_id, is_granted, granted_at) VALUES
-- 个人中心（只读）
(5, 100, true, CURRENT_TIMESTAMP),
(5, 101, true, CURRENT_TIMESTAMP);

\echo '角色权限关联数据插入完成！'