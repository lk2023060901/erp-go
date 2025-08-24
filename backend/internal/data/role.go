package data

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"erp-system/internal/biz"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/lib/pq"
)

// roleRepo 角色仓储实现
type roleRepo struct {
	data *Data
	log  *log.Helper
}

// NewRoleRepo 创建角色仓储
func NewRoleRepo(data *Data, logger log.Logger) biz.RoleRepo {
	return &roleRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

// CreateRole 创建角色
func (r *roleRepo) CreateRole(ctx context.Context, role *biz.Role) (*biz.Role, error) {
	var id int32
	query := `
		INSERT INTO roles (name, code, description, is_system_role, is_enabled, sort_order, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id`

	err := r.data.db.QueryRowContext(ctx, query,
		role.Name, role.Code, role.Description, role.IsSystemRole,
		role.IsEnabled, role.SortOrder, role.CreatedAt, role.UpdatedAt,
	).Scan(&id)

	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			if pqErr.Code == "23505" { // unique constraint violation
				if pqErr.Constraint == "roles_code_key" {
					return nil, biz.ErrRoleCodeExists
				}
				if pqErr.Constraint == "roles_name_key" {
					return nil, biz.ErrRoleNameExists
				}
			}
		}
		r.log.Errorf("failed to create role: %v", err)
		return nil, err
	}

	role.ID = id
	return role, nil
}

// GetRole 获取角色
func (r *roleRepo) GetRole(ctx context.Context, id int32) (*biz.Role, error) {
	var role biz.Role
	query := `
		SELECT id, name, code, description, is_system_role, is_enabled, sort_order, created_at, updated_at
		FROM roles WHERE id = $1`

	err := r.data.db.QueryRowContext(ctx, query, id).Scan(
		&role.ID, &role.Name, &role.Code, &role.Description,
		&role.IsSystemRole, &role.IsEnabled, &role.SortOrder,
		&role.CreatedAt, &role.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("role not found")
		}
		r.log.Errorf("failed to get role: %v", err)
		return nil, err
	}

	return &role, nil
}

// UpdateRole 更新角色
func (r *roleRepo) UpdateRole(ctx context.Context, role *biz.Role) (*biz.Role, error) {
	query := `
		UPDATE roles 
		SET name = $2, description = $3, is_enabled = $4, sort_order = $5, updated_at = $6
		WHERE id = $1`

	role.UpdatedAt = time.Now()
	_, err := r.data.db.ExecContext(ctx, query,
		role.ID, role.Name, role.Description, role.IsEnabled,
		role.SortOrder, role.UpdatedAt,
	)

	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			if pqErr.Code == "23505" { // unique constraint violation
				if pqErr.Constraint == "roles_name_key" {
					return nil, biz.ErrRoleNameExists
				}
			}
		}
		r.log.Errorf("failed to update role: %v", err)
		return nil, err
	}

	return role, nil
}

// DeleteRole 删除角色
func (r *roleRepo) DeleteRole(ctx context.Context, id int32) error {
	// 检查是否为系统角色
	role, err := r.GetRole(ctx, id)
	if err != nil {
		return err
	}
	
	if role.IsSystemRole {
		return biz.ErrCannotDeleteSystemRole
	}

	// 检查是否有用户使用此角色
	var userCount int32
	countQuery := "SELECT COUNT(*) FROM user_roles WHERE role_id = $1"
	if err := r.data.db.QueryRowContext(ctx, countQuery, id).Scan(&userCount); err != nil {
		r.log.Errorf("failed to check role usage: %v", err)
		return err
	}

	if userCount > 0 {
		return biz.ErrRoleInUse
	}

	// 开启事务
	tx, err := r.data.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// 删除角色权限关联
	_, err = tx.ExecContext(ctx, "DELETE FROM role_permissions WHERE role_id = $1", id)
	if err != nil {
		r.log.Errorf("failed to delete role permissions: %v", err)
		return err
	}

	// 删除角色
	_, err = tx.ExecContext(ctx, "DELETE FROM roles WHERE id = $1", id)
	if err != nil {
		r.log.Errorf("failed to delete role: %v", err)
		return err
	}

	return tx.Commit()
}

