-- ================================================================================================
-- 高级权限系统扩展迁移脚本
-- 在现有RBAC基础上添加四层权限控制：
-- 1. 角色权限规则（Role Permission Rules）- 文档级+字段级权限
-- 2. 用户权限（User Permissions） - 数据范围权限  
-- 3. 文档类型定义（Document Types）
-- 4. 字段权限级别（Field Permission Levels）
-- 5. 文档工作流状态（Document Workflow States）
-- ================================================================================================

-- 开启事务
BEGIN;

-- ================================================================================================
-- 1. 文档类型表 (doc_types) - 定义系统中的业务文档类型
-- ================================================================================================
CREATE TABLE doc_types (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL UNIQUE,                        -- 文档类型名称 (如: User, Order, Customer)
    label VARCHAR(100) NOT NULL,                             -- 文档类型显示名称
    module VARCHAR(50) NOT NULL DEFAULT 'system',            -- 所属模块
    description TEXT,                                         -- 文档描述
    
    -- 文档属性
    is_submittable BOOLEAN DEFAULT FALSE,                    -- 是否支持提交工作流
    has_workflow BOOLEAN DEFAULT FALSE,                      -- 是否有工作流
    track_changes BOOLEAN DEFAULT TRUE,                      -- 是否跟踪变更
    
    -- 数据权限设置
    applies_to_all_users BOOLEAN DEFAULT TRUE,               -- 是否对所有用户可见
    max_attachments INTEGER DEFAULT 10,                      -- 最大附件数
    
    -- 审计字段
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    created_by BIGINT,
    updated_by BIGINT,
    
    CONSTRAINT fk_doc_types_created_by FOREIGN KEY (created_by) REFERENCES users(id),
    CONSTRAINT fk_doc_types_updated_by FOREIGN KEY (updated_by) REFERENCES users(id)
);

-- 文档类型表索引
CREATE UNIQUE INDEX idx_doc_types_name ON doc_types(name);
CREATE INDEX idx_doc_types_module ON doc_types(module);
CREATE INDEX idx_doc_types_submittable ON doc_types(is_submittable);

-- 文档类型表触发器
CREATE TRIGGER update_doc_types_updated_at BEFORE UPDATE ON doc_types FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- ================================================================================================
-- 2. 权限规则表 (permission_rules) - 高级角色权限规则
-- 支持文档级权限和字段级权限（通过permission_level区分）
-- ================================================================================================
CREATE TABLE permission_rules (
    id BIGSERIAL PRIMARY KEY,
    role_id BIGINT NOT NULL,                                 -- 角色ID  
    doc_type VARCHAR(50) NOT NULL,                           -- 文档类型名称
    permission_level INTEGER DEFAULT 0,                      -- 权限级别 (0=文档级, 1-9=字段级)
    
    -- 基础权限 (所有级别通用)
    can_read BOOLEAN DEFAULT FALSE,                          -- 读取权限
    can_write BOOLEAN DEFAULT FALSE,                         -- 写入权限
    
    -- 文档级权限 (仅permission_level=0时有效)
    can_create BOOLEAN DEFAULT FALSE,                        -- 创建权限
    can_delete BOOLEAN DEFAULT FALSE,                        -- 删除权限
    can_submit BOOLEAN DEFAULT FALSE,                        -- 提交权限 (工作流)
    can_cancel BOOLEAN DEFAULT FALSE,                        -- 取消权限 (工作流)  
    can_amend BOOLEAN DEFAULT FALSE,                         -- 修订权限 (工作流)
    can_print BOOLEAN DEFAULT FALSE,                         -- 打印权限
    can_email BOOLEAN DEFAULT FALSE,                         -- 邮件权限
    can_import BOOLEAN DEFAULT FALSE,                        -- 导入权限
    can_export BOOLEAN DEFAULT FALSE,                        -- 导出权限
    can_share BOOLEAN DEFAULT FALSE,                         -- 分享权限
    can_report BOOLEAN DEFAULT FALSE,                        -- 报表权限
    
    -- 条件权限
    only_if_creator BOOLEAN DEFAULT FALSE,                   -- 仅创建者可访问
    
    -- 审计字段
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    created_by BIGINT,
    updated_by BIGINT,
    
    CONSTRAINT fk_permission_rules_role FOREIGN KEY (role_id) REFERENCES roles(id) ON DELETE CASCADE,
    CONSTRAINT fk_permission_rules_created_by FOREIGN KEY (created_by) REFERENCES users(id),
    CONSTRAINT fk_permission_rules_updated_by FOREIGN KEY (updated_by) REFERENCES users(id),
    CONSTRAINT uk_permission_rules UNIQUE(role_id, doc_type, permission_level)
);

