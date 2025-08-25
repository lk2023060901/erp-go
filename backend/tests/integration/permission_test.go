package integration

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"erp-system/internal/biz"
	"erp-system/internal/cache"
	"erp-system/internal/conf"
	"erp-system/internal/data"
	"erp-system/internal/middleware"
	"erp-system/internal/service"

	"github.com/go-kratos/kratos/v2/log"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// getProjectRoot 获取项目根目录路径
func getProjectRoot() string {
	// 获取当前文件的绝对路径
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		// 如果无法获取当前文件路径，使用工作目录
		wd, err := os.Getwd()
		if err != nil {
			return "."
		}
		// 从当前工作目录向上查找项目根目录
		for {
			if _, err := os.Stat(filepath.Join(wd, "go.mod")); err == nil {
				return wd
			}
			parent := filepath.Dir(wd)
			if parent == wd {
				break
			}
			wd = parent
		}
		return "."
	}
	
	// 从当前文件路径向上查找项目根目录
	dir := filepath.Dir(filename)
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	
	return "."
}

// TestPermissionSystemIntegration 权限系统集成测试
func TestPermissionSystemIntegration(t *testing.T) {
	// 设置测试数据库
	db, dbPath, cleanup := setupTestDatabase(t)
	defer cleanup()

	// 插入测试基础数据
	setupTestBaseData(t, db)

	// 创建依赖
	logger := log.DefaultLogger
	permissionCache := cache.NewMemoryPermissionCache(logger)

	// 使用相同数据库路径创建 Data 实例
	dataInstance, cleanup2, err := data.NewData(&conf.Data{
		Database: &conf.Database{Driver: "sqlite3", Source: dbPath},
		Redis:    &conf.Redis{Addr: "localhost:6379", Password: "", DB: 0},
	}, logger)
	require.NoError(t, err)
	defer cleanup2()

	// 创建Repository
	permissionRepo := data.NewPermissionRepo(dataInstance, logger)
	cachedPermissionRepo := data.NewCachedPermissionRepo(permissionRepo, permissionCache, logger)

	// 创建Usecase
	permissionUsecase := biz.NewPermissionUsecase(cachedPermissionRepo, logger)

	// 创建Service
	permissionService := service.NewPermissionService(permissionUsecase, logger)

	// 创建带有超级管理员权限的测试上下文
	ctx := context.Background()
	ctx = middleware.SetUserIDToContext(ctx, 1)
	ctx = middleware.SetUsernameToContext(ctx, "testadmin")
	ctx = middleware.SetUserEmailToContext(ctx, "admin@test.com")
	ctx = middleware.SetUserRolesToContext(ctx, []string{"SUPER_ADMIN"})

	// 运行集成测试套件
	t.Run("DocType Management", func(t *testing.T) {
		testDocTypeManagement(t, ctx, permissionService)
	})

	t.Run("Permission Rule Management", func(t *testing.T) {
		testPermissionRuleManagement(t, ctx, permissionService, db)
	})

	t.Run("User Permission Management", func(t *testing.T) {
		testUserPermissionManagement(t, ctx, permissionService, db)
	})

	t.Run("Permission Checking", func(t *testing.T) {
		testPermissionChecking(t, ctx, permissionService, db)
	})
}

// setupTestDatabase 设置测试数据库
func setupTestDatabase(t *testing.T) (*sql.DB, string, func()) {
	// 创建临时数据库文件
	dbPath := fmt.Sprintf("/tmp/test_permission_%d.db", time.Now().Unix())

	db, err := sql.Open("sqlite3", dbPath)
	require.NoError(t, err)

	// 获取项目根路径
	projectRoot := getProjectRoot()

	// 先手动执行核心系统表结构 SQL
	coreSchemaSQL, err := os.ReadFile(filepath.Join(projectRoot, "database", "schema", "core_system.sql"))
	require.NoError(t, err)

	_, err = db.Exec(string(coreSchemaSQL))
	require.NoError(t, err)

	// 执行权限系统表结构 SQL
	permissionSchemaSQL, err := os.ReadFile(filepath.Join(projectRoot, "database", "schema", "permission_system.sql"))
	require.NoError(t, err)

	_, err = db.Exec(string(permissionSchemaSQL))
	require.NoError(t, err)

	// 执行核心系统种子数据 SQL
	coreSeedSQL, err := os.ReadFile(filepath.Join(projectRoot, "database", "data", "core_seed.sql"))
	require.NoError(t, err)

	_, err = db.Exec(string(coreSeedSQL))
	require.NoError(t, err)

	// 执行权限系统种子数据 SQL
	permissionSeedSQL, err := os.ReadFile(filepath.Join(projectRoot, "database", "data", "permission_seed.sql"))
	require.NoError(t, err)

	_, err = db.Exec(string(permissionSeedSQL))
	require.NoError(t, err)

	cleanup := func() {
		db.Close()
		os.Remove(dbPath)
	}

	return db, dbPath, cleanup
}