// ListRoles 角色列表
func (r *roleRepo) ListRoles(ctx context.Context, page, size int32, search string, isEnabled *bool, sortField, sortOrder string) ([]*biz.Role, int32, error) {
	offset := (page - 1) * size
	var roles []*biz.Role
	var total int32

	// 构建WHERE条件
	whereClause := ""
	var whereArgs []interface{}
	argIndex := 0

	if search != "" {
		whereClause = "WHERE name ILIKE $" + fmt.Sprintf("%d", argIndex+1) + " OR code ILIKE $" + fmt.Sprintf("%d", argIndex+2)
		whereArgs = append(whereArgs, "%"+search+"%", "%"+search+"%")
		argIndex += 2
	}

	if isEnabled != nil {
		if whereClause == "" {
			whereClause = "WHERE is_enabled = $" + fmt.Sprintf("%d", argIndex+1)
		} else {
			whereClause += " AND is_enabled = $" + fmt.Sprintf("%d", argIndex+1)
		}
		whereArgs = append(whereArgs, *isEnabled)
		argIndex++
	}

	// 查询总数
	countQuery := "SELECT COUNT(*) FROM roles " + whereClause
	err := r.data.db.QueryRowContext(ctx, countQuery, whereArgs...).Scan(&total)
	if err != nil {
		r.log.Errorf("failed to count roles: %v", err)
		return nil, 0, err
	}

	// 构建ORDER BY子句
	orderByClause := "ORDER BY sort_order, created_at DESC" // 默认排序
	if sortField != "" && sortOrder != "" {
		// 映射前端字段到数据库字段
		dbField := sortField
		switch sortField {
		case "name":
			dbField = "name"
		case "code":
			dbField = "code"
		case "is_system_role":
			dbField = "is_system_role"
		case "is_enabled":
			dbField = "is_enabled"
		case "sort_order":
			dbField = "sort_order"
		case "created_at":
			dbField = "created_at"
		case "updated_at":
			dbField = "updated_at"
		default:
			// 如果字段无效，使用默认排序
			dbField = ""
		}
		
		if dbField != "" {
			if sortOrder == "asc" || sortOrder == "desc" {
				orderByClause = fmt.Sprintf("ORDER BY %s %s, created_at DESC", dbField, strings.ToUpper(sortOrder))
			}
		}
	}

	// 查询数据
	query := fmt.Sprintf(`
		SELECT id, name, code, description, is_system_role, is_enabled, sort_order, created_at, updated_at
		FROM roles 
		%s
		%s 
		LIMIT $%d OFFSET $%d`, whereClause, orderByClause, argIndex+1, argIndex+2)

	// 添加分页参数
	whereArgs = append(whereArgs, size, offset)

	rows, err := r.data.db.QueryContext(ctx, query, whereArgs...)
	if err != nil {
		r.log.Errorf("failed to list roles: %v", err)
		return nil, 0, err
	}
	defer rows.Close()

	for rows.Next() {
		var role biz.Role
		err := rows.Scan(
			&role.ID, &role.Name, &role.Code, &role.Description,
			&role.IsSystemRole, &role.IsEnabled, &role.SortOrder,
			&role.CreatedAt, &role.UpdatedAt,
		)
		if err != nil {
			r.log.Errorf("failed to scan role: %v", err)
			return nil, 0, err
		}

		roles = append(roles, &role)
	}

	return roles, total, nil
}

