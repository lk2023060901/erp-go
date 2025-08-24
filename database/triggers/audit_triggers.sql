-- ================================================================
-- 审计触发器脚本
-- ================================================================
-- 自动记录数据变更的审计日志

\echo '创建审计触发器...'

-- 创建通用审计函数
CREATE OR REPLACE FUNCTION audit_trigger_function()
RETURNS TRIGGER AS $$
DECLARE
    old_data JSONB;
    new_data JSONB;
    operation_desc TEXT;
    resource_type TEXT;
    resource_id BIGINT;
BEGIN
    -- 确定资源类型
    resource_type := TG_TABLE_NAME;
    
    -- 确定操作描述和资源ID
    IF TG_OP = 'DELETE' THEN
        old_data := row_to_json(OLD)::JSONB;
        resource_id := (OLD.id)::BIGINT;
        operation_desc := '删除' || resource_type;
    ELSIF TG_OP = 'INSERT' THEN
        new_data := row_to_json(NEW)::JSONB;
        resource_id := (NEW.id)::BIGINT;
        operation_desc := '创建' || resource_type;
    ELSIF TG_OP = 'UPDATE' THEN
        old_data := row_to_json(OLD)::JSONB;
        new_data := row_to_json(NEW)::JSONB;
        resource_id := (NEW.id)::BIGINT;
        operation_desc := '更新' || resource_type;
    END IF;
    
    -- 插入审计日志（如果operation_logs表存在）
    BEGIN
        INSERT INTO operation_logs (
            user_id,
            username,
            operation_type,
            resource_type,
            resource_id,
            operation_desc,
            request_params,
            request_body,
            created_at
        ) VALUES (
            COALESCE(
                CASE WHEN TG_OP = 'DELETE' THEN NULL 
                     ELSE (NEW.updated_by)::BIGINT 
                END,
                CASE WHEN TG_OP = 'DELETE' THEN NULL 
                     ELSE (NEW.created_by)::BIGINT 
                END
            ),
            NULL, -- 用户名将在应用层填充
            TG_OP,
            resource_type,
            resource_id,
            operation_desc,
            CASE WHEN old_data IS NOT NULL 
                 THEN jsonb_build_object('old_data', old_data)::TEXT
                 ELSE NULL 
            END,
            CASE WHEN new_data IS NOT NULL 
                 THEN jsonb_build_object('new_data', new_data)::TEXT
                 ELSE NULL 
            END,
            CURRENT_TIMESTAMP
        );
    EXCEPTION WHEN OTHERS THEN
        -- 审计失败不应影响业务操作，记录错误但继续
        RAISE WARNING '审计日志记录失败: %', SQLERRM;
    END;
    
    -- 返回适当的记录
    IF TG_OP = 'DELETE' THEN
        RETURN OLD;
    ELSE
        RETURN NEW;
    END IF;
END;
$$ LANGUAGE plpgsql;

\echo '为核心表创建审计触发器...'

-- 用户表审计触发器
CREATE TRIGGER trigger_audit_users
    AFTER INSERT OR UPDATE OR DELETE ON users
    FOR EACH ROW
    EXECUTE FUNCTION audit_trigger_function();

-- 角色表审计触发器
CREATE TRIGGER trigger_audit_roles
    AFTER INSERT OR UPDATE OR DELETE ON roles
    FOR EACH ROW
    EXECUTE FUNCTION audit_trigger_function();

-- 权限表审计触发器
CREATE TRIGGER trigger_audit_permissions
    AFTER INSERT OR UPDATE OR DELETE ON permissions
    FOR EACH ROW
    EXECUTE FUNCTION audit_trigger_function();

-- 用户角色关联表审计触发器
CREATE TRIGGER trigger_audit_user_roles
    AFTER INSERT OR UPDATE OR DELETE ON user_roles
    FOR EACH ROW
    EXECUTE FUNCTION audit_trigger_function();

-- 角色权限关联表审计触发器
CREATE TRIGGER trigger_audit_role_permissions
    AFTER INSERT OR UPDATE OR DELETE ON role_permissions
    FOR EACH ROW
    EXECUTE FUNCTION audit_trigger_function();

-- 组织表审计触发器
CREATE TRIGGER trigger_audit_organizations
    AFTER INSERT OR UPDATE OR DELETE ON organizations
    FOR EACH ROW
    EXECUTE FUNCTION audit_trigger_function();

-- 用户组织关联表审计触发器
CREATE TRIGGER trigger_audit_user_organizations
    AFTER INSERT OR UPDATE OR DELETE ON user_organizations
    FOR EACH ROW
    EXECUTE FUNCTION audit_trigger_function();

\echo '创建数据完整性触发器...'

-- 用户删除级联处理触发器
CREATE OR REPLACE FUNCTION cascade_user_deletion()
RETURNS TRIGGER AS $$
BEGIN
    -- 软删除用户时，同时处理关联数据
    IF NEW.deleted_at IS NOT NULL AND OLD.deleted_at IS NULL THEN
        -- 停用用户的角色关联
        UPDATE user_roles 
        SET is_active = false, updated_at = CURRENT_TIMESTAMP
        WHERE user_id = NEW.id AND is_active = true;
        
        -- 停用用户的组织关联
        UPDATE user_organizations 
        SET is_active = false, updated_at = CURRENT_TIMESTAMP
        WHERE user_id = NEW.id AND is_active = true;
        
        -- 停用用户的会话
        UPDATE user_sessions 
        SET is_active = false, updated_at = CURRENT_TIMESTAMP
        WHERE user_id = NEW.id AND is_active = true;
    END IF;
    
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_cascade_user_deletion
    AFTER UPDATE ON users
    FOR EACH ROW
    EXECUTE FUNCTION cascade_user_deletion();

-- 角色删除级联处理触发器
CREATE OR REPLACE FUNCTION cascade_role_deletion()
RETURNS TRIGGER AS $$
BEGIN
    -- 软删除角色时，同时处理关联数据
    IF NEW.deleted_at IS NOT NULL AND OLD.deleted_at IS NULL THEN
        -- 停用角色的用户关联
        UPDATE user_roles 
        SET is_active = false, updated_at = CURRENT_TIMESTAMP
        WHERE role_id = NEW.id AND is_active = true;
    END IF;
    
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_cascade_role_deletion
    AFTER UPDATE ON roles
    FOR EACH ROW
    EXECUTE FUNCTION cascade_role_deletion();

\echo '所有审计触发器创建完成！'