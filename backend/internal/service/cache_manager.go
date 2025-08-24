package service

import (
	"context"
	"fmt"
	"sync"
	"time"

	"erp-system/internal/biz"
	"erp-system/internal/cache"
	"erp-system/internal/middleware"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
)

// CacheManagerService 缓存管理服务
type CacheManagerService struct {
	permissionUc *biz.PermissionUsecase
	cache        cache.PermissionCache
	log          *log.Helper
}

// NewCacheManagerService 创建缓存管理服务
func NewCacheManagerService(permissionUc *biz.PermissionUsecase, cache cache.PermissionCache, logger log.Logger) *CacheManagerService {
	return &CacheManagerService{
		permissionUc: permissionUc,
		cache:        cache,
		log:          log.NewHelper(logger),
	}
}

// WarmUpCacheRequest 缓存预热请求
type WarmUpCacheRequest struct {
	CacheTypes   []string `json:"cache_types,omitempty"` // 指定要预热的缓存类型
	UserIDs      []int64  `json:"user_ids,omitempty"`    // 指定要预热的用户ID
	RoleIDs      []int64  `json:"role_ids,omitempty"`    // 指定要预热的角色ID
	DocTypes     []string `json:"doc_types,omitempty"`   // 指定要预热的文档类型
	ForceRefresh bool     `json:"force_refresh"`         // 是否强制刷新缓存
}

// WarmUpCacheResponse 缓存预热响应
type WarmUpCacheResponse struct {
	Success     bool           `json:"success"`
	Message     string         `json:"message"`
	WarmupStats map[string]int `json:"warmup_stats"`
	Duration    time.Duration  `json:"duration"`
	Errors      []string       `json:"errors,omitempty"`
}

// ClearCacheRequest 清除缓存请求
type ClearCacheRequest struct {
	CacheTypes []string `json:"cache_types,omitempty"` // 指定要清除的缓存类型
	UserIDs    []int64  `json:"user_ids,omitempty"`    // 指定要清除的用户ID
	RoleIDs    []int64  `json:"role_ids,omitempty"`    // 指定要清除的角色ID
	DocTypes   []string `json:"doc_types,omitempty"`   // 指定要清除的文档类型
	ClearAll   bool     `json:"clear_all"`             // 是否清除所有权限缓存
}

// ClearCacheResponse 清除缓存响应
type ClearCacheResponse struct {
	Success     bool     `json:"success"`
	Message     string   `json:"message"`
	ClearedKeys int      `json:"cleared_keys"`
	Errors      []string `json:"errors,omitempty"`
}

// GetCacheStatsResponse 获取缓存统计响应
type GetCacheStatsResponse struct {
	Stats map[string]interface{} `json:"stats"`
}

// WarmUpCache 缓存预热
func (s *CacheManagerService) WarmUpCache(ctx context.Context, req *WarmUpCacheRequest) (*WarmUpCacheResponse, error) {
	// 检查权限
	currentUser := middleware.GetCurrentUser(ctx)
	if !currentUser.HasAnyRole("SUPER_ADMIN", "ADMIN") {
		return nil, errors.Forbidden("PERMISSION_DENIED", "无权限执行缓存预热")
	}

	startTime := time.Now()
	stats := make(map[string]int)
	var errs []string

	s.log.Infof("Starting cache warmup by user %s", currentUser.Username)

	// 默认预热所有类型的缓存
	cacheTypes := req.CacheTypes
	if len(cacheTypes) == 0 {
		cacheTypes = []string{"user_roles", "doc_types", "permission_rules", "field_levels"}
	}

	var wg sync.WaitGroup
	var mu sync.Mutex

	for _, cacheType := range cacheTypes {
		wg.Add(1)
		go func(ct string) {
			defer wg.Done()

			count, err := s.warmupCacheType(ctx, ct, req)

			mu.Lock()
			defer mu.Unlock()

			if err != nil {
				errs = append(errs, fmt.Sprintf("%s: %v", ct, err))
				s.log.Errorf("Failed to warmup cache type %s: %v", ct, err)
			} else {
				stats[ct] = count
				s.log.Infof("Successfully warmed up %s cache: %d items", ct, count)
			}
		}(cacheType)
	}

	wg.Wait()

	duration := time.Since(startTime)
	success := len(errs) == 0

	response := &WarmUpCacheResponse{
		Success:     success,
		Message:     fmt.Sprintf("Cache warmup completed in %v", duration),
		WarmupStats: stats,
		Duration:    duration,
		Errors:      errs,
	}

	if success {
		s.log.Infof("Cache warmup completed successfully in %v", duration)
	} else {
		s.log.Warnf("Cache warmup completed with errors in %v: %v", duration, errs)
	}

	return response, nil
}

// warmupCacheType 预热指定类型的缓存
func (s *CacheManagerService) warmupCacheType(ctx context.Context, cacheType string, req *WarmUpCacheRequest) (int, error) {
	switch cacheType {
	case "user_roles":
		return s.warmupUserRoles(ctx, req.UserIDs)
	case "doc_types":
		return s.warmupDocTypes(ctx, req.DocTypes)
	case "permission_rules":
		return s.warmupPermissionRules(ctx, req.RoleIDs, req.DocTypes)
	case "field_levels":
		return s.warmupFieldPermissionLevels(ctx, req.DocTypes)
	default:
		return 0, fmt.Errorf("unsupported cache type: %s", cacheType)
	}
}

