-- ================================================================
-- 角色权限关联表 (role_permissions)
-- ================================================================
-- 角色与权限的多对多关系，支持拒绝权限

DROP TABLE IF EXISTS role_permissions CASCADE;

CREATE TABLE role_permissions (
    id BIGSERIAL PRIMARY KEY,
    role_id BIGINT NOT NULL REFERENCES roles(id) ON DELETE CASCADE COMMENT '角色ID',
    permission_id BIGINT NOT NULL REFERENCES permissions(id) ON DELETE CASCADE COMMENT '权限ID',
    is_granted BOOLEAN DEFAULT true COMMENT '是否授权',
    granted_by BIGINT COMMENT '授权人ID',
    granted_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '授权时间',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '更新时间',
    
    CONSTRAINT uk_role_permissions UNIQUE(role_id, permission_id)
);

-- 添加注释
COMMENT ON TABLE role_permissions IS '角色权限关联表 - 角色与权限的多对多关系';

-- 创建索引
CREATE INDEX idx_role_permissions_role ON role_permissions(role_id, is_granted);
CREATE INDEX idx_role_permissions_permission ON role_permissions(permission_id, is_granted);
CREATE INDEX idx_role_permissions_check ON role_permissions(role_id, is_granted, permission_id);
CREATE INDEX idx_role_permissions_granted_by ON role_permissions(granted_by);
CREATE INDEX idx_role_permissions_created_at ON role_permissions(created_at);

-- 创建更新时间触发器
CREATE OR REPLACE FUNCTION update_role_permissions_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER trigger_role_permissions_updated_at
    BEFORE UPDATE ON role_permissions
    FOR EACH ROW
    EXECUTE FUNCTION update_role_permissions_updated_at();