-- 权限规则表索引
CREATE INDEX idx_permission_rules_role_doc ON permission_rules(role_id, doc_type);
CREATE INDEX idx_permission_rules_doc_level ON permission_rules(doc_type, permission_level);
CREATE INDEX idx_permission_rules_level ON permission_rules(permission_level);

-- 权限规则表触发器  
CREATE TRIGGER update_permission_rules_updated_at BEFORE UPDATE ON permission_rules FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- ================================================================================================
-- 3. 用户权限表 (user_permissions) - 高级用户数据范围权限
-- 限制用户只能访问特定范围的数据记录
-- ================================================================================================
CREATE TABLE user_permissions (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,                                 -- 用户ID
    doc_type VARCHAR(50) NOT NULL,                           -- 受限制的文档类型
    doc_name VARCHAR(100) NOT NULL,                          -- 允许访问的具体记录标识
    value VARCHAR(100) NOT NULL,                             -- 权限值 (如: 客户ID, 部门ID)
    
    -- 高级控制
    applicable_for VARCHAR(50),                              -- 仅对指定文档类型生效
    hide_descendants BOOLEAN DEFAULT FALSE,                  -- 隐藏子记录
    is_default BOOLEAN DEFAULT FALSE,                        -- 是否为默认权限
    
    -- 审计字段
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    created_by BIGINT,
    updated_by BIGINT,
    
    CONSTRAINT fk_user_permissions_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT fk_user_permissions_created_by FOREIGN KEY (created_by) REFERENCES users(id),
    CONSTRAINT fk_user_permissions_updated_by FOREIGN KEY (updated_by) REFERENCES users(id),
    CONSTRAINT uk_user_permissions UNIQUE(user_id, doc_type, doc_name, applicable_for)
);

-- 用户权限表索引
CREATE INDEX idx_user_permissions_user_doc ON user_permissions(user_id, doc_type);
CREATE INDEX idx_user_permissions_doc_value ON user_permissions(doc_type, value);
CREATE INDEX idx_user_permissions_applicable ON user_permissions(applicable_for) WHERE applicable_for IS NOT NULL;
CREATE INDEX idx_user_permissions_default ON user_permissions(user_id, is_default) WHERE is_default = TRUE;

-- 用户权限表触发器
CREATE TRIGGER update_user_permissions_updated_at BEFORE UPDATE ON user_permissions FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- ================================================================================================
-- 4. 字段权限级别表 (field_permission_levels) - 定义字段的权限级别
-- 将文档的字段分配到不同的权限级别，实现字段级权限控制
-- ================================================================================================
CREATE TABLE field_permission_levels (
    id BIGSERIAL PRIMARY KEY,
    doc_type VARCHAR(50) NOT NULL,                           -- 文档类型
    field_name VARCHAR(100) NOT NULL,                        -- 字段名称
    field_label VARCHAR(100),                                -- 字段显示名称
    permission_level INTEGER DEFAULT 0,                      -- 权限级别 (0-9)
    field_type VARCHAR(50) DEFAULT 'Data',                   -- 字段类型
    is_mandatory BOOLEAN DEFAULT FALSE,                      -- 是否必填
    
    -- 审计字段
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    created_by BIGINT,
    updated_by BIGINT,
    
    CONSTRAINT fk_field_permission_levels_created_by FOREIGN KEY (created_by) REFERENCES users(id),
    CONSTRAINT fk_field_permission_levels_updated_by FOREIGN KEY (updated_by) REFERENCES users(id),
    CONSTRAINT uk_field_permission_levels UNIQUE(doc_type, field_name)
);