// GetRolePermissions 获取角色权限
func (r *roleRepo) GetRolePermissions(ctx context.Context, roleID int32) ([]*biz.Permission, error) {
	query := `
		SELECT p.id, p.parent_id, p.name, p.code, p.resource, p.action, p.module,
		       p.description, p.is_menu, p.path, p.sort_order,
		       p.created_at, p.updated_at
		FROM permissions p
		INNER JOIN role_permissions rp ON p.id = rp.permission_id
		WHERE rp.role_id = $1
		ORDER BY p.sort_order`

	rows, err := r.data.db.QueryContext(ctx, query, roleID)
	if err != nil {
		r.log.Errorf("failed to get role permissions: %v", err)
		return nil, err
	}
	defer rows.Close()

	var permissions []*biz.Permission
	for rows.Next() {
		var permission biz.Permission
		var parentID sql.NullInt32
		var path, description sql.NullString

		err := rows.Scan(
			&permission.ID, &parentID, &permission.Name, &permission.Code,
			&permission.Resource, &permission.Action, &permission.Module,
			&description, &permission.IsMenu, &path,
			&permission.SortOrder, &permission.CreatedAt,
			&permission.UpdatedAt,
		)
		if err != nil {
			r.log.Errorf("failed to scan permission: %v", err)
			return nil, err
		}

		if parentID.Valid {
			permission.ParentID = &parentID.Int32
		}
		
		// 处理可能为NULL的字符串字段
		if description.Valid {
			permission.Description = description.String
		}
		if path.Valid {
			permission.Path = path.String
		}

		permissions = append(permissions, &permission)
	}

	return permissions, nil
}

// AssignPermissions 分配权限
func (r *roleRepo) AssignPermissions(ctx context.Context, roleID int32, permissionIDs []int32) error {
	// 开启事务
	tx, err := r.data.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// 删除现有权限
	_, err = tx.ExecContext(ctx, "DELETE FROM role_permissions WHERE role_id = $1", roleID)
	if err != nil {
		r.log.Errorf("failed to delete role permissions: %v", err)
		return err
	}

	// 添加新权限
	for _, permissionID := range permissionIDs {
		_, err = tx.ExecContext(ctx, 
			"INSERT INTO role_permissions (role_id, permission_id, created_at) VALUES ($1, $2, $3)",
			roleID, permissionID, time.Now())
		if err != nil {
			r.log.Errorf("failed to assign permission: %v", err)
			return err
		}
	}

	return tx.Commit()
}

// GetRoleByCode 根据编码获取角色
func (r *roleRepo) GetRoleByCode(ctx context.Context, code string) (*biz.Role, error) {
	var role biz.Role
	query := `
		SELECT id, name, code, description, is_system_role, is_enabled, sort_order, created_at, updated_at
		FROM roles WHERE code = $1`

	err := r.data.db.QueryRowContext(ctx, query, code).Scan(
		&role.ID, &role.Name, &role.Code, &role.Description,
		&role.IsSystemRole, &role.IsEnabled, &role.SortOrder,
		&role.CreatedAt, &role.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("role not found")
		}
		r.log.Errorf("failed to get role by code: %v", err)
		return nil, err
	}

	return &role, nil
}

// GetEnabledRoles 获取启用的角色
func (r *roleRepo) GetEnabledRoles(ctx context.Context) ([]*biz.Role, error) {
	query := `
		SELECT id, name, code, description, is_system_role, is_enabled, sort_order, created_at, updated_at
		FROM roles 
		WHERE is_enabled = true
		ORDER BY sort_order, created_at`

	rows, err := r.data.db.QueryContext(ctx, query)
	if err != nil {
		r.log.Errorf("failed to get enabled roles: %v", err)
		return nil, err
	}
	defer rows.Close()

	var roles []*biz.Role
	for rows.Next() {
		var role biz.Role
		err := rows.Scan(
			&role.ID, &role.Name, &role.Code, &role.Description,
			&role.IsSystemRole, &role.IsEnabled, &role.SortOrder,
			&role.CreatedAt, &role.UpdatedAt,
		)
		if err != nil {
			r.log.Errorf("failed to scan role: %v", err)
			return nil, err
		}

		roles = append(roles, &role)
	}

	return roles, nil
}