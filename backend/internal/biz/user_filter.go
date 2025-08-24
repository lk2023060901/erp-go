package biz

import (
	"context"
	"errors"
	"time"

	"github.com/go-kratos/kratos/v2/log"
)

// UserFilter 用户过滤器实体
type UserFilter struct {
	ID               int32                  `json:"id"`
	UserID           int32                  `json:"user_id"`
	ModuleType       string                 `json:"module_type"`
	FilterName       string                 `json:"filter_name"`
	FilterConditions map[string]interface{} `json:"filter_conditions"`
	SortConfig       map[string]interface{} `json:"sort_config,omitempty"`
	IsDefault        bool                   `json:"is_default"`
	IsPublic         bool                   `json:"is_public"`
	CreatedAt        time.Time              `json:"created_at"`
	UpdatedAt        time.Time              `json:"updated_at"`
}

// FilterCondition 过滤条件
type FilterCondition struct {
	Field    string      `json:"field"`
	Operator string      `json:"operator"`
	Value    interface{} `json:"value"`
}

// FilterSort 排序配置
type FilterSort struct {
	Field     string `json:"field"`
	Direction string `json:"direction"` // "asc" or "desc"
}

// 错误定义
var (
	ErrFilterNotFound    = errors.New("filter not found")
	ErrFilterNameExists  = errors.New("filter name already exists")
	ErrInvalidModuleType = errors.New("invalid module type")
)

// UserFilterRepo 用户过滤器仓储接口
type UserFilterRepo interface {
	CreateFilter(ctx context.Context, filter *UserFilter) (*UserFilter, error)
	GetFilter(ctx context.Context, id int32) (*UserFilter, error)
	UpdateFilter(ctx context.Context, filter *UserFilter) (*UserFilter, error)
	DeleteFilter(ctx context.Context, id int32) error
	ListFilters(ctx context.Context, userID int32, moduleType string) ([]*UserFilter, error)
	SetDefaultFilter(ctx context.Context, userID int32, moduleType string, filterID int32) error
}

// UserFilterUsecase 用户过滤器用例
type UserFilterUsecase struct {
	repo UserFilterRepo
	log  *log.Helper
}

// NewUserFilterUsecase 创建用户过滤器用例
func NewUserFilterUsecase(repo UserFilterRepo, logger log.Logger) *UserFilterUsecase {
	return &UserFilterUsecase{
		repo: repo,
		log:  log.NewHelper(logger),
	}
}

// CreateFilter 创建过滤器
func (uc *UserFilterUsecase) CreateFilter(ctx context.Context, filter *UserFilter) (*UserFilter, error) {
	// 验证模块类型
	if !uc.isValidModuleType(filter.ModuleType) {
		return nil, ErrInvalidModuleType
	}

	// 如果设置为默认过滤器，需要先清除该用户在该模块的其他默认过滤器
	if filter.IsDefault {
		err := uc.repo.SetDefaultFilter(ctx, filter.UserID, filter.ModuleType, 0) // 先清除所有
		if err != nil {
			uc.log.Errorf("failed to clear default filters: %v", err)
			return nil, err
		}
	}

	return uc.repo.CreateFilter(ctx, filter)
}

// GetFilter 获取过滤器
func (uc *UserFilterUsecase) GetFilter(ctx context.Context, id int32) (*UserFilter, error) {
	return uc.repo.GetFilter(ctx, id)
}

// UpdateFilter 更新过滤器
func (uc *UserFilterUsecase) UpdateFilter(ctx context.Context, filter *UserFilter) (*UserFilter, error) {
	// 验证模块类型
	if !uc.isValidModuleType(filter.ModuleType) {
		return nil, ErrInvalidModuleType
	}

	// 如果设置为默认过滤器，需要先清除该用户在该模块的其他默认过滤器
	if filter.IsDefault {
		err := uc.repo.SetDefaultFilter(ctx, filter.UserID, filter.ModuleType, filter.ID)
		if err != nil {
			uc.log.Errorf("failed to set default filter: %v", err)
			return nil, err
		}
	}

	return uc.repo.UpdateFilter(ctx, filter)
}

// DeleteFilter 删除过滤器
func (uc *UserFilterUsecase) DeleteFilter(ctx context.Context, id int32) error {
	return uc.repo.DeleteFilter(ctx, id)
}

// ListFilters 获取过滤器列表
func (uc *UserFilterUsecase) ListFilters(ctx context.Context, userID int32, moduleType string) ([]*UserFilter, error) {
	if !uc.isValidModuleType(moduleType) {
		return nil, ErrInvalidModuleType
	}

	return uc.repo.ListFilters(ctx, userID, moduleType)
}

// SetDefaultFilter 设置默认过滤器
func (uc *UserFilterUsecase) SetDefaultFilter(ctx context.Context, userID int32, moduleType string, filterID int32) error {
	if !uc.isValidModuleType(moduleType) {
		return ErrInvalidModuleType
	}

	return uc.repo.SetDefaultFilter(ctx, userID, moduleType, filterID)
}

// isValidModuleType 验证模块类型
func (uc *UserFilterUsecase) isValidModuleType(moduleType string) bool {
	validTypes := []string{"users", "roles", "permissions", "organizations", "operation_logs"}
	for _, validType := range validTypes {
		if moduleType == validType {
			return true
		}
	}
	return false
}

// BuildFilterQuery 构建过滤查询条件
func (uc *UserFilterUsecase) BuildFilterQuery(filter *UserFilter) (string, []interface{}) {
	// TODO: 根据过滤条件构建SQL查询
	// 这里可以实现通用的查询构建逻辑
	return "", nil
}