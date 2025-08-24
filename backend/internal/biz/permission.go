package biz

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/go-kratos/kratos/v2/log"
)

// ================================================================
// Frappe式权限系统数据模型
// ================================================================

// DocType 文档类型
type DocType struct {
	ID             int64                  `json:"id"`
	Name           string                 `json:"name"`                    // 文档类型名称
	Label          string                 `json:"label"`                   // 显示名称
	Module         string                 `json:"module"`                  // 所属模块
	Description    string                 `json:"description,omitempty"`   // 描述
	IsSubmittable  bool                   `json:"is_submittable"`          // 是否支持提交工作流
	IsChildTable   bool                   `json:"is_child_table"`          // 是否是子表
	HasWorkflow    bool                   `json:"has_workflow"`            // 是否有工作流
	TrackChanges   bool                   `json:"track_changes"`           // 是否跟踪变更
	AppliesToAll   bool                   `json:"applies_to_all_users"`    // 是否对所有用户可见
	MaxAttachments int                    `json:"max_attachments"`         // 最大附件数
	Permissions    map[string]interface{} `json:"permissions,omitempty"`   // 权限设置JSON
	NamingRule     string                 `json:"naming_rule,omitempty"`   // 命名规则
	TitleField     string                 `json:"title_field,omitempty"`   // 标题字段
	SearchFields   []string               `json:"search_fields,omitempty"` // 搜索字段
	SortField      string                 `json:"sort_field,omitempty"`    // 排序字段
	SortOrder      string                 `json:"sort_order,omitempty"`    // 排序方向
	Version        int                    `json:"version"`                 // 版本号
	CreatedAt      time.Time              `json:"created_at"`
	UpdatedAt      time.Time              `json:"updated_at"`
	CreatedBy      *int64                 `json:"created_by,omitempty"`
	UpdatedBy      *int64                 `json:"updated_by,omitempty"`
}

// Validate 验证DocType数据的完整性和正确性
func (d *DocType) Validate() error {
	if strings.TrimSpace(d.Name) == "" {
		return fmt.Errorf("DocType name is required")
	}

	if strings.TrimSpace(d.Label) == "" {
		return fmt.Errorf("DocType label is required")
	}

	if strings.TrimSpace(d.Module) == "" {
		return fmt.Errorf("DocType module is required")
	}

	if d.MaxAttachments < 0 {
		return fmt.Errorf("MaxAttachments cannot be negative")
	}

	// 验证命名规则
	if d.NamingRule != "" {
		validNamingRules := []string{"autoname", "prompt", "field", "series", "random", "expression"}
		isValid := false
		for _, rule := range validNamingRules {
			if d.NamingRule == rule {
				isValid = true
				break
			}
		}
		if !isValid {
			return fmt.Errorf("invalid naming rule: %s", d.NamingRule)
		}
	}

	// 验证排序方向
	if d.SortOrder != "" && d.SortOrder != "ASC" && d.SortOrder != "DESC" {
		return fmt.Errorf("sort order must be ASC or DESC")
	}

	return nil
}

// PermissionRule 权限规则
type PermissionRule struct {
	ID              int64  `json:"id"`
	RoleID          int64  `json:"role_id"`          // 角色ID
	DocType         string `json:"doc_type"`         // 文档类型
	PermissionLevel int    `json:"permission_level"` // 权限级别 (0=文档级, 1-9=字段级)

	// 基础权限
	CanRead  bool `json:"can_read"`  // 读取权限
	CanWrite bool `json:"can_write"` // 写入权限

	// 文档级权限 (仅permission_level=0时有效)
	CanCreate             bool `json:"can_create"`               // 创建权限
	CanDelete             bool `json:"can_delete"`               // 删除权限
	CanSubmit             bool `json:"can_submit"`               // 提交权限
	CanCancel             bool `json:"can_cancel"`               // 取消权限
	CanAmend              bool `json:"can_amend"`                // 修订权限
	CanPrint              bool `json:"can_print"`                // 打印权限
	CanEmail              bool `json:"can_email"`                // 邮件权限
	CanImport             bool `json:"can_import"`               // 导入权限
	CanExport             bool `json:"can_export"`               // 导出权限
	CanShare              bool `json:"can_share"`                // 分享权限
	CanReport             bool `json:"can_report"`               // 报表权限
	CanSetUserPermissions bool `json:"can_set_user_permissions"` // 设置用户权限

	// 条件权限
	OnlyIfCreator bool `json:"only_if_creator"` // 仅创建者可访问

	// 审计字段
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	CreatedBy *int64    `json:"created_by,omitempty"`
	UpdatedBy *int64    `json:"updated_by,omitempty"`

	// 关联对象
	Role *Role `json:"role,omitempty"`
}

