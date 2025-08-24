package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"erp-system/internal/biz"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-redis/redis/v8"
)

// PermissionCache 权限缓存接口
type PermissionCache interface {
	// 用户权限缓存
	SetUserPermissions(ctx context.Context, userID int64, permissions []string, ttl time.Duration) error
	GetUserPermissions(ctx context.Context, userID int64) ([]string, error)
	DeleteUserPermissions(ctx context.Context, userID int64) error

	// 用户角色缓存
	SetUserRoles(ctx context.Context, userID int64, roles []string, ttl time.Duration) error
	GetUserRoles(ctx context.Context, userID int64) ([]string, error)
	DeleteUserRoles(ctx context.Context, userID int64) error

	// 权限规则缓存
	SetPermissionRules(ctx context.Context, roleID int64, docType string, rules []*biz.PermissionRule, ttl time.Duration) error
	GetPermissionRules(ctx context.Context, roleID int64, docType string) ([]*biz.PermissionRule, error)
	DeletePermissionRules(ctx context.Context, roleID int64, docType string) error

	// 用户权限级别缓存
	SetUserPermissionLevel(ctx context.Context, userID int64, docType string, level int, ttl time.Duration) error
	GetUserPermissionLevel(ctx context.Context, userID int64, docType string) (int, error)
	DeleteUserPermissionLevel(ctx context.Context, userID int64, docType string) error

	// 字段权限级别缓存
	SetFieldPermissionLevels(ctx context.Context, docType string, levels map[string]int, ttl time.Duration) error
	GetFieldPermissionLevels(ctx context.Context, docType string) (map[string]int, error)
	DeleteFieldPermissionLevels(ctx context.Context, docType string) error

	// 文档类型缓存
	SetDocType(ctx context.Context, name string, docType *biz.DocType, ttl time.Duration) error
	GetDocType(ctx context.Context, name string) (*biz.DocType, error)
	DeleteDocType(ctx context.Context, name string) error

	// 批量清除缓存
	ClearUserCache(ctx context.Context, userID int64) error
	ClearRoleCache(ctx context.Context, roleID int64) error
	ClearDocTypeCache(ctx context.Context, docType string) error
	ClearAllPermissionCache(ctx context.Context) error

	// 缓存统计
	GetCacheStats(ctx context.Context) (map[string]interface{}, error)
}

// RedisPermissionCache Redis权限缓存实现
type RedisPermissionCache struct {
	client *redis.Client
	logger *log.Helper
	prefix string
}

// NewRedisPermissionCache 创建Redis权限缓存
func NewRedisPermissionCache(client *redis.Client, logger log.Logger) PermissionCache {
	return &RedisPermissionCache{
		client: client,
		logger: log.NewHelper(logger),
		prefix: "perm:",
	}
}

// 缓存键生成函数
func (c *RedisPermissionCache) userPermissionsKey(userID int64) string {
	return fmt.Sprintf("%suser_perms:%d", c.prefix, userID)
}

func (c *RedisPermissionCache) userRolesKey(userID int64) string {
	return fmt.Sprintf("%suser_roles:%d", c.prefix, userID)
}

func (c *RedisPermissionCache) permissionRulesKey(roleID int64, docType string) string {
	return fmt.Sprintf("%srole_rules:%d:%s", c.prefix, roleID, docType)
}

func (c *RedisPermissionCache) userPermissionLevelKey(userID int64, docType string) string {
	return fmt.Sprintf("%suser_level:%d:%s", c.prefix, userID, docType)
}

func (c *RedisPermissionCache) fieldPermissionLevelsKey(docType string) string {
	return fmt.Sprintf("%sfield_levels:%s", c.prefix, docType)
}

func (c *RedisPermissionCache) docTypeKey(name string) string {
	return fmt.Sprintf("%sdoctype:%s", c.prefix, name)
}

// 用户权限缓存实现
func (c *RedisPermissionCache) SetUserPermissions(ctx context.Context, userID int64, permissions []string, ttl time.Duration) error {
	key := c.userPermissionsKey(userID)
	data, err := json.Marshal(permissions)
	if err != nil {
		return fmt.Errorf("failed to marshal user permissions: %w", err)
	}

	return c.client.Set(ctx, key, data, ttl).Err()
}

func (c *RedisPermissionCache) GetUserPermissions(ctx context.Context, userID int64) ([]string, error) {
	key := c.userPermissionsKey(userID)
	data, err := c.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil // 缓存未命中
		}
		return nil, fmt.Errorf("failed to get user permissions from cache: %w", err)
	}

	var permissions []string
	if err := json.Unmarshal([]byte(data), &permissions); err != nil {
		return nil, fmt.Errorf("failed to unmarshal user permissions: %w", err)
	}

	return permissions, nil
}

func (c *RedisPermissionCache) DeleteUserPermissions(ctx context.Context, userID int64) error {
	key := c.userPermissionsKey(userID)
	return c.client.Del(ctx, key).Err()
}

