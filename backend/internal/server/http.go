package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
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
	permissionService *service.PermissionService
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
	permissionService *service.PermissionService,
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
		permissionService:   permissionService,
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

	// 传统权限管理路由已禁用 - 使用ERP权限系统
	// permissions := authenticated.PathPrefix("/permissions").Subrouter()
	// permissions.HandleFunc("", s.handleListPermissions).Methods("GET", "OPTIONS")
	// permissions.HandleFunc("", s.handleCreatePermission).Methods("POST", "OPTIONS")
	// permissions.HandleFunc("/{id:[0-9]+}", s.handleGetPermission).Methods("GET", "OPTIONS")
	// permissions.HandleFunc("/{id:[0-9]+}", s.handleUpdatePermission).Methods("PUT", "OPTIONS")
	// permissions.HandleFunc("/{id:[0-9]+}", s.handleDeletePermission).Methods("DELETE", "OPTIONS")
	// permissions.HandleFunc("/tree", s.handleGetPermissionTree).Methods("GET", "OPTIONS")
	// permissions.HandleFunc("/modules", s.handleGetPermissionModules).Methods("GET", "OPTIONS")
	// permissions.HandleFunc("/sync-api", s.handleSyncApiPermissions).Methods("POST", "OPTIONS")
	// permissions.HandleFunc("/menus", s.handleGetUserMenus).Methods("GET", "OPTIONS")
	// permissions.HandleFunc("/check", s.handleCheckUserPermission).Methods("POST", "OPTIONS")

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

	// DocType管理路由（标准系统）
	docTypes := authenticated.PathPrefix("/doctypes").Subrouter()
	docTypes.HandleFunc("", s.handleGetDocTypeList).Methods("GET", "OPTIONS")
	docTypes.HandleFunc("/modules", s.handleGetDocTypeModules).Methods("GET", "OPTIONS")
	docTypes.HandleFunc("", s.handleCreateDocTypeManagement).Methods("POST", "OPTIONS")
	docTypes.HandleFunc("/{id:[0-9]+}", s.handleGetDocTypeManagement).Methods("GET", "OPTIONS")
	docTypes.HandleFunc("/{id:[0-9]+}", s.handleUpdateDocTypeManagement).Methods("PUT", "OPTIONS")
	docTypes.HandleFunc("/{id:[0-9]+}", s.handleDeleteDocTypeManagement).Methods("DELETE", "OPTIONS")

	// ERP文档权限系统路由
	erpPermissions := authenticated.PathPrefix("/erp-permissions").Subrouter()
	erpPermissions.HandleFunc("/doctypes", s.handleGetDocTypes).Methods("GET", "OPTIONS")
	erpPermissions.HandleFunc("/doctypes", s.handleCreateDocType).Methods("POST", "OPTIONS")
	erpPermissions.HandleFunc("/doctypes/{name}", s.handleGetDocType).Methods("GET", "OPTIONS")
	erpPermissions.HandleFunc("/doctypes/{name}", s.handleUpdateDocType).Methods("PUT", "OPTIONS")
	erpPermissions.HandleFunc("/doctypes/{name}", s.handleDeleteDocType).Methods("DELETE", "OPTIONS")
	erpPermissions.HandleFunc("/permission-rules", s.handleGetPermissionRules).Methods("GET", "OPTIONS")
	erpPermissions.HandleFunc("/permission-rules", s.handleCreatePermissionRule).Methods("POST", "OPTIONS")
	erpPermissions.HandleFunc("/permission-rules/{id:[0-9]+}", s.handleGetPermissionRule).Methods("GET", "OPTIONS")
	erpPermissions.HandleFunc("/permission-rules/{id:[0-9]+}", s.handleUpdatePermissionRule).Methods("PUT", "OPTIONS")
	erpPermissions.HandleFunc("/permission-rules/{id:[0-9]+}", s.handleDeletePermissionRule).Methods("DELETE", "OPTIONS")

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

// ========== ERP权限系统处理函数 ==========

