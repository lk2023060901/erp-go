/**
 * 认证相关类型定义
 */

// 用户信息类型
export interface User {
  id: number;
  username: string;
  email: string;
  first_name: string;
  last_name: string;
  phone?: string;
  gender?: 'M' | 'F' | 'O';
  birth_date?: string;
  avatar_url?: string;
  is_enabled: boolean;
  phone_verified: boolean;
  two_factor_enabled: boolean;
  last_login_time?: string;
  last_login_ip?: string;
  created_at: string;
  updated_at: string;
}

// 角色信息类型
export interface Role {
  id: number;
  name: string;
  code: string;
  description?: string;
  is_system_role: boolean;
  is_enabled: boolean;
  sort_order: number;
  created_at: string;
  updated_at: string;
}

// 权限信息类型
export interface Permission {
  id: number;
  parent_id?: number;
  name: string;
  code: string;
  resource: string;
  action: string;
  module: string;
  description?: string;
  is_menu: boolean;
  is_button: boolean;
  is_api: boolean;
  menu_url?: string;
  menu_icon?: string;
  api_path?: string;
  api_method?: string;
  level: number;
  sort_order: number;
  is_enabled: boolean;
  created_at: string;
  updated_at: string;
}

// 组织信息类型
export interface Organization {
  id: number;
  parent_id?: number;
  name: string;
  code: string;
  type: string;
  description?: string;
  level: number;
  path: string;
  leader_id?: number;
  phone?: string;
  email?: string;
  address?: string;
  data_isolation: boolean;
  sort_order: number;
  is_enabled: boolean;
  created_at: string;
  updated_at: string;
}

// 登录请求类型
export interface LoginRequest {
  username: string;
  password: string;
  captcha?: string;
  two_factor_code?: string;
  remember_me?: boolean;
}

// 登录响应类型
export interface LoginResponse {
  access_token: string;
  refresh_token: string;
  expires_in: number;
  user: User;
}

// 用户资料响应类型
export interface ProfileResponse {
  id: number;
  username: string;
  email: string;
  first_name: string;
  last_name: string;
  phone: string;
  gender: string;
  avatar_url: string;
  is_active: boolean;
  two_factor_enabled: boolean;
  last_login_at: string;
  created_at: string;
  roles: string[];
  permissions: string[];
}

// 修改密码请求类型
export interface ChangePasswordRequest {
  old_password: string;
  new_password: string;
  confirm_password: string;
}

// JWT Token Claims 类型
export interface JWTClaims {
  user_id: number;
  username: string;
  email: string;
  roles: string[];
  session_id: string;
  token_type: 'access' | 'refresh';
  exp: number;
  iat: number;
  iss: string;
}

// 认证状态类型
export interface AuthState {
  isAuthenticated: boolean;
  user: User | null;
  roles: Role[];
  permissions: string[];
  accessToken: string | null;
  refreshToken: string | null;
  loading: boolean;
  error: string | null;
}

// 权限检查选项类型
export interface PermissionCheckOptions {
  // 需要的权限列表（任意一个满足即可）
  permissions?: string[];
  // 需要的角色列表（任意一个满足即可）
  roles?: string[];
  // 是否需要所有权限都满足
  requireAll?: boolean;
  // 超级管理员是否跳过检查
  skipSuperAdmin?: boolean;
  // 自定义检查函数
  customCheck?: (user: User, permissions: string[], roles: Role[]) => boolean;
}

// 菜单项类型
export interface MenuItem {
  id: number;
  parent_id?: number;
  name: string;
  code: string;
  url?: string;
  icon?: string;
  sort_order: number;
  children?: MenuItem[];
  meta?: {
    requiresAuth?: boolean;
    permissions?: string[];
    roles?: string[];
    keepAlive?: boolean;
    hidden?: boolean;
  };
}

// 路由权限配置类型
export interface RoutePermission {
  path: string;
  permissions?: string[];
  roles?: string[];
  requiresAuth?: boolean;
  fallbackPath?: string;
}

// 按钮权限配置类型
export interface ButtonPermission {
  key: string;
  permission?: string;
  roles?: string[];
  visible?: boolean;
  disabled?: boolean;
}

// 权限验证结果类型
export interface PermissionResult {
  allowed: boolean;
  reason?: string;
}

// 会话信息类型
export interface SessionInfo {
  id: string;
  user_id: number;
  device_type: string;
  ip_address: string;
  user_agent: string;
  location?: string;
  is_active: boolean;
  last_activity_at: string;
  expires_at: string;
  created_at: string;
}

