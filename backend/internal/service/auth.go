package service

import (
	"context"
	"time"

	"erp-system/internal/biz"
	"erp-system/internal/middleware"
	"erp-system/internal/pkg"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/uuid"
)

// AuthService 认证服务
type AuthService struct {
	userUc *biz.UserUsecase
	jwtMgr *pkg.JWTManager
	pwdMgr *pkg.PasswordManager
	log    *log.Helper
}

// NewAuthService 创建认证服务
func NewAuthService(
	userUc *biz.UserUsecase,
	jwtMgr *pkg.JWTManager,
	pwdMgr *pkg.PasswordManager,
	logger log.Logger,
) *AuthService {
	return &AuthService{
		userUc: userUc,
		jwtMgr: jwtMgr,
		pwdMgr: pwdMgr,
		log:    log.NewHelper(logger),
	}
}

// LoginRequest 登录请求
type LoginRequest struct {
	Username      string `json:"username" validate:"required,min=3,max=32"`
	Password      string `json:"password" validate:"required,min=8"`
	TwoFactorCode string `json:"two_factor_code,omitempty"`
	RememberMe    bool   `json:"remember_me"`
	ClientIP      string `json:"-"`
	UserAgent     string `json:"-"`
}

// LoginResponse 登录响应
type LoginResponse struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresIn    int64     `json:"expires_in"`
	TokenType    string    `json:"token_type"`
	User         *UserInfo `json:"user"`
}

// RegisterRequest 注册请求
type RegisterRequest struct {
	Username        string `json:"username" validate:"required,min=3,max=32"`
	Email           string `json:"email" validate:"required,email"`
	Password        string `json:"password" validate:"required,min=8"`
	ConfirmPassword string `json:"confirm_password" validate:"required,eqfield=Password"`
	FirstName       string `json:"first_name" validate:"required,min=1,max=50"`
	LastName        string `json:"last_name" validate:"required,min=1,max=50"`
	Phone           string `json:"phone" validate:"omitempty,len=11"`
	Gender          string `json:"gender" validate:"omitempty,oneof=MALE FEMALE OTHER"`
	InvitationCode  string `json:"invitation_code,omitempty"`
}

// RegisterResponse 注册响应
type RegisterResponse struct {
	User    *UserInfo `json:"user"`
	Message string    `json:"message"`
}

// RefreshTokenRequest 刷新令牌请求
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

// RefreshTokenResponse 刷新令牌响应
type RefreshTokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int64  `json:"expires_in"`
	TokenType   string `json:"token_type"`
}

// LogoutRequest 登出请求
type LogoutRequest struct {
	RefreshToken string `json:"refresh_token,omitempty"`
}

// ChangePasswordRequest 修改密码请求
type ChangePasswordRequest struct {
	OldPassword     string `json:"old_password" validate:"required"`
	NewPassword     string `json:"new_password" validate:"required,min=8"`
	ConfirmPassword string `json:"confirm_password" validate:"required,eqfield=NewPassword"`
}

// UserInfo 用户信息
type UserInfo struct {
	ID               int32      `json:"id"`
	Username         string     `json:"username"`
	Email            string     `json:"email"`
	FirstName        string     `json:"first_name"`
	LastName         string     `json:"last_name"`
	Phone            string     `json:"phone"`
	Gender           string     `json:"gender"`
	BirthDate        *time.Time `json:"birth_date"`
	AvatarURL        string     `json:"avatar_url"`
	IsActive         bool       `json:"is_active"`
	TwoFactorEnabled bool       `json:"two_factor_enabled"`
	LastLoginAt      time.Time  `json:"last_login_at"`
	CreatedAt        time.Time  `json:"created_at"`
	Roles            []string   `json:"roles"`
	Permissions      []string   `json:"permissions"`
}

// ToUserInfo 将 biz.User 转换为 UserInfo
func ToUserInfo(user *biz.User, roles []string, permissions []string) *UserInfo {
	return &UserInfo{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Phone:     user.Phone,
		Gender:    user.Gender,
		BirthDate: func() *time.Time {
			if user.BirthDate.IsZero() {
				return nil
			}
			return &user.BirthDate
		}(),
		AvatarURL:        user.AvatarURL,
		IsActive:         user.IsActive,
		TwoFactorEnabled: user.TwoFactorEnabled,
		LastLoginAt:      user.LastLoginAt,
		CreatedAt:        user.CreatedAt,
		Roles:            roles,
		Permissions:      permissions,
	}
}