// setupTestBaseData 插入测试基础数据
func setupTestBaseData(t *testing.T, db *sql.DB) {
	// 插入测试角色
	_, err := db.Exec(`
		INSERT OR IGNORE INTO roles (id, name, code, description) VALUES 
		(1, '超级管理员', 'SUPER_ADMIN', '系统超级管理员'),
		(2, '普通用户', 'USER', '普通用户')
	`)
	require.NoError(t, err)

	// 插入测试用户
	_, err = db.Exec(`
		INSERT OR IGNORE INTO users (id, username, email, first_name, last_name, password) VALUES 
		(1, 'admin', 'admin@test.com', 'Admin', 'User', 'hashed_password'),
		(2, 'testuser', 'user@test.com', 'Test', 'User', 'hashed_password')
	`)
	require.NoError(t, err)

	// 建立用户角色关系
	_, err = db.Exec(`
		INSERT OR IGNORE INTO user_roles (user_id, role_id) VALUES 
		(1, 1),
		(2, 2)
	`)
	require.NoError(t, err)
}

// testDocTypeManagement 测试DocType管理
func testDocTypeManagement(t *testing.T, ctx context.Context, svc *service.PermissionService) {
	// 创建DocType
	createReq := &service.CreateDocTypeRequest{
		Name:        "TestDocument",
		Label:       "测试文档",
		Module:      "TestModule",
		Description: "这是一个测试文档类型",
	}

	docType, err := svc.CreateDocType(ctx, createReq)
	require.NoError(t, err)
	assert.Equal(t, "TestDocument", docType.Name)
	assert.Equal(t, "TestModule", docType.Module)

	// 获取DocType
	getDocType, err := svc.GetDocType(ctx, docType.Name)
	require.NoError(t, err)
	assert.Equal(t, docType.ID, getDocType.ID)
	assert.Equal(t, docType.Name, getDocType.Name)

	// 更新DocType
	updateReq := &service.UpdateDocTypeRequest{
		Name: docType.Name,
		CreateDocTypeRequest: service.CreateDocTypeRequest{
			Label:       docType.Label,
			Module:      docType.Module,
			Description: "这是一个更新的测试文档类型",
		},
	}

	updatedDocType, err := svc.UpdateDocType(ctx, updateReq)
	require.NoError(t, err)
	assert.Equal(t, "这是一个更新的测试文档类型", updatedDocType.Description)

	// 列出DocTypes
	listReq := &service.ListDocTypesRequest{
		Module: "TestModule",
	}

	docTypes, err := svc.ListDocTypes(ctx, listReq)
	require.NoError(t, err)
	assert.Greater(t, len(docTypes.DocTypes), 0)

	// 删除DocType
	err = svc.DeleteDocType(ctx, docType.Name)
	require.NoError(t, err)

	// 验证删除
	_, err = svc.GetDocType(ctx, docType.Name)
	assert.Error(t, err)
}

// testPermissionRuleManagement 测试权限规则管理
func testPermissionRuleManagement(t *testing.T, ctx context.Context, svc *service.PermissionService, db *sql.DB) {
	// 创建权限规则
	createReq := &service.CreatePermissionRuleRequest{
		RoleID:          1,
		DocType:         "TestDocument",
		PermissionLevel: 0,
		CanRead:         true,
		CanWrite:        true,
		CanCreate:       true,
	}

	rule, err := svc.CreatePermissionRule(ctx, createReq)
	require.NoError(t, err)
	assert.Equal(t, int64(1), rule.RoleID)
	assert.Equal(t, "TestDocument", rule.DocType)
	assert.True(t, rule.CanRead)
	assert.True(t, rule.CanWrite)
	assert.True(t, rule.CanCreate)

	// 获取权限规则
	getRule, err := svc.GetPermissionRule(ctx, rule.ID)
	require.NoError(t, err)
	assert.Equal(t, rule.ID, getRule.ID)
	assert.Equal(t, rule.RoleID, getRule.RoleID)

	// 更新权限规则
	updateReq := &service.UpdatePermissionRuleRequest{
		ID: rule.ID,
		CreatePermissionRuleRequest: service.CreatePermissionRuleRequest{
			RoleID:          rule.RoleID,
			DocType:         rule.DocType,
			PermissionLevel: rule.PermissionLevel,
			CanRead:         true,
			CanWrite:        true,
			CanCreate:       true,
			CanDelete:       true, // 新增删除权限
		},
	}

	updatedRule, err := svc.UpdatePermissionRule(ctx, updateReq)
	require.NoError(t, err)
	assert.True(t, updatedRule.CanDelete)

	// 列出权限规则
	listReq := &service.ListPermissionRulesRequest{
		RoleID:  1,
		DocType: "TestDocument",
	}

	rules, err := svc.ListPermissionRules(ctx, listReq)
	require.NoError(t, err)
	assert.Greater(t, len(rules.Rules), 0)

	// 删除权限规则
	err = svc.DeletePermissionRule(ctx, rule.ID)
	require.NoError(t, err)

	// 验证删除
	_, err = svc.GetPermissionRule(ctx, rule.ID)
	assert.Error(t, err)
}

