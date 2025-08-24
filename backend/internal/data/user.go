package data

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"erp-system/internal/biz"
	"erp-system/internal/pkg"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/lib/pq"
)

// userRepo 用户仓储实现
type userRepo struct {
	data *Data
	log  *log.Helper
	pm   *pkg.PasswordManager
}

// NewUserRepo 创建用户仓储
func NewUserRepo(data *Data, logger log.Logger) biz.UserRepo {
	return &userRepo{
		data: data,
		log:  log.NewHelper(logger),
		pm:   pkg.NewPasswordManager(),
	}
}

// CreateUser 创建用户
func (r *userRepo) CreateUser(ctx context.Context, user *biz.User) (*biz.User, error) {
	var id int32
	query := `
		INSERT INTO users (username, email, password_hash, salt, first_name, last_name, phone, gender, birth_date, 
		                  avatar_url, is_enabled, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
		RETURNING id`

	// 处理可选的gender字段
	var gender interface{} = nil
	if user.Gender != "" {
		gender = user.Gender
	}

	// 处理可选的phone字段
	var phone interface{} = nil
	if user.Phone != "" {
		phone = user.Phone
	}

	err := r.data.db.QueryRowContext(ctx, query,
		user.Username, user.Email, user.Password, "", user.FirstName, user.LastName,
		phone, gender, user.BirthDate, user.AvatarURL, user.IsActive,
		user.CreatedAt, user.UpdatedAt,
	).Scan(&id)

	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			if pqErr.Code == "23505" { // unique constraint violation
				if strings.Contains(pqErr.Detail, "username") {
					return nil, biz.ErrUsernameExists
				}
				if strings.Contains(pqErr.Detail, "email") {
					return nil, biz.ErrEmailExists
				}
			}
		}
		r.log.Errorf("failed to create user: %v", err)
		return nil, err
	}

	user.ID = id
	return user, nil
}

// GetUser 获取用户
func (r *userRepo) GetUser(ctx context.Context, id int32) (*biz.User, error) {
	var user biz.User
	var phone, gender, avatarURL, twoFactorSecret, lastLoginIP sql.NullString
	var birthDate, lastLoginAt sql.NullTime

	query := `
		SELECT id, username, email, password_hash, first_name, last_name, phone, gender, birth_date,
		       avatar_url, is_enabled, two_factor_enabled, two_factor_secret,
		       last_login_time, last_login_ip, login_count, created_at, updated_at
		FROM users WHERE id = $1`

	err := r.data.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID, &user.Username, &user.Email, &user.Password, &user.FirstName,
		&user.LastName, &phone, &gender, &birthDate, &avatarURL,
		&user.IsActive, &user.TwoFactorEnabled, &twoFactorSecret,
		&lastLoginAt, &lastLoginIP, &user.LoginCount, &user.CreatedAt, &user.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		r.log.Errorf("failed to get user: %v", err)
		return nil, err
	}

	// 处理可选字段
	if phone.Valid {
		user.Phone = phone.String
	}
	if gender.Valid {
		user.Gender = gender.String
	}
	if avatarURL.Valid {
		user.AvatarURL = avatarURL.String
	}
	if twoFactorSecret.Valid {
		user.TwoFactorSecret = twoFactorSecret.String
	}
	if lastLoginIP.Valid {
		user.LastLoginIP = lastLoginIP.String
	}
	if birthDate.Valid {
		user.BirthDate = birthDate.Time
	}
	if lastLoginAt.Valid {
		user.LastLoginAt = lastLoginAt.Time
	}

	return &user, nil
}