// Login 用户登录
func (s *AuthService) Login(ctx context.Context, req *LoginRequest) (*LoginResponse, error) {
	s.log.Infof("User login attempt: %s", req.Username)

	// 验证用户名格式
	if err := pkg.ValidateUsername(req.Username); err != nil {
		return nil, errors.BadRequest("INVALID_USERNAME", err.Error())
	}

	// 登录时不需要验证密码强度

	// 获取用户
	user, err := s.userUc.GetUserByUsername(ctx, req.Username)
	if err != nil {
		s.log.Warnf("User not found: %s", req.Username)
		return nil, errors.Unauthorized("INVALID_CREDENTIALS", "用户名或密码错误")
	}

	// 验证密码
	if !s.userUc.ValidatePassword(user.Password, req.Password) {
		s.log.Warnf("Invalid password for user: %s", req.Username)
		return nil, errors.Unauthorized("INVALID_CREDENTIALS", "用户名或密码错误")
	}

	// 检查账户状态
	if !user.IsActive {
		s.log.Warnf("Account disabled: %s", req.Username)
		return nil, errors.Forbidden("ACCOUNT_DISABLED", "账户已被禁用")
	}

	// 验证2FA（如果启用）
	if user.TwoFactorEnabled {
		if req.TwoFactorCode == "" {
			return nil, errors.BadRequest("MISSING_2FA_CODE", "需要二次验证码")
		}

		if !s.userUc.ValidateTwoFactor(ctx, user.ID, req.TwoFactorCode) {
			s.log.Warnf("Invalid 2FA code for user: %s", req.Username)
			return nil, errors.Unauthorized("INVALID_2FA_CODE", "二次验证码错误")
		}
	}

	// 获取用户角色和权限
	roles, err := s.userUc.GetUserRoles(ctx, user.ID)
	if err != nil {
		s.log.Errorf("Failed to get user roles: %v", err)
		return nil, errors.InternalServer("INTERNAL_ERROR", "系统错误")
	}

	permissions, err := s.userUc.GetUserPermissions(ctx, user.ID)
	if err != nil {
		s.log.Errorf("Failed to get user permissions: %v", err)
		return nil, errors.InternalServer("INTERNAL_ERROR", "系统错误")
	}

	// 生成会话ID
	sessionID := uuid.New().String()

	// 转换角色和权限为字符串数组
	roleStrs := make([]string, len(roles))
	for i, role := range roles {
		roleStrs[i] = role.Code
	}

	// 生成访问令牌
	tokenDuration := time.Hour * 2 // 2小时
	if req.RememberMe {
		tokenDuration = time.Hour * 24 * 7 // 7天
	}

	accessToken, err := s.jwtMgr.Generate(
		int64(user.ID), user.Username, user.Email,
		roleStrs, permissions, sessionID, "access",
	)
	if err != nil {
		s.log.Errorf("Failed to generate access token: %v", err)
		return nil, errors.InternalServer("TOKEN_GENERATION_ERROR", "令牌生成失败")
	}

	// 生成刷新令牌
	refreshToken, err := s.jwtMgr.GenerateRefreshToken(int64(user.ID), user.Username, sessionID)
	if err != nil {
		s.log.Errorf("Failed to generate refresh token: %v", err)
		return nil, errors.InternalServer("TOKEN_GENERATION_ERROR", "令牌生成失败")
	}

	// 更新登录信息
	if err := s.userUc.UpdateLoginInfo(ctx, user.ID, req.ClientIP); err != nil {
		s.log.Warnf("Failed to update login info: %v", err)
	}

	// TODO: 保存会话信息到Redis

	s.log.Infof("User login successful: %s", req.Username)

	return &LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int64(tokenDuration.Seconds()),
		TokenType:    "Bearer",
		User:         ToUserInfo(user, roleStrs, permissions),
	}, nil
}