// testUserPermissionManagement 测试用户权限管理
func testUserPermissionManagement(t *testing.T, ctx context.Context, svc *service.PermissionService, db *sql.DB) {
	// 创建用户权限
	createReq := &service.CreateUserPermissionRequest{
		UserID:          2,
		DocType:         "TestDocument",
		DocumentName:    "testdoc",
		Condition:       "test-condition",
		ApplicableFor:   stringPtr("TestDocument"),
		HideDescendants: false,
		IsDefault:       true,
	}

	userPerm, err := svc.CreateUserPermission(ctx, createReq)
	require.NoError(t, err)
	assert.Equal(t, int64(2), userPerm.UserID)
	assert.Equal(t, "TestDocument", userPerm.DocType)
	assert.Equal(t, "testdoc", userPerm.DocumentName)

	// 获取用户权限
	getUserPerm, err := svc.GetUserPermission(ctx, userPerm.ID)
	require.NoError(t, err)
	assert.Equal(t, userPerm.ID, getUserPerm.ID)

	// 更新用户权限
	updateReq := &service.CreateUserPermissionRequest{
		UserID:          userPerm.UserID,
		DocType:         userPerm.DocType,
		DocumentName:    userPerm.DocumentName,
		Condition:       userPerm.Condition,
		ApplicableFor:   stringPtr("TestDocument"),
		HideDescendants: true, // 更新为隐藏子项
		IsDefault:       true,
	}

	updatedUserPerm, err := svc.UpdateUserPermission(ctx, userPerm.ID, updateReq)
	require.NoError(t, err)
	assert.NotNil(t, updatedUserPerm)

	// 列出用户权限
	listReq := &service.ListUserPermissionsRequest{
		UserID:  2,
		DocType: "TestDocument",
		Page:    1,
		Size:    10,
	}

	userPerms, err := svc.ListUserPermissions(ctx, listReq)
	require.NoError(t, err)
	assert.Greater(t, len(userPerms.UserPermissions), 0)

	// 删除用户权限
	err = svc.DeleteUserPermission(ctx, userPerm.ID)
	require.NoError(t, err)

	// 验证删除
	_, err = svc.GetUserPermission(ctx, userPerm.ID)
	assert.Error(t, err)
}

// testPermissionChecking 测试权限检查功能
func testPermissionChecking(t *testing.T, ctx context.Context, svc *service.PermissionService, db *sql.DB) {
	// 首先创建一个权限规则供测试
	_, err := db.Exec(`
		INSERT OR REPLACE INTO permission_rules 
		(role, document_type, permission_level, read, write, [create]) 
		VALUES (1, 'TestDocument', 0, true, true, true)
	`)
	require.NoError(t, err)

	// 测试文档权限检查
	checkReq := &service.CheckDocumentPermissionRequest{
		UserID:     1,
		DocType:    "TestDocument",
		Permission: "read",
	}

	result, err := svc.CheckDocumentPermission(ctx, checkReq)
	require.NoError(t, err)
	assert.True(t, result.HasPermission)

	// 测试无权限的情况
	checkReq2 := &service.CheckDocumentPermissionRequest{
		UserID:     2, // 普通用户没有TestDocument DocType的权限
		DocType:    "TestDocument",
		Permission: "write",
	}

	result2, err := svc.CheckDocumentPermission(ctx, checkReq2)
	require.NoError(t, err)
	assert.False(t, result2.HasPermission)

	// 测试文档过滤
	filterReq := &service.FilterDocumentsByPermissionRequest{
		UserID:       1,
		DocumentType: "TestDocument",
		Documents:    []map[string]interface{}{{"id": 1, "name": "test"}},
	}

	filteredDocs, err := svc.FilterDocumentsByPermission(ctx, filterReq)
	require.NoError(t, err)
	assert.NotNil(t, filteredDocs)
}

