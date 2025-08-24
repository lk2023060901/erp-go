-- ================================================================
-- 默认组织数据
-- ================================================================
-- 创建默认的组织架构

\echo '创建默认组织架构...'

INSERT INTO organizations (id, parent_id, name, code, type, description, level, path, is_enabled, sort_order, created_at) VALUES
(1, NULL, '总公司', 'ROOT', 'company', '公司总部', 1, '/1/', true, 100, CURRENT_TIMESTAMP),
(2, 1, '技术部', 'TECH', 'department', '技术研发部门', 2, '/1/2/', true, 110, CURRENT_TIMESTAMP),
(3, 1, '人力资源部', 'HR', 'department', '人力资源管理部门', 2, '/1/3/', true, 120, CURRENT_TIMESTAMP),
(4, 1, '财务部', 'FINANCE', 'department', '财务管理部门', 2, '/1/4/', true, 130, CURRENT_TIMESTAMP),
(5, 2, '前端开发组', 'FRONTEND', 'team', '前端开发团队', 3, '/1/2/5/', true, 111, CURRENT_TIMESTAMP),
(6, 2, '后端开发组', 'BACKEND', 'team', '后端开发团队', 3, '/1/2/6/', true, 112, CURRENT_TIMESTAMP),
(7, 2, '测试组', 'QA', 'team', '质量保证团队', 3, '/1/2/7/', true, 113, CURRENT_TIMESTAMP);

-- 重置组织表序列
SELECT setval('organizations_id_seq', (SELECT MAX(id) FROM organizations));

\echo '默认组织架构创建完成！'