// DocType相关
func (s *HTTPServer) handleGetDocTypes(w http.ResponseWriter, r *http.Request) {
	// 模拟DocType数据
	docTypes := []map[string]interface{}{
		{
			"id":                     1,
			"name":                   "User",
			"label":                  "用户",
			"module":                 "Core",
			"description":            "系统用户信息",
			"is_submittable":         false,
			"is_child_table":         false,
			"has_workflow":           false,
			"track_changes":          true,
			"applies_to_all_users":   true,
			"max_attachments":        10,
			"naming_rule":            "field:username",
			"title_field":            "full_name",
			"search_fields":          []string{"username", "email", "full_name"},
			"sort_field":             "created_at",
			"sort_order":             "DESC",
			"version":                1,
			"created_at":             "2024-01-01T00:00:00Z",
			"updated_at":             "2024-01-01T00:00:00Z",
		},
		{
			"id":                     2,
			"name":                   "Role",
			"label":                  "角色",
			"module":                 "Core",
			"description":            "系统角色信息",
			"is_submittable":         false,
			"is_child_table":         false,
			"has_workflow":           false,
			"track_changes":          true,
			"applies_to_all_users":   true,
			"max_attachments":        5,
			"naming_rule":            "field:name",
			"title_field":            "name",
			"search_fields":          []string{"name", "description"},
			"sort_field":             "sort_order",
			"sort_order":             "ASC",
			"version":                1,
			"created_at":             "2024-01-01T00:00:00Z",
			"updated_at":             "2024-01-01T00:00:00Z",
		},
		{
			"id":                     3,
			"name":                   "Permission",
			"label":                  "权限",
			"module":                 "Core",
			"description":            "系统权限信息",
			"is_submittable":         false,
			"is_child_table":         false,
			"has_workflow":           false,
			"track_changes":          true,
			"applies_to_all_users":   true,
			"max_attachments":        0,
			"naming_rule":            "field:name",
			"title_field":            "name",
			"search_fields":          []string{"name", "resource", "action"},
			"sort_field":             "created_at",
			"sort_order":             "DESC",
			"version":                1,
			"created_at":             "2024-01-01T00:00:00Z",
			"updated_at":             "2024-01-01T00:00:00Z",
		},
		{
			"id":                     4,
			"name":                   "Organization",
			"label":                  "组织",
			"module":                 "Core",
			"description":            "组织架构信息",
			"is_submittable":         false,
			"is_child_table":         false,
			"has_workflow":           false,
			"track_changes":          true,
			"applies_to_all_users":   true,
			"max_attachments":        20,
			"naming_rule":            "field:name",
			"title_field":            "name",
			"search_fields":          []string{"name", "code"},
			"sort_field":             "sort_order",
			"sort_order":             "ASC",
			"version":                1,
			"created_at":             "2024-01-01T00:00:00Z",
			"updated_at":             "2024-01-01T00:00:00Z",
		},
	}

	// 分页参数处理
	page := 1
	size := 10
	if pageStr := r.URL.Query().Get("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}
	if sizeStr := r.URL.Query().Get("size"); sizeStr != "" {
		if s, err := strconv.Atoi(sizeStr); err == nil && s > 0 && s <= 100 {
			size = s
		}
	}

	// 模块过滤
	moduleFilter := r.URL.Query().Get("module")
	filteredDocTypes := docTypes
	if moduleFilter != "" {
		var filtered []map[string]interface{}
		for _, dt := range docTypes {
			if dt["module"] == moduleFilter {
				filtered = append(filtered, dt)
			}
		}
		filteredDocTypes = filtered
	}

	responseData := map[string]interface{}{
		"doc_types": filteredDocTypes,
		"total":     len(filteredDocTypes),
		"page":      page,
		"size":      size,
	}

	s.sendResponse(w, http.StatusOK, responseData)
}

func (s *HTTPServer) handleCreateDocType(w http.ResponseWriter, r *http.Request) {
	s.sendResponse(w, http.StatusNotImplemented, map[string]string{"message": "Not implemented yet"})
}

func (s *HTTPServer) handleGetDocType(w http.ResponseWriter, r *http.Request) {
	s.sendResponse(w, http.StatusNotImplemented, map[string]string{"message": "Not implemented yet"})
}

func (s *HTTPServer) handleUpdateDocType(w http.ResponseWriter, r *http.Request) {
	s.sendResponse(w, http.StatusNotImplemented, map[string]string{"message": "Not implemented yet"})
}

func (s *HTTPServer) handleDeleteDocType(w http.ResponseWriter, r *http.Request) {
	s.sendResponse(w, http.StatusNotImplemented, map[string]string{"message": "Not implemented yet"})
}

// PermissionRule相关
func (s *HTTPServer) handleGetPermissionRules(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	
	// 解析查询参数
	roleIDStr := r.URL.Query().Get("role_id")
	docType := r.URL.Query().Get("doc_type")
	
	var roleID int64
	var err error
	if roleIDStr != "" {
		roleID, err = strconv.ParseInt(roleIDStr, 10, 64)
		if err != nil {
			s.sendError(w, errors.BadRequest("INVALID_PARAMETER", "角色ID无效"))
			return
		}
	}
	
	// 构造请求
	req := &service.ListPermissionRulesRequest{
		RoleID:  roleID,
		DocType: docType,
	}
	
	// 调用权限服务获取权限规则列表
	response, err := s.permissionService.ListPermissionRules(ctx, req)
	if err != nil {
		s.sendError(w, err)
		return
	}
	
	s.sendResponse(w, http.StatusOK, response)
}

func (s *HTTPServer) handleCreatePermissionRule(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	
	// 解析请求体
	var req service.CreatePermissionRuleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.sendError(w, errors.BadRequest("INVALID_JSON", "请求数据格式不正确"))
		return
	}
	
	// 调用权限服务创建权限规则
	rule, err := s.permissionService.CreatePermissionRule(ctx, &req)
	if err != nil {
		s.sendError(w, err)
		return
	}
	
	s.sendResponse(w, http.StatusCreated, rule)
}

