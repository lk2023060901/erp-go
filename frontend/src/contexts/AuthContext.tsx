/**
 * 认证上下文
 */

import React, { createContext, ReactNode } from 'react';
import { useAuthState } from '../hooks/useAuth';
import type { AuthContextType } from '../types/auth';

// 创建认证上下文
export const AuthContext = createContext<AuthContextType | undefined>(undefined);

// 认证提供者属性
interface AuthProviderProps {
  children: ReactNode;
}

/**
 * 认证提供者组件
 * 
 * 用法：
 * ```tsx
 * function App() {
 *   return (
 *     <AuthProvider>
 *       <Router>
 *         <Routes>
 *           // ... 路由配置
 *         </Routes>
 *       </Router>
 *     </AuthProvider>
 *   );
 * }
 * ```
 */
export function AuthProvider({ children }: AuthProviderProps) {
  const authState = useAuthState();

  return (
    <AuthContext.Provider value={authState}>
      {children}
    </AuthContext.Provider>
  );
}

/**
 * 权限上下文（从认证上下文派生）
 */
import type { PermissionContextType } from '../types/auth';

export const PermissionContext = createContext<PermissionContextType | undefined>(undefined);

interface PermissionProviderProps {
  children: ReactNode;
}

/**
 * 权限提供者组件
 * 
 * 注意：这个组件需要在 AuthProvider 内部使用
 */
export function PermissionProvider({ children }: PermissionProviderProps) {
  const authContext = React.useContext(AuthContext);
  
  if (!authContext) {
    throw new Error('PermissionProvider must be used within AuthProvider');
  }

  const {
    user,
    roles,
    permissions,
    hasPermission,
    hasRole,
    checkPermission,
    loading,
  } = authContext;

  // 权限相关方法
  const hasAnyPermission = React.useCallback((permissionList: string[]) => {
    return permissionList.some(permission => hasPermission(permission));
  }, [hasPermission]);

  const hasAllPermissions = React.useCallback((permissionList: string[]) => {
    return permissionList.every(permission => hasPermission(permission));
  }, [hasPermission]);

  const refreshPermissions = React.useCallback(async () => {
    // 这里可以调用刷新权限的方法
    // 具体实现取决于认证服务的设计
    console.warn('refreshPermissions not implemented');
  }, []);

  const permissionContextValue: PermissionContextType = {
    user,
    roles,
    permissions,
    hasPermission,
    hasRole,
    hasAnyPermission,
    hasAllPermissions,
    checkPermission,
    refreshPermissions,
    isLoading: loading,
  };

  return (
    <PermissionContext.Provider value={permissionContextValue}>
      {children}
    </PermissionContext.Provider>
  );
}

/**
 * 组合提供者组件
 */
interface CombinedProviderProps {
  children: ReactNode;
}

export function CombinedAuthProvider({ children }: CombinedProviderProps) {
  return (
    <AuthProvider>
      <PermissionProvider>
        {children}
      </PermissionProvider>
    </AuthProvider>
  );
}

/**
 * 认证状态调试组件
 */
interface AuthDebuggerProps {
  position?: 'top-left' | 'top-right' | 'bottom-left' | 'bottom-right';
  show?: boolean;
}

export function AuthDebugger({ 
  position = 'bottom-right', 
  show = import.meta.env.DEV 
}: AuthDebuggerProps) {
  const authContext = React.useContext(AuthContext);
  
  if (!authContext || !show) {
    return null;
  }

  const { isAuthenticated, user, roles, permissions, loading, error } = authContext;

  const positionStyles = {
    'top-left': { top: 10, left: 10 },
    'top-right': { top: 10, right: 10 },
    'bottom-left': { bottom: 10, left: 10 },
    'bottom-right': { bottom: 10, right: 10 },
  };

  return (
    <div
      style={{
        position: 'fixed',
        ...positionStyles[position],
        backgroundColor: 'rgba(0, 0, 0, 0.8)',
        color: 'white',
        padding: '10px',
        borderRadius: '4px',
        fontSize: '12px',
        fontFamily: 'monospace',
        zIndex: 9999,
        maxWidth: '300px',
        maxHeight: '400px',
        overflow: 'auto',
      }}
    >
      <div><strong>Auth Debug Info</strong></div>
      <div>Authenticated: {isAuthenticated ? '✅' : '❌'}</div>
      <div>Loading: {loading ? '⏳' : '✅'}</div>
      {error && <div style={{ color: 'red' }}>Error: {error}</div>}
      {user && (
        <>
          <div><strong>User:</strong></div>
          <div>ID: {user.id}</div>
          <div>Username: {user.username}</div>
          <div>Email: {user.email}</div>
        </>
      )}
      {roles.length > 0 && (
        <>
          <div><strong>Roles:</strong></div>
          {roles.map(role => (
            <div key={role.id}>- {role.name}</div>
          ))}
        </>
      )}
      {permissions.length > 0 && (
        <>
          <div><strong>Permissions ({permissions.length}):</strong></div>
          <div style={{ maxHeight: '150px', overflow: 'auto' }}>
            {permissions.map(permission => (
              <div key={permission} style={{ fontSize: '10px' }}>- {permission}</div>
            ))}
          </div>
        </>
      )}
    </div>
  );
}

/**
 * HOC: 为组件提供认证上下文
 */
export function withAuth<P extends object>(Component: React.ComponentType<P>) {
  const WrappedComponent = (props: P) => {
    const authContext = React.useContext(AuthContext);
    
    if (!authContext) {
      throw new Error(`${Component.displayName || Component.name} must be used within AuthProvider`);
    }

    return <Component {...props} />;
  };

  WrappedComponent.displayName = `withAuth(${Component.displayName || Component.name})`;
  
  return WrappedComponent;
}

/**
 * HOC: 为组件提供权限上下文
 */
export function withPermission<P extends object>(Component: React.ComponentType<P>) {
  const WrappedComponent = (props: P) => {
    const permissionContext = React.useContext(PermissionContext);
    
    if (!permissionContext) {
      throw new Error(`${Component.displayName || Component.name} must be used within PermissionProvider`);
    }

    return <Component {...props} />;
  };

  WrappedComponent.displayName = `withPermission(${Component.displayName || Component.name})`;
  
  return WrappedComponent;
}