/**
 * 认证服务
 */

import type { 
  LoginRequest, 
  LoginResponse, 
  ProfileResponse, 
  ChangePasswordRequest, 
  User
} from '../types/auth';

// API 基础配置
const API_BASE_URL = (import.meta.env.VITE_API_URL || 'http://localhost:58080') + '/api/v1';

// 调试输出（可在生产环境中移除）
if (import.meta.env.DEV) {
  console.log('🔍 API Configuration Debug:');
  console.log('VITE_API_URL from env:', import.meta.env.VITE_API_URL);
  console.log('Final API_BASE_URL:', API_BASE_URL);
}

// HTTP 客户端配置
class ApiClient {
  private baseURL: string;

  constructor(baseURL: string) {
    this.baseURL = baseURL;
  }

  // 获取认证头
  private getAuthHeaders(): HeadersInit {
    const token = localStorage.getItem('access_token');
    return {
      'Content-Type': 'application/json',
      ...(token && { 'Authorization': `Bearer ${token}` }),
    };
  }

  // 通用请求方法
  private async request<T>(
    endpoint: string, 
    options: RequestInit = {}
  ): Promise<T> {
    const url = `${this.baseURL}${endpoint}`;
    
    // 调试输出（可在生产环境中移除）
    if (import.meta.env.DEV) {
      console.log('🚀 API Request Debug:', {
        baseURL: this.baseURL,
        endpoint: endpoint,
        finalURL: url
      });
    }
    
    const response = await fetch(url, {
      ...options,
      headers: {
        ...this.getAuthHeaders(),
        ...options.headers,
      },
    });

    if (!response.ok) {
      const errorData = await response.json().catch(() => null);
      console.error('API Error Response:', JSON.stringify(errorData, null, 2));
      const error = new Error(
        errorData?.error?.message || errorData?.message || `HTTP ${response.status}: ${response.statusText}`
      );
      (error as any).code = errorData?.error?.code || errorData?.code || response.status;
      (error as any).details = errorData?.error?.details || errorData?.details;
      throw error;
    }

    const data = await response.json();
    
    if (!data.success) {
      const error = new Error(data.error?.message || data.message || 'API request failed');
      (error as any).code = data.error?.code || 'API_ERROR';
      throw error;
    }

    return data.data;
  }

  // GET 请求
  get<T>(endpoint: string, params?: Record<string, any>): Promise<T> {
    let url = endpoint;
    if (params) {
      // 过滤掉 undefined 和 null 值
      const filteredParams = Object.entries(params)
        .filter(([_, value]) => value !== undefined && value !== null)
        .reduce((acc, [key, value]) => ({ ...acc, [key]: value }), {});
      
      if (Object.keys(filteredParams).length > 0) {
        url = `${endpoint}?${new URLSearchParams(filteredParams).toString()}`;
      }
    }
    return this.request<T>(url, { method: 'GET' });
  }

  // POST 请求
  post<T>(endpoint: string, data?: any): Promise<T> {
    return this.request<T>(endpoint, {
      method: 'POST',
      body: data ? JSON.stringify(data) : undefined,
    });
  }

  // PUT 请求
  put<T>(endpoint: string, data?: any): Promise<T> {
    return this.request<T>(endpoint, {
      method: 'PUT',
      body: data ? JSON.stringify(data) : undefined,
    });
  }

  // DELETE 请求
  delete<T>(endpoint: string): Promise<T> {
    return this.request<T>(endpoint, { method: 'DELETE' });
  }

  // PATCH 请求
  patch<T>(endpoint: string, data?: any): Promise<T> {
    return this.request<T>(endpoint, {
      method: 'PATCH',
      body: data ? JSON.stringify(data) : undefined,
    });
  }
}

// 创建API客户端实例
const apiClient = new ApiClient(API_BASE_URL);

/**
 * 用户注册
 */
export async function register(data: {
  username: string;
  email: string;
  password: string;
  confirm_password: string;
  first_name: string;
  last_name: string;
  phone?: string;
  gender?: string;
}): Promise<{ user: User; message: string }> {
  try {
    const response = await apiClient.post<{ user: User; message: string }>('/auth/register', data);
    return response;
  } catch (error) {
    console.error('Register failed:', error);
    throw new Error('注册失败，请重试');
  }
}

/**
 * 用户登录
 */