// Validate 验证PermissionRule数据的完整性和正确性
func (r *PermissionRule) Validate() error {
	if r.RoleID <= 0 {
		return fmt.Errorf("RoleID must be positive")
	}

	if strings.TrimSpace(r.DocType) == "" {
		return fmt.Errorf("DocType is required")
	}

	if r.PermissionLevel < 0 || r.PermissionLevel > 10 {
		return fmt.Errorf("PermissionLevel must be between 0 and 10")
	}

	// 文档级权限验证（permission_level = 0）
	if r.PermissionLevel == 0 {
		// 至少需要一个基础权限
		if !r.CanRead && !r.CanWrite {
			return fmt.Errorf("at least one of CanRead or CanWrite must be true for document level permissions")
		}
	}

	// 字段级权限验证（permission_level > 0）
	if r.PermissionLevel > 0 {
		// 字段级权限不应该有文档级的权限设置
		if r.CanCreate || r.CanDelete || r.CanSubmit || r.CanCancel || r.CanAmend ||
			r.CanPrint || r.CanEmail || r.CanImport || r.CanExport || r.CanShare || r.CanReport {
			return fmt.Errorf("field level permissions (level > 0) should only use CanRead and CanWrite")
		}
	}

	return nil
}

// UserPermission 用户权限 - 数据范围权限
type UserPermission struct {
	ID              int64     `json:"id"`
	UserID          int64     `json:"user_id"`                  // 用户ID
	DocType         string    `json:"doc_type"`                 // 受限制的文档类型
	DocName         string    `json:"doc_name"`                 // 允许访问的具体记录标识
	Value           string    `json:"value"`                    // 权限值
	ApplicableFor   *string   `json:"applicable_for,omitempty"` // 仅对指定文档类型生效
	HideDescendants bool      `json:"hide_descendants"`         // 隐藏子记录
	IsDefault       bool      `json:"is_default"`               // 是否为默认权限
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
	CreatedBy       *int64    `json:"created_by,omitempty"`
	UpdatedBy       *int64    `json:"updated_by,omitempty"`

	// 关联对象
	User *User `json:"user,omitempty"`
}

// Validate 验证UserPermission数据的完整性和正确性
func (u *UserPermission) Validate() error {
	if u.UserID <= 0 {
		return fmt.Errorf("UserID must be positive")
	}

	if strings.TrimSpace(u.DocType) == "" {
		return fmt.Errorf("DocType is required")
	}

	if strings.TrimSpace(u.DocName) == "" {
		return fmt.Errorf("DocName is required")
	}

	if strings.TrimSpace(u.Value) == "" {
		return fmt.Errorf("Value is required")
	}

	// 验证ApplicableFor如果不为空，必须是有效的DocType
	if u.ApplicableFor != nil && strings.TrimSpace(*u.ApplicableFor) == "" {
		return fmt.Errorf("ApplicableFor cannot be empty string if provided")
	}

	return nil
}

// FieldPermissionLevel 字段权限级别
type FieldPermissionLevel struct {
	ID              int64     `json:"id"`
	DocType         string    `json:"doc_type"`              // 文档类型
	FieldName       string    `json:"field_name"`            // 字段名称
	FieldLabel      string    `json:"field_label,omitempty"` // 字段显示名称
	PermissionLevel int       `json:"permission_level"`      // 权限级别
	FieldType       string    `json:"field_type"`            // 字段类型
	IsMandatory     bool      `json:"is_mandatory"`          // 是否必填
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
	CreatedBy       *int64    `json:"created_by,omitempty"`
	UpdatedBy       *int64    `json:"updated_by,omitempty"`
}

// Validate 验证FieldPermissionLevel数据的完整性和正确性
func (f *FieldPermissionLevel) Validate() error {
	if strings.TrimSpace(f.DocType) == "" {
		return fmt.Errorf("DocType is required")
	}

	if strings.TrimSpace(f.FieldName) == "" {
		return fmt.Errorf("FieldName is required")
	}

	if f.PermissionLevel < 0 || f.PermissionLevel > 10 {
		return fmt.Errorf("PermissionLevel must be between 0 and 10")
	}

	// 验证字段类型
	if strings.TrimSpace(f.FieldType) != "" {
		validTypes := []string{
			"Data", "Text", "Long Text", "HTML", "Markdown",
			"Int", "Float", "Currency", "Percent",
			"Date", "Datetime", "Time",
			"Select", "Link", "Dynamic Link", "Table", "Check",
			"Small Text", "Text Editor", "Code", "Password",
			"Attach", "Attach Image", "Signature", "Color",
			"Barcode", "Geolocation", "Duration", "Rating",
		}

		isValidType := false
		for _, validType := range validTypes {
			if f.FieldType == validType {
				isValidType = true
				break
			}
		}
		if !isValidType {
			return fmt.Errorf("invalid field type: %s", f.FieldType)
		}
	}

	return nil
}

