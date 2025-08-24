/**
 * 认证相关 Hook
 */

import { useState, useEffect, useCallback, useContext } from 'react';
import type { 
  User, 
 
  LoginRequest, 
  ChangePasswordRequest,
  AuthState,
  PermissionCheckOptions,
  PermissionResult
} from '../types/auth';
import { checkPermission, TokenUtils, PermissionCache } from '../utils/permissions';
import { AuthContext } from '../contexts/AuthContext';
import * as authService from '../services/authService';

// 本地存储键名
const ACCESS_TOKEN_KEY = 'access_token';
const REFRESH_TOKEN_KEY = 'refresh_token';
const USER_INFO_KEY = 'user_info';

/**
 * 认证状态管理 Hook
 */
export function useAuth() {
  const context = useContext(AuthContext);
  
  if (!context) {
    throw new Error('useAuth must be used within an AuthProvider');
  }
  
  return context;
}

/**
 * 底层认证状态管理 Hook（用于 AuthProvider）
 */
export function useAuthState() {
  const [state, setState] = useState<AuthState>({
    isAuthenticated: false,
    user: null,
    roles: [],
    permissions: [],
    accessToken: null,
    refreshToken: null,
    loading: true,
    error: null,
  });

  // 从本地存储初始化状态
  const initializeAuth = useCallback(async () => {
    try {
      const accessToken = localStorage.getItem(ACCESS_TOKEN_KEY);
      const refreshToken = localStorage.getItem(REFRESH_TOKEN_KEY);
      const userInfoStr = localStorage.getItem(USER_INFO_KEY);

      if (accessToken && !TokenUtils.isTokenExpired(accessToken)) {
        // Token有效，获取用户信息
        const userInfo = userInfoStr ? JSON.parse(userInfoStr) : null;
        
        if (userInfo) {
          // 获取用户权限信息
          const profile = await authService.getProfile();
          
          // 构造用户对象
          const user: User = {
            id: profile.id,
            username: profile.username,
            email: profile.email,
            first_name: profile.first_name,
            last_name: profile.last_name,
            phone: profile.phone,
            gender: profile.gender as 'M' | 'F' | 'O',
            avatar_url: profile.avatar_url,
            is_enabled: profile.is_active,
            phone_verified: false,
            two_factor_enabled: profile.two_factor_enabled,
            last_login_time: profile.last_login_at,
            created_at: profile.created_at,
            updated_at: profile.created_at,
          };
          
          setState(prev => ({
            ...prev,
            isAuthenticated: true,
            user,
            roles: [], // 角色对象数组需要从其他API获取
            permissions: profile.permissions,
            accessToken,
            refreshToken,
            loading: false,
          }));
        } else {
          // 用户信息不存在，清除Token
          clearAuth();
        }
      } else if (refreshToken && !TokenUtils.isTokenExpired(refreshToken)) {
        // AccessToken过期，尝试刷新
        await refreshAccessToken();
      } else {
        // 所有Token都无效，清除状态
        clearAuth();
      }
    } catch (error) {
      console.error('Initialize auth failed:', error);
      clearAuth();
    }
  }, []);

  // 登录
  const login = useCallback(async (credentials: LoginRequest) => {
    try {
      setState(prev => ({ ...prev, loading: true, error: null }));

      const response = await authService.login(credentials);
      
      // 保存Token到本地存储
      localStorage.setItem(ACCESS_TOKEN_KEY, response.access_token);
      localStorage.setItem(REFRESH_TOKEN_KEY, response.refresh_token);
      localStorage.setItem(USER_INFO_KEY, JSON.stringify(response.user));

      // 获取完整的用户权限信息
      const profile = await authService.getProfile();

      // 构造用户对象
      const user: User = {
        id: profile.id,
        username: profile.username,
        email: profile.email,
        first_name: profile.first_name,
        last_name: profile.last_name,
        phone: profile.phone,
        gender: profile.gender as 'M' | 'F' | 'O',
        avatar_url: profile.avatar_url,
        is_enabled: profile.is_active,
        phone_verified: false,
        two_factor_enabled: profile.two_factor_enabled,
        last_login_time: profile.last_login_at,
        created_at: profile.created_at,
        updated_at: profile.created_at,
      };

      setState(prev => ({
        ...prev,
        isAuthenticated: true,
        user,
        roles: [], // 角色对象数组需要从其他API获取
        permissions: profile.permissions,
        accessToken: response.access_token,
        refreshToken: response.refresh_token,
        loading: false,
      }));

      // 清除权限缓存
      PermissionCache.clear();
    } catch (error: any) {
      setState(prev => ({
        ...prev,
        loading: false,
        error: error.message || 'Login failed',
      }));
      throw error;
    }
  }, []);

  // 登出
  const logout = useCallback(async () => {
    try {
      setState(prev => ({ ...prev, loading: true }));

      // 调用登出API
      if (state.refreshToken) {
        await authService.logout(state.refreshToken);
      }
    } catch (error) {
      console.warn('Logout API call failed:', error);
    } finally {
      clearAuth();
    }
  }, [state.refreshToken]);

  // 刷新Token
  const refreshAccessToken = useCallback(async (): Promise<void> => {
    try {
      const refreshToken = localStorage.getItem(REFRESH_TOKEN_KEY);
      if (!refreshToken) {
        throw new Error('No refresh token available');
      }

      const response = await authService.refreshToken(refreshToken);

      // 更新Token
      localStorage.setItem(ACCESS_TOKEN_KEY, response.access_token);
      localStorage.setItem(REFRESH_TOKEN_KEY, response.refresh_token);

      setState(prev => ({
        ...prev,
        accessToken: response.access_token,
        refreshToken: response.refresh_token,
        loading: false,
      }));

    } catch (error) {
      console.error('Refresh token failed:', error);
      clearAuth();
      throw error;
    }
  }, []);

  // 更新用户资料
  const updateProfile = useCallback(async (data: Partial<User>) => {
    try {
      setState(prev => ({ ...prev, loading: true, error: null }));

      const updatedUser = await authService.updateProfile(data);

      setState(prev => ({
        ...prev,
        user: updatedUser,
        loading: false,
      }));

      // 更新本地存储
      localStorage.setItem(USER_INFO_KEY, JSON.stringify(updatedUser));
    } catch (error: any) {
      setState(prev => ({
        ...prev,
        loading: false,
        error: error.message || 'Update profile failed',
      }));
      throw error;
    }
  }, []);

  // 修改密码
  const changePassword = useCallback(async (data: ChangePasswordRequest) => {
    try {
      setState(prev => ({ ...prev, loading: true, error: null }));

      await authService.changePassword(data);

      setState(prev => ({ ...prev, loading: false }));
    } catch (error: any) {
      setState(prev => ({
        ...prev,
        loading: false,
        error: error.message || 'Change password failed',
      }));
      throw error;
    }
  }, []);

  // 权限检查
  const checkUserPermission = useCallback((options: PermissionCheckOptions): PermissionResult => {
    return checkPermission(state.user, state.permissions, state.roles, options);
  }, [state.user, state.permissions, state.roles]);

  // 检查是否拥有权限
  const hasPermission = useCallback((permission: string | string[]): boolean => {
    return checkUserPermission({
      permissions: Array.isArray(permission) ? permission : [permission]
    }).allowed;
  }, [checkUserPermission]);

  // 检查是否拥有角色
  const hasRole = useCallback((role: string | string[]): boolean => {
    return checkUserPermission({
      roles: Array.isArray(role) ? role : [role]
    }).allowed;
  }, [checkUserPermission]);

  // 清除认证状态
  const clearAuth = useCallback(() => {
    localStorage.removeItem(ACCESS_TOKEN_KEY);
    localStorage.removeItem(REFRESH_TOKEN_KEY);
    localStorage.removeItem(USER_INFO_KEY);
    
    setState({
      isAuthenticated: false,
      user: null,
      roles: [],
      permissions: [],
      accessToken: null,
      refreshToken: null,
      loading: false,
      error: null,
    });

    // 清除权限缓存
    PermissionCache.clear();
  }, []);

  // 清除错误
  const clearError = useCallback(() => {
    setState(prev => ({ ...prev, error: null }));
  }, []);

  // 自动刷新Token
  useEffect(() => {
    if (!state.accessToken || !state.isAuthenticated) return;

    const checkTokenExpiration = () => {
      if (TokenUtils.isTokenExpiringSoon(state.accessToken!, 5)) {
        refreshAccessToken().catch(console.error);
      }
    };

    // 立即检查一次
    checkTokenExpiration();

    // 每分钟检查一次
    const interval = setInterval(checkTokenExpiration, 60 * 1000);

    return () => clearInterval(interval);
  }, [state.accessToken, state.isAuthenticated, refreshAccessToken]);

  // 初始化认证状态
  useEffect(() => {
    initializeAuth();
  }, [initializeAuth]);

  return {
    // 状态
    ...state,
    
    // 方法
    login,
    logout,
    refreshToken: refreshAccessToken,
    updateProfile,
    changePassword,
    checkPermission: checkUserPermission,
    hasPermission,
    hasRole,
    clearError,
  };
}