// 在线用户类型
export interface OnlineUser {
  user_id: number;
  username: string;
  first_name: string;
  last_name: string;
  avatar_url?: string;
  ip_address: string;
  user_agent: string;
  device_type: string;
  location?: string;
  login_time: string;
  last_activity: string;
  session_count: number;
}

// 双重认证设置类型
export interface TwoFactorAuth {
  enabled: boolean;
  secret_key?: string;
  qr_code_url?: string;
  backup_codes?: string[];
}

// 权限树节点类型
export interface PermissionTreeNode {
  permission: Permission;
  children: PermissionTreeNode[];
  checked: boolean;
  expanded: boolean;
}

// 组织树节点类型
export interface OrganizationTreeNode {
  organization: Organization;
  children: OrganizationTreeNode[];
  member_count: number;
  expanded: boolean;
}

// API 响应基础类型
export interface ApiResponse<T = any> {
  code: number;
  message: string;
  data: T;
  timestamp: number;
}

// 分页请求类型
export interface PageRequest {
  page: number;
  page_size: number;
  keyword?: string;
  sort_by?: string;
  sort_order?: 'asc' | 'desc';
}

// 分页响应类型
export interface PageResponse<T = any> {
  items: T[];
  total: number;
  page: number;
  page_size: number;
}

// 错误类型
export interface ApiError {
  code: string;
  message: string;
  details?: Record<string, any>;
}

// 权限上下文类型
export interface PermissionContextType {
  user: User | null;
  roles: Role[];
  permissions: string[];
  hasPermission: (permission: string | string[], options?: PermissionCheckOptions) => boolean;
  hasRole: (role: string | string[]) => boolean;
  hasAnyPermission: (permissions: string[]) => boolean;
  hasAllPermissions: (permissions: string[]) => boolean;
  checkPermission: (options: PermissionCheckOptions) => PermissionResult;
  refreshPermissions: () => Promise<void>;
  isLoading: boolean;
}

// 认证上下文类型
export interface AuthContextType {
  // 状态
  isAuthenticated: boolean;
  user: User | null;
  roles: Role[];
  permissions: string[];
  loading: boolean;
  error: string | null;

  // 方法
  login: (credentials: LoginRequest) => Promise<void>;
  logout: () => Promise<void>;
  refreshToken: () => Promise<void>;
  updateProfile: (data: Partial<User>) => Promise<void>;
  changePassword: (data: ChangePasswordRequest) => Promise<void>;
  checkPermission: (options: PermissionCheckOptions) => PermissionResult;
  hasPermission: (permission: string | string[]) => boolean;
  hasRole: (role: string | string[]) => boolean;
  clearError: () => void;
}

// ========== 用户编辑表单相关类型 ==========

// 表单模式类型
export type FormMode = 'create' | 'edit' | 'profile';

// 字段权限类型
export interface FieldPermissions {
  readonly: boolean;
  required: boolean;
  visible: boolean;
}

// 用户编辑表单数据类型
export interface UserFormData {
  // 基本资料
  username: string; // 创建模式可编辑，编辑模式只读
  email: string; // 创建模式可编辑，编辑模式只读
  first_name: string;
  last_name: string;
  is_active: boolean;
  avatar?: File;
  
  // 详细信息
  gender?: 'M' | 'F' | 'O';
  birth_date?: string;
  phone?: string;
  address?: string;
  emergency_contact?: string;
  emergency_phone?: string;
  biography?: string;
  
  // 权限信息
  selected_permissions: string[];
  permission_notes?: string;
}

// 权限模块类型
export interface PermissionModule {
  id: string;
  name: string;
  description: string;
  icon: React.ReactNode;
  permissions: PermissionItem[];
  expanded?: boolean;
}

// 权限项类型
export interface PermissionItem {
  id: string;
  name: string;
  code: string;
  description?: string;
  checked?: boolean;
}

// 用户编辑表单属性类型
export interface UserEditFormProps {
  user?: User | null; // 编辑模式时传入用户数据
  mode: FormMode; // 使用模式
  onSubmit: (data: UserFormData) => Promise<void>;
  onCancel?: () => void;
  showActionButtons?: boolean; // 是否显示操作按钮
  initialTab?: string; // 初始选项卡
}

// 选项卡属性基础类型
export interface TabProps {
  mode: FormMode;
  initialData?: Partial<UserFormData>;
  onDataChange?: (data: Partial<UserFormData>) => void;
}

// 头像上传属性类型
export interface AvatarUploadProps {
  value?: string; // 当前头像URL
  onChange?: (file: File | null) => void;
  disabled?: boolean;
}

// 字符计数器属性类型
export interface CharacterCounterProps {
  current: number;
  max: number;
  showWarning?: boolean;
  warningThreshold?: number;
}