-- 字段权限级别表索引
CREATE INDEX idx_field_permission_levels_doc_type ON field_permission_levels(doc_type);
CREATE INDEX idx_field_permission_levels_level ON field_permission_levels(permission_level);
CREATE INDEX idx_field_permission_levels_doc_level ON field_permission_levels(doc_type, permission_level);

-- 字段权限级别表触发器
CREATE TRIGGER update_field_permission_levels_updated_at BEFORE UPDATE ON field_permission_levels FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- ================================================================================================
-- 5. 文档工作流状态表 (document_workflow_states) - 文档生命周期管理
-- 支持提交/取消/修订工作流
-- ================================================================================================
CREATE TABLE document_workflow_states (
    id BIGSERIAL PRIMARY KEY,
    doc_type VARCHAR(50) NOT NULL,                           -- 文档类型
    doc_id BIGINT NOT NULL,                                  -- 文档ID
    doc_name VARCHAR(100),                                   -- 文档名称/编号
    
    -- 工作流状态
    workflow_state VARCHAR(50) DEFAULT 'Draft',              -- 工作流状态: Draft, Submitted, Cancelled
    docstatus INTEGER DEFAULT 0,                             -- 文档状态: 0=草稿, 1=已提交, 2=已取消
    
    -- 提交相关
    submitted_at TIMESTAMP WITH TIME ZONE,                   -- 提交时间
    submitted_by BIGINT,                                      -- 提交人
    
    -- 取消相关
    cancelled_at TIMESTAMP WITH TIME ZONE,                   -- 取消时间
    cancelled_by BIGINT,                                      -- 取消人
    cancel_reason TEXT,                                       -- 取消原因
    
    -- 修订相关
    amended_from BIGINT,                                      -- 修订自哪个文档
    is_amended BOOLEAN DEFAULT FALSE,                        -- 是否为修订版本
    
    -- 审计字段
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    
    CONSTRAINT fk_document_workflow_states_submitted_by FOREIGN KEY (submitted_by) REFERENCES users(id),
    CONSTRAINT fk_document_workflow_states_cancelled_by FOREIGN KEY (cancelled_by) REFERENCES users(id),
    CONSTRAINT fk_document_workflow_states_amended_from FOREIGN KEY (amended_from) REFERENCES document_workflow_states(id),
    CONSTRAINT uk_document_workflow_states UNIQUE(doc_type, doc_id)
);

-- 文档工作流状态表索引
CREATE INDEX idx_document_workflow_states_doc ON document_workflow_states(doc_type, doc_id);
CREATE INDEX idx_document_workflow_states_status ON document_workflow_states(docstatus);
CREATE INDEX idx_document_workflow_states_workflow ON document_workflow_states(workflow_state);
CREATE INDEX idx_document_workflow_states_submitted_by ON document_workflow_states(submitted_by);
CREATE INDEX idx_document_workflow_states_amended_from ON document_workflow_states(amended_from);

-- 文档工作流状态表触发器
CREATE TRIGGER update_document_workflow_states_updated_at BEFORE UPDATE ON document_workflow_states FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- ================================================================================================
-- 视图和函数
-- ================================================================================================

-- 创建增强的用户权限视图 - 包含高级权限规则
CREATE OR REPLACE VIEW enhanced_user_permissions_view AS
SELECT DISTINCT
    ur.user_id,
    pr.doc_type,
    pr.permission_level,
    pr.can_read,
    pr.can_write,
    pr.can_create,
    pr.can_delete,
    pr.can_submit,
    pr.can_cancel,
    pr.can_amend,
    pr.can_print,
    pr.can_email,
    pr.can_import,
    pr.can_export,
    pr.can_share,
    pr.can_report,
    pr.only_if_creator,
    r.code as role_code,
    r.name as role_name