// 用户角色缓存实现
func (c *RedisPermissionCache) SetUserRoles(ctx context.Context, userID int64, roles []string, ttl time.Duration) error {
	key := c.userRolesKey(userID)
	data, err := json.Marshal(roles)
	if err != nil {
		return fmt.Errorf("failed to marshal user roles: %w", err)
	}

	return c.client.Set(ctx, key, data, ttl).Err()
}

func (c *RedisPermissionCache) GetUserRoles(ctx context.Context, userID int64) ([]string, error) {
	key := c.userRolesKey(userID)
	data, err := c.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil // 缓存未命中
		}
		return nil, fmt.Errorf("failed to get user roles from cache: %w", err)
	}

	var roles []string
	if err := json.Unmarshal([]byte(data), &roles); err != nil {
		return nil, fmt.Errorf("failed to unmarshal user roles: %w", err)
	}

	return roles, nil
}

func (c *RedisPermissionCache) DeleteUserRoles(ctx context.Context, userID int64) error {
	key := c.userRolesKey(userID)
	return c.client.Del(ctx, key).Err()
}

// 权限规则缓存实现
func (c *RedisPermissionCache) SetPermissionRules(ctx context.Context, roleID int64, docType string, rules []*biz.PermissionRule, ttl time.Duration) error {
	key := c.permissionRulesKey(roleID, docType)
	data, err := json.Marshal(rules)
	if err != nil {
		return fmt.Errorf("failed to marshal permission rules: %w", err)
	}

	return c.client.Set(ctx, key, data, ttl).Err()
}

func (c *RedisPermissionCache) GetPermissionRules(ctx context.Context, roleID int64, docType string) ([]*biz.PermissionRule, error) {
	key := c.permissionRulesKey(roleID, docType)
	data, err := c.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil // 缓存未命中
		}
		return nil, fmt.Errorf("failed to get permission rules from cache: %w", err)
	}

	var rules []*biz.PermissionRule
	if err := json.Unmarshal([]byte(data), &rules); err != nil {
		return nil, fmt.Errorf("failed to unmarshal permission rules: %w", err)
	}

	return rules, nil
}

func (c *RedisPermissionCache) DeletePermissionRules(ctx context.Context, roleID int64, docType string) error {
	key := c.permissionRulesKey(roleID, docType)
	return c.client.Del(ctx, key).Err()
}

// 用户权限级别缓存实现
func (c *RedisPermissionCache) SetUserPermissionLevel(ctx context.Context, userID int64, docType string, level int, ttl time.Duration) error {
	key := c.userPermissionLevelKey(userID, docType)
	return c.client.Set(ctx, key, level, ttl).Err()
}

func (c *RedisPermissionCache) GetUserPermissionLevel(ctx context.Context, userID int64, docType string) (int, error) {
	key := c.userPermissionLevelKey(userID, docType)
	result, err := c.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return -1, nil // 缓存未命中，返回-1表示未缓存
		}
		return -1, fmt.Errorf("failed to get user permission level from cache: %w", err)
	}

	level, err := strconv.Atoi(result)
	if err != nil {
		return -1, fmt.Errorf("failed to parse permission level: %w", err)
	}

	return level, nil
}

func (c *RedisPermissionCache) DeleteUserPermissionLevel(ctx context.Context, userID int64, docType string) error {
	key := c.userPermissionLevelKey(userID, docType)
	return c.client.Del(ctx, key).Err()
}

// 字段权限级别缓存实现
func (c *RedisPermissionCache) SetFieldPermissionLevels(ctx context.Context, docType string, levels map[string]int, ttl time.Duration) error {
	key := c.fieldPermissionLevelsKey(docType)
	data, err := json.Marshal(levels)
	if err != nil {
		return fmt.Errorf("failed to marshal field permission levels: %w", err)
	}

	return c.client.Set(ctx, key, data, ttl).Err()
}

func (c *RedisPermissionCache) GetFieldPermissionLevels(ctx context.Context, docType string) (map[string]int, error) {
	key := c.fieldPermissionLevelsKey(docType)
	data, err := c.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil // 缓存未命中
		}
		return nil, fmt.Errorf("failed to get field permission levels from cache: %w", err)
	}

	var levels map[string]int
	if err := json.Unmarshal([]byte(data), &levels); err != nil {
		return nil, fmt.Errorf("failed to unmarshal field permission levels: %w", err)
	}

	return levels, nil
}

func (c *RedisPermissionCache) DeleteFieldPermissionLevels(ctx context.Context, docType string) error {
	key := c.fieldPermissionLevelsKey(docType)
	return c.client.Del(ctx, key).Err()
}

// 文档类型缓存实现
func (c *RedisPermissionCache) SetDocType(ctx context.Context, name string, docType *biz.DocType, ttl time.Duration) error {
	key := c.docTypeKey(name)
	data, err := json.Marshal(docType)
	if err != nil {
		return fmt.Errorf("failed to marshal doctype: %w", err)
	}

	return c.client.Set(ctx, key, data, ttl).Err()
}

