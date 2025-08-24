package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"erp-system/internal/biz"

	"github.com/go-kratos/kratos/v2/log"
)

// MemoryPermissionCache 内存权限缓存实现
type MemoryPermissionCache struct {
	cache  sync.Map
	logger *log.Helper
	prefix string
}

// NewMemoryPermissionCache 创建内存权限缓存
func NewMemoryPermissionCache(logger log.Logger) PermissionCache {
	return &MemoryPermissionCache{
		cache:  sync.Map{},
		logger: log.NewHelper(logger),
		prefix: "perm:",
	}
}

// 缓存键生成函数
func (c *MemoryPermissionCache) userPermissionsKey(userID int64) string {
	return fmt.Sprintf("%suser_perms:%d", c.prefix, userID)
}

func (c *MemoryPermissionCache) userRolesKey(userID int64) string {
	return fmt.Sprintf("%suser_roles:%d", c.prefix, userID)
}

func (c *MemoryPermissionCache) permissionRulesKey(roleID int64, docType string) string {
	return fmt.Sprintf("%srole_rules:%d:%s", c.prefix, roleID, docType)
}

func (c *MemoryPermissionCache) userPermissionLevelKey(userID int64, docType string) string {
	return fmt.Sprintf("%suser_level:%d:%s", c.prefix, userID, docType)
}

func (c *MemoryPermissionCache) fieldPermissionLevelsKey(docType string) string {
	return fmt.Sprintf("%sfield_levels:%s", c.prefix, docType)
}

func (c *MemoryPermissionCache) docTypeKey(name string) string {
	return fmt.Sprintf("%sdoctype:%s", c.prefix, name)
}

// 用户权限缓存实现
func (c *MemoryPermissionCache) SetUserPermissions(ctx context.Context, userID int64, permissions []string, ttl time.Duration) error {
	key := c.userPermissionsKey(userID)
	data, err := json.Marshal(permissions)
	if err != nil {
		return fmt.Errorf("failed to marshal user permissions: %w", err)
	}
	
	c.cache.Store(key, string(data))
	return nil
}

func (c *MemoryPermissionCache) GetUserPermissions(ctx context.Context, userID int64) ([]string, error) {
	key := c.userPermissionsKey(userID)
	value, ok := c.cache.Load(key)
	if !ok {
		return nil, nil // 缓存未命中
	}
	
	data, ok := value.(string)
	if !ok {
		return nil, fmt.Errorf("invalid cached data type for user permissions")
	}
	
	var permissions []string
	if err := json.Unmarshal([]byte(data), &permissions); err != nil {
		return nil, fmt.Errorf("failed to unmarshal user permissions: %w", err)
	}
	
	return permissions, nil
}

func (c *MemoryPermissionCache) DeleteUserPermissions(ctx context.Context, userID int64) error {
	key := c.userPermissionsKey(userID)
	c.cache.Delete(key)
	return nil
}

// 用户角色缓存实现
func (c *MemoryPermissionCache) SetUserRoles(ctx context.Context, userID int64, roles []string, ttl time.Duration) error {
	key := c.userRolesKey(userID)
	data, err := json.Marshal(roles)
	if err != nil {
		return fmt.Errorf("failed to marshal user roles: %w", err)
	}
	
	c.cache.Store(key, string(data))
	return nil
}

func (c *MemoryPermissionCache) GetUserRoles(ctx context.Context, userID int64) ([]string, error) {
	key := c.userRolesKey(userID)
	value, ok := c.cache.Load(key)
	if !ok {
		return nil, nil // 缓存未命中
	}
	
	data, ok := value.(string)
	if !ok {
		return nil, fmt.Errorf("invalid cached data type for user roles")
	}
	
	var roles []string
	if err := json.Unmarshal([]byte(data), &roles); err != nil {
		return nil, fmt.Errorf("failed to unmarshal user roles: %w", err)
	}
	
	return roles, nil
}

func (c *MemoryPermissionCache) DeleteUserRoles(ctx context.Context, userID int64) error {
	key := c.userRolesKey(userID)
	c.cache.Delete(key)
	return nil
}

