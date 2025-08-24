package data

import (
	"context"
	"time"

	"erp-system/internal/biz"
	"erp-system/internal/cache"

	"github.com/go-kratos/kratos/v2/log"
)

// CachedPermissionRepo 带缓存的权限仓库
type CachedPermissionRepo struct {
	repo  biz.PermissionRepo
	cache cache.PermissionCache
	log   *log.Helper
	
	// 缓存TTL配置
	userPermissionTTL    time.Duration
	userRoleTTL          time.Duration
	permissionRuleTTL    time.Duration
	fieldPermissionTTL   time.Duration
	docTypeTTL           time.Duration
	userPermissionLevelTTL time.Duration
}

// NewCachedPermissionRepo 创建带缓存的权限仓库
func NewCachedPermissionRepo(repo biz.PermissionRepo, cache cache.PermissionCache, logger log.Logger) biz.PermissionRepo {
	return &CachedPermissionRepo{
		repo:  repo,
		cache: cache,
		log:   log.NewHelper(logger),
		
		// 设置默认缓存TTL
		userPermissionTTL:      15 * time.Minute,
		userRoleTTL:            30 * time.Minute,
		permissionRuleTTL:      1 * time.Hour,
		fieldPermissionTTL:     2 * time.Hour,
		docTypeTTL:             4 * time.Hour,
		userPermissionLevelTTL: 10 * time.Minute,
	}
}

// 文档类型管理 - 带缓存
func (r *CachedPermissionRepo) CreateDocType(ctx context.Context, docType *biz.DocType) (*biz.DocType, error) {
	result, err := r.repo.CreateDocType(ctx, docType)
	if err != nil {
		return nil, err
	}
	
	// 缓存新创建的文档类型
	if err := r.cache.SetDocType(ctx, result.Name, result, r.docTypeTTL); err != nil {
		r.log.Warnf("Failed to cache created doctype %s: %v", result.Name, err)
	}
	
	return result, nil
}

func (r *CachedPermissionRepo) GetDocType(ctx context.Context, name string) (*biz.DocType, error) {
	// 先从缓存获取
	cached, err := r.cache.GetDocType(ctx, name)
	if err != nil {
		r.log.Warnf("Failed to get doctype %s from cache: %v", name, err)
	} else if cached != nil {
		return cached, nil
	}
	
	// 缓存未命中，从数据库获取
	docType, err := r.repo.GetDocType(ctx, name)
	if err != nil {
		return nil, err
	}
	
	// 缓存结果
	if err := r.cache.SetDocType(ctx, name, docType, r.docTypeTTL); err != nil {
		r.log.Warnf("Failed to cache doctype %s: %v", name, err)
	}
	
	return docType, nil
}

func (r *CachedPermissionRepo) UpdateDocType(ctx context.Context, docType *biz.DocType) (*biz.DocType, error) {
	result, err := r.repo.UpdateDocType(ctx, docType)
	if err != nil {
		return nil, err
	}
	
	// 更新缓存
	if err := r.cache.SetDocType(ctx, result.Name, result, r.docTypeTTL); err != nil {
		r.log.Warnf("Failed to update cached doctype %s: %v", result.Name, err)
	}
	
	// 清除相关缓存
	if err := r.cache.ClearDocTypeCache(ctx, result.Name); err != nil {
		r.log.Warnf("Failed to clear doctype cache for %s: %v", result.Name, err)
	}
	
	return result, nil
}

func (r *CachedPermissionRepo) DeleteDocType(ctx context.Context, name string) error {
	err := r.repo.DeleteDocType(ctx, name)
	if err != nil {
		return err
	}
	
	// 清除缓存
	if err := r.cache.DeleteDocType(ctx, name); err != nil {
		r.log.Warnf("Failed to delete cached doctype %s: %v", name, err)
	}
	
	// 清除相关缓存
	if err := r.cache.ClearDocTypeCache(ctx, name); err != nil {
		r.log.Warnf("Failed to clear doctype cache for %s: %v", name, err)
	}
	
	return nil
}

