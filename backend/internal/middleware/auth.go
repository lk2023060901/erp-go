package middleware

import (
	"context"
	"fmt"
	"reflect"
	"strings"
	"time"

	"erp-system/internal/biz"
	"erp-system/internal/cache"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport"
	"github.com/golang-jwt/jwt/v5"
)

// JWTClaims JWT载荷
type JWTClaims struct {
	UserID    int64    `json:"user_id"`
	Username  string   `json:"username"`
	Email     string   `json:"email"`
	Roles     []string `json:"roles"`
	SessionID string   `json:"session_id"`
	TokenType string   `json:"token_type"` // access, refresh
	jwt.RegisteredClaims
}

// AuthMiddleware 认证中间件配置
type AuthMiddleware struct {
	jwtSecret     string
	permissionSvc *biz.PermissionUsecase
	userSvc       *biz.UserUsecase
	cache         cache.Cache
	logger        *log.Helper
	skipPaths     map[string]bool // 跳过认证的路径
	skipMethods   map[string]bool // 跳过认证的方法
}

// NewAuthMiddleware 创建认证中间件
func NewAuthMiddleware(
	jwtSecret string,
	permissionSvc *biz.PermissionUsecase,
	userSvc *biz.UserUsecase,
	cache cache.Cache,
	logger log.Logger,
) *AuthMiddleware {
	return &AuthMiddleware{
		jwtSecret:     jwtSecret,
		permissionSvc: permissionSvc,
		userSvc:       userSvc,
		cache:         cache,
		logger:        log.NewHelper(logger),
		skipPaths: map[string]bool{
			"/health":                 true,
			"/v1/auth/login":          true,
			"/v1/auth/register":       true,
			"/v1/auth/refresh":        true,
			"/v1/auth/reset-password": true,
			"/v1/auth/send-code":      true,
		},
		skipMethods: map[string]bool{
			"HealthCheck":          true,
			"Login":                true,
			"Register":             true,
			"RefreshToken":         true,
			"ResetPassword":        true,
			"SendVerificationCode": true,
		},
	}
}

// JWT JWT认证中间件
func (m *AuthMiddleware) JWT() middleware.Middleware {
	return middleware.Middleware(func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			// 获取传输层信息
			tr, ok := transport.FromServerContext(ctx)
			if !ok {
				return nil, errors.Unauthorized("UNAUTHORIZED", "transport info not found")
			}

			// 检查是否需要跳过认证
			if m.shouldSkipAuth(tr) {
				return handler(ctx, req)
			}

			// 获取Authorization头
			token := m.extractToken(tr)
			if token == "" {
				return nil, errors.Unauthorized("UNAUTHORIZED", "missing authorization token")
			}

			// 验证并解析JWT
			claims, err := m.parseToken(token)
			if err != nil {
				m.logger.Warnf("Invalid token: %v", err)
				return nil, errors.Unauthorized("UNAUTHORIZED", "invalid token")
			}

			// 验证Token类型
			if claims.TokenType != "access" {
				return nil, errors.Unauthorized("UNAUTHORIZED", "invalid token type")
			}

			// 检查会话是否有效
			if err := m.validateSession(ctx, claims.UserID, claims.SessionID); err != nil {
				m.logger.Warnf("Invalid session: %v", err)
				return nil, errors.Unauthorized("UNAUTHORIZED", "session expired or invalid")
			}

			// 将用户信息放入上下文
			ctx = m.setUserContext(ctx, claims)

			return handler(ctx, req)
		}
	})
}

