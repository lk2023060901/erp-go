package service

import (
	"context"
	"time"

	"erp-system/internal/biz"
	"erp-system/internal/middleware"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
)

// RoleService 角色服务
type RoleService struct {
	roleUc *biz.RoleUsecase
	log    *log.Helper
}

// NewRoleService 创建角色服务
func NewRoleService(roleUc *biz.RoleUsecase, logger log.Logger) *RoleService {
	return &RoleService{
		roleUc: roleUc,
		log:    log.NewHelper(logger),
	}
}

// CreateRoleRequest 创建角色请求
type CreateRoleRequest struct {
	Name        string  `json:"name" validate:"required,min=2,max=50"`
	Code        string  `json:"code" validate:"required,min=2,max=50"`
	Description string  `json:"description"`
	IsEnabled   bool    `json:"is_enabled"`
	SortOrder   int32   `json:"sort_order"`
	PermissionIDs []int32 `json:"permission_ids"`
}

// UpdateRoleRequest 更新角色请求
type UpdateRoleRequest struct {
	ID          int32   `json:"id" validate:"required"`
	Name        string  `json:"name" validate:"required,min=2,max=50"`
	Description string  `json:"description"`
	IsEnabled   bool    `json:"is_enabled"`
	SortOrder   int32   `json:"sort_order"`
	PermissionIDs []int32 `json:"permission_ids"`
}

// RoleListRequest 角色列表请求
type RoleListRequest struct {
	Page      int32  `json:"page" validate:"min=1"`
	Size      int32  `json:"size" validate:"min=1,max=100"`
	Search    string `json:"search"`
	IsEnabled *bool  `json:"is_enabled"`
	SortField string `json:"sort_field"`
	SortOrder string `json:"sort_order" validate:"omitempty,oneof=asc desc"`
}

// RoleListResponse 角色列表响应
type RoleListResponse struct {
	Roles []*RoleInfo `json:"roles"`
	Total int32       `json:"total"`
	Page  int32       `json:"page"`
	Size  int32       `json:"size"`
}

// AssignPermissionsRequest 分配权限请求
type AssignPermissionsRequest struct {
	RoleID        int32   `json:"role_id" validate:"required"`
	PermissionIDs []int32 `json:"permission_ids" validate:"required"`
}

// RoleInfo 角色信息
type RoleInfo struct {
	ID           int32              `json:"id"`
	Name         string             `json:"name"`
	Code         string             `json:"code"`
	Description  string             `json:"description"`
	IsSystemRole bool               `json:"is_system_role"`
	IsEnabled    bool               `json:"is_enabled"`
	SortOrder    int32              `json:"sort_order"`
	CreatedAt    time.Time          `json:"created_at"`
	UpdatedAt    time.Time          `json:"updated_at"`
	Permissions  []*PermissionInfo  `json:"permissions,omitempty"`
}

// PermissionInfo 权限信息
type PermissionInfo struct {
	ID          int32              `json:"id"`
	ParentID    *int32             `json:"parent_id"`
	Name        string             `json:"name"`
	Code        string             `json:"code"`
	Resource    string             `json:"resource"`
	Action      string             `json:"action"`
	Module      string             `json:"module"`
	Description string             `json:"description"`
	IsMenu      bool               `json:"is_menu"`
	Path        string             `json:"path"`
	SortOrder   int32              `json:"sort_order"`
	CreatedAt   time.Time          `json:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at"`
	Children    []*PermissionInfo  `json:"children,omitempty"`
}

