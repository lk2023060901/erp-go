package service

import (
	"context"
	"time"

	"erp-system/internal/biz"
	"erp-system/internal/middleware"
	"erp-system/internal/pkg"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
)

// UserService 用户服务
type UserService struct {
	userUc *biz.UserUsecase
	pwdMgr *pkg.PasswordManager
	log    *log.Helper
}

// NewUserService 创建用户服务
func NewUserService(
	userUc *biz.UserUsecase,
	pwdMgr *pkg.PasswordManager,
	logger log.Logger,
) *UserService {
	return &UserService{
		userUc: userUc,
		pwdMgr: pwdMgr,
		log:    log.NewHelper(logger),
	}
}

// CreateUserRequest 创建用户请求
type CreateUserRequest struct {
	Username    string   `json:"username" validate:"required,min=3,max=32"`
	Email       string   `json:"email" validate:"required,email"`
	Password    string   `json:"password" validate:"required,min=8"`
	FirstName   string   `json:"first_name" validate:"omitempty,min=1,max=50"`
	LastName    string   `json:"last_name" validate:"omitempty,min=1,max=50"`
	Phone       string   `json:"phone" validate:"omitempty,len=11"`
	Gender      string   `json:"gender" validate:"omitempty,oneof=MALE FEMALE OTHER"`
	IsActive    bool     `json:"is_active"`
	RoleIDs     []int32  `json:"role_ids" validate:"omitempty,min=1"`
}

// UpdateUserRequest 更新用户请求
type UpdateUserRequest struct {
	ID         int32   `json:"id" validate:"required"`
	Email      string  `json:"email" validate:"required,email"`
	FirstName  string  `json:"first_name" validate:"required,min=1,max=50"`
	LastName   string  `json:"last_name" validate:"required,min=1,max=50"`
	Phone      string  `json:"phone" validate:"omitempty,len=11"`
	Gender     string  `json:"gender" validate:"omitempty,oneof=MALE FEMALE OTHER"`
	BirthDate  string  `json:"birth_date" validate:"omitempty"`
	AvatarURL  string  `json:"avatar_url"`
	IsActive   bool    `json:"is_active"`
	RoleIDs    []int32 `json:"role_ids"`
}

// UserListRequest 用户列表请求
type UserListRequest struct {
	Page             int32                  `json:"page" validate:"min=1"`
	Size             int32                  `json:"size" validate:"min=1,max=100"`
	Search           string                 `json:"search"` // 兼容旧版搜索
	FilterConditions map[string]interface{} `json:"filter_conditions"` // 新的过滤条件
	SortConfig       map[string]interface{} `json:"sort_config"`       // 排序配置
	FilterID         int32                  `json:"filter_id"`         // 使用保存的过滤器ID
}

// UserListResponse 用户列表响应
type UserListResponse struct {
	Users []*UserInfo `json:"users"`
	Total int32       `json:"total"`
	Page  int32       `json:"page"`
	Size  int32       `json:"size"`
}

// AssignRolesRequest 分配角色请求
type AssignRolesRequest struct {
	UserID  int32   `json:"user_id" validate:"required"`
	RoleIDs []int32 `json:"role_ids" validate:"required"`
}

// ResetPasswordRequest 重置密码请求
type ResetPasswordRequest struct {
	UserID      int32  `json:"user_id" validate:"required"`
	NewPassword string `json:"new_password" validate:"required,min=8"`
}

// ToggleTwoFactorRequest 切换2FA请求
type ToggleTwoFactorRequest struct {
	UserID int32 `json:"user_id" validate:"required"`
	Enable bool  `json:"enable"`
}

