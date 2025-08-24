-- ================================================================
-- 用户表 (users)
-- ================================================================
-- 存储用户基本信息和认证相关数据
-- 支持用户名/邮箱登录，包含完整的个人信息字段

DROP TABLE IF EXISTS users CASCADE;

CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    username VARCHAR(50) NOT NULL UNIQUE COMMENT '用户名',
    email VARCHAR(100) NOT NULL UNIQUE COMMENT '邮箱地址',
    first_name VARCHAR(50) NOT NULL COMMENT '姓',
    last_name VARCHAR(50) NOT NULL COMMENT '名',
    phone VARCHAR(20) COMMENT '手机号',
    gender CHAR(1) CHECK (gender IN ('M', 'F', 'O')) COMMENT '性别: M=男, F=女, O=其他',
    birth_date DATE COMMENT '出生日期',
    password_hash VARCHAR(255) NOT NULL COMMENT '密码哈希',
    salt VARCHAR(100) NOT NULL COMMENT '密码盐值',
    avatar_url VARCHAR(500) COMMENT '头像URL',
    status SMALLINT DEFAULT 1 CHECK (status IN (0, 1, -1)) COMMENT '状态: 1=正常, 0=禁用, -1=删除',
    is_enabled BOOLEAN DEFAULT true COMMENT '账户启用状态',
    phone_verified BOOLEAN DEFAULT false COMMENT '手机验证状态',
    two_factor_enabled BOOLEAN DEFAULT false COMMENT '双重验证开关',
    two_factor_secret VARCHAR(32) COMMENT '2FA密钥',
    login_count INTEGER DEFAULT 0 COMMENT '登录次数',
    first_login_time TIMESTAMP COMMENT '首次登录时间',
    last_login_time TIMESTAMP COMMENT '最后登录时间',
    last_login_ip INET COMMENT '最后登录IP',
    current_login_time TIMESTAMP COMMENT '当前登录时间',
    current_login_ip INET COMMENT '当前登录IP',
    password_changed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '密码修改时间',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '更新时间',
    deleted_at TIMESTAMP COMMENT '删除时间'
);

-- 添加注释
COMMENT ON TABLE users IS '用户表 - 存储用户基本信息和认证相关数据';

-- 创建索引
CREATE UNIQUE INDEX idx_users_username ON users(username) WHERE deleted_at IS NULL;
CREATE UNIQUE INDEX idx_users_email ON users(email) WHERE deleted_at IS NULL;
CREATE INDEX idx_users_status ON users(status);
CREATE INDEX idx_users_created_at ON users(created_at);
CREATE INDEX idx_users_phone ON users(phone) WHERE phone IS NOT NULL AND deleted_at IS NULL;

-- 创建更新时间触发器
CREATE OR REPLACE FUNCTION update_users_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER trigger_users_updated_at
    BEFORE UPDATE ON users
    FOR EACH ROW
    EXECUTE FUNCTION update_users_updated_at();