func (r *CachedPermissionRepo) ListDocTypes(ctx context.Context, module string) ([]*biz.DocType, error) {
	// 文档类型列表不缓存，因为过滤条件多样且不频繁调用
	return r.repo.ListDocTypes(ctx, module)
}

// 权限规则管理 - 带缓存
func (r *CachedPermissionRepo) CreatePermissionRule(ctx context.Context, rule *biz.PermissionRule) (*biz.PermissionRule, error) {
	result, err := r.repo.CreatePermissionRule(ctx, rule)
	if err != nil {
		return nil, err
	}
	
	// 清除相关角色的权限规则缓存
	if err := r.cache.DeletePermissionRules(ctx, result.RoleID, result.DocType); err != nil {
		r.log.Warnf("Failed to clear permission rule cache for role %d doctype %s: %v", result.RoleID, result.DocType, err)
	}
	
	// 清除角色相关缓存
	if err := r.cache.ClearRoleCache(ctx, result.RoleID); err != nil {
		r.log.Warnf("Failed to clear role cache for role %d: %v", result.RoleID, err)
	}
	
	return result, nil
}

func (r *CachedPermissionRepo) GetPermissionRule(ctx context.Context, id int64) (*biz.PermissionRule, error) {
	// 单个权限规则不缓存，因为通常通过角色+文档类型查询
	return r.repo.GetPermissionRule(ctx, id)
}

func (r *CachedPermissionRepo) ListPermissionRules(ctx context.Context, roleID int64, docType string) ([]*biz.PermissionRule, error) {
	// 先从缓存获取
	cached, err := r.cache.GetPermissionRules(ctx, roleID, docType)
	if err != nil {
		r.log.Warnf("Failed to get permission rules from cache for role %d doctype %s: %v", roleID, docType, err)
	} else if cached != nil {
		return cached, nil
	}
	
	// 缓存未命中，从数据库获取
	rules, err := r.repo.ListPermissionRules(ctx, roleID, docType)
	if err != nil {
		return nil, err
	}
	
	// 缓存结果
	if err := r.cache.SetPermissionRules(ctx, roleID, docType, rules, r.permissionRuleTTL); err != nil {
		r.log.Warnf("Failed to cache permission rules for role %d doctype %s: %v", roleID, docType, err)
	}
	
	return rules, nil
}

func (r *CachedPermissionRepo) UpdatePermissionRule(ctx context.Context, rule *biz.PermissionRule) (*biz.PermissionRule, error) {
	result, err := r.repo.UpdatePermissionRule(ctx, rule)
	if err != nil {
		return nil, err
	}
	
	// 清除相关缓存
	if err := r.cache.DeletePermissionRules(ctx, result.RoleID, result.DocType); err != nil {
		r.log.Warnf("Failed to clear permission rule cache for role %d doctype %s: %v", result.RoleID, result.DocType, err)
	}
	
	// 清除角色相关缓存
	if err := r.cache.ClearRoleCache(ctx, result.RoleID); err != nil {
		r.log.Warnf("Failed to clear role cache for role %d: %v", result.RoleID, err)
	}
	
	return result, nil
}

func (r *CachedPermissionRepo) DeletePermissionRule(ctx context.Context, id int64) error {
	// 先获取权限规则信息，用于清除缓存
	rule, err := r.repo.GetPermissionRule(ctx, id)
	if err != nil {
		return err
	}
	
	// 删除权限规则
	err = r.repo.DeletePermissionRule(ctx, id)
	if err != nil {
		return err
	}
	
	// 清除相关缓存
	if rule != nil {
		if err := r.cache.DeletePermissionRules(ctx, rule.RoleID, rule.DocType); err != nil {
			r.log.Warnf("Failed to clear permission rule cache for role %d doctype %s: %v", rule.RoleID, rule.DocType, err)
		}
		
		if err := r.cache.ClearRoleCache(ctx, rule.RoleID); err != nil {
			r.log.Warnf("Failed to clear role cache for role %d: %v", rule.RoleID, err)
		}
	}
	
	return nil
}

