package cache

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"erp-system/internal/biz"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestRedis(t *testing.T) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       15, // 使用测试数据库
	})

	// 测试连接
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	err := client.Ping(ctx).Err()
	if err != nil {
		t.Skip("Redis not available for testing")
	}

	// 清空测试数据库
	err = client.FlushDB(ctx).Err()
	require.NoError(t, err)

	return client
}

func TestNewRedisPermissionCache(t *testing.T) {
	client := setupTestRedis(t)
	defer client.Close()

	logger := log.DefaultLogger
	cache := NewRedisPermissionCache(client, logger)

	assert.NotNil(t, cache)
	redisCache, ok := cache.(*RedisPermissionCache)
	assert.True(t, ok)
	assert.Equal(t, "perm:", redisCache.prefix)
}

func TestUserPermissionsCache(t *testing.T) {
	client := setupTestRedis(t)
	defer client.Close()

	logger := log.DefaultLogger
	cache := NewRedisPermissionCache(client, logger)
	ctx := context.Background()

	userID := int64(123)
	permissions := []string{"user.create", "user.read", "user.update"}
	ttl := 5 * time.Minute

	// 测试设置用户权限
	err := cache.SetUserPermissions(ctx, userID, permissions, ttl)
	assert.NoError(t, err)

	// 测试获取用户权限
	result, err := cache.GetUserPermissions(ctx, userID)
	assert.NoError(t, err)
	assert.Equal(t, permissions, result)

	// 测试删除用户权限
	err = cache.DeleteUserPermissions(ctx, userID)
	assert.NoError(t, err)

	// 验证权限已删除
	result, err = cache.GetUserPermissions(ctx, userID)
	assert.NoError(t, err)
	assert.Nil(t, result)
}

func TestUserRolesCache(t *testing.T) {
	client := setupTestRedis(t)
	defer client.Close()

	logger := log.DefaultLogger
	cache := NewRedisPermissionCache(client, logger)
	ctx := context.Background()

	userID := int64(456)
	roles := []string{"admin", "manager"}
	ttl := 10 * time.Minute

	// 测试设置用户角色
	err := cache.SetUserRoles(ctx, userID, roles, ttl)
	assert.NoError(t, err)

	// 测试获取用户角色
	result, err := cache.GetUserRoles(ctx, userID)
	assert.NoError(t, err)
	assert.Equal(t, roles, result)

	// 测试删除用户角色
	err = cache.DeleteUserRoles(ctx, userID)
	assert.NoError(t, err)

	// 验证角色已删除
	result, err = cache.GetUserRoles(ctx, userID)
	assert.NoError(t, err)
	assert.Nil(t, result)
}

func TestPermissionRulesCache(t *testing.T) {
	client := setupTestRedis(t)
	defer client.Close()

	logger := log.DefaultLogger
	cache := NewRedisPermissionCache(client, logger)
	ctx := context.Background()

	roleID := int64(789)
	docType := "User"
	rules := []*biz.PermissionRule{
		{
			ID:         1,
			RoleID:     roleID,
			DocType:    docType,
			Permission: "read",
			Rule:       "user == owner",
		},
		{
			ID:         2,
			RoleID:     roleID,
			DocType:    docType,
			Permission: "write",
			Rule:       "user.department == doc.department",
		},
	}
	ttl := 15 * time.Minute

	// 测试设置权限规则
	err := cache.SetPermissionRules(ctx, roleID, docType, rules, ttl)
	assert.NoError(t, err)

	// 测试获取权限规则
	result, err := cache.GetPermissionRules(ctx, roleID, docType)
	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, rules[0].Permission, result[0].Permission)
	assert.Equal(t, rules[1].Rule, result[1].Rule)

	// 测试删除权限规则
	err = cache.DeletePermissionRules(ctx, roleID, docType)
	assert.NoError(t, err)

	// 验证规则已删除
	result, err = cache.GetPermissionRules(ctx, roleID, docType)
	assert.NoError(t, err)
	assert.Nil(t, result)
}

func TestUserPermissionLevelCache(t *testing.T) {
	client := setupTestRedis(t)
	defer client.Close()

	logger := log.DefaultLogger
	cache := NewRedisPermissionCache(client, logger)
	ctx := context.Background()

	userID := int64(111)
	docType := "Document"
	level := 5
	ttl := 20 * time.Minute

	// 测试设置用户权限级别
	err := cache.SetUserPermissionLevel(ctx, userID, docType, level, ttl)
	assert.NoError(t, err)

	// 测试获取用户权限级别
	result, err := cache.GetUserPermissionLevel(ctx, userID, docType)
	assert.NoError(t, err)
	assert.Equal(t, level, result)

	// 测试删除用户权限级别
	err = cache.DeleteUserPermissionLevel(ctx, userID, docType)
	assert.NoError(t, err)

	// 验证级别已删除
	result, err = cache.GetUserPermissionLevel(ctx, userID, docType)
	assert.NoError(t, err)
	assert.Equal(t, -1, result) // 未缓存返回-1
}

func TestFieldPermissionLevelsCache(t *testing.T) {
	client := setupTestRedis(t)
	defer client.Close()

	logger := log.DefaultLogger
	cache := NewRedisPermissionCache(client, logger)
	ctx := context.Background()

	docType := "User"
	levels := map[string]int{
		"name":     1,
		"email":    2,
		"password": 9,
		"salary":   8,
	}
	ttl := 30 * time.Minute

	// 测试设置字段权限级别
	err := cache.SetFieldPermissionLevels(ctx, docType, levels, ttl)
	assert.NoError(t, err)

	// 测试获取字段权限级别
	result, err := cache.GetFieldPermissionLevels(ctx, docType)
	assert.NoError(t, err)
	assert.Equal(t, levels, result)

	// 测试删除字段权限级别
	err = cache.DeleteFieldPermissionLevels(ctx, docType)
	assert.NoError(t, err)

	// 验证级别已删除
	result, err = cache.GetFieldPermissionLevels(ctx, docType)
	assert.NoError(t, err)
	assert.Nil(t, result)
}

