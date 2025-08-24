-- ================================================================
-- 组织表 (organizations)
-- ================================================================
-- 支持无限层级的组织架构，用于部门管理和数据隔离

DROP TABLE IF EXISTS organizations CASCADE;

CREATE TABLE organizations (
    id BIGSERIAL PRIMARY KEY,
    parent_id BIGINT REFERENCES organizations(id) COMMENT '父组织ID',
    name VARCHAR(100) NOT NULL COMMENT '组织名称',
    code VARCHAR(50) NOT NULL UNIQUE COMMENT '组织编码',
    type VARCHAR(20) DEFAULT 'department' COMMENT '组织类型',
    description TEXT COMMENT '组织描述',
    level INTEGER DEFAULT 1 COMMENT '组织层级',
    path VARCHAR(500) COMMENT '组织路径',
    leader_id BIGINT COMMENT '负责人ID',
    phone VARCHAR(20) COMMENT '联系电话',
    email VARCHAR(100) COMMENT '联系邮箱',
    address TEXT COMMENT '办公地址',
    data_isolation BOOLEAN DEFAULT false COMMENT '数据隔离开关',
    sort_order INTEGER DEFAULT 0 COMMENT '排序权重',
    is_enabled BOOLEAN DEFAULT true COMMENT '组织启用状态',
    created_by BIGINT COMMENT '创建人ID',
    updated_by BIGINT COMMENT '更新人ID',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '更新时间',
    deleted_at TIMESTAMP COMMENT '删除时间'
);

-- 添加注释
COMMENT ON TABLE organizations IS '组织表 - 支持无限层级的组织架构';

-- 创建索引
CREATE UNIQUE INDEX idx_organizations_code ON organizations(code) WHERE deleted_at IS NULL;
CREATE INDEX idx_organizations_parent ON organizations(parent_id);
CREATE INDEX idx_organizations_name ON organizations(name);
CREATE INDEX idx_organizations_path ON organizations USING GIN (path gin_trgm_ops);
CREATE INDEX idx_organizations_level ON organizations(level);
CREATE INDEX idx_organizations_leader ON organizations(leader_id);
CREATE INDEX idx_organizations_enabled ON organizations(is_enabled);

-- 创建组织路径更新触发器
CREATE OR REPLACE FUNCTION update_organization_path()
RETURNS TRIGGER AS $$
BEGIN
    -- 更新路径
    IF NEW.parent_id IS NULL THEN
        NEW.path = '/' || NEW.id::text || '/';
        NEW.level = 1;
    ELSE
        SELECT path || NEW.id::text || '/', level + 1
        INTO NEW.path, NEW.level
        FROM organizations 
        WHERE id = NEW.parent_id;
    END IF;
    
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER trigger_organization_path
    BEFORE INSERT OR UPDATE ON organizations
    FOR EACH ROW
    EXECUTE FUNCTION update_organization_path();