// 用户权限管理 - 带缓存
func (r *CachedPermissionRepo) CreateUserPermission(ctx context.Context, userPermission *biz.UserPermission) (*biz.UserPermission, error) {
	result, err := r.repo.CreateUserPermission(ctx, userPermission)
	if err != nil {
		return nil, err
	}
	
	// 清除用户相关缓存
	if err := r.cache.ClearUserCache(ctx, result.UserID); err != nil {
		r.log.Warnf("Failed to clear user cache for user %d: %v", result.UserID, err)
	}
	
	return result, nil
}

func (r *CachedPermissionRepo) GetUserPermission(ctx context.Context, id int64) (*biz.UserPermission, error) {
	return r.repo.GetUserPermission(ctx, id)
}

func (r *CachedPermissionRepo) UpdateUserPermission(ctx context.Context, userPermission *biz.UserPermission) (*biz.UserPermission, error) {
	result, err := r.repo.UpdateUserPermission(ctx, userPermission)
	if err != nil {
		return nil, err
	}
	
	// 清除用户相关缓存
	if err := r.cache.ClearUserCache(ctx, result.UserID); err != nil {
		r.log.Warnf("Failed to clear user cache for user %d: %v", result.UserID, err)
	}
	
	return result, nil
}

func (r *CachedPermissionRepo) DeleteUserPermission(ctx context.Context, id int64) error {
	// 先获取用户权限信息，用于清除缓存
	userPerm, err := r.repo.GetUserPermission(ctx, id)
	if err != nil {
		return err
	}
	
	// 删除用户权限
	err = r.repo.DeleteUserPermission(ctx, id)
	if err != nil {
		return err
	}
	
	// 清除用户相关缓存
	if userPerm != nil {
		if err := r.cache.ClearUserCache(ctx, userPerm.UserID); err != nil {
			r.log.Warnf("Failed to clear user cache for user %d: %v", userPerm.UserID, err)
		}
	}
	
	return nil
}

func (r *CachedPermissionRepo) ListUserPermissions(ctx context.Context, userID int64, docType string, page, size int32) ([]*biz.UserPermission, error) {
	// 用户权限列表不缓存，因为分页参数多样
	return r.repo.ListUserPermissions(ctx, userID, docType, page, size)
}

func (r *CachedPermissionRepo) GetUserPermissionsCount(ctx context.Context, userID int64, docType string) (int32, error) {
	return r.repo.GetUserPermissionsCount(ctx, userID, docType)
}

// 字段权限级别管理 - 带缓存
func (r *CachedPermissionRepo) CreateFieldPermissionLevel(ctx context.Context, field *biz.FieldPermissionLevel) (*biz.FieldPermissionLevel, error) {
	result, err := r.repo.CreateFieldPermissionLevel(ctx, field)
	if err != nil {
		return nil, err
	}
	
	// 清除字段权限级别缓存
	if err := r.cache.DeleteFieldPermissionLevels(ctx, result.DocType); err != nil {
		r.log.Warnf("Failed to clear field permission levels cache for doctype %s: %v", result.DocType, err)
	}
	
	return result, nil
}

func (r *CachedPermissionRepo) GetFieldPermissionLevel(ctx context.Context, id int64) (*biz.FieldPermissionLevel, error) {
	return r.repo.GetFieldPermissionLevel(ctx, id)
}

func (r *CachedPermissionRepo) UpdateFieldPermissionLevel(ctx context.Context, field *biz.FieldPermissionLevel) (*biz.FieldPermissionLevel, error) {
	result, err := r.repo.UpdateFieldPermissionLevel(ctx, field)
	if err != nil {
		return nil, err
	}
	
	// 清除字段权限级别缓存
	if err := r.cache.DeleteFieldPermissionLevels(ctx, result.DocType); err != nil {
		r.log.Warnf("Failed to clear field permission levels cache for doctype %s: %v", result.DocType, err)
	}
	
	return result, nil
}

