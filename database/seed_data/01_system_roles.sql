-- ================================================================
-- 系统默认角色数据
-- ================================================================
-- 插入系统预定义的角色

\echo '插入系统默认角色...'

INSERT INTO roles (id, name, code, description, is_system_role, is_enabled, sort_order, created_at) VALUES
(1, '超级管理员', 'SUPER_ADMIN', '系统超级管理员，拥有所有权限', true, true, 100, CURRENT_TIMESTAMP),
(2, '系统管理员', 'ADMIN', '系统管理员，拥有大部分管理权限', true, true, 90, CURRENT_TIMESTAMP),
(3, '部门管理员', 'DEPT_ADMIN', '部门管理员，拥有部门内管理权限', true, true, 80, CURRENT_TIMESTAMP),
(4, '普通用户', 'USER', '普通用户，基础权限', true, true, 50, CURRENT_TIMESTAMP),
(5, '访客用户', 'GUEST', '访客用户，只读权限', true, true, 10, CURRENT_TIMESTAMP);

-- 重置角色表序列
SELECT setval('roles_id_seq', (SELECT MAX(id) FROM roles));

\echo '系统默认角色插入完成！'