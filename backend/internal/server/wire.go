//go:build wireinject
// +build wireinject

package server

import (
	"erp-system/internal/biz"
	"erp-system/internal/conf"
	"erp-system/internal/data"
	"erp-system/internal/pkg"
	"erp-system/internal/service"
	"time"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
)

// ProviderSet 是所有提供者的集合
var ProviderSet = wire.NewSet(
	// Data layer
	data.ProviderSet,
	
	// Business logic layer
	biz.NewUserUsecase,
	biz.NewRoleUsecase, 
	biz.NewPermissionUsecase,
	biz.NewOrganizationUsecase,
	biz.NewAuditUsecase,
	
	// Service layer
	service.NewAuthService,
	service.NewUserService,
	service.NewRoleService,
	// service.NewPermissionService,  // Temporarily disabled
	service.NewOrganizationService,
	service.NewSystemService,
	
	// Infrastructure
	pkg.NewPasswordManager,
	NewJWTManager,
	
	// Servers
	NewHTTPServer,
	NewGRPCServer,
)

// NewJWTManager 创建JWT管理器
func NewJWTManager(c *conf.Data) *pkg.JWTManager {
	// 从配置中获取JWT密钥，如果没有配置则使用默认值
	secretKey := "dev-jwt-secret-key-for-testing-only"
	if c.Jwt != nil && c.Jwt.SecretKey != "" {
		secretKey = c.Jwt.SecretKey
	}
	
	// 设置访问令牌过期时间（2小时）
	tokenDuration := time.Hour * 2
	if c.Jwt != nil && c.Jwt.AccessTokenExpire > 0 {
		tokenDuration = time.Duration(c.Jwt.AccessTokenExpire) * time.Second
	}
	
	return pkg.NewJWTManager(secretKey, tokenDuration)
}

// InitializeApp 初始化应用
func InitializeApp(*conf.Server, *conf.Data, log.Logger) (*kratos.App, func(), error) {
	panic(wire.Build(ProviderSet, newApp))
}

// newApp 创建Kratos应用实例
func newApp(logger log.Logger, hs *HTTPServer, gs *GRPCServer) *kratos.App {
	return kratos.New(
		kratos.Name("erp-system"),
		kratos.Version("v1.0.0"),
		kratos.Logger(logger),
		kratos.Server(
			hs.Server,
			gs.Server,
		),
	)
}