// CreateUser 创建用户
func (s *UserService) CreateUser(ctx context.Context, req *CreateUserRequest) (*UserInfo, error) {
	// 检查权限
	currentUser := middleware.GetCurrentUser(ctx)
	if !currentUser.HasAnyRole("SUPER_ADMIN", "ADMIN", "USER_MANAGER") {
		return nil, errors.Forbidden("PERMISSION_DENIED", "无权限创建用户")
	}

	s.log.Infof("Creating user: %s by %s", req.Username, currentUser.Username)

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
		return nil, errors.InternalServer("INTERNAL_ERROR", "密码加密失败")
	}

	// 创建用户
	user := &biz.User{
		Username:   req.Username,
		Email:      req.Email,
		Password:   hashedPassword,
		FirstName:  req.FirstName,
		LastName:   req.LastName,
		Phone:      req.Phone,
		Gender:     req.Gender,
		IsActive:   req.IsActive,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	createdUser, err := s.userUc.CreateUser(ctx, user)
	if err != nil {
		s.log.Errorf("Failed to create user: %v", err)
		return nil, errors.InternalServer("INTERNAL_ERROR", "用户创建失败")
	}

	// 分配角色
	if err := s.userUc.AssignRoles(ctx, createdUser.ID, req.RoleIDs); err != nil {
		s.log.Errorf("Failed to assign roles: %v", err)
		return nil, errors.InternalServer("INTERNAL_ERROR", "角色分配失败")
	}

	// 获取角色和权限
	roles, _ := s.userUc.GetUserRoles(ctx, createdUser.ID)
	permissions, _ := s.userUc.GetUserPermissions(ctx, createdUser.ID)

	roleStrs := make([]string, len(roles))
	for i, role := range roles {
		roleStrs[i] = role.Code
	}

	s.log.Infof("User created successfully: %s", req.Username)

	return ToUserInfo(createdUser, roleStrs, permissions), nil
}

// GetUser 获取用户详情
func (s *UserService) GetUser(ctx context.Context, userID int32) (*UserInfo, error) {
	// 检查权限
	currentUser := middleware.GetCurrentUser(ctx)
	if !currentUser.HasAnyRole("SUPER_ADMIN", "ADMIN", "USER_MANAGER") && currentUser.ID != int64(userID) {
		return nil, errors.Forbidden("PERMISSION_DENIED", "无权限查看该用户")
	}

	// 获取用户
	user, err := s.userUc.GetUser(ctx, userID)
	if err != nil {
		return nil, errors.NotFound("USER_NOT_FOUND", "用户不存在")
	}

	// 获取角色和权限
	roles, _ := s.userUc.GetUserRoles(ctx, user.ID)
	permissions, _ := s.userUc.GetUserPermissions(ctx, user.ID)

	roleStrs := make([]string, len(roles))
	for i, role := range roles {
		roleStrs[i] = role.Code
	}

	return ToUserInfo(user, roleStrs, permissions), nil
}

// UpdateUser 更新用户
func (s *UserService) UpdateUser(ctx context.Context, req *UpdateUserRequest) (*UserInfo, error) {
	// 检查权限
	currentUser := middleware.GetCurrentUser(ctx)
	if !currentUser.HasAnyRole("SUPER_ADMIN", "ADMIN", "USER_MANAGER") && currentUser.ID != int64(req.ID) {
		return nil, errors.Forbidden("PERMISSION_DENIED", "无权限修改该用户")
	}

	s.log.Infof("Updating user: %d by %s", req.ID, currentUser.Username)

	// 获取原用户信息
	user, err := s.userUc.GetUser(ctx, req.ID)
	if err != nil {
		return nil, errors.NotFound("USER_NOT_FOUND", "用户不存在")
	}

	// 验证邮箱格式
	if !pkg.ValidateEmail(req.Email) {
		return nil, errors.BadRequest("INVALID_EMAIL", "邮箱格式不正确")
	}

	// 验证手机号格式（如果提供）
	if req.Phone != "" && !pkg.ValidatePhone(req.Phone) {
		return nil, errors.BadRequest("INVALID_PHONE", "手机号格式不正确")
	}

	// 检查邮箱是否被其他用户使用
	if existingUser, err := s.userUc.GetUserByEmail(ctx, req.Email); err == nil && existingUser.ID != req.ID {
		return nil, errors.BadRequest("EMAIL_EXISTS", "邮箱已被使用")
	}

	// 更新用户信息
	user.Email = req.Email
	user.FirstName = req.FirstName
	user.LastName = req.LastName
	user.Phone = req.Phone
	user.Gender = req.Gender
	user.AvatarURL = req.AvatarURL
	user.UpdatedAt = time.Now()
	
	// 处理生日字段
	if req.BirthDate != "" {
		if birthDate, err := time.Parse("2006-01-02", req.BirthDate); err == nil {
			user.BirthDate = birthDate
		}
	}

	// 只有管理员可以修改状态
	if currentUser.HasAnyRole("SUPER_ADMIN", "ADMIN", "USER_MANAGER") {
		user.IsActive = req.IsActive
	}

	updatedUser, err := s.userUc.UpdateUser(ctx, user)
	if err != nil {
		s.log.Errorf("Failed to update user: %v", err)
		return nil, errors.InternalServer("INTERNAL_ERROR", "用户更新失败")
	}

	// 更新角色（只有管理员可以）
	if len(req.RoleIDs) > 0 && currentUser.HasAnyRole("SUPER_ADMIN", "ADMIN", "USER_MANAGER") {
		if err := s.userUc.AssignRoles(ctx, updatedUser.ID, req.RoleIDs); err != nil {
			s.log.Errorf("Failed to assign roles: %v", err)
		}
	}

	// 获取最新的角色和权限
	roles, _ := s.userUc.GetUserRoles(ctx, updatedUser.ID)
	permissions, _ := s.userUc.GetUserPermissions(ctx, updatedUser.ID)

	roleStrs := make([]string, len(roles))
	for i, role := range roles {
		roleStrs[i] = role.Code
	}

	s.log.Infof("User updated successfully: %d", req.ID)

	return ToUserInfo(updatedUser, roleStrs, permissions), nil
}

