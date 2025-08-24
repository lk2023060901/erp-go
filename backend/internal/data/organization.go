package data

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"erp-system/internal/biz"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/lib/pq"
)

// organizationRepo 组织仓储实现
type organizationRepo struct {
	data *Data
	log  *log.Helper
}

// NewOrganizationRepo 创建组织仓储
func NewOrganizationRepo(data *Data, logger log.Logger) biz.OrganizationRepo {
	return &organizationRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

// CreateOrganization 创建组织
func (r *organizationRepo) CreateOrganization(ctx context.Context, org *biz.Organization) (*biz.Organization, error) {
	var id int32
	query := `
		INSERT INTO organizations (parent_id, name, code, description, is_enabled, sort_order, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id`

	err := r.data.db.QueryRowContext(ctx, query,
		org.ParentID, org.Name, org.Code, org.Description,
		org.IsEnabled, org.SortOrder, org.CreatedAt, org.UpdatedAt,
	).Scan(&id)

	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			if pqErr.Code == "23505" { // unique constraint violation
				if pqErr.Constraint == "organizations_code_key" {
					return nil, biz.ErrOrganizationCodeExists
				}
			}
		}
		r.log.Errorf("failed to create organization: %v", err)
		return nil, err
	}

	org.ID = id
	return org, nil
}

// GetOrganization 获取组织
func (r *organizationRepo) GetOrganization(ctx context.Context, id int32) (*biz.Organization, error) {
	var org biz.Organization
	var parentID sql.NullInt32

	query := `
		SELECT id, parent_id, name, code, description, is_enabled, sort_order, created_at, updated_at
		FROM organizations WHERE id = $1`

	err := r.data.db.QueryRowContext(ctx, query, id).Scan(
		&org.ID, &parentID, &org.Name, &org.Code, &org.Description,
		&org.IsEnabled, &org.SortOrder, &org.CreatedAt, &org.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("organization not found")
		}
		r.log.Errorf("failed to get organization: %v", err)
		return nil, err
	}

	if parentID.Valid {
		org.ParentID = &parentID.Int32
	}

	return &org, nil
}

// UpdateOrganization 更新组织
func (r *organizationRepo) UpdateOrganization(ctx context.Context, org *biz.Organization) (*biz.Organization, error) {
	query := `
		UPDATE organizations 
		SET parent_id = $2, name = $3, description = $4, is_enabled = $5, sort_order = $6, updated_at = $7
		WHERE id = $1`

	org.UpdatedAt = time.Now()
	_, err := r.data.db.ExecContext(ctx, query,
		org.ID, org.ParentID, org.Name, org.Description,
		org.IsEnabled, org.SortOrder, org.UpdatedAt,
	)

	if err != nil {
		r.log.Errorf("failed to update organization: %v", err)
		return nil, err
	}

	return org, nil
}

// DeleteOrganization 删除组织
func (r *organizationRepo) DeleteOrganization(ctx context.Context, id int32) error {
	// 检查是否有子组织
	var childCount int32
	childQuery := "SELECT COUNT(*) FROM organizations WHERE parent_id = $1"
	if err := r.data.db.QueryRowContext(ctx, childQuery, id).Scan(&childCount); err != nil {
		r.log.Errorf("failed to check child organizations: %v", err)
		return err
	}

	if childCount > 0 {
		return biz.ErrOrganizationHasChildren
	}

	// 检查是否有用户
	var userCount int32
	userQuery := "SELECT COUNT(*) FROM user_organizations WHERE organization_id = $1"
	if err := r.data.db.QueryRowContext(ctx, userQuery, id).Scan(&userCount); err != nil {
		r.log.Errorf("failed to check organization users: %v", err)
		return err
	}

	if userCount > 0 {
		return biz.ErrOrganizationHasUsers
	}

	// 删除组织
	query := "DELETE FROM organizations WHERE id = $1"
	_, err := r.data.db.ExecContext(ctx, query, id)
	if err != nil {
		r.log.Errorf("failed to delete organization: %v", err)
		return err
	}

	return nil
}