func (r *CachedPermissionRepo) DeleteFieldPermissionLevel(ctx context.Context, id int64) error {
	// 先获取字段权限级别信息，用于清除缓存
	field, err := r.repo.GetFieldPermissionLevel(ctx, id)
	if err != nil {
		return err
	}
	
	// 删除字段权限级别
	err = r.repo.DeleteFieldPermissionLevel(ctx, id)
	if err != nil {
		return err
	}
	
	// 清除相关缓存
	if field != nil {
		if err := r.cache.DeleteFieldPermissionLevels(ctx, field.DocType); err != nil {
			r.log.Warnf("Failed to clear field permission levels cache for doctype %s: %v", field.DocType, err)
		}
	}
	
	return nil
}

func (r *CachedPermissionRepo) ListFieldPermissionLevels(ctx context.Context, docType string, page, size int32) ([]*biz.FieldPermissionLevel, error) {
	// 字段权限级别列表不缓存，因为分页参数多样
	return r.repo.ListFieldPermissionLevels(ctx, docType, page, size)
}

func (r *CachedPermissionRepo) GetFieldPermissionLevelsCount(ctx context.Context, docType string) (int32, error) {
	return r.repo.GetFieldPermissionLevelsCount(ctx, docType)
}

// 文档工作流状态管理 - 不缓存，因为状态变化频繁
func (r *CachedPermissionRepo) CreateDocumentWorkflowState(ctx context.Context, state *biz.DocumentWorkflowState) (*biz.DocumentWorkflowState, error) {
	return r.repo.CreateDocumentWorkflowState(ctx, state)
}

func (r *CachedPermissionRepo) GetDocumentWorkflowState(ctx context.Context, stateID int64) (*biz.DocumentWorkflowState, error) {
	return r.repo.GetDocumentWorkflowState(ctx, stateID)
}

func (r *CachedPermissionRepo) UpdateDocumentWorkflowState(ctx context.Context, state *biz.DocumentWorkflowState) (*biz.DocumentWorkflowState, error) {
	return r.repo.UpdateDocumentWorkflowState(ctx, state)
}

func (r *CachedPermissionRepo) DeleteDocumentWorkflowState(ctx context.Context, stateID int64) error {
	return r.repo.DeleteDocumentWorkflowState(ctx, stateID)
}

func (r *CachedPermissionRepo) ListDocumentWorkflowStates(ctx context.Context, docType, documentName, state string, userID int64, page, size int32) ([]*biz.DocumentWorkflowState, error) {
	return r.repo.ListDocumentWorkflowStates(ctx, docType, documentName, state, userID, page, size)
}

func (r *CachedPermissionRepo) GetDocumentWorkflowStatesCount(ctx context.Context, docType, documentName, state string, userID int64) (int32, error) {
	return r.repo.GetDocumentWorkflowStatesCount(ctx, docType, documentName, state, userID)
}

// 权限检查和查询 - 带缓存
func (r *CachedPermissionRepo) CheckDocumentPermission(ctx context.Context, req *biz.PermissionCheckRequest) (bool, error) {
	// 权限检查结果不缓存，因为参数复杂且安全性要求高
	return r.repo.CheckDocumentPermission(ctx, req)
}

func (r *CachedPermissionRepo) CheckPermission(ctx context.Context, userID int64, documentType, action string, permissionLevel int) (bool, error) {
	// 权限检查结果不缓存，因为安全性要求高
	return r.repo.CheckPermission(ctx, userID, documentType, action, permissionLevel)
}

func (r *CachedPermissionRepo) GetUserPermissionLevel(ctx context.Context, userID int64, documentType string) (int, error) {
	// 先从缓存获取
	level, err := r.cache.GetUserPermissionLevel(ctx, userID, documentType)
	if err != nil {
		r.log.Warnf("Failed to get user permission level from cache for user %d doctype %s: %v", userID, documentType, err)
	} else if level >= 0 {
		return level, nil
	}
	
	// 缓存未命中，从数据库获取
	level, err = r.repo.GetUserPermissionLevel(ctx, userID, documentType)
	if err != nil {
		return 0, err
	}
	
	// 缓存结果
	if err := r.cache.SetUserPermissionLevel(ctx, userID, documentType, level, r.userPermissionLevelTTL); err != nil {
		r.log.Warnf("Failed to cache user permission level for user %d doctype %s: %v", userID, documentType, err)
	}
	
	return level, nil
}

