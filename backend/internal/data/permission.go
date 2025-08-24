package data

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"erp-system/internal/biz"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/lib/pq"
)

type permissionRepo struct {
	data *Data
	log  *log.Helper
}

func NewPermissionRepo(data *Data, logger log.Logger) biz.PermissionRepo {
	return &permissionRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

// DocType Operations
func (r *permissionRepo) CreateDocType(ctx context.Context, docType *biz.DocType) (*biz.DocType, error) {
	var id int64
	query := `
		INSERT INTO doc_types (name, label, module, is_submittable, is_child_table, permissions, naming_rule, 
		                      title_field, search_fields, sort_field, sort_order, description,
		                      created_at, updated_at, version)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)
		RETURNING id`

	permissionsJSON, _ := json.Marshal(docType.Permissions)
	searchFieldsJSON, _ := json.Marshal(docType.SearchFields)

	err := r.data.db.QueryRowContext(ctx, query,
		docType.Name, docType.Label, docType.Module, docType.IsSubmittable, docType.IsChildTable,
		permissionsJSON, docType.NamingRule, docType.TitleField, searchFieldsJSON,
		docType.SortField, docType.SortOrder, docType.Description,
		docType.CreatedAt, docType.UpdatedAt, docType.Version,
	).Scan(&id)

	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			if pqErr.Code == "23505" {
				return nil, fmt.Errorf("doctype already exists")
			}
		}
		r.log.Errorf("failed to create doctype: %v", err)
		return nil, err
	}

	docType.ID = id
	return docType, nil
}

func (r *permissionRepo) GetDocType(ctx context.Context, name string) (*biz.DocType, error) {
	var docType biz.DocType
	var permissionsJSON, searchFieldsJSON []byte
	var titleField, sortField, description sql.NullString

	query := `
		SELECT id, name, module, is_submittable, is_child_table, permissions, naming_rule,
		       title_field, search_fields, sort_field, sort_order, description,
		       created_at, updated_at, version
		FROM doc_types WHERE name = $1`

	err := r.data.db.QueryRowContext(ctx, query, name).Scan(
		&docType.ID, &docType.Name, &docType.Module, &docType.IsSubmittable,
		&docType.IsChildTable, &permissionsJSON, &docType.NamingRule,
		&titleField, &searchFieldsJSON, &sortField, &docType.SortOrder,
		&description, &docType.CreatedAt, &docType.UpdatedAt, &docType.Version,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("doctype not found")
		}
		r.log.Errorf("failed to get doctype: %v", err)
		return nil, err
	}

	// 解析JSON字段
	json.Unmarshal(permissionsJSON, &docType.Permissions)
	json.Unmarshal(searchFieldsJSON, &docType.SearchFields)

	if titleField.Valid {
		docType.TitleField = titleField.String
	}
	if sortField.Valid {
		docType.SortField = sortField.String
	}
	if description.Valid {
		docType.Description = description.String
	}

	return &docType, nil
}

func (r *permissionRepo) UpdateDocType(ctx context.Context, docType *biz.DocType) (*biz.DocType, error) {
	query := `
		UPDATE doc_types 
		SET module = $1, is_submittable = $2, is_child_table = $3, permissions = $4,
		    naming_rule = $5, title_field = $6, search_fields = $7, sort_field = $8,
		    sort_order = $9, description = $10, updated_at = $11, version = version + 1
		WHERE name = $12`

	permissionsJSON, _ := json.Marshal(docType.Permissions)
	searchFieldsJSON, _ := json.Marshal(docType.SearchFields)
	docType.UpdatedAt = time.Now()

	_, err := r.data.db.ExecContext(ctx, query,
		docType.Module, docType.IsSubmittable, docType.IsChildTable,
		permissionsJSON, docType.NamingRule, docType.TitleField,
		searchFieldsJSON, docType.SortField, docType.SortOrder,
		docType.Description, docType.UpdatedAt, docType.Name,
	)

	if err != nil {
		r.log.Errorf("failed to update doctype: %v", err)
		return nil, err
	}

	return docType, nil
}

func (r *permissionRepo) DeleteDocType(ctx context.Context, name string) error {
	// 检查是否有相关权限规则
	var count int
	err := r.data.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM permission_rules WHERE document_type = $1", name).Scan(&count)
	if err != nil {
		return err
	}
	if count > 0 {
		return fmt.Errorf("cannot delete doctype with existing permission rules")
	}

	// 删除DocType
	_, err = r.data.db.ExecContext(ctx, "DELETE FROM doc_types WHERE name = $1", name)
	if err != nil {
		r.log.Errorf("failed to delete doctype: %v", err)
		return err
	}

	return nil
}