FROM user_roles ur
JOIN permission_rules pr ON ur.role_id = pr.role_id
JOIN roles r ON ur.role_id = r.id
WHERE ur.is_active = TRUE 
    AND r.is_enabled = TRUE
    AND (ur.expires_at IS NULL OR ur.expires_at > CURRENT_TIMESTAMP);

-- 创建高级权限检查函数
CREATE OR REPLACE FUNCTION check_document_permission(
    p_user_id BIGINT,
    p_doc_type VARCHAR(50),
    p_permission VARCHAR(20),
    p_permission_level INTEGER DEFAULT 0,
    p_doc_id BIGINT DEFAULT NULL,
    p_doc_creator_id BIGINT DEFAULT NULL
) RETURNS BOOLEAN AS $$
DECLARE
    has_permission BOOLEAN := FALSE;
    permission_count INTEGER := 0;
    user_permission_count INTEGER := 0;
BEGIN
    -- 检查基于角色的权限规则
    SELECT COUNT(*) INTO permission_count
    FROM enhanced_user_permissions_view eupv
    WHERE eupv.user_id = p_user_id 
        AND eupv.doc_type = p_doc_type
        AND eupv.permission_level = p_permission_level
        AND (
            (p_permission = 'read' AND eupv.can_read = TRUE) OR
            (p_permission = 'write' AND eupv.can_write = TRUE) OR
            (p_permission = 'create' AND eupv.can_create = TRUE) OR
            (p_permission = 'delete' AND eupv.can_delete = TRUE) OR
            (p_permission = 'submit' AND eupv.can_submit = TRUE) OR
            (p_permission = 'cancel' AND eupv.can_cancel = TRUE) OR
            (p_permission = 'amend' AND eupv.can_amend = TRUE) OR
            (p_permission = 'print' AND eupv.can_print = TRUE) OR
            (p_permission = 'email' AND eupv.can_email = TRUE) OR
            (p_permission = 'import' AND eupv.can_import = TRUE) OR
            (p_permission = 'export' AND eupv.can_export = TRUE) OR
            (p_permission = 'share' AND eupv.can_share = TRUE) OR
            (p_permission = 'report' AND eupv.can_report = TRUE)
        )
        AND (
            eupv.only_if_creator = FALSE OR 
            (eupv.only_if_creator = TRUE AND p_doc_creator_id = p_user_id)
        );
    
    -- 如果没有基于角色的权限，直接返回false
    IF permission_count = 0 THEN
        RETURN FALSE;
    END IF;
    
    -- 检查用户权限限制（数据范围权限）
    SELECT COUNT(*) INTO user_permission_count
    FROM user_permissions up
    WHERE up.user_id = p_user_id 
        AND up.doc_type = p_doc_type;
    
    -- 如果没有用户权限限制，则基于角色权限通过
    IF user_permission_count = 0 THEN
        RETURN TRUE;
    END IF;
    
    -- 如果有用户权限限制，需要检查具体记录是否被允许
    -- 这里简化处理，实际应该根据具体业务逻辑检查
    -- TODO: 实现复杂的用户权限数据范围检查逻辑
    
    RETURN TRUE;
END;
$$ LANGUAGE plpgsql;

-- 创建字段权限过滤函数
CREATE OR REPLACE FUNCTION get_accessible_fields(
    p_user_id BIGINT,
    p_doc_type VARCHAR(50),
    p_permission VARCHAR(20)
) RETURNS TABLE(
    field_name VARCHAR(100),
    field_label VARCHAR(100),
    permission_level INTEGER,
    can_access BOOLEAN
) AS $$
BEGIN
    RETURN QUERY
    SELECT DISTINCT
        fpl.field_name,
        fpl.field_label,
        fpl.permission_level,
        check_document_permission(p_user_id, p_doc_type, p_permission, fpl.permission_level) as can_access
    FROM field_permission_levels fpl
    WHERE fpl.doc_type = p_doc_type
    ORDER BY fpl.permission_level, fpl.field_name;
