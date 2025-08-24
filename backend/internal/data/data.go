package data

import (
	"database/sql"
	"erp-system/internal/conf"
	"os"
	"strings"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-redis/redis/v8"
	"github.com/google/wire"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(NewData, NewUserRepo, NewRoleRepo, NewPermissionRepo, NewOrganizationRepo, NewSessionRepo, NewAuditRepo)

// Data .
type Data struct {
	db    *sql.DB
	redis *redis.Client
}

// NewData .
func NewData(c *conf.Data, logger log.Logger) (*Data, func(), error) {
	l := log.NewHelper(logger)

	// 数据库连接
	db, err := sql.Open(c.Database.Driver, c.Database.Source)
	if err != nil {
		l.Errorf("failed to open database: %v", err)
		return nil, nil, err
	}

	if err := db.Ping(); err != nil {
		l.Errorf("failed to ping database: %v", err)
		return nil, nil, err
	}

	// SQLite数据库初始化
	if c.Database.Driver == "sqlite3" {
		if err := initSQLiteDB(db, l); err != nil {
			l.Errorf("failed to initialize SQLite database: %v", err)
			return nil, nil, err
		}
	}

	// Redis连接
	rdb := redis.NewClient(&redis.Options{
		Addr:     c.Redis.Addr,
		Password: c.Redis.Password,
		DB:       int(c.Redis.DB),
	})

	d := &Data{
		db:    db,
		redis: rdb,
	}

	cleanup := func() {
		l.Info("closing the data resources")
		if err := d.db.Close(); err != nil {
			l.Errorf("failed to close database: %v", err)
		}
		if err := d.redis.Close(); err != nil {
			l.Errorf("failed to close redis: %v", err)
		}
	}

	return d, cleanup, nil
}

// initSQLiteDB 初始化SQLite数据库
func initSQLiteDB(db *sql.DB, log *log.Helper) error {
	// 执行核心系统表结构
	coreSchemaFile := "../database/schema/core_system.sql"
	coreSchemaBytes, err := os.ReadFile(coreSchemaFile)
	if err != nil {
		log.Warnf("core_system.sql not found, creating basic tables manually")
		return initBasicTables(db, log)
	}

	// 执行核心系统SQL脚本
	if err := executeSQLScript(db, string(coreSchemaBytes), log); err != nil {
		return err
	}

	// 执行权限系统表结构
	permissionSchemaFile := "../database/schema/permission_system.sql"
	permissionSchemaBytes, err := os.ReadFile(permissionSchemaFile)
	if err == nil {
		if err := executeSQLScript(db, string(permissionSchemaBytes), log); err != nil {
			log.Warnf("Failed to execute permission system schema: %v", err)
		}
	}

	// 执行核心系统种子数据
	coreSeedFile := "../database/data/core_seed.sql"
	coreSeedBytes, err := os.ReadFile(coreSeedFile)
	if err == nil {
		if err := executeSQLScript(db, string(coreSeedBytes), log); err != nil {
			log.Warnf("Failed to execute core seed data: %v", err)
		}
	}

	// 执行权限系统种子数据
	permissionSeedFile := "../database/data/permission_seed.sql"
	permissionSeedBytes, err := os.ReadFile(permissionSeedFile)
	if err == nil {
		if err := executeSQLScript(db, string(permissionSeedBytes), log); err != nil {
			log.Warnf("Failed to execute permission seed data: %v", err)
		}
	}

	log.Info("SQLite database initialized successfully")
	return nil
}

// executeSQLScript 执行SQL脚本
func executeSQLScript(db *sql.DB, sqlScript string, log *log.Helper) error {
	statements := strings.Split(sqlScript, ";")

	for _, stmt := range statements {
		stmt = strings.TrimSpace(stmt)
		if stmt == "" || strings.HasPrefix(stmt, "--") {
			continue
		}

		if _, err := db.Exec(stmt); err != nil {
			log.Errorf("failed to execute SQL statement: %s, error: %v", stmt, err)
			return err
		}
	}

	log.Info("SQLite database initialized successfully")
	return nil
}

// initBasicTables 创建基础表（如果找不到init_db.sql文件）
func initBasicTables(db *sql.DB, log *log.Helper) error {
	// 基础用户表
	createUsersTable := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username VARCHAR(50) UNIQUE NOT NULL,
		email VARCHAR(100) UNIQUE NOT NULL,
		password VARCHAR(255) NOT NULL,
		first_name VARCHAR(50) NOT NULL,
		last_name VARCHAR(50) NOT NULL,
		phone VARCHAR(20),
		gender VARCHAR(1),
		birth_date DATE,
		avatar_url VARCHAR(500),
		is_active BOOLEAN DEFAULT true,
		two_factor_enabled BOOLEAN DEFAULT false,
		two_factor_secret VARCHAR(255),
		last_login_at DATETIME,
		last_login_ip VARCHAR(45),
		login_count INTEGER DEFAULT 0,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	)`

	if _, err := db.Exec(createUsersTable); err != nil {
		return err
	}

	log.Info("Basic SQLite tables created successfully")
	return nil
}