func (r *CachedPermissionRepo) GetUserEnhancedPermissions(ctx context.Context, userID int64, docType string) ([]*biz.EnhancedUserPermission, error) {
	// 增强用户权限不缓存，因为结构复杂且使用频率不高
	return r.repo.GetUserEnhancedPermissions(ctx, userID, docType)
}

func (r *CachedPermissionRepo) GetAccessibleFields(ctx context.Context, req *biz.FieldPermissionRequest) ([]*biz.AccessibleField, error) {
	// 字段访问权限不缓存，因为参数复杂
	return r.repo.GetAccessibleFields(ctx, req)
}

func (r *CachedPermissionRepo) FilterDocumentsByPermission(ctx context.Context, userID int64, documentType string, documents []map[string]interface{}) ([]map[string]interface{}, error) {
	// 文档过滤不缓存，因为数据动态性强
	return r.repo.FilterDocumentsByPermission(ctx, userID, documentType, documents)
}

func (r *CachedPermissionRepo) GetUserRoles(ctx context.Context, userID int64) ([]string, error) {
	// 先从缓存获取
	cached, err := r.cache.GetUserRoles(ctx, userID)
	if err != nil {
		r.log.Warnf("Failed to get user roles from cache for user %d: %v", userID, err)
	} else if cached != nil {
		return cached, nil
	}
	
	// 缓存未命中，从数据库获取
	roles, err := r.repo.GetUserRoles(ctx, userID)
	if err != nil {
		return nil, err
	}
	
	// 缓存结果
	if err := r.cache.SetUserRoles(ctx, userID, roles, r.userRoleTTL); err != nil {
		r.log.Warnf("Failed to cache user roles for user %d: %v", userID, err)
	}
	
	return roles, nil
}

// 批量操作 - 清除相关缓存
func (r *CachedPermissionRepo) BatchCreatePermissionRules(ctx context.Context, rules []*biz.PermissionRule) error {
	err := r.repo.BatchCreatePermissionRules(ctx, rules)
	if err != nil {
		return err
	}
	
	// 清除相关缓存
	roleDocTypeMap := make(map[int64]map[string]bool)
	for _, rule := range rules {
		if roleDocTypeMap[rule.RoleID] == nil {
			roleDocTypeMap[rule.RoleID] = make(map[string]bool)
		}
		roleDocTypeMap[rule.RoleID][rule.DocType] = true
	}
	
	for roleID, docTypes := range roleDocTypeMap {
		// 清除角色缓存
		if err := r.cache.ClearRoleCache(ctx, roleID); err != nil {
			r.log.Warnf("Failed to clear role cache for role %d: %v", roleID, err)
		}
		
		// 清除权限规则缓存
		for docType := range docTypes {
			if err := r.cache.DeletePermissionRules(ctx, roleID, docType); err != nil {
				r.log.Warnf("Failed to clear permission rule cache for role %d doctype %s: %v", roleID, docType, err)
			}
		}
	}
	
	return nil
}

func (r *CachedPermissionRepo) BatchCreateUserPermissions(ctx context.Context, permissions []*biz.UserPermission) error {
	err := r.repo.BatchCreateUserPermissions(ctx, permissions)
	if err != nil {
		return err
	}
	
	// 清除用户缓存
	userIDs := make(map[int64]bool)
	for _, perm := range permissions {
		userIDs[perm.UserID] = true
	}
	
	for userID := range userIDs {
		if err := r.cache.ClearUserCache(ctx, userID); err != nil {
			r.log.Warnf("Failed to clear user cache for user %d: %v", userID, err)
		}
	}
	
	return nil
}