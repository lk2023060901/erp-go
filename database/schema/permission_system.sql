-- ================================================================
-- 高级权限系统表结构定义
-- 支持多层级权限控制：文档级、字段级、用户级、角色级
-- ================================================================

-- 文档类型表 (DocType) 
CREATE TABLE IF NOT EXISTS doc_types (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name VARCHAR(140) UNIQUE NOT NULL,           -- 文档类型名称
    label VARCHAR(140) NOT NULL,                 -- 显示名称
    module VARCHAR(140) NOT NULL,                -- 所属模块
    description TEXT,                            -- 描述
    is_submittable BOOLEAN DEFAULT false,        -- 是否支持提交工作流
    is_child_table BOOLEAN DEFAULT false,        -- 是否是子表
    has_workflow BOOLEAN DEFAULT false,          -- 是否有工作流
    track_changes BOOLEAN DEFAULT false,         -- 是否跟踪变更
    applies_to_all_users BOOLEAN DEFAULT false,  -- 是否对所有用户可见
    max_attachments INTEGER DEFAULT 0,           -- 最大附件数
    permissions TEXT,                            -- 权限设置JSON
    naming_rule VARCHAR(140),                    -- 命名规则
    title_field VARCHAR(140),                    -- 标题字段
    search_fields TEXT,                          -- 搜索字段JSON数组
    sort_field VARCHAR(140),                     -- 排序字段
    sort_order VARCHAR(10) DEFAULT 'ASC',        -- 排序方向
    version INTEGER DEFAULT 1,                   -- 版本号
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    created_by INTEGER,
    updated_by INTEGER
);

-- 权限规则表 (Permission Rules)
CREATE TABLE IF NOT EXISTS permission_rules (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    role INTEGER NOT NULL,                       -- 角色ID  
    document_type VARCHAR(140) NOT NULL,         -- 文档类型
    permission_level INTEGER DEFAULT 0,          -- 权限级别 (0=文档级, 1-10=字段级)
    
    -- 权限设置
    read BOOLEAN DEFAULT false,                  -- 读权限
    write BOOLEAN DEFAULT false,                 -- 写权限
    [create] BOOLEAN DEFAULT false,              -- 创建权限
    [delete] BOOLEAN DEFAULT false,              -- 删除权限
    submit BOOLEAN DEFAULT false,                -- 提交权限
    cancel BOOLEAN DEFAULT false,                -- 取消权限
    amend BOOLEAN DEFAULT false,                 -- 修正权限
    print BOOLEAN DEFAULT false,                 -- 打印权限
    email BOOLEAN DEFAULT false,                 -- 邮件权限
    import BOOLEAN DEFAULT false,                -- 导入权限
    export BOOLEAN DEFAULT false,                -- 导出权限
    share BOOLEAN DEFAULT false,                 -- 分享权限
    report BOOLEAN DEFAULT false,                -- 报表权限
    set_user_permissions BOOLEAN DEFAULT false,  -- 设置用户权限
    
    -- 条件设置
    if_owner BOOLEAN DEFAULT false,              -- 仅限创建者
    match TEXT,                                  -- 匹配条件
    select_condition TEXT,                       -- 选择条件
    delete_condition TEXT,                       -- 删除条件
    amend_condition TEXT,                        -- 修正条件
    
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    created_by INTEGER,
    updated_by INTEGER,
    
    UNIQUE(role, document_type, permission_level)
);

-- 用户权限表 (User Permissions)
CREATE TABLE IF NOT EXISTS user_permissions (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,                    -- 用户ID
    doc_type VARCHAR(140) NOT NULL,              -- 文档类型
    doc_name VARCHAR(140),                       -- 文档名称(可选，用于特定记录权限)
    permission_value VARCHAR(500),               -- 权限值
    applicable_for VARCHAR(140),                 -- 适用于哪个DocType
    hide_descendants BOOLEAN DEFAULT false,      -- 隐藏子项
    is_default BOOLEAN DEFAULT false,            -- 是否默认权限
    
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    created_by INTEGER,
    updated_by INTEGER,
    
    UNIQUE(user_id, doc_type, doc_name, applicable_for)
);

-- 字段权限级别表 (Field Permission Levels)
CREATE TABLE IF NOT EXISTS field_permission_levels (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    document_type VARCHAR(140) NOT NULL,         -- 文档类型
    field_name VARCHAR(140) NOT NULL,            -- 字段名称
    field_label VARCHAR(140),                    -- 字段显示名称
    permission_level INTEGER NOT NULL,           -- 权限级别
    field_type VARCHAR(50),                      -- 字段类型
    is_mandatory BOOLEAN DEFAULT false,          -- 是否必填
    
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    created_by INTEGER,
    updated_by INTEGER,
    
    UNIQUE(document_type, field_name)
);

-- 文档工作流状态表 (Document Workflow States)
CREATE TABLE IF NOT EXISTS document_workflow_states (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    document_type VARCHAR(140) NOT NULL,         -- 文档类型
    document_name VARCHAR(140),                  -- 文档名称
    workflow_state VARCHAR(140) NOT NULL,       -- 工作流状态
    workflow_action VARCHAR(140),               -- 工作流动作
    
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    created_by INTEGER
);

-- ================================================================
-- 权限系统索引创建
-- ================================================================

-- DocType索引
CREATE INDEX IF NOT EXISTS idx_doc_types_name ON doc_types(name);
CREATE INDEX IF NOT EXISTS idx_doc_types_module ON doc_types(module);

-- Permission Rules索引  
CREATE INDEX IF NOT EXISTS idx_permission_rules_role_doctype ON permission_rules(role, document_type);
CREATE INDEX IF NOT EXISTS idx_permission_rules_doctype ON permission_rules(document_type);
CREATE INDEX IF NOT EXISTS idx_permission_rules_level ON permission_rules(permission_level);

-- User Permissions索引
CREATE INDEX IF NOT EXISTS idx_user_permissions_user_doctype ON user_permissions(user_id, doc_type);
CREATE INDEX IF NOT EXISTS idx_user_permissions_doctype ON user_permissions(doc_type);

-- Field Permission Levels索引
CREATE INDEX IF NOT EXISTS idx_field_permission_levels_doctype ON field_permission_levels(document_type);
CREATE INDEX IF NOT EXISTS idx_field_permission_levels_field ON field_permission_levels(document_type, field_name);

-- Document Workflow States索引
CREATE INDEX IF NOT EXISTS idx_doc_workflow ON document_workflow_states(document_type, document_name);