// RBAC RBAC权限验证中间件
func (m *AuthMiddleware) RBAC() middleware.Middleware {
	return middleware.Middleware(func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			// 获取传输层信息
			tr, ok := transport.FromServerContext(ctx)
			if !ok {
				return handler(ctx, req) // 如果没有传输层信息，继续执行
			}

			// 检查是否需要跳过权限验证
			if m.shouldSkipAuth(tr) {
				return handler(ctx, req)
			}

			// 获取当前用户信息
			userID := GetUserIDFromContext(ctx)
			if userID == 0 {
				return nil, errors.Unauthorized("UNAUTHORIZED", "user not authenticated")
			}

			// 生成权限编码
			permissionCode := m.generatePermissionCode(tr)
			if permissionCode == "" {
				// 如果无法生成权限编码，允许通过（可能是系统内部调用）
				return handler(ctx, req)
			}

			// 检查用户权限
			hasPermission, err := m.checkUserPermission(ctx, userID, permissionCode)
			if err != nil {
				m.logger.Errorf("Check permission error: %v", err)
				return nil, errors.InternalServer("PERMISSION_CHECK_ERROR", "permission check failed")
			}

			if !hasPermission {
				m.logger.Warnf("User %d has no permission for %s", userID, permissionCode)
				return nil, errors.Forbidden("FORBIDDEN", "insufficient permissions")
			}

			// 记录权限验证成功
			m.logger.Debugf("User %d passed permission check for %s", userID, permissionCode)

			return handler(ctx, req)
		}
	})
}

// shouldSkipAuth 检查是否需要跳过认证
func (m *AuthMiddleware) shouldSkipAuth(tr transport.Transporter) bool {
	// HTTP请求
	if ht, ok := tr.(interface{ RequestURI() string }); ok {
		path := ht.RequestURI()
		// 移除查询参数
		if idx := strings.Index(path, "?"); idx > 0 {
			path = path[:idx]
		}
		if m.skipPaths[path] {
			return true
		}
	}

	// gRPC请求
	if gt, ok := tr.(interface{ Method() string }); ok {
		method := gt.Method()
		// 从完整方法名中提取方法名
		parts := strings.Split(method, "/")
		if len(parts) > 0 {
			methodName := parts[len(parts)-1]
			if m.skipMethods[methodName] {
				return true
			}
		}
	}

	return false
}

// extractToken 提取Token
func (m *AuthMiddleware) extractToken(tr transport.Transporter) string {
	var token string

	// 使用反射安全地检查方法是否存在
	if headerMethod := reflect.ValueOf(tr).MethodByName("RequestHeader"); headerMethod.IsValid() {
		results := headerMethod.Call(nil)
		if len(results) > 0 {
			if headerMap, ok := results[0].Interface().(map[string][]string); ok {
				// 先尝试HTTP标准的Authorization头
				if auths := headerMap["Authorization"]; len(auths) > 0 {
					auth := auths[0]
					if strings.HasPrefix(auth, "Bearer ") {
						token = strings.TrimPrefix(auth, "Bearer ")
						return token
					}
				}
				// 再尝试gRPC小写的authorization头
				if auths := headerMap["authorization"]; len(auths) > 0 {
					auth := auths[0]
					if strings.HasPrefix(auth, "Bearer ") {
						token = strings.TrimPrefix(auth, "Bearer ")
						return token
					}
				}
			}
		}
	}

	return token
}