// DocumentWorkflowState 文档工作流状态
type DocumentWorkflowState struct {
	ID            int64   `json:"id"`
	DocType       string  `json:"doc_type"`           // 文档类型
	DocID         int64   `json:"doc_id"`             // 文档ID
	DocName       *string `json:"doc_name,omitempty"` // 文档名称/编号
	WorkflowState string  `json:"workflow_state"`     // 工作流状态
	DocStatus     int     `json:"docstatus"`          // 文档状态: 0=草稿, 1=已提交, 2=已取消

	// 提交相关
	SubmittedAt *time.Time `json:"submitted_at,omitempty"` // 提交时间
	SubmittedBy *int64     `json:"submitted_by,omitempty"` // 提交人

	// 取消相关
	CancelledAt  *time.Time `json:"cancelled_at,omitempty"`  // 取消时间
	CancelledBy  *int64     `json:"cancelled_by,omitempty"`  // 取消人
	CancelReason *string    `json:"cancel_reason,omitempty"` // 取消原因

	// 修订相关
	AmendedFrom *int64 `json:"amended_from,omitempty"` // 修订自哪个文档
	IsAmended   bool   `json:"is_amended"`             // 是否为修订版本

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// 关联对象
	SubmittedByUser *User `json:"submitted_by_user,omitempty"`
	CancelledByUser *User `json:"cancelled_by_user,omitempty"`
}

// EnhancedUserPermission 增强用户权限视图数据
type EnhancedUserPermission struct {
	UserID          int64  `json:"user_id"`
	DocType         string `json:"doc_type"`
	PermissionLevel int    `json:"permission_level"`
	RoleCode        string `json:"role_code"`
	RoleName        string `json:"role_name"`

	// 权限列表
	CanRead   bool `json:"can_read"`
	CanWrite  bool `json:"can_write"`
	CanCreate bool `json:"can_create"`
	CanDelete bool `json:"can_delete"`
	CanSubmit bool `json:"can_submit"`
	CanCancel bool `json:"can_cancel"`
	CanAmend  bool `json:"can_amend"`
	CanPrint  bool `json:"can_print"`
	CanEmail  bool `json:"can_email"`
	CanImport bool `json:"can_import"`
	CanExport bool `json:"can_export"`
	CanShare  bool `json:"can_share"`
	CanReport bool `json:"can_report"`

	OnlyIfCreator bool `json:"only_if_creator"`
}

// AccessibleField 可访问字段
type AccessibleField struct {
	FieldName       string `json:"field_name"`
	FieldLabel      string `json:"field_label"`
	PermissionLevel int    `json:"permission_level"`
	CanAccess       bool   `json:"can_access"`
}

// ================================================================
// 权限检查请求和响应
// ================================================================

// PermissionCheckRequest 权限检查请求
type PermissionCheckRequest struct {
	UserID          int64  `json:"user_id"`
	DocType         string `json:"doc_type"`
	Permission      string `json:"permission"`
	PermissionLevel int    `json:"permission_level"`
	DocID           *int64 `json:"doc_id,omitempty"`
	DocCreatorID    *int64 `json:"doc_creator_id,omitempty"`
}

// Validate 验证PermissionCheckRequest数据的完整性和正确性
func (r *PermissionCheckRequest) Validate() error {
	if r.UserID <= 0 {
		return fmt.Errorf("UserID must be positive")
	}

	if strings.TrimSpace(r.DocType) == "" {
		return fmt.Errorf("DocType is required")
	}

	if strings.TrimSpace(r.Permission) == "" {
		return fmt.Errorf("Permission is required")
	}

	if r.PermissionLevel < 0 || r.PermissionLevel > 10 {
		return fmt.Errorf("PermissionLevel must be between 0 and 10")
	}

	// 验证权限类型
	validPermissions := []string{
		"read", "write", "create", "delete", "submit", "cancel", "amend",
		"print", "email", "import", "export", "share", "report",
	}

	isValidPermission := false
	for _, validPerm := range validPermissions {
		if r.Permission == validPerm {
			isValidPermission = true
			break
		}
	}
	if !isValidPermission {
		return fmt.Errorf("invalid permission type: %s", r.Permission)
	}

	// 如果提供了DocID，必须为正数
	if r.DocID != nil && *r.DocID <= 0 {
		return fmt.Errorf("DocID must be positive if provided")
	}

	// 如果提供了DocCreatorID，必须为正数
	if r.DocCreatorID != nil && *r.DocCreatorID <= 0 {
		return fmt.Errorf("DocCreatorID must be positive if provided")
	}

	return nil
}

// PermissionCheckResponse 权限检查响应
type PermissionCheckResponse struct {
	HasPermission bool   `json:"has_permission"`
	Reason        string `json:"reason,omitempty"`
}

// FieldPermissionRequest 字段权限请求
type FieldPermissionRequest struct {
	UserID     int64  `json:"user_id"`
	DocType    string `json:"doc_type"`
	Permission string `json:"permission"`
}