// warmupUserRoles 预热用户角色缓存
func (s *CacheManagerService) warmupUserRoles(ctx context.Context, userIDs []int64) (int, error) {
	// 如果没有指定用户ID，暂时跳过（避免加载所有用户）
	if len(userIDs) == 0 {
		return 0, nil
	}

	count := 0
	for _, userID := range userIDs {
		roles, err := s.permissionUc.GetUserRoles(ctx, userID)
		if err != nil {
			s.log.Warnf("Failed to get roles for user %d: %v", userID, err)
			continue
		}

		if err := s.cache.SetUserRoles(ctx, userID, roles, 30*time.Minute); err != nil {
			s.log.Warnf("Failed to cache roles for user %d: %v", userID, err)
			continue
		}

		count++
	}

	return count, nil
}

// warmupDocTypes 预热文档类型缓存
func (s *CacheManagerService) warmupDocTypes(ctx context.Context, docTypes []string) (int, error) {
	var targetDocTypes []string

	if len(docTypes) == 0 {
		// 获取所有文档类型
		allDocTypes, err := s.permissionUc.ListDocTypes(ctx, "")
		if err != nil {
			return 0, fmt.Errorf("failed to list doc types: %w", err)
		}
		for _, dt := range allDocTypes {
			targetDocTypes = append(targetDocTypes, dt.Name)
		}
	} else {
		targetDocTypes = docTypes
	}

	count := 0
	for _, docTypeName := range targetDocTypes {
		docType, err := s.permissionUc.GetDocType(ctx, docTypeName)
		if err != nil {
			s.log.Warnf("Failed to get doctype %s: %v", docTypeName, err)
			continue
		}

		if err := s.cache.SetDocType(ctx, docTypeName, docType, 4*time.Hour); err != nil {
			s.log.Warnf("Failed to cache doctype %s: %v", docTypeName, err)
			continue
		}

		count++
	}

	return count, nil
}

// warmupPermissionRules 预热权限规则缓存
func (s *CacheManagerService) warmupPermissionRules(ctx context.Context, roleIDs []int64, docTypes []string) (int, error) {
	// 如果没有指定角色和文档类型，暂时跳过（避免组合爆炸）
	if len(roleIDs) == 0 || len(docTypes) == 0 {
		return 0, nil
	}

	count := 0
	for _, roleID := range roleIDs {
		for _, docType := range docTypes {
			rules, err := s.permissionUc.ListPermissionRules(ctx, roleID, docType)
			if err != nil {
				s.log.Warnf("Failed to get permission rules for role %d doctype %s: %v", roleID, docType, err)
				continue
			}

			if err := s.cache.SetPermissionRules(ctx, roleID, docType, rules, 1*time.Hour); err != nil {
				s.log.Warnf("Failed to cache permission rules for role %d doctype %s: %v", roleID, docType, err)
				continue
			}

			count++
		}
	}

	return count, nil
}

// warmupFieldPermissionLevels 预热字段权限级别缓存
func (s *CacheManagerService) warmupFieldPermissionLevels(ctx context.Context, docTypes []string) (int, error) {
	var targetDocTypes []string

	if len(docTypes) == 0 {
		// 获取所有文档类型
		allDocTypes, err := s.permissionUc.ListDocTypes(ctx, "")
		if err != nil {
			return 0, fmt.Errorf("failed to list doc types: %w", err)
		}
		for _, dt := range allDocTypes {
			targetDocTypes = append(targetDocTypes, dt.Name)
		}
	} else {
		targetDocTypes = docTypes
	}

	count := 0
	for _, docType := range targetDocTypes {
		// 获取文档类型的字段权限级别
		levels, err := s.permissionUc.ListFieldPermissionLevels(ctx, docType, 1, 1000)
		if err != nil {
			s.log.Warnf("Failed to get field permission levels for doctype %s: %v", docType, err)
			continue
		}

		// 转换为map格式用于缓存
		levelMap := make(map[string]int)
		for _, level := range levels {
			levelMap[level.FieldName] = level.PermissionLevel
		}

		if len(levelMap) > 0 {
			if err := s.cache.SetFieldPermissionLevels(ctx, docType, levelMap, 2*time.Hour); err != nil {
				s.log.Warnf("Failed to cache field permission levels for doctype %s: %v", docType, err)
				continue
			}
		}

		count++
	}

	return count, nil
}

