-- ================================================================
-- 高级权限系统种子数据
-- 包含：基础文档类型、权限规则、字段权限级别等数据
-- ================================================================

-- 插入基础DocType数据
INSERT OR IGNORE INTO doc_types (name, label, module, description) VALUES 
('User', '用户', 'Core', '系统用户文档类型'),
('Role', '角色', 'Core', '系统角色文档类型'),
('DocType', '文档类型', 'Core', 'DocType管理'),
('Permission Rule', '权限规则', 'Core', '权限规则管理'),
('User Permission', '用户权限', 'Core', '用户权限管理'),
('System Settings', '系统设置', 'Core', '系统配置管理');

-- ================================================================
-- 插入基础权限规则 (为超级管理员角色)
-- ================================================================

-- 为超级管理员角色(ID=1)分配所有核心DocType的完整权限
INSERT OR IGNORE INTO permission_rules (
    role, document_type, permission_level,
    read, write, [create], [delete], submit, cancel, amend,
    print, email, import, export, share, report
) VALUES 
(1, 'User', 0, 1, 1, 1, 1, 0, 0, 0, 1, 1, 1, 1, 1, 1),
(1, 'Role', 0, 1, 1, 1, 1, 0, 0, 0, 1, 1, 1, 1, 1, 1),
(1, 'DocType', 0, 1, 1, 1, 1, 0, 0, 0, 1, 1, 1, 1, 1, 1),
(1, 'Permission Rule', 0, 1, 1, 1, 1, 0, 0, 0, 1, 1, 1, 1, 1, 1),
(1, 'User Permission', 0, 1, 1, 1, 1, 0, 0, 0, 1, 1, 1, 1, 1, 1),
(1, 'System Settings', 0, 1, 1, 1, 1, 0, 0, 0, 1, 1, 1, 1, 1, 1);

-- ================================================================
-- 插入基础字段权限级别
-- ================================================================

-- User DocType的字段权限级别
INSERT OR IGNORE INTO field_permission_levels (document_type, field_name, field_label, permission_level, field_type, is_mandatory) VALUES 
('User', 'username', '用户名', 1, 'Data', true),
('User', 'email', '邮箱', 1, 'Data', true),
('User', 'first_name', '名', 2, 'Data', true),
('User', 'last_name', '姓', 2, 'Data', true),
('User', 'password', '密码', 5, 'Password', true),
('User', 'phone', '电话', 3, 'Data', false),
('User', 'is_active', '启用状态', 4, 'Check', false);

-- Role DocType的字段权限级别
INSERT OR IGNORE INTO field_permission_levels (document_type, field_name, field_label, permission_level, field_type, is_mandatory) VALUES 
('Role', 'name', '角色名称', 1, 'Data', true),
('Role', 'code', '角色代码', 1, 'Data', true),
('Role', 'description', '描述', 2, 'Text', false),
('Role', 'is_system_role', '系统角色', 4, 'Check', false);

-- ================================================================
-- 创建权限检查视图
-- ================================================================

-- 权限检查视图 (简化权限查询)
CREATE VIEW IF NOT EXISTS user_permission_check AS
SELECT 
    ur.user_id,
    pr.document_type,
    pr.read,
    pr.write,
    pr.[create],
    pr.[delete],
    pr.submit,
    pr.cancel,
    pr.amend,
    pr.print,
    pr.email,
    pr.import,
    pr.export,
    pr.share,
    pr.report
FROM user_roles ur
INNER JOIN permission_rules pr ON ur.role_id = pr.role
WHERE pr.permission_level = 0;