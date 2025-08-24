package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"erp-system/internal/conf"
	"erp-system/internal/middleware"
	"erp-system/internal/pkg"
	"erp-system/internal/service"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	khttp "github.com/go-kratos/kratos/v2/transport/http"
	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
)

// HTTPServer HTTP服务器
type HTTPServer struct {
	Server      *khttp.Server
	authService *service.AuthService
	userService *service.UserService
	roleService *service.RoleService
	// permissionService *service.PermissionService  // Temporarily disabled - using Frappe permission system
	organizationService *service.OrganizationService
	systemService       *service.SystemService
	jwtSecret           string
	log                 *log.Helper
}

// NewHTTPServer 创建HTTP服务器
func NewHTTPServer(
	c *conf.Server,
	data *conf.Data,
	authService *service.AuthService,
	userService *service.UserService,
	roleService *service.RoleService,
	// permissionService *service.PermissionService,  // Temporarily disabled
	organizationService *service.OrganizationService,
	systemService *service.SystemService,
	logger log.Logger,
) *HTTPServer {
	var opts = []khttp.ServerOption{}

	if c.Http.Network != "" {
		opts = append(opts, khttp.Network(c.Http.Network))
	}
	if c.Http.Addr != "" {
		opts = append(opts, khttp.Address(c.Http.Addr))
	}
	if c.Http.Timeout != "" {
		timeout := c.Http.GetTimeout()
		opts = append(opts, khttp.Timeout(timeout))
	}

	srv := khttp.NewServer(opts...)

	// 创建自定义的HTTP服务器实例
	jwtSecret := "dev-jwt-secret-key-for-testing-only"
	if data.Jwt != nil && data.Jwt.SecretKey != "" {
		jwtSecret = data.Jwt.SecretKey
	}

	httpSrv := &HTTPServer{
		Server:      srv,
		authService: authService,
		userService: userService,
		roleService: roleService,
		// permissionService:   permissionService,  // Temporarily disabled
		organizationService: organizationService,
		systemService:       systemService,
		jwtSecret:           jwtSecret,
		log:                 log.NewHelper(logger),
	}

	// 注册路由
	httpSrv.registerRoutes()

	return httpSrv
}