func (r *permissionRepo) ListDocTypes(ctx context.Context, module string) ([]*biz.DocType, error) {
	query := `
		SELECT id, name, module, is_submittable, is_child_table, permissions, naming_rule,
		       title_field, search_fields, sort_field, sort_order, description,
		       created_at, updated_at, version
		FROM doc_types
		WHERE ($1 = '' OR module = $1)
		ORDER BY module, name`

	rows, err := r.data.db.QueryContext(ctx, query, module)
	if err != nil {
		r.log.Errorf("failed to list doctypes: %v", err)
		return nil, err
	}
	defer rows.Close()

	var docTypes []*biz.DocType
	for rows.Next() {
		var docType biz.DocType
		var permissionsJSON, searchFieldsJSON []byte
		var titleField, sortField, description sql.NullString

		err := rows.Scan(
			&docType.ID, &docType.Name, &docType.Module, &docType.IsSubmittable,
			&docType.IsChildTable, &permissionsJSON, &docType.NamingRule,
			&titleField, &searchFieldsJSON, &sortField, &docType.SortOrder,
			&description, &docType.CreatedAt, &docType.UpdatedAt, &docType.Version,
		)
		if err != nil {
			r.log.Errorf("failed to scan doctype: %v", err)
			return nil, err
		}

		// 解析JSON字段
		json.Unmarshal(permissionsJSON, &docType.Permissions)
		json.Unmarshal(searchFieldsJSON, &docType.SearchFields)

		if titleField.Valid {
			docType.TitleField = titleField.String
		}
		if sortField.Valid {
			docType.SortField = sortField.String
		}
		if description.Valid {
			docType.Description = description.String
		}

		docTypes = append(docTypes, &docType)
	}

	if err = rows.Err(); err != nil {
		r.log.Errorf("failed to iterate doctypes: %v", err)
		return nil, err
	}

	return docTypes, nil
}

// PermissionRule Operations
func (r *permissionRepo) CreatePermissionRule(ctx context.Context, rule *biz.PermissionRule) (*biz.PermissionRule, error) {
	var id int64
	query := `
		INSERT INTO permission_rules (role, document_type, permission_level, read, write, [create],
		                            [delete], submit, cancel, amend, report, export, import,
		                            set_user_permissions, if_owner, match, select_condition,
		                            delete_condition, amend_condition, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21)
		RETURNING id`

	err := r.data.db.QueryRowContext(ctx, query,
		rule.RoleID, rule.DocType, rule.PermissionLevel, rule.CanRead, rule.CanWrite,
		rule.CanCreate, rule.CanDelete, rule.CanSubmit, rule.CanCancel, rule.CanAmend,
		rule.CanReport, rule.CanExport, rule.CanImport, false,
		rule.OnlyIfCreator, "", "", "",
		"", rule.CreatedAt, rule.UpdatedAt,
	).Scan(&id)

	if err != nil {
		r.log.Errorf("failed to create permission rule: %v", err)
		return nil, err
	}

	rule.ID = id
	return rule, nil
}

func (r *permissionRepo) GetPermissionRule(ctx context.Context, id int64) (*biz.PermissionRule, error) {
	var rule biz.PermissionRule
	var selectCondition, deleteCondition, amendCondition sql.NullString

	query := `
		SELECT id, role, document_type, permission_level, read, write, [create],
		       [delete], submit, cancel, amend, report, export, import,
		       set_user_permissions, if_owner, match, select_condition,
		       delete_condition, amend_condition, created_at, updated_at
		FROM permission_rules WHERE id = $1`

	var roleID int64
	var ifOwner, setUserPermissions bool
	var match string
	err := r.data.db.QueryRowContext(ctx, query, id).Scan(
		&rule.ID, &roleID, &rule.DocType, &rule.PermissionLevel,
		&rule.CanRead, &rule.CanWrite, &rule.CanCreate, &rule.CanDelete,
		&rule.CanSubmit, &rule.CanCancel, &rule.CanAmend, &rule.CanReport,
		&rule.CanExport, &rule.CanImport, &setUserPermissions,
		&ifOwner, &match, &selectCondition,
		&deleteCondition, &amendCondition, &rule.CreatedAt, &rule.UpdatedAt,
	)
	rule.RoleID = roleID
	rule.OnlyIfCreator = ifOwner
	rule.CanSetUserPermissions = setUserPermissions

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("permission rule not found")
		}
		r.log.Errorf("failed to get permission rule: %v", err)
		return nil, err
	}

	// 注意：selectCondition, deleteCondition, amendCondition 字段
	// 在当前的PermissionRule结构体中不存在，这里暂时忽略

	return &rule, nil
}