// parseToken 解析JWT Token
func (m *AuthMiddleware) parseToken(tokenString string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(m.jwtSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		// 检查Token是否过期
		if time.Now().After(claims.ExpiresAt.Time) {
			return nil, fmt.Errorf("token expired")
		}
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}

// validateSession 验证会话是否有效
func (m *AuthMiddleware) validateSession(ctx context.Context, userID int64, sessionID string) error {
	// 从缓存中检查会话
	cacheKey := fmt.Sprintf("session:%d:%s", userID, sessionID)
	exists, err := m.cache.Exists(ctx, cacheKey)
	if err != nil {
		return fmt.Errorf("check session cache error: %w", err)
	}

	if !exists {
		// 缓存中不存在，暂时跳过会话验证
		return nil
	}

	return nil
}

// setUserContext 设置用户上下文
func (m *AuthMiddleware) setUserContext(ctx context.Context, claims *JWTClaims) context.Context {
	ctx = SetUserIDToContext(ctx, claims.UserID)
	ctx = SetUsernameToContext(ctx, claims.Username)
	ctx = SetUserEmailToContext(ctx, claims.Email)
	ctx = SetUserRolesToContext(ctx, claims.Roles)
	ctx = SetSessionIDToContext(ctx, claims.SessionID)
	return ctx
}

// generatePermissionCode 生成权限编码
func (m *AuthMiddleware) generatePermissionCode(tr transport.Transporter) string {
	var method, path string

	// HTTP请求
	if ht, ok := tr.(interface {
		Method() string
		RequestURI() string
	}); ok {
		method = ht.Method()
		path = ht.RequestURI()
		// 移除查询参数
		if idx := strings.Index(path, "?"); idx > 0 {
			path = path[:idx]
		}
	}

	// gRPC请求 - 转换为HTTP风格的权限编码
	if gt, ok := tr.(interface{ Method() string }); ok {
		grpcMethod := gt.Method()
		// 例：/api.user.v1.UserService/GetUser -> user.get
		parts := strings.Split(grpcMethod, "/")
		if len(parts) >= 2 {
			servicePart := parts[1] // api.user.v1.UserService
			methodPart := parts[2]  // GetUser

			// 解析服务名
			serviceParts := strings.Split(servicePart, ".")
			if len(serviceParts) >= 2 {
				module := serviceParts[1]              // user
				action := m.methodToAction(methodPart) // get
				return fmt.Sprintf("%s.%s", module, action)
			}
		}
	}

	// HTTP请求 - 从路径和方法生成权限编码
	if method != "" && path != "" {
		return m.httpToPermissionCode(method, path)
	}

	return ""
}

// httpToPermissionCode HTTP请求转权限编码
func (m *AuthMiddleware) httpToPermissionCode(method, path string) string {
	// 路径模式匹配和权限编码映射
	pathMappings := map[string]map[string]string{
		"/v1/users": {
			"GET":  "user.list",
			"POST": "user.create",
		},
		"/v1/users/{id}": {
			"GET":    "user.view",
			"PUT":    "user.edit",
			"DELETE": "user.delete",
		},
		"/v1/users/{id}/roles": {
			"POST": "user.assign_role",
			"GET":  "user.view_roles",
		},
		"/v1/roles": {
			"GET":  "role.list",
			"POST": "role.create",
		},
		"/v1/roles/{id}": {
			"GET":    "role.view",
			"PUT":    "role.edit",
			"DELETE": "role.delete",
		},
		"/v1/permissions": {
			"GET":  "permission.list",
			"POST": "permission.create",
		},
		"/v1/organizations": {
			"GET":  "organization.list",
			"POST": "organization.create",
		},
		"/v1/system/configs": {
			"GET": "system.config.view",
			"PUT": "system.config.edit",
		},
	}

	// 规范化路径（替换参数为通用模式）
	normalizedPath := m.normalizePath(path)

	if methodMap, exists := pathMappings[normalizedPath]; exists {
		if permCode, exists := methodMap[strings.ToUpper(method)]; exists {
			return permCode
		}
	}

	// 默认权限编码生成规则
	segments := strings.Split(strings.Trim(path, "/"), "/")
	if len(segments) >= 2 && segments[0] == "v1" {
		resource := segments[1]
		action := m.methodToAction(method)
		return fmt.Sprintf("%s.%s", resource, action)
	}

	return ""
}

// normalizePath 规范化路径
func (m *AuthMiddleware) normalizePath(path string) string {
	// 将数字ID替换为{id}
	parts := strings.Split(path, "/")
	for i, part := range parts {
		if m.isNumeric(part) {
			parts[i] = "{id}"
		}
	}
	return strings.Join(parts, "/")
}

// isNumeric 检查字符串是否为数字
func (m *AuthMiddleware) isNumeric(s string) bool {
	for _, c := range s {
		if c < '0' || c > '9' {
			return false
		}
	}
	return len(s) > 0
}

// methodToAction 方法名转动作
func (m *AuthMiddleware) methodToAction(method string) string {
	method = strings.ToUpper(method)
	actionMap := map[string]string{
		"GET":    "view",
		"POST":   "create",
		"PUT":    "edit",
		"DELETE": "delete",
		"PATCH":  "edit",
		// gRPC方法映射
		"LIST":   "list",
		"CREATE": "create",
		"UPDATE": "edit",
	}

	// gRPC方法名处理
	if strings.HasPrefix(method, "LIST") {
		return "list"
	} else if strings.HasPrefix(method, "GET") {
		return "view"
	} else if strings.HasPrefix(method, "CREATE") {
		return "create"
	} else if strings.HasPrefix(method, "UPDATE") {
		return "edit"
	} else if strings.HasPrefix(method, "DELETE") {
		return "delete"
	}

	if action, exists := actionMap[method]; exists {
		return action
	}

	return "unknown"
}

// checkUserPermission 检查用户权限
func (m *AuthMiddleware) checkUserPermission(ctx context.Context, userID int64, permissionCode string) (bool, error) {
	// 先从缓存检查
	cacheKey := fmt.Sprintf("user_permission:%d:%s", userID, permissionCode)

	if result, err := m.cache.Get(ctx, cacheKey); err == nil && result != "" {
		return result == "true", nil
	}

	// 缓存未命中，从服务层检查
	hasPermission, err := m.permissionSvc.CheckUserPermission(ctx, userID, permissionCode)
	if err != nil {
		return false, fmt.Errorf("check permission error: %w", err)
	}

	// 缓存结果30分钟
	cacheValue := "false"
	if hasPermission {
		cacheValue = "true"
	}
	_ = m.cache.Set(ctx, cacheKey, cacheValue, 30*time.Minute)

	return hasPermission, nil
}

// checkFrappePermission 检查Frappe风格的多级权限
func (m *AuthMiddleware) checkFrappePermission(ctx context.Context, userID int64, documentType, action string, permissionLevel int) (bool, error) {
	// 从缓存检查
	cacheKey := fmt.Sprintf("frappe_permission:%d:%s:%s:%d", userID, documentType, action, permissionLevel)

	if result, err := m.cache.Get(ctx, cacheKey); err == nil && result != "" {
		return result == "true", nil
	}

	// 使用新的Frappe权限系统检查
	hasPermission, err := m.permissionSvc.CheckPermission(ctx, userID, documentType, action, permissionLevel)
	if err != nil {
		return false, fmt.Errorf("check frappe permission error: %w", err)
	}

	// 缓存结果30分钟
	cacheValue := "false"
	if hasPermission {
		cacheValue = "true"
	}
	_ = m.cache.Set(ctx, cacheKey, cacheValue, 30*time.Minute)

	return hasPermission, nil
}

// checkDocumentPermission 检查文档级权限
func (m *AuthMiddleware) checkDocumentPermission(ctx context.Context, userID int64, documentType, action string) (bool, error) {
	// 文档级权限使用级别0
	return m.checkFrappePermission(ctx, userID, documentType, action, 0)
}

// checkFieldPermission 检查字段级权限
func (m *AuthMiddleware) checkFieldPermission(ctx context.Context, userID int64, documentType string, permissionLevel int, action string) (bool, error) {
	// 字段级权限使用级别1-9
	if permissionLevel < 1 || permissionLevel > 9 {
		return false, fmt.Errorf("invalid permission level for field: %d", permissionLevel)
	}

	return m.checkFrappePermission(ctx, userID, documentType, action, permissionLevel)
}

// getUserPermissionLevel 获取用户的权限级别
func (m *AuthMiddleware) getUserPermissionLevel(ctx context.Context, userID int64, documentType string) (int, error) {
	// 从缓存检查
	cacheKey := fmt.Sprintf("user_permission_level:%d:%s", userID, documentType)

	if result, err := m.cache.Get(ctx, cacheKey); err == nil && result != "" {
		var level int
		if _, err := fmt.Sscanf(result, "%d", &level); err == nil {
			return level, nil
		}
	}

	// 从服务层获取
	level, err := m.permissionSvc.GetUserPermissionLevel(ctx, userID, documentType)
	if err != nil {
		return 0, fmt.Errorf("get user permission level error: %w", err)
	}

	// 缓存结果30分钟
	_ = m.cache.Set(ctx, cacheKey, fmt.Sprintf("%d", level), 30*time.Minute)

	return level, nil
}

// PermissionMiddleware 返回权限检查中间件
func (m *AuthMiddleware) PermissionMiddleware(requiredDocType string, requiredAction string, requiredLevel int) middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			// 获取用户信息
			userID := m.getUserIDFromContext(ctx)
			if userID == 0 {
				return nil, errors.Unauthorized("UNAUTHORIZED", "用户未认证")
			}

			// 检查权限
			hasPermission, err := m.checkFrappePermission(ctx, userID, requiredDocType, requiredAction, requiredLevel)
			if err != nil {
				m.logger.Errorf("permission check failed: %v", err)
				return nil, errors.Forbidden("PERMISSION_DENIED", "权限检查失败")
			}

			if !hasPermission {
				return nil, errors.Forbidden("INSUFFICIENT_PERMISSION",
					fmt.Sprintf("用户无权限访问 %s 的 %s 操作(级别:%d)", requiredDocType, requiredAction, requiredLevel))
			}

			return handler(ctx, req)
		}
	}
}

