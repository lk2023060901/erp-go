-- ================================================================
-- 操作日志表 (operation_logs)
-- ================================================================
-- 记录用户操作审计日志，支持按时间分片

DROP TABLE IF EXISTS operation_logs CASCADE;

CREATE TABLE operation_logs (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT COMMENT '操作用户ID',
    username VARCHAR(50) COMMENT '用户名快照',
    operation_type VARCHAR(50) NOT NULL COMMENT '操作类型',
    resource_type VARCHAR(50) COMMENT '资源类型',
    resource_id BIGINT COMMENT '资源ID',
    operation_desc TEXT COMMENT '操作描述',
    request_method VARCHAR(10) COMMENT 'HTTP方法',
    request_url VARCHAR(500) COMMENT '请求URL',
    request_params TEXT COMMENT '请求参数',
    request_body TEXT COMMENT '请求体',
    response_status INTEGER COMMENT '响应状态码',
    response_body TEXT COMMENT '响应内容',
    ip_address INET COMMENT 'IP地址',
    user_agent TEXT COMMENT '用户代理',
    execution_time INTEGER COMMENT '执行时间(ms)',
    is_success BOOLEAN DEFAULT true COMMENT '是否成功',
    error_message TEXT COMMENT '错误信息',
    session_id VARCHAR(128) COMMENT '会话ID',
    trace_id VARCHAR(64) COMMENT '追踪ID',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间'
);

-- 添加注释
COMMENT ON TABLE operation_logs IS '操作日志表 - 记录用户操作审计日志';

-- 创建索引
CREATE INDEX idx_operation_logs_user ON operation_logs(user_id, created_at DESC);
CREATE INDEX idx_operation_logs_type ON operation_logs(operation_type, created_at DESC);
CREATE INDEX idx_operation_logs_resource ON operation_logs(resource_type, resource_id);
CREATE INDEX idx_operation_logs_ip ON operation_logs(ip_address, created_at DESC);
CREATE INDEX idx_operation_logs_success ON operation_logs(is_success, created_at DESC);
CREATE INDEX idx_operation_logs_created_at ON operation_logs(created_at DESC);
CREATE INDEX idx_operation_logs_session ON operation_logs(session_id);
CREATE INDEX idx_operation_logs_trace ON operation_logs(trace_id);

-- 创建按月分区的操作日志表模板
-- 实际使用时需要根据业务量创建对应的分区表