func (r *permissionRepo) UpdatePermissionRule(ctx context.Context, rule *biz.PermissionRule) (*biz.PermissionRule, error) {
	query := `
		UPDATE permission_rules 
		SET role = $1, document_type = $2, permission_level = $3, read = $4, write = $5,
		    [create] = $6, [delete] = $7, submit = $8, cancel = $9, amend = $10,
		    report = $11, export = $12, import = $13, set_user_permissions = $14,
		    if_owner = $15, match = $16, select_condition = $17, delete_condition = $18,
		    amend_condition = $19, updated_at = $20
		WHERE id = $21`

	rule.UpdatedAt = time.Now()
	_, err := r.data.db.ExecContext(ctx, query,
		rule.RoleID, rule.DocType, rule.PermissionLevel, rule.CanRead, rule.CanWrite,
		rule.CanCreate, rule.CanDelete, rule.CanSubmit, rule.CanCancel, rule.CanAmend,
		rule.CanReport, rule.CanExport, rule.CanImport, false,
		rule.OnlyIfCreator, "", "", "",
		"", rule.UpdatedAt, rule.ID,
	)

	if err != nil {
		r.log.Errorf("failed to update permission rule: %v", err)
		return nil, err
	}

	return rule, nil
}

func (r *permissionRepo) DeletePermissionRule(ctx context.Context, id int64) error {
	_, err := r.data.db.ExecContext(ctx, "DELETE FROM permission_rules WHERE id = $1", id)
	if err != nil {
		r.log.Errorf("failed to delete permission rule: %v", err)
		return err
	}
	return nil
}

func (r *permissionRepo) ListPermissionRules(ctx context.Context, roleID int64, docType string) ([]*biz.PermissionRule, error) {
	query := `
		SELECT id, role, document_type, permission_level, read, write, [create],
		       [delete], submit, cancel, amend, report, export, import,
		       set_user_permissions, if_owner, match, select_condition,
		       delete_condition, amend_condition, created_at, updated_at
		FROM permission_rules
		WHERE ($1 = 0 OR role = $1) AND ($2 = '' OR document_type = $2)
		ORDER BY document_type, permission_level, role`

	rows, err := r.data.db.QueryContext(ctx, query, roleID, docType)
	if err != nil {
		r.log.Errorf("failed to list permission rules: %v", err)
		return nil, err
	}
	defer rows.Close()

	var rules []*biz.PermissionRule
	for rows.Next() {
		var rule biz.PermissionRule
		var selectCondition, deleteCondition, amendCondition sql.NullString

		var roleID int64
		var ifOwner, setUserPermissions bool
		var match string
		err := rows.Scan(
			&rule.ID, &roleID, &rule.DocType, &rule.PermissionLevel,
			&rule.CanRead, &rule.CanWrite, &rule.CanCreate, &rule.CanDelete,
			&rule.CanSubmit, &rule.CanCancel, &rule.CanAmend, &rule.CanReport,
			&rule.CanExport, &rule.CanImport, &setUserPermissions,
			&ifOwner, &match, &selectCondition,
			&deleteCondition, &amendCondition, &rule.CreatedAt, &rule.UpdatedAt,
		)
		if err != nil {
			r.log.Errorf("failed to scan permission rule: %v", err)
			return nil, err
		}

		rule.RoleID = roleID
		rule.OnlyIfCreator = ifOwner
		rule.CanSetUserPermissions = setUserPermissions
		// 忽略不存在的字段selectCondition, deleteCondition, amendCondition

		rules = append(rules, &rule)
	}

	if err = rows.Err(); err != nil {
		r.log.Errorf("failed to iterate permission rules: %v", err)
		return nil, err
	}

	return rules, nil
}

func (r *permissionRepo) BatchCreatePermissionRules(ctx context.Context, rules []*biz.PermissionRule) error {
	if len(rules) == 0 {
		return nil
	}

	// 使用事务处理批量插入
	tx, err := r.data.db.BeginTx(ctx, nil)
	if err != nil {
		r.log.Errorf("failed to begin transaction: %v", err)
		return err
	}
	defer tx.Rollback()

	query := `
		INSERT INTO permission_rules (role, document_type, permission_level, read, write, create,
		                            delete, submit, cancel, amend, report, export, import,
		                            set_user_permissions, if_owner, match, select_condition,
		                            delete_condition, amend_condition, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21)`

	for _, rule := range rules {
		_, err = tx.ExecContext(ctx, query,
			rule.RoleID, rule.DocType, rule.PermissionLevel, rule.CanRead, rule.CanWrite,
			rule.CanCreate, rule.CanDelete, rule.CanSubmit, rule.CanCancel, rule.CanAmend,
			rule.CanReport, rule.CanExport, rule.CanImport, false,
			rule.OnlyIfCreator, "", "", "",
			"", rule.CreatedAt, rule.UpdatedAt,
		)
		if err != nil {
			r.log.Errorf("failed to batch create permission rule: %v", err)
			return err
		}
	}

	if err = tx.Commit(); err != nil {
		r.log.Errorf("failed to commit transaction: %v", err)
		return err
	}

	return nil
}