// CreateRole 创建角色
func (s *RoleService) CreateRole(ctx context.Context, req *CreateRoleRequest) (*RoleInfo, error) {
	// 检查权限
	currentUser := middleware.GetCurrentUser(ctx)
	if !currentUser.HasAnyRole("SUPER_ADMIN", "ADMIN") {
		return nil, errors.Forbidden("PERMISSION_DENIED", "无权限创建角色")
	}

	s.log.Infof("Creating role: %s by %s", req.Name, currentUser.Username)

	// 创建角色
	role := &biz.Role{
		Name:        req.Name,
		Code:        req.Code,
		Description: req.Description,
		IsSystemRole: false, // 用户创建的角色都不是系统角色
		IsEnabled:   req.IsEnabled,
		SortOrder:   req.SortOrder,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	createdRole, err := s.roleUc.CreateRole(ctx, role)
	if err != nil {
		s.log.Errorf("Failed to create role: %v", err)
		if err == biz.ErrRoleCodeExists {
			return nil, errors.BadRequest("ROLE_CODE_EXISTS", "角色编码已存在")
		}
		if err == biz.ErrRoleNameExists {
			return nil, errors.BadRequest("ROLE_NAME_EXISTS", "角色名称已存在")
		}
		return nil, errors.InternalServer("INTERNAL_ERROR", "角色创建失败")
	}

	// 分配权限
	if len(req.PermissionIDs) > 0 {
		if err := s.roleUc.AssignPermissions(ctx, createdRole.ID, req.PermissionIDs); err != nil {
			s.log.Errorf("Failed to assign permissions: %v", err)
		}
	}

	// 获取权限
	permissions, _ := s.roleUc.GetRolePermissions(ctx, createdRole.ID)

	s.log.Infof("Role created successfully: %s", req.Name)

	return &RoleInfo{
		ID:           createdRole.ID,
		Name:         createdRole.Name,
		Code:         createdRole.Code,
		Description:  createdRole.Description,
		IsSystemRole: createdRole.IsSystemRole,
		IsEnabled:    createdRole.IsEnabled,
		SortOrder:    createdRole.SortOrder,
		CreatedAt:    createdRole.CreatedAt,
		UpdatedAt:    createdRole.UpdatedAt,
		Permissions:  s.convertToPermissionInfos(permissions),
	}, nil
}

// GetRole 获取角色详情
func (s *RoleService) GetRole(ctx context.Context, roleID int32) (*RoleInfo, error) {
	// 检查权限
	currentUser := middleware.GetCurrentUser(ctx)
	if !currentUser.HasAnyRole("SUPER_ADMIN", "ADMIN", "ROLE_MANAGER") {
		return nil, errors.Forbidden("PERMISSION_DENIED", "无权限查看角色")
	}

	// 获取角色
	role, err := s.roleUc.GetRole(ctx, roleID)
	if err != nil {
		return nil, errors.NotFound("ROLE_NOT_FOUND", "角色不存在")
	}

	// 获取权限
	permissions, _ := s.roleUc.GetRolePermissions(ctx, role.ID)

	return &RoleInfo{
		ID:           role.ID,
		Name:         role.Name,
		Code:         role.Code,
		Description:  role.Description,
		IsSystemRole: role.IsSystemRole,
		IsEnabled:    role.IsEnabled,
		SortOrder:    role.SortOrder,
		CreatedAt:    role.CreatedAt,
		UpdatedAt:    role.UpdatedAt,
		Permissions:  s.convertToPermissionInfos(permissions),
	}, nil
}

// UpdateRole 更新角色
func (s *RoleService) UpdateRole(ctx context.Context, req *UpdateRoleRequest) (*RoleInfo, error) {
	// 检查权限
	currentUser := middleware.GetCurrentUser(ctx)
	if !currentUser.HasAnyRole("SUPER_ADMIN", "ADMIN") {
		return nil, errors.Forbidden("PERMISSION_DENIED", "无权限修改角色")
	}

	s.log.Infof("Updating role: %d by %s", req.ID, currentUser.Username)

	// 获取原角色信息
	role, err := s.roleUc.GetRole(ctx, req.ID)
	if err != nil {
		return nil, errors.NotFound("ROLE_NOT_FOUND", "角色不存在")
	}

	// 检查是否为系统角色
	if role.IsSystemRole {
		return nil, errors.BadRequest("CANNOT_UPDATE_SYSTEM_ROLE", "不能修改系统角色")
	}

	// 更新角色信息
	role.Name = req.Name
	role.Description = req.Description
	role.IsEnabled = req.IsEnabled
	role.SortOrder = req.SortOrder
	role.UpdatedAt = time.Now()

	updatedRole, err := s.roleUc.UpdateRole(ctx, role)
	if err != nil {
		s.log.Errorf("Failed to update role: %v", err)
		if err == biz.ErrRoleNameExists {
			return nil, errors.BadRequest("ROLE_NAME_EXISTS", "角色名称已存在")
		}
		return nil, errors.InternalServer("INTERNAL_ERROR", "角色更新失败")
	}

	// 更新权限
	if len(req.PermissionIDs) > 0 {
		if err := s.roleUc.AssignPermissions(ctx, updatedRole.ID, req.PermissionIDs); err != nil {
			s.log.Errorf("Failed to assign permissions: %v", err)
		}
	}

	// 获取最新权限
	permissions, _ := s.roleUc.GetRolePermissions(ctx, updatedRole.ID)

	s.log.Infof("Role updated successfully: %d", req.ID)

	return &RoleInfo{
		ID:           updatedRole.ID,
		Name:         updatedRole.Name,
		Code:         updatedRole.Code,
		Description:  updatedRole.Description,
		IsSystemRole: updatedRole.IsSystemRole,
		IsEnabled:    updatedRole.IsEnabled,
		SortOrder:    updatedRole.SortOrder,
		CreatedAt:    updatedRole.CreatedAt,
		UpdatedAt:    updatedRole.UpdatedAt,
		Permissions:  s.convertToPermissionInfos(permissions),
	}, nil
}

// DeleteRole 删除角色
func (s *RoleService) DeleteRole(ctx context.Context, roleID int32) error {
	// 检查权限
	currentUser := middleware.GetCurrentUser(ctx)
	if !currentUser.HasAnyRole("SUPER_ADMIN", "ADMIN") {
		return errors.Forbidden("PERMISSION_DENIED", "无权限删除角色")
	}

	s.log.Infof("Deleting role: %d by %s", roleID, currentUser.Username)

	// 删除角色
	if err := s.roleUc.DeleteRole(ctx, roleID); err != nil {
		s.log.Errorf("Failed to delete role: %v", err)
		if err == biz.ErrCannotDeleteSystemRole {
			return errors.BadRequest("CANNOT_DELETE_SYSTEM_ROLE", "不能删除系统角色")
		}
		if err == biz.ErrRoleInUse {
			return errors.BadRequest("ROLE_IN_USE", "角色正在使用中，无法删除")
		}
		return errors.InternalServer("INTERNAL_ERROR", "角色删除失败")
	}

	s.log.Infof("Role deleted successfully: %d", roleID)
	return nil
}

// ListRoles 获取角色列表
func (s *RoleService) ListRoles(ctx context.Context, req *RoleListRequest) (*RoleListResponse, error) {
	// 检查权限
	currentUser := middleware.GetCurrentUser(ctx)
	if !currentUser.HasAnyRole("SUPER_ADMIN", "ADMIN", "ROLE_MANAGER") {
		return nil, errors.Forbidden("PERMISSION_DENIED", "无权限查看角色列表")
	}

	// 设置默认值
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Size <= 0 {
		req.Size = 10
	}

	// 获取角色列表
	roles, total, err := s.roleUc.ListRoles(ctx, req.Page, req.Size, req.Search, req.IsEnabled, req.SortField, req.SortOrder)
	if err != nil {
		s.log.Errorf("Failed to list roles: %v", err)
		return nil, errors.InternalServer("INTERNAL_ERROR", "角色列表获取失败")
	}

	// 转换为响应格式
	roleInfos := make([]*RoleInfo, len(roles))
	for i, role := range roles {
		roleInfos[i] = &RoleInfo{
			ID:           role.ID,
			Name:         role.Name,
			Code:         role.Code,
			Description:  role.Description,
			IsSystemRole: role.IsSystemRole,
			IsEnabled:    role.IsEnabled,
			SortOrder:    role.SortOrder,
			CreatedAt:    role.CreatedAt,
			UpdatedAt:    role.UpdatedAt,
		}
	}

	return &RoleListResponse{
		Roles: roleInfos,
		Total: total,
		Page:  req.Page,
		Size:  req.Size,
	}, nil
}

// AssignPermissions 分配角色权限
func (s *RoleService) AssignPermissions(ctx context.Context, req *AssignPermissionsRequest) error {
	// 检查权限
	currentUser := middleware.GetCurrentUser(ctx)
	if !currentUser.HasAnyRole("SUPER_ADMIN", "ADMIN") {
		return errors.Forbidden("PERMISSION_DENIED", "无权限分配权限")
	}

	s.log.Infof("Assigning permissions to role: %d by %s", req.RoleID, currentUser.Username)

	// 检查角色是否存在
	if _, err := s.roleUc.GetRole(ctx, req.RoleID); err != nil {
		return errors.NotFound("ROLE_NOT_FOUND", "角色不存在")
	}

	// 分配权限
	if err := s.roleUc.AssignPermissions(ctx, req.RoleID, req.PermissionIDs); err != nil {
		s.log.Errorf("Failed to assign permissions: %v", err)
		return errors.InternalServer("INTERNAL_ERROR", "权限分配失败")
	}

	s.log.Infof("Permissions assigned successfully to role: %d", req.RoleID)
	return nil
}

// GetEnabledRoles 获取启用的角色列表
func (s *RoleService) GetEnabledRoles(ctx context.Context) ([]*RoleInfo, error) {
	// 检查权限
	currentUser := middleware.GetCurrentUser(ctx)
	if !currentUser.IsAuthenticated() {
		return nil, errors.Unauthorized("NOT_AUTHENTICATED", "用户未认证")
	}

	// 获取启用的角色
	roles, err := s.roleUc.GetEnabledRoles(ctx)
	if err != nil {
		s.log.Errorf("Failed to get enabled roles: %v", err)
		return nil, errors.InternalServer("INTERNAL_ERROR", "角色列表获取失败")
	}

	// 转换为响应格式
	roleInfos := make([]*RoleInfo, len(roles))
	for i, role := range roles {
		roleInfos[i] = &RoleInfo{
			ID:           role.ID,
			Name:         role.Name,
			Code:         role.Code,
			Description:  role.Description,
			IsSystemRole: role.IsSystemRole,
			IsEnabled:    role.IsEnabled,
			SortOrder:    role.SortOrder,
			CreatedAt:    role.CreatedAt,
			UpdatedAt:    role.UpdatedAt,
		}
	}

	return roleInfos, nil
}

// convertToPermissionInfos 转换权限为响应格式
func (s *RoleService) convertToPermissionInfos(permissions []*biz.Permission) []*PermissionInfo {
	if len(permissions) == 0 {
		return nil
	}

	infos := make([]*PermissionInfo, len(permissions))
	for i, perm := range permissions {
		infos[i] = &PermissionInfo{
			ID:          perm.ID,
			ParentID:    perm.ParentID,
			Name:        perm.Name,
			Code:        perm.Code,
			Resource:    perm.Resource,
			Action:      perm.Action,
			Module:      perm.Module,
			Description: perm.Description,
			IsMenu:      perm.IsMenu,
			Path:        perm.Path,
			SortOrder:   perm.SortOrder,
			CreatedAt:   perm.CreatedAt,
			UpdatedAt:   perm.UpdatedAt,
			Children:    s.convertToPermissionInfos(perm.Children),
		}
	}

	return infos
}