// Register 用户注册
func (s *AuthService) Register(ctx context.Context, req *RegisterRequest) (*RegisterResponse, error) {
	s.log.Infof("User registration attempt: %s", req.Username)

	// 验证用户名格式
	if err := pkg.ValidateUsername(req.Username); err != nil {
		return nil, errors.BadRequest("INVALID_USERNAME", err.Error())
	}

	// 验证邮箱格式
	if !pkg.ValidateEmail(req.Email) {
		return nil, errors.BadRequest("INVALID_EMAIL", "邮箱格式不正确")
	}

	// 验证密码强度
	if err := s.pwdMgr.ValidatePasswordStrength(req.Password); err != nil {
		return nil, errors.BadRequest("WEAK_PASSWORD", err.Error())
	}

	// 验证手机号格式（如果提供）
	if req.Phone != "" && !pkg.ValidatePhone(req.Phone) {
		return nil, errors.BadRequest("INVALID_PHONE", "手机号格式不正确")
	}

	// 检查用户名是否存在
	if _, err := s.userUc.GetUserByUsername(ctx, req.Username); err == nil {
		return nil, errors.BadRequest("USERNAME_EXISTS", "用户名已存在")
	}

	// 检查邮箱是否存在
	if _, err := s.userUc.GetUserByEmail(ctx, req.Email); err == nil {
		return nil, errors.BadRequest("EMAIL_EXISTS", "邮箱已存在")
	}

	// 密码加密
	hashedPassword, err := s.userUc.HashPassword(req.Password)
	if err != nil {
		s.log.Errorf("Failed to hash password: %v", err)
		return nil, errors.InternalServer("INTERNAL_ERROR", "系统错误")
	}

	// 创建用户
	user := &biz.User{
		Username:  req.Username,
		Email:     req.Email,
		Password:  hashedPassword,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Phone:     req.Phone,
		Gender:    req.Gender,
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	createdUser, err := s.userUc.CreateUser(ctx, user)
	if err != nil {
		s.log.Errorf("Failed to create user: %v", err)
		return nil, errors.InternalServer("INTERNAL_ERROR", "用户创建失败")
	}

	// TODO: 发送邮箱验证邮件
	// TODO: 分配默认角色

	s.log.Infof("User registration successful: %s", req.Username)

	return &RegisterResponse{
		User: &UserInfo{
			ID:        createdUser.ID,
			Username:  createdUser.Username,
			Email:     createdUser.Email,
			FirstName: createdUser.FirstName,
			LastName:  createdUser.LastName,
			Phone:     createdUser.Phone,
			Gender:    createdUser.Gender,
			AvatarURL: createdUser.AvatarURL,
			IsActive:  createdUser.IsActive,
			CreatedAt: createdUser.CreatedAt,
		},
		Message: "注册成功，请检查邮箱完成验证",
	}, nil
}

// RefreshToken 刷新访问令牌
func (s *AuthService) RefreshToken(ctx context.Context, req *RefreshTokenRequest) (*RefreshTokenResponse, error) {
	// 解析刷新令牌
	claims, err := s.jwtMgr.ParseToken(req.RefreshToken)
	if err != nil {
		return nil, errors.Unauthorized("INVALID_REFRESH_TOKEN", "无效的刷新令牌")
	}

	// 检查令牌类型
	if claims.TokenType != "refresh" {
		return nil, errors.Unauthorized("INVALID_TOKEN_TYPE", "令牌类型错误")
	}

	// 获取用户信息
	user, err := s.userUc.GetUser(ctx, int32(claims.UserID))
	if err != nil {
		return nil, errors.Unauthorized("USER_NOT_FOUND", "用户不存在")
	}

	// 检查账户状态
	if !user.IsActive {
		return nil, errors.Forbidden("ACCOUNT_DISABLED", "账户已被禁用")
	}

	// 获取用户角色和权限
	roles, err := s.userUc.GetUserRoles(ctx, user.ID)
	if err != nil {
		s.log.Errorf("Failed to get user roles: %v", err)
		return nil, errors.InternalServer("INTERNAL_ERROR", "系统错误")
	}

	permissions, err := s.userUc.GetUserPermissions(ctx, user.ID)
	if err != nil {
		s.log.Errorf("Failed to get user permissions: %v", err)
		return nil, errors.InternalServer("INTERNAL_ERROR", "系统错误")
	}

	// 转换角色为字符串数组
	roleStrs := make([]string, len(roles))
	for i, role := range roles {
		roleStrs[i] = role.Code
	}

	// 生成新的访问令牌
	tokenDuration := time.Hour * 2 // 2小时
	accessToken, err := s.jwtMgr.Generate(
		int64(user.ID), user.Username, user.Email,
		roleStrs, permissions, claims.SessionID, "access",
	)
	if err != nil {
		s.log.Errorf("Failed to generate access token: %v", err)
		return nil, errors.InternalServer("TOKEN_GENERATION_ERROR", "令牌生成失败")
	}

	return &RefreshTokenResponse{
		AccessToken: accessToken,
		ExpiresIn:   int64(tokenDuration.Seconds()),
		TokenType:   "Bearer",
	}, nil
}

// Logout 用户登出
func (s *AuthService) Logout(ctx context.Context, req *LogoutRequest) error {
	// 获取当前用户信息
	currentUser := middleware.GetCurrentUser(ctx)
	if !currentUser.IsAuthenticated() {
		return errors.Unauthorized("NOT_AUTHENTICATED", "用户未认证")
	}

	s.log.Infof("User logout: %s", currentUser.Username)

	// TODO: 清除Redis中的会话信息
	// TODO: 将刷新令牌加入黑名单

	s.log.Infof("User logout successful: %s", currentUser.Username)
	return nil
}

// ChangePassword 修改密码
func (s *AuthService) ChangePassword(ctx context.Context, req *ChangePasswordRequest) error {
	// 获取当前用户信息
	currentUser := middleware.GetCurrentUser(ctx)
	if !currentUser.IsAuthenticated() {
		return errors.Unauthorized("NOT_AUTHENTICATED", "用户未认证")
	}

	// 获取用户详细信息
	user, err := s.userUc.GetUser(ctx, int32(currentUser.ID))
	if err != nil {
		return errors.NotFound("USER_NOT_FOUND", "用户不存在")
	}

	// 验证旧密码
	if !s.userUc.ValidatePassword(user.Password, req.OldPassword) {
		return errors.BadRequest("INVALID_OLD_PASSWORD", "原密码错误")
	}

	// 验证新密码强度
	if err := s.pwdMgr.ValidatePasswordStrength(req.NewPassword); err != nil {
		return errors.BadRequest("WEAK_PASSWORD", err.Error())
	}

	// 检查新密码是否与旧密码相同
	if req.OldPassword == req.NewPassword {
		return errors.BadRequest("SAME_PASSWORD", "新密码不能与原密码相同")
	}

	// 加密新密码
	hashedPassword, err := s.userUc.HashPassword(req.NewPassword)
	if err != nil {
		s.log.Errorf("Failed to hash password: %v", err)
		return errors.InternalServer("INTERNAL_ERROR", "系统错误")
	}

	// 更新密码
	user.Password = hashedPassword
	user.UpdatedAt = time.Now()

	if _, err := s.userUc.UpdateUser(ctx, user); err != nil {
		s.log.Errorf("Failed to update password: %v", err)
		return errors.InternalServer("INTERNAL_ERROR", "密码更新失败")
	}

	// TODO: 强制用户重新登录（清除所有会话）

	s.log.Infof("Password changed for user: %s", currentUser.Username)
	return nil
}

// GetProfile 获取用户资料
func (s *AuthService) GetProfile(ctx context.Context) (*UserInfo, error) {
	// 获取当前用户信息
	currentUser := middleware.GetCurrentUser(ctx)
	if !currentUser.IsAuthenticated() {
		return nil, errors.Unauthorized("NOT_AUTHENTICATED", "用户未认证")
	}

	// 获取用户详细信息
	user, err := s.userUc.GetUser(ctx, int32(currentUser.ID))
	if err != nil {
		return nil, errors.NotFound("USER_NOT_FOUND", "用户不存在")
	}

	// 获取用户角色和权限
	roles, err := s.userUc.GetUserRoles(ctx, user.ID)
	if err != nil {
		s.log.Errorf("Failed to get user roles: %v", err)
		return nil, errors.InternalServer("INTERNAL_ERROR", "系统错误")
	}

	permissions, err := s.userUc.GetUserPermissions(ctx, user.ID)
	if err != nil {
		s.log.Errorf("Failed to get user permissions: %v", err)
		return nil, errors.InternalServer("INTERNAL_ERROR", "系统错误")
	}

	// 转换角色为字符串数组
	roleStrs := make([]string, len(roles))
	for i, role := range roles {
		roleStrs[i] = role.Code
	}

	return &UserInfo{
		ID:               user.ID,
		Username:         user.Username,
		Email:            user.Email,
		FirstName:        user.FirstName,
		LastName:         user.LastName,
		Phone:            user.Phone,
		Gender:           user.Gender,
		AvatarURL:        user.AvatarURL,
		IsActive:         user.IsActive,
		TwoFactorEnabled: user.TwoFactorEnabled,
		LastLoginAt:      user.LastLoginAt,
		CreatedAt:        user.CreatedAt,
		Roles:            roleStrs,
		Permissions:      permissions,
	}, nil
}