// DeleteUser 删除用户
func (s *UserService) DeleteUser(ctx context.Context, userID int32) error {
	// 检查权限
	currentUser := middleware.GetCurrentUser(ctx)
	s.log.Infof("User %s (roles: %v) attempting to delete user %d", currentUser.Username, currentUser.Roles, userID)
	if !currentUser.HasAnyRole("SUPER_ADMIN") {
		s.log.Warnf("Permission denied: User %s does not have SUPER_ADMIN role", currentUser.Username)
		return errors.Forbidden("PERMISSION_DENIED", "无权限删除用户")
	}

	// 不能删除自己
	if currentUser.ID == int64(userID) {
		return errors.BadRequest("CANNOT_DELETE_SELF", "不能删除自己")
	}

	s.log.Infof("Deleting user: %d by %s", userID, currentUser.Username)

	// 检查用户是否存在
	if _, err := s.userUc.GetUser(ctx, userID); err != nil {
		return errors.NotFound("USER_NOT_FOUND", "用户不存在")
	}

	// 删除用户
	if err := s.userUc.DeleteUser(ctx, userID); err != nil {
		s.log.Errorf("Failed to delete user: %v", err)
		return errors.InternalServer("INTERNAL_ERROR", "用户删除失败")
	}

	s.log.Infof("User deleted successfully: %d", userID)
	return nil
}

// ListUsers 获取用户列表
func (s *UserService) ListUsers(ctx context.Context, req *UserListRequest) (*UserListResponse, error) {
	// 检查权限
	currentUser := middleware.GetCurrentUser(ctx)
	if !currentUser.HasAnyRole("SUPER_ADMIN", "ADMIN", "USER_MANAGER") {
		return nil, errors.Forbidden("PERMISSION_DENIED", "无权限查看用户列表")
	}

	// 设置默认值
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Size <= 0 {
		req.Size = 10
	}

	// 获取用户列表
	users, total, err := s.userUc.ListUsers(ctx, req.Page, req.Size, req.Search)
	if err != nil {
		s.log.Errorf("Failed to list users: %v", err)
		return nil, errors.InternalServer("INTERNAL_ERROR", "用户列表获取失败")
	}

	// 转换为响应格式
	userInfos := make([]*UserInfo, len(users))
	for i, user := range users {
		// 获取角色和权限
		roles, _ := s.userUc.GetUserRoles(ctx, user.ID)
		permissions, _ := s.userUc.GetUserPermissions(ctx, user.ID)
		
		roleStrs := make([]string, len(roles))
		for j, role := range roles {
			roleStrs[j] = role.Code
		}

		userInfos[i] = ToUserInfo(user, roleStrs, permissions)
	}

	return &UserListResponse{
		Users: userInfos,
		Total: total,
		Page:  req.Page,
		Size:  req.Size,
	}, nil
}

