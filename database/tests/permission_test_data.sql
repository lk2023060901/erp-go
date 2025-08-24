-- ================================================================
-- 权限系统测试数据
-- 用于测试高级权限系统功能的数据
-- ================================================================

-- 插入更多测试角色
INSERT OR IGNORE INTO roles (name, code, description, is_system_role) VALUES 
('财务专员', 'FINANCE_STAFF', '财务部门专员', false),
('人事专员', 'HR_STAFF', '人事部门专员', false),
('客服专员', 'SERVICE_STAFF', '客服部门专员', false);

-- 插入更多测试权限
INSERT OR IGNORE INTO permissions (name, code, resource, action, module, description, is_menu, is_api) VALUES
('客户管理', 'customer.list', 'customer', 'list', 'customer', '查看客户列表', true, true),
('创建客户', 'customer.create', 'customer', 'create', 'customer', '创建新客户', false, true),
('编辑客户', 'customer.edit', 'customer', 'edit', 'customer', '编辑客户信息', false, true),
('删除客户', 'customer.delete', 'customer', 'delete', 'customer', '删除客户', false, true),
('财务报表', 'finance.report', 'finance', 'report', 'finance', '查看财务报表', true, true),
('人事管理', 'hr.list', 'hr', 'list', 'hr', '人事管理', true, true);

-- 为新角色分配权限
INSERT OR IGNORE INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id FROM roles r, permissions p 
WHERE r.code = 'FINANCE_STAFF' AND p.code IN ('finance.report', 'customer.list');

INSERT OR IGNORE INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id FROM roles r, permissions p 
WHERE r.code = 'HR_STAFF' AND p.code IN ('hr.list', 'user.list', 'user.edit');

INSERT OR IGNORE INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id FROM roles r, permissions p 
WHERE r.code = 'SERVICE_STAFF' AND p.code IN ('customer.list', 'customer.edit');

-- 插入测试文档类型
INSERT OR IGNORE INTO doc_types (name, label, module, description, is_submittable, has_workflow) VALUES
('Customer', '客户', 'crm', '客户管理', false, false),
('Invoice', '发票', 'accounting', '财务发票', true, true),
('Purchase Order', '采购订单', 'procurement', '采购管理', true, true),
('Quotation', '报价单', 'sales', '销售报价', true, true);

-- 为Customer文档类型创建字段权限级别
INSERT OR IGNORE INTO field_permission_levels (doc_type, field_name, field_label, permission_level, field_type, is_mandatory) VALUES
-- Level 0: 基础客户信息
('Customer', 'name', '客户名称', 0, 'Data', true),
('Customer', 'code', '客户编码', 0, 'Data', true),
('Customer', 'contact_person', '联系人', 0, 'Data', false),
('Customer', 'phone', '电话', 0, 'Data', false),
('Customer', 'email', '邮箱', 0, 'Data', false),
('Customer', 'address', '地址', 0, 'Text', false),

-- Level 1: 商业信息
('Customer', 'credit_limit', '信用额度', 1, 'Currency', false),
('Customer', 'payment_terms', '付款条件', 1, 'Data', false),
('Customer', 'tax_id', '税号', 1, 'Data', false),

-- Level 2: 财务敏感信息
('Customer', 'outstanding_amount', '应收账款', 2, 'Currency', false),
('Customer', 'credit_score', '信用评分', 2, 'Int', false),
('Customer', 'risk_level', '风险级别', 2, 'Select', false);

-- 创建Customer文档的权限规则
INSERT OR IGNORE INTO permission_rules (role, document_type, permission_level, read, write, [create], [delete], print, email, export, report) VALUES
-- 超级管理员：所有权限
(1, 'Customer', 0, 1, 1, 1, 1, 1, 1, 1, 1),
(1, 'Customer', 1, 1, 1, 0, 0, 1, 1, 1, 1),
(1, 'Customer', 2, 1, 1, 0, 0, 1, 1, 1, 1),

-- 客服专员：只能查看和编辑基础信息
(6, 'Customer', 0, 1, 1, 1, 0, 1, 1, 0, 0),
(6, 'Customer', 1, 1, 0, 0, 0, 0, 0, 0, 0),

-- 财务专员：可以查看所有级别，编辑财务信息
(4, 'Customer', 0, 1, 0, 0, 0, 1, 1, 1, 1),
(4, 'Customer', 1, 1, 1, 0, 0, 1, 1, 1, 1),
(4, 'Customer', 2, 1, 1, 0, 0, 1, 1, 1, 1);

-- 创建用户权限示例（数据范围限制）
INSERT OR IGNORE INTO user_permissions (user_id, doc_type, doc_name, permission_value, applicable_for) VALUES
-- 限制某个用户只能看到特定客户
(4, 'Customer', 'CUST-001', '华为技术有限公司', NULL),
(4, 'Customer', 'CUST-002', '腾讯科技有限公司', NULL),
(5, 'Customer', 'CUST-003', '阿里巴巴集团', NULL),
(5, 'Customer', 'CUST-004', '百度在线网络技术有限公司', NULL);