package service

import (
	"context"
	"fmt"
	"time"

	"erp-system/internal/biz"
	"erp-system/internal/middleware"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
)

// FrappePermissionService Frappe权限管理服务
type FrappePermissionService struct {
	permissionUc biz.PermissionUsecaseInterface
	log          *log.Helper
}

// NewFrappePermissionService 创建Frappe权限服务
func NewFrappePermissionService(permissionUc biz.PermissionUsecaseInterface, logger log.Logger) *FrappePermissionService {
	return &FrappePermissionService{
		permissionUc: permissionUc,
		log:          log.NewHelper(logger),
	}
}

// ================================================================
// DocType 文档类型管理 API
// ================================================================

// CreateDocTypeRequest 创建文档类型请求
type CreateDocTypeRequest struct {
	Name           string                 `json:"name" validate:"required,min=2,max=100"`
	Label          string                 `json:"label" validate:"required,min=2,max=100"`
	Module         string                 `json:"module" validate:"required"`
	Description    string                 `json:"description"`
	IsSubmittable  bool                   `json:"is_submittable"`
	IsChildTable   bool                   `json:"is_child_table"`
	HasWorkflow    bool                   `json:"has_workflow"`
	TrackChanges   bool                   `json:"track_changes"`
	AppliesToAll   bool                   `json:"applies_to_all_users"`
	MaxAttachments int                    `json:"max_attachments"`
	Permissions    map[string]interface{} `json:"permissions"`
	NamingRule     string                 `json:"naming_rule"`
	TitleField     string                 `json:"title_field"`
	SearchFields   []string               `json:"search_fields"`
	SortField      string                 `json:"sort_field"`
	SortOrder      string                 `json:"sort_order"`
}

// UpdateDocTypeRequest 更新文档类型请求
type UpdateDocTypeRequest struct {
	Name string `json:"name" validate:"required"`
	CreateDocTypeRequest
}

// DocTypeInfo 文档类型信息
type DocTypeInfo struct {
	ID             int64                  `json:"id"`
	Name           string                 `json:"name"`
	Label          string                 `json:"label"`
	Module         string                 `json:"module"`
	Description    string                 `json:"description"`
	IsSubmittable  bool                   `json:"is_submittable"`
	IsChildTable   bool                   `json:"is_child_table"`
	HasWorkflow    bool                   `json:"has_workflow"`
	TrackChanges   bool                   `json:"track_changes"`
	AppliesToAll   bool                   `json:"applies_to_all_users"`
	MaxAttachments int                    `json:"max_attachments"`
	Permissions    map[string]interface{} `json:"permissions"`
	NamingRule     string                 `json:"naming_rule"`
	TitleField     string                 `json:"title_field"`
	SearchFields   []string               `json:"search_fields"`
	SortField      string                 `json:"sort_field"`
	SortOrder      string                 `json:"sort_order"`
	Version        int                    `json:"version"`
	CreatedAt      time.Time              `json:"created_at"`
	UpdatedAt      time.Time              `json:"updated_at"`
}

