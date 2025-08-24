package service

import (
	"context"
	"time"

	"erp-system/internal/biz"
	"erp-system/internal/middleware"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
)

// OrganizationService 组织服务
type OrganizationService struct {
	orgUc *biz.OrganizationUsecase
	log   *log.Helper
}

// NewOrganizationService 创建组织服务
func NewOrganizationService(orgUc *biz.OrganizationUsecase, logger log.Logger) *OrganizationService {
	return &OrganizationService{
		orgUc: orgUc,
		log:   log.NewHelper(logger),
	}
}

// CreateOrganizationRequest 创建组织请求
type CreateOrganizationRequest struct {
	ParentID    *int32 `json:"parent_id"`
	Name        string `json:"name" validate:"required,min=2,max=50"`
	Code        string `json:"code" validate:"required,min=2,max=50"`
	Description string `json:"description"`
	IsEnabled   bool   `json:"is_enabled"`
	SortOrder   int32  `json:"sort_order"`
}

// UpdateOrganizationRequest 更新组织请求
type UpdateOrganizationRequest struct {
	ID          int32  `json:"id" validate:"required"`
	ParentID    *int32 `json:"parent_id"`
	Name        string `json:"name" validate:"required,min=2,max=50"`
	Description string `json:"description"`
	IsEnabled   bool   `json:"is_enabled"`
	SortOrder   int32  `json:"sort_order"`
}

// AssignUsersRequest 分配用户请求
type AssignUsersRequest struct {
	OrganizationID int32   `json:"organization_id" validate:"required"`
	UserIDs        []int32 `json:"user_ids" validate:"required"`
}

// OrganizationInfo 组织信息
type OrganizationInfo struct {
	ID          int32               `json:"id"`
	ParentID    *int32              `json:"parent_id"`
	Name        string              `json:"name"`
	Code        string              `json:"code"`
	Description string              `json:"description"`
	IsEnabled   bool                `json:"is_enabled"`
	SortOrder   int32               `json:"sort_order"`
	CreatedAt   time.Time           `json:"created_at"`
	UpdatedAt   time.Time           `json:"updated_at"`
	Children    []*OrganizationInfo `json:"children,omitempty"`
	Users       []*UserInfo         `json:"users,omitempty"`
}

