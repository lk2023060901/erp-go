-- ================================================================
-- 用户角色关联表 (user_roles)
-- ================================================================
-- 用户与角色的多对多关系，支持临时权限

DROP TABLE IF EXISTS user_roles CASCADE;

CREATE TABLE user_roles (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE COMMENT '用户ID',
    role_id BIGINT NOT NULL REFERENCES roles(id) ON DELETE CASCADE COMMENT '角色ID',
    expires_at TIMESTAMP COMMENT '过期时间',
    is_active BOOLEAN DEFAULT true COMMENT '激活状态',
    granted_by BIGINT COMMENT '授权人ID',
    granted_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '授权时间',
    revoked_by BIGINT COMMENT '撤销人ID',
    revoked_at TIMESTAMP COMMENT '撤销时间',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '更新时间',
    
    CONSTRAINT uk_user_roles UNIQUE(user_id, role_id)
);

-- 添加注释
COMMENT ON TABLE user_roles IS '用户角色关联表 - 用户与角色的多对多关系';

-- 创建索引
CREATE INDEX idx_user_roles_user ON user_roles(user_id, is_active);
CREATE INDEX idx_user_roles_role ON user_roles(role_id, is_active);
CREATE INDEX idx_user_roles_expires ON user_roles(expires_at) WHERE expires_at IS NOT NULL;
CREATE INDEX idx_user_roles_permission_check ON user_roles(user_id, is_active, expires_at);
CREATE INDEX idx_user_roles_granted_by ON user_roles(granted_by);
CREATE INDEX idx_user_roles_created_at ON user_roles(created_at);

-- 创建更新时间触发器
CREATE OR REPLACE FUNCTION update_user_roles_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    
    -- 如果is_active从true变为false，记录撤销信息
    IF OLD.is_active = true AND NEW.is_active = false AND NEW.revoked_at IS NULL THEN
        NEW.revoked_at = CURRENT_TIMESTAMP;
    END IF;
    
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER trigger_user_roles_updated_at
    BEFORE UPDATE ON user_roles
    FOR EACH ROW
    EXECUTE FUNCTION update_user_roles_updated_at();