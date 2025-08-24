-- ================================================================
-- 用户会话表 (user_sessions)
-- ================================================================
-- 存储用户登录会话信息，支持JWT Token管理

DROP TABLE IF EXISTS user_sessions CASCADE;

CREATE TABLE user_sessions (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE COMMENT '用户ID',
    session_id VARCHAR(128) NOT NULL UNIQUE COMMENT '会话ID',
    token_hash VARCHAR(255) NOT NULL COMMENT 'Token哈希',
    refresh_token_hash VARCHAR(255) COMMENT '刷新Token哈希',
    device_type VARCHAR(20) DEFAULT 'web' COMMENT '设备类型',
    device_id VARCHAR(128) COMMENT '设备ID',
    ip_address INET COMMENT 'IP地址',
    user_agent TEXT COMMENT '用户代理',
    location VARCHAR(200) COMMENT '地理位置',
    is_active BOOLEAN DEFAULT true COMMENT '会话状态',
    last_activity_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '最后活跃时间',
    expires_at TIMESTAMP NOT NULL COMMENT '过期时间',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '更新时间'
);

-- 添加注释
COMMENT ON TABLE user_sessions IS '用户会话表 - 存储用户登录会话信息';

-- 创建索引
CREATE INDEX idx_user_sessions_user ON user_sessions(user_id, is_active);
CREATE INDEX idx_user_sessions_token ON user_sessions(token_hash);
CREATE INDEX idx_user_sessions_refresh ON user_sessions(refresh_token_hash);
CREATE INDEX idx_user_sessions_device ON user_sessions(user_id, device_type, device_id);
CREATE INDEX idx_user_sessions_expires ON user_sessions(expires_at);
CREATE INDEX idx_user_sessions_activity ON user_sessions(last_activity_at DESC);
CREATE INDEX idx_user_sessions_created_at ON user_sessions(created_at DESC);

-- 创建会话清理功能
CREATE OR REPLACE FUNCTION cleanup_expired_sessions()
RETURNS INTEGER AS $$
DECLARE
    deleted_count INTEGER;
BEGIN
    -- 删除过期会话
    DELETE FROM user_sessions 
    WHERE expires_at < CURRENT_TIMESTAMP 
    OR (is_active = false AND updated_at < CURRENT_TIMESTAMP - INTERVAL '7 days');
    
    GET DIAGNOSTICS deleted_count = ROW_COUNT;
    RETURN deleted_count;
END;
$$ LANGUAGE plpgsql;

-- 创建更新时间触发器
CREATE OR REPLACE FUNCTION update_user_sessions_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER trigger_user_sessions_updated_at
    BEFORE UPDATE ON user_sessions
    FOR EACH ROW
    EXECUTE FUNCTION update_user_sessions_updated_at();