// CreateDocType 创建文档类型
func (s *FrappePermissionService) CreateDocType(ctx context.Context, req *CreateDocTypeRequest) (*DocTypeInfo, error) {
	docType := &biz.DocType{
		Name:           req.Name,
		Label:          req.Label,
		Module:         req.Module,
		Description:    req.Description,
		IsSubmittable:  req.IsSubmittable,
		IsChildTable:   req.IsChildTable,
		HasWorkflow:    req.HasWorkflow,
		TrackChanges:   req.TrackChanges,
		AppliesToAll:   req.AppliesToAll,
		MaxAttachments: req.MaxAttachments,
		Permissions:    req.Permissions,
		NamingRule:     req.NamingRule,
		TitleField:     req.TitleField,
		SearchFields:   req.SearchFields,
		SortField:      req.SortField,
		SortOrder:      req.SortOrder,
		Version:        1,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	// 验证数据完整性
	if err := docType.Validate(); err != nil {
		s.log.Errorf("DocType validation failed: %v", err)
		return nil, errors.BadRequest("INVALID_DATA", fmt.Sprintf("数据验证失败: %v", err))
	}

	createdDocType, err := s.permissionUc.CreateDocType(ctx, docType)
	if err != nil {
		s.log.Errorf("Failed to create doctype: %v", err)
		return nil, errors.InternalServer("INTERNAL_ERROR", "文档类型创建失败")
	}

	s.log.Infof("DocType created successfully: %s", req.Name)

	return &DocTypeInfo{
		ID:             createdDocType.ID,
		Name:           createdDocType.Name,
		Label:          createdDocType.Label,
		Module:         createdDocType.Module,
		Description:    createdDocType.Description,
		IsSubmittable:  createdDocType.IsSubmittable,
		IsChildTable:   createdDocType.IsChildTable,
		HasWorkflow:    createdDocType.HasWorkflow,
		TrackChanges:   createdDocType.TrackChanges,
		AppliesToAll:   createdDocType.AppliesToAll,
		MaxAttachments: createdDocType.MaxAttachments,
		Permissions:    createdDocType.Permissions,
		NamingRule:     createdDocType.NamingRule,
		TitleField:     createdDocType.TitleField,
		SearchFields:   createdDocType.SearchFields,
		SortField:      createdDocType.SortField,
		SortOrder:      createdDocType.SortOrder,
		Version:        createdDocType.Version,
		CreatedAt:      createdDocType.CreatedAt,
		UpdatedAt:      createdDocType.UpdatedAt,
	}, nil
}

// GetDocType 获取文档类型
func (s *FrappePermissionService) GetDocType(ctx context.Context, name string) (*DocTypeInfo, error) {
	if name == "" {
		return nil, errors.BadRequest("INVALID_PARAMETER", "文档类型名称不能为空")
	}

	docType, err := s.permissionUc.GetDocType(ctx, name)
	if err != nil {
		s.log.Errorf("Failed to get doctype: %v", err)
		return nil, errors.NotFound("NOT_FOUND", "文档类型不存在")
	}

	return &DocTypeInfo{
		ID:             docType.ID,
		Name:           docType.Name,
		Label:          docType.Label,
		Module:         docType.Module,
		Description:    docType.Description,
		IsSubmittable:  docType.IsSubmittable,
		IsChildTable:   docType.IsChildTable,
		HasWorkflow:    docType.HasWorkflow,
		TrackChanges:   docType.TrackChanges,
		AppliesToAll:   docType.AppliesToAll,
		MaxAttachments: docType.MaxAttachments,
		Permissions:    docType.Permissions,
		NamingRule:     docType.NamingRule,
		TitleField:     docType.TitleField,
		SearchFields:   docType.SearchFields,
		SortField:      docType.SortField,
		SortOrder:      docType.SortOrder,
		Version:        docType.Version,
		CreatedAt:      docType.CreatedAt,
		UpdatedAt:      docType.UpdatedAt,
	}, nil
}

// UpdateDocType 更新文档类型
func (s *FrappePermissionService) UpdateDocType(ctx context.Context, req *UpdateDocTypeRequest) (*DocTypeInfo, error) {
	if req.Name == "" {
		return nil, errors.BadRequest("INVALID_PARAMETER", "文档类型名称不能为空")
	}

	docType := &biz.DocType{
		Name:           req.Name,
		Label:          req.Label,
		Module:         req.Module,
		Description:    req.Description,
		IsSubmittable:  req.IsSubmittable,
		IsChildTable:   req.IsChildTable,
		HasWorkflow:    req.HasWorkflow,
		TrackChanges:   req.TrackChanges,
		AppliesToAll:   req.AppliesToAll,
		MaxAttachments: req.MaxAttachments,
		Permissions:    req.Permissions,
		NamingRule:     req.NamingRule,
		TitleField:     req.TitleField,
		SearchFields:   req.SearchFields,
		SortField:      req.SortField,
		SortOrder:      req.SortOrder,
		UpdatedAt:      time.Now(),
	}

	updatedDocType, err := s.permissionUc.UpdateDocType(ctx, docType)
	if err != nil {
		s.log.Errorf("Failed to update doctype: %v", err)
		return nil, errors.InternalServer("INTERNAL_ERROR", "文档类型更新失败")
	}

	s.log.Infof("DocType updated successfully: %s", req.Name)

	return &DocTypeInfo{
		ID:             updatedDocType.ID,
		Name:           updatedDocType.Name,
		Label:          updatedDocType.Label,
		Module:         updatedDocType.Module,
		Description:    updatedDocType.Description,
		IsSubmittable:  updatedDocType.IsSubmittable,
		IsChildTable:   updatedDocType.IsChildTable,
		HasWorkflow:    updatedDocType.HasWorkflow,
		TrackChanges:   updatedDocType.TrackChanges,
		AppliesToAll:   updatedDocType.AppliesToAll,
		MaxAttachments: updatedDocType.MaxAttachments,
		Permissions:    updatedDocType.Permissions,
		NamingRule:     updatedDocType.NamingRule,
		TitleField:     updatedDocType.TitleField,
		SearchFields:   updatedDocType.SearchFields,
		SortField:      updatedDocType.SortField,
		SortOrder:      updatedDocType.SortOrder,
		Version:        updatedDocType.Version,
		CreatedAt:      updatedDocType.CreatedAt,
		UpdatedAt:      updatedDocType.UpdatedAt,
	}, nil
}

// DeleteDocType 删除文档类型
func (s *FrappePermissionService) DeleteDocType(ctx context.Context, name string) error {
	if name == "" {
		return errors.BadRequest("INVALID_PARAMETER", "文档类型名称不能为空")
	}

	err := s.permissionUc.DeleteDocType(ctx, name)
	if err != nil {
		s.log.Errorf("Failed to delete doctype: %v", err)
		return errors.InternalServer("INTERNAL_ERROR", "文档类型删除失败")
	}

	s.log.Infof("DocType deleted successfully: %s", name)
	return nil
}

// ListDocTypesRequest 列出文档类型请求
type ListDocTypesRequest struct {
	Module string `json:"module"`
}

// ListDocTypesResponse 列出文档类型响应
type ListDocTypesResponse struct {
	DocTypes []*DocTypeInfo `json:"doc_types"`
	Total    int            `json:"total"`
}

// ListDocTypes 列出文档类型
func (s *FrappePermissionService) ListDocTypes(ctx context.Context, req *ListDocTypesRequest) (*ListDocTypesResponse, error) {
	docTypes, err := s.permissionUc.ListDocTypes(ctx, req.Module)
	if err != nil {
		s.log.Errorf("Failed to list doctypes: %v", err)
		return nil, errors.InternalServer("INTERNAL_ERROR", "获取文档类型列表失败")
	}

	var docTypeInfos []*DocTypeInfo
	for _, docType := range docTypes {
		docTypeInfos = append(docTypeInfos, &DocTypeInfo{
			ID:             docType.ID,
			Name:           docType.Name,
			Label:          docType.Label,
			Module:         docType.Module,
			Description:    docType.Description,
			IsSubmittable:  docType.IsSubmittable,
			IsChildTable:   docType.IsChildTable,
			HasWorkflow:    docType.HasWorkflow,
			TrackChanges:   docType.TrackChanges,
			AppliesToAll:   docType.AppliesToAll,
			MaxAttachments: docType.MaxAttachments,
			Permissions:    docType.Permissions,
			NamingRule:     docType.NamingRule,
			TitleField:     docType.TitleField,
			SearchFields:   docType.SearchFields,
			SortField:      docType.SortField,
			SortOrder:      docType.SortOrder,
			Version:        docType.Version,
			CreatedAt:      docType.CreatedAt,
			UpdatedAt:      docType.UpdatedAt,
		})
	}

	return &ListDocTypesResponse{
		DocTypes: docTypeInfos,
		Total:    len(docTypeInfos),
	}, nil
}

// ================================================================
// PermissionRule 权限规则管理 API
// ================================================================

// CreatePermissionRuleRequest 创建权限规则请求
type CreatePermissionRuleRequest struct {
	RoleID          int64  `json:"role_id" validate:"required"`
	DocType         string `json:"doc_type" validate:"required"`
	PermissionLevel int    `json:"permission_level" validate:"min=0,max=9"`

	// 基础权限
	CanRead  bool `json:"can_read"`
	CanWrite bool `json:"can_write"`

	// 文档级权限 (仅permission_level=0时有效)
	CanCreate             bool `json:"can_create"`
	CanDelete             bool `json:"can_delete"`
	CanSubmit             bool `json:"can_submit"`
	CanCancel             bool `json:"can_cancel"`
	CanAmend              bool `json:"can_amend"`
	CanPrint              bool `json:"can_print"`
	CanEmail              bool `json:"can_email"`
	CanImport             bool `json:"can_import"`
	CanExport             bool `json:"can_export"`
	CanShare              bool `json:"can_share"`
	CanReport             bool `json:"can_report"`
	CanSetUserPermissions bool `json:"can_set_user_permissions"`

	// 条件权限
	OnlyIfCreator bool `json:"only_if_creator"`
}

// Validate validates the CreatePermissionRuleRequest
func (req *CreatePermissionRuleRequest) Validate() error {
	if req.RoleID <= 0 {
		return fmt.Errorf("role_id is required and must be positive")
	}
	if req.DocType == "" {
		return fmt.Errorf("doc_type is required")
	}
	if req.PermissionLevel < 0 || req.PermissionLevel > 9 {
		return fmt.Errorf("permission_level must be between 0 and 9")
	}
	return nil
}

// UpdatePermissionRuleRequest 更新权限规则请求
type UpdatePermissionRuleRequest struct {
	ID int64 `json:"id" validate:"required"`
	CreatePermissionRuleRequest
}

// PermissionRuleInfo 权限规则信息
type PermissionRuleInfo struct {
	ID              int64  `json:"id"`
	RoleID          int64  `json:"role_id"`
	DocType         string `json:"doc_type"`
	PermissionLevel int    `json:"permission_level"`

	// 基础权限
	CanRead  bool `json:"can_read"`
	CanWrite bool `json:"can_write"`

	// 文档级权限
	CanCreate             bool `json:"can_create"`
	CanDelete             bool `json:"can_delete"`
	CanSubmit             bool `json:"can_submit"`
	CanCancel             bool `json:"can_cancel"`
	CanAmend              bool `json:"can_amend"`
	CanPrint              bool `json:"can_print"`
	CanEmail              bool `json:"can_email"`
	CanImport             bool `json:"can_import"`
	CanExport             bool `json:"can_export"`
	CanShare              bool `json:"can_share"`
	CanReport             bool `json:"can_report"`
	CanSetUserPermissions bool `json:"can_set_user_permissions"`

	// 条件权限
	OnlyIfCreator bool `json:"only_if_creator"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// CreatePermissionRule 创建权限规则
func (s *FrappePermissionService) CreatePermissionRule(ctx context.Context, req *CreatePermissionRuleRequest) (*PermissionRuleInfo, error) {
	// 先验证请求数据
	if err := req.Validate(); err != nil {
		s.log.Errorf("CreatePermissionRuleRequest validation failed: %v", err)
		return nil, errors.BadRequest("INVALID_DATA", fmt.Sprintf("请求数据验证失败: %v", err))
	}

	rule := &biz.PermissionRule{
		RoleID:                req.RoleID,
		DocType:               req.DocType,
		PermissionLevel:       req.PermissionLevel,
		CanRead:               req.CanRead,
		CanWrite:              req.CanWrite,
		CanCreate:             req.CanCreate,
		CanDelete:             req.CanDelete,
		CanSubmit:             req.CanSubmit,
		CanCancel:             req.CanCancel,
		CanAmend:              req.CanAmend,
		CanPrint:              req.CanPrint,
		CanEmail:              req.CanEmail,
		CanImport:             req.CanImport,
		CanExport:             req.CanExport,
		CanShare:              req.CanShare,
		CanReport:             req.CanReport,
		CanSetUserPermissions: req.CanSetUserPermissions,
		OnlyIfCreator:         req.OnlyIfCreator,
		CreatedAt:             time.Now(),
		UpdatedAt:             time.Now(),
	}

	// 验证业务实体数据
	if err := rule.Validate(); err != nil {
		s.log.Errorf("PermissionRule validation failed: %v", err)
		return nil, errors.BadRequest("INVALID_DATA", fmt.Sprintf("权限规则验证失败: %v", err))
	}

	createdRule, err := s.permissionUc.CreatePermissionRule(ctx, rule)
	if err != nil {
		s.log.Errorf("Failed to create permission rule: %v", err)
		return nil, errors.InternalServer("INTERNAL_ERROR", "权限规则创建失败")
	}

	s.log.Infof("Permission rule created successfully for role %d, doctype %s", req.RoleID, req.DocType)

	return &PermissionRuleInfo{
		ID:              createdRule.ID,
		RoleID:          createdRule.RoleID,
		DocType:         createdRule.DocType,
		PermissionLevel: createdRule.PermissionLevel,
		CanRead:         createdRule.CanRead,
		CanWrite:        createdRule.CanWrite,
		CanCreate:       createdRule.CanCreate,
		CanDelete:       createdRule.CanDelete,
		CanSubmit:       createdRule.CanSubmit,
		CanCancel:       createdRule.CanCancel,
		CanAmend:        createdRule.CanAmend,
		CanPrint:        createdRule.CanPrint,
		CanEmail:        createdRule.CanEmail,
		CanImport:       createdRule.CanImport,
		CanExport:       createdRule.CanExport,
		CanShare:        createdRule.CanShare,
		CanReport:       createdRule.CanReport,
		OnlyIfCreator:   createdRule.OnlyIfCreator,
		CreatedAt:       createdRule.CreatedAt,
		UpdatedAt:       createdRule.UpdatedAt,
	}, nil
}

// GetPermissionRule 获取权限规则
func (s *FrappePermissionService) GetPermissionRule(ctx context.Context, id int64) (*PermissionRuleInfo, error) {
	if id <= 0 {
		return nil, errors.BadRequest("INVALID_PARAMETER", "权限规则ID无效")
	}

	rule, err := s.permissionUc.GetPermissionRule(ctx, id)
	if err != nil {
		s.log.Errorf("Failed to get permission rule: %v", err)
		return nil, errors.NotFound("NOT_FOUND", "权限规则不存在")
	}

	return &PermissionRuleInfo{
		ID:              rule.ID,
		RoleID:          rule.RoleID,
		DocType:         rule.DocType,
		PermissionLevel: rule.PermissionLevel,
		CanRead:         rule.CanRead,
		CanWrite:        rule.CanWrite,
		CanCreate:       rule.CanCreate,
		CanDelete:       rule.CanDelete,
		CanSubmit:       rule.CanSubmit,
		CanCancel:       rule.CanCancel,
		CanAmend:        rule.CanAmend,
		CanPrint:        rule.CanPrint,
		CanEmail:        rule.CanEmail,
		CanImport:       rule.CanImport,
		CanExport:       rule.CanExport,
		CanShare:        rule.CanShare,
		CanReport:       rule.CanReport,
		OnlyIfCreator:   rule.OnlyIfCreator,
		CreatedAt:       rule.CreatedAt,
		UpdatedAt:       rule.UpdatedAt,
	}, nil
}

// UpdatePermissionRule 更新权限规则
func (s *FrappePermissionService) UpdatePermissionRule(ctx context.Context, req *UpdatePermissionRuleRequest) (*PermissionRuleInfo, error) {
	if req.ID <= 0 {
		return nil, errors.BadRequest("INVALID_PARAMETER", "权限规则ID无效")
	}

	rule := &biz.PermissionRule{
		ID:              req.ID,
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
		UpdatedAt:       time.Now(),
	}

	updatedRule, err := s.permissionUc.UpdatePermissionRule(ctx, rule)
	if err != nil {
		s.log.Errorf("Failed to update permission rule: %v", err)
		return nil, errors.InternalServer("INTERNAL_ERROR", "权限规则更新失败")
	}

	s.log.Infof("Permission rule updated successfully: %d", req.ID)

	return &PermissionRuleInfo{
		ID:              updatedRule.ID,
		RoleID:          updatedRule.RoleID,
		DocType:         updatedRule.DocType,
		PermissionLevel: updatedRule.PermissionLevel,
		CanRead:         updatedRule.CanRead,
		CanWrite:        updatedRule.CanWrite,
		CanCreate:       updatedRule.CanCreate,
		CanDelete:       updatedRule.CanDelete,
		CanSubmit:       updatedRule.CanSubmit,
		CanCancel:       updatedRule.CanCancel,
		CanAmend:        updatedRule.CanAmend,
		CanPrint:        updatedRule.CanPrint,
		CanEmail:        updatedRule.CanEmail,
		CanImport:       updatedRule.CanImport,
		CanExport:       updatedRule.CanExport,
		CanShare:        updatedRule.CanShare,
		CanReport:       updatedRule.CanReport,
		OnlyIfCreator:   updatedRule.OnlyIfCreator,
		CreatedAt:       updatedRule.CreatedAt,
		UpdatedAt:       updatedRule.UpdatedAt,
	}, nil
}

// DeletePermissionRule 删除权限规则
func (s *FrappePermissionService) DeletePermissionRule(ctx context.Context, id int64) error {
	if id <= 0 {
		return errors.BadRequest("INVALID_PARAMETER", "权限规则ID无效")
	}

	err := s.permissionUc.DeletePermissionRule(ctx, id)
	if err != nil {
		s.log.Errorf("Failed to delete permission rule: %v", err)
		return errors.InternalServer("INTERNAL_ERROR", "权限规则删除失败")
	}

	s.log.Infof("Permission rule deleted successfully: %d", id)
	return nil
}

// ListPermissionRulesRequest 列出权限规则请求
type ListPermissionRulesRequest struct {
	RoleID  int64  `json:"role_id"`
	DocType string `json:"doc_type"`
}

// ListPermissionRulesResponse 列出权限规则响应
type ListPermissionRulesResponse struct {
	Rules []*PermissionRuleInfo `json:"rules"`
	Total int                   `json:"total"`
}

// ListPermissionRules 列出权限规则
func (s *FrappePermissionService) ListPermissionRules(ctx context.Context, req *ListPermissionRulesRequest) (*ListPermissionRulesResponse, error) {
	rules, err := s.permissionUc.ListPermissionRules(ctx, req.RoleID, req.DocType)
	if err != nil {
		s.log.Errorf("Failed to list permission rules: %v", err)
		return nil, errors.InternalServer("INTERNAL_ERROR", "获取权限规则列表失败")
	}

	var ruleInfos []*PermissionRuleInfo
	for _, rule := range rules {
		ruleInfos = append(ruleInfos, &PermissionRuleInfo{
			ID:                    rule.ID,
			RoleID:                rule.RoleID,
			DocType:               rule.DocType,
			PermissionLevel:       rule.PermissionLevel,
			CanRead:               rule.CanRead,
			CanWrite:              rule.CanWrite,
			CanCreate:             rule.CanCreate,
			CanDelete:             rule.CanDelete,
			CanSubmit:             rule.CanSubmit,
			CanCancel:             rule.CanCancel,
			CanAmend:              rule.CanAmend,
			CanPrint:              rule.CanPrint,
			CanEmail:              rule.CanEmail,
			CanImport:             rule.CanImport,
			CanExport:             rule.CanExport,
			CanShare:              rule.CanShare,
			CanReport:             rule.CanReport,
			CanSetUserPermissions: rule.CanSetUserPermissions,
			OnlyIfCreator:         rule.OnlyIfCreator,
			CreatedAt:             rule.CreatedAt,
			UpdatedAt:             rule.UpdatedAt,
		})
	}

	return &ListPermissionRulesResponse{
		Rules: ruleInfos,
		Total: len(ruleInfos),
	}, nil
}

// BatchCreatePermissionRulesRequest 批量创建权限规则请求
type BatchCreatePermissionRulesRequest struct {
	Rules []*CreatePermissionRuleRequest `json:"rules" validate:"required,min=1"`
}

// BatchCreatePermissionRules 批量创建权限规则
func (s *FrappePermissionService) BatchCreatePermissionRules(ctx context.Context, req *BatchCreatePermissionRulesRequest) error {
	if len(req.Rules) == 0 {
		return errors.BadRequest("INVALID_PARAMETER", "权限规则列表不能为空")
	}

	var rules []*biz.PermissionRule
	for _, ruleReq := range req.Rules {
		rule := &biz.PermissionRule{
			RoleID:          ruleReq.RoleID,
			DocType:         ruleReq.DocType,
			PermissionLevel: ruleReq.PermissionLevel,
			CanRead:         ruleReq.CanRead,
			CanWrite:        ruleReq.CanWrite,
			CanCreate:       ruleReq.CanCreate,
			CanDelete:       ruleReq.CanDelete,
			CanSubmit:       ruleReq.CanSubmit,
			CanCancel:       ruleReq.CanCancel,
			CanAmend:        ruleReq.CanAmend,
			CanPrint:        ruleReq.CanPrint,
			CanEmail:        ruleReq.CanEmail,
			CanImport:       ruleReq.CanImport,
			CanExport:       ruleReq.CanExport,
			CanShare:        ruleReq.CanShare,
			CanReport:       ruleReq.CanReport,
			OnlyIfCreator:   ruleReq.OnlyIfCreator,
			CreatedAt:       time.Now(),
			UpdatedAt:       time.Now(),
		}
		rules = append(rules, rule)
	}

	err := s.permissionUc.BatchCreatePermissionRules(ctx, rules)
	if err != nil {
		s.log.Errorf("Failed to batch create permission rules: %v", err)
		return errors.InternalServer("INTERNAL_ERROR", "批量创建权限规则失败")
	}

	s.log.Infof("Batch created %d permission rules successfully", len(rules))
	return nil
}

// ================================================================
// 权限检查和查询 API
// ================================================================

// CheckPermissionRequest 权限检查请求
type CheckPermissionRequest struct {
	UserID          int64  `json:"user_id" validate:"required"`
	DocType         string `json:"doc_type" validate:"required"`
	Permission      string `json:"permission" validate:"required"`
	PermissionLevel int    `json:"permission_level" validate:"min=0,max=9"`
}

// CheckPermissionResponse 权限检查响应
type CheckPermissionResponse struct {
	HasPermission bool `json:"has_permission"`
}

// CheckPermission 检查权限
func (s *FrappePermissionService) CheckPermission(ctx context.Context, req *CheckPermissionRequest) (*CheckPermissionResponse, error) {
	hasPermission, err := s.permissionUc.CheckPermission(ctx, req.UserID, req.DocType, req.Permission, req.PermissionLevel)
	if err != nil {
		s.log.Errorf("Failed to check permission: %v", err)
		return nil, errors.InternalServer("INTERNAL_ERROR", "权限检查失败")
	}

	return &CheckPermissionResponse{
		HasPermission: hasPermission,
	}, nil
}

// CheckDocumentPermissionRequest 文档权限检查请求
type CheckDocumentPermissionRequest struct {
	UserID     int64  `json:"user_id" validate:"required"`
	DocType    string `json:"doc_type" validate:"required"`
	Permission string `json:"permission" validate:"required"`
	DocID      *int64 `json:"doc_id,omitempty"`
}

// CheckDocumentPermission 检查文档权限
func (s *FrappePermissionService) CheckDocumentPermission(ctx context.Context, req *CheckDocumentPermissionRequest) (*CheckPermissionResponse, error) {
	permissionReq := &biz.PermissionCheckRequest{
		UserID:          req.UserID,
		DocType:         req.DocType,
		Permission:      req.Permission,
		PermissionLevel: 0, // 默认文档级权限
		DocID:           req.DocID,
	}

	// 验证权限检查请求
	if err := permissionReq.Validate(); err != nil {
		s.log.Errorf("PermissionCheckRequest validation failed: %v", err)
		return nil, errors.BadRequest("INVALID_DATA", fmt.Sprintf("权限检查请求验证失败: %v", err))
	}

	permissionResponse, err := s.permissionUc.CheckDocumentPermission(ctx, permissionReq)
	if err != nil {
		s.log.Errorf("Failed to check document permission: %v", err)
		return nil, errors.InternalServer("INTERNAL_ERROR", "文档权限检查失败")
	}

	return &CheckPermissionResponse{
		HasPermission: permissionResponse.HasPermission,
	}, nil
}

// GetUserPermissionLevelRequest 获取用户权限级别请求
type GetUserPermissionLevelRequest struct {
	UserID  int64  `json:"user_id" validate:"required"`
	DocType string `json:"doc_type" validate:"required"`
}

// GetUserPermissionLevelResponse 获取用户权限级别响应
type GetUserPermissionLevelResponse struct {
	PermissionLevel int `json:"permission_level"`
}

// GetUserPermissionLevel 获取用户权限级别
func (s *FrappePermissionService) GetUserPermissionLevel(ctx context.Context, req *GetUserPermissionLevelRequest) (*GetUserPermissionLevelResponse, error) {
	level, err := s.permissionUc.GetUserPermissionLevel(ctx, req.UserID, req.DocType)
	if err != nil {
		s.log.Errorf("Failed to get user permission level: %v", err)
		return nil, errors.InternalServer("INTERNAL_ERROR", "获取用户权限级别失败")
	}

	return &GetUserPermissionLevelResponse{
		PermissionLevel: level,
	}, nil
}

// GetUserEnhancedPermissionsRequest 获取用户增强权限请求
type GetUserEnhancedPermissionsRequest struct {
	UserID  int64  `json:"user_id" validate:"required"`
	DocType string `json:"doc_type" validate:"required"`
}

// GetUserEnhancedPermissionsResponse 获取用户增强权限响应
type GetUserEnhancedPermissionsResponse struct {
	Permissions []*biz.EnhancedUserPermission `json:"permissions"`
}

// GetUserEnhancedPermissions 获取用户增强权限
func (s *FrappePermissionService) GetUserEnhancedPermissions(ctx context.Context, req *GetUserEnhancedPermissionsRequest) (*GetUserEnhancedPermissionsResponse, error) {
	permissions, err := s.permissionUc.GetUserEnhancedPermissions(ctx, req.UserID, req.DocType)
	if err != nil {
		s.log.Errorf("Failed to get user enhanced permissions: %v", err)
		return nil, errors.InternalServer("INTERNAL_ERROR", "获取用户增强权限失败")
	}

	return &GetUserEnhancedPermissionsResponse{
		Permissions: permissions,
	}, nil
}

// GetAccessibleFieldsRequest 获取可访问字段请求
type GetAccessibleFieldsRequest struct {
	UserID  int64  `json:"user_id" validate:"required"`
	DocType string `json:"doc_type" validate:"required"`
}

// GetAccessibleFieldsResponse 获取可访问字段响应
type GetAccessibleFieldsResponse struct {
	Fields []*biz.AccessibleField `json:"fields"`
}

// GetAccessibleFields 获取可访问字段
func (s *FrappePermissionService) GetAccessibleFields(ctx context.Context, req *GetAccessibleFieldsRequest) (*GetAccessibleFieldsResponse, error) {
	fieldReq := &biz.FieldPermissionRequest{
		UserID:  req.UserID,
		DocType: req.DocType,
	}

	fields, err := s.permissionUc.GetAccessibleFields(ctx, fieldReq)
	if err != nil {
		s.log.Errorf("Failed to get accessible fields: %v", err)
		return nil, errors.InternalServer("INTERNAL_ERROR", "获取可访问字段失败")
	}

	return &GetAccessibleFieldsResponse{
		Fields: fields.Fields,
	}, nil
}

// FilterDocumentsByPermissionRequest 按权限过滤文档请求
type FilterDocumentsByPermissionRequest struct {
	UserID       int64                    `json:"user_id" validate:"required"`
	DocumentType string                   `json:"document_type" validate:"required"`
	Documents    []map[string]interface{} `json:"documents" validate:"required"`
}

// FilterDocumentsByPermissionResponse 按权限过滤文档响应
type FilterDocumentsByPermissionResponse struct {
	FilteredDocuments []map[string]interface{} `json:"filtered_documents"`
}

// FilterDocumentsByPermission 按权限过滤文档
func (s *FrappePermissionService) FilterDocumentsByPermission(ctx context.Context, req *FilterDocumentsByPermissionRequest) (*FilterDocumentsByPermissionResponse, error) {
	filteredDocs, err := s.permissionUc.FilterDocumentsByPermission(ctx, req.UserID, req.DocumentType, req.Documents)
	if err != nil {
		s.log.Errorf("Failed to filter documents by permission: %v", err)
		return nil, errors.InternalServer("INTERNAL_ERROR", "按权限过滤文档失败")
	}

	return &FilterDocumentsByPermissionResponse{
		FilteredDocuments: filteredDocs,
	}, nil
}

// GetUserRolesRequest 获取用户角色请求
type GetUserRolesRequest struct {
	UserID int64 `json:"user_id" validate:"required"`
}

// GetUserRolesResponse 获取用户角色响应
type GetUserRolesResponse struct {
	Roles []string `json:"roles"`
}

// GetUserRoles 获取用户角色
func (s *FrappePermissionService) GetUserRoles(ctx context.Context, req *GetUserRolesRequest) (*GetUserRolesResponse, error) {
	roles, err := s.permissionUc.GetUserRoles(ctx, req.UserID)
	if err != nil {
		s.log.Errorf("Failed to get user roles: %v", err)
		return nil, errors.InternalServer("INTERNAL_ERROR", "获取用户角色失败")
	}

	return &GetUserRolesResponse{
		Roles: roles,
	}, nil
}

// UserPermission APIs

// CreateUserPermissionRequest 创建用户权限请求
type CreateUserPermissionRequest struct {
	UserID          int64   `json:"user_id" validate:"required"`
	DocType         string  `json:"doc_type" validate:"required"`
	DocumentName    string  `json:"document_name" validate:"required"`
	Condition       string  `json:"condition"`
	ApplicableFor   *string `json:"applicable_for,omitempty"`
	HideDescendants bool    `json:"hide_descendants"`
	IsDefault       bool    `json:"is_default"`
}

// UserPermissionInfo 用户权限信息
type UserPermissionInfo struct {
	ID           int64     `json:"id"`
	UserID       int64     `json:"user_id"`
	DocType      string    `json:"doc_type"`
	DocumentName string    `json:"document_name"`
	Condition    string    `json:"condition"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// ListUserPermissionsRequest 获取用户权限列表请求
type ListUserPermissionsRequest struct {
	UserID  int64  `form:"user_id"`
	DocType string `form:"doc_type"`
	Page    int32  `form:"page"`
	Size    int32  `form:"size"`
}

// ListUserPermissionsResponse 获取用户权限列表响应
type ListUserPermissionsResponse struct {
	UserPermissions []*UserPermissionInfo `json:"user_permissions"`
	Total           int32                 `json:"total"`
	Page            int32                 `json:"page"`
	Size            int32                 `json:"size"`
}

// CreateUserPermission 创建用户权限
func (s *FrappePermissionService) CreateUserPermission(ctx context.Context, req *CreateUserPermissionRequest) (*UserPermissionInfo, error) {
	// 检查权限
	currentUser := middleware.GetCurrentUser(ctx)
	if !currentUser.HasAnyRole("SUPER_ADMIN", "ADMIN", "PERMISSION_MANAGER") {
		return nil, errors.Forbidden("PERMISSION_DENIED", "无权限创建用户权限")
	}

	// 创建业务层请求对象
	bizReq := &biz.CreateUserPermissionRequest{
		UserID:          req.UserID,
		DocType:         req.DocType,
		DocName:         req.DocumentName,
		Value:           req.Condition,
		ApplicableFor:   req.ApplicableFor,
		HideDescendants: req.HideDescendants,
		IsDefault:       req.IsDefault,
	}

	// 验证请求数据
	if err := bizReq.Validate(); err != nil {
		s.log.Errorf("CreateUserPermissionRequest validation failed: %v", err)
		return nil, errors.BadRequest("INVALID_DATA", fmt.Sprintf("请求数据验证失败: %v", err))
	}

	// 创建用户权限
	userPermission := &biz.UserPermission{
		UserID:          req.UserID,
		DocType:         req.DocType,
		DocName:         req.DocumentName,
		Value:           req.Condition,
		ApplicableFor:   req.ApplicableFor,
		HideDescendants: req.HideDescendants,
		IsDefault:       req.IsDefault,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	// 验证业务实体数据
	if err := userPermission.Validate(); err != nil {
		s.log.Errorf("UserPermission validation failed: %v", err)
		return nil, errors.BadRequest("INVALID_DATA", fmt.Sprintf("用户权限验证失败: %v", err))
	}

	createdUserPermission, err := s.permissionUc.CreateUserPermission(ctx, userPermission)
	if err != nil {
		s.log.Errorf("Failed to create user permission: %v", err)
		return nil, errors.InternalServer("INTERNAL_ERROR", "用户权限创建失败")
	}

	s.log.Infof("User permission created successfully for user %d", req.UserID)

	return s.convertToUserPermissionInfo(createdUserPermission), nil
}

// GetUserPermission 获取用户权限详情
func (s *FrappePermissionService) GetUserPermission(ctx context.Context, permissionID int64) (*UserPermissionInfo, error) {
	// 检查权限
	currentUser := middleware.GetCurrentUser(ctx)
	if !currentUser.HasAnyRole("SUPER_ADMIN", "ADMIN", "PERMISSION_MANAGER") {
		return nil, errors.Forbidden("PERMISSION_DENIED", "无权限查看用户权限")
	}

	// 获取用户权限
	userPermission, err := s.permissionUc.GetUserPermission(ctx, permissionID)
	if err != nil {
		return nil, errors.NotFound("USER_PERMISSION_NOT_FOUND", "用户权限不存在")
	}

	return s.convertToUserPermissionInfo(userPermission), nil
}

// UpdateUserPermission 更新用户权限
func (s *FrappePermissionService) UpdateUserPermission(ctx context.Context, permissionID int64, req *CreateUserPermissionRequest) (*UserPermissionInfo, error) {
	// 检查权限
	currentUser := middleware.GetCurrentUser(ctx)
	if !currentUser.HasAnyRole("SUPER_ADMIN", "ADMIN", "PERMISSION_MANAGER") {
		return nil, errors.Forbidden("PERMISSION_DENIED", "无权限修改用户权限")
	}

	// 获取原用户权限
	userPermission, err := s.permissionUc.GetUserPermission(ctx, permissionID)
	if err != nil {
		return nil, errors.NotFound("USER_PERMISSION_NOT_FOUND", "用户权限不存在")
	}

	// 更新用户权限信息
	userPermission.UserID = req.UserID
	userPermission.DocType = req.DocType
	userPermission.DocName = req.DocumentName
	userPermission.Value = req.Condition
	userPermission.UpdatedAt = time.Now()

	updatedUserPermission, err := s.permissionUc.UpdateUserPermission(ctx, userPermission)
	if err != nil {
		s.log.Errorf("Failed to update user permission: %v", err)
		return nil, errors.InternalServer("INTERNAL_ERROR", "用户权限更新失败")
	}

	s.log.Infof("User permission updated successfully: %d", permissionID)

	return s.convertToUserPermissionInfo(updatedUserPermission), nil
}

// DeleteUserPermission 删除用户权限
func (s *FrappePermissionService) DeleteUserPermission(ctx context.Context, permissionID int64) error {
	// 检查权限
	currentUser := middleware.GetCurrentUser(ctx)
	if !currentUser.HasAnyRole("SUPER_ADMIN", "ADMIN", "PERMISSION_MANAGER") {
		return errors.Forbidden("PERMISSION_DENIED", "无权限删除用户权限")
	}

	// 删除用户权限
	if err := s.permissionUc.DeleteUserPermission(ctx, permissionID); err != nil {
		s.log.Errorf("Failed to delete user permission: %v", err)
		return errors.InternalServer("INTERNAL_ERROR", "用户权限删除失败")
	}

	s.log.Infof("User permission deleted successfully: %d", permissionID)
	return nil
}

// ListUserPermissions 获取用户权限列表
func (s *FrappePermissionService) ListUserPermissions(ctx context.Context, req *ListUserPermissionsRequest) (*ListUserPermissionsResponse, error) {
	// 检查权限
	currentUser := middleware.GetCurrentUser(ctx)
	if !currentUser.HasAnyRole("SUPER_ADMIN", "ADMIN", "PERMISSION_MANAGER") {
		return nil, errors.Forbidden("PERMISSION_DENIED", "无权限查看用户权限列表")
	}

	// 设置默认值
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Size <= 0 || req.Size > 100 {
		req.Size = 20
	}

	// 获取用户权限列表
	userPermissions, err := s.permissionUc.ListUserPermissions(ctx, req.UserID, req.DocType, req.Page, req.Size)
	if err != nil {
		s.log.Errorf("Failed to list user permissions: %v", err)
		return nil, errors.InternalServer("INTERNAL_ERROR", "用户权限列表获取失败")
	}

	// 获取总数
	total, err := s.permissionUc.GetUserPermissionsCount(ctx, req.UserID, req.DocType)
	if err != nil {
		s.log.Errorf("Failed to get user permissions count: %v", err)
		return nil, errors.InternalServer("INTERNAL_ERROR", "用户权限总数获取失败")
	}

	// 转换为响应格式
	infos := make([]*UserPermissionInfo, len(userPermissions))
	for i, userPermission := range userPermissions {
		infos[i] = s.convertToUserPermissionInfo(userPermission)
	}

	return &ListUserPermissionsResponse{
		UserPermissions: infos,
		Total:           total,
		Page:            req.Page,
		Size:            req.Size,
	}, nil
}

// FieldPermissionLevel APIs

// CreateFieldPermissionLevelRequest 创建字段权限级别请求
type CreateFieldPermissionLevelRequest struct {
	DocType         string `json:"doc_type" validate:"required"`
	FieldName       string `json:"field_name" validate:"required"`
	PermissionLevel int    `json:"permission_level" validate:"required,min=0,max=9"`
}

// FieldPermissionLevelInfo 字段权限级别信息
type FieldPermissionLevelInfo struct {
	ID              int64     `json:"id"`
	DocType         string    `json:"doc_type"`
	FieldName       string    `json:"field_name"`
	PermissionLevel int       `json:"permission_level"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// ListFieldPermissionLevelsRequest 获取字段权限级别列表请求
type ListFieldPermissionLevelsRequest struct {
	DocType string `form:"doc_type"`
	Page    int32  `form:"page"`
	Size    int32  `form:"size"`
}

// ListFieldPermissionLevelsResponse 获取字段权限级别列表响应
type ListFieldPermissionLevelsResponse struct {
	FieldPermissionLevels []*FieldPermissionLevelInfo `json:"field_permission_levels"`
	Total                 int32                       `json:"total"`
	Page                  int32                       `json:"page"`
	Size                  int32                       `json:"size"`
}

// CreateFieldPermissionLevel 创建字段权限级别
func (s *FrappePermissionService) CreateFieldPermissionLevel(ctx context.Context, req *CreateFieldPermissionLevelRequest) (*FieldPermissionLevelInfo, error) {
	// 检查权限
	currentUser := middleware.GetCurrentUser(ctx)
	if !currentUser.HasAnyRole("SUPER_ADMIN", "ADMIN") {
		return nil, errors.Forbidden("PERMISSION_DENIED", "无权限创建字段权限级别")
	}

	// 验证权限级别
	if req.PermissionLevel < 0 || req.PermissionLevel > 9 {
		return nil, errors.BadRequest("INVALID_PERMISSION_LEVEL", "权限级别必须在0-9之间")
	}

	// 创建字段权限级别
	fieldPermissionLevel := &biz.FieldPermissionLevel{
		DocType:         req.DocType,
		FieldName:       req.FieldName,
		PermissionLevel: req.PermissionLevel,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	createdFieldPermissionLevel, err := s.permissionUc.CreateFieldPermissionLevel(ctx, fieldPermissionLevel)
	if err != nil {
		s.log.Errorf("Failed to create field permission level: %v", err)
		return nil, errors.InternalServer("INTERNAL_ERROR", "字段权限级别创建失败")
	}

	s.log.Infof("Field permission level created successfully: %s.%s", req.DocType, req.FieldName)

	return s.convertToFieldPermissionLevelInfo(createdFieldPermissionLevel), nil
}

// GetFieldPermissionLevel 获取字段权限级别详情
func (s *FrappePermissionService) GetFieldPermissionLevel(ctx context.Context, levelID int64) (*FieldPermissionLevelInfo, error) {
	// 检查权限
	currentUser := middleware.GetCurrentUser(ctx)
	if !currentUser.HasAnyRole("SUPER_ADMIN", "ADMIN", "PERMISSION_MANAGER") {
		return nil, errors.Forbidden("PERMISSION_DENIED", "无权限查看字段权限级别")
	}

	// 获取字段权限级别
	fieldPermissionLevel, err := s.permissionUc.GetFieldPermissionLevel(ctx, levelID)
	if err != nil {
		return nil, errors.NotFound("FIELD_PERMISSION_LEVEL_NOT_FOUND", "字段权限级别不存在")
	}

	return s.convertToFieldPermissionLevelInfo(fieldPermissionLevel), nil
}

// UpdateFieldPermissionLevel 更新字段权限级别
func (s *FrappePermissionService) UpdateFieldPermissionLevel(ctx context.Context, levelID int64, req *CreateFieldPermissionLevelRequest) (*FieldPermissionLevelInfo, error) {
	// 检查权限
	currentUser := middleware.GetCurrentUser(ctx)
	if !currentUser.HasAnyRole("SUPER_ADMIN", "ADMIN") {
		return nil, errors.Forbidden("PERMISSION_DENIED", "无权限修改字段权限级别")
	}

	// 验证权限级别
	if req.PermissionLevel < 0 || req.PermissionLevel > 9 {
		return nil, errors.BadRequest("INVALID_PERMISSION_LEVEL", "权限级别必须在0-9之间")
	}

	// 获取原字段权限级别
	fieldPermissionLevel, err := s.permissionUc.GetFieldPermissionLevel(ctx, levelID)
	if err != nil {
		return nil, errors.NotFound("FIELD_PERMISSION_LEVEL_NOT_FOUND", "字段权限级别不存在")
	}

	// 更新字段权限级别信息
	fieldPermissionLevel.DocType = req.DocType
	fieldPermissionLevel.FieldName = req.FieldName
	fieldPermissionLevel.PermissionLevel = req.PermissionLevel
	fieldPermissionLevel.UpdatedAt = time.Now()

	updatedFieldPermissionLevel, err := s.permissionUc.UpdateFieldPermissionLevel(ctx, fieldPermissionLevel)
	if err != nil {
		s.log.Errorf("Failed to update field permission level: %v", err)
		return nil, errors.InternalServer("INTERNAL_ERROR", "字段权限级别更新失败")
	}

	s.log.Infof("Field permission level updated successfully: %d", levelID)

	return s.convertToFieldPermissionLevelInfo(updatedFieldPermissionLevel), nil
}

// DeleteFieldPermissionLevel 删除字段权限级别
func (s *FrappePermissionService) DeleteFieldPermissionLevel(ctx context.Context, levelID int64) error {
	// 检查权限
	currentUser := middleware.GetCurrentUser(ctx)
	if !currentUser.HasAnyRole("SUPER_ADMIN", "ADMIN") {
		return errors.Forbidden("PERMISSION_DENIED", "无权限删除字段权限级别")
	}

	// 删除字段权限级别
	if err := s.permissionUc.DeleteFieldPermissionLevel(ctx, levelID); err != nil {
		s.log.Errorf("Failed to delete field permission level: %v", err)
		return errors.InternalServer("INTERNAL_ERROR", "字段权限级别删除失败")
	}

	s.log.Infof("Field permission level deleted successfully: %d", levelID)
	return nil
}

// ListFieldPermissionLevels 获取字段权限级别列表
func (s *FrappePermissionService) ListFieldPermissionLevels(ctx context.Context, req *ListFieldPermissionLevelsRequest) (*ListFieldPermissionLevelsResponse, error) {
	// 检查权限
	currentUser := middleware.GetCurrentUser(ctx)
	if !currentUser.HasAnyRole("SUPER_ADMIN", "ADMIN", "PERMISSION_MANAGER") {
		return nil, errors.Forbidden("PERMISSION_DENIED", "无权限查看字段权限级别列表")
	}

	// 设置默认值
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Size <= 0 || req.Size > 100 {
		req.Size = 20
	}

	// 获取字段权限级别列表
	fieldPermissionLevels, err := s.permissionUc.ListFieldPermissionLevels(ctx, req.DocType, req.Page, req.Size)
	if err != nil {
		s.log.Errorf("Failed to list field permission levels: %v", err)
		return nil, errors.InternalServer("INTERNAL_ERROR", "字段权限级别列表获取失败")
	}

	// 获取总数
	total, err := s.permissionUc.GetFieldPermissionLevelsCount(ctx, req.DocType)
	if err != nil {
		s.log.Errorf("Failed to get field permission levels count: %v", err)
		return nil, errors.InternalServer("INTERNAL_ERROR", "字段权限级别总数获取失败")
	}

	// 转换为响应格式
	infos := make([]*FieldPermissionLevelInfo, len(fieldPermissionLevels))
	for i, fieldPermissionLevel := range fieldPermissionLevels {
		infos[i] = s.convertToFieldPermissionLevelInfo(fieldPermissionLevel)
	}

	return &ListFieldPermissionLevelsResponse{
		FieldPermissionLevels: infos,
		Total:                 total,
		Page:                  req.Page,
		Size:                  req.Size,
	}, nil
}

// DocumentWorkflowState APIs

// CreateDocumentWorkflowStateRequest 创建文档工作流状态请求
type CreateDocumentWorkflowStateRequest struct {
	DocType       string  `json:"doc_type" validate:"required"`
	DocID         int64   `json:"doc_id" validate:"required"`
	DocName       *string `json:"doc_name,omitempty"`
	WorkflowState string  `json:"workflow_state" validate:"required"`
	DocStatus     int     `json:"doc_status"`
}

// DocumentWorkflowStateInfo 文档工作流状态信息
type DocumentWorkflowStateInfo struct {
	ID            int64      `json:"id"`
	DocType       string     `json:"doc_type"`
	DocID         int64      `json:"doc_id"`
	DocName       *string    `json:"doc_name,omitempty"`
	WorkflowState string     `json:"workflow_state"`
	DocStatus     int        `json:"doc_status"`
	SubmittedAt   *time.Time `json:"submitted_at,omitempty"`
	SubmittedBy   *int64     `json:"submitted_by,omitempty"`
	CancelledAt   *time.Time `json:"cancelled_at,omitempty"`
	CancelledBy   *int64     `json:"cancelled_by,omitempty"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}

// ListDocumentWorkflowStatesRequest 获取文档工作流状态列表请求
type ListDocumentWorkflowStatesRequest struct {
	DocType       string `form:"doc_type"`
	DocID         int64  `form:"doc_id"`
	WorkflowState string `form:"workflow_state"`
	DocStatus     int    `form:"doc_status"`
	Page          int32  `form:"page"`
	Size          int32  `form:"size"`
}

// ListDocumentWorkflowStatesResponse 获取文档工作流状态列表响应
type ListDocumentWorkflowStatesResponse struct {
	DocumentWorkflowStates []*DocumentWorkflowStateInfo `json:"document_workflow_states"`
	Total                  int32                        `json:"total"`
	Page                   int32                        `json:"page"`
	Size                   int32                        `json:"size"`
}

// CreateDocumentWorkflowState 创建文档工作流状态
func (s *FrappePermissionService) CreateDocumentWorkflowState(ctx context.Context, req *CreateDocumentWorkflowStateRequest) (*DocumentWorkflowStateInfo, error) {
	// 检查权限
	currentUser := middleware.GetCurrentUser(ctx)
	if !currentUser.HasAnyRole("SUPER_ADMIN", "ADMIN") {
		return nil, errors.Forbidden("PERMISSION_DENIED", "无权限创建文档工作流状态")
	}

	// 验证请求
	if req.DocID <= 0 {
		return nil, errors.BadRequest("INVALID_DOC_ID", "无效的文档ID")
	}

	// 创建文档工作流状态
	documentWorkflowState := &biz.DocumentWorkflowState{
		DocType:       req.DocType,
		DocID:         req.DocID,
		DocName:       req.DocName,
		WorkflowState: req.WorkflowState,
		DocStatus:     req.DocStatus,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	createdDocumentWorkflowState, err := s.permissionUc.CreateDocumentWorkflowState(ctx, documentWorkflowState)
	if err != nil {
		s.log.Errorf("Failed to create document workflow state: %v", err)
		return nil, errors.InternalServer("INTERNAL_ERROR", "文档工作流状态创建失败")
	}

	s.log.Infof("Document workflow state created successfully: %s.%d", req.DocType, req.DocID)

	return s.convertToDocumentWorkflowStateInfo(createdDocumentWorkflowState), nil
}

// GetDocumentWorkflowState 获取文档工作流状态详情
func (s *FrappePermissionService) GetDocumentWorkflowState(ctx context.Context, stateID int64) (*DocumentWorkflowStateInfo, error) {
	// 检查权限
	currentUser := middleware.GetCurrentUser(ctx)
	if !currentUser.HasAnyRole("SUPER_ADMIN", "ADMIN", "PERMISSION_MANAGER") {
		return nil, errors.Forbidden("PERMISSION_DENIED", "无权限查看文档工作流状态")
	}

	// 获取文档工作流状态
	documentWorkflowState, err := s.permissionUc.GetDocumentWorkflowState(ctx, stateID)
	if err != nil {
		return nil, errors.NotFound("DOCUMENT_WORKFLOW_STATE_NOT_FOUND", "文档工作流状态不存在")
	}

	return s.convertToDocumentWorkflowStateInfo(documentWorkflowState), nil
}

// UpdateDocumentWorkflowState 更新文档工作流状态
func (s *FrappePermissionService) UpdateDocumentWorkflowState(ctx context.Context, stateID int64, req *CreateDocumentWorkflowStateRequest) (*DocumentWorkflowStateInfo, error) {
	// 检查权限
	currentUser := middleware.GetCurrentUser(ctx)
	if !currentUser.HasAnyRole("SUPER_ADMIN", "ADMIN") {
		return nil, errors.Forbidden("PERMISSION_DENIED", "无权限修改文档工作流状态")
	}

	// 验证请求
	if req.DocID <= 0 {
		return nil, errors.BadRequest("INVALID_DOC_ID", "无效的文档ID")
	}

	// 获取原文档工作流状态
	documentWorkflowState, err := s.permissionUc.GetDocumentWorkflowState(ctx, stateID)
	if err != nil {
		return nil, errors.NotFound("DOCUMENT_WORKFLOW_STATE_NOT_FOUND", "文档工作流状态不存在")
	}

	// 更新文档工作流状态信息
	documentWorkflowState.DocType = req.DocType
	documentWorkflowState.DocID = req.DocID
	documentWorkflowState.DocName = req.DocName
	documentWorkflowState.WorkflowState = req.WorkflowState
	documentWorkflowState.DocStatus = req.DocStatus
	documentWorkflowState.UpdatedAt = time.Now()

	updatedDocumentWorkflowState, err := s.permissionUc.UpdateDocumentWorkflowState(ctx, documentWorkflowState)
	if err != nil {
		s.log.Errorf("Failed to update document workflow state: %v", err)
		return nil, errors.InternalServer("INTERNAL_ERROR", "文档工作流状态更新失败")
	}

	s.log.Infof("Document workflow state updated successfully: %d", stateID)

	return s.convertToDocumentWorkflowStateInfo(updatedDocumentWorkflowState), nil
}

// DeleteDocumentWorkflowState 删除文档工作流状态
func (s *FrappePermissionService) DeleteDocumentWorkflowState(ctx context.Context, stateID int64) error {
	// 检查权限
	currentUser := middleware.GetCurrentUser(ctx)
	if !currentUser.HasAnyRole("SUPER_ADMIN", "ADMIN") {
		return errors.Forbidden("PERMISSION_DENIED", "无权限删除文档工作流状态")
	}

	// 删除文档工作流状态
	if err := s.permissionUc.DeleteDocumentWorkflowState(ctx, stateID); err != nil {
		s.log.Errorf("Failed to delete document workflow state: %v", err)
		return errors.InternalServer("INTERNAL_ERROR", "文档工作流状态删除失败")
	}

	s.log.Infof("Document workflow state deleted successfully: %d", stateID)
	return nil
}

// ListDocumentWorkflowStates 获取文档工作流状态列表
func (s *FrappePermissionService) ListDocumentWorkflowStates(ctx context.Context, req *ListDocumentWorkflowStatesRequest) (*ListDocumentWorkflowStatesResponse, error) {
	// 检查权限
	currentUser := middleware.GetCurrentUser(ctx)
	if !currentUser.HasAnyRole("SUPER_ADMIN", "ADMIN", "PERMISSION_MANAGER") {
		return nil, errors.Forbidden("PERMISSION_DENIED", "无权限查看文档工作流状态列表")
	}

	// 设置默认值
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Size <= 0 || req.Size > 100 {
		req.Size = 20
	}

	// 获取文档工作流状态列表
	// Note: The method signature needs to be updated to match the new parameters
	documentWorkflowStates, err := s.permissionUc.ListDocumentWorkflowStates(ctx, req.DocType, "", req.WorkflowState, req.DocID, req.Page, req.Size)
	if err != nil {
		s.log.Errorf("Failed to list document workflow states: %v", err)
		return nil, errors.InternalServer("INTERNAL_ERROR", "文档工作流状态列表获取失败")
	}

	// 获取总数
	total, err := s.permissionUc.GetDocumentWorkflowStatesCount(ctx, req.DocType, "", req.WorkflowState, req.DocID)
	if err != nil {
		s.log.Errorf("Failed to get document workflow states count: %v", err)
		return nil, errors.InternalServer("INTERNAL_ERROR", "文档工作流状态总数获取失败")
	}

	// 转换为响应格式
	infos := make([]*DocumentWorkflowStateInfo, len(documentWorkflowStates))
	for i, documentWorkflowState := range documentWorkflowStates {
		infos[i] = s.convertToDocumentWorkflowStateInfo(documentWorkflowState)
	}

	return &ListDocumentWorkflowStatesResponse{
		DocumentWorkflowStates: infos,
		Total:                  total,
		Page:                   req.Page,
		Size:                   req.Size,
	}, nil
}

// 辅助转换函数

func (s *FrappePermissionService) convertToUserPermissionInfo(userPermission *biz.UserPermission) *UserPermissionInfo {
	return &UserPermissionInfo{
		ID:           userPermission.ID,
		UserID:       userPermission.UserID,
		DocType:      userPermission.DocType,
		DocumentName: userPermission.DocName,
		Condition:    userPermission.Value,
		CreatedAt:    userPermission.CreatedAt,
		UpdatedAt:    userPermission.UpdatedAt,
	}
}

func (s *FrappePermissionService) convertToFieldPermissionLevelInfo(fieldPermissionLevel *biz.FieldPermissionLevel) *FieldPermissionLevelInfo {
	return &FieldPermissionLevelInfo{
		ID:              fieldPermissionLevel.ID,
		DocType:         fieldPermissionLevel.DocType,
		FieldName:       fieldPermissionLevel.FieldName,
		PermissionLevel: fieldPermissionLevel.PermissionLevel,
		CreatedAt:       fieldPermissionLevel.CreatedAt,
		UpdatedAt:       fieldPermissionLevel.UpdatedAt,
	}
}

func (s *FrappePermissionService) convertToDocumentWorkflowStateInfo(documentWorkflowState *biz.DocumentWorkflowState) *DocumentWorkflowStateInfo {
	return &DocumentWorkflowStateInfo{
		ID:            documentWorkflowState.ID,
		DocType:       documentWorkflowState.DocType,
		DocID:         documentWorkflowState.DocID,
		DocName:       documentWorkflowState.DocName,
		WorkflowState: documentWorkflowState.WorkflowState,
		DocStatus:     documentWorkflowState.DocStatus,
		SubmittedAt:   documentWorkflowState.SubmittedAt,
		SubmittedBy:   documentWorkflowState.SubmittedBy,
		CancelledAt:   documentWorkflowState.CancelledAt,
		CancelledBy:   documentWorkflowState.CancelledBy,
		CreatedAt:     documentWorkflowState.CreatedAt,
		UpdatedAt:     documentWorkflowState.UpdatedAt,
	}
}