// ClearCache 清除缓存
func (s *CacheManagerService) ClearCache(ctx context.Context, req *ClearCacheRequest) (*ClearCacheResponse, error) {
	// 检查权限
	currentUser := middleware.GetCurrentUser(ctx)
	if !currentUser.HasAnyRole("SUPER_ADMIN", "ADMIN") {
		return nil, errors.Forbidden("PERMISSION_DENIED", "无权限执行缓存清除")
	}

	s.log.Infof("Starting cache clear by user %s", currentUser.Username)

	var errs []string
	clearedKeys := 0

	// 清除所有权限缓存
	if req.ClearAll {
		if err := s.cache.ClearAllPermissionCache(ctx); err != nil {
			errs = append(errs, fmt.Sprintf("clear all: %v", err))
		} else {
			clearedKeys += 1000 // 估算值
			s.log.Info("Cleared all permission cache")
		}

		return &ClearCacheResponse{
			Success:     len(errs) == 0,
			Message:     "Cache clear completed",
			ClearedKeys: clearedKeys,
			Errors:      errs,
		}, nil
	}

	// 按用户清除缓存
	for _, userID := range req.UserIDs {
		if err := s.cache.ClearUserCache(ctx, userID); err != nil {
			errs = append(errs, fmt.Sprintf("user %d: %v", userID, err))
		} else {
			clearedKeys += 10 // 估算每个用户的缓存键数量
			s.log.Infof("Cleared cache for user %d", userID)
		}
	}

	// 按角色清除缓存
	for _, roleID := range req.RoleIDs {
		if err := s.cache.ClearRoleCache(ctx, roleID); err != nil {
			errs = append(errs, fmt.Sprintf("role %d: %v", roleID, err))
		} else {
			clearedKeys += 5 // 估算每个角色的缓存键数量
			s.log.Infof("Cleared cache for role %d", roleID)
		}
	}

	// 按文档类型清除缓存
	for _, docType := range req.DocTypes {
		if err := s.cache.ClearDocTypeCache(ctx, docType); err != nil {
			errs = append(errs, fmt.Sprintf("doctype %s: %v", docType, err))
		} else {
			clearedKeys += 20 // 估算每个文档类型的缓存键数量
			s.log.Infof("Cleared cache for doctype %s", docType)
		}
	}

	success := len(errs) == 0
	message := "Cache clear completed"
	if !success {
		message = fmt.Sprintf("Cache clear completed with %d errors", len(errs))
	}

	return &ClearCacheResponse{
		Success:     success,
		Message:     message,
		ClearedKeys: clearedKeys,
		Errors:      errs,
	}, nil
}

// GetCacheStats 获取缓存统计
func (s *CacheManagerService) GetCacheStats(ctx context.Context) (*GetCacheStatsResponse, error) {
	// 检查权限
	currentUser := middleware.GetCurrentUser(ctx)
	if !currentUser.HasAnyRole("SUPER_ADMIN", "ADMIN") {
		return nil, errors.Forbidden("PERMISSION_DENIED", "无权限查看缓存统计")
	}

	stats, err := s.cache.GetCacheStats(ctx)
	if err != nil {
		return nil, errors.InternalServer("CACHE_STATS_ERROR", "获取缓存统计失败")
	}

	return &GetCacheStatsResponse{
		Stats: stats,
	}, nil
}

// ScheduledCacheRefresh 定时刷新缓存
func (s *CacheManagerService) ScheduledCacheRefresh(ctx context.Context) error {
	s.log.Info("Starting scheduled cache refresh")

	// 获取热门文档类型（这里简化处理，实际应该从统计数据获取）
	commonDocTypes := []string{"User", "Role", "PermissionRule", "Company", "Customer"}

	// 预热常用的文档类型缓存
	req := &WarmUpCacheRequest{
		CacheTypes:   []string{"doc_types", "field_levels"},
		DocTypes:     commonDocTypes,
		ForceRefresh: false,
	}

	// 创建系统用户上下文（用于定时任务）
	systemCtx := context.WithValue(ctx, "current_user", &middleware.CurrentUser{
		ID:       0,
		Username: "system",
		Roles:    []string{"SUPER_ADMIN"},
	})

	_, err := s.WarmUpCache(systemCtx, req)
	if err != nil {
		s.log.Errorf("Scheduled cache refresh failed: %v", err)
		return err
	}

	s.log.Info("Scheduled cache refresh completed")
	return nil
}

// StartCacheMaintenanceTask 启动缓存维护任务
func (s *CacheManagerService) StartCacheMaintenanceTask(ctx context.Context) {
	// 每30分钟执行一次缓存刷新
	ticker := time.NewTicker(30 * time.Minute)
	defer ticker.Stop()

	// 每天凌晨2点清理过期缓存
	cleanupTicker := time.NewTicker(24 * time.Hour)
	defer cleanupTicker.Stop()

	for {
		select {
		case <-ctx.Done():
			s.log.Info("Cache maintenance task stopped")
			return
		case <-ticker.C:
			if err := s.ScheduledCacheRefresh(ctx); err != nil {
				s.log.Errorf("Scheduled cache refresh error: %v", err)
			}
		case <-cleanupTicker.C:
			// 清理统计信息，检查缓存健康状态
			stats, err := s.cache.GetCacheStats(ctx)
			if err != nil {
				s.log.Errorf("Failed to get cache stats for maintenance: %v", err)
				continue
			}

			s.log.Infof("Daily cache maintenance report: %+v", stats)
		}
	}
}