// GetUserByUsername 根据用户名获取用户
func (r *userRepo) GetUserByUsername(ctx context.Context, username string) (*biz.User, error) {
	var user biz.User
	var phone, gender, avatarURL, twoFactorSecret, lastLoginIP sql.NullString
	var birthDate, lastLoginAt sql.NullTime

	query := `
		SELECT id, username, email, password_hash, first_name, last_name, phone, gender, birth_date,
		       avatar_url, is_enabled, two_factor_enabled, two_factor_secret,
		       last_login_time, last_login_ip, login_count, created_at, updated_at
		FROM users WHERE username = $1`

	err := r.data.db.QueryRowContext(ctx, query, username).Scan(
		&user.ID, &user.Username, &user.Email, &user.Password, &user.FirstName,
		&user.LastName, &phone, &gender, &birthDate, &avatarURL,
		&user.IsActive, &user.TwoFactorEnabled, &twoFactorSecret,
		&lastLoginAt, &lastLoginIP, &user.LoginCount, &user.CreatedAt, &user.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		r.log.Errorf("failed to get user by username: %v", err)
		return nil, err
	}

	// 处理可选字段
	if phone.Valid {
		user.Phone = phone.String
	}
	if gender.Valid {
		user.Gender = gender.String
	}
	if avatarURL.Valid {
		user.AvatarURL = avatarURL.String
	}
	if twoFactorSecret.Valid {
		user.TwoFactorSecret = twoFactorSecret.String
	}
	if lastLoginIP.Valid {
		user.LastLoginIP = lastLoginIP.String
	}
	if birthDate.Valid {
		user.BirthDate = birthDate.Time
	}
	if lastLoginAt.Valid {
		user.LastLoginAt = lastLoginAt.Time
	}

	return &user, nil
}

// GetUserByEmail 根据邮箱获取用户
func (r *userRepo) GetUserByEmail(ctx context.Context, email string) (*biz.User, error) {
	var user biz.User
	var phone, gender, avatarURL, twoFactorSecret, lastLoginIP sql.NullString
	var birthDate, lastLoginAt sql.NullTime

	query := `
		SELECT id, username, email, password_hash, first_name, last_name, phone, gender, birth_date,
		       avatar_url, is_enabled, two_factor_enabled, two_factor_secret,
		       last_login_time, last_login_ip, login_count, created_at, updated_at
		FROM users WHERE email = $1`

	err := r.data.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID, &user.Username, &user.Email, &user.Password, &user.FirstName,
		&user.LastName, &phone, &gender, &birthDate, &avatarURL,
		&user.IsActive, &user.TwoFactorEnabled, &twoFactorSecret,
		&lastLoginAt, &lastLoginIP, &user.LoginCount, &user.CreatedAt, &user.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		r.log.Errorf("failed to get user by email: %v", err)
		return nil, err
	}

	// 处理NULL字段
	if phone.Valid {
		user.Phone = phone.String
	}
	if gender.Valid {
		user.Gender = gender.String
	}
	if avatarURL.Valid {
		user.AvatarURL = avatarURL.String
	}
	if twoFactorSecret.Valid {
		user.TwoFactorSecret = twoFactorSecret.String
	}
	if lastLoginIP.Valid {
		user.LastLoginIP = lastLoginIP.String
	}
	if birthDate.Valid {
		user.BirthDate = birthDate.Time
	}
	if lastLoginAt.Valid {
		user.LastLoginAt = lastLoginAt.Time
	}

	return &user, nil
}

// UpdateUser 更新用户
func (r *userRepo) UpdateUser(ctx context.Context, user *biz.User) (*biz.User, error) {
	query := `
		UPDATE users 
		SET email = $1, first_name = $2, last_name = $3, phone = $4, gender = $5, 
		    birth_date = $6, avatar_url = $7, is_enabled = $8, updated_at = $9
		WHERE id = $10`

	user.UpdatedAt = time.Now()
	_, err := r.data.db.ExecContext(ctx, query,
		user.Email, user.FirstName, user.LastName, user.Phone,
		user.Gender, user.BirthDate, user.AvatarURL, user.IsActive,
		user.UpdatedAt, user.ID,
	)

	if err != nil {
		r.log.Errorf("failed to update user: %v", err)
		return nil, err
	}

	return user, nil
}

// UpdatePassword 更新用户密码
func (r *userRepo) UpdatePassword(ctx context.Context, id int32, hashedPassword string) error {
	query := "UPDATE users SET password_hash = $1, updated_at = CURRENT_TIMESTAMP WHERE id = $2"
	_, err := r.data.db.ExecContext(ctx, query, hashedPassword, id)
	if err != nil {
		r.log.Errorf("failed to update password: %v", err)
		return err
	}
	return nil
}

// DeleteUser 删除用户
func (r *userRepo) DeleteUser(ctx context.Context, id int32) error {
	query := "DELETE FROM users WHERE id = $1"
	_, err := r.data.db.ExecContext(ctx, query, id)
	if err != nil {
		r.log.Errorf("failed to delete user: %v", err)
		return err
	}
	return nil
}

