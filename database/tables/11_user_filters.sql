-- ================================================================
-- 用户过滤器表
-- ================================================================
-- 存储用户自定义的过滤器配置，支持多模块复用
-- 创建时间: 2024-01-20
-- 版本: 1.0

-- 删除表（如果存在）
DROP TABLE IF EXISTS user_filters CASCADE;

-- 创建用户过滤器表
CREATE TABLE user_filters (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    module_type VARCHAR(50) NOT NULL, -- 模块类型：users, roles, permissions, organizations等
    filter_name VARCHAR(100) NOT NULL, -- 过滤器名称
    filter_conditions JSONB NOT NULL DEFAULT '{}', -- 过滤条件（JSON格式）
    sort_config JSONB, -- 排序配置（JSON格式）
    is_default BOOLEAN DEFAULT FALSE, -- 是否为默认过滤器
    is_public BOOLEAN DEFAULT FALSE, -- 是否为公共过滤器（所有用户可见）
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- 创建索引
CREATE INDEX idx_user_filters_user_module ON user_filters(user_id, module_type);
CREATE INDEX idx_user_filters_module_type ON user_filters(module_type);
CREATE INDEX idx_user_filters_user_default ON user_filters(user_id, is_default) WHERE is_default = TRUE;
CREATE INDEX idx_user_filters_public ON user_filters(module_type, is_public) WHERE is_public = TRUE;
CREATE INDEX idx_user_filters_created_at ON user_filters(created_at);

-- 创建更新时间触发器
CREATE OR REPLACE FUNCTION update_user_filters_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER trigger_user_filters_updated_at
    BEFORE UPDATE ON user_filters
    FOR EACH ROW
    EXECUTE FUNCTION update_user_filters_updated_at();

-- 添加约束
ALTER TABLE user_filters ADD CONSTRAINT unique_user_module_filter_name 
    UNIQUE (user_id, module_type, filter_name);

-- 添加检查约束
ALTER TABLE user_filters ADD CONSTRAINT check_module_type 
    CHECK (module_type IN ('users', 'roles', 'permissions', 'organizations', 'operation_logs'));

-- 添加注释
COMMENT ON TABLE user_filters IS '用户过滤器配置表，支持多模块复用';
COMMENT ON COLUMN user_filters.id IS '主键ID';
COMMENT ON COLUMN user_filters.user_id IS '用户ID，关联users表';
COMMENT ON COLUMN user_filters.module_type IS '模块类型：users, roles, permissions, organizations等';
COMMENT ON COLUMN user_filters.filter_name IS '过滤器名称';
COMMENT ON COLUMN user_filters.filter_conditions IS '过滤条件，JSON格式存储';
COMMENT ON COLUMN user_filters.sort_config IS '排序配置，JSON格式存储';
COMMENT ON COLUMN user_filters.is_default IS '是否为用户在该模块的默认过滤器';
COMMENT ON COLUMN user_filters.is_public IS '是否为公共过滤器（所有用户可见）';
COMMENT ON COLUMN user_filters.created_at IS '创建时间';
COMMENT ON COLUMN user_filters.updated_at IS '更新时间';

-- 插入示例数据
INSERT INTO user_filters (user_id, module_type, filter_name, filter_conditions, sort_config, is_default) 
VALUES 
(1, 'users', '活跃用户', 
 '{"conditions": [{"field": "is_active", "operator": "equals", "value": true}], "logic": "AND"}',
 '{"field": "username", "direction": "asc"}',
 false),
(1, 'users', '管理员用户', 
 '{"conditions": [{"field": "roles", "operator": "contains", "value": "admin"}], "logic": "AND"}',
 '{"field": "created_at", "direction": "desc"}',
 false),
(1, 'users', '近期创建用户', 
 '{"conditions": [{"field": "created_at", "operator": "greater_than", "value": "2024-01-01"}], "logic": "AND"}',
 '{"field": "created_at", "direction": "desc"}',
 false);

\echo '✓ 用户过滤器表创建完成'