func (s *HTTPServer) handleGetPermissionRule(w http.ResponseWriter, r *http.Request) {
	s.sendResponse(w, http.StatusNotImplemented, map[string]string{"message": "Not implemented yet"})
}

func (s *HTTPServer) handleUpdatePermissionRule(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	
	// 从URL路径中获取权限规则ID
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		s.sendError(w, errors.BadRequest("INVALID_PARAMETER", "权限规则ID无效"))
		return
	}
	
	// 解析请求体
	var reqBody service.CreatePermissionRuleRequest
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		s.sendError(w, errors.BadRequest("INVALID_JSON", "请求数据格式不正确"))
		return
	}
	
	// 构造更新请求
	req := service.UpdatePermissionRuleRequest{
		ID: id,
		CreatePermissionRuleRequest: reqBody,
	}
	
	// 调用权限服务更新权限规则
	rule, err := s.permissionService.UpdatePermissionRule(ctx, &req)
	if err != nil {
		s.sendError(w, err)
		return
	}
	
	s.sendResponse(w, http.StatusOK, rule)
}

func (s *HTTPServer) handleDeletePermissionRule(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	
	// 从URL路径中获取权限规则ID
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		s.sendError(w, errors.BadRequest("INVALID_PARAMETER", "权限规则ID无效"))
		return
	}
	
	// 调用权限服务删除权限规则
	err = s.permissionService.DeletePermissionRule(ctx, id)
	if err != nil {
		s.sendError(w, err)
		return
	}
	
	s.sendResponse(w, http.StatusOK, map[string]string{"message": "权限规则删除成功"})
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

// ========== DocType管理系统处理函数 ==========