func (r *permissionRepo) BatchCreateUserPermissions(ctx context.Context, userPerms []*biz.UserPermission) error {
	if len(userPerms) == 0 {
		return nil
	}

	// 使用事务处理批量插入
	tx, err := r.data.db.BeginTx(ctx, nil)
	if err != nil {
		r.log.Errorf("failed to begin transaction: %v", err)
		return err
	}
	defer tx.Rollback()

	query := `
		INSERT INTO user_permissions (user, allow, for_value, document_type, is_default,
		                            created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)`

	for _, userPerm := range userPerms {
		_, err = tx.ExecContext(ctx, query,
			userPerm.UserID, userPerm.Value, userPerm.DocName,
			userPerm.DocType, userPerm.IsDefault,
			userPerm.CreatedAt, userPerm.UpdatedAt,
		)
		if err != nil {
			r.log.Errorf("failed to batch create user permission: %v", err)
			return err
		}
	}

	if err = tx.Commit(); err != nil {
		r.log.Errorf("failed to commit transaction: %v", err)
		return err
	}

	return nil
}

// UserPermission Operations
func (r *permissionRepo) CreateUserPermission(ctx context.Context, userPerm *biz.UserPermission) (*biz.UserPermission, error) {
	var id int64
	query := `
		INSERT INTO user_permissions (user, allow, for_value, document_type, is_default,
		                            created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id`

	err := r.data.db.QueryRowContext(ctx, query,
		userPerm.UserID, userPerm.Value, userPerm.DocName,
		userPerm.DocType, userPerm.IsDefault,
		userPerm.CreatedAt, userPerm.UpdatedAt,
	).Scan(&id)

	if err != nil {
		r.log.Errorf("failed to create user permission: %v", err)
		return nil, err
	}

	userPerm.ID = id
	return userPerm, nil
}

func (r *permissionRepo) GetUserPermission(ctx context.Context, id int64) (*biz.UserPermission, error) {
	var userPerm biz.UserPermission

	query := `
		SELECT id, user, allow, for_value, document_type, is_default,
		       created_at, updated_at
		FROM user_permissions WHERE id = $1`

	err := r.data.db.QueryRowContext(ctx, query, id).Scan(
		&userPerm.ID, &userPerm.UserID, &userPerm.Value, &userPerm.DocName,
		&userPerm.DocType, &userPerm.IsDefault,
		&userPerm.CreatedAt, &userPerm.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user permission not found")
		}
		r.log.Errorf("failed to get user permission: %v", err)
		return nil, err
	}

	return &userPerm, nil
}

func (r *permissionRepo) UpdateUserPermission(ctx context.Context, userPerm *biz.UserPermission) (*biz.UserPermission, error) {
	query := `
		UPDATE user_permissions 
		SET user = $1, allow = $2, for_value = $3, document_type = $4,
		    is_default = $5, updated_at = $6
		WHERE id = $7`

	userPerm.UpdatedAt = time.Now()
	_, err := r.data.db.ExecContext(ctx, query,
		userPerm.UserID, userPerm.Value, userPerm.DocName,
		userPerm.DocType, userPerm.IsDefault,
		userPerm.UpdatedAt, userPerm.ID,
	)

	if err != nil {
		r.log.Errorf("failed to update user permission: %v", err)
		return nil, err
	}

	return userPerm, nil
}

func (r *permissionRepo) DeleteUserPermission(ctx context.Context, id int64) error {
	_, err := r.data.db.ExecContext(ctx, "DELETE FROM user_permissions WHERE id = $1", id)
	if err != nil {
		r.log.Errorf("failed to delete user permission: %v", err)
		return err
	}
	return nil
}

