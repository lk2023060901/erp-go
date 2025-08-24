-- ================================================================
-- 用户组织关联表 (user_organizations)
-- ================================================================
-- 用户与组织的关联关系，支持单主组织模式

DROP TABLE IF EXISTS user_organizations CASCADE;

CREATE TABLE user_organizations (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE COMMENT '用户ID',
    org_id BIGINT NOT NULL REFERENCES organizations(id) ON DELETE CASCADE COMMENT '组织ID',
    position VARCHAR(100) COMMENT '职位',
    is_primary BOOLEAN DEFAULT false COMMENT '是否主组织',
    is_leader BOOLEAN DEFAULT false COMMENT '是否负责人',
    join_date DATE DEFAULT CURRENT_DATE COMMENT '入职日期',
    leave_date DATE COMMENT '离职日期',
    is_active BOOLEAN DEFAULT true COMMENT '激活状态',
    created_by BIGINT COMMENT '创建人ID',
    updated_by BIGINT COMMENT '更新人ID',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '更新时间',
    
    CONSTRAINT uk_user_organizations UNIQUE(user_id, org_id)
);

-- 添加注释
COMMENT ON TABLE user_organizations IS '用户组织关联表 - 用户与组织的关联关系';

-- 创建索引
CREATE INDEX idx_user_organizations_user ON user_organizations(user_id, is_active);
CREATE INDEX idx_user_organizations_org ON user_organizations(org_id, is_active);
CREATE INDEX idx_user_organizations_primary ON user_organizations(user_id, is_primary) WHERE is_primary = true;
CREATE INDEX idx_user_organizations_leader ON user_organizations(org_id, is_leader) WHERE is_leader = true;
CREATE INDEX idx_user_organizations_position ON user_organizations(position);
CREATE INDEX idx_user_organizations_join_date ON user_organizations(join_date);

-- 确保用户只有一个主组织的触发器
CREATE OR REPLACE FUNCTION ensure_single_primary_org()
RETURNS TRIGGER AS $$
BEGIN
    -- 如果设置为主组织，将该用户的其他组织关系设为非主组织
    IF NEW.is_primary = true THEN
        UPDATE user_organizations 
        SET is_primary = false, updated_at = CURRENT_TIMESTAMP
        WHERE user_id = NEW.user_id 
        AND id != COALESCE(NEW.id, 0)
        AND is_primary = true;
    END IF;
    
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER trigger_single_primary_org
    BEFORE INSERT OR UPDATE ON user_organizations
    FOR EACH ROW
    EXECUTE FUNCTION ensure_single_primary_org();