/**
 * 权限检查 Hook
 */
export function usePermission() {
  const { user, permissions, roles, checkPermission: check } = useAuth();

  const checkPermission = useCallback((options: PermissionCheckOptions) => {
    return check(options);
  }, [check]);

  const hasPermission = useCallback((permission: string | string[], requireAll = false) => {
    const permissions = Array.isArray(permission) ? permission : [permission];
    return checkPermission({ permissions, requireAll }).allowed;
  }, [checkPermission]);

  const hasRole = useCallback((role: string | string[]) => {
    const roles = Array.isArray(role) ? role : [role];
    return checkPermission({ roles }).allowed;
  }, [checkPermission]);

  const hasAnyPermission = useCallback((permissions: string[]) => {
    return checkPermission({ permissions, requireAll: false }).allowed;
  }, [checkPermission]);

  const hasAllPermissions = useCallback((permissions: string[]) => {
    return checkPermission({ permissions, requireAll: true }).allowed;
  }, [checkPermission]);

  return {
    user,
    permissions,
    roles,
    checkPermission,
    hasPermission,
    hasRole,
    hasAnyPermission,
    hasAllPermissions,
  };
}

/**
 * 在线状态检查 Hook
 */
export function useOnlineStatus() {
  const [isOnline, setIsOnline] = useState(navigator.onLine);

  useEffect(() => {
    const handleOnline = () => setIsOnline(true);
    const handleOffline = () => setIsOnline(false);

    window.addEventListener('online', handleOnline);
    window.addEventListener('offline', handleOffline);

    return () => {
      window.removeEventListener('online', handleOnline);
      window.removeEventListener('offline', handleOffline);
    };
  }, []);

  return isOnline;
}