func (c *RedisPermissionCache) GetDocType(ctx context.Context, name string) (*biz.DocType, error) {
	key := c.docTypeKey(name)
	data, err := c.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil // 缓存未命中
		}
		return nil, fmt.Errorf("failed to get doctype from cache: %w", err)
	}

	var docType biz.DocType
	if err := json.Unmarshal([]byte(data), &docType); err != nil {
		return nil, fmt.Errorf("failed to unmarshal doctype: %w", err)
	}

	return &docType, nil
}

func (c *RedisPermissionCache) DeleteDocType(ctx context.Context, name string) error {
	key := c.docTypeKey(name)
	return c.client.Del(ctx, key).Err()
}

// 批量清除缓存实现
func (c *RedisPermissionCache) ClearUserCache(ctx context.Context, userID int64) error {
	// 清除用户相关的所有缓存
	pattern := fmt.Sprintf("%suser_*:%d*", c.prefix, userID)
	keys, err := c.client.Keys(ctx, pattern).Result()
	if err != nil {
		return fmt.Errorf("failed to find user cache keys: %w", err)
	}

	if len(keys) > 0 {
		return c.client.Del(ctx, keys...).Err()
	}

	return nil
}

func (c *RedisPermissionCache) ClearRoleCache(ctx context.Context, roleID int64) error {
	// 清除角色相关的所有缓存
	pattern := fmt.Sprintf("%srole_*:%d:*", c.prefix, roleID)
	keys, err := c.client.Keys(ctx, pattern).Result()
	if err != nil {
		return fmt.Errorf("failed to find role cache keys: %w", err)
	}

	if len(keys) > 0 {
		return c.client.Del(ctx, keys...).Err()
	}

	return nil
}

func (c *RedisPermissionCache) ClearDocTypeCache(ctx context.Context, docType string) error {
	// 清除文档类型相关的所有缓存
	patterns := []string{
		fmt.Sprintf("%sdoctype:%s", c.prefix, docType),
		fmt.Sprintf("%sfield_levels:%s", c.prefix, docType),
		fmt.Sprintf("%s*:%s", c.prefix, docType),
		fmt.Sprintf("%s*:*:%s", c.prefix, docType),
	}

	var allKeys []string
	for _, pattern := range patterns {
		keys, err := c.client.Keys(ctx, pattern).Result()
		if err != nil {
			c.logger.Warnf("Failed to find cache keys with pattern %s: %v", pattern, err)
			continue
		}
		allKeys = append(allKeys, keys...)
	}

	if len(allKeys) > 0 {
		return c.client.Del(ctx, allKeys...).Err()
	}

	return nil
}

func (c *RedisPermissionCache) ClearAllPermissionCache(ctx context.Context) error {
	// 清除所有权限相关缓存
	pattern := fmt.Sprintf("%s*", c.prefix)
	keys, err := c.client.Keys(ctx, pattern).Result()
	if err != nil {
		return fmt.Errorf("failed to find all permission cache keys: %w", err)
	}

	if len(keys) > 0 {
		// 分批删除，避免单次删除过多键值导致Redis阻塞
		batchSize := 100
		for i := 0; i < len(keys); i += batchSize {
			end := i + batchSize
			if end > len(keys) {
				end = len(keys)
			}

			if err := c.client.Del(ctx, keys[i:end]...).Err(); err != nil {
				c.logger.Errorf("Failed to delete cache keys batch %d-%d: %v", i, end-1, err)
				return err
			}
		}
	}

	return nil
}

// 缓存统计和监控方法
func (c *RedisPermissionCache) GetCacheStats(ctx context.Context) (map[string]interface{}, error) {
	stats := make(map[string]interface{})

	// 获取各类缓存的键数量
	patterns := map[string]string{
		"user_permissions": fmt.Sprintf("%suser_perms:*", c.prefix),
		"user_roles":       fmt.Sprintf("%suser_roles:*", c.prefix),
		"permission_rules": fmt.Sprintf("%srole_rules:*", c.prefix),
		"user_levels":      fmt.Sprintf("%suser_level:*", c.prefix),
		"field_levels":     fmt.Sprintf("%sfield_levels:*", c.prefix),
		"doctypes":         fmt.Sprintf("%sdoctype:*", c.prefix),
	}

	for name, pattern := range patterns {
		keys, err := c.client.Keys(ctx, pattern).Result()
		if err != nil {
			c.logger.Warnf("Failed to count keys for pattern %s: %v", pattern, err)
			stats[name] = -1
		} else {
			stats[name] = len(keys)
		}
	}

	// 获取Redis内存使用情况
	info, err := c.client.Info(ctx, "memory").Result()
	if err == nil {
		lines := strings.Split(info, "\r\n")
		for _, line := range lines {
			if strings.HasPrefix(line, "used_memory:") {
				stats["redis_memory_used"] = strings.TrimPrefix(line, "used_memory:")
			} else if strings.HasPrefix(line, "used_memory_human:") {
				stats["redis_memory_human"] = strings.TrimPrefix(line, "used_memory_human:")
			}
		}
	}

	return stats, nil
}
