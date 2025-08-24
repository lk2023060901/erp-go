-- ================================================================
-- 系统配置表 (system_configs)
-- ================================================================
-- 存储系统运行时配置参数

DROP TABLE IF EXISTS system_configs CASCADE;

CREATE TABLE system_configs (
    id BIGSERIAL PRIMARY KEY,
    config_key VARCHAR(100) NOT NULL UNIQUE COMMENT '配置键',
    config_value TEXT COMMENT '配置值',
    config_type VARCHAR(20) DEFAULT 'string' COMMENT '配置类型',
    description TEXT COMMENT '配置描述',
    is_system BOOLEAN DEFAULT false COMMENT '是否系统配置',
    is_encrypted BOOLEAN DEFAULT false COMMENT '是否加密存储',
    group_name VARCHAR(50) COMMENT '配置分组',
    sort_order INTEGER DEFAULT 0 COMMENT '排序权重',
    is_enabled BOOLEAN DEFAULT true COMMENT '配置启用状态',
    created_by BIGINT COMMENT '创建人ID',
    updated_by BIGINT COMMENT '更新人ID',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '更新时间'
);

-- 添加注释
COMMENT ON TABLE system_configs IS '系统配置表 - 存储系统运行时配置参数';

-- 创建索引
CREATE UNIQUE INDEX idx_system_configs_key ON system_configs(config_key);
CREATE INDEX idx_system_configs_group ON system_configs(group_name);
CREATE INDEX idx_system_configs_type ON system_configs(config_type);
CREATE INDEX idx_system_configs_enabled ON system_configs(is_enabled);
CREATE INDEX idx_system_configs_system ON system_configs(is_system);

-- 创建更新时间触发器
CREATE OR REPLACE FUNCTION update_system_configs_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER trigger_system_configs_updated_at
    BEFORE UPDATE ON system_configs
    FOR EACH ROW
    EXECUTE FUNCTION update_system_configs_updated_at();