export async function login(credentials: LoginRequest): Promise<LoginResponse> {
  try {
    const response = await apiClient.post<LoginResponse>('/auth/login', credentials);
    return response;
  } catch (error) {
    console.error('Login failed:', error);
    throw new Error('登录失败，请检查用户名和密码');
  }
}

/**
 * 用户登出
 */
export async function logout(refreshToken: string): Promise<void> {
  try {
    await apiClient.post('/auth/logout', { refresh_token: refreshToken });
  } catch (error) {
    console.warn('Logout API call failed:', error);
    // 登出失败不影响客户端清理状态
  }
}

/**
 * 刷新访问令牌
 */
export async function refreshToken(refreshToken: string): Promise<{ 
  access_token: string; 
  refresh_token: string; 
  expires_in: number; 
}> {
  try {
    const response = await apiClient.post('/auth/refresh', {
      refresh_token: refreshToken,
    });
    return response as { access_token: string; refresh_token: string; expires_in: number; };
  } catch (error) {
    console.error('Refresh token failed:', error);
    throw new Error('令牌刷新失败，请重新登录');
  }
}

/**
 * 获取当前用户资料
 */
export async function getProfile(): Promise<ProfileResponse> {
  try {
    const response = await apiClient.get<ProfileResponse>('/auth/profile');
    return response;
  } catch (error) {
    console.error('Get profile failed:', error);
    throw new Error('获取用户信息失败');
  }
}

/**
 * 更新用户资料
 */
export async function updateProfile(data: Partial<User>): Promise<User> {
  try {
    const response = await apiClient.put<{ user: User }>('/auth/profile', data);
    return response.user;
  } catch (error) {
    console.error('Update profile failed:', error);
    throw new Error('更新用户信息失败');
  }
}

/**
 * 修改密码
 */
export async function changePassword(data: ChangePasswordRequest): Promise<void> {
  try {
    await apiClient.put('/auth/password', data);
  } catch (error) {
    console.error('Change password failed:', error);
    throw new Error('修改密码失败');
  }
}

/**
 * 发送验证码
 */
export async function sendVerificationCode(email: string, type: string): Promise<void> {
  try {
    await apiClient.post('/auth/send-code', { email, type });
  } catch (error) {
    console.error('Send verification code failed:', error);
    throw new Error('发送验证码失败');
  }
}

/**
 * 重置密码
 */
export async function resetPassword(data: {
  email: string;
  verification_code: string;
  new_password: string;
}): Promise<void> {
  try {
    await apiClient.post('/auth/reset-password', data);
  } catch (error) {
    console.error('Reset password failed:', error);
    throw new Error('重置密码失败');
  }
}

/**
 * 启用双重认证
 */
export async function enableTwoFactor(code: string): Promise<{
  secret_key: string;
  qr_code_url: string;
  backup_codes: string[];
}> {
  try {
    const response = await apiClient.post('/auth/enable-2fa', { code });
    return response as { secret_key: string; qr_code_url: string; backup_codes: string[]; };
  } catch (error) {
    console.error('Enable 2FA failed:', error);
    throw new Error('启用双重认证失败');
  }
}

/**
 * 禁用双重认证
 */
export async function disableTwoFactor(password: string, code: string): Promise<void> {
  try {
    await apiClient.post('/auth/disable-2fa', { password, code });
  } catch (error) {
    console.error('Disable 2FA failed:', error);
    throw new Error('禁用双重认证失败');
  }
}

/**
 * 验证双重认证
 */
export async function verifyTwoFactor(code: string): Promise<{ valid: boolean }> {
  try {
    const response = await apiClient.post<{ valid: boolean }>('/auth/verify-2fa', { code });
    return response;
  } catch (error) {
    console.error('Verify 2FA failed:', error);
    throw new Error('双重认证验证失败');
  }
}

/**
 * 上传用户头像
 */
export async function uploadAvatar(file: File): Promise<{ avatar_url: string }> {
  try {
    const formData = new FormData();
    formData.append('avatar', file);

    const token = localStorage.getItem('access_token');
    const response = await fetch(`${API_BASE_URL}/auth/avatar`, {
      method: 'POST',
      headers: {
        ...(token && { 'Authorization': `Bearer ${token}` }),
      },
      body: formData,
    });

    if (!response.ok) {
      throw new Error(`HTTP ${response.status}: ${response.statusText}`);
    }

    const data = await response.json();
    return data.data;
  } catch (error) {
    console.error('Upload avatar failed:', error);
    throw new Error('上传头像失败');
  }
}