func (r *permissionRepo) ListUserPermissions(ctx context.Context, userID int64, docType string, page, size int32) ([]*biz.UserPermission, error) {
	query := `
		SELECT id, user, allow, for_value, document_type, is_default,
		       created_at, updated_at
		FROM user_permissions
		WHERE ($1 = 0 OR user = $1) AND ($2 = '' OR document_type = $2)
		ORDER BY document_type, user, allow
		LIMIT $3 OFFSET $4`

	offset := (page - 1) * size
	rows, err := r.data.db.QueryContext(ctx, query, userID, docType, size, offset)
	if err != nil {
		r.log.Errorf("failed to list user permissions: %v", err)
		return nil, err
	}
	defer rows.Close()

	var userPerms []*biz.UserPermission
	for rows.Next() {
		var userPerm biz.UserPermission

		err := rows.Scan(
			&userPerm.ID, &userPerm.UserID, &userPerm.Value, &userPerm.DocName,
			&userPerm.DocType, &userPerm.IsDefault,
			&userPerm.CreatedAt, &userPerm.UpdatedAt,
		)
		if err != nil {
			r.log.Errorf("failed to scan user permission: %v", err)
			return nil, err
		}

		userPerms = append(userPerms, &userPerm)
	}

	if err = rows.Err(); err != nil {
		r.log.Errorf("failed to iterate user permissions: %v", err)
		return nil, err
	}

	return userPerms, nil
}

func (r *permissionRepo) GetUserPermissionsCount(ctx context.Context, userID int64, docType string) (int32, error) {
	query := `
		SELECT COUNT(*)
		FROM user_permissions
		WHERE ($1 = 0 OR user = $1) AND ($2 = '' OR document_type = $2)`

	var count int32
	err := r.data.db.QueryRowContext(ctx, query, userID, docType).Scan(&count)
	if err != nil {
		r.log.Errorf("failed to count user permissions: %v", err)
		return 0, err
	}

	return count, nil
}

// FieldPermissionLevel Operations
func (r *permissionRepo) CreateFieldPermissionLevel(ctx context.Context, fieldPerm *biz.FieldPermissionLevel) (*biz.FieldPermissionLevel, error) {
	var id int64
	query := `
		INSERT INTO field_permission_levels (document_type, field_name, permission_level,
		                                   read, write, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id`

	err := r.data.db.QueryRowContext(ctx, query,
		fieldPerm.DocType, fieldPerm.FieldName, fieldPerm.PermissionLevel,
		false, false, // 暂时设置为false，因为当前结构体没有Read/Write字段
		fieldPerm.CreatedAt, fieldPerm.UpdatedAt,
	).Scan(&id)

	if err != nil {
		r.log.Errorf("failed to create field permission level: %v", err)
		return nil, err
	}

	fieldPerm.ID = id
	return fieldPerm, nil
}

func (r *permissionRepo) GetFieldPermissionLevel(ctx context.Context, id int64) (*biz.FieldPermissionLevel, error) {
	var fieldPerm biz.FieldPermissionLevel

	query := `
		SELECT id, document_type, field_name, permission_level, read, write,
		       created_at, updated_at
		FROM field_permission_levels WHERE id = $1`

	var read, write bool
	err := r.data.db.QueryRowContext(ctx, query, id).Scan(
		&fieldPerm.ID, &fieldPerm.DocType, &fieldPerm.FieldName,
		&fieldPerm.PermissionLevel, &read, &write,
		&fieldPerm.CreatedAt, &fieldPerm.UpdatedAt,
	)
	// 忽略read/write字段，因为当前结构体没有这些字段

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("field permission level not found")
		}
		r.log.Errorf("failed to get field permission level: %v", err)
		return nil, err
	}

	return &fieldPerm, nil
}

func (r *permissionRepo) ListFieldPermissionLevels(ctx context.Context, docType string, page, size int32) ([]*biz.FieldPermissionLevel, error) {
	query := `
		SELECT id, document_type, field_name, permission_level, read, write,
		       created_at, updated_at
		FROM field_permission_levels
		WHERE ($1 = '' OR document_type = $1)
		ORDER BY document_type, field_name, permission_level
		LIMIT $2 OFFSET $3`

	offset := (page - 1) * size
	rows, err := r.data.db.QueryContext(ctx, query, docType, size, offset)
	if err != nil {
		r.log.Errorf("failed to list field permission levels: %v", err)
		return nil, err
	}
	defer rows.Close()

	var fieldPerms []*biz.FieldPermissionLevel
	for rows.Next() {
		var fieldPerm biz.FieldPermissionLevel

		var read, write bool
		err := rows.Scan(
			&fieldPerm.ID, &fieldPerm.DocType, &fieldPerm.FieldName,
			&fieldPerm.PermissionLevel, &read, &write,
			&fieldPerm.CreatedAt, &fieldPerm.UpdatedAt,
		)
		// 忽略read/write字段，因为当前结构体没有这些字段
		if err != nil {
			r.log.Errorf("failed to scan field permission level: %v", err)
			return nil, err
		}

		fieldPerms = append(fieldPerms, &fieldPerm)
	}

	if err = rows.Err(); err != nil {
		r.log.Errorf("failed to iterate field permission levels: %v", err)
		return nil, err
	}

	return fieldPerms, nil
}