// registerRoutes 注册路由
func (s *HTTPServer) registerRoutes() {
	// 创建路由器
	router := mux.NewRouter()

	// API版本前缀
	v1 := router.PathPrefix("/api/v1").Subrouter()

	// 认证相关路由（无需认证）
	auth := v1.PathPrefix("/auth").Subrouter()
	auth.HandleFunc("/login", s.handleLogin).Methods("POST", "OPTIONS")
	auth.HandleFunc("/register", s.handleRegister).Methods("POST", "OPTIONS")
	auth.HandleFunc("/refresh", s.handleRefreshToken).Methods("POST", "OPTIONS")

	// 需要认证的路由
	authenticated := v1.NewRoute().Subrouter()
	// 这里应该使用标准HTTP中间件，而不是Kratos中间件
	authenticated.Use(s.jwtMiddleware)

	// 认证相关路由
	authenticated.HandleFunc("/auth/logout", s.handleLogout).Methods("POST", "OPTIONS")
	authenticated.HandleFunc("/auth/profile", s.handleGetProfile).Methods("GET", "OPTIONS")
	authenticated.HandleFunc("/auth/change-password", s.handleChangePassword).Methods("POST", "OPTIONS")

	// 用户管理路由
	users := authenticated.PathPrefix("/users").Subrouter()
	users.HandleFunc("", s.handleListUsers).Methods("GET", "OPTIONS")
	users.HandleFunc("", s.handleCreateUser).Methods("POST", "OPTIONS")
	users.HandleFunc("/{id:[0-9]+}", s.handleGetUser).Methods("GET", "OPTIONS")
	users.HandleFunc("/{id:[0-9]+}", s.handleUpdateUser).Methods("PUT", "OPTIONS")
	users.HandleFunc("/{id:[0-9]+}", s.handleDeleteUser).Methods("DELETE", "OPTIONS")
	users.HandleFunc("/{id:[0-9]+}/roles", s.handleAssignUserRoles).Methods("POST", "OPTIONS")
	users.HandleFunc("/{id:[0-9]+}/reset-password", s.handleResetUserPassword).Methods("POST", "OPTIONS")
	users.HandleFunc("/{id:[0-9]+}/toggle-2fa", s.handleToggleUser2FA).Methods("POST", "OPTIONS")

	// 角色管理路由
	roles := authenticated.PathPrefix("/roles").Subrouter()
	roles.HandleFunc("", s.handleListRoles).Methods("GET", "OPTIONS")
	roles.HandleFunc("", s.handleCreateRole).Methods("POST", "OPTIONS")
	roles.HandleFunc("/{id:[0-9]+}", s.handleGetRole).Methods("GET", "OPTIONS")
	roles.HandleFunc("/{id:[0-9]+}", s.handleUpdateRole).Methods("PUT", "OPTIONS")
	roles.HandleFunc("/{id:[0-9]+}", s.handleDeleteRole).Methods("DELETE", "OPTIONS")
	roles.HandleFunc("/{id:[0-9]+}/permissions", s.handleAssignRolePermissions).Methods("POST", "OPTIONS")
	roles.HandleFunc("/enabled", s.handleGetEnabledRoles).Methods("GET", "OPTIONS")

	// 权限管理路由
	permissions := authenticated.PathPrefix("/permissions").Subrouter()
	permissions.HandleFunc("", s.handleListPermissions).Methods("GET", "OPTIONS")
	permissions.HandleFunc("", s.handleCreatePermission).Methods("POST", "OPTIONS")
	permissions.HandleFunc("/{id:[0-9]+}", s.handleGetPermission).Methods("GET", "OPTIONS")
	permissions.HandleFunc("/{id:[0-9]+}", s.handleUpdatePermission).Methods("PUT", "OPTIONS")
	permissions.HandleFunc("/{id:[0-9]+}", s.handleDeletePermission).Methods("DELETE", "OPTIONS")
	permissions.HandleFunc("/tree", s.handleGetPermissionTree).Methods("GET", "OPTIONS")
	permissions.HandleFunc("/modules", s.handleGetPermissionModules).Methods("GET", "OPTIONS")
	permissions.HandleFunc("/sync-api", s.handleSyncApiPermissions).Methods("POST", "OPTIONS")
	permissions.HandleFunc("/menus", s.handleGetUserMenus).Methods("GET", "OPTIONS")
	permissions.HandleFunc("/check", s.handleCheckUserPermission).Methods("POST", "OPTIONS")

	// 组织管理路由
	orgs := authenticated.PathPrefix("/organizations").Subrouter()
	orgs.HandleFunc("", s.handleCreateOrganization).Methods("POST", "OPTIONS")
	orgs.HandleFunc("/{id:[0-9]+}", s.handleGetOrganization).Methods("GET", "OPTIONS")
	orgs.HandleFunc("/{id:[0-9]+}", s.handleUpdateOrganization).Methods("PUT", "OPTIONS")
	orgs.HandleFunc("/{id:[0-9]+}", s.handleDeleteOrganization).Methods("DELETE", "OPTIONS")
	orgs.HandleFunc("/tree", s.handleGetOrganizationTree).Methods("GET", "OPTIONS")
	orgs.HandleFunc("/enabled", s.handleGetEnabledOrganizations).Methods("GET", "OPTIONS")
	orgs.HandleFunc("/{id:[0-9]+}/users", s.handleAssignOrganizationUsers).Methods("POST", "OPTIONS")

	// 系统管理路由
	system := authenticated.PathPrefix("/system").Subrouter()
	system.HandleFunc("/logs", s.handleGetOperationLogs).Methods("GET", "OPTIONS")
	system.HandleFunc("/logs/{id:[0-9]+}", s.handleGetOperationLog).Methods("GET", "OPTIONS")
	system.HandleFunc("/statistics", s.handleGetOperationStatistics).Methods("GET", "OPTIONS")
	system.HandleFunc("/active-users", s.handleGetTopActiveUsers).Methods("GET", "OPTIONS")
	system.HandleFunc("/cleanup-logs", s.handleCleanupLogs).Methods("POST", "OPTIONS")
	system.HandleFunc("/info", s.handleGetSystemInfo).Methods("GET", "OPTIONS")
	system.HandleFunc("/dashboard", s.handleGetDashboardData).Methods("GET", "OPTIONS")

	// 健康检查
	router.HandleFunc("/health", s.handleHealth).Methods("GET")

	// 添加CORS中间件
	router.Use(CorsMiddleware())

	// 注册路由到Kratos HTTP服务器
	s.Server.HandlePrefix("/", router)
}

