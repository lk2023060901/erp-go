package service

import (
	"context"

	"erp-system/internal/biz"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
)

// UserFilterService 用户过滤器服务
type UserFilterService struct {
	filterUc *biz.UserFilterUsecase
	log      *log.Helper
}

// NewUserFilterService 创建用户过滤器服务
func NewUserFilterService(
	filterUc *biz.UserFilterUsecase,
	logger log.Logger,
) *UserFilterService {
	return &UserFilterService{
		filterUc: filterUc,
		log:      log.NewHelper(logger),
	}
}

// CreateFilterRequest 创建过滤器请求
type CreateFilterRequest struct {
	ModuleType       string                 `json:"module_type" validate:"required,oneof=users roles permissions organizations operation_logs"`
	FilterName       string                 `json:"filter_name" validate:"required,min=1,max=100"`
	FilterConditions map[string]interface{} `json:"filter_conditions" validate:"required"`
	SortConfig       map[string]interface{} `json:"sort_config,omitempty"`
	IsDefault        bool                   `json:"is_default"`
	IsPublic         bool                   `json:"is_public"`
}

// UpdateFilterRequest 更新过滤器请求
type UpdateFilterRequest struct {
	FilterName       string                 `json:"filter_name" validate:"required,min=1,max=100"`
	FilterConditions map[string]interface{} `json:"filter_conditions" validate:"required"`
	SortConfig       map[string]interface{} `json:"sort_config,omitempty"`
	IsDefault        bool                   `json:"is_default"`
	IsPublic         bool                   `json:"is_public"`
}

// FilterResponse 过滤器响应
type FilterResponse struct {
	ID               int32                  `json:"id"`
	UserID           int32                  `json:"user_id"`
	ModuleType       string                 `json:"module_type"`
	FilterName       string                 `json:"filter_name"`
	FilterConditions map[string]interface{} `json:"filter_conditions"`
	SortConfig       map[string]interface{} `json:"sort_config,omitempty"`
	IsDefault        bool                   `json:"is_default"`
	IsPublic         bool                   `json:"is_public"`
	CreatedAt        string                 `json:"created_at"`
	UpdatedAt        string                 `json:"updated_at"`
}

// CreateFilter 创建过滤器
func (s *UserFilterService) CreateFilter(ctx context.Context, req *CreateFilterRequest) (*FilterResponse, error) {
	// 从上下文获取当前用户ID
	userID := int32(1) // 临时硬编码，实际项目中应该从上下文获取
	// userID, err := middleware.GetUserID(ctx)
	// if err != nil {
	//     return nil, errors.Unauthorized("USER_AUTH_FAILED", "用户认证失败")
	// }

	// 创建过滤器实体
	filter := &biz.UserFilter{
		UserID:           userID,
		ModuleType:       req.ModuleType,
		FilterName:       req.FilterName,
		FilterConditions: req.FilterConditions,
		SortConfig:       req.SortConfig,
		IsDefault:        req.IsDefault,
		IsPublic:         req.IsPublic,
	}

	// 创建过滤器
	createdFilter, err := s.filterUc.CreateFilter(ctx, filter)
	if err != nil {
		s.log.Errorf("failed to create filter: %v", err)
		return nil, errors.InternalServer("CREATE_FILTER_FAILED", "创建过滤器失败")
	}

	return s.convertToResponse(createdFilter), nil
}

// GetFilter 获取过滤器
func (s *UserFilterService) GetFilter(ctx context.Context, id int32) (*FilterResponse, error) {
	filter, err := s.filterUc.GetFilter(ctx, id)
	if err != nil {
		if err == biz.ErrFilterNotFound {
			return nil, errors.NotFound("FILTER_NOT_FOUND", "过滤器不存在")
		}
		s.log.Errorf("failed to get filter: %v", err)
		return nil, errors.InternalServer("GET_FILTER_FAILED", "获取过滤器失败")
	}

	// 验证用户权限（只能访问自己的过滤器或公共过滤器）
	userID := int32(1) // 临时硬编码，实际项目中应该从上下文获取

	if filter.UserID != userID && !filter.IsPublic {
		return nil, errors.Forbidden("ACCESS_DENIED", "无权访问此过滤器")
	}

	return s.convertToResponse(filter), nil
}