func (r *permissionRepo) DeleteFieldPermissionLevel(ctx context.Context, id int64) error {
	_, err := r.data.db.ExecContext(ctx, "DELETE FROM field_permission_levels WHERE id = $1", id)
	if err != nil {
		r.log.Errorf("failed to delete field permission level: %v", err)
		return err
	}
	return nil
}

func (r *permissionRepo) UpdateFieldPermissionLevel(ctx context.Context, fieldPerm *biz.FieldPermissionLevel) (*biz.FieldPermissionLevel, error) {
	query := `
		UPDATE field_permission_levels 
		SET document_type = $1, field_name = $2, permission_level = $3,
		    read = $4, write = $5, updated_at = $6
		WHERE id = $7`

	fieldPerm.UpdatedAt = time.Now()
	_, err := r.data.db.ExecContext(ctx, query,
		fieldPerm.DocType, fieldPerm.FieldName, fieldPerm.PermissionLevel,
		false, false, fieldPerm.UpdatedAt, fieldPerm.ID,
	)

	if err != nil {
		r.log.Errorf("failed to update field permission level: %v", err)
		return nil, err
	}

	return fieldPerm, nil
}

// DocumentWorkflowState Operations
func (r *permissionRepo) CreateDocumentWorkflowState(ctx context.Context, state *biz.DocumentWorkflowState) (*biz.DocumentWorkflowState, error) {
	var id int64
	query := `
		INSERT INTO document_workflow_states (document_type, document_name, workflow_state,
		                                    workflow_action, created_by, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id`

	var docName string
	if state.DocName != nil {
		docName = *state.DocName
	} else {
		docName = fmt.Sprintf("%d", state.DocID)
	}

	err := r.data.db.QueryRowContext(ctx, query,
		state.DocType, docName, state.WorkflowState,
		"", state.SubmittedBy, state.CreatedAt, // 使用空字符串作为workflow_action，使用SubmittedBy作为created_by
	).Scan(&id)

	if err != nil {
		r.log.Errorf("failed to create document workflow state: %v", err)
		return nil, err
	}

	state.ID = id
	return state, nil
}

func (r *permissionRepo) GetDocumentWorkflowState(ctx context.Context, stateID int64) (*biz.DocumentWorkflowState, error) {
	var state biz.DocumentWorkflowState
	var workflowAction sql.NullString

	query := `
		SELECT id, document_type, document_name, workflow_state, workflow_action,
		       created_by, created_at
		FROM document_workflow_states 
		WHERE id = $1`

	var docName string
	var createdBy sql.NullInt64
	err := r.data.db.QueryRowContext(ctx, query, stateID).Scan(
		&state.ID, &state.DocType, &docName,
		&state.WorkflowState, &workflowAction,
		&createdBy, &state.CreatedAt,
	)

	if docName != "" {
		state.DocName = &docName
	}
	if createdBy.Valid {
		state.SubmittedBy = &createdBy.Int64
	}

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("document workflow state not found")
		}
		r.log.Errorf("failed to get document workflow state: %v", err)
		return nil, err
	}

	// WorkflowAction字段在当前结构体中不存在，暂时忽略
	_ = workflowAction

	return &state, nil
}

func (r *permissionRepo) DeleteDocumentWorkflowState(ctx context.Context, stateID int64) error {
	_, err := r.data.db.ExecContext(ctx,
		"DELETE FROM document_workflow_states WHERE id = $1",
		stateID)
	if err != nil {
		r.log.Errorf("failed to delete document workflow state: %v", err)
		return err
	}
	return nil
}

func (r *permissionRepo) UpdateDocumentWorkflowState(ctx context.Context, state *biz.DocumentWorkflowState) (*biz.DocumentWorkflowState, error) {
	// 对于工作流状态，我们通常是插入新记录而不是更新，因为需要保留历史记录
	return r.CreateDocumentWorkflowState(ctx, state)
}