// GetOrganizationTree 获取组织树
func (r *organizationRepo) GetOrganizationTree(ctx context.Context) ([]*biz.Organization, error) {
	// 获取所有组织
	orgs, err := r.ListAllOrganizations(ctx)
	if err != nil {
		return nil, err
	}

	// 构建组织映射
	orgMap := make(map[int32]*biz.Organization)
	for _, org := range orgs {
		orgMap[org.ID] = org
	}

	// 构建树形结构
	var rootOrgs []*biz.Organization
	for _, org := range orgs {
		if org.ParentID == nil {
			// 根节点
			rootOrgs = append(rootOrgs, org)
		} else {
			// 子节点，添加到父节点
			if parent, exists := orgMap[*org.ParentID]; exists {
				if parent.Children == nil {
					parent.Children = make([]*biz.Organization, 0)
				}
				parent.Children = append(parent.Children, org)
			}
		}
	}

	return rootOrgs, nil
}

// ListAllOrganizations 获取所有组织
func (r *organizationRepo) ListAllOrganizations(ctx context.Context) ([]*biz.Organization, error) {
	query := `
		SELECT id, parent_id, name, code, description, is_enabled, sort_order, created_at, updated_at
		FROM organizations 
		ORDER BY sort_order, created_at`

	rows, err := r.data.db.QueryContext(ctx, query)
	if err != nil {
		r.log.Errorf("failed to list all organizations: %v", err)
		return nil, err
	}
	defer rows.Close()

	var organizations []*biz.Organization
	for rows.Next() {
		var org biz.Organization
		var parentID sql.NullInt32

		err := rows.Scan(
			&org.ID, &parentID, &org.Name, &org.Code, &org.Description,
			&org.IsEnabled, &org.SortOrder, &org.CreatedAt, &org.UpdatedAt,
		)
		if err != nil {
			r.log.Errorf("failed to scan organization: %v", err)
			return nil, err
		}

		if parentID.Valid {
			org.ParentID = &parentID.Int32
		}

		organizations = append(organizations, &org)
	}

	return organizations, nil
}

// GetOrganizationByCode 根据编码获取组织
func (r *organizationRepo) GetOrganizationByCode(ctx context.Context, code string) (*biz.Organization, error) {
	var org biz.Organization
	var parentID sql.NullInt32

	query := `
		SELECT id, parent_id, name, code, description, is_enabled, sort_order, created_at, updated_at
		FROM organizations WHERE code = $1`

	err := r.data.db.QueryRowContext(ctx, query, code).Scan(
		&org.ID, &parentID, &org.Name, &org.Code, &org.Description,
		&org.IsEnabled, &org.SortOrder, &org.CreatedAt, &org.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("organization not found")
		}
		r.log.Errorf("failed to get organization by code: %v", err)
		return nil, err
	}

	if parentID.Valid {
		org.ParentID = &parentID.Int32
	}

	return &org, nil
}

// GetChildOrganizations 获取子组织
func (r *organizationRepo) GetChildOrganizations(ctx context.Context, parentID int32) ([]*biz.Organization, error) {
	query := `
		SELECT id, parent_id, name, code, description, is_enabled, sort_order, created_at, updated_at
		FROM organizations 
		WHERE parent_id = $1
		ORDER BY sort_order, created_at`

	rows, err := r.data.db.QueryContext(ctx, query, parentID)
	if err != nil {
		r.log.Errorf("failed to get child organizations: %v", err)
		return nil, err
	}
	defer rows.Close()

	var organizations []*biz.Organization
	for rows.Next() {
		var org biz.Organization
		var childParentID sql.NullInt32

		err := rows.Scan(
			&org.ID, &childParentID, &org.Name, &org.Code, &org.Description,
			&org.IsEnabled, &org.SortOrder, &org.CreatedAt, &org.UpdatedAt,
		)
		if err != nil {
			r.log.Errorf("failed to scan child organization: %v", err)
			return nil, err
		}

		if childParentID.Valid {
			org.ParentID = &childParentID.Int32
		}

		organizations = append(organizations, &org)
	}

	return organizations, nil
}