// 权限规则缓存实现
func (c *MemoryPermissionCache) SetPermissionRules(ctx context.Context, roleID int64, docType string, rules []*biz.PermissionRule, ttl time.Duration) error {
	key := c.permissionRulesKey(roleID, docType)
	data, err := json.Marshal(rules)
	if err != nil {
		return fmt.Errorf("failed to marshal permission rules: %w", err)
	}
	
	c.cache.Store(key, string(data))
	return nil
}

func (c *MemoryPermissionCache) GetPermissionRules(ctx context.Context, roleID int64, docType string) ([]*biz.PermissionRule, error) {
	key := c.permissionRulesKey(roleID, docType)
	value, ok := c.cache.Load(key)
	if !ok {
		return nil, nil // 缓存未命中
	}
	
	data, ok := value.(string)
	if !ok {
		return nil, fmt.Errorf("invalid cached data type for permission rules")
	}
	
	var rules []*biz.PermissionRule
	if err := json.Unmarshal([]byte(data), &rules); err != nil {
		return nil, fmt.Errorf("failed to unmarshal permission rules: %w", err)
	}
	
	return rules, nil
}

func (c *MemoryPermissionCache) DeletePermissionRules(ctx context.Context, roleID int64, docType string) error {
	key := c.permissionRulesKey(roleID, docType)
	c.cache.Delete(key)
	return nil
}

// 用户权限级别缓存实现
func (c *MemoryPermissionCache) SetUserPermissionLevel(ctx context.Context, userID int64, docType string, level int, ttl time.Duration) error {
	key := c.userPermissionLevelKey(userID, docType)
	c.cache.Store(key, level)
	return nil
}

func (c *MemoryPermissionCache) GetUserPermissionLevel(ctx context.Context, userID int64, docType string) (int, error) {
	key := c.userPermissionLevelKey(userID, docType)
	value, ok := c.cache.Load(key)
	if !ok {
		return -1, nil // 缓存未命中，返回-1表示未缓存
	}
	
	level, ok := value.(int)
	if !ok {
		return -1, fmt.Errorf("invalid cached data type for user permission level")
	}
	
	return level, nil
}

func (c *MemoryPermissionCache) DeleteUserPermissionLevel(ctx context.Context, userID int64, docType string) error {
	key := c.userPermissionLevelKey(userID, docType)
	c.cache.Delete(key)
	return nil
}

// 字段权限级别缓存实现
func (c *MemoryPermissionCache) SetFieldPermissionLevels(ctx context.Context, docType string, levels map[string]int, ttl time.Duration) error {
	key := c.fieldPermissionLevelsKey(docType)
	data, err := json.Marshal(levels)
	if err != nil {
		return fmt.Errorf("failed to marshal field permission levels: %w", err)
	}
	
	c.cache.Store(key, string(data))
	return nil
}

func (c *MemoryPermissionCache) GetFieldPermissionLevels(ctx context.Context, docType string) (map[string]int, error) {
	key := c.fieldPermissionLevelsKey(docType)
	value, ok := c.cache.Load(key)
	if !ok {
		return nil, nil // 缓存未命中
	}
	
	data, ok := value.(string)
	if !ok {
		return nil, fmt.Errorf("invalid cached data type for field permission levels")
	}
	
	var levels map[string]int
	if err := json.Unmarshal([]byte(data), &levels); err != nil {
		return nil, fmt.Errorf("failed to unmarshal field permission levels: %w", err)
	}
	
	return levels, nil
}

func (c *MemoryPermissionCache) DeleteFieldPermissionLevels(ctx context.Context, docType string) error {
	key := c.fieldPermissionLevelsKey(docType)
	c.cache.Delete(key)
	return nil
}

// 文档类型缓存实现
func (c *MemoryPermissionCache) SetDocType(ctx context.Context, name string, docType *biz.DocType, ttl time.Duration) error {
	key := c.docTypeKey(name)
	data, err := json.Marshal(docType)
	if err != nil {
		return fmt.Errorf("failed to marshal doctype: %w", err)
	}
	
	c.cache.Store(key, string(data))
	return nil
}

