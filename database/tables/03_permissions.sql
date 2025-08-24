-- ================================================================
-- 权限表 (permissions)
-- ================================================================
-- 存储权限定义，支持层级结构和三种权限类型

DROP TABLE IF EXISTS permissions CASCADE;

CREATE TABLE permissions (
    id BIGSERIAL PRIMARY KEY,
    parent_id BIGINT REFERENCES permissions(id) COMMENT '父权限ID',
    name VARCHAR(100) NOT NULL COMMENT '权限名称',
    code VARCHAR(100) NOT NULL UNIQUE COMMENT '权限编码',
    resource VARCHAR(50) NOT NULL COMMENT '资源名称',
    action VARCHAR(50) NOT NULL COMMENT '操作类型',
    module VARCHAR(50) NOT NULL COMMENT '所属模块',
    description TEXT COMMENT '权限描述',
    is_menu BOOLEAN DEFAULT false COMMENT '是否为菜单权限',
    is_button BOOLEAN DEFAULT false COMMENT '是否为按钮权限',
    is_api BOOLEAN DEFAULT false COMMENT '是否为API权限',
    menu_url VARCHAR(200) COMMENT '菜单URL',
    menu_icon VARCHAR(50) COMMENT '菜单图标',
    api_path VARCHAR(200) COMMENT 'API路径',
    api_method VARCHAR(10) COMMENT 'HTTP方法',
    level INTEGER DEFAULT 1 COMMENT '权限层级',
    sort_order INTEGER DEFAULT 0 COMMENT '排序权重',
    is_enabled BOOLEAN DEFAULT true COMMENT '权限启用状态',
    created_by BIGINT COMMENT '创建人ID',
    updated_by BIGINT COMMENT '更新人ID',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '更新时间',
    deleted_at TIMESTAMP COMMENT '删除时间'
);

-- 添加注释
COMMENT ON TABLE permissions IS '权限表 - 存储权限定义，支持层级结构';

-- 创建索引
CREATE UNIQUE INDEX idx_permissions_code ON permissions(code) WHERE deleted_at IS NULL;
CREATE INDEX idx_permissions_parent ON permissions(parent_id);
CREATE INDEX idx_permissions_resource_action ON permissions(resource, action);
CREATE INDEX idx_permissions_module ON permissions(module);
CREATE INDEX idx_permissions_menu ON permissions(is_menu) WHERE is_menu = true;
CREATE INDEX idx_permissions_button ON permissions(is_button) WHERE is_button = true;
CREATE INDEX idx_permissions_api ON permissions(is_api) WHERE is_api = true;
CREATE INDEX idx_permissions_api_path ON permissions(api_path, api_method) WHERE is_api = true;
CREATE INDEX idx_permissions_level ON permissions(level);
CREATE INDEX idx_permissions_enabled ON permissions(is_enabled);

-- 创建更新时间触发器
CREATE OR REPLACE FUNCTION update_permissions_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER trigger_permissions_updated_at
    BEFORE UPDATE ON permissions
    FOR EACH ROW
    EXECUTE FUNCTION update_permissions_updated_at();