// GetEnabledOrganizations 获取启用的组织
func (r *organizationRepo) GetEnabledOrganizations(ctx context.Context) ([]*biz.Organization, error) {
	query := `
		SELECT id, parent_id, name, code, description, is_enabled, sort_order, created_at, updated_at
		FROM organizations 
		WHERE is_enabled = true
		ORDER BY sort_order, created_at`

	rows, err := r.data.db.QueryContext(ctx, query)
	if err != nil {
		r.log.Errorf("failed to get enabled organizations: %v", err)
		return nil, err
	}
	defer rows.Close()

	var organizations []*biz.Organization
	for rows.Next() {
		var org biz.Organization
		var parentID sql.NullInt32

		err := rows.Scan(
			&org.ID, &parentID, &org.Name, &org.Code, &org.Description,
			&org.IsEnabled, &org.SortOrder, &org.CreatedAt, &org.UpdatedAt,
		)
		if err != nil {
			r.log.Errorf("failed to scan enabled organization: %v", err)
			return nil, err
		}

		if parentID.Valid {
			org.ParentID = &parentID.Int32
		}

		organizations = append(organizations, &org)
	}

	return organizations, nil
}

// GetOrganizationUsers 获取组织用户
func (r *organizationRepo) GetOrganizationUsers(ctx context.Context, orgID int32) ([]*biz.User, error) {
	query := `
		SELECT u.id, u.username, u.email, u.first_name, u.last_name, u.phone, u.gender,
		       u.avatar_url, u.is_active, u.two_factor_enabled,
		       u.last_login_at, u.last_login_ip, u.login_count, u.created_at, u.updated_at
		FROM users u
		INNER JOIN user_organizations uo ON u.id = uo.user_id
		WHERE uo.organization_id = $1 AND u.is_active = true
		ORDER BY u.created_at DESC`

	rows, err := r.data.db.QueryContext(ctx, query, orgID)
	if err != nil {
		r.log.Errorf("failed to get organization users: %v", err)
		return nil, err
	}
	defer rows.Close()

	var users []*biz.User
	for rows.Next() {
		var user biz.User
		var lastLoginAt sql.NullTime

		err := rows.Scan(
			&user.ID, &user.Username, &user.Email, &user.FirstName, &user.LastName,
			&user.Phone, &user.Gender, &user.AvatarURL, &user.IsActive,
			&user.TwoFactorEnabled, &lastLoginAt, &user.LastLoginIP, &user.LoginCount,
			&user.CreatedAt, &user.UpdatedAt,
		)
		if err != nil {
			r.log.Errorf("failed to scan organization user: %v", err)
			return nil, err
		}

		if lastLoginAt.Valid {
			user.LastLoginAt = lastLoginAt.Time
		}

		users = append(users, &user)
	}

	return users, nil
}

// AssignUsers 分配用户到组织
func (r *organizationRepo) AssignUsers(ctx context.Context, orgID int32, userIDs []int32) error {
	// 开启事务
	tx, err := r.data.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// 删除现有关联
	_, err = tx.ExecContext(ctx, "DELETE FROM user_organizations WHERE organization_id = $1", orgID)
	if err != nil {
		r.log.Errorf("failed to delete organization users: %v", err)
		return err
	}

	// 添加新关联
	for _, userID := range userIDs {
		_, err = tx.ExecContext(ctx, 
			"INSERT INTO user_organizations (user_id, organization_id, created_at) VALUES ($1, $2, $3)",
			userID, orgID, time.Now())
		if err != nil {
			r.log.Errorf("failed to assign user to organization: %v", err)
			return err
		}
	}

	return tx.Commit()
}