// handleGetDocTypeList 获取DocType列表（用于DocType管理页面）
func (s *HTTPServer) handleGetDocTypeList(w http.ResponseWriter, r *http.Request) {
	// 解析查询参数
	page := 1
	pageSize := 20
	sortBy := "created_at"
	sortOrder := "desc"
	
	if pageStr := r.URL.Query().Get("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}
	if sizeStr := r.URL.Query().Get("page_size"); sizeStr != "" {
		if s, err := strconv.Atoi(sizeStr); err == nil && s > 0 && s <= 100 {
			pageSize = s
		}
	}
	if sort := r.URL.Query().Get("sort_by"); sort != "" {
		sortBy = sort
	}
	if order := r.URL.Query().Get("sort_order"); order != "" {
		sortOrder = order
	}

	// 模拟DocType数据
	docTypes := []map[string]interface{}{
		{
			"id":           1,
			"name":         "User",
			"label":        "用户",
			"module":       "Core",
			"description":  "系统用户管理",
			"is_enabled":   true,
			"is_custom":    false,
			"has_workflow": false,
			"track_changes": true,
			"version":      1,
			"created_at":   "2024-01-01T00:00:00Z",
			"updated_at":   "2024-01-01T00:00:00Z",
		},
		{
			"id":           2,
			"name":         "Role",
			"label":        "角色",
			"module":       "Core", 
			"description":  "系统角色管理",
			"is_enabled":   true,
			"is_custom":    false,
			"has_workflow": false,
			"track_changes": false,
			"version":      1,
			"created_at":   "2024-01-02T00:00:00Z",
			"updated_at":   "2024-01-02T00:00:00Z",
		},
		{
			"id":           3,
			"name":         "Permission",
			"label":        "权限",
			"module":       "Core",
			"description":  "系统权限管理",
			"is_enabled":   true,
			"is_custom":    false,
			"has_workflow": true,
			"track_changes": true,
			"version":      1,
			"created_at":   "2024-01-03T00:00:00Z",
			"updated_at":   "2024-01-03T00:00:00Z",
		},
		{
			"id":           4,
			"name":         "Organization",
			"label":        "组织",
			"module":       "Core",
			"description":  "组织架构管理",
			"is_enabled":   true,
			"is_custom":    false,
			"has_workflow": false,
			"track_changes": true,
			"version":      1,
			"created_at":   "2024-01-04T00:00:00Z",
			"updated_at":   "2024-01-04T00:00:00Z",
		},
		{
			"id":           5,
			"name":         "Product",
			"label":        "产品",
			"module":       "Inventory",
			"description":  "产品信息管理",
			"is_enabled":   true,
			"is_custom":    true,
			"has_workflow": true,
			"track_changes": false,
			"version":      1,
			"created_at":   "2024-01-05T00:00:00Z",
			"updated_at":   "2024-01-05T00:00:00Z",
		},
		{
			"id":           6,
			"name":         "Customer",
			"label":        "客户",
			"module":       "CRM",
			"description":  "客户关系管理",
			"is_enabled":   true,
			"is_custom":    true,
			"has_workflow": false,
			"track_changes": false,
			"version":      1,
			"created_at":   "2024-01-06T00:00:00Z",
			"updated_at":   "2024-01-06T00:00:00Z",
		},
	}

	s.sendResponse(w, http.StatusOK, map[string]interface{}{
		"items":      docTypes,
		"total":      len(docTypes),
		"page":       page,
		"page_size":  pageSize,
		"sort_by":    sortBy,
		"sort_order": sortOrder,
	})
}

// handleGetDocTypeModules 获取DocType模块列表
func (s *HTTPServer) handleGetDocTypeModules(w http.ResponseWriter, r *http.Request) {
	modules := []map[string]interface{}{
		{
			"name":        "Core",
			"label":       "核心模块",
			"description": "系统核心功能模块",
			"count":       4,
		},
		{
			"name":        "Inventory",
			"label":       "库存管理",
			"description": "产品和库存管理模块",
			"count":       1,
		},
		{
			"name":        "CRM",
			"label":       "客户关系管理",
			"description": "客户和销售管理模块",
			"count":       1,
		},
	}

	s.sendResponse(w, http.StatusOK, map[string]interface{}{
		"modules": modules,
		"total":   len(modules),
	})
}

// handleCreateDocTypeManagement 创建DocType
func (s *HTTPServer) handleCreateDocTypeManagement(w http.ResponseWriter, r *http.Request) {
	s.sendResponse(w, http.StatusCreated, map[string]interface{}{
		"doctype": map[string]interface{}{
			"id":          7,
			"name":        "NewDocType",
			"label":       "新文档类型",
			"module":      "Custom",
			"description": "用户创建的新文档类型",
			"is_enabled":  true,
			"is_custom":   true,
			"version":     1,
			"created_at":  time.Now().Format(time.RFC3339),
			"updated_at":  time.Now().Format(time.RFC3339),
		},
	})
}

// handleGetDocTypeManagement 获取单个DocType
func (s *HTTPServer) handleGetDocTypeManagement(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	
	id, err := strconv.Atoi(idStr)
	if err != nil {
		s.sendResponse(w, http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"error": map[string]interface{}{
				"code":    "INVALID_PARAMETER",
				"message": "无效的DocType ID",
			},
		})
		return
	}

	docType := map[string]interface{}{
		"id":          id,
		"name":        "User",
		"label":       "用户",
		"module":      "Core",
		"description": "系统用户管理",
		"is_enabled":  true,
		"is_custom":   false,
		"version":     1,
		"created_at":  "2024-01-01T00:00:00Z",
		"updated_at":  "2024-01-01T00:00:00Z",
	}

	s.sendResponse(w, http.StatusOK, map[string]interface{}{
		"doctype": docType,
	})
}

// handleUpdateDocTypeManagement 更新DocType
func (s *HTTPServer) handleUpdateDocTypeManagement(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	
	id, err := strconv.Atoi(idStr)
	if err != nil {
		s.sendResponse(w, http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"error": map[string]interface{}{
				"code":    "INVALID_PARAMETER",
				"message": "无效的DocType ID",
			},
		})
		return
	}

	s.sendResponse(w, http.StatusOK, map[string]interface{}{
		"doctype": map[string]interface{}{
			"id":         id,
			"name":       "UpdatedDocType",
			"label":      "更新后的文档类型",
			"updated_at": time.Now().Format(time.RFC3339),
		},
	})
}

// handleDeleteDocTypeManagement 删除DocType
func (s *HTTPServer) handleDeleteDocTypeManagement(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	
	id, err := strconv.Atoi(idStr)
	if err != nil {
		s.sendResponse(w, http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"error": map[string]interface{}{
				"code":    "INVALID_PARAMETER",
				"message": "无效的DocType ID",
			},
		})
		return
	}

	s.sendResponse(w, http.StatusOK, map[string]interface{}{
		"message": fmt.Sprintf("DocType %d 已删除", id),
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