// stringPtr 返回字符串指针
func stringPtr(s string) *string {
	return &s
}

// TestPermissionSystemLoadTest 负载测试
func TestPermissionSystemLoadTest(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过负载测试")
	}

	// 设置测试数据库
	db, dbPath, cleanup := setupTestDatabase(t)
	defer cleanup()

	// 插入测试基础数据
	setupTestBaseData(t, db)

	// 创建依赖
	logger := log.DefaultLogger
	permissionCache := cache.NewMemoryPermissionCache(logger)
	dataInstance, cleanup2, err := data.NewData(&conf.Data{
		Database: &conf.Database{Driver: "sqlite3", Source: dbPath},
		Redis:    &conf.Redis{Addr: "localhost:6379", Password: "", DB: 0},
	}, logger)
	require.NoError(t, err)
	defer cleanup2()

	permissionRepo := data.NewPermissionRepo(dataInstance, logger)
	cachedPermissionRepo := data.NewCachedPermissionRepo(permissionRepo, permissionCache, logger)

	permissionUsecase := biz.NewPermissionUsecase(cachedPermissionRepo, logger)
	permissionService := service.NewPermissionService(permissionUsecase, logger)

	// 创建带有超级管理员权限的测试上下文
	ctx := context.Background()
	ctx = middleware.SetUserIDToContext(ctx, 1)
	ctx = middleware.SetUsernameToContext(ctx, "testadmin")
	ctx = middleware.SetUserEmailToContext(ctx, "admin@test.com")
	ctx = middleware.SetUserRolesToContext(ctx, []string{"SUPER_ADMIN"})

	// 并发权限检查测试
	t.Run("Concurrent Permission Checks", func(t *testing.T) {
		// 先创建权限规则
		_, err := db.Exec(`
			INSERT OR REPLACE INTO permission_rules 
			(role, document_type, permission_level, read) 
			VALUES (1, 'User', 0, 1)
		`)
		require.NoError(t, err)

		// 并发执行权限检查
		const numConcurrent = 50
		results := make(chan bool, numConcurrent)
		errors := make(chan error, numConcurrent)

		for i := 0; i < numConcurrent; i++ {
			go func() {
				checkReq := &service.CheckDocumentPermissionRequest{
					UserID:     1,
					DocType:    "User",
					Permission: "read",
				}

				result, err := permissionService.CheckDocumentPermission(ctx, checkReq)
				if err != nil {
					errors <- err
				} else {
					results <- result.HasPermission
				}
			}()
		}

		// 收集结果
		successCount := 0
		errorCount := 0

		for i := 0; i < numConcurrent; i++ {
			select {
			case hasPermission := <-results:
				if hasPermission {
					successCount++
				}
			case err := <-errors:
				t.Logf("权限检查错误: %v", err)
				errorCount++
			case <-time.After(10 * time.Second):
				t.Fatal("权限检查超时")
			}
		}

		t.Logf("并发权限检查结果: 成功 %d, 错误 %d", successCount, errorCount)
		assert.Greater(t, successCount, numConcurrent/2) // 至少一半成功
	})
}