func (r *permissionRepo) ListDocumentWorkflowStates(ctx context.Context, docType, documentName, state string, userID int64, page, size int32) ([]*biz.DocumentWorkflowState, error) {
	query := `
		SELECT id, document_type, document_name, workflow_state,
		       created_by, created_at
		FROM document_workflow_states 
		WHERE ($1 = '' OR document_type = $1)
		  AND ($2 = '' OR document_name = $2)
		  AND ($3 = '' OR workflow_state = $3)
		  AND ($4 = 0 OR created_by = $4)
		ORDER BY created_at DESC
		LIMIT $5 OFFSET $6`

	offset := (page - 1) * size
	rows, err := r.data.db.QueryContext(ctx, query, docType, documentName, state, userID, size, offset)
	if err != nil {
		r.log.Errorf("failed to list document workflow states: %v", err)
		return nil, err
	}
	defer rows.Close()

	var states []*biz.DocumentWorkflowState
	for rows.Next() {
		var workflowState biz.DocumentWorkflowState
		var docName string
		var createdBy sql.NullInt64

		err := rows.Scan(
			&workflowState.ID, &workflowState.DocType, &docName,
			&workflowState.WorkflowState, &createdBy, &workflowState.CreatedAt,
		)
		if err != nil {
			r.log.Errorf("failed to scan document workflow state: %v", err)
			return nil, err
		}

		if docName != "" {
			workflowState.DocName = &docName
		}
		if createdBy.Valid {
			workflowState.SubmittedBy = &createdBy.Int64
		}
		workflowState.UpdatedAt = workflowState.CreatedAt

		states = append(states, &workflowState)
	}

	if err = rows.Err(); err != nil {
		r.log.Errorf("failed to iterate document workflow states: %v", err)
		return nil, err
	}

	return states, nil
}

// Permission Checking Operations
func (r *permissionRepo) CheckPermission(ctx context.Context, userID int64, documentType, action string, permissionLevel int) (bool, error) {
	query := `
		SELECT CASE WHEN COUNT(*) > 0 THEN 1 ELSE 0 END
		FROM user_roles ur
		INNER JOIN permission_rules pr ON ur.role_id = pr.role
		WHERE ur.user_id = $1 
		  AND pr.document_type = $2
		  AND pr.permission_level = $4
		  AND (
		    ($3 = 'read' AND pr.read = 1) OR
		    ($3 = 'write' AND pr.write = 1) OR
		    ($3 = 'create' AND pr.[create] = 1) OR
		    ($3 = 'delete' AND pr.[delete] = 1) OR
		    ($3 = 'submit' AND pr.submit = 1) OR
		    ($3 = 'cancel' AND pr.cancel = 1) OR
		    ($3 = 'amend' AND pr.amend = 1) OR
		    ($3 = 'print' AND pr.print = 1) OR
		    ($3 = 'email' AND pr.email = 1) OR
		    ($3 = 'import' AND pr.import = 1) OR
		    ($3 = 'export' AND pr.export = 1) OR
		    ($3 = 'share' AND pr.share = 1) OR
		    ($3 = 'report' AND pr.report = 1)
		  )`

	var hasPermission bool
	err := r.data.db.QueryRowContext(ctx, query, userID, documentType, action, permissionLevel).Scan(&hasPermission)
	if err != nil {
		r.log.Errorf("failed to check permission: %v", err)
		return false, err
	}

	return hasPermission, nil
}

func (r *permissionRepo) GetDocumentWorkflowStatesCount(ctx context.Context, docType, documentName, state string, userID int64) (int32, error) {
	query := `
		SELECT COUNT(*) 
		FROM document_workflow_states dws
		WHERE dws.document_type = $1 
		AND ($2 = '' OR dws.document_name = $2)
		AND ($3 = '' OR dws.workflow_state = $3)
		AND ($4 = 0 OR dws.created_by = $4)`

	var count int32
	err := r.data.db.QueryRowContext(ctx, query, docType, documentName, state, userID).Scan(&count)
	if err != nil {
		r.log.Errorf("failed to count document workflow states: %v", err)
		return 0, err
	}

	return count, nil
}

func (r *permissionRepo) GetFieldPermissionLevelsCount(ctx context.Context, docType string) (int32, error) {
	query := `
		SELECT COUNT(*) 
		FROM field_permission_levels 
		WHERE doc_type = $1`

	var count int32
	err := r.data.db.QueryRowContext(ctx, query, docType).Scan(&count)
	if err != nil {
		r.log.Errorf("failed to count field permission levels: %v", err)
		return 0, err
	}

	return count, nil
}

func (r *permissionRepo) GetUserPermissionLevel(ctx context.Context, userID int64, documentType string) (int, error) {
	query := `
		SELECT COALESCE(MIN(pr.permission_level), 0)
		FROM user_roles ur
		INNER JOIN permission_rules pr ON ur.role_id = pr.role
		WHERE ur.user_id = $1 
		  AND pr.document_type = $2
		  AND (pr.read = 1 OR pr.write = 1)`

	var level int
	err := r.data.db.QueryRowContext(ctx, query, userID, documentType).Scan(&level)
	if err != nil {
		r.log.Errorf("failed to get user permission level: %v", err)
		return 0, err
	}

	return level, nil
}