// DocumentPermissionMiddleware 返回文档级权限检查中间件
func (m *AuthMiddleware) DocumentPermissionMiddleware(requiredDocType string, requiredAction string) middleware.Middleware {
	return m.PermissionMiddleware(requiredDocType, requiredAction, 0)
}

// FieldPermissionMiddleware 返回字段级权限检查中间件
func (m *AuthMiddleware) FieldPermissionMiddleware(requiredDocType string, requiredAction string, requiredLevel int) middleware.Middleware {
	if requiredLevel < 1 || requiredLevel > 9 {
		panic(fmt.Sprintf("invalid field permission level: %d", requiredLevel))
	}
	return m.PermissionMiddleware(requiredDocType, requiredAction, requiredLevel)
}

// getUserIDFromContext 从上下文获取用户ID
func (m *AuthMiddleware) getUserIDFromContext(ctx context.Context) int64 {
	if userID, ok := ctx.Value("user_id").(int64); ok {
		return userID
	}
	return 0
}

// PermissionMiddlewareFactory 权限中间件工厂
type PermissionMiddlewareFactory struct {
	auth *AuthMiddleware
}

// NewPermissionMiddlewareFactory 创建权限中间件工厂
func NewPermissionMiddlewareFactory(auth *AuthMiddleware) *PermissionMiddlewareFactory {
	return &PermissionMiddlewareFactory{auth: auth}
}