// FieldPermissionResponse 字段权限响应
type FieldPermissionResponse struct {
	Fields []*AccessibleField `json:"fields"`
}

// ================================================================
// 权限管理请求
// ================================================================

// CreatePermissionRuleRequest 创建权限规则请求
type CreatePermissionRuleRequest struct {
	RoleID          int64  `json:"role_id"`
	DocType         string `json:"doc_type"`
	PermissionLevel int    `json:"permission_level"`

	// 权限设置
	CanRead   bool `json:"can_read"`
	CanWrite  bool `json:"can_write"`
	CanCreate bool `json:"can_create"`
	CanDelete bool `json:"can_delete"`
	CanSubmit bool `json:"can_submit"`
	CanCancel bool `json:"can_cancel"`
	CanAmend  bool `json:"can_amend"`
	CanPrint  bool `json:"can_print"`
	CanEmail  bool `json:"can_email"`
	CanImport bool `json:"can_import"`
	CanExport bool `json:"can_export"`
	CanShare  bool `json:"can_share"`
	CanReport bool `json:"can_report"`

	OnlyIfCreator bool `json:"only_if_creator"`
}

// Validate 验证CreatePermissionRuleRequest数据的完整性和正确性
func (r *CreatePermissionRuleRequest) Validate() error {
	if r.RoleID <= 0 {
		return fmt.Errorf("RoleID must be positive")
	}

	if strings.TrimSpace(r.DocType) == "" {
		return fmt.Errorf("DocType is required")
	}

	if r.PermissionLevel < 0 || r.PermissionLevel > 10 {
		return fmt.Errorf("PermissionLevel must be between 0 and 10")
	}

	// 文档级权限验证（permission_level = 0）
	if r.PermissionLevel == 0 {
		// 至少需要一个基础权限
		if !r.CanRead && !r.CanWrite {
			return fmt.Errorf("at least one of CanRead or CanWrite must be true for document level permissions")
		}
	}

	// 字段级权限验证（permission_level > 0）
	if r.PermissionLevel > 0 {
		// 字段级权限不应该有文档级的权限设置
		if r.CanCreate || r.CanDelete || r.CanSubmit || r.CanCancel || r.CanAmend ||
			r.CanPrint || r.CanEmail || r.CanImport || r.CanExport || r.CanShare || r.CanReport {
			return fmt.Errorf("field level permissions (level > 0) should only use CanRead and CanWrite")
		}
	}

	return nil
}

// UpdatePermissionRuleRequest 更新权限规则请求
type UpdatePermissionRuleRequest struct {
	ID int64 `json:"id"`
	CreatePermissionRuleRequest
}

// CreateUserPermissionRequest 创建用户权限请求
type CreateUserPermissionRequest struct {
	UserID          int64   `json:"user_id"`
	DocType         string  `json:"doc_type"`
	DocName         string  `json:"doc_name"`
	Value           string  `json:"value"`
	ApplicableFor   *string `json:"applicable_for,omitempty"`
	HideDescendants bool    `json:"hide_descendants"`
	IsDefault       bool    `json:"is_default"`
}

// Validate 验证CreateUserPermissionRequest数据的完整性和正确性
func (r *CreateUserPermissionRequest) Validate() error {
	if r.UserID <= 0 {
		return fmt.Errorf("UserID must be positive")
	}

	if strings.TrimSpace(r.DocType) == "" {
		return fmt.Errorf("DocType is required")
	}

	if strings.TrimSpace(r.DocName) == "" {
		return fmt.Errorf("DocName is required")
	}

	if strings.TrimSpace(r.Value) == "" {
		return fmt.Errorf("Value is required")
	}

	// 验证ApplicableFor如果不为空，必须是有效的DocType
	if r.ApplicableFor != nil && strings.TrimSpace(*r.ApplicableFor) == "" {
		return fmt.Errorf("ApplicableFor cannot be empty string if provided")
	}

	return nil
}

// UpdateUserPermissionRequest 更新用户权限请求
type UpdateUserPermissionRequest struct {
	ID int64 `json:"id"`
	CreateUserPermissionRequest
}

// ================================================================
// Repository接口定义
// ================================================================