/**
 * 检查用户权限
 */
export async function checkPermission(permission: string): Promise<{ has_permission: boolean }> {
  try {
    const response = await apiClient.post<{ has_permission: boolean }>('/permissions/check', {
      permission_code: permission,
    });
    return response;
  } catch (error) {
    console.error('Check permission failed:', error);
    throw new Error('权限检查失败');
  }
}

/**
 * 批量检查权限
 */
export async function checkPermissions(permissions: string[]): Promise<Record<string, boolean>> {
  try {
    const response = await apiClient.post<{ permissions: Record<string, boolean> }>('/permissions/batch-check', {
      permission_codes: permissions,
    });
    return response.permissions;
  } catch (error) {
    console.error('Batch check permissions failed:', error);
    throw new Error('批量权限检查失败');
  }
}

/**
 * 获取用户菜单
 */
export async function getUserMenus(): Promise<any[]> {
  try {
    const response = await apiClient.get<{ menus: any[] }>('/permissions/menus');
    return response.menus;
  } catch (error) {
    console.error('Get user menus failed:', error);
    throw new Error('获取用户菜单失败');
  }
}

/**
 * HTTP 拦截器设置
 */
class HttpInterceptors {
  private static tokenRefreshPromise: Promise<string> | null = null;

  /**
   * 响应拦截器 - 处理Token过期
   */
  static setupResponseInterceptor() {
    // 保存原始fetch
    const originalFetch = window.fetch;

    window.fetch = async (input: RequestInfo | URL, init?: RequestInit): Promise<Response> => {
      const response = await originalFetch(input, init);

      // 如果返回401且不是登录或刷新Token接口
      if (response.status === 401) {
        const url = typeof input === 'string' ? input : (input instanceof URL ? input.href : (input as Request).url);
        const isAuthEndpoint = url.includes('/auth/login') || url.includes('/auth/refresh');

        if (!isAuthEndpoint) {
          // 尝试刷新Token
          try {
            await this.refreshTokenIfNeeded();
            // 重新发送原请求
            const newHeaders = {
              ...init?.headers,
              'Authorization': `Bearer ${localStorage.getItem('access_token')}`,
            };
            return originalFetch(input, { ...init, headers: newHeaders });
          } catch (error) {
            // 刷新失败，清除认证状态
            localStorage.removeItem('access_token');
            localStorage.removeItem('refresh_token');
            localStorage.removeItem('user_info');
            window.location.href = '/login';
          }
        }
      }

      return response;
    };
  }

  /**
   * 刷新Token（防止并发）
   */
  private static async refreshTokenIfNeeded(): Promise<string> {
    // 如果已经在刷新Token，等待结果
    if (this.tokenRefreshPromise) {
      return this.tokenRefreshPromise;
    }

    const refreshTokenValue = localStorage.getItem('refresh_token');
    if (!refreshTokenValue) {
      throw new Error('No refresh token');
    }

    this.tokenRefreshPromise = refreshToken(refreshTokenValue)
      .then(response => {
        localStorage.setItem('access_token', response.access_token);
        localStorage.setItem('refresh_token', response.refresh_token);
        return response.access_token;
      })
      .finally(() => {
        this.tokenRefreshPromise = null;
      });

    return this.tokenRefreshPromise;
  }
}

// 初始化拦截器
HttpInterceptors.setupResponseInterceptor();

/**
 * 工具函数：检查API是否可用
 */
export async function checkApiHealth(): Promise<boolean> {
  try {
    const response = await fetch(`${API_BASE_URL}/health`);
    return response.ok;
  } catch (error) {
    console.error('API health check failed:', error);
    return false;
  }
}

/**
 * 工具函数：获取API版本信息
 */
export async function getApiVersion(): Promise<{ version: string; build_time: string }> {
  try {
    const response = await apiClient.get<{ version: string; build_time: string }>('/system/info');
    return response;
  } catch (error) {
    console.error('Get API version failed:', error);
    throw new Error('获取API版本信息失败');
  }
}

// 导出API客户端实例（供其他服务使用）
export { apiClient };