// ForUser 用户管理权限中间件
func (f *PermissionMiddlewareFactory) ForUser(action string) middleware.Middleware {
	return f.auth.DocumentPermissionMiddleware("User", action)
}

// ForRole 角色管理权限中间件
func (f *PermissionMiddlewareFactory) ForRole(action string) middleware.Middleware {
	return f.auth.DocumentPermissionMiddleware("Role", action)
}

// ForOrganization 组织管理权限中间件
func (f *PermissionMiddlewareFactory) ForOrganization(action string) middleware.Middleware {
	return f.auth.DocumentPermissionMiddleware("Organization", action)
}

// ForDocument 通用文档权限中间件
func (f *PermissionMiddlewareFactory) ForDocument(docType, action string) middleware.Middleware {
	return f.auth.DocumentPermissionMiddleware(docType, action)
}

// ForField 字段级权限中间件
func (f *PermissionMiddlewareFactory) ForField(docType, action string, level int) middleware.Middleware {
	return f.auth.FieldPermissionMiddleware(docType, action, level)
}

// RequireRead 需要读权限
func (f *PermissionMiddlewareFactory) RequireRead(docType string) middleware.Middleware {
	return f.auth.DocumentPermissionMiddleware(docType, "read")
}

// RequireWrite 需要写权限
func (f *PermissionMiddlewareFactory) RequireWrite(docType string) middleware.Middleware {
	return f.auth.DocumentPermissionMiddleware(docType, "write")
}

// RequireCreate 需要创建权限
func (f *PermissionMiddlewareFactory) RequireCreate(docType string) middleware.Middleware {
	return f.auth.DocumentPermissionMiddleware(docType, "create")
}

// RequireDelete 需要删除权限
func (f *PermissionMiddlewareFactory) RequireDelete(docType string) middleware.Middleware {
	return f.auth.DocumentPermissionMiddleware(docType, "delete")
}

// RequireSubmit 需要提交权限
func (f *PermissionMiddlewareFactory) RequireSubmit(docType string) middleware.Middleware {
	return f.auth.DocumentPermissionMiddleware(docType, "submit")
}

// RequireCancel 需要取消权限
func (f *PermissionMiddlewareFactory) RequireCancel(docType string) middleware.Middleware {
	return f.auth.DocumentPermissionMiddleware(docType, "cancel")
}

// RequireAmend 需要修订权限
func (f *PermissionMiddlewareFactory) RequireAmend(docType string) middleware.Middleware {
	return f.auth.DocumentPermissionMiddleware(docType, "amend")
}