// PermissionRepo 权限仓库接口
type PermissionRepo interface {
	// 文档类型管理
	CreateDocType(ctx context.Context, docType *DocType) (*DocType, error)
	UpdateDocType(ctx context.Context, docType *DocType) (*DocType, error)
	GetDocType(ctx context.Context, name string) (*DocType, error)
	ListDocTypes(ctx context.Context, module string) ([]*DocType, error)
	DeleteDocType(ctx context.Context, name string) error

	// 权限规则管理
	CreatePermissionRule(ctx context.Context, rule *PermissionRule) (*PermissionRule, error)
	UpdatePermissionRule(ctx context.Context, rule *PermissionRule) (*PermissionRule, error)
	GetPermissionRule(ctx context.Context, id int64) (*PermissionRule, error)
	ListPermissionRules(ctx context.Context, roleID int64, docType string) ([]*PermissionRule, error)
	DeletePermissionRule(ctx context.Context, id int64) error

	// 用户权限管理
	CreateUserPermission(ctx context.Context, userPermission *UserPermission) (*UserPermission, error)
	UpdateUserPermission(ctx context.Context, userPermission *UserPermission) (*UserPermission, error)
	GetUserPermission(ctx context.Context, id int64) (*UserPermission, error)
	ListUserPermissions(ctx context.Context, userID int64, docType string, page, size int32) ([]*UserPermission, error)
	GetUserPermissionsCount(ctx context.Context, userID int64, docType string) (int32, error)
	DeleteUserPermission(ctx context.Context, id int64) error

	// 字段权限级别管理
	CreateFieldPermissionLevel(ctx context.Context, field *FieldPermissionLevel) (*FieldPermissionLevel, error)
	UpdateFieldPermissionLevel(ctx context.Context, field *FieldPermissionLevel) (*FieldPermissionLevel, error)
	GetFieldPermissionLevel(ctx context.Context, id int64) (*FieldPermissionLevel, error)
	ListFieldPermissionLevels(ctx context.Context, docType string, page, size int32) ([]*FieldPermissionLevel, error)
	GetFieldPermissionLevelsCount(ctx context.Context, docType string) (int32, error)
	DeleteFieldPermissionLevel(ctx context.Context, id int64) error

	// 文档工作流状态管理
	CreateDocumentWorkflowState(ctx context.Context, state *DocumentWorkflowState) (*DocumentWorkflowState, error)
	UpdateDocumentWorkflowState(ctx context.Context, state *DocumentWorkflowState) (*DocumentWorkflowState, error)
	GetDocumentWorkflowState(ctx context.Context, stateID int64) (*DocumentWorkflowState, error)
	ListDocumentWorkflowStates(ctx context.Context, docType, documentName, state string, userID int64, page, size int32) ([]*DocumentWorkflowState, error)
	GetDocumentWorkflowStatesCount(ctx context.Context, docType, documentName, state string, userID int64) (int32, error)
	DeleteDocumentWorkflowState(ctx context.Context, stateID int64) error

	// 权限检查和查询
	CheckDocumentPermission(ctx context.Context, req *PermissionCheckRequest) (bool, error)
	CheckPermission(ctx context.Context, userID int64, documentType, action string, permissionLevel int) (bool, error)
	GetUserPermissionLevel(ctx context.Context, userID int64, documentType string) (int, error)
	GetUserEnhancedPermissions(ctx context.Context, userID int64, docType string) ([]*EnhancedUserPermission, error)
	GetAccessibleFields(ctx context.Context, req *FieldPermissionRequest) ([]*AccessibleField, error)
	FilterDocumentsByPermission(ctx context.Context, userID int64, documentType string, documents []map[string]interface{}) ([]map[string]interface{}, error)
	GetUserRoles(ctx context.Context, userID int64) ([]string, error)

	// 批量操作
	BatchCreatePermissionRules(ctx context.Context, rules []*PermissionRule) error
	BatchCreateUserPermissions(ctx context.Context, permissions []*UserPermission) error
}

// ================================================================
// 权限管理用例接口和实现
// ================================================================