func TestDocTypeCache(t *testing.T) {
	client := setupTestRedis(t)
	defer client.Close()

	logger := log.DefaultLogger
	cache := NewRedisPermissionCache(client, logger)
	ctx := context.Background()

	name := "TestDocType"
	docType := &biz.DocType{
		ID:          1,
		Name:        name,
		DisplayName: "Test Document Type",
		Description: "Test description",
		Fields: []biz.DocField{
			{
				Name:        "title",
				DisplayName: "Title",
				Type:        "string",
				Required:    true,
				ReadLevel:   1,
				WriteLevel:  2,
			},
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	ttl := 25 * time.Minute

	// 测试设置文档类型
	err := cache.SetDocType(ctx, name, docType, ttl)
	assert.NoError(t, err)

	// 测试获取文档类型
	result, err := cache.GetDocType(ctx, name)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, docType.Name, result.Name)
	assert.Equal(t, docType.DisplayName, result.DisplayName)
	assert.Len(t, result.Fields, 1)
	assert.Equal(t, "title", result.Fields[0].Name)

	// 测试删除文档类型
	err = cache.DeleteDocType(ctx, name)
	assert.NoError(t, err)

	// 验证类型已删除
	result, err = cache.GetDocType(ctx, name)
	assert.NoError(t, err)
	assert.Nil(t, result)
}

func TestClearUserCache(t *testing.T) {
	client := setupTestRedis(t)
	defer client.Close()

	logger := log.DefaultLogger
	cache := NewRedisPermissionCache(client, logger)
	ctx := context.Background()

	userID := int64(999)
	
	// 设置用户相关的多个缓存
	err := cache.SetUserPermissions(ctx, userID, []string{"perm1"}, time.Hour)
	assert.NoError(t, err)
	
	err = cache.SetUserRoles(ctx, userID, []string{"role1"}, time.Hour)
	assert.NoError(t, err)
	
	err = cache.SetUserPermissionLevel(ctx, userID, "TestDoc", 5, time.Hour)
	assert.NoError(t, err)

	// 清除用户缓存
	err = cache.ClearUserCache(ctx, userID)
	assert.NoError(t, err)

	// 验证所有用户相关缓存已清除
	perms, err := cache.GetUserPermissions(ctx, userID)
	assert.NoError(t, err)
	assert.Nil(t, perms)

	roles, err := cache.GetUserRoles(ctx, userID)
	assert.NoError(t, err)
	assert.Nil(t, roles)

	level, err := cache.GetUserPermissionLevel(ctx, userID, "TestDoc")
	assert.NoError(t, err)
	assert.Equal(t, -1, level)
}

func TestCacheStats(t *testing.T) {
	client := setupTestRedis(t)
	defer client.Close()

	logger := log.DefaultLogger
	cache := NewRedisPermissionCache(client, logger)
	ctx := context.Background()

	// 设置一些测试数据
	err := cache.SetUserPermissions(ctx, 1, []string{"perm1"}, time.Hour)
	assert.NoError(t, err)
	
	err = cache.SetUserRoles(ctx, 1, []string{"role1"}, time.Hour)
	assert.NoError(t, err)

	// 获取缓存统计
	stats, err := cache.GetCacheStats(ctx)
	assert.NoError(t, err)
	assert.NotNil(t, stats)
	
	// 验证统计信息包含预期字段
	assert.Contains(t, stats, "user_permissions")
	assert.Contains(t, stats, "user_roles")
	assert.Contains(t, stats, "permission_rules")
	assert.Contains(t, stats, "user_levels")
	assert.Contains(t, stats, "field_levels")
	assert.Contains(t, stats, "doctypes")

	// 验证计数正确
	userPermsCount, ok := stats["user_permissions"].(int)
	assert.True(t, ok)
	assert.GreaterOrEqual(t, userPermsCount, 1)

	userRolesCount, ok := stats["user_roles"].(int)
	assert.True(t, ok)
	assert.GreaterOrEqual(t, userRolesCount, 1)
}

func TestCacheKeyGeneration(t *testing.T) {
	client := setupTestRedis(t)
	defer client.Close()

	logger := log.DefaultLogger
	redisCache := NewRedisPermissionCache(client, logger).(*RedisPermissionCache)

	// 测试缓存键生成
	userID := int64(123)
	roleID := int64(456)
	docType := "TestDoc"

	userPermsKey := redisCache.userPermissionsKey(userID)
	assert.Equal(t, "perm:user_perms:123", userPermsKey)

	userRolesKey := redisCache.userRolesKey(userID)
	assert.Equal(t, "perm:user_roles:123", userRolesKey)

	rulesKey := redisCache.permissionRulesKey(roleID, docType)
	assert.Equal(t, "perm:role_rules:456:TestDoc", rulesKey)

	levelKey := redisCache.userPermissionLevelKey(userID, docType)
	assert.Equal(t, "perm:user_level:123:TestDoc", levelKey)

	fieldLevelsKey := redisCache.fieldPermissionLevelsKey(docType)
	assert.Equal(t, "perm:field_levels:TestDoc", fieldLevelsKey)

	docTypeKey := redisCache.docTypeKey("MyDoc")
	assert.Equal(t, "perm:doctype:MyDoc", docTypeKey)
}