// JSON响应结构
type APIResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message,omitempty"`
	Error   *ErrorInfo  `json:"error,omitempty"`
}

type ErrorInfo struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

// 发送JSON响应
func (s *HTTPServer) sendResponse(w http.ResponseWriter, status int, data interface{}) {
	response := APIResponse{
		Success: status < 400,
		Data:    data,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(response)
}

// 发送错误响应
func (s *HTTPServer) sendError(w http.ResponseWriter, err error) {
	var status int
	var code string
	var message string

	if e := errors.FromError(err); e != nil {
		status = int(e.Code)
		code = e.Reason
		message = e.Message
	} else {
		status = http.StatusInternalServerError
		code = "INTERNAL_ERROR"
		message = "内部服务器错误"
	}

	// 根据Kratos错误码映射HTTP状态码
	switch {
	case status == 400:
		status = http.StatusBadRequest
	case status == 401:
		status = http.StatusUnauthorized
	case status == 403:
		status = http.StatusForbidden
	case status == 404:
		status = http.StatusNotFound
	case status >= 500:
		status = http.StatusInternalServerError
	}

	response := APIResponse{
		Success: false,
		Error: &ErrorInfo{
			Code:    code,
			Message: message,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(response)
}

// 解析JSON请求体
func (s *HTTPServer) parseJSON(r *http.Request, v interface{}) error {
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(v); err != nil {
		return errors.BadRequest("INVALID_JSON", "无效的JSON格式")
	}
	return nil
}

// 获取客户端IP
func (s *HTTPServer) getClientIP(r *http.Request) string {
	// 检查X-Forwarded-For头
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		ips := strings.Split(xff, ",")
		return strings.TrimSpace(ips[0])
	}

	// 检查X-Real-IP头
	if xri := r.Header.Get("X-Real-IP"); xri != "" {
		return xri
	}

	// 使用RemoteAddr
	ip := r.RemoteAddr
	if idx := strings.LastIndex(ip, ":"); idx != -1 {
		ip = ip[:idx]
	}

	return ip
}

// jwtMiddleware JWT中间件
func (s *HTTPServer) jwtMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 从Authorization头获取token
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			s.sendError(w, errors.Unauthorized("UNAUTHORIZED", "missing authorization header"))
			return
		}

		// 验证Bearer前缀
		if !strings.HasPrefix(authHeader, "Bearer ") {
			s.sendError(w, errors.Unauthorized("UNAUTHORIZED", "invalid authorization header format"))
			return
		}

		// 提取token
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// 解析JWT
		token, err := jwt.ParseWithClaims(tokenString, &pkg.CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(s.jwtSecret), nil
		})

		if err != nil {
			s.log.Warnf("Invalid token: %v", err)
			s.sendError(w, errors.Unauthorized("UNAUTHORIZED", "invalid token"))
			return
		}

		claims, ok := token.Claims.(*pkg.CustomClaims)
		if !ok || !token.Valid {
			s.sendError(w, errors.Unauthorized("UNAUTHORIZED", "invalid token claims"))
			return
		}

		// 检查token类型
		if claims.TokenType != "access" {
			s.sendError(w, errors.Unauthorized("UNAUTHORIZED", "invalid token type"))
			return
		}

		// 检查token是否过期
		if time.Now().After(claims.ExpiresAt.Time) {
			s.sendError(w, errors.Unauthorized("UNAUTHORIZED", "token expired"))
			return
		}

		// 将用户信息添加到请求上下文
		ctx := r.Context()
		ctx = middleware.SetUserIDToContext(ctx, claims.UserID)
		ctx = middleware.SetUsernameToContext(ctx, claims.Username)
		ctx = middleware.SetUserEmailToContext(ctx, claims.Email)
		ctx = middleware.SetUserRolesToContext(ctx, claims.Roles)
		ctx = middleware.SetSessionIDToContext(ctx, claims.SessionID)

		// 继续执行
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// handleLogin 处理登录
func (s *HTTPServer) handleLogin(w http.ResponseWriter, r *http.Request) {
	var req service.LoginRequest
	if err := s.parseJSON(r, &req); err != nil {
		s.sendError(w, err)
		return
	}

	// 设置客户端信息
	req.ClientIP = s.getClientIP(r)
	req.UserAgent = r.Header.Get("User-Agent")

	// 调用服务层
	resp, err := s.authService.Login(r.Context(), &req)
	if err != nil {
		s.sendError(w, err)
		return
	}

	s.sendResponse(w, http.StatusOK, resp)
}