// PermissionUsecaseInterface 权限管理用例接口
type PermissionUsecaseInterface interface {
	// 文档类型管理
	CreateDocType(ctx context.Context, docType *DocType) (*DocType, error)
	GetDocType(ctx context.Context, name string) (*DocType, error)
	ListDocTypes(ctx context.Context, module string) ([]*DocType, error)
	UpdateDocType(ctx context.Context, docType *DocType) (*DocType, error)
	DeleteDocType(ctx context.Context, name string) error

	// 权限规则管理
	CreatePermissionRule(ctx context.Context, rule *PermissionRule) (*PermissionRule, error)
	CreatePermissionRuleFromRequest(ctx context.Context, req *CreatePermissionRuleRequest) (*PermissionRule, error)
	GetPermissionRule(ctx context.Context, id int64) (*PermissionRule, error)
	ListPermissionRules(ctx context.Context, roleID int64, docType string) ([]*PermissionRule, error)
	UpdatePermissionRule(ctx context.Context, rule *PermissionRule) (*PermissionRule, error)
	DeletePermissionRule(ctx context.Context, id int64) error
	BatchCreatePermissionRules(ctx context.Context, rules []*PermissionRule) error

	// 用户权限管理
	CreateUserPermission(ctx context.Context, userPermission *UserPermission) (*UserPermission, error)
	CreateUserPermissionFromRequest(ctx context.Context, req *CreateUserPermissionRequest) (*UserPermission, error)
	GetUserPermission(ctx context.Context, id int64) (*UserPermission, error)
	ListUserPermissions(ctx context.Context, userID int64, docType string, page, size int32) ([]*UserPermission, error)
	GetUserPermissionsCount(ctx context.Context, userID int64, docType string) (int32, error)
	UpdateUserPermission(ctx context.Context, userPerm *UserPermission) (*UserPermission, error)
	DeleteUserPermission(ctx context.Context, id int64) error
	BatchCreateUserPermissions(ctx context.Context, permissions []*UserPermission) error

	// 字段权限级别管理
	CreateFieldPermissionLevel(ctx context.Context, fieldPerm *FieldPermissionLevel) (*FieldPermissionLevel, error)
	GetFieldPermissionLevel(ctx context.Context, id int64) (*FieldPermissionLevel, error)
	UpdateFieldPermissionLevel(ctx context.Context, fieldPerm *FieldPermissionLevel) (*FieldPermissionLevel, error)
	DeleteFieldPermissionLevel(ctx context.Context, id int64) error
	ListFieldPermissionLevels(ctx context.Context, docType string, page, size int32) ([]*FieldPermissionLevel, error)
	GetFieldPermissionLevelsCount(ctx context.Context, docType string) (int32, error)

	// 文档工作流状态管理
	CreateDocumentWorkflowState(ctx context.Context, state *DocumentWorkflowState) (*DocumentWorkflowState, error)
	GetDocumentWorkflowState(ctx context.Context, stateID int64) (*DocumentWorkflowState, error)
	UpdateDocumentWorkflowState(ctx context.Context, state *DocumentWorkflowState) (*DocumentWorkflowState, error)
	DeleteDocumentWorkflowState(ctx context.Context, stateID int64) error
	ListDocumentWorkflowStates(ctx context.Context, docType, documentName, state string, userID int64, page, size int32) ([]*DocumentWorkflowState, error)
	GetDocumentWorkflowStatesCount(ctx context.Context, docType, documentName, state string, userID int64) (int32, error)

	// 权限检查和查询
	CheckDocumentPermission(ctx context.Context, req *PermissionCheckRequest) (*PermissionCheckResponse, error)
	CheckUserPermission(ctx context.Context, userID int64, permissionCode string) (bool, error)
	CheckPermission(ctx context.Context, userID int64, documentType, action string, permissionLevel int) (bool, error)
	GetUserPermissionLevel(ctx context.Context, userID int64, documentType string) (int, error)
	GetUserEnhancedPermissions(ctx context.Context, userID int64, docType string) ([]*EnhancedUserPermission, error)
	GetAccessibleFields(ctx context.Context, req *FieldPermissionRequest) (*FieldPermissionResponse, error)
	FilterDocumentsByPermission(ctx context.Context, userID int64, documentType string, documents []map[string]interface{}) ([]map[string]interface{}, error)
	GetUserRoles(ctx context.Context, userID int64) ([]string, error)
}

// PermissionUsecase 权限管理用例
type PermissionUsecase struct {
	repo   PermissionRepo
	logger *log.Helper
}

// 确保 PermissionUsecase 实现了 PermissionUsecaseInterface 接口
var _ PermissionUsecaseInterface = (*PermissionUsecase)(nil)

// NewPermissionUsecase 创建权限用例
func NewPermissionUsecase(repo PermissionRepo, logger log.Logger) *PermissionUsecase {
	return &PermissionUsecase{
		repo:   repo,
		logger: log.NewHelper(logger),
	}
}

// 文档类型管理
func (uc *PermissionUsecase) CreateDocType(ctx context.Context, req *DocType) (*DocType, error) {
	req.CreatedAt = time.Now()
	req.UpdatedAt = time.Now()
	return uc.repo.CreateDocType(ctx, req)
}

func (uc *PermissionUsecase) GetDocType(ctx context.Context, name string) (*DocType, error) {
	return uc.repo.GetDocType(ctx, name)
}

func (uc *PermissionUsecase) ListDocTypes(ctx context.Context, module string) ([]*DocType, error) {
	return uc.repo.ListDocTypes(ctx, module)
}

// 权限规则管理
func (uc *PermissionUsecase) CreatePermissionRule(ctx context.Context, rule *PermissionRule) (*PermissionRule, error) {
	rule.CreatedAt = time.Now()
	rule.UpdatedAt = time.Now()
	return uc.repo.CreatePermissionRule(ctx, rule)
}