// BenchmarkPermissionCheck 性能基准测试
func BenchmarkPermissionCheck(b *testing.B) {
	// 设置测试数据库
	db, dbPath, cleanup := setupBenchmarkDatabase(b)
	defer cleanup()

	// 创建依赖
	logger := log.DefaultLogger
	permissionCache := cache.NewMemoryPermissionCache(logger)
	dataInstance, cleanup2, err := data.NewData(&conf.Data{
		Database: &conf.Database{Driver: "sqlite3", Source: dbPath},
		Redis:    &conf.Redis{Addr: "localhost:6379", Password: "", DB: 0},
	}, logger)
	if err != nil {
		b.Fatal(err)
	}
	defer cleanup2()

	permissionRepo := data.NewPermissionRepo(dataInstance, logger)
	cachedPermissionRepo := data.NewCachedPermissionRepo(permissionRepo, permissionCache, logger)

	permissionUsecase := biz.NewPermissionUsecase(cachedPermissionRepo, logger)
	permissionService := service.NewPermissionService(permissionUsecase, logger)

	ctx := context.Background()

	// 创建权限规则用于基准测试
	_, err = db.Exec(`
		INSERT OR REPLACE INTO permission_rules 
		(role, document_type, permission_level, read) 
		VALUES (1, 'User', 0, true)
	`)
	if err != nil {
		b.Fatal(err)
	}

	checkReq := &service.CheckDocumentPermissionRequest{
		UserID:     1,
		DocType:    "User",
		Permission: "read",
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, err := permissionService.CheckDocumentPermission(ctx, checkReq)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}

// setupBenchmarkDatabase 设置基准测试数据库
func setupBenchmarkDatabase(b *testing.B) (*sql.DB, string, func()) {
	dbPath := fmt.Sprintf("/tmp/benchmark_permission_%d.db", time.Now().Unix())

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		b.Fatal(err)
	}

	// 获取项目根路径
	projectRoot := getProjectRoot()

	// 尝试使用完整的schema文件，失败则使用基础表创建
	coreSchemaSQL, err := os.ReadFile(filepath.Join(projectRoot, "database", "schema", "core_system.sql"))
	if err == nil {
		// 执行完整的core schema
		_, err = db.Exec(string(coreSchemaSQL))
		if err != nil {
			b.Fatal(err)
		}

		// 执行权限系统schema
		permissionSchemaSQL, err := os.ReadFile(filepath.Join(projectRoot, "database", "schema", "permission_system.sql"))
		if err == nil {
			_, err = db.Exec(string(permissionSchemaSQL))
			if err != nil {
				b.Fatal(err)
			}
		}

		// 尝试执行种子数据（如果失败则手动插入基础数据）
		seedDataInserted := false
		coreSeedSQL, err := os.ReadFile("../../../database/data/core_seed.sql")
		if err == nil {
			if db.Exec(string(coreSeedSQL)); err == nil {
				permissionSeedSQL, err := os.ReadFile("../../../database/data/permission_seed.sql")
				if err == nil {
					if db.Exec(string(permissionSeedSQL)); err == nil {
						seedDataInserted = true
					}
				}
			}
		}

		// 如果种子数据插入失败，手动插入基础数据
		if !seedDataInserted {
			_, err = db.Exec(`
				INSERT OR IGNORE INTO roles (id, name, code, description) VALUES (1, '管理员', 'ADMIN', '系统管理员');
				INSERT OR IGNORE INTO users (id, username, email, first_name, last_name, password) 
				VALUES (1, 'admin', 'admin@test.com', 'Admin', 'User', 'password');
				INSERT OR IGNORE INTO user_roles (user_id, role_id) VALUES (1, 1);
			`)
			if err != nil {
				b.Fatal(err)
			}
		}
	} else {
		// 回退到基础表创建逻辑
		createBasicTables := `
		CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			username VARCHAR(50) UNIQUE NOT NULL,
			email VARCHAR(100) UNIQUE NOT NULL,
			password VARCHAR(255) NOT NULL,
			first_name VARCHAR(50) NOT NULL,
			last_name VARCHAR(50) NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		);

		CREATE TABLE IF NOT EXISTS roles (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name VARCHAR(50) UNIQUE NOT NULL,
			code VARCHAR(50) UNIQUE NOT NULL,
			description TEXT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		);

		CREATE TABLE IF NOT EXISTS user_roles (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER NOT NULL,
			role_id INTEGER NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			UNIQUE(user_id, role_id)
		);

		CREATE TABLE IF NOT EXISTS permission_rules (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			role INTEGER NOT NULL,
			document_type VARCHAR(140) NOT NULL,
			permission_level INTEGER DEFAULT 0,
			read BOOLEAN DEFAULT false,
			write BOOLEAN DEFAULT false,
			[create] BOOLEAN DEFAULT false,
			[delete] BOOLEAN DEFAULT false,
			submit BOOLEAN DEFAULT false,
			cancel BOOLEAN DEFAULT false,
			amend BOOLEAN DEFAULT false,
			print BOOLEAN DEFAULT false,
			email BOOLEAN DEFAULT false,
			import BOOLEAN DEFAULT false,
			export BOOLEAN DEFAULT false,
			share BOOLEAN DEFAULT false,
			report BOOLEAN DEFAULT false,
			set_user_permissions BOOLEAN DEFAULT false,
			if_owner BOOLEAN DEFAULT false,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			UNIQUE(role, document_type, permission_level)
		);

		INSERT OR IGNORE INTO roles (id, name, code, description) VALUES (1, '管理员', 'ADMIN', '系统管理员');
		INSERT OR IGNORE INTO users (id, username, email, first_name, last_name, password) 
		VALUES (1, 'admin', 'admin@test.com', 'Admin', 'User', 'password');
		INSERT OR IGNORE INTO user_roles (user_id, role_id) VALUES (1, 1);
		`

		_, err = db.Exec(createBasicTables)
		if err != nil {
			b.Fatal(err)
		}
	}

	cleanup := func() {
		db.Close()
		os.Remove(dbPath)
	}

	return db, dbPath, cleanup
}