// handleRegister 处理注册
func (s *HTTPServer) handleRegister(w http.ResponseWriter, r *http.Request) {
	var req service.RegisterRequest
	if err := s.parseJSON(r, &req); err != nil {
		s.sendError(w, err)
		return
	}

	// 调用服务层
	resp, err := s.authService.Register(r.Context(), &req)
	if err != nil {
		s.sendError(w, err)
		return
	}

	s.sendResponse(w, http.StatusCreated, resp)
}

// handleRefreshToken 处理刷新令牌
func (s *HTTPServer) handleRefreshToken(w http.ResponseWriter, r *http.Request) {
	var req service.RefreshTokenRequest
	if err := s.parseJSON(r, &req); err != nil {
		s.sendError(w, err)
		return
	}

	// 调用服务层
	resp, err := s.authService.RefreshToken(r.Context(), &req)
	if err != nil {
		s.sendError(w, err)
		return
	}

	s.sendResponse(w, http.StatusOK, resp)
}

// handleLogout 处理登出
func (s *HTTPServer) handleLogout(w http.ResponseWriter, r *http.Request) {
	var req service.LogoutRequest
	if err := s.parseJSON(r, &req); err != nil {
		s.sendError(w, err)
		return
	}

	// 调用服务层
	if err := s.authService.Logout(r.Context(), &req); err != nil {
		s.sendError(w, err)
		return
	}

	s.sendResponse(w, http.StatusOK, map[string]string{
		"message": "登出成功",
	})
}

// handleGetProfile 处理获取用户资料
func (s *HTTPServer) handleGetProfile(w http.ResponseWriter, r *http.Request) {
	// 调用服务层
	resp, err := s.authService.GetProfile(r.Context())
	if err != nil {
		s.sendError(w, err)
		return
	}

	s.sendResponse(w, http.StatusOK, resp)
}

// handleChangePassword 处理修改密码
func (s *HTTPServer) handleChangePassword(w http.ResponseWriter, r *http.Request) {
	var req service.ChangePasswordRequest
	if err := s.parseJSON(r, &req); err != nil {
		s.sendError(w, err)
		return
	}

	// 调用服务层
	if err := s.authService.ChangePassword(r.Context(), &req); err != nil {
		s.sendError(w, err)
		return
	}

	s.sendResponse(w, http.StatusOK, map[string]string{
		"message": "密码修改成功",
	})
}

// handleHealth 处理健康检查
func (s *HTTPServer) handleHealth(w http.ResponseWriter, r *http.Request) {
	health := map[string]interface{}{
		"status":    "healthy",
		"timestamp": time.Now().Unix(),
		"service":   "erp-system",
		"version":   "1.0.0",
	}

	s.sendResponse(w, http.StatusOK, health)
}

// CORS中间件
func CorsMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// 设置CORS头
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
			w.Header().Set("Access-Control-Max-Age", "86400")

			// 处理预检请求
			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
