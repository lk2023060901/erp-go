-- ================================================================
-- 角色表 (roles)
-- ================================================================
-- 存储角色定义和属性，支持系统预定义角色

DROP TABLE IF EXISTS roles CASCADE;

CREATE TABLE roles (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL COMMENT '角色名称',
    code VARCHAR(50) NOT NULL UNIQUE COMMENT '角色编码',
    description TEXT COMMENT '角色描述',
    is_system_role BOOLEAN DEFAULT false COMMENT '是否为系统预定义角色',
    is_enabled BOOLEAN DEFAULT true COMMENT '角色启用状态',
    sort_order INTEGER DEFAULT 0 COMMENT '排序权重',
    created_by BIGINT COMMENT '创建人ID',
    updated_by BIGINT COMMENT '更新人ID',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '更新时间',
    deleted_at TIMESTAMP COMMENT '删除时间'
);

-- 添加注释
COMMENT ON TABLE roles IS '角色表 - 存储角色定义和属性';

-- 创建索引
CREATE UNIQUE INDEX idx_roles_code ON roles(code) WHERE deleted_at IS NULL;
CREATE INDEX idx_roles_name ON roles(name);
CREATE INDEX idx_roles_is_system ON roles(is_system_role);
CREATE INDEX idx_roles_is_enabled ON roles(is_enabled);
CREATE INDEX idx_roles_created_at ON roles(created_at);

-- 创建更新时间触发器
CREATE OR REPLACE FUNCTION update_roles_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER trigger_roles_updated_at
    BEFORE UPDATE ON roles
    FOR EACH ROW
    EXECUTE FUNCTION update_roles_updated_at();