func (uc *PermissionUsecase) CreatePermissionRuleFromRequest(ctx context.Context, req *CreatePermissionRuleRequest) (*PermissionRule, error) {
	rule := &PermissionRule{
		RoleID:          req.RoleID,
		DocType:         req.DocType,
		PermissionLevel: req.PermissionLevel,
		CanRead:         req.CanRead,
		CanWrite:        req.CanWrite,
		CanCreate:       req.CanCreate,
		CanDelete:       req.CanDelete,
		CanSubmit:       req.CanSubmit,
		CanCancel:       req.CanCancel,
		CanAmend:        req.CanAmend,
		CanPrint:        req.CanPrint,
		CanEmail:        req.CanEmail,
		CanImport:       req.CanImport,
		CanExport:       req.CanExport,
		CanShare:        req.CanShare,
		CanReport:       req.CanReport,
		OnlyIfCreator:   req.OnlyIfCreator,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	return uc.repo.CreatePermissionRule(ctx, rule)
}

func (uc *PermissionUsecase) ListPermissionRules(ctx context.Context, roleID int64, docType string) ([]*PermissionRule, error) {
	return uc.repo.ListPermissionRules(ctx, roleID, docType)
}

// 用户权限管理
func (uc *PermissionUsecase) CreateUserPermission(ctx context.Context, userPermission *UserPermission) (*UserPermission, error) {
	userPermission.CreatedAt = time.Now()
	userPermission.UpdatedAt = time.Now()
	return uc.repo.CreateUserPermission(ctx, userPermission)
}

func (uc *PermissionUsecase) CreateUserPermissionFromRequest(ctx context.Context, req *CreateUserPermissionRequest) (*UserPermission, error) {
	userPermission := &UserPermission{
		UserID:          req.UserID,
		DocType:         req.DocType,
		DocName:         req.DocName,
		Value:           req.Value,
		ApplicableFor:   req.ApplicableFor,
		HideDescendants: req.HideDescendants,
		IsDefault:       req.IsDefault,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	return uc.repo.CreateUserPermission(ctx, userPermission)
}

// 权限检查
func (uc *PermissionUsecase) CheckDocumentPermission(ctx context.Context, req *PermissionCheckRequest) (*PermissionCheckResponse, error) {
	hasPermission, err := uc.repo.CheckDocumentPermission(ctx, req)
	if err != nil {
		return nil, err
	}

	response := &PermissionCheckResponse{
		HasPermission: hasPermission,
	}

	if !hasPermission {
		response.Reason = "权限不足"
	}

	return response, nil
}

func (uc *PermissionUsecase) GetAccessibleFields(ctx context.Context, req *FieldPermissionRequest) (*FieldPermissionResponse, error) {
	fields, err := uc.repo.GetAccessibleFields(ctx, req)
	if err != nil {
		return nil, err
	}

	return &FieldPermissionResponse{
		Fields: fields,
	}, nil
}

// CheckUserPermission 检查用户权限（向后兼容方法）
func (uc *PermissionUsecase) CheckUserPermission(ctx context.Context, userID int64, permissionCode string) (bool, error) {
	// 对于旧的权限代码检查，我们可以将其映射到新的权限系统
	// 这里先实现一个简单的检查，后续可以根据需要完善
	return uc.repo.CheckPermission(ctx, userID, "User", "read", 0)
}

// CheckPermission 检查Frappe风格的权限
func (uc *PermissionUsecase) CheckPermission(ctx context.Context, userID int64, documentType, action string, permissionLevel int) (bool, error) {
	return uc.repo.CheckPermission(ctx, userID, documentType, action, permissionLevel)
}

// GetUserPermissionLevel 获取用户权限级别
func (uc *PermissionUsecase) GetUserPermissionLevel(ctx context.Context, userID int64, documentType string) (int, error) {
	return uc.repo.GetUserPermissionLevel(ctx, userID, documentType)
}

// ================================================================
// 额外的业务逻辑方法实现
// ================================================================

// DocType 管理方法 - 附加实现
func (uc *PermissionUsecase) UpdateDocType(ctx context.Context, docType *DocType) (*DocType, error) {
	docType.UpdatedAt = time.Now()
	return uc.repo.UpdateDocType(ctx, docType)
}

func (uc *PermissionUsecase) DeleteDocType(ctx context.Context, name string) error {
	return uc.repo.DeleteDocType(ctx, name)
}

// PermissionRule 管理方法 - 附加实现
func (uc *PermissionUsecase) GetPermissionRule(ctx context.Context, id int64) (*PermissionRule, error) {
	return uc.repo.GetPermissionRule(ctx, id)
}

func (uc *PermissionUsecase) UpdatePermissionRule(ctx context.Context, rule *PermissionRule) (*PermissionRule, error) {
	rule.UpdatedAt = time.Now()
	return uc.repo.UpdatePermissionRule(ctx, rule)
}

func (uc *PermissionUsecase) DeletePermissionRule(ctx context.Context, id int64) error {
	return uc.repo.DeletePermissionRule(ctx, id)
}

func (uc *PermissionUsecase) BatchCreatePermissionRules(ctx context.Context, rules []*PermissionRule) error {
	return uc.repo.BatchCreatePermissionRules(ctx, rules)
}

// UserPermission 管理方法 - 附加实现
func (uc *PermissionUsecase) GetUserPermission(ctx context.Context, id int64) (*UserPermission, error) {
	return uc.repo.GetUserPermission(ctx, id)
}

func (uc *PermissionUsecase) UpdateUserPermission(ctx context.Context, userPerm *UserPermission) (*UserPermission, error) {
	userPerm.UpdatedAt = time.Now()
	return uc.repo.UpdateUserPermission(ctx, userPerm)
}

func (uc *PermissionUsecase) DeleteUserPermission(ctx context.Context, id int64) error {
	return uc.repo.DeleteUserPermission(ctx, id)
}

func (uc *PermissionUsecase) ListUserPermissions(ctx context.Context, userID int64, docType string, page, size int32) ([]*UserPermission, error) {
	return uc.repo.ListUserPermissions(ctx, userID, docType, page, size)
}

func (uc *PermissionUsecase) GetUserPermissionsCount(ctx context.Context, userID int64, docType string) (int32, error) {
	return uc.repo.GetUserPermissionsCount(ctx, userID, docType)
}

func (uc *PermissionUsecase) BatchCreateUserPermissions(ctx context.Context, permissions []*UserPermission) error {
	return uc.repo.BatchCreateUserPermissions(ctx, permissions)
}

// FieldPermissionLevel 管理方法
func (uc *PermissionUsecase) CreateFieldPermissionLevel(ctx context.Context, fieldPerm *FieldPermissionLevel) (*FieldPermissionLevel, error) {
	return uc.repo.CreateFieldPermissionLevel(ctx, fieldPerm)
}

func (uc *PermissionUsecase) GetFieldPermissionLevel(ctx context.Context, id int64) (*FieldPermissionLevel, error) {
	return uc.repo.GetFieldPermissionLevel(ctx, id)
}

func (uc *PermissionUsecase) UpdateFieldPermissionLevel(ctx context.Context, fieldPerm *FieldPermissionLevel) (*FieldPermissionLevel, error) {
	return uc.repo.UpdateFieldPermissionLevel(ctx, fieldPerm)
}

func (uc *PermissionUsecase) DeleteFieldPermissionLevel(ctx context.Context, id int64) error {
	return uc.repo.DeleteFieldPermissionLevel(ctx, id)
}

func (uc *PermissionUsecase) ListFieldPermissionLevels(ctx context.Context, docType string, page, size int32) ([]*FieldPermissionLevel, error) {
	return uc.repo.ListFieldPermissionLevels(ctx, docType, page, size)
}

func (uc *PermissionUsecase) GetFieldPermissionLevelsCount(ctx context.Context, docType string) (int32, error) {
	return uc.repo.GetFieldPermissionLevelsCount(ctx, docType)
}

// DocumentWorkflowState 管理方法
func (uc *PermissionUsecase) CreateDocumentWorkflowState(ctx context.Context, state *DocumentWorkflowState) (*DocumentWorkflowState, error) {
	return uc.repo.CreateDocumentWorkflowState(ctx, state)
}

func (uc *PermissionUsecase) GetDocumentWorkflowState(ctx context.Context, stateID int64) (*DocumentWorkflowState, error) {
	return uc.repo.GetDocumentWorkflowState(ctx, stateID)
}

func (uc *PermissionUsecase) UpdateDocumentWorkflowState(ctx context.Context, state *DocumentWorkflowState) (*DocumentWorkflowState, error) {
	return uc.repo.UpdateDocumentWorkflowState(ctx, state)
}

func (uc *PermissionUsecase) DeleteDocumentWorkflowState(ctx context.Context, stateID int64) error {
	return uc.repo.DeleteDocumentWorkflowState(ctx, stateID)
}

func (uc *PermissionUsecase) ListDocumentWorkflowStates(ctx context.Context, docType, documentName, state string, userID int64, page, size int32) ([]*DocumentWorkflowState, error) {
	return uc.repo.ListDocumentWorkflowStates(ctx, docType, documentName, state, userID, page, size)
}

func (uc *PermissionUsecase) GetDocumentWorkflowStatesCount(ctx context.Context, docType, documentName, state string, userID int64) (int32, error) {
	return uc.repo.GetDocumentWorkflowStatesCount(ctx, docType, documentName, state, userID)
}

// 权限检查和查询方法 - 附加实现

func (uc *PermissionUsecase) GetUserEnhancedPermissions(ctx context.Context, userID int64, docType string) ([]*EnhancedUserPermission, error) {
	return uc.repo.GetUserEnhancedPermissions(ctx, userID, docType)
}

func (uc *PermissionUsecase) FilterDocumentsByPermission(ctx context.Context, userID int64, documentType string, documents []map[string]interface{}) ([]map[string]interface{}, error) {
	return uc.repo.FilterDocumentsByPermission(ctx, userID, documentType, documents)
}

func (uc *PermissionUsecase) GetUserRoles(ctx context.Context, userID int64) ([]string, error) {
	return uc.repo.GetUserRoles(ctx, userID)
}