// AssignRoles 分配用户角色
func (s *UserService) AssignRoles(ctx context.Context, req *AssignRolesRequest) error {
	// 检查权限
	currentUser := middleware.GetCurrentUser(ctx)
	if !currentUser.HasAnyRole("SUPER_ADMIN", "ADMIN", "USER_MANAGER") {
		return errors.Forbidden("PERMISSION_DENIED", "无权限分配角色")
	}

	s.log.Infof("Assigning roles to user: %d by %s", req.UserID, currentUser.Username)

	// 检查用户是否存在
	if _, err := s.userUc.GetUser(ctx, req.UserID); err != nil {
		return errors.NotFound("USER_NOT_FOUND", "用户不存在")
	}

	// 分配角色
	if err := s.userUc.AssignRoles(ctx, req.UserID, req.RoleIDs); err != nil {
		s.log.Errorf("Failed to assign roles: %v", err)
		return errors.InternalServer("INTERNAL_ERROR", "角色分配失败")
	}

	s.log.Infof("Roles assigned successfully to user: %d", req.UserID)
	return nil
}

// ResetPassword 重置用户密码
func (s *UserService) ResetPassword(ctx context.Context, req *ResetPasswordRequest) error {
	// 检查权限
	currentUser := middleware.GetCurrentUser(ctx)
	if !currentUser.HasAnyRole("SUPER_ADMIN", "ADMIN") {
		return errors.Forbidden("PERMISSION_DENIED", "无权限重置密码")
	}

	s.log.Infof("Resetting password for user: %d by %s", req.UserID, currentUser.Username)

	// 验证密码强度
	if err := s.pwdMgr.ValidatePasswordStrength(req.NewPassword); err != nil {
		return errors.BadRequest("WEAK_PASSWORD", err.Error())
	}

	// 验证用户存在
	_, err := s.userUc.GetUser(ctx, req.UserID)
	if err != nil {
		return errors.NotFound("USER_NOT_FOUND", "用户不存在")
	}

	// 加密新密码
	hashedPassword, err := s.userUc.HashPassword(req.NewPassword)
	if err != nil {
		s.log.Errorf("Failed to hash password: %v", err)
		return errors.InternalServer("INTERNAL_ERROR", "密码加密失败")
	}

	// 更新密码
	if err := s.userUc.UpdatePassword(ctx, req.UserID, hashedPassword); err != nil {
		s.log.Errorf("Failed to reset password: %v", err)
		return errors.InternalServer("INTERNAL_ERROR", "密码重置失败")
	}

	// TODO: 发送密码重置通知

	s.log.Infof("Password reset successfully for user: %d", req.UserID)
	return nil
}

// ToggleTwoFactor 切换用户2FA
func (s *UserService) ToggleTwoFactor(ctx context.Context, req *ToggleTwoFactorRequest) error {
	// 检查权限
	currentUser := middleware.GetCurrentUser(ctx)
	if !currentUser.HasAnyRole("SUPER_ADMIN", "ADMIN") && currentUser.ID != int64(req.UserID) {
		return errors.Forbidden("PERMISSION_DENIED", "无权限修改2FA设置")
	}

	s.log.Infof("Toggling 2FA for user: %d, enable: %v", req.UserID, req.Enable)

	if req.Enable {
		// 生成2FA密钥
		secret := s.pwdMgr.GenerateTOTPSecret()
		if err := s.userUc.EnableTwoFactor(ctx, req.UserID, secret); err != nil {
			s.log.Errorf("Failed to enable 2FA: %v", err)
			return errors.InternalServer("INTERNAL_ERROR", "2FA启用失败")
		}
	} else {
		if err := s.userUc.DisableTwoFactor(ctx, req.UserID); err != nil {
			s.log.Errorf("Failed to disable 2FA: %v", err)
			return errors.InternalServer("INTERNAL_ERROR", "2FA禁用失败")
		}
	}

	s.log.Infof("2FA toggled successfully for user: %d", req.UserID)
	return nil
}