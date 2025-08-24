package biz

import (
	"context"
	"time"

	"github.com/go-kratos/kratos/v2/log"
)

// User 用户实体
type User struct {
	ID               int32     `json:"id"`
	Username         string    `json:"username"`
	Email            string    `json:"email"`
	Password         string    `json:"-"` // 不在JSON中暴露
	FirstName        string    `json:"first_name"`
	LastName         string    `json:"last_name"`
	Phone            string    `json:"phone"`
	Gender           string    `json:"gender"`
	BirthDate        time.Time `json:"birth_date"`
	AvatarURL        string    `json:"avatar_url"`
	IsActive         bool      `json:"is_active"`
	TwoFactorEnabled bool      `json:"two_factor_enabled"`
	TwoFactorSecret  string    `json:"-"` // 不在JSON中暴露
	LastLoginAt      time.Time `json:"last_login_at"`
	LastLoginIP      string    `json:"last_login_ip"`
	LoginCount       int32     `json:"login_count"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`

	// 关联数据
	Roles         []*Role         `json:"roles,omitempty"`
	Organizations []*Organization `json:"organizations,omitempty"`
	Permissions   []string        `json:"permissions,omitempty"`
}

// Role 角色实体
type Role struct {
	ID           int32     `json:"id"`
	Name         string    `json:"name"`
	Code         string    `json:"code"`
	Description  string    `json:"description"`
	IsSystemRole bool      `json:"is_system_role"`
	IsEnabled    bool      `json:"is_enabled"`
	SortOrder    int32     `json:"sort_order"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`

	Permissions []*Permission `json:"permissions,omitempty"`
}

// Permission 权限实体
type Permission struct {
	ID          int32     `json:"id"`
	ParentID    *int32    `json:"parent_id"`
	Name        string    `json:"name"`
	Code        string    `json:"code"`
	Resource    string    `json:"resource"`
	Action      string    `json:"action"`
	Module      string    `json:"module"`
	Description string    `json:"description"`
	IsMenu      bool      `json:"is_menu"`
	Path        string    `json:"path"`
	SortOrder   int32     `json:"sort_order"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	Children []*Permission `json:"children,omitempty"`
}

// Organization 组织实体
type Organization struct {
	ID          int32     `json:"id"`
	ParentID    *int32    `json:"parent_id"`
	Name        string    `json:"name"`
	Code        string    `json:"code"`
	Description string    `json:"description"`
	IsEnabled   bool      `json:"is_enabled"`
	SortOrder   int32     `json:"sort_order"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	Children []*Organization `json:"children,omitempty"`
	Users    []*User         `json:"users,omitempty"`
}

// UserSession 用户会话实体
type UserSession struct {
	ID           string    `json:"id"`
	UserID       int32     `json:"user_id"`
	DeviceType   string    `json:"device_type"`
	IPAddress    string    `json:"ip_address"`
	UserAgent    string    `json:"user_agent"`
	Location     string    `json:"location"`
	IsActive     bool      `json:"is_active"`
	LastActivity time.Time `json:"last_activity_at"`
	ExpiresAt    time.Time `json:"expires_at"`
	CreatedAt    time.Time `json:"created_at"`

	User *User `json:"user,omitempty"`
}

// OperationLog 操作日志实体
type OperationLog struct {
	ID            int32     `json:"id"`
	UserID        *int32    `json:"user_id"`
	Username      string    `json:"username"`
	Action        string    `json:"action"`
	Resource      string    `json:"resource"`
	ResourceID    string    `json:"resource_id"`
	Description   string    `json:"description"`
	IPAddress     string    `json:"ip_address"`
	UserAgent     string    `json:"user_agent"`
	RequestData   string    `json:"request_data"`
	ResponseData  string    `json:"response_data"`
	Status        string    `json:"status"`
	ErrorMessage  string    `json:"error_message"`
	ExecutionTime int32     `json:"execution_time"`
	CreatedAt     time.Time `json:"created_at"`
}

// LoginRequest 登录请求
type LoginRequest struct {
	Username      string `json:"username"`
	Password      string `json:"password"`
	TwoFactorCode string `json:"two_factor_code,omitempty"`
	RememberMe    bool   `json:"remember_me"`
}