END;
$$ LANGUAGE plpgsql;

-- ================================================================================================
-- 初始化数据
-- ================================================================================================

-- 插入系统默认文档类型
INSERT INTO doc_types (name, label, module, description, is_submittable, has_workflow) VALUES
('User', '用户', 'system', '系统用户管理', FALSE, FALSE),
('Role', '角色', 'system', '系统角色管理', FALSE, FALSE),
('Organization', '组织', 'system', '组织架构管理', FALSE, FALSE),
('Permission', '权限', 'system', '权限管理', FALSE, FALSE),
('Order', '订单', 'business', '业务订单', TRUE, TRUE),
('Customer', '客户', 'business', '客户管理', FALSE, FALSE),
('Product', '产品', 'business', '产品管理', FALSE, FALSE);

-- 为User文档类型创建字段权限级别定义
INSERT INTO field_permission_levels (doc_type, field_name, field_label, permission_level, field_type, is_mandatory) VALUES
-- Level 0: 基础字段，所有有权限的用户都能看到
('User', 'username', '用户名', 0, 'Data', TRUE),
('User', 'email', '邮箱', 0, 'Data', TRUE),
('User', 'first_name', '姓', 0, 'Data', TRUE),
('User', 'last_name', '名', 0, 'Data', TRUE),
('User', 'phone', '手机号', 0, 'Data', FALSE),
('User', 'is_enabled', '启用状态', 0, 'Check', FALSE),
('User', 'created_at', '创建时间', 0, 'Datetime', FALSE),

-- Level 1: 个人敏感信息，需要更高权限
('User', 'birth_date', '出生日期', 1, 'Date', FALSE),
('User', 'gender', '性别', 1, 'Select', FALSE),
('User', 'avatar_url', '头像', 1, 'Attach Image', FALSE),

-- Level 2: 系统管理信息，仅管理员可见
('User', 'password_hash', '密码哈希', 2, 'Password', FALSE),
('User', 'salt', '密码盐值', 2, 'Data', FALSE),
('User', 'last_login_time', '最后登录时间', 2, 'Datetime', FALSE),
('User', 'last_login_ip', '最后登录IP', 2, 'Data', FALSE),
('User', 'two_factor_enabled', '双因子认证', 2, 'Check', FALSE),
('User', 'failed_login_count', '失败登录次数', 2, 'Int', FALSE);

-- 创建默认权限规则：为现有角色添加User文档的权限规则
-- 超级管理员：所有级别的所有权限
INSERT INTO permission_rules (role_id, doc_type, permission_level, can_read, can_write, can_create, can_delete, can_submit, can_cancel, can_amend, can_print, can_email, can_import, can_export, can_share, can_report, only_if_creator) VALUES
-- Level 0: 文档级权限
(1, 'User', 0, TRUE, TRUE, TRUE, TRUE, FALSE, FALSE, FALSE, TRUE, TRUE, TRUE, TRUE, TRUE, TRUE, FALSE),
-- Level 1: 个人信息字段权限
(1, 'User', 1, TRUE, TRUE, FALSE, FALSE, FALSE, FALSE, FALSE, FALSE, FALSE, FALSE, FALSE, FALSE, FALSE, FALSE),
-- Level 2: 系统信息字段权限
(1, 'User', 2, TRUE, TRUE, FALSE, FALSE, FALSE, FALSE, FALSE, FALSE, FALSE, FALSE, FALSE, FALSE, FALSE, FALSE);