// UpdateFilter 更新过滤器
func (s *UserFilterService) UpdateFilter(ctx context.Context, id int32, req *UpdateFilterRequest) (*FilterResponse, error) {
	// 获取当前用户ID
	userID := int32(1) // 临时硬编码，实际项目中应该从上下文获取

	// 获取现有过滤器
	existingFilter, err := s.filterUc.GetFilter(ctx, id)
	if err != nil {
		if err == biz.ErrFilterNotFound {
			return nil, errors.NotFound("FILTER_NOT_FOUND", "过滤器不存在")
		}
		return nil, errors.InternalServer("GET_FILTER_FAILED", "获取过滤器失败")
	}

	// 验证权限（只能更新自己的过滤器）
	if existingFilter.UserID != userID {
		return nil, errors.Forbidden("ACCESS_DENIED", "无权修改此过滤器")
	}

	// 更新过滤器数据
	existingFilter.FilterName = req.FilterName
	existingFilter.FilterConditions = req.FilterConditions
	existingFilter.SortConfig = req.SortConfig
	existingFilter.IsDefault = req.IsDefault
	existingFilter.IsPublic = req.IsPublic

	// 执行更新
	updatedFilter, err := s.filterUc.UpdateFilter(ctx, existingFilter)
	if err != nil {
		s.log.Errorf("failed to update filter: %v", err)
		return nil, errors.InternalServer("UPDATE_FILTER_FAILED", "更新过滤器失败")
	}

	return s.convertToResponse(updatedFilter), nil
}

// DeleteFilter 删除过滤器
func (s *UserFilterService) DeleteFilter(ctx context.Context, id int32) error {
	// 获取当前用户ID
	userID := int32(1) // 临时硬编码，实际项目中应该从上下文获取

	// 获取过滤器验证权限
	filter, err := s.filterUc.GetFilter(ctx, id)
	if err != nil {
		if err == biz.ErrFilterNotFound {
			return errors.NotFound("FILTER_NOT_FOUND", "过滤器不存在")
		}
		return errors.InternalServer("GET_FILTER_FAILED", "获取过滤器失败")
	}

	// 验证权限（只能删除自己的过滤器）
	if filter.UserID != userID {
		return errors.Forbidden("ACCESS_DENIED", "无权删除此过滤器")
	}

	// 执行删除
	err = s.filterUc.DeleteFilter(ctx, id)
	if err != nil {
		s.log.Errorf("failed to delete filter: %v", err)
		return errors.InternalServer("DELETE_FILTER_FAILED", "删除过滤器失败")
	}

	return nil
}

// ListFilters 获取过滤器列表
func (s *UserFilterService) ListFilters(ctx context.Context, moduleType string) ([]*FilterResponse, error) {
	// 获取当前用户ID
	userID := int32(1) // 临时硬编码，实际项目中应该从上下文获取

	// 获取过滤器列表
	filters, err := s.filterUc.ListFilters(ctx, userID, moduleType)
	if err != nil {
		if err == biz.ErrInvalidModuleType {
			return nil, errors.BadRequest("INVALID_MODULE_TYPE", "无效的模块类型")
		}
		s.log.Errorf("failed to list filters: %v", err)
		return nil, errors.InternalServer("LIST_FILTERS_FAILED", "获取过滤器列表失败")
	}

	// 转换响应
	responses := make([]*FilterResponse, len(filters))
	for i, filter := range filters {
		responses[i] = s.convertToResponse(filter)
	}

	return responses, nil
}

// SetDefaultFilter 设置默认过滤器
func (s *UserFilterService) SetDefaultFilter(ctx context.Context, moduleType string, filterID int32) error {
	// 获取当前用户ID
	userID := int32(1) // 临时硬编码，实际项目中应该从上下文获取

	// 验证过滤器权限
	filter, err := s.filterUc.GetFilter(ctx, filterID)
	if err != nil {
		if err == biz.ErrFilterNotFound {
			return errors.NotFound("FILTER_NOT_FOUND", "过滤器不存在")
		}
		return errors.InternalServer("GET_FILTER_FAILED", "获取过滤器失败")
	}

	// 验证权限（只能设置自己的过滤器为默认）
	if filter.UserID != userID {
		return errors.Forbidden("ACCESS_DENIED", "无权设置此过滤器为默认")
	}

	// 设置默认过滤器
	err = s.filterUc.SetDefaultFilter(ctx, userID, moduleType, filterID)
	if err != nil {
		if err == biz.ErrInvalidModuleType {
			return errors.BadRequest("INVALID_MODULE_TYPE", "无效的模块类型")
		}
		s.log.Errorf("failed to set default filter: %v", err)
		return errors.InternalServer("SET_DEFAULT_FILTER_FAILED", "设置默认过滤器失败")
	}

	return nil
}

// convertToResponse 转换为响应格式
func (s *UserFilterService) convertToResponse(filter *biz.UserFilter) *FilterResponse {
	return &FilterResponse{
		ID:               filter.ID,
		UserID:           filter.UserID,
		ModuleType:       filter.ModuleType,
		FilterName:       filter.FilterName,
		FilterConditions: filter.FilterConditions,
		SortConfig:       filter.SortConfig,
		IsDefault:        filter.IsDefault,
		IsPublic:         filter.IsPublic,
		CreatedAt:        filter.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:        filter.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}