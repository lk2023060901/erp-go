/**
 * 路由守卫组件
 */

import { ReactNode, useEffect, useState } from 'react';
import { Navigate, useLocation } from 'react-router-dom';
import { useAuth, usePermission } from '../hooks/useAuth';
import type { PermissionCheckOptions } from '../types/auth';

interface RouteGuardProps {
  children: ReactNode;
  // 权限配置
  permissions?: string[];
  roles?: string[];
  requireAll?: boolean;
  requireAuth?: boolean;
  skipSuperAdmin?: boolean;
  customCheck?: PermissionCheckOptions['customCheck'];
  // 重定向配置
  redirectTo?: string;
  loginPath?: string;
  unauthorizedPath?: string;
  // 加载配置
  loading?: ReactNode;
  // 调试配置
  debug?: boolean;
}

/**
 * 路由守卫组件
 * 
 * 用法示例：
 * ```tsx
 * // 需要登录的路由
 * <RouteGuard requireAuth>
 *   <Dashboard />
 * </RouteGuard>
 * 
 * // 需要特定权限的路由
 * <RouteGuard permissions={['user.manage']} redirectTo="/403">
 *   <UserManagement />
 * </RouteGuard>
 * 
 * // 需要特定角色的路由
 * <RouteGuard roles={['ADMIN']} loginPath="/login">
 *   <AdminPanel />
 * </RouteGuard>
 * ```
 */
export function RouteGuard({
  children,
  permissions,
  roles,
  requireAll = false,
  requireAuth = true,
  skipSuperAdmin = true,
  customCheck,
  redirectTo,
  loginPath = '/login',
  unauthorizedPath = '/403',
  loading = <div>Loading...</div>,
  debug = false,
}: RouteGuardProps) {
  const { isAuthenticated, loading: authLoading, user } = useAuth();
  const { checkPermission } = usePermission();
  const location = useLocation();
  const [checked, setChecked] = useState(false);

  useEffect(() => {
    // 延迟检查，等待认证状态确定
    const timer = setTimeout(() => setChecked(true), 100);
    return () => clearTimeout(timer);
  }, []);

  // 调试信息
  if (debug) {
    console.group(`RouteGuard Debug - ${location.pathname}`);
    console.log('Authenticated:', isAuthenticated);
    console.log('User:', user);
    console.log('Auth Loading:', authLoading);
    console.log('Checked:', checked);
    console.log('Required Auth:', requireAuth);
    console.log('Required Permissions:', permissions);
    console.log('Required Roles:', roles);
    console.groupEnd();
  }

  // 等待认证状态确定
  if (authLoading || !checked) {
    return <>{loading}</>;
  }

  // 需要认证但用户未登录
  if (requireAuth && !isAuthenticated) {
    return <Navigate to={loginPath} state={{ from: location }} replace />;
  }

  // 如果不需要认证，直接渲染
  if (!requireAuth && !permissions && !roles && !customCheck) {
    return <>{children}</>;
  }

  // 执行权限检查
  if (isAuthenticated && (permissions || roles || customCheck)) {
    const checkOptions: PermissionCheckOptions = {
      permissions,
      roles,
      requireAll,
      skipSuperAdmin,
      customCheck,
    };

    const result = checkPermission(checkOptions);

    if (!result.allowed) {
      const redirectPath = redirectTo || unauthorizedPath;
      return (
        <Navigate 
          to={redirectPath} 
          state={{ 
            from: location,
            reason: result.reason,
            requiredPermissions: permissions,
            requiredRoles: roles,
          }} 
          replace 
        />
      );
    }
  }

  return <>{children}</>;
}

/**
 * 认证路由守卫（简化版）
 */
interface AuthGuardProps {
  children: ReactNode;
  loginPath?: string;
  loading?: ReactNode;
}

export function AuthGuard({ 
  children, 
  loginPath = '/login',
  loading = <div>Loading...</div>
}: AuthGuardProps) {
  return (
    <RouteGuard
      requireAuth={true}
      loginPath={loginPath}
      loading={loading}
    >
      {children}
    </RouteGuard>
  );
}

/**
 * 游客路由守卫（未登录用户专用）
 */
interface GuestGuardProps {
  children: ReactNode;
  redirectTo?: string;
  loading?: ReactNode;
}

export function GuestGuard({ 
  children, 
  redirectTo = '/dashboard',
  loading = <div>Loading...</div>
}: GuestGuardProps) {
  const { isAuthenticated, loading: authLoading } = useAuth();

  if (authLoading) {
    return <>{loading}</>;
  }

  if (isAuthenticated) {
    return <Navigate to={redirectTo} replace />;
  }

  return <>{children}</>;
}

/**
 * 角色路由守卫
 */
interface RoleRouteGuardProps {
  children: ReactNode;
  roles: string | string[];
  fallbackPath?: string;
  loading?: ReactNode;
}

export function RoleRouteGuard({ 
  children, 
  roles,
  fallbackPath = '/403',
  loading = <div>Loading...</div>
}: RoleRouteGuardProps) {
  return (
    <RouteGuard
      roles={Array.isArray(roles) ? roles : [roles]}
      unauthorizedPath={fallbackPath}
      loading={loading}
    >
      {children}
    </RouteGuard>
  );
}

/**
 * 管理员路由守卫
 */
interface AdminRouteGuardProps {
  children: ReactNode;
  fallbackPath?: string;
  loading?: ReactNode;
}

export function AdminRouteGuard({ 
  children, 
  fallbackPath = '/403',
  loading = <div>Loading...</div>
}: AdminRouteGuardProps) {
  return (
    <RoleRouteGuard
      roles={['SUPER_ADMIN', 'ADMIN']}
      fallbackPath={fallbackPath}
      loading={loading}
    >
      {children}
    </RoleRouteGuard>
  );
}

/**
 * 条件路由组件
 */
interface ConditionalRouteProps {
  condition: boolean;
  children: ReactNode;
  fallback?: ReactNode;
  redirectTo?: string;
}

export function ConditionalRoute({ 
  condition, 
  children, 
  fallback, 
  redirectTo 
}: ConditionalRouteProps) {
  if (!condition) {
    if (redirectTo) {
      return <Navigate to={redirectTo} replace />;
    }
    return <>{fallback}</>;
  }

  return <>{children}</>;
}

/**
 * 延迟路由守卫（用于需要等待异步数据的路由）
 */
interface LazyRouteGuardProps {
  children: ReactNode;
  checkAsync: () => Promise<boolean>;
  loading?: ReactNode;
  fallbackPath?: string;
  cacheKey?: string;
  cacheDuration?: number; // 缓存时长（毫秒）
}

export function LazyRouteGuard({
  children,
  checkAsync,
  loading = <div>Loading...</div>,
  fallbackPath = '/403',
  cacheKey,
  cacheDuration = 5 * 60 * 1000, // 默认5分钟
}: LazyRouteGuardProps) {
  const [state, setState] = useState<{
    loading: boolean;
    allowed: boolean | null;
  }>({
    loading: true,
    allowed: null,
  });

  useEffect(() => {
    let mounted = true;

    const performCheck = async () => {
      // 检查缓存
      if (cacheKey) {
        const cached = localStorage.getItem(`route_check_${cacheKey}`);
        if (cached) {
          const { result, timestamp } = JSON.parse(cached);
          if (Date.now() - timestamp < cacheDuration) {
            setState({ loading: false, allowed: result });
            return;
          }
        }
      }

      try {
        const result = await checkAsync();
        
        if (mounted) {
          setState({ loading: false, allowed: result });
          
          // 缓存结果
          if (cacheKey) {
            localStorage.setItem(`route_check_${cacheKey}`, JSON.stringify({
              result,
              timestamp: Date.now(),
            }));
          }
        }
      } catch (error) {
        console.error('Route check failed:', error);
        if (mounted) {
          setState({ loading: false, allowed: false });
        }
      }
    };

    performCheck();

    return () => {
      mounted = false;
    };
  }, [checkAsync, cacheKey, cacheDuration]);

  if (state.loading) {
    return <>{loading}</>;
  }

  if (!state.allowed) {
    return <Navigate to={fallbackPath} replace />;
  }

  return <>{children}</>;
}

/**
 * 多条件路由守卫
 */
interface MultiConditionRouteGuardProps {
  children: ReactNode;
  conditions: Array<{
    check: () => boolean | Promise<boolean>;
    fallbackPath?: string;
    loading?: ReactNode;
  }>;
  operator?: 'AND' | 'OR'; // 条件运算符
  loading?: ReactNode;
}

export function MultiConditionRouteGuard({
  children,
  conditions,
  operator = 'AND',
  loading = <div>Loading...</div>,
}: MultiConditionRouteGuardProps) {
  const [state, setState] = useState<{
    loading: boolean;
    results: boolean[];
    currentCheck: number;
  }>({
    loading: true,
    results: [],
    currentCheck: 0,
  });

  useEffect(() => {
    let mounted = true;

    const performChecks = async () => {
      const results: boolean[] = [];

      for (let i = 0; i < conditions.length; i++) {
        if (!mounted) break;

        setState(prev => ({ ...prev, currentCheck: i }));

        try {
          const result = await Promise.resolve(conditions[i].check());
          results.push(result);

          // 如果是AND操作且某个条件失败，可以提前结束
          if (operator === 'AND' && !result) {
            break;
          }
          // 如果是OR操作且某个条件成功，可以提前结束
          if (operator === 'OR' && result) {
            break;
          }
        } catch (error) {
          console.error(`Condition ${i} check failed:`, error);
          results.push(false);
          if (operator === 'AND') {
            break;
          }
        }
      }

      if (mounted) {
        setState({
          loading: false,
          results,
          currentCheck: conditions.length,
        });
      }
    };

    performChecks();

    return () => {
      mounted = false;
    };
  }, [conditions, operator]);

  if (state.loading) {
    const currentCondition = conditions[state.currentCheck];
    return <>{currentCondition?.loading || loading}</>;
  }

  // 计算最终结果
  const finalResult = operator === 'AND' 
    ? state.results.every(Boolean)
    : state.results.some(Boolean);

  if (!finalResult) {
    // 找到第一个失败的条件的fallback路径
    const failedIndex = operator === 'AND' 
      ? state.results.findIndex(result => !result)
      : state.results.length - 1; // OR操作中，如果所有都失败，使用最后一个

    const failedCondition = conditions[failedIndex];
    const fallbackPath = failedCondition?.fallbackPath || '/403';
    
    return <Navigate to={fallbackPath} replace />;
  }

  return <>{children}</>;
}