/**
 * 自动登出 Hook
 */
export function useAutoLogout(timeoutMinutes: number = 30) {
  const { logout, isAuthenticated } = useAuth();
  const [timeLeft, setTimeLeft] = useState(timeoutMinutes * 60);

  useEffect(() => {
    if (!isAuthenticated) return;

    let timer: number;
    let warningTimer: number;

    const resetTimer = () => {
      clearTimeout(timer);
      clearTimeout(warningTimer);
      setTimeLeft(timeoutMinutes * 60);

      // 设置警告计时器（剩余5分钟时警告）
      warningTimer = setTimeout(() => {
        setTimeLeft(5 * 60);
      }, (timeoutMinutes - 5) * 60 * 1000);

      // 设置自动登出计时器
      timer = setTimeout(() => {
        logout();
      }, timeoutMinutes * 60 * 1000);
    };

    // 监听用户活动
    const events = ['mousedown', 'mousemove', 'keypress', 'scroll', 'touchstart'];
    const resetTimerHandler = () => resetTimer();

    events.forEach(event => {
      document.addEventListener(event, resetTimerHandler, true);
    });

    // 初始化计时器
    resetTimer();

    return () => {
      clearTimeout(timer);
      clearTimeout(warningTimer);
      events.forEach(event => {
        document.removeEventListener(event, resetTimerHandler, true);
      });
    };
  }, [isAuthenticated, logout, timeoutMinutes]);

  return {
    timeLeft,
    isWarning: timeLeft <= 5 * 60, // 剩余5分钟内为警告状态
  };
}