func (r *permissionRepo) FilterDocumentsByPermission(ctx context.Context, userID int64, documentType string, documents []map[string]interface{}) ([]map[string]interface{}, error) {
	// 获取用户权限级别
	permissionLevel, err := r.GetUserPermissionLevel(ctx, userID, documentType)
	if err != nil {
		return nil, err
	}

	// 获取字段权限配置 (使用较大的limit来获取所有字段权限)
	fieldPermissions, err := r.ListFieldPermissionLevels(ctx, documentType, 1, 1000)
	if err != nil {
		return nil, err
	}

	// 过滤文档字段
	var filteredDocuments []map[string]interface{}
	for _, doc := range documents {
		filteredDoc := make(map[string]interface{})
		for field, value := range doc {
			canRead := true

			// 检查字段权限
			for _, fp := range fieldPermissions {
				if fp.FieldName == field && fp.PermissionLevel > permissionLevel {
					canRead = false // 当前结构体没有Read字段，默认为false
					break
				}
			}

			if canRead {
				filteredDoc[field] = value
			}
		}
		filteredDocuments = append(filteredDocuments, filteredDoc)
	}

	return filteredDocuments, nil
}

func (r *permissionRepo) GetUserRoles(ctx context.Context, userID int64) ([]string, error) {
	query := `
		SELECT role_name 
		FROM user_roles ur 
		JOIN roles r ON ur.role_id = r.id 
		WHERE ur.user_id = $1`

	rows, err := r.data.db.QueryContext(ctx, query, userID)
	if err != nil {
		r.log.Errorf("failed to get user roles: %v", err)
		return nil, err
	}
	defer rows.Close()

	var roles []string
	for rows.Next() {
		var role string
		if err := rows.Scan(&role); err != nil {
			r.log.Errorf("failed to scan user role: %v", err)
			return nil, err
		}
		roles = append(roles, role)
	}

	if err = rows.Err(); err != nil {
		r.log.Errorf("failed to iterate user roles: %v", err)
		return nil, err
	}

	return roles, nil
}

// CheckDocumentPermission 检查文档权限
func (r *permissionRepo) CheckDocumentPermission(ctx context.Context, req *biz.PermissionCheckRequest) (bool, error) {
	// 根据请求检查权限
	return r.CheckPermission(ctx, req.UserID, req.DocType, req.Permission, 0)
}

// GetAccessibleFields 获取可访问字段
func (r *permissionRepo) GetAccessibleFields(ctx context.Context, req *biz.FieldPermissionRequest) ([]*biz.AccessibleField, error) {
	// 获取用户权限级别
	userLevel, err := r.GetUserPermissionLevel(ctx, req.UserID, req.DocType)
	if err != nil {
		return nil, err
	}

	// 获取字段权限配置 (使用较大的limit来获取所有字段权限)
	fieldPermissions, err := r.ListFieldPermissionLevels(ctx, req.DocType, 1, 1000)
	if err != nil {
		return nil, err
	}

	var accessibleFields []*biz.AccessibleField
	for _, fp := range fieldPermissions {
		accessible := &biz.AccessibleField{
			FieldName:       fp.FieldName,
			CanAccess:       userLevel >= fp.PermissionLevel,
			PermissionLevel: fp.PermissionLevel,
		}
		accessibleFields = append(accessibleFields, accessible)
	}

	return accessibleFields, nil
}

// GetUserEnhancedPermissions 获取用户增强权限
func (r *permissionRepo) GetUserEnhancedPermissions(ctx context.Context, userID int64, docType string) ([]*biz.EnhancedUserPermission, error) {
	// 获取用户权限 (使用较大的limit来获取所有用户权限)
	userPerms, err := r.ListUserPermissions(ctx, userID, docType, 1, 1000)
	if err != nil {
		return nil, err
	}

	var enhancedPerms []*biz.EnhancedUserPermission
	for _, perm := range userPerms {
		enhanced := &biz.EnhancedUserPermission{
			UserID:          perm.UserID,
			DocType:         perm.DocType,
			PermissionLevel: 0,  // 默认权限级别
			RoleCode:        "", // 需要从用户角色中获取
			RoleName:        "", // 需要从用户角色中获取

			// 权限列表，这里设置为默认值，实际应该从权限规则中计算
			CanRead:       true,
			CanWrite:      false,
			CanCreate:     false,
			CanDelete:     false,
			CanSubmit:     false,
			CanCancel:     false,
			CanAmend:      false,
			CanPrint:      false,
			CanEmail:      false,
			CanImport:     false,
			CanExport:     false,
			CanShare:      false,
			CanReport:     false,
			OnlyIfCreator: false,
		}
		enhancedPerms = append(enhancedPerms, enhanced)
	}

	return enhancedPerms, nil
}