// CreateOrganization 创建组织
func (s *OrganizationService) CreateOrganization(ctx context.Context, req *CreateOrganizationRequest) (*OrganizationInfo, error) {
	// 检查权限
	currentUser := middleware.GetCurrentUser(ctx)
	if !currentUser.HasAnyRole("SUPER_ADMIN", "ADMIN", "ORG_MANAGER") {
		return nil, errors.Forbidden("PERMISSION_DENIED", "无权限创建组织")
	}

	s.log.Infof("Creating organization: %s by %s", req.Name, currentUser.Username)

	// 创建组织
	org := &biz.Organization{
		ParentID:    req.ParentID,
		Name:        req.Name,
		Code:        req.Code,
		Description: req.Description,
		IsEnabled:   req.IsEnabled,
		SortOrder:   req.SortOrder,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	createdOrg, err := s.orgUc.CreateOrganization(ctx, org)
	if err != nil {
		s.log.Errorf("Failed to create organization: %v", err)
		if err == biz.ErrOrganizationCodeExists {
			return nil, errors.BadRequest("ORGANIZATION_CODE_EXISTS", "组织编码已存在")
		}
		return nil, errors.InternalServer("INTERNAL_ERROR", "组织创建失败")
	}

	s.log.Infof("Organization created successfully: %s", req.Name)

	return &OrganizationInfo{
		ID:          createdOrg.ID,
		ParentID:    createdOrg.ParentID,
		Name:        createdOrg.Name,
		Code:        createdOrg.Code,
		Description: createdOrg.Description,
		IsEnabled:   createdOrg.IsEnabled,
		SortOrder:   createdOrg.SortOrder,
		CreatedAt:   createdOrg.CreatedAt,
		UpdatedAt:   createdOrg.UpdatedAt,
	}, nil
}

// GetOrganization 获取组织详情
func (s *OrganizationService) GetOrganization(ctx context.Context, orgID int32) (*OrganizationInfo, error) {
	// 检查权限
	currentUser := middleware.GetCurrentUser(ctx)
	if !currentUser.HasAnyRole("SUPER_ADMIN", "ADMIN", "ORG_MANAGER", "USER_MANAGER") {
		return nil, errors.Forbidden("PERMISSION_DENIED", "无权限查看组织")
	}

	// 获取组织
	org, err := s.orgUc.GetOrganization(ctx, orgID)
	if err != nil {
		return nil, errors.NotFound("ORGANIZATION_NOT_FOUND", "组织不存在")
	}

	// 获取组织用户
	users, _ := s.orgUc.GetOrganizationUsers(ctx, org.ID)

	// 转换用户信息
	var userInfos []*UserInfo
	if len(users) > 0 {
		userInfos = make([]*UserInfo, len(users))
		for i, user := range users {
			userInfos[i] = &UserInfo{
				ID:               user.ID,
				Username:         user.Username,
				Email:            user.Email,
				FirstName:        user.FirstName,
				LastName:         user.LastName,
				Phone:            user.Phone,
				Gender:           user.Gender,
				AvatarURL:        user.AvatarURL,
				IsActive:         user.IsActive,
				TwoFactorEnabled: user.TwoFactorEnabled,
				LastLoginAt:      user.LastLoginAt,
				CreatedAt:        user.CreatedAt,
			}
		}
	}

	return &OrganizationInfo{
		ID:          org.ID,
		ParentID:    org.ParentID,
		Name:        org.Name,
		Code:        org.Code,
		Description: org.Description,
		IsEnabled:   org.IsEnabled,
		SortOrder:   org.SortOrder,
		CreatedAt:   org.CreatedAt,
		UpdatedAt:   org.UpdatedAt,
		Users:       userInfos,
	}, nil
}

// UpdateOrganization 更新组织
func (s *OrganizationService) UpdateOrganization(ctx context.Context, req *UpdateOrganizationRequest) (*OrganizationInfo, error) {
	// 检查权限
	currentUser := middleware.GetCurrentUser(ctx)
	if !currentUser.HasAnyRole("SUPER_ADMIN", "ADMIN", "ORG_MANAGER") {
		return nil, errors.Forbidden("PERMISSION_DENIED", "无权限修改组织")
	}

	s.log.Infof("Updating organization: %d by %s", req.ID, currentUser.Username)

	// 获取原组织信息
	org, err := s.orgUc.GetOrganization(ctx, req.ID)
	if err != nil {
		return nil, errors.NotFound("ORGANIZATION_NOT_FOUND", "组织不存在")
	}

	// 更新组织信息
	org.ParentID = req.ParentID
	org.Name = req.Name
	org.Description = req.Description
	org.IsEnabled = req.IsEnabled
	org.SortOrder = req.SortOrder
	org.UpdatedAt = time.Now()

	updatedOrg, err := s.orgUc.UpdateOrganization(ctx, org)
	if err != nil {
		s.log.Errorf("Failed to update organization: %v", err)
		return nil, errors.InternalServer("INTERNAL_ERROR", "组织更新失败")
	}

	s.log.Infof("Organization updated successfully: %d", req.ID)

	return &OrganizationInfo{
		ID:          updatedOrg.ID,
		ParentID:    updatedOrg.ParentID,
		Name:        updatedOrg.Name,
		Code:        updatedOrg.Code,
		Description: updatedOrg.Description,
		IsEnabled:   updatedOrg.IsEnabled,
		SortOrder:   updatedOrg.SortOrder,
		CreatedAt:   updatedOrg.CreatedAt,
		UpdatedAt:   updatedOrg.UpdatedAt,
	}, nil
}

// DeleteOrganization 删除组织
func (s *OrganizationService) DeleteOrganization(ctx context.Context, orgID int32) error {
	// 检查权限
	currentUser := middleware.GetCurrentUser(ctx)
	if !currentUser.HasAnyRole("SUPER_ADMIN", "ADMIN") {
		return errors.Forbidden("PERMISSION_DENIED", "无权限删除组织")
	}

	s.log.Infof("Deleting organization: %d by %s", orgID, currentUser.Username)

	// 删除组织
	if err := s.orgUc.DeleteOrganization(ctx, orgID); err != nil {
		s.log.Errorf("Failed to delete organization: %v", err)
		if err == biz.ErrOrganizationHasChildren {
			return errors.BadRequest("ORGANIZATION_HAS_CHILDREN", "组织存在子组织，无法删除")
		}
		if err == biz.ErrOrganizationHasUsers {
			return errors.BadRequest("ORGANIZATION_HAS_USERS", "组织存在用户，无法删除")
		}
		return errors.InternalServer("INTERNAL_ERROR", "组织删除失败")
	}

	s.log.Infof("Organization deleted successfully: %d", orgID)
	return nil
}

// GetOrganizationTree 获取组织树
func (s *OrganizationService) GetOrganizationTree(ctx context.Context) ([]*OrganizationInfo, error) {
	// 检查权限
	currentUser := middleware.GetCurrentUser(ctx)
	if !currentUser.HasAnyRole("SUPER_ADMIN", "ADMIN", "ORG_MANAGER", "USER_MANAGER") {
		return nil, errors.Forbidden("PERMISSION_DENIED", "无权限查看组织树")
	}

	// 获取组织树
	orgs, err := s.orgUc.GetOrganizationTree(ctx)
	if err != nil {
		s.log.Errorf("Failed to get organization tree: %v", err)
		return nil, errors.InternalServer("INTERNAL_ERROR", "组织树获取失败")
	}

	// 转换为响应格式
	return s.convertToOrganizationInfos(orgs), nil
}

// AssignUsers 分配用户到组织
func (s *OrganizationService) AssignUsers(ctx context.Context, req *AssignUsersRequest) error {
	// 检查权限
	currentUser := middleware.GetCurrentUser(ctx)
	if !currentUser.HasAnyRole("SUPER_ADMIN", "ADMIN", "ORG_MANAGER") {
		return errors.Forbidden("PERMISSION_DENIED", "无权限分配用户")
	}

	s.log.Infof("Assigning users to organization: %d by %s", req.OrganizationID, currentUser.Username)

	// 检查组织是否存在
	if _, err := s.orgUc.GetOrganization(ctx, req.OrganizationID); err != nil {
		return errors.NotFound("ORGANIZATION_NOT_FOUND", "组织不存在")
	}

	// 分配用户
	if err := s.orgUc.AssignUsers(ctx, req.OrganizationID, req.UserIDs); err != nil {
		s.log.Errorf("Failed to assign users: %v", err)
		return errors.InternalServer("INTERNAL_ERROR", "用户分配失败")
	}

	s.log.Infof("Users assigned successfully to organization: %d", req.OrganizationID)
	return nil
}

// GetEnabledOrganizations 获取启用的组织列表
func (s *OrganizationService) GetEnabledOrganizations(ctx context.Context) ([]*OrganizationInfo, error) {
	// 检查权限
	currentUser := middleware.GetCurrentUser(ctx)
	if !currentUser.IsAuthenticated() {
		return nil, errors.Unauthorized("NOT_AUTHENTICATED", "用户未认证")
	}

	// 获取启用的组织
	orgs, err := s.orgUc.GetEnabledOrganizations(ctx)
	if err != nil {
		s.log.Errorf("Failed to get enabled organizations: %v", err)
		return nil, errors.InternalServer("INTERNAL_ERROR", "组织列表获取失败")
	}

	// 转换为响应格式
	return s.convertToOrganizationInfos(orgs), nil
}

// convertToOrganizationInfos 转换组织为响应格式
func (s *OrganizationService) convertToOrganizationInfos(orgs []*biz.Organization) []*OrganizationInfo {
	if len(orgs) == 0 {
		return nil
	}

	infos := make([]*OrganizationInfo, len(orgs))
	for i, org := range orgs {
		infos[i] = &OrganizationInfo{
			ID:          org.ID,
			ParentID:    org.ParentID,
			Name:        org.Name,
			Code:        org.Code,
			Description: org.Description,
			IsEnabled:   org.IsEnabled,
			SortOrder:   org.SortOrder,
			CreatedAt:   org.CreatedAt,
			UpdatedAt:   org.UpdatedAt,
			Children:    s.convertToOrganizationInfos(org.Children),
		}

		// 转换用户信息
		if len(org.Users) > 0 {
			infos[i].Users = make([]*UserInfo, len(org.Users))
			for j, user := range org.Users {
				infos[i].Users[j] = &UserInfo{
					ID:               user.ID,
					Username:         user.Username,
					Email:            user.Email,
					FirstName:        user.FirstName,
					LastName:         user.LastName,
					Phone:            user.Phone,
					Gender:           user.Gender,
					AvatarURL:        user.AvatarURL,
					IsActive:         user.IsActive,
					TwoFactorEnabled: user.TwoFactorEnabled,
					LastLoginAt:      user.LastLoginAt,
					CreatedAt:        user.CreatedAt,
				}
			}
		}
	}

	return infos
}