// RegisterRequest 注册请求
type RegisterRequest struct {
	Username  string `json:"username"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Phone     string `json:"phone"`
	Gender    string `json:"gender"`
}

// LoginResponse 登录响应
type LoginResponse struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresIn    int64     `json:"expires_in"`
	TokenType    string    `json:"token_type"`
	User         *User     `json:"user"`
}

// UserRepo 用户仓储接口
type UserRepo interface {
	// 基础CRUD
	CreateUser(ctx context.Context, user *User) (*User, error)
	GetUser(ctx context.Context, id int32) (*User, error)
	GetUserByUsername(ctx context.Context, username string) (*User, error)
	GetUserByEmail(ctx context.Context, email string) (*User, error)
	UpdateUser(ctx context.Context, user *User) (*User, error)
	UpdatePassword(ctx context.Context, id int32, hashedPassword string) error
	DeleteUser(ctx context.Context, id int32) error
	ListUsers(ctx context.Context, page, size int32, search string) ([]*User, int32, error)
	ListUsersWithFilter(ctx context.Context, options interface{}) ([]*User, int32, error)

	// 认证相关
	ValidatePassword(hashedPassword, password string) bool
	HashPassword(password string) (string, error)
	UpdateLoginInfo(ctx context.Context, userID int32, ip string) error

	// 角色权限
	GetUserRoles(ctx context.Context, userID int32) ([]*Role, error)
	GetUserPermissions(ctx context.Context, userID int32) ([]string, error)
	AssignRoles(ctx context.Context, userID int32, roleIDs []int32) error

	// 2FA
	EnableTwoFactor(ctx context.Context, userID int32, secret string) error
	DisableTwoFactor(ctx context.Context, userID int32) error
	ValidateTwoFactor(ctx context.Context, userID int32, code string) bool
}

// RoleRepo 角色仓储接口
type RoleRepo interface {
	CreateRole(ctx context.Context, role *Role) (*Role, error)
	GetRole(ctx context.Context, id int32) (*Role, error)
	UpdateRole(ctx context.Context, role *Role) (*Role, error)
	DeleteRole(ctx context.Context, id int32) error
	ListRoles(ctx context.Context, page, size int32, search string, isEnabled *bool, sortField, sortOrder string) ([]*Role, int32, error)
	
	GetRolePermissions(ctx context.Context, roleID int32) ([]*Permission, error)
	AssignPermissions(ctx context.Context, roleID int32, permissionIDs []int32) error
	GetEnabledRoles(ctx context.Context) ([]*Role, error)
}


// OrganizationRepo 组织仓储接口
type OrganizationRepo interface {
	CreateOrganization(ctx context.Context, org *Organization) (*Organization, error)
	GetOrganization(ctx context.Context, id int32) (*Organization, error)
	UpdateOrganization(ctx context.Context, org *Organization) (*Organization, error)
	DeleteOrganization(ctx context.Context, id int32) error
	GetOrganizationTree(ctx context.Context) ([]*Organization, error)
	GetOrganizationUsers(ctx context.Context, orgID int32) ([]*User, error)
	AssignUsers(ctx context.Context, orgID int32, userIDs []int32) error
	GetEnabledOrganizations(ctx context.Context) ([]*Organization, error)
}

// SessionRepo 会话仓储接口
type SessionRepo interface {
	CreateSession(ctx context.Context, session *UserSession) (*UserSession, error)
	GetSession(ctx context.Context, sessionID string) (*UserSession, error)
	GetUserSession(ctx context.Context, userID int64, sessionID string) (*UserSession, error)
	UpdateSessionActivity(ctx context.Context, sessionID string) error
	DeactivateSession(ctx context.Context, sessionID string) error
	DeactivateUserSessions(ctx context.Context, userID int64) error
	ListUserSessions(ctx context.Context, userID int64) ([]*UserSession, error)
	CleanupExpiredSessions(ctx context.Context) error
}

// UserUsecase 用户用例接口
type UserUsecase struct {
	repo UserRepo
	log  *log.Helper
}

// 添加用户用例方法
func (uc *UserUsecase) CreateUser(ctx context.Context, user *User) (*User, error) {
	return uc.repo.CreateUser(ctx, user)
}

func (uc *UserUsecase) GetUser(ctx context.Context, id int32) (*User, error) {
	return uc.repo.GetUser(ctx, id)
}

func (uc *UserUsecase) GetUserByUsername(ctx context.Context, username string) (*User, error) {
	return uc.repo.GetUserByUsername(ctx, username)
}

func (uc *UserUsecase) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	return uc.repo.GetUserByEmail(ctx, email)
}

func (uc *UserUsecase) UpdateUser(ctx context.Context, user *User) (*User, error) {
	return uc.repo.UpdateUser(ctx, user)
}

func (uc *UserUsecase) UpdatePassword(ctx context.Context, id int32, hashedPassword string) error {
	return uc.repo.UpdatePassword(ctx, id, hashedPassword)
}

func (uc *UserUsecase) DeleteUser(ctx context.Context, id int32) error {
	return uc.repo.DeleteUser(ctx, id)
}

func (uc *UserUsecase) ListUsers(ctx context.Context, page, size int32, search string) ([]*User, int32, error) {
	return uc.repo.ListUsers(ctx, page, size, search)
}

func (uc *UserUsecase) ValidatePassword(hashedPassword, password string) bool {
	return uc.repo.ValidatePassword(hashedPassword, password)
}

func (uc *UserUsecase) HashPassword(password string) (string, error) {
	return uc.repo.HashPassword(password)
}

func (uc *UserUsecase) UpdateLoginInfo(ctx context.Context, userID int32, ip string) error {
	return uc.repo.UpdateLoginInfo(ctx, userID, ip)
}

func (uc *UserUsecase) GetUserRoles(ctx context.Context, userID int32) ([]*Role, error) {
	return uc.repo.GetUserRoles(ctx, userID)
}

func (uc *UserUsecase) GetUserPermissions(ctx context.Context, userID int32) ([]string, error) {
	return uc.repo.GetUserPermissions(ctx, userID)
}

func (uc *UserUsecase) AssignRoles(ctx context.Context, userID int32, roleIDs []int32) error {
	return uc.repo.AssignRoles(ctx, userID, roleIDs)
}

func (uc *UserUsecase) EnableTwoFactor(ctx context.Context, userID int32, secret string) error {
	return uc.repo.EnableTwoFactor(ctx, userID, secret)
}

func (uc *UserUsecase) DisableTwoFactor(ctx context.Context, userID int32) error {
	return uc.repo.DisableTwoFactor(ctx, userID)
}

func (uc *UserUsecase) ValidateTwoFactor(ctx context.Context, userID int32, code string) bool {
	return uc.repo.ValidateTwoFactor(ctx, userID, code)
}

// RoleUsecase 角色用例
type RoleUsecase struct {
	repo RoleRepo
	log  *log.Helper
}

// NewRoleUsecase 创建角色用例
func NewRoleUsecase(repo RoleRepo, logger log.Logger) *RoleUsecase {
	return &RoleUsecase{
		repo: repo,
		log:  log.NewHelper(logger),
	}
}

func (uc *RoleUsecase) CreateRole(ctx context.Context, role *Role) (*Role, error) {
	return uc.repo.CreateRole(ctx, role)
}

func (uc *RoleUsecase) GetRole(ctx context.Context, id int32) (*Role, error) {
	return uc.repo.GetRole(ctx, id)
}

func (uc *RoleUsecase) UpdateRole(ctx context.Context, role *Role) (*Role, error) {
	return uc.repo.UpdateRole(ctx, role)
}

func (uc *RoleUsecase) DeleteRole(ctx context.Context, id int32) error {
	return uc.repo.DeleteRole(ctx, id)
}

func (uc *RoleUsecase) ListRoles(ctx context.Context, page, size int32, search string, isEnabled *bool, sortField, sortOrder string) ([]*Role, int32, error) {
	return uc.repo.ListRoles(ctx, page, size, search, isEnabled, sortField, sortOrder)
}

func (uc *RoleUsecase) GetRolePermissions(ctx context.Context, roleID int32) ([]*Permission, error) {
	return uc.repo.GetRolePermissions(ctx, roleID)
}

func (uc *RoleUsecase) AssignPermissions(ctx context.Context, roleID int32, permissionIDs []int32) error {
	return uc.repo.AssignPermissions(ctx, roleID, permissionIDs)
}

func (uc *RoleUsecase) GetEnabledRoles(ctx context.Context) ([]*Role, error) {
	return uc.repo.GetEnabledRoles(ctx)
}



// OrganizationUsecase 组织用例
type OrganizationUsecase struct {
	repo OrganizationRepo
	log  *log.Helper
}

// NewOrganizationUsecase 创建组织用例
func NewOrganizationUsecase(repo OrganizationRepo, logger log.Logger) *OrganizationUsecase {
	return &OrganizationUsecase{
		repo: repo,
		log:  log.NewHelper(logger),
	}
}

func (uc *OrganizationUsecase) CreateOrganization(ctx context.Context, org *Organization) (*Organization, error) {
	return uc.repo.CreateOrganization(ctx, org)
}

func (uc *OrganizationUsecase) GetOrganization(ctx context.Context, id int32) (*Organization, error) {
	return uc.repo.GetOrganization(ctx, id)
}

func (uc *OrganizationUsecase) UpdateOrganization(ctx context.Context, org *Organization) (*Organization, error) {
	return uc.repo.UpdateOrganization(ctx, org)
}

func (uc *OrganizationUsecase) DeleteOrganization(ctx context.Context, id int32) error {
	return uc.repo.DeleteOrganization(ctx, id)
}

func (uc *OrganizationUsecase) GetOrganizationTree(ctx context.Context) ([]*Organization, error) {
	return uc.repo.GetOrganizationTree(ctx)
}

func (uc *OrganizationUsecase) GetOrganizationUsers(ctx context.Context, orgID int32) ([]*User, error) {
	return uc.repo.GetOrganizationUsers(ctx, orgID)
}

func (uc *OrganizationUsecase) AssignUsers(ctx context.Context, orgID int32, userIDs []int32) error {
	return uc.repo.AssignUsers(ctx, orgID, userIDs)
}

func (uc *OrganizationUsecase) GetEnabledOrganizations(ctx context.Context) ([]*Organization, error) {
	return uc.repo.GetEnabledOrganizations(ctx)
}

// 操作日志列表请求
type OperationLogListRequest struct {
	Page      int32      `json:"page"`
	Size      int32      `json:"size"`
	UserID    *int32     `json:"user_id"`
	Username  string     `json:"username"`
	Action    string     `json:"action"`
	Resource  string     `json:"resource"`
	Status    string     `json:"status"`
	StartTime time.Time  `json:"start_time"`
	EndTime   time.Time  `json:"end_time"`
}

// 操作统计信息
type OperationStatistics struct {
	TotalOperations      int32            `json:"total_operations"`
	SuccessOperations    int32            `json:"success_operations"`
	FailedOperations     int32            `json:"failed_operations"`
	AverageExecutionTime float64          `json:"average_execution_time"`
	ActiveUsers          int32            `json:"active_users"`
	ActionDistribution   map[string]int32 `json:"action_distribution"`
}

// 用户活动信息
type UserActivity struct {
	UserID         int32  `json:"user_id"`
	Username       string `json:"username"`
	Email          string `json:"email"`
	OperationCount int32  `json:"operation_count"`
}

// AuditRepo 审计仓储接口
type AuditRepo interface {
	CreateOperationLog(ctx context.Context, log *OperationLog) error
	GetOperationLog(ctx context.Context, id int32) (*OperationLog, error)
	ListOperationLogs(ctx context.Context, req *OperationLogListRequest) ([]*OperationLog, int32, error)
	DeleteOperationLogs(ctx context.Context, beforeTime time.Time) (int64, error)
	GetOperationStatistics(ctx context.Context, startTime, endTime time.Time) (*OperationStatistics, error)
	GetTopActiveUsers(ctx context.Context, startTime, endTime time.Time, limit int32) ([]*UserActivity, error)
}

// AuditUsecase 审计用例
type AuditUsecase struct {
	repo AuditRepo
	log  *log.Helper
}

// NewAuditUsecase 创建审计用例
func NewAuditUsecase(repo AuditRepo, logger log.Logger) *AuditUsecase {
	return &AuditUsecase{
		repo: repo,
		log:  log.NewHelper(logger),
	}
}

func (uc *AuditUsecase) CreateOperationLog(ctx context.Context, log *OperationLog) error {
	return uc.repo.CreateOperationLog(ctx, log)
}

func (uc *AuditUsecase) GetOperationLog(ctx context.Context, id int32) (*OperationLog, error) {
	return uc.repo.GetOperationLog(ctx, id)
}

func (uc *AuditUsecase) ListOperationLogs(ctx context.Context, req *OperationLogListRequest) ([]*OperationLog, int32, error) {
	return uc.repo.ListOperationLogs(ctx, req)
}

func (uc *AuditUsecase) DeleteOperationLogs(ctx context.Context, beforeTime time.Time) (int64, error) {
	return uc.repo.DeleteOperationLogs(ctx, beforeTime)
}

func (uc *AuditUsecase) GetOperationStatistics(ctx context.Context, startTime, endTime time.Time) (*OperationStatistics, error) {
	return uc.repo.GetOperationStatistics(ctx, startTime, endTime)
}

func (uc *AuditUsecase) GetTopActiveUsers(ctx context.Context, startTime, endTime time.Time, limit int32) ([]*UserActivity, error) {
	return uc.repo.GetTopActiveUsers(ctx, startTime, endTime, limit)
}

// NewUserUsecase 创建用户用例
func NewUserUsecase(repo UserRepo, logger log.Logger) *UserUsecase {
	return &UserUsecase{
		repo: repo,
		log:  log.NewHelper(logger),
	}
}

// Login 用户登录
func (uc *UserUsecase) Login(ctx context.Context, req *LoginRequest) (*LoginResponse, error) {
	// 获取用户
	user, err := uc.repo.GetUserByUsername(ctx, req.Username)
	if err != nil {
		uc.log.Errorf("failed to get user: %v", err)
		return nil, err
	}

	// 验证密码
	if !uc.repo.ValidatePassword(user.Password, req.Password) {
		uc.log.Warnf("invalid password for user: %s", req.Username)
		return nil, ErrInvalidCredentials
	}

	// 检查账户状态
	if !user.IsActive {
		return nil, ErrAccountDisabled
	}

	// 验证2FA（如果启用）
	if user.TwoFactorEnabled && req.TwoFactorCode != "" {
		if !uc.repo.ValidateTwoFactor(ctx, user.ID, req.TwoFactorCode) {
			return nil, ErrInvalidTwoFactor
		}
	}

	// 生成JWT令牌（这里需要实现JWT服务）
	// TODO: 实现JWT令牌生成
	
	return &LoginResponse{
		AccessToken:  "mock-access-token",
		RefreshToken: "mock-refresh-token",
		ExpiresIn:    3600,
		TokenType:    "Bearer",
		User:         user,
	}, nil
}

// Register 用户注册
func (uc *UserUsecase) Register(ctx context.Context, req *RegisterRequest) (*User, error) {
	// 检查用户名是否存在
	if _, err := uc.repo.GetUserByUsername(ctx, req.Username); err == nil {
		return nil, ErrUsernameExists
	}

	// 检查邮箱是否存在
	if _, err := uc.repo.GetUserByEmail(ctx, req.Email); err == nil {
		return nil, ErrEmailExists
	}

	// 密码加密
	hashedPassword, err := uc.repo.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	// 创建用户
	user := &User{
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

	return uc.repo.CreateUser(ctx, user)
}

// 定义错误
var (
	ErrInvalidCredentials = &BizError{Code: 401, Message: "Invalid username or password"}
	ErrAccountDisabled    = &BizError{Code: 403, Message: "Account is disabled"}
	ErrInvalidTwoFactor   = &BizError{Code: 401, Message: "Invalid two-factor authentication code"}
	ErrUsernameExists     = &BizError{Code: 400, Message: "Username already exists"}
	ErrEmailExists        = &BizError{Code: 400, Message: "Email already exists"}
	
	// 角色相关错误
	ErrRoleCodeExists        = &BizError{Code: 400, Message: "Role code already exists"}
	ErrRoleNameExists        = &BizError{Code: 400, Message: "Role name already exists"}
	ErrCannotDeleteSystemRole = &BizError{Code: 400, Message: "Cannot delete system role"}
	ErrRoleInUse             = &BizError{Code: 400, Message: "Role is in use"}
	
	// 权限相关错误
	ErrPermissionCodeExists  = &BizError{Code: 400, Message: "Permission code already exists"}
	ErrPermissionHasChildren = &BizError{Code: 400, Message: "Permission has child permissions"}
	ErrPermissionInUse       = &BizError{Code: 400, Message: "Permission is in use"}
	
	// 组织相关错误
	ErrOrganizationCodeExists  = &BizError{Code: 400, Message: "Organization code already exists"}
	ErrOrganizationHasChildren = &BizError{Code: 400, Message: "Organization has child organizations"}
	ErrOrganizationHasUsers    = &BizError{Code: 400, Message: "Organization has users"}
)

// BizError 业务错误
type BizError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (e *BizError) Error() string {
	return e.Message
}