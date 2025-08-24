package middleware

import (
	"context"
)

// 上下文键类型
type contextKey string

const (
	// 用户相关上下文键
	UserIDKey    contextKey = "user_id"
	UsernameKey  contextKey = "username"
	UserEmailKey contextKey = "user_email"
	UserRolesKey contextKey = "user_roles"
	SessionIDKey contextKey = "session_id"

	// 权限相关上下文键
	PermissionsKey contextKey = "permissions"

	// 请求相关上下文键
	TraceIDKey   contextKey = "trace_id"
	RequestIDKey contextKey = "request_id"
	ClientIPKey  contextKey = "client_ip"
	UserAgentKey contextKey = "user_agent"
)

// SetUserIDToContext 设置用户ID到上下文
func SetUserIDToContext(ctx context.Context, userID int64) context.Context {
	return context.WithValue(ctx, UserIDKey, userID)
}

// GetUserIDFromContext 从上下文获取用户ID
func GetUserIDFromContext(ctx context.Context) int64 {
	if userID, ok := ctx.Value(UserIDKey).(int64); ok {
		return userID
	}
	return 0
}

// SetUsernameToContext 设置用户名到上下文
func SetUsernameToContext(ctx context.Context, username string) context.Context {
	return context.WithValue(ctx, UsernameKey, username)
}

// GetUsernameFromContext 从上下文获取用户名
func GetUsernameFromContext(ctx context.Context) string {
	if username, ok := ctx.Value(UsernameKey).(string); ok {
		return username
	}
	return ""
}

// SetUserEmailToContext 设置用户邮箱到上下文
func SetUserEmailToContext(ctx context.Context, email string) context.Context {
	return context.WithValue(ctx, UserEmailKey, email)
}

// GetUserEmailFromContext 从上下文获取用户邮箱
func GetUserEmailFromContext(ctx context.Context) string {
	if email, ok := ctx.Value(UserEmailKey).(string); ok {
		return email
	}
	return ""
}

// SetUserRolesToContext 设置用户角色到上下文
func SetUserRolesToContext(ctx context.Context, roles []string) context.Context {
	return context.WithValue(ctx, UserRolesKey, roles)
}

// GetUserRolesFromContext 从上下文获取用户角色
func GetUserRolesFromContext(ctx context.Context) []string {
	if roles, ok := ctx.Value(UserRolesKey).([]string); ok {
		return roles
	}
	return []string{}
}

// SetSessionIDToContext 设置会话ID到上下文
func SetSessionIDToContext(ctx context.Context, sessionID string) context.Context {
	return context.WithValue(ctx, SessionIDKey, sessionID)
}

// GetSessionIDFromContext 从上下文获取会话ID
func GetSessionIDFromContext(ctx context.Context) string {
	if sessionID, ok := ctx.Value(SessionIDKey).(string); ok {
		return sessionID
	}
	return ""
}

// SetPermissionsToContext 设置用户权限到上下文
func SetPermissionsToContext(ctx context.Context, permissions []string) context.Context {
	return context.WithValue(ctx, PermissionsKey, permissions)
}

// GetPermissionsFromContext 从上下文获取用户权限
func GetPermissionsFromContext(ctx context.Context) []string {
	if permissions, ok := ctx.Value(PermissionsKey).([]string); ok {
		return permissions
	}
	return []string{}
}

// SetTraceIDToContext 设置追踪ID到上下文
func SetTraceIDToContext(ctx context.Context, traceID string) context.Context {
	return context.WithValue(ctx, TraceIDKey, traceID)
}

// GetTraceIDFromContext 从上下文获取追踪ID
func GetTraceIDFromContext(ctx context.Context) string {
	if traceID, ok := ctx.Value(TraceIDKey).(string); ok {
		return traceID
	}
	return ""
}

// SetRequestIDToContext 设置请求ID到上下文
func SetRequestIDToContext(ctx context.Context, requestID string) context.Context {
	return context.WithValue(ctx, RequestIDKey, requestID)
}

// GetRequestIDFromContext 从上下文获取请求ID
func GetRequestIDFromContext(ctx context.Context) string {
	if requestID, ok := ctx.Value(RequestIDKey).(string); ok {
		return requestID
	}
	return ""
}

// SetClientIPToContext 设置客户端IP到上下文
func SetClientIPToContext(ctx context.Context, clientIP string) context.Context {
	return context.WithValue(ctx, ClientIPKey, clientIP)
}

// GetClientIPFromContext 从上下文获取客户端IP
func GetClientIPFromContext(ctx context.Context) string {
	if clientIP, ok := ctx.Value(ClientIPKey).(string); ok {
		return clientIP
	}
	return ""
}

// SetUserAgentToContext 设置用户代理到上下文
func SetUserAgentToContext(ctx context.Context, userAgent string) context.Context {
	return context.WithValue(ctx, UserAgentKey, userAgent)
}

// GetUserAgentFromContext 从上下文获取用户代理
func GetUserAgentFromContext(ctx context.Context) string {
	if userAgent, ok := ctx.Value(UserAgentKey).(string); ok {
		return userAgent
	}
	return ""
}

// GetCurrentUser 获取当前用户信息
func GetCurrentUser(ctx context.Context) *CurrentUser {
	return &CurrentUser{
		ID:       GetUserIDFromContext(ctx),
		Username: GetUsernameFromContext(ctx),
		Email:    GetUserEmailFromContext(ctx),
		Roles:    GetUserRolesFromContext(ctx),
		SessionID: GetSessionIDFromContext(ctx),
	}
}

// CurrentUser 当前用户信息
type CurrentUser struct {
	ID        int64    `json:"id"`
	Username  string   `json:"username"`
	Email     string   `json:"email"`
	Roles     []string `json:"roles"`
	SessionID string   `json:"session_id"`
}

// IsAuthenticated 检查用户是否已认证
func (u *CurrentUser) IsAuthenticated() bool {
	return u.ID > 0 && u.Username != ""
}

// HasRole 检查用户是否拥有指定角色
func (u *CurrentUser) HasRole(role string) bool {
	for _, r := range u.Roles {
		if r == role {
			return true
		}
	}
	return false
}

// HasAnyRole 检查用户是否拥有任意一个指定角色
func (u *CurrentUser) HasAnyRole(roles ...string) bool {
	for _, role := range roles {
		if u.HasRole(role) {
			return true
		}
	}
	return false
}

// IsSuperAdmin 检查用户是否为超级管理员
func (u *CurrentUser) IsSuperAdmin() bool {
	return u.HasRole("SUPER_ADMIN")
}

// IsAdmin 检查用户是否为管理员
func (u *CurrentUser) IsAdmin() bool {
	return u.HasAnyRole("SUPER_ADMIN", "ADMIN")
}