func (c *MemoryPermissionCache) GetDocType(ctx context.Context, name string) (*biz.DocType, error) {
	key := c.docTypeKey(name)
	value, ok := c.cache.Load(key)
	if !ok {
		return nil, nil // 缓存未命中
	}
	
	data, ok := value.(string)
	if !ok {
		return nil, fmt.Errorf("invalid cached data type for doctype")
	}
	
	var docType biz.DocType
	if err := json.Unmarshal([]byte(data), &docType); err != nil {
		return nil, fmt.Errorf("failed to unmarshal doctype: %w", err)
	}
	
	return &docType, nil
}

func (c *MemoryPermissionCache) DeleteDocType(ctx context.Context, name string) error {
	key := c.docTypeKey(name)
	c.cache.Delete(key)
	return nil
}

// 批量清除缓存实现
func (c *MemoryPermissionCache) ClearUserCache(ctx context.Context, userID int64) error {
	// 清除用户相关的所有缓存
	pattern := fmt.Sprintf("%suser_", c.prefix)
	c.cache.Range(func(key, value interface{}) bool {
		if k, ok := key.(string); ok {
			if len(k) > len(pattern) && k[:len(pattern)] == pattern {
				c.cache.Delete(key)
			}
		}
		return true
	})
	
	return nil
}

func (c *MemoryPermissionCache) ClearRoleCache(ctx context.Context, roleID int64) error {
	// 清除角色相关的所有缓存
	pattern := fmt.Sprintf("%srole_", c.prefix)
	c.cache.Range(func(key, value interface{}) bool {
		if k, ok := key.(string); ok {
			if len(k) > len(pattern) && k[:len(pattern)] == pattern {
				c.cache.Delete(key)
			}
		}
		return true
	})
	
	return nil
}

func (c *MemoryPermissionCache) ClearDocTypeCache(ctx context.Context, docType string) error {
	// 清除文档类型相关的所有缓存
	c.cache.Range(func(key, value interface{}) bool {
		if k, ok := key.(string); ok {
			if k == c.docTypeKey(docType) || k == c.fieldPermissionLevelsKey(docType) {
				c.cache.Delete(key)
			}
		}
		return true
	})
	
	return nil
}

func (c *MemoryPermissionCache) ClearAllPermissionCache(ctx context.Context) error {
	// 清除所有权限相关缓存
	c.cache.Range(func(key, value interface{}) bool {
		if k, ok := key.(string); ok {
			if len(k) > len(c.prefix) && k[:len(c.prefix)] == c.prefix {
				c.cache.Delete(key)
			}
		}
		return true
	})
	
	return nil
}

// 缓存统计和监控方法
func (c *MemoryPermissionCache) GetCacheStats(ctx context.Context) (map[string]interface{}, error) {
	stats := make(map[string]interface{})
	
	// 统计各类缓存的数量
	counts := map[string]int{
		"user_permissions": 0,
		"user_roles":       0,
		"permission_rules": 0,
		"user_levels":      0,
		"field_levels":     0,
		"doctypes":         0,
		"total":           0,
	}
	
	c.cache.Range(func(key, value interface{}) bool {
		if k, ok := key.(string); ok {
			counts["total"]++
			switch {
			case len(k) > len(c.prefix+"user_perms:") && k[:len(c.prefix+"user_perms:")] == c.prefix+"user_perms:":
				counts["user_permissions"]++
			case len(k) > len(c.prefix+"user_roles:") && k[:len(c.prefix+"user_roles:")] == c.prefix+"user_roles:":
				counts["user_roles"]++
			case len(k) > len(c.prefix+"role_rules:") && k[:len(c.prefix+"role_rules:")] == c.prefix+"role_rules:":
				counts["permission_rules"]++
			case len(k) > len(c.prefix+"user_level:") && k[:len(c.prefix+"user_level:")] == c.prefix+"user_level:":
				counts["user_levels"]++
			case len(k) > len(c.prefix+"field_levels:") && k[:len(c.prefix+"field_levels:")] == c.prefix+"field_levels:":
				counts["field_levels"]++
			case len(k) > len(c.prefix+"doctype:") && k[:len(c.prefix+"doctype:")] == c.prefix+"doctype:":
				counts["doctypes"]++
			}
		}
		return true
	})
	
	for k, v := range counts {
		stats[k] = v
	}
	
	stats["cache_type"] = "memory"
	
	return stats, nil
}