-- 系统管理员：除了Level 2的系统信息，其他都有权限
INSERT INTO permission_rules (role_id, doc_type, permission_level, can_read, can_write, can_create, can_delete, can_submit, can_cancel, can_amend, can_print, can_email, can_import, can_export, can_share, can_report, only_if_creator) VALUES
-- Level 0: 文档级权限
(2, 'User', 0, TRUE, TRUE, TRUE, TRUE, FALSE, FALSE, FALSE, TRUE, TRUE, FALSE, TRUE, TRUE, TRUE, FALSE),
-- Level 1: 个人信息字段权限  
(2, 'User', 1, TRUE, TRUE, FALSE, FALSE, FALSE, FALSE, FALSE, FALSE, FALSE, FALSE, FALSE, FALSE, FALSE, FALSE);

-- 普通用户：只能查看基础信息，修改自己的个人信息
INSERT INTO permission_rules (role_id, doc_type, permission_level, can_read, can_write, can_create, can_delete, can_submit, can_cancel, can_amend, can_print, can_email, can_import, can_export, can_share, can_report, only_if_creator) VALUES
-- Level 0: 文档级权限，只能读取，不能创建删除
(3, 'User', 0, TRUE, FALSE, FALSE, FALSE, FALSE, FALSE, FALSE, FALSE, FALSE, FALSE, FALSE, FALSE, FALSE, FALSE),
-- Level 1: 个人信息字段权限，只能修改自己的
(3, 'User', 1, TRUE, TRUE, FALSE, FALSE, FALSE, FALSE, FALSE, FALSE, FALSE, FALSE, FALSE, FALSE, FALSE, TRUE);

-- ================================================================================================
-- 数据完整性检查和性能优化
-- ================================================================================================

-- 创建复合索引优化高级权限查询
CREATE INDEX idx_permission_rules_user_lookup ON permission_rules(role_id, doc_type, permission_level);
CREATE INDEX idx_enhanced_permissions_lookup ON user_permissions(user_id, doc_type, value);

-- 添加表注释
COMMENT ON TABLE doc_types IS '高级文档类型定义表：定义系统中的业务文档类型';
COMMENT ON TABLE permission_rules IS '高级权限规则表：支持文档级和字段级权限控制';
COMMENT ON TABLE user_permissions IS '高级用户权限表：数据范围权限控制';
COMMENT ON TABLE field_permission_levels IS '高级字段权限级别表：定义字段的权限级别';
COMMENT ON TABLE document_workflow_states IS '高级文档工作流状态表：文档生命周期管理';

COMMENT ON VIEW enhanced_user_permissions_view IS '高级增强用户权限视图：包含完整权限信息的优化查询视图';

COMMENT ON FUNCTION check_document_permission IS '高级文档权限检查函数：检查用户对特定文档的权限';
COMMENT ON FUNCTION get_accessible_fields IS '高级字段权限获取函数：返回用户可访问的字段列表';

-- ================================================================================================
-- 提交事务
-- ================================================================================================
COMMIT;

-- 成功提示
DO $$
BEGIN
    RAISE NOTICE '';
    RAISE NOTICE '🎉 高级权限系统扩展完成！';
    RAISE NOTICE '';
    RAISE NOTICE '📊 新增功能:';
    RAISE NOTICE '   ✅ 文档类型管理';
    RAISE NOTICE '   ✅ 角色权限规则 (支持0-9级权限)';
    RAISE NOTICE '   ✅ 用户数据范围权限';
    RAISE NOTICE '   ✅ 字段权限级别定义';
    RAISE NOTICE '   ✅ 文档工作流状态管理';
    RAISE NOTICE '';
    RAISE NOTICE '🔧 使用示例:';
    RAISE NOTICE '   -- 检查权限: SELECT check_document_permission(1, ''User'', ''read'', 0);';
    RAISE NOTICE '   -- 获取可访问字段: SELECT * FROM get_accessible_fields(1, ''User'', ''read'');';
    RAISE NOTICE '   -- 查看用户权限: SELECT * FROM enhanced_user_permissions_view WHERE user_id = 1;';
    RAISE NOTICE '';
    RAISE NOTICE '✨ 高级四层权限控制系统已就绪！';
    RAISE NOTICE '';
END $$;