// UserListOptions 用户列表查询选项
type UserListOptions struct {
	Page             int32
	Size             int32
	Search           string                 // 兼容旧版搜索
	FilterConditions map[string]interface{} // 新的过滤条件
	SortConfig       map[string]interface{} // 排序配置
}

// ListUsersWithOptions 用户列表（新版本）
func (r *userRepo) ListUsersWithOptions(ctx context.Context, options *UserListOptions) ([]*biz.User, int32, error) {
	offset := (options.Page - 1) * options.Size
	var users []*biz.User
	var total int32

	// 构建查询条件
	whereClause, args := r.buildUserFilterQuery(options)
	
	// 构建排序条件
	orderClause := r.buildUserSortQuery(options.SortConfig)

	// 查询总数
	countQuery := "SELECT COUNT(*) FROM users " + whereClause
	err := r.data.db.QueryRowContext(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		r.log.Errorf("failed to count users: %v", err)
		return nil, 0, err
	}

	// 查询数据
	query := `
		SELECT id, username, email, first_name, last_name, phone, gender, birth_date,
		       avatar_url, is_enabled, two_factor_enabled,
		       last_login_time, last_login_ip, login_count, created_at, updated_at
		FROM users ` + whereClause + ` ` + orderClause + ` 
		LIMIT $` + fmt.Sprintf("%d", len(args)+1) + ` OFFSET $` + fmt.Sprintf("%d", len(args)+2)

	args = append(args, options.Size, offset)

	rows, err := r.data.db.QueryContext(ctx, query, args...)
	if err != nil {
		r.log.Errorf("failed to list users: %v", err)
		return nil, 0, err
	}
	defer rows.Close()

	for rows.Next() {
		var user biz.User
		var phone, gender, avatarURL, lastLoginIP sql.NullString
		var birthDate, lastLoginAt sql.NullTime

		err := rows.Scan(
			&user.ID, &user.Username, &user.Email, &user.FirstName, &user.LastName,
			&phone, &gender, &birthDate, &avatarURL, &user.IsActive,
			&user.TwoFactorEnabled, &lastLoginAt, &lastLoginIP,
			&user.LoginCount, &user.CreatedAt, &user.UpdatedAt,
		)
		if err != nil {
			r.log.Errorf("failed to scan user: %v", err)
			return nil, 0, err
		}

		if phone.Valid {
			user.Phone = phone.String
		}
		if gender.Valid {
			user.Gender = gender.String
		}
		if avatarURL.Valid {
			user.AvatarURL = avatarURL.String
		}
		if lastLoginIP.Valid {
			user.LastLoginIP = lastLoginIP.String
		}
		if birthDate.Valid {
			user.BirthDate = birthDate.Time
		}
		if lastLoginAt.Valid {
			user.LastLoginAt = lastLoginAt.Time
		}

		users = append(users, &user)
	}

	return users, total, nil
}

// buildUserFilterQuery 构建用户过滤查询条件
func (r *userRepo) buildUserFilterQuery(options *UserListOptions) (string, []interface{}) {
	whereClause := "WHERE 1=1"
	args := []interface{}{}
	
	// 兼容旧版搜索
	if options.Search != "" {
		whereClause += fmt.Sprintf(" AND (username ILIKE $%d OR email ILIKE $%d OR first_name ILIKE $%d OR last_name ILIKE $%d)", 
			len(args)+1, len(args)+2, len(args)+3, len(args)+4)
		searchPattern := "%" + options.Search + "%"
		args = append(args, searchPattern, searchPattern, searchPattern, searchPattern)
	}
	
	// 处理新的过滤条件
	if options.FilterConditions != nil {
		if conditions, ok := options.FilterConditions["conditions"].([]interface{}); ok {
			for _, condition := range conditions {
				if condMap, ok := condition.(map[string]interface{}); ok {
					field := condMap["field"].(string)
					operator := condMap["operator"].(string)
					value := condMap["value"]
					
					switch operator {
					case "equals":
						whereClause += fmt.Sprintf(" AND %s = $%d", field, len(args)+1)
						args = append(args, value)
					case "contains":
						whereClause += fmt.Sprintf(" AND %s ILIKE $%d", field, len(args)+1)
						args = append(args, "%"+fmt.Sprintf("%v", value)+"%")
					case "starts_with":
						whereClause += fmt.Sprintf(" AND %s ILIKE $%d", field, len(args)+1)
						args = append(args, fmt.Sprintf("%v", value)+"%")
					case "ends_with":
						whereClause += fmt.Sprintf(" AND %s ILIKE $%d", field, len(args)+1)
						args = append(args, "%"+fmt.Sprintf("%v", value))
					case "greater_than":
						whereClause += fmt.Sprintf(" AND %s > $%d", field, len(args)+1)
						args = append(args, value)
					case "less_than":
						whereClause += fmt.Sprintf(" AND %s < $%d", field, len(args)+1)
						args = append(args, value)
					case "in":
						if valueArr, ok := value.([]interface{}); ok && len(valueArr) > 0 {
							placeholders := make([]string, len(valueArr))
							for i, v := range valueArr {
								placeholders[i] = fmt.Sprintf("$%d", len(args)+1)
								args = append(args, v)
							}
							whereClause += fmt.Sprintf(" AND %s IN (%s)", field, strings.Join(placeholders, ","))
						}
					case "not_in":
						if valueArr, ok := value.([]interface{}); ok && len(valueArr) > 0 {
							placeholders := make([]string, len(valueArr))
							for i, v := range valueArr {
								placeholders[i] = fmt.Sprintf("$%d", len(args)+1)
								args = append(args, v)
							}
							whereClause += fmt.Sprintf(" AND %s NOT IN (%s)", field, strings.Join(placeholders, ","))
						}
					}
				}
			}
		}
	}
	
	return whereClause, args
}

// buildUserSortQuery 构建用户排序查询条件
func (r *userRepo) buildUserSortQuery(sortConfig map[string]interface{}) string {
	orderClause := "ORDER BY created_at DESC" // 默认排序
	
	if sortConfig != nil {
		if field, ok := sortConfig["field"].(string); ok {
			direction := "ASC"
			if dir, ok := sortConfig["direction"].(string); ok && strings.ToUpper(dir) == "DESC" {
				direction = "DESC"
			}
			
			// 验证字段名以防止SQL注入
			allowedFields := []string{"id", "username", "email", "first_name", "last_name", "created_at", "updated_at", "last_login_time"}
			fieldAllowed := false
			for _, allowedField := range allowedFields {
				if field == allowedField {
					fieldAllowed = true
					break
				}
			}
			
			if fieldAllowed {
				orderClause = fmt.Sprintf("ORDER BY %s %s", field, direction)
			}
		}
	}
	
	return orderClause
}

// ListUsers 用户列表（原始版本，保持接口兼容性）
func (r *userRepo) ListUsers(ctx context.Context, page, size int32, search string) ([]*biz.User, int32, error) {
	options := &UserListOptions{
		Page:   page,
		Size:   size,
		Search: search,
	}
	return r.ListUsersWithOptions(ctx, options)
}

// ListUsersWithFilter 兼容接口的新方法
func (r *userRepo) ListUsersWithFilter(ctx context.Context, options interface{}) ([]*biz.User, int32, error) {
	if opts, ok := options.(*UserListOptions); ok {
		return r.ListUsersWithOptions(ctx, opts)
	}
	
	// 如果不是预期的类型，返回空结果
	return []*biz.User{}, 0, nil
}

// ValidatePassword 验证密码
func (r *userRepo) ValidatePassword(hashedPassword, password string) bool {
	return r.pm.ValidatePassword(hashedPassword, password)
}

// HashPassword 密码加密
func (r *userRepo) HashPassword(password string) (string, error) {
	return r.pm.HashPassword(password)
}

// UpdateLoginInfo 更新登录信息
func (r *userRepo) UpdateLoginInfo(ctx context.Context, userID int32, ip string) error {
	query := `
		UPDATE users 
		SET last_login_time = $1, last_login_ip = $2, login_count = login_count + 1
		WHERE id = $3`

	_, err := r.data.db.ExecContext(ctx, query, time.Now(), ip, userID)
	if err != nil {
		r.log.Errorf("failed to update login info: %v", err)
		return err
	}
	return nil
}

// GetUserRoles 获取用户角色
func (r *userRepo) GetUserRoles(ctx context.Context, userID int32) ([]*biz.Role, error) {
	query := `
		SELECT r.id, r.name, r.code, r.description, r.is_system_role, r.is_enabled, 
		       r.sort_order, r.created_at, r.updated_at
		FROM roles r
		INNER JOIN user_roles ur ON r.id = ur.role_id
		WHERE ur.user_id = $1 AND r.is_enabled = true
		ORDER BY r.sort_order`

	rows, err := r.data.db.QueryContext(ctx, query, userID)
	if err != nil {
		r.log.Errorf("failed to get user roles: %v", err)
		return nil, err
	}
	defer rows.Close()

	var roles []*biz.Role
	for rows.Next() {
		var role biz.Role
		err := rows.Scan(
			&role.ID, &role.Name, &role.Code, &role.Description, &role.IsSystemRole,
			&role.IsEnabled, &role.SortOrder, &role.CreatedAt, &role.UpdatedAt,
		)
		if err != nil {
			r.log.Errorf("failed to scan role: %v", err)
			return nil, err
		}
		roles = append(roles, &role)
	}

	return roles, nil
}

// GetUserPermissions 获取用户权限
func (r *userRepo) GetUserPermissions(ctx context.Context, userID int32) ([]string, error) {
	query := `
		SELECT DISTINCT p.code
		FROM permissions p
		INNER JOIN role_permissions rp ON p.id = rp.permission_id
		INNER JOIN user_roles ur ON rp.role_id = ur.role_id
		WHERE ur.user_id = $1`

	rows, err := r.data.db.QueryContext(ctx, query, userID)
	if err != nil {
		r.log.Errorf("failed to get user permissions: %v", err)
		return nil, err
	}
	defer rows.Close()

	var permissions []string
	for rows.Next() {
		var permission string
		err := rows.Scan(&permission)
		if err != nil {
			r.log.Errorf("failed to scan permission: %v", err)
			return nil, err
		}
		permissions = append(permissions, permission)
	}

	return permissions, nil
}

// AssignRoles 分配角色
func (r *userRepo) AssignRoles(ctx context.Context, userID int32, roleIDs []int32) error {
	r.log.Infof("AssignRoles called with userID: %d, roleIDs: %v", userID, roleIDs)
	tx, err := r.data.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// 删除现有角色
	_, err = tx.ExecContext(ctx, "DELETE FROM user_roles WHERE user_id = $1", userID)
	if err != nil {
		r.log.Errorf("failed to delete user roles: %v", err)
		return err
	}

	// 添加新角色
	for _, roleID := range roleIDs {
		_, err = tx.ExecContext(ctx, 
			"INSERT INTO user_roles (user_id, role_id, granted_at) VALUES ($1, $2, $3)",
			userID, roleID, time.Now())
		if err != nil {
			r.log.Errorf("failed to assign role: %v", err)
			return err
		}
	}

	return tx.Commit()
}

// EnableTwoFactor 启用2FA
func (r *userRepo) EnableTwoFactor(ctx context.Context, userID int32, secret string) error {
	query := `UPDATE users SET two_factor_enabled = true, two_factor_secret = $1 WHERE id = $2`
	_, err := r.data.db.ExecContext(ctx, query, secret, userID)
	if err != nil {
		r.log.Errorf("failed to enable two factor: %v", err)
		return err
	}
	return nil
}

// DisableTwoFactor 禁用2FA
func (r *userRepo) DisableTwoFactor(ctx context.Context, userID int32) error {
	query := `UPDATE users SET two_factor_enabled = false, two_factor_secret = '' WHERE id = $1`
	_, err := r.data.db.ExecContext(ctx, query, userID)
	if err != nil {
		r.log.Errorf("failed to disable two factor: %v", err)
		return err
	}
	return nil
}

// ValidateTwoFactor 验证2FA
func (r *userRepo) ValidateTwoFactor(ctx context.Context, userID int32, code string) bool {
	// TODO: 实现TOTP验证逻辑
	// 这里可以使用第三方库如 github.com/pquerna/otp
	r.log.Info("Two factor validation not implemented yet